package instance

import (
	"context"
	// #nosec
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	haproxy "github.com/haproxytech/client-native/v6/configuration/options"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"github.com/six-group/haproxy-operator/pkg/defaults"
	"github.com/six-group/haproxy-operator/pkg/utils"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *Reconciler) reconcileConfig(ctx context.Context, instance *proxyv1alpha1.Instance, listens *configv1alpha1.ListenList, frontends *configv1alpha1.FrontendList, backends *configv1alpha1.BackendList, resolvers *configv1alpha1.ResolverList) (string, error) {
	logger := log.FromContext(ctx)

	config, err := r.generateHAPProxyConfiguration(ctx, instance, listens, frontends, backends, resolvers)
	if err != nil {
		return "", err
	}

	certificates, err := r.generateCertificates(ctx, instance, listens, frontends, backends)
	if err != nil {
		return "", err
	}

	envs, err := r.generateEnvs(ctx, instance, listens)
	if err != nil {
		return "", err
	}

	mappings, err := r.generateBackendMappingFiles(ctx, instance, frontends)
	if err != nil {
		return "", err
	}

	errorFiles, err := r.generateErrorFiles(ctx, instance, frontends, backends)
	if err != nil {
		return "", err
	}

	customCerts, err := r.generateCustomCertificatesFile(ctx, instance, frontends, listens)
	if err != nil {
		return "", err
	}

	aclValueFiles := r.generateACLValuesFiles(ctx, listens, frontends, backends)

	configSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetConfigSecretName(instance),
			Namespace: instance.Namespace,
		},
	}
	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, configSecret, func() error {
		if err := controllerutil.SetOwnerReference(instance, configSecret, r.Scheme); err != nil {
			return err
		}

		configSecret.Data = map[string][]byte{
			filepath.Base(haproxy.DefaultConfigurationFile): []byte(config),
		}

		if hasLocalLoggingTarget(instance) {
			configSecret.Data["rsyslog.conf"] = []byte(fmt.Sprintf(utils.RsyslogConfigFormat, instance.Spec.Configuration.Global.Logging.Address))
		}

		for file, certificate := range certificates {
			configSecret.Data[filepath.Base(file)] = []byte(certificate)
		}

		if len(envs) > 0 {
			configSecret.Data["env"] = []byte(strings.Join(envs, "/n"))
		}

		for file, data := range mappings {
			configSecret.Data[filepath.Base(file)] = []byte(data)
		}

		for file, data := range errorFiles {
			configSecret.Data[filepath.Base(file)] = []byte(data)
		}

		for file, data := range customCerts {
			configSecret.Data[filepath.Base(file)] = []byte(data)
		}

		for file, data := range aclValueFiles {
			configSecret.Data[filepath.Base(file)] = []byte(data)
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	if result != controllerutil.OperationResultNone {
		logger.Info(fmt.Sprintf("Object %s", result), "secret", configSecret.Name)
	}

	cs := generateChecksum(configSecret)

	return cs, nil
}

// #nosec
func generateChecksum(secret *corev1.Secret) string {
	var b []byte
	keys := make([]string, 0, len(secret.Data))
	for k := range secret.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b = append(b, secret.Data[k]...)
	}

	hash := md5.Sum(b)
	return hex.EncodeToString(hash[:])
}

func (r *Reconciler) generateHAPProxyConfiguration(ctx context.Context, instance *proxyv1alpha1.Instance, listens *configv1alpha1.ListenList, frontends *configv1alpha1.FrontendList, backends *configv1alpha1.BackendList, resolvers *configv1alpha1.ResolverList) (string, error) {
	p, err := parser.New()
	if err != nil {
		return "", err
	}

	nameKindMap := make(map[string]string)

	if err := instance.AddToParser(p); err != nil {
		return "", err
	}

	for i := range listens.Items {
		listen := &listens.Items[i]

		if err = checkNameKind(nameKindMap, listen); err == nil {
			err = listen.AddToParser(p)
		}

		if err != nil {
			listen.Status.Phase = configv1alpha1.StatusPhaseInternalError
			listen.Status.Error = err.Error()
			return "", multierr.Combine(err, r.Status().Update(ctx, listen))
		}

	}

	for i := range frontends.Items {
		frontend := &frontends.Items[i]

		if err = checkNameKind(nameKindMap, frontend); err == nil {
			err = frontend.AddToParser(p)
		}

		if err != nil {
			frontend.Status.Phase = configv1alpha1.StatusPhaseInternalError
			frontend.Status.Error = err.Error()
			return "", multierr.Combine(err, r.Status().Update(ctx, frontend))
		}
	}

	for i := range backends.Items {
		backend := &backends.Items[i]

		if err = checkNameKind(nameKindMap, backend); err == nil {
			err = backend.AddToParser(p)
		}

		if err != nil {
			backend.Status.Phase = configv1alpha1.StatusPhaseInternalError
			backend.Status.Error = err.Error()
			return "", multierr.Combine(err, r.Status().Update(ctx, backend))
		}
	}

	for i := range resolvers.Items {
		resolver := &resolvers.Items[i]

		if err = checkNameKind(nameKindMap, resolver); err == nil {
			err = resolver.AddToParser(p)
		}

		if err != nil {
			resolver.Status.Phase = configv1alpha1.StatusPhaseInternalError
			resolver.Status.Error = err.Error()
			return "", multierr.Combine(err, r.Status().Update(ctx, resolver))
		}
	}

	if instance.Spec.Metrics != nil {
		if err := instance.Spec.Metrics.AddToParser(p); err != nil {
			return "", err
		}
	}

	return p.String(), nil
}

func (r *Reconciler) generateEnvs(ctx context.Context, instance *proxyv1alpha1.Instance, listens *configv1alpha1.ListenList) ([]string, error) {
	var envs []string

	for i := range listens.Items {
		listen := listens.Items[i]

		if listen.Spec.HTTPRequest != nil {
			for _, headers := range listen.Spec.HTTPRequest.SetHeader {
				envValues, err := r.headerEnvValue(ctx, instance, headers, listen)
				if err != nil {
					return nil, err
				}
				envs = append(envs, envValues...)
			}
			for _, headers := range listen.Spec.HTTPRequest.AddHeader {
				envValues, err := r.headerEnvValue(ctx, instance, headers, listen)
				if err != nil {
					return nil, err
				}
				envs = append(envs, envValues...)
			}
		}
	}

	return envs, nil
}

func (r *Reconciler) headerEnvValue(ctx context.Context, instance *proxyv1alpha1.Instance, headers configv1alpha1.HTTPHeaderRule, listen configv1alpha1.Listen) ([]string, error) {
	var envs []string

	if headers.Value.Env != nil {
		value := headers.Value.Env.Value

		if headers.Value.Env.ValueFrom != nil {
			ref := headers.Value.Env.ValueFrom.SecretKeyRef
			if ref != nil {
				secret := &corev1.Secret{}
				if err := r.Client.Get(ctx, client.ObjectKey{Name: ref.Name, Namespace: instance.Namespace}, secret); err != nil {
					listen.Status.Phase = configv1alpha1.StatusPhaseInternalError
					listen.Status.Error = err.Error()
					return nil, multierr.Combine(err, r.Status().Update(ctx, &listen))
				}

				bytes, ok := secret.Data[ref.Key]
				if !ok {
					err := fmt.Errorf("key %s not found in HTTP header secret: %s/%s", ref.Key, instance.Namespace, ref.Name)
					listen.Status.Phase = configv1alpha1.StatusPhaseInternalError
					listen.Status.Error = err.Error()
					return nil, multierr.Combine(err, r.Status().Update(ctx, &listen))
				}
				value = string(bytes)
			}
		}

		envs = append(envs, fmt.Sprintf("%s=%s", headers.Value.Env.Name, value))
	}
	return envs, nil
}

func (r *Reconciler) generateBackendMappingFiles(ctx context.Context, instance *proxyv1alpha1.Instance, frontends *configv1alpha1.FrontendList) (map[string]string, error) {
	files := map[string]string{}

	for i := range frontends.Items {
		frontend := frontends.Items[i]

		for _, rules := range frontend.Spec.BackendSwitching {
			if rules.Backend.RegexMapping != nil {
				labelSelector := rules.Backend.RegexMapping.LabelSelector
				selector, err := metav1.LabelSelectorAsSelector(&labelSelector)
				if err != nil {
					frontend.Status.Phase = configv1alpha1.StatusPhaseInternalError
					frontend.Status.Error = err.Error()
					return files, multierr.Combine(err, r.Status().Update(ctx, &frontend))
				}

				backends := &configv1alpha1.BackendList{}
				if err = r.Client.List(ctx, backends, client.MatchingLabelsSelector{Selector: selector}, client.InNamespace(instance.Namespace)); err != nil {
					frontend.Status.Phase = configv1alpha1.StatusPhaseInternalError
					frontend.Status.Error = err.Error()
					return files, multierr.Combine(err, r.Status().Update(ctx, &frontend))
				}

				var mappings []string
				for _, backend := range backends.Items {
					if backend.Spec.HostRegex == "" {
						err := fmt.Errorf("regex not found in backend: %s/%s", backend.Namespace, backend.Name)
						frontend.Status.Phase = configv1alpha1.StatusPhaseInternalError
						frontend.Status.Error = err.Error()
						return files, multierr.Combine(err, r.Status().Update(ctx, &frontend))
					}
					mappings = append(mappings, fmt.Sprintf("^%s$ %s", strings.TrimPrefix(strings.TrimSuffix(backend.Spec.HostRegex, "$"), "^"), backend.Name))
				}

				sort.Sort(sort.Reverse(sort.StringSlice(mappings)))
				files[rules.Backend.RegexMapping.FilePath()] = strings.Join(mappings, "\n")
			}
		}
	}

	return files, nil
}

func (r *Reconciler) generateErrorFiles(ctx context.Context, instance *proxyv1alpha1.Instance, frontends *configv1alpha1.FrontendList, backends *configv1alpha1.BackendList) (map[string]string, error) {
	files := map[string]string{}

	var list []configv1alpha1.StaticHTTPFile
	for _, frontend := range frontends.Items {
		for _, ef := range frontend.Spec.ErrorFiles {
			list = append(list, ef.File)
		}
	}
	for _, backend := range backends.Items {
		for _, ef := range backend.Spec.ErrorFiles {
			list = append(list, ef.File)
		}
	}

	for _, file := range list {
		if file.Value != nil {
			files[file.FilePath()] = *file.Value
			continue
		}
		if file.ValueFrom.ConfigMapKeyRef != nil {
			configmap := &corev1.ConfigMap{}
			if err := r.Client.Get(ctx, client.ObjectKey{Name: file.ValueFrom.ConfigMapKeyRef.Name, Namespace: instance.Namespace}, configmap); err != nil {
				return files, err
			}

			data, ok := configmap.Data[file.ValueFrom.ConfigMapKeyRef.Key]
			if !ok {
				return files, fmt.Errorf("key %s not found in HTTP static file configmap: %s/%s", file.ValueFrom.ConfigMapKeyRef, instance.Namespace, file.ValueFrom.ConfigMapKeyRef.Name)
			}

			files[file.FilePath()] = strings.TrimSpace(data)
		}
	}

	return files, nil
}

func (r *Reconciler) generateACLValuesFiles(_ context.Context, listens *configv1alpha1.ListenList, frontends *configv1alpha1.FrontendList, backends *configv1alpha1.BackendList) map[string]string {
	files := map[string]string{}

	for _, frontend := range frontends.Items {
		for _, acl := range frontend.Spec.ACL {
			if len(acl.Values) > defaults.MaxLineArgs-3 {
				sort.Strings(acl.Values)
				files[acl.FilePath()] = strings.Join(acl.Values, "\n")
			}
		}
	}

	for _, backend := range backends.Items {
		for _, acl := range backend.Spec.ACL {
			if len(acl.Values) > defaults.MaxLineArgs-3 {
				sort.Strings(acl.Values)
				files[acl.FilePath()] = strings.Join(acl.Values, "\n")
			}
		}
	}

	for _, listen := range listens.Items {
		for _, acl := range listen.Spec.ACL {
			if len(acl.Values) > defaults.MaxLineArgs-3 {
				sort.Strings(acl.Values)
				files[acl.FilePath()] = strings.Join(acl.Values, "\n")
			}
		}
	}

	return files
}

func checkNameKind(nameKindMap map[string]string, object client.Object) error {
	if val, ok := nameKindMap[object.GetName()]; ok {
		return fmt.Errorf("name %s already used by resource of kind %s", object.GetName(), val)
	}
	nameKindMap[object.GetName()] = object.GetObjectKind().GroupVersionKind().Kind
	return nil
}
