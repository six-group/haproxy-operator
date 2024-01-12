package instance

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sort"
	"text/template"

	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"github.com/six-group/haproxy-operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

type initScriptData struct {
	Host string
	IP   string
	File string
}

func (r *Reconciler) reconcileStatefulSet(ctx context.Context, instance *proxyv1alpha1.Instance) error {
	logger := log.FromContext(ctx)

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-haproxy", instance.Name),
			Namespace: instance.Namespace,
		},
	}

	// FIXME OSCP-4269 workaround to change podManagementPolicy
	_ = r.Get(ctx, client.ObjectKeyFromObject(statefulset), statefulset)
	if statefulset.Spec.PodManagementPolicy == appsv1.OrderedReadyPodManagement {
		_ = r.Delete(ctx, statefulset)
		logger.Info("Delete stateful set to change podManagementPolicy")
	}

	// cannot avoid update triggered at startup
	// too many properties are added by the system (spec.template etc)
	result, err := controllerutil.CreateOrPatch(ctx, r.Client, statefulset, func() error {
		if err := controllerutil.SetOwnerReference(instance, statefulset, r.Scheme); err != nil {
			return err
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
					Containers: []corev1.Container{
						{
							Name:            "haproxy",
							Image:           utils.StringOrDefault(instance.Spec.Image, "haproxy:latest"),
							ImagePullPolicy: instance.Spec.ImagePullPolicy,
							Env: []corev1.EnvVar{
								{Name: "HAPROXY_SOCKET", Value: "/var/lib/haproxy/run/haproxy.sock"},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "haproxy-run",
									MountPath: "/var/lib/haproxy/run",
								},
								{
									Name:      "haproxy-config",
									MountPath: "/usr/local/etc/haproxy",
								},
							},
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
									SecretName: utils.GetConfigSecretName(instance),
								},
							},
						},
					},
				},
			},
		}

		if hasLocalLoggingTarget(instance) {
			volumes := []corev1.Volume{
				{
					Name: "rsyslog-config",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: utils.GetConfigSecretName(instance),
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
				MountPath: "/var/lib/rsyslog", // FIXME use filepath.Dir()
			}
			statefulset.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulset.Spec.Template.Spec.Containers[0].VolumeMounts, mount)

			container := corev1.Container{
				Name:            "logs",
				Image:           utils.GetRsyslogImage(),
				ImagePullPolicy: instance.Spec.ImagePullPolicy,
				Command:         []string{"/sbin/rsyslogd", "-n", "-i", "/tmp/rsyslog.pid", "-f", "/etc/rsyslog/rsyslog.conf"},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "rsyslog-run",
						MountPath: "/var/lib/rsyslog", // FIXME use filepath.Dir()
					},
					{
						Name:      "rsyslog-config",
						MountPath: "/etc/rsyslog",
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
				ImagePullPolicy: instance.Spec.ImagePullPolicy,
				Command:         []string{"/bin/sh", "-c"},
				Args:            []string{script},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "haproxy-run",
						MountPath: "/var/lib/haproxy/run",
					},
				},
			})

			statefulset.Spec.Template.Spec.Containers[0].Env = append(statefulset.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
				Name:  "ENV_FILE",
				Value: file,
			})
		}

		return nil
	})
	if err != nil {
		return err
	}
	if result != controllerutil.OperationResultNone {
		logger.Info(fmt.Sprintf("Object %s", result), "statefulset", statefulset.Name)
	}

	return nil
}

func hasLocalLoggingTarget(instance *proxyv1alpha1.Instance) bool {
	config := instance.Spec.Configuration.Global.Logging
	return config != nil && config.Enabled && net.ParseIP(config.Address) == nil
}
