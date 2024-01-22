package instance_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"github.com/six-group/haproxy-operator/controllers/instance"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Reconcile", Label("controller"), func() {
	Context("Reconcile", func() {
		var (
			scheme                                                                        *runtime.Scheme
			ctx                                                                           context.Context
			proxy                                                                         *proxyv1alpha1.Instance
			frontend, frontendCustomCerts, frontendCustomCerts2, frontendCustomCertsEmpty *configv1alpha1.Frontend
			backend, backend2                                                             *configv1alpha1.Backend
			listen                                                                        *configv1alpha1.Listen
			resolver                                                                      *configv1alpha1.Resolver
			initObjs                                                                      []client.Object
		)

		customCert := "Certificate"
		customCertCA := "CAcertificate"
		customCertKey := "Key"
		customCert2 := "Certificate2"
		customCertCA2 := "CAcertificate2"
		customCertKey2 := "Key2"

		BeforeEach(func() {
			scheme = runtime.NewScheme()
			Ω(clientgoscheme.AddToScheme(scheme)).ShouldNot(HaveOccurred())
			Ω(configv1alpha1.AddToScheme(scheme)).ShouldNot(HaveOccurred())
			Ω(proxyv1alpha1.AddToScheme(scheme)).ShouldNot(HaveOccurred())

			ctx = context.Background()

			labels := map[string]string{
				"label-test": "ok",
			}

			dur, _ := time.ParseDuration("30s")

			proxy = &proxyv1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar-foo",
					Namespace: "foo",
					UID:       uuid.NewUUID(),
				},
				Spec: proxyv1alpha1.InstanceSpec{
					Configuration: proxyv1alpha1.Configuration{
						Global: proxyv1alpha1.GlobalConfiguration{
							Logging: &proxyv1alpha1.GlobalLoggingConfiguration{
								Enabled:      true,
								Address:      "/var/lib/rsyslog/rsyslog.sock",
								Facility:     "local0",
								SendHostname: ptr.To(true),
							},
							HardStopAfter: &dur,
						},
					},
					Labels: labels,
				},
			}

			frontend = &configv1alpha1.Frontend{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-front",
					Namespace: "foo",
					Labels:    labels,
				},
			}

			val := "123456789"
			pemCert := strings.Join([]string{customCertKey, customCert, customCertCA}, "\n\n")

			frontendCustomCerts = &configv1alpha1.Frontend{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fe-https-tls-termination",
					Namespace: "foo",
					Labels:    labels,
				},
				Spec: configv1alpha1.FrontendSpec{
					Binds: []configv1alpha1.Bind{
						{
							Address:     "unix@/var/lib/haproxy/run/local.sock",
							Port:        9443,
							Name:        "https",
							AcceptProxy: ptr.To(true),
							Hidden:      ptr.To(true),
							SSL: &configv1alpha1.SSL{
								Enabled: true,
							},
							SSLCertificateList: &configv1alpha1.CertificateList{
								Name: "cert_list",
								Elements: []configv1alpha1.CertificateListElement{
									{
										Certificate: configv1alpha1.SSLCertificate{
											Name:  "route.name4",
											Value: &pemCert,
										},
										SNIFilter: "route.host4",
										Alpn:      []string{"h2", "http/1.0"},
									},
								},
							},
						},
					},
				},
			}

			frontendCustomCerts2 = &configv1alpha1.Frontend{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fe-https-tls-termination2",
					Namespace: "foo",
					Labels:    labels,
				},
				Spec: configv1alpha1.FrontendSpec{
					Binds: []configv1alpha1.Bind{
						{
							Address:     "unix@/var/lib/haproxy/run/local.sock",
							Port:        9443,
							Name:        "https",
							AcceptProxy: ptr.To(true),
							Hidden:      ptr.To(true),
							SSL: &configv1alpha1.SSL{
								Enabled: true,
							},
							SSLCertificateList: &configv1alpha1.CertificateList{
								Name:          "cert_list",
								LabelSelector: metav1.SetAsLabelSelector(labels),
							},
						},
					},
				},
			}

			frontendCustomCertsEmpty = &configv1alpha1.Frontend{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fe-https-tls-termination-empty",
					Namespace: "foo",
					Labels:    labels,
				},
				Spec: configv1alpha1.FrontendSpec{
					Binds: []configv1alpha1.Bind{
						{
							Address:     "unix@/var/lib/haproxy/run/local.sock",
							Port:        9443,
							Name:        "https",
							AcceptProxy: ptr.To(true),
							Hidden:      ptr.To(true),
							SSL: &configv1alpha1.SSL{
								Enabled: true,
							},
						},
					},
				},
			}

			backend = &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-back",
					Namespace: "foo",
					Labels:    labels,
				},
				Spec: configv1alpha1.BackendSpec{
					HostCertificate: &configv1alpha1.CertificateListElement{
						Certificate: configv1alpha1.SSLCertificate{
							Name:  "route.name",
							Value: &pemCert,
						},
						SNIFilter: "route.host",
					},
					Servers: []configv1alpha1.Server{
						{
							Name:    "server",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SSL: &configv1alpha1.SSL{
									Enabled: true,
									Verify:  "required",
									CACertificate: &configv1alpha1.SSLCertificate{
										Name:  "test-ca.crt",
										Value: &val,
									},
									Alpn: []string{"h2", "http/1.0"},
								},
								VerifyHost: "routername.namespace.svc",
								Weight:     ptr.To(int64(256)),
								Check: &configv1alpha1.Check{
									Enabled: true,
									Inter:   &metav1.Duration{Duration: 5 * time.Second},
								},
							},
						},
					},
				},
			}

			pemCert2 := strings.Join([]string{customCertKey2, customCert2, customCertCA2}, "\n\n")

			backend2 = &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-back2",
					Namespace: "foo",
					Labels:    labels,
				},
				Spec: configv1alpha1.BackendSpec{
					HostCertificate: &configv1alpha1.CertificateListElement{
						Certificate: configv1alpha1.SSLCertificate{
							Name:  "route.name2",
							Value: &pemCert2,
						},
						SNIFilter: "route.host2",
						Alpn:      []string{"h2", "http/1.0"},
					},
					Servers: []configv1alpha1.Server{
						{
							Name:    "server",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SSL: &configv1alpha1.SSL{
									Enabled: true,
									Verify:  "required",
									CACertificate: &configv1alpha1.SSLCertificate{
										Name:  "test-ca.crt",
										Value: &val,
									},
									Alpn: []string{"h2", "http/1.0"},
								},
								VerifyHost: "routername.namespace.svc",
								Weight:     ptr.To(int64(256)),
								Check: &configv1alpha1.Check{
									Enabled: true,
									Inter:   &metav1.Duration{Duration: 5 * time.Second},
								},
							},
						},
					},
				},
			}

			listen = &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-listen",
					Namespace: "foo",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: proxyv1alpha1.GroupVersion.String(),
							Kind:       "Instance",
							Name:       proxy.Name,
							UID:        proxy.UID,
						},
					},
				},
				Spec: configv1alpha1.ListenSpec{
					HostCertificate: &configv1alpha1.CertificateListElement{
						Certificate: configv1alpha1.SSLCertificate{
							Name:  "route.name.tcp",
							Value: &pemCert2,
						},
						SNIFilter: "route.host.tcp",
						Alpn:      []string{"h2", "http/1.0"},
					},
					Binds: []configv1alpha1.Bind{
						{
							Address:     "${BIND_ADDRESS}",
							Port:        int64(20005),
							Name:        fmt.Sprintf("tcp-%d", 20005),
							AcceptProxy: ptr.To(true),
							Hidden:      ptr.To(true),
							SSL: &configv1alpha1.SSL{
								Enabled: true,
							},
							SSLCertificateList: &configv1alpha1.CertificateList{
								Name: "cert_list",
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"config.haproxy.com/frontend": "li-tcp",
									},
								},
							},
						},
					},
					Servers: []configv1alpha1.Server{
						{
							Name:    "routeName",
							Address: fmt.Sprintf("%s.%s.svc.cluster.local", "routeName", "routeNamespace"),
							Port:    8443,
							ServerParams: configv1alpha1.ServerParams{
								SSL: &configv1alpha1.SSL{
									Enabled: true,
									Verify:  "required",
									Alpn:    []string{"http/1.1", "h2"},
								},
								Weight:     ptr.To(int64(256)),
								VerifyHost: "routeName" + "." + "routeName" + ".svc",
								InitAddr:   ptr.To("none"),
								Check: &configv1alpha1.Check{
									Enabled: true,
									Inter:   &metav1.Duration{Duration: 500 * time.Millisecond},
								},
								Resolvers: &corev1.LocalObjectReference{
									Name: fmt.Sprintf("dns-%s", "routeNamespace"),
								},
							},
						},
					},
				},
			}

			resolver = &configv1alpha1.Resolver{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar-foo-res",
					Namespace: "foo",
					UID:       uuid.NewUUID(),
					Labels:    labels,
				},
				Spec: configv1alpha1.ResolverSpec{
					ParseResolvConf: ptr.To(true),
					Hold: &configv1alpha1.Hold{
						Nx:    &metav1.Duration{Duration: 500 * time.Millisecond},
						Valid: &metav1.Duration{Duration: 1 * time.Second},
					},
				},
			}

			initObjs = []client.Object{proxy, frontend, frontendCustomCerts, frontendCustomCerts2, frontendCustomCertsEmpty, backend, backend2, listen, resolver}
		})

		It("should deploy haproxy instance", func() {
			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjs...).WithStatusSubresource(initObjs...).Build()
			r := instance.Reconciler{
				Client: cli,
				Scheme: scheme,
			}
			result, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{Name: proxy.Name, Namespace: proxy.Namespace}})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())

			Ω(cli.Get(ctx, client.ObjectKeyFromObject(proxy), proxy)).ShouldNot(HaveOccurred())
			Ω(proxy.Status.Phase).Should(Equal(proxyv1alpha1.InstancePhaseRunning))
			Ω(proxy.Status.Error).Should(BeEmpty())

			secret := &corev1.Secret{}
			Ω(cli.Get(ctx, client.ObjectKey{Namespace: proxy.Namespace, Name: "bar-foo-haproxy-config"}, secret)).ShouldNot(HaveOccurred())
			Ω(string(secret.Data["haproxy.cfg"])).Should(Equal("\nglobal\n  hard-stop-after 30000\n  log /var/lib/rsyslog/rsyslog.sock local0\n  log-send-hostname\n\ndefaults unnamed_defaults_1\n\n" +
				"resolvers bar-foo-res\n  hold nx 500\n  hold valid 1000\n  timeout resolve 1000\n  timeout retry 1000\n  parse-resolv-conf\n  resolve_retries 3\n\n" +
				"frontend fe-https-tls-termination\n  bind unix@/var/lib/haproxy/run/local.sock:9443 name https ssl accept-proxy crt-list /usr/local/etc/haproxy/cert_list.map\n\n" +
				"frontend fe-https-tls-termination-empty\n  bind unix@/var/lib/haproxy/run/local.sock:9443 name https ssl accept-proxy\n\n" +
				"frontend fe-https-tls-termination2\n  bind unix@/var/lib/haproxy/run/local.sock:9443 name https ssl accept-proxy crt-list /usr/local/etc/haproxy/cert_list.map\n\n" +
				"frontend foo-front\n\n" +
				"frontend foo-listen\n  bind ${BIND_ADDRESS}:20005 name tcp-20005 ssl accept-proxy crt-list /usr/local/etc/haproxy/cert_list.map\n  default_backend foo-listen\n\n" +
				"backend foo-back\n  server server localhost:80 check ssl alpn h2,http/1.0 ca-file /usr/local/etc/haproxy/test-ca.crt inter 5000 verify required verifyhost routername.namespace.svc weight 256\n\n" +
				"backend foo-back2\n  server server localhost:80 check ssl alpn h2,http/1.0 ca-file /usr/local/etc/haproxy/test-ca.crt inter 5000 verify required verifyhost routername.namespace.svc weight 256\n\n" +
				"backend foo-listen\n  server routeName routeName.routeNamespace.svc.cluster.local:8443 check ssl alpn http/1.1,h2 init-addr none inter 500 resolvers dns-routeNamespace verify required verifyhost routeName.routeName.svc weight 256\n"))

			statefulSet := &appsv1.StatefulSet{}
			Ω(cli.Get(ctx, client.ObjectKey{Namespace: proxy.Namespace, Name: "bar-foo-haproxy"}, statefulSet)).ShouldNot(HaveOccurred())
			Ω(statefulSet.Spec.Template.ObjectMeta.Labels["app.kubernetes.io/name"]).Should(Equal(proxy.Name + "-haproxy"))
			Ω(statefulSet.Spec.Template.ObjectMeta.Labels["label-test"]).Should(Equal("ok"))
		})

		It("same resource names error", func() {
			backend.Kind = "Backend"
			backend.Name = "foo"

			frontend.Kind = "Frontend"
			frontend.Name = "foo"

			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjs...).WithStatusSubresource(initObjs...).Build()
			r := instance.Reconciler{
				Client: cli,
				Scheme: scheme,
			}
			result, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{Name: proxy.Name, Namespace: proxy.Namespace}})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())

			Ω(cli.Get(ctx, client.ObjectKeyFromObject(proxy), proxy)).ShouldNot(HaveOccurred())
			Ω(proxy.Status.Phase).Should(Equal(proxyv1alpha1.InstancePhaseInternalError))
			Ω(proxy.Status.Error).Should(Equal("name foo already used by resource of kind Frontend"))

			backendRes := &configv1alpha1.Backend{}
			Ω(cli.Get(ctx, client.ObjectKey{Namespace: proxy.Namespace, Name: backend.Name}, backendRes)).ShouldNot(HaveOccurred())
			Ω(backendRes.Status.Error).Should(Equal(proxy.Status.Error))
		})

		It("should set status to pending if there is no listens", func() {
			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(proxy).WithStatusSubresource(proxy).Build()
			r := instance.Reconciler{
				Client: cli,
				Scheme: scheme,
			}
			result, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{Name: proxy.Name, Namespace: proxy.Namespace}})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())

			Ω(cli.Get(ctx, client.ObjectKeyFromObject(proxy), proxy)).ShouldNot(HaveOccurred())
			Ω(proxy.Status.Phase).Should(Equal(proxyv1alpha1.InstancePhasePending))
			Ω(proxy.Status.Error).ShouldNot(BeEmpty())
		})
		It("should create custom certs", func() {
			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjs...).WithStatusSubresource(initObjs...).Build()
			r := instance.Reconciler{
				Client: cli,
				Scheme: scheme,
			}
			result, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{Name: proxy.Name, Namespace: proxy.Namespace}})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())

			Ω(cli.Get(ctx, client.ObjectKeyFromObject(proxy), proxy)).ShouldNot(HaveOccurred())
			Ω(proxy.Status.Phase).Should(Equal(proxyv1alpha1.InstancePhaseRunning))
			Ω(proxy.Status.Error).Should(BeEmpty())

			secret := &corev1.Secret{}
			Ω(cli.Get(ctx, client.ObjectKey{Namespace: proxy.Namespace, Name: "bar-foo-haproxy-config"}, secret)).ShouldNot(HaveOccurred())
			Ω(string(secret.Data["haproxy.cfg"])).Should(Equal("\nglobal\n  hard-stop-after 30000\n  log /var/lib/rsyslog/rsyslog.sock local0\n  log-send-hostname\n\ndefaults unnamed_defaults_1\n\n" +
				"resolvers bar-foo-res\n  hold nx 500\n  hold valid 1000\n  timeout resolve 1000\n  timeout retry 1000\n  parse-resolv-conf\n  resolve_retries 3\n\n" +
				"frontend fe-https-tls-termination\n  bind unix@/var/lib/haproxy/run/local.sock:9443 name https ssl accept-proxy crt-list /usr/local/etc/haproxy/cert_list.map\n\n" +
				"frontend fe-https-tls-termination-empty\n  bind unix@/var/lib/haproxy/run/local.sock:9443 name https ssl accept-proxy\n\n" +
				"frontend fe-https-tls-termination2\n  bind unix@/var/lib/haproxy/run/local.sock:9443 name https ssl accept-proxy crt-list /usr/local/etc/haproxy/cert_list.map\n\n" +
				"frontend foo-front\n\n" +
				"frontend foo-listen\n  bind ${BIND_ADDRESS}:20005 name tcp-20005 ssl accept-proxy crt-list /usr/local/etc/haproxy/cert_list.map\n  default_backend foo-listen\n\n" +
				"backend foo-back\n  server server localhost:80 check ssl alpn h2,http/1.0 ca-file /usr/local/etc/haproxy/test-ca.crt inter 5000 verify required verifyhost routername.namespace.svc weight 256\n\n" +
				"backend foo-back2\n  server server localhost:80 check ssl alpn h2,http/1.0 ca-file /usr/local/etc/haproxy/test-ca.crt inter 5000 verify required verifyhost routername.namespace.svc weight 256\n\n" +
				"backend foo-listen\n  server routeName routeName.routeNamespace.svc.cluster.local:8443 check ssl alpn http/1.1,h2 init-addr none inter 500 resolvers dns-routeNamespace verify required verifyhost routeName.routeName.svc weight 256\n"))

			Ω(string(secret.Data["cert_list.map"])).Should(Equal(
				"/usr/local/etc/haproxy/route.name.crt  route.host \n" +
					"/usr/local/etc/haproxy/route.name.tcp.crt [alpn h2,http/1.0] route.host.tcp \n" +
					"/usr/local/etc/haproxy/route.name2.crt [alpn h2,http/1.0] route.host2 \n" +
					"/usr/local/etc/haproxy/route.name4.crt [alpn h2,http/1.0] route.host4 \n",
			),
			)

			Ω(string(secret.Data["route.name.crt"])).Should(Equal("Key\n\nCertificate\n\nCAcertificate"))
			Ω(string(secret.Data["route.name2.crt"])).Should(Equal("Key2\n\nCertificate2\n\nCAcertificate2"))
			Ω(string(secret.Data["route.name.tcp.crt"])).Should(Equal("Key2\n\nCertificate2\n\nCAcertificate2"))
			Ω(string(secret.Data["route.name4.crt"])).Should(Equal("Key\n\nCertificate\n\nCAcertificate"))
		})
	})
})
