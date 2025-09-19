package instance

import (
	"context"
	"fmt"
	"sort"
	"strings"

	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Reconciler) generateCertificates(ctx context.Context, instance *proxyv1alpha1.Instance, listens *configv1alpha1.ListenList, frontends *configv1alpha1.FrontendList, backends *configv1alpha1.BackendList) (map[string]string, error) {
	certificates := map[string]string{}

	for idx := range instance.Spec.Configuration.Global.AdditionalCertificates {
		certificate := instance.Spec.Configuration.Global.AdditionalCertificates[idx]

		data, err := r.loadSSLCertificateValueData(ctx, instance, &certificate)
		if err != nil {
			instance.Status.Phase = proxyv1alpha1.InstancePhaseInternalError
			instance.Status.Error = err.Error()
			return certificates, multierr.Combine(err, r.Status().Update(ctx, instance))
		}

		certificates[certificate.FilePath()] = data
	}

	for i := range listens.Items {
		listen := listens.Items[i]

		for _, certificate := range extractSLCCertificatesFromFrontend(listen.ToFrontend()) {
			data, err := r.loadSSLCertificateValueData(ctx, instance, certificate)
			if err != nil {
				listen.Status.Phase = configv1alpha1.StatusPhaseInternalError
				listen.Status.Error = err.Error()
				return certificates, multierr.Combine(err, r.Status().Update(ctx, &listen))
			}

			certificates[certificate.FilePath()] = data
		}

		for _, certificate := range extractSLCCertificatesFromBackend(listen.ToBackend()) {
			data, err := r.loadSSLCertificateValueData(ctx, instance, certificate)
			if err != nil {
				listen.Status.Phase = configv1alpha1.StatusPhaseInternalError
				listen.Status.Error = err.Error()
				return certificates, multierr.Combine(err, r.Status().Update(ctx, &listen))
			}

			certificates[certificate.FilePath()] = data
		}
	}

	for i := range frontends.Items {
		frontend := frontends.Items[i]

		for _, certificate := range extractSLCCertificatesFromFrontend(&frontend) {
			data, err := r.loadSSLCertificateValueData(ctx, instance, certificate)
			if err != nil {
				frontend.Status.Phase = configv1alpha1.StatusPhaseInternalError
				frontend.Status.Error = err.Error()
				return certificates, multierr.Combine(err, r.Status().Update(ctx, &frontend))
			}

			certificates[certificate.FilePath()] = data
		}
	}

	for i := range backends.Items {
		backend := backends.Items[i]

		for _, certificate := range extractSLCCertificatesFromBackend(&backend) {
			data, err := r.loadSSLCertificateValueData(ctx, instance, certificate)
			if err != nil {
				backend.Status.Phase = configv1alpha1.StatusPhaseInternalError
				backend.Status.Error = err.Error()
				return certificates, multierr.Combine(err, r.Status().Update(ctx, &backend))
			}

			certificates[certificate.FilePath()] = data
		}
	}

	return certificates, nil
}

func (r *Reconciler) generateCustomCertificatesFile(ctx context.Context, instance *proxyv1alpha1.Instance, frontends *configv1alpha1.FrontendList, listens *configv1alpha1.ListenList) (map[string]string, error) {
	files := map[string]string{}
	var mappings []string

	for i := range frontends.Items {
		frontend := frontends.Items[i]

		for _, bind := range frontend.Spec.Binds {
			if bind.SSLCertificateList != nil {
				var elements []configv1alpha1.CertificateListElement
				if len(bind.SSLCertificateList.Elements) > 0 {
					elements = append(elements, bind.SSLCertificateList.Elements...)
				}

				if bind.SSLCertificateList.LabelSelector != nil {
					selector, err := metav1.LabelSelectorAsSelector(bind.SSLCertificateList.LabelSelector)
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

					for _, backend := range backends.Items {
						if backend.Spec.HostCertificate != nil {
							elements = append(elements, *backend.Spec.HostCertificate)
						}
					}
				}

				for _, element := range elements {
					data, err := r.loadSSLCertificateValueData(ctx, instance, &element.Certificate)
					if err != nil {
						frontend.Status.Phase = configv1alpha1.StatusPhaseInternalError
						frontend.Status.Error = err.Error()
						return nil, multierr.Combine(err, r.Status().Update(ctx, &frontend))
					}
					files[element.Certificate.FilePath()] = data

					var alpn string
					if len(element.Alpn) > 0 {
						if element.Ocsp {
							if element.OcspFile != nil {
								files[element.OcspFile.FilePath()] = *element.OcspFile.Value
								alpn = fmt.Sprintf("[alpn %s %s %s %s]", strings.Join(element.Alpn, ","), "ocsp-update on", "ocsp", element.OcspFile.FilePath())
							} else {
								alpn = fmt.Sprintf("[alpn %s %s]", strings.Join(element.Alpn, ","), "ocsp-update on")
							}
						} else {
							alpn = fmt.Sprintf("[alpn %s]", strings.Join(element.Alpn, ","))
						}
					}

					mappings = append(mappings, strings.Join([]string{element.Certificate.FilePath(), alpn, element.SNIFilter, "\n"}, " "))
				}

				sort.Strings(mappings)
				files[bind.SSLCertificateList.FilePath()] = strings.Join(mappings, "")
			}
		}
	}

	for i := range listens.Items {
		listen := listens.Items[i]

		for _, bind := range listen.Spec.Binds {
			if bind.SSLCertificateList != nil {
				var elements []configv1alpha1.CertificateListElement
				if len(bind.SSLCertificateList.Elements) > 0 {
					elements = append(elements, bind.SSLCertificateList.Elements...)
				}

				if listen.Spec.HostCertificate != nil {
					elements = append(elements, *listen.Spec.HostCertificate)
				}

				for _, element := range elements {
					data, err := r.loadSSLCertificateValueData(ctx, instance, &element.Certificate)
					if err != nil {
						listen.Status.Phase = configv1alpha1.StatusPhaseInternalError
						listen.Status.Error = err.Error()
						return nil, multierr.Combine(err, r.Status().Update(ctx, &listen))
					}
					files[element.Certificate.FilePath()] = data

					var alpn string
					if len(element.Alpn) > 0 {
						if element.Ocsp {
							if element.OcspFile != nil {
								files[element.OcspFile.FilePath()] = *element.OcspFile.Value
								alpn = fmt.Sprintf("[alpn %s %s %s %s]", strings.Join(element.Alpn, ","), "ocsp-update on", "ocsp", element.OcspFile.FilePath())
							} else {
								alpn = fmt.Sprintf("[alpn %s %s]", strings.Join(element.Alpn, ","), "ocsp-update on")
							}
						} else {
							alpn = fmt.Sprintf("[alpn %s]", strings.Join(element.Alpn, ","))
						}
					}

					mappings = append(mappings, strings.Join([]string{element.Certificate.FilePath(), alpn, element.SNIFilter, "\n"}, " "))
				}

				sort.Strings(mappings)
				files[bind.SSLCertificateList.FilePath()] = strings.Join(mappings, "")
			}
		}
	}

	return files, nil
}

func (r *Reconciler) loadSSLCertificateValueData(ctx context.Context, instance *proxyv1alpha1.Instance, certificate *configv1alpha1.SSLCertificate) (string, error) {
	if certificate.Value != nil {
		return *certificate.Value, nil
	}

	var items []string

	for _, ref := range certificate.ValueFrom {
		if ref.ConfigMapKeyRef != nil {
			configmap := &corev1.ConfigMap{}
			if err := r.Client.Get(ctx, client.ObjectKey{Name: ref.ConfigMapKeyRef.Name, Namespace: instance.Namespace}, configmap); err != nil {
				return "", err
			}

			data, ok := configmap.Data[ref.ConfigMapKeyRef.Key]
			if !ok {
				return "", fmt.Errorf("key %s not found in SSL certrifcate configmap: %s/%s", ref.ConfigMapKeyRef.Key, instance.Namespace, ref.ConfigMapKeyRef.Name)
			}

			items = append(items, strings.TrimSpace(data))
		}

		if ref.SecretKeyRef != nil {
			secret := &corev1.Secret{}
			if err := r.Client.Get(ctx, client.ObjectKey{Name: ref.SecretKeyRef.Name, Namespace: instance.Namespace}, secret); err != nil {
				return "", err
			}

			data, ok := secret.Data[ref.SecretKeyRef.Key]
			if !ok {
				return "", fmt.Errorf("key %s not found in SSL certrifcate secret: %s/%s", ref.SecretKeyRef.Key, instance.Namespace, ref.SecretKeyRef.Name)
			}

			items = append(items, strings.TrimSpace(string(data)))
		}

		if ref.SecretKeyExternalRef != nil {
			secret := &corev1.Secret{}
			if err := r.Client.Get(ctx, client.ObjectKey{Name: ref.SecretKeyExternalRef.Name, Namespace: ref.SecretKeyExternalRef.Namespace}, secret); err != nil {
				return "", err
			}

			data, ok := secret.Data[ref.SecretKeyExternalRef.Key]
			if !ok {
				return "", fmt.Errorf("key %s not found in SSL certrifcate secret: %s/%s", ref.SecretKeyExternalRef.Key, ref.SecretKeyExternalRef.Namespace, ref.SecretKeyExternalRef.Name)
			}

			items = append(items, strings.TrimSpace(string(data)))
		}
	}

	return strings.Join(items, "\n"), nil
}

func extractSLCCertificatesFromFrontend(frontend *configv1alpha1.Frontend) []*configv1alpha1.SSLCertificate {
	var certificates []*configv1alpha1.SSLCertificate

	for _, bind := range frontend.Spec.Binds {
		if bind.SSL == nil {
			continue
		}

		if bind.SSL.Certificate != nil {
			certificates = append(certificates, bind.SSL.Certificate)
		}
		if bind.SSL.CACertificate != nil {
			certificates = append(certificates, bind.SSL.CACertificate)
		}
	}

	return certificates
}

func extractSLCCertificatesFromBackend(backend *configv1alpha1.Backend) []*configv1alpha1.SSLCertificate {
	var certificates []*configv1alpha1.SSLCertificate

	for _, server := range backend.Spec.Servers {
		if server.SSL == nil {
			continue
		}

		if server.SSL.Certificate != nil {
			certificates = append(certificates, server.SSL.Certificate)
		}
		if server.SSL.CACertificate != nil {
			certificates = append(certificates, server.SSL.CACertificate)
		}
	}

	return certificates
}
