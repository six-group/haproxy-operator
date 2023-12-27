package v1alpha1_test

import (
	"time"

	parser "github.com/haproxytech/config-parser/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var simpleBackend = `
backend foo
`

var _ = Describe("Backend", Label("type"), func() {
	Context("AddToParser", func() {
		var p parser.Parser
		BeforeEach(func() {
			var err error
			p, err = parser.New()
			Ω(err).ShouldNot(HaveOccurred())
		})
		// valid
		It("should create backend/frontend", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(Equal(simpleBackend))
		})
		It("should set mode http", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						Mode: "http",
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("mode http"))
		})
		It("should set option forwardfor", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						Forwardfor: &configv1alpha1.Forwardfor{
							Enabled: true,
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("option forwardfor"))
		})
		It("should set option redispatch", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.BackendSpec{
					Redispatch: pointer.Bool(true),
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("option redispatch"))
		})
		It("should set hash-type", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.BackendSpec{
					HashType: &configv1alpha1.HashType{
						Method:   "consistent",
						Function: "djb2",
						Modifier: "avalanche",
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("hash-type consistent djb2 avalanche"))
		})

		It("should set ssl parameters", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.BackendSpec{
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
										Name: "test-ca.crt",
									},
									Alpn: []string{"h2", "http/1.0"},
								},
								Weight: pointer.Int64(256),
								Check: &configv1alpha1.Check{
									Enabled: true,
									Inter:   &metav1.Duration{Duration: 5 * time.Second},
								},
								VerifyHost: "routername.namespace.svc",
								Cookie:     true,
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("ssl alpn h2,http/1.0 ca-file /usr/local/etc/haproxy/test-ca.crt cookie 1c3c2192e2912699ccd31119b162666a inter 5000 verify required verifyhost routername.namespace.svc weight 256"))
		})

		It("should set option http-request deny", func() {
			var notFound int64 = 404
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPRequest: &configv1alpha1.HTTPRequestRules{
							Deny: &configv1alpha1.Deny{
								Rule: configv1alpha1.Rule{
									ConditionType: "if",
									Condition:     "{ var(my-ip) -m ip 127.0.0.0/8 10.0.0.0/8 }",
								},
								Enabled: true,
							},
							DenyStatus: &notFound,
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("http-request deny deny_status 404 if { var(my-ip) -m ip 127.0.0.0/8 10.0.0.0/8 }\n"))
		})

		It("should set option http-pretend-keepalive", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPPretendKeepalive: pointer.Bool(true),
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("option http-pretend-keepalive\n"))
		})

		It("should set option forwardfor", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						Forwardfor: &configv1alpha1.Forwardfor{
							Enabled: true,
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("option forwardfor\n"))
		})
		It("should set cookie", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "set_cookie"},
				Spec: configv1alpha1.BackendSpec{
					Cookie: &configv1alpha1.Cookie{
						Name: "test",
						Mode: configv1alpha1.CookieMode{
							Rewrite: true,
						},
						Indirect: pointer.Bool(true),
						NoCache:  pointer.Bool(true),
						PostOnly: pointer.Bool(true),
						Preserve: pointer.Bool(true),
						HTTPOnly: pointer.Bool(true),
						Secure:   pointer.Bool(true),
						Domain: []string{
							"domain1", ".openshift",
						},

						MaxIdle:   120,
						MaxLife:   45,
						Attribute: []string{"SameSite=None"},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("cookie 098f6bcd4621d373cade4e832627b4f6 domain domain1 domain .openshift attr SameSite=None httponly indirect maxidle 120 maxlife 45 nocache postonly preserve rewrite secure\n"))
		})
		It("an error should appear for selecting more than one cookie mode", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "set_cookie"},
				Spec: configv1alpha1.BackendSpec{
					Cookie: &configv1alpha1.Cookie{
						Name: "test",
						Mode: configv1alpha1.CookieMode{
							Rewrite: true,
							Insert:  true,
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).Should(HaveOccurred())
		})
		It("should set ssl parameters", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.BackendSpec{
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
										Name: "test-ca.crt",
									},
									Alpn: []string{"h2", "http/1.0"},
								},
								Weight: pointer.Int64(256),
								Check: &configv1alpha1.Check{
									Enabled: true,
									Inter:   &metav1.Duration{Duration: 5 * time.Second},
								},
								VerifyHost: "routername.namespace.svc",
								Cookie:     true,
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("ssl alpn h2,http/1.0 ca-file /usr/local/etc/haproxy/test-ca.crt cookie 1c3c2192e2912699ccd31119b162666a inter 5000 verify required verifyhost routername.namespace.svc weight 256"))
		})
		It("should set option http-request redirect location", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPRequest: &configv1alpha1.HTTPRequestRules{
							Redirect: []configv1alpha1.Redirect{
								{
									Rule: configv1alpha1.Rule{
										ConditionType: "unless",
										Condition:     "has_www",
									},
									Code: pointer.Int64(301),
									Type: configv1alpha1.RedirectType{
										Location: true,
									},
									Value: "host",
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("http-request redirect location host code 301 unless has_www\n"))
		})
		It("should set option http-request redirect prefix", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPRequest: &configv1alpha1.HTTPRequestRules{
							Redirect: []configv1alpha1.Redirect{
								{
									Rule: configv1alpha1.Rule{
										ConditionType: "unless",
										Condition:     "begins_with_api",
									},
									Code: pointer.Int64(301),
									Type: configv1alpha1.RedirectType{
										Prefix: true,
									},
									Value: "/api/v2",
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("http-request redirect prefix /api/v2 code 301 unless begins_with_api\n"))
		})
		It("should set option http-request redirect scheme", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPRequest: &configv1alpha1.HTTPRequestRules{
							Redirect: []configv1alpha1.Redirect{
								{
									Rule: configv1alpha1.Rule{
										ConditionType: "unless",
										Condition:     "is_https",
									},
									Type: configv1alpha1.RedirectType{
										Scheme: true,
									},
									Value: "https",
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("http-request redirect scheme https unless is_https\n"))
		})
		It("an error should appear for selecting more than one redirect type", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPRequest: &configv1alpha1.HTTPRequestRules{
							Redirect: []configv1alpha1.Redirect{
								{
									Type: configv1alpha1.RedirectType{
										Scheme:   true,
										Location: true,
									},
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).Should(HaveOccurred())
		})
		It("should set option http-request redirect with options", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPRequest: &configv1alpha1.HTTPRequestRules{
							Redirect: []configv1alpha1.Redirect{
								{
									Type: configv1alpha1.RedirectType{
										Location: true,
									},
									Value: "https",
									Option: &configv1alpha1.RedirectOption{
										DropQuery:   true,
										AppendSlash: true,
										SetCookie: &configv1alpha1.RedirectCookie{
											Name:  "classic",
											Value: "=1",
										},
										ClearCookie: &configv1alpha1.RedirectCookie{
											Name:  "classic",
											Value: "=",
										},
									},
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("http-request redirect location https drop-query append-slash set-cookie CLASSIC=1 clear-cookie CLASSIC=\n"))
		})
		It("should set sendProxy with proxy protocol v1", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					Servers: []configv1alpha1.Server{
						{
							Name:    "server1",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SendProxyV2: &configv1alpha1.ProxyProtocol{
									V1: true,
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("server server1 localhost:80 send-proxy\n"))
		})
		It("should set sendProxy with proxy protocol v2ssl", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					Servers: []configv1alpha1.Server{
						{
							Name:    "server1",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SendProxyV2: &configv1alpha1.ProxyProtocol{
									V2SSL: true,
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("server server1 localhost:80 send-proxy-v2-ssl\n"))
		})
		It("an error should appear for selecting more than one proxy protocol", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					Servers: []configv1alpha1.Server{
						{
							Name:    "server1",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SendProxyV2: &configv1alpha1.ProxyProtocol{
									V2SSL: true,
									V1:    true,
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).Should(HaveOccurred())
		})
		It("should not set send proxy if no proxy protocol has been defined", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					Servers: []configv1alpha1.Server{
						{
							Name:    "server1",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SendProxyV2: &configv1alpha1.ProxyProtocol{},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).Should(HaveOccurred())
		})
		It("should set sendProxy with proxy protocol v2 and options", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					Servers: []configv1alpha1.Server{
						{
							Name:    "server1",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SendProxyV2: &configv1alpha1.ProxyProtocol{
									V2: &configv1alpha1.ProxyProtocolV2{
										Enabled: true,
										Options: &configv1alpha1.ProxyProtocolV2Options{
											CertCn:   true,
											CertSig:  true,
											UniqueID: true,
											Ssl:      true,
										},
									},
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("server server1 localhost:80 send-proxy-v2 proxy-v2-options ssl,cert-cn,cert-sig,unique-id\n"))
		})
		It("should set sendProxy with proxy protocol v2 without options", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					Servers: []configv1alpha1.Server{
						{
							Name:    "server1",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SendProxyV2: &configv1alpha1.ProxyProtocol{
									V2: &configv1alpha1.ProxyProtocolV2{
										Enabled: true,
									},
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("server server1 localhost:80 send-proxy-v2\n"))
		})
		It("should not proxyProtocolOptions if proxyProtocol is not V2", func() {
			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{Name: "openshift_default"},
				Spec: configv1alpha1.BackendSpec{
					Servers: []configv1alpha1.Server{
						{
							Name:    "server1",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SendProxyV2: &configv1alpha1.ProxyProtocol{
									V2SSLCN: true,
									V2: &configv1alpha1.ProxyProtocolV2{
										Options: &configv1alpha1.ProxyProtocolV2Options{
											CertCn:   true,
											CertSig:  true,
											UniqueID: true,
											Ssl:      true,
										},
									},
								},
							},
						},
					},
				},
			}
			Ω(backend.AddToParser(p)).Should(HaveOccurred())
		})
	})
})
