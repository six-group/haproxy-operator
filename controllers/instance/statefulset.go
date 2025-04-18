package instance

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"path/filepath"
	"sort"
	"text/template"

	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"github.com/six-group/haproxy-operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const initContainerScript = `
if [ "$HOSTNAME" = "{{.Host}}" ]
then
  i=0
  while [ $(ip a show to '{{.IP}}' | wc -l) -eq 0 ]
  do
    ((i=i+1))
    if [ "$i" -gt "20" ]
      then echo 'timeout waiting for IP {{.IP}}, aborting'
      exit 1
    fi
    echo 'waiting for IP {{.IP}} to be assigned...'
    sleep 5
  done

  echo 'IP {{.IP}} assignment verified, waiting 5 seconds before continuing...'

  sleep 5

  echo -n "BIND_ADDRESS={{.IP}}" > {{.File}}
  cat {{.File}}
  exit 0
fi

`

const (
	MemoryRequest = "256Mi"
	CPURequest    = "100m"
	MemoryLimit   = "512Mi"
	CPULimit      = "2000m"
)

type initScriptData struct {
	Host string
	IP   string
	File string
}

func (r *Reconciler) reconcileStatefulSet(ctx context.Context, instance *proxyv1alpha1.Instance, checksum string) error {
	logger := log.FromContext(ctx)

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-haproxy", instance.Name),
			Namespace: instance.Namespace,
		},
	}

	var create bool
	err := r.Get(ctx, client.ObjectKeyFromObject(statefulset), statefulset)
	if err != nil {
		if errors.IsNotFound(err) {
			create = true
		} else {
			return err
		}
	}

	oldObj := statefulset.DeepCopy()

	_ = r.Get(ctx, client.ObjectKeyFromObject(statefulset), statefulset)
	if statefulset.Spec.PodManagementPolicy == appsv1.OrderedReadyPodManagement {
		_ = r.Delete(ctx, statefulset)
		logger.Info("Delete stateful set to change podManagementPolicy")
	}

	if err := controllerutil.SetOwnerReference(instance, statefulset, r.Scheme); err != nil {
		return err
	}

	envVars := []corev1.EnvVar{{Name: "HAPROXY_SOCKET", Value: "/var/lib/haproxy/run/haproxy.sock"}}
	for k, v := range instance.Spec.Env {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}

	imagePullPolicy := corev1.PullIfNotPresent
	if instance.Spec.ImagePullPolicy != "" {
		imagePullPolicy = instance.Spec.ImagePullPolicy
	}

	statefulset.Spec = appsv1.StatefulSetSpec{
		Replicas: &instance.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: utils.GetAppSelectorLabels(instance),
		},
		PodManagementPolicy: appsv1.ParallelPodManagement,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      utils.GetPodLabels(instance),
				Annotations: map[string]string{},
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: instance.Spec.ServiceAccountName,
				ImagePullSecrets:   instance.Spec.ImagePullSecrets,
				Containers: []corev1.Container{
					{
						Name:            "haproxy",
						Image:           utils.StringOrDefault(instance.Spec.Image, "haproxy:latest"),
						ImagePullPolicy: imagePullPolicy,
						Env:             envVars,
						Resources:       getResources(instance),
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "haproxy-run",
								MountPath: filepath.Dir("/var/lib/haproxy/run/"),
							},
							{
								Name:      "haproxy-config",
								MountPath: filepath.Dir("/usr/local/etc/haproxy/"),
							},
						},
						ReadinessProbe: instance.Spec.ReadinessProbe,
						LivenessProbe:  instance.Spec.LivenessProbe,
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "haproxy-run",
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{},
						},
					},
					{
						Name: "haproxy-config",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								DefaultMode: ptr.To(int32(420)),
								SecretName:  utils.GetConfigSecretName(instance),
							},
						},
					},
				},
			},
		},
	}

	if instance.Spec.RolloutOnConfigChange {
		statefulset.Spec.Template.ObjectMeta.Annotations["checksum/config"] = checksum
	}

	if hasLocalLoggingTarget(instance) {
		volumes := []corev1.Volume{
			{
				Name: "rsyslog-config",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						DefaultMode: ptr.To(int32(420)),
						SecretName:  utils.GetConfigSecretName(instance),
						Items: []corev1.KeyToPath{
							{
								Key:  "rsyslog.conf",
								Path: "rsyslog.conf",
							},
						},
					},
				},
			},
			{
				Name: "rsyslog-run",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		}
		statefulset.Spec.Template.Spec.Volumes = append(statefulset.Spec.Template.Spec.Volumes, volumes...)

		mount := corev1.VolumeMount{
			Name:      "rsyslog-run",
			MountPath: filepath.Dir("/var/lib/rsyslog/"),
		}
		statefulset.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulset.Spec.Template.Spec.Containers[0].VolumeMounts, mount)

		container := corev1.Container{
			Name:            "logs",
			Image:           utils.GetRsyslogImage(),
			ImagePullPolicy: imagePullPolicy,
			Command:         []string{"/sbin/rsyslogd", "-n", "-i", "/tmp/rsyslog.pid", "-f", "/etc/rsyslog/rsyslog.conf"},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "rsyslog-run",
					MountPath: filepath.Dir("/var/lib/rsyslog/"),
				},
				{
					Name:      "rsyslog-config",
					MountPath: filepath.Dir("/etc/rsyslog/"),
				},
			},
		}
		statefulset.Spec.Template.Spec.Containers = append(statefulset.Spec.Template.Spec.Containers, container)
		statefulset.Spec.Template.Spec.Containers = append(statefulset.Spec.Template.Spec.Containers, instance.Spec.Sidecars...)
	}

	if instance.Spec.Network.HostNetwork {
		statefulset.Spec.Template.Spec.HostNetwork = true
		statefulset.Spec.Template.Spec.DNSPolicy = corev1.DNSClusterFirstWithHostNet
	}

	if instance.Spec.Placement != nil {
		statefulset.Spec.Template.Spec.NodeSelector = instance.Spec.Placement.NodeSelector
		statefulset.Spec.Template.Spec.TopologySpreadConstraints = instance.Spec.Placement.TopologySpreadConstraints

		for idx := range statefulset.Spec.Template.Spec.TopologySpreadConstraints {
			statefulset.Spec.Template.Spec.TopologySpreadConstraints[idx].LabelSelector = statefulset.Spec.Selector
		}
	}

	if ptr.Deref(instance.Spec.AllowPrivilegedPorts, false) {
		statefulset.Spec.Template.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{
			Privileged: ptr.To(true),
		}
	}

	if len(instance.Spec.Network.HostIPs) > 0 {
		file := "/var/lib/haproxy/run/env"

		var hosts []string
		for host := range instance.Spec.Network.HostIPs {
			hosts = append(hosts, host)
		}
		sort.Strings(hosts)

		script := ""
		for _, host := range hosts {
			data := initScriptData{
				Host: host,
				IP:   instance.Spec.Network.HostIPs[host],
				File: file,
			}

			tmpl, err := template.New("initScript").Parse(initContainerScript)
			if err != nil {
				return err
			}
			var s bytes.Buffer
			err = tmpl.Execute(&s, data)
			if err != nil {
				return err
			}

			script += s.String()
		}
		script += "exit 1\n"

		statefulset.Spec.Template.Spec.InitContainers = append(statefulset.Spec.Template.Spec.InitContainers, corev1.Container{
			Name:            "setup-env",
			Image:           utils.GetHelperImage(),
			ImagePullPolicy: imagePullPolicy,
			Command:         []string{"/bin/sh", "-c"},
			Args:            []string{script},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "haproxy-run",
					MountPath: filepath.Dir("/var/lib/haproxy/run/"),
				},
			},
		})

		statefulset.Spec.Template.Spec.Containers[0].Env = append(statefulset.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "ENV_FILE",
			Value: file,
		})
	}

	if needsUpdate(oldObj, statefulset) {
		if create {
			err = r.Create(ctx, statefulset)
			logger.Info("created", "statefulset", statefulset.Name)
			return err
		}
		err = r.Update(ctx, statefulset)
		logger.Info("updated", "statefulset", statefulset.Name)
		return err
	}

	return nil
}

func needsUpdate(old, new *appsv1.StatefulSet) bool {
	oldCpy := old.DeepCopy()
	newCpy := new.DeepCopy()

	removeIrrelevantProperties(oldCpy)
	removeIrrelevantProperties(newCpy)

	return !equality.Semantic.DeepEqual(oldCpy.Spec, newCpy.Spec) || !equality.Semantic.DeepEqual(oldCpy.OwnerReferences, newCpy.OwnerReferences)
}

func removeIrrelevantProperties(ss *appsv1.StatefulSet) {
	ss.Spec.UpdateStrategy = appsv1.StatefulSetUpdateStrategy{}
	ss.Spec.RevisionHistoryLimit = nil
	ss.Spec.PersistentVolumeClaimRetentionPolicy = nil

	ss.Spec.Template.Spec.SchedulerName = ""
	ss.Spec.Template.Spec.DeprecatedServiceAccount = ""
	ss.Spec.Template.Spec.RestartPolicy = ""
	ss.Spec.Template.Spec.TerminationGracePeriodSeconds = nil
	ss.Spec.Template.Spec.SecurityContext = nil
	ss.Spec.Template.Spec.SecurityContext = nil

	for i := range ss.Spec.Template.Spec.InitContainers {
		ss.Spec.Template.Spec.InitContainers[i].TerminationMessagePath = ""
		ss.Spec.Template.Spec.InitContainers[i].TerminationMessagePolicy = ""
	}
	for i := range ss.Spec.Template.Spec.Containers {
		ss.Spec.Template.Spec.Containers[i].TerminationMessagePath = ""
		ss.Spec.Template.Spec.Containers[i].TerminationMessagePolicy = ""
	}
}

func getResources(instance *proxyv1alpha1.Instance) corev1.ResourceRequirements {
	resources := corev1.ResourceRequirements{}
	if instance.Spec.Resources == nil {
		resources.Requests = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(MemoryRequest),
			corev1.ResourceCPU:    resource.MustParse(CPURequest),
		}
		resources.Limits = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(MemoryLimit),
			corev1.ResourceCPU:    resource.MustParse(CPULimit),
		}
	} else {
		resources = *instance.Spec.Resources
	}
	return resources
}

func hasLocalLoggingTarget(instance *proxyv1alpha1.Instance) bool {
	config := instance.Spec.Configuration.Global.Logging
	return config != nil && config.Enabled && net.ParseIP(config.Address) == nil
}
