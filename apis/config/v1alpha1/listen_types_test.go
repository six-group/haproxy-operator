package v1alpha1_test

import (
	"time"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

var simple = `
frontend foo
  default_backend foo

backend foo
`

var _ = Describe("Listen", Label("type"), func() {
	Context("AddToParser", func() {
		var p parser.Parser
		BeforeEach(func() {
			var err error
			p, err = parser.New()
			Ω(err).ShouldNot(HaveOccurred())
		})
		// valid
		It("should create backend/frontend", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(Equal(simple))
		})
		It("should set mode tcp", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						Mode: "tcp",
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("mode tcp"))
		})
		It("should set set-header and add-header", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPResponse: &configv1alpha1.HTTPResponseRules{
							SetHeader: []configv1alpha1.HTTPHeaderRule{
								{
									Rule: configv1alpha1.Rule{
										ConditionType: "if",
										Condition:     "!{ ssl_fc }",
									},
									Name: "Strict-Transport-Security",
									Value: configv1alpha1.HTTPHeaderValue{
										Str: ptr.To("max-age=16000000; includeSubDomains; preload;"),
									},
								},
							},
						},
						HTTPRequest: &configv1alpha1.HTTPRequestRules{
							SetHeader: []configv1alpha1.HTTPHeaderRule{
								{
									Rule: configv1alpha1.Rule{
										ConditionType: "if",
										Condition:     "!{ ssl_fc }",
									},
									Name: "X-Forwarded-Host",
									Value: configv1alpha1.HTTPHeaderValue{
										Str: ptr.To("local.com"),
									},
								},
								{
									Name: "X-Forwarded-Port",
									Value: configv1alpha1.HTTPHeaderValue{
										Str: ptr.To("8055"),
									},
								},
							},
							SetPath: []configv1alpha1.HTTPPathRule{
								{
									Rule: configv1alpha1.Rule{
										ConditionType: "if",
										Condition:     "!{ ssl_fc }",
									},
									Value: "/metrics",
								},
							},
							AddHeader: []configv1alpha1.HTTPHeaderRule{
								{
									Rule: configv1alpha1.Rule{
										ConditionType: "unless",
										Condition:     "{ ssl_fc_alpn -i h2 }",
									},
									Name: "SOAPAction",
									Value: configv1alpha1.HTTPHeaderValue{
										Str: ptr.To("\"urn:mediate\""),
									},
								},
								{
									Name: "X-Forwarded-Proto",
									Value: configv1alpha1.HTTPHeaderValue{
										Str:    ptr.To("s"),
										Format: ptr.To("http%s"),
									},
								},
								{
									Name: "X-Forwarded-Proto-Version",
									Value: configv1alpha1.HTTPHeaderValue{
										Env: &corev1.EnvVar{
											Name: "PROTO_VERSION",
										},
										Format: ptr.To("\"%s\""),
									},
								},
							},
							DelHeader: []configv1alpha1.HTTPDeleteHeaderRule{
								{
									Name:   "regex",
									Method: "str",
								},
								{
									Name: "Proxy",
								},
							},
						},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("http-request set-header X-Forwarded-Host local.com if !{ ssl_fc }"))
			Ω(p.String()).Should(ContainSubstring("http-request set-header X-Forwarded-Port 8055"))
			Ω(p.String()).Should(ContainSubstring("http-request add-header SOAPAction \"urn:mediate\" unless { ssl_fc_alpn -i h2 }"))
			Ω(p.String()).Should(ContainSubstring("http-request add-header X-Forwarded-Proto https"))
			Ω(p.String()).Should(ContainSubstring("http-request add-header X-Forwarded-Proto-Version \"${PROTO_VERSION}\""))
			Ω(p.String()).Should(ContainSubstring("http-request set-path /metrics if !{ ssl_fc }"))
			Ω(p.String()).Should(ContainSubstring("http-request del-header Proxy"))
			Ω(p.String()).Should(ContainSubstring("http-request del-header regex -m str"))
			Ω(p.String()).Should(ContainSubstring("http-response set-header Strict-Transport-Security max-age=16000000; includeSubDomains; preload; if !{ ssl_fc }"))
		})
		It("should create binds", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					Binds: []configv1alpha1.Bind{
						{Name: "bind01", Port: 80},
						{Name: "bind02", Port: 81},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("bind :80 name bind01"))
			Ω(p.String()).Should(ContainSubstring("bind :81 name bind02"))
		})
		It("should create servers", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					Servers: []configv1alpha1.Server{
						{Name: "server01", Port: 80, Address: "localhost"},
						{Name: "server02", Port: 81, Address: "localhost"},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("server server01 localhost:80"))
			Ω(p.String()).Should(ContainSubstring("server server02 localhost:81"))
		})
		It("should create server templates", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					ServerTemplates: []configv1alpha1.ServerTemplate{
						{FQDN: "google.com", NumMin: ptr.To(int64(1)), Num: 3, Port: 80, Prefix: "srv", ServerParams: configv1alpha1.ServerParams{Check: &configv1alpha1.Check{Enabled: true}}},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("server-template srv 1-3 google.com:80 check"))
		})
		It("should set option forwardfor", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						Forwardfor: &configv1alpha1.Forwardfor{
							Enabled: true,
						},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("option forwardfor"))
		})
		It("should set http log option", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPLog: ptr.To(true),
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("option httplog"))
		})
		It("should set option redispatch", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					Redispatch: ptr.To(true),
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("option redispatch"))
		})
		It("should set hash-type", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					HashType: &configv1alpha1.HashType{
						Method:   "consistent",
						Function: "djb2",
						Modifier: "avalanche",
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("hash-type consistent djb2 avalanche"))
		})
		It("should set ssl parameters", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					Binds: []configv1alpha1.Bind{
						{
							Name: "bind",
							Port: 80,
							SSL: &configv1alpha1.SSL{
								Enabled:    true,
								MinVersion: "SSLv3",
								Verify:     "required",
								CACertificate: &configv1alpha1.SSLCertificate{
									Name: "test-ca.crt",
								},
								Certificate: &configv1alpha1.SSLCertificate{
									Name: "test.crt",
								},
							},
							SSLCertificateList: &configv1alpha1.CertificateList{
								Name: "cert_list.map",
							},
						},
					},
					Servers: []configv1alpha1.Server{
						{
							Name:    "server",
							Port:    80,
							Address: "localhost",
							ServerParams: configv1alpha1.ServerParams{
								SSL: &configv1alpha1.SSL{
									Enabled:    true,
									MinVersion: "TLSv1.3",
									Verify:     "none",
									CACertificate: &configv1alpha1.SSLCertificate{
										Name: "test-ca.crt",
									},
									Certificate: &configv1alpha1.SSLCertificate{
										Name: "test.crt",
									},
									SNI: "str(localhost)",
								},
								Weight: ptr.To(int64(256)),
								Check: &configv1alpha1.Check{
									Enabled: true,
									Inter:   &metav1.Duration{Duration: 5 * time.Second},
								},
								Cookie: true,
							},
						},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("crt /usr/local/etc/haproxy/test.crt ca-file /usr/local/etc/haproxy/test-ca.crt ssl verify required crt-list /usr/local/etc/haproxy/cert_list.map ssl-min-ver SSLv3"))
			Ω(p.String()).Should(ContainSubstring("ssl ca-file /usr/local/etc/haproxy/test-ca.crt cookie 1c3c2192e2912699ccd31119b162666a crt /usr/local/etc/haproxy/test.crt inter 5000 sni str(localhost) ssl-min-ver TLSv1.3 weight 256"))
		})
		It("should set load balancer algorithm", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					Balance: &configv1alpha1.Balance{
						Algorithm: "roundrobin",
					},
					Binds: []configv1alpha1.Bind{
						{Name: "bind", Port: 80},
					},
					Servers: []configv1alpha1.Server{
						{Name: "server", Port: 80, Address: "localhost"},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("balance roundrobin"))
		})
		It("should set load balancer algorithm", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						ACL: []configv1alpha1.ACL{
							{Name: "whitelist", Criterion: "src", Values: []string{"10.0.0.1/32", "10.0.0.2/32"}},
						},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("acl whitelist src 10.0.0.1/32 10.0.0.2/32"))
		})
		It("should set error files", func() {
			fileCntnt403 := "HTTP/1.1 403 Not Authorized\n\n<head><meta charset=\"utf-8\"/><title>403 Not Authorized</title></head></html>"
			fileCntnt404 := "HTTP/1.1 404 Not Found\n\n<head><meta charset=\"utf-8\"/><title>404 Not Found</title></head></html>"
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						ErrorFiles: []*configv1alpha1.ErrorFile{
							{
								Code: 403,
								File: configv1alpha1.StaticHTTPFile{
									Name:  "error-file-403.http",
									Value: &fileCntnt403,
								},
							},
							{
								Code: 404,
								File: configv1alpha1.StaticHTTPFile{
									Name:  "error-file-404.http",
									Value: &fileCntnt404,
								},
							},
						},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			cfgDump := p.String()
			Ω(cfgDump).Should(ContainSubstring("errorfile 403 /usr/local/etc/haproxy/error-file-403.http"))
			Ω(cfgDump).Should(ContainSubstring("errorfile 404 /usr/local/etc/haproxy/error-file-404.http"))
		})
		It("should set timeouts", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						Timeouts: map[string]metav1.Duration{
							"client": {Duration: 30 * time.Second},
							"tunnel": {Duration: 1 * time.Hour},
							"server": {Duration: 30 * time.Second},
						},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("timeout client 30000"))
			Ω(p.String()).Should(ContainSubstring("timeout tunnel 3600000"))
			Ω(p.String()).Should(ContainSubstring("timeout server 30000"))
		})
		// invalid
		It("should fail if name is not defined", func() {
			listen := &configv1alpha1.Listen{}
			Ω(listen.AddToParser(p)).Should(HaveOccurred())
		})
		It("should fail if port is undefined for bind", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					Binds: []configv1alpha1.Bind{
						{Name: "bind"},
					},
				},
			}
			Ω(listen.AddToParser(p)).Should(HaveOccurred())
		})
		It("should fail if address is undefined for server", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					Servers: []configv1alpha1.Server{
						{Name: "server", Port: 80},
					},
				},
			}
			Ω(listen.AddToParser(p)).Should(HaveOccurred())
		})
		It("should fail if mode is invalid", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						Mode: "udp",
					},
				},
			}
			Ω(listen.AddToParser(p)).Should(HaveOccurred())
		})
		It("should fail if ssl is invalid", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ListenSpec{
					Binds: []configv1alpha1.Bind{
						{
							Name: "bind",
							Port: 80,
							SSL: &configv1alpha1.SSL{
								Enabled:    true,
								MinVersion: "XSSLv3",
								Verify:     "disabled",
								CACertificate: &configv1alpha1.SSLCertificate{
									Name: "test-ca.crt",
								},
								Certificate: &configv1alpha1.SSLCertificate{
									Name: "test.crt",
								},
							},
						},
					},
				},
			}
			Ω(listen.AddToParser(p)).Should(HaveOccurred())
		})
		It("should set cookie", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "set_cookie"},
				Spec: configv1alpha1.ListenSpec{
					Cookie: &configv1alpha1.Cookie{
						Name: "cookie_name",
						Mode: configv1alpha1.CookieMode{
							Rewrite: true,
						},
						Indirect:  ptr.To(true),
						NoCache:   ptr.To(true),
						PostOnly:  ptr.To(true),
						Preserve:  ptr.To(true),
						HTTPOnly:  ptr.To(true),
						Secure:    ptr.To(true),
						Domain:    []string{"domain1", ".openshift"},
						MaxIdle:   120,
						MaxLife:   45,
						Attribute: []string{"SameSite=None"},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("cookie e3cb9741ffde596f46710a5d7e3ec587 domain domain1 domain .openshift attr SameSite=None httponly indirect maxidle 120 maxlife 45 nocache postonly preserve rewrite secure\n"))
		})
		It("should set http-request return", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{Name: "return"},
				Spec: configv1alpha1.ListenSpec{
					BaseSpec: configv1alpha1.BaseSpec{
						HTTPRequest: &configv1alpha1.HTTPRequestRules{
							Return: &configv1alpha1.HTTPReturn{
								Status: ptr.To(int64(200)),
								Content: configv1alpha1.HTTPReturnContent{
									Type:   "text/plain",
									Format: "lf-string",
									Value:  "Hello World",
								},
							},
						},
					},
				},
			}
			Ω(listen.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("http-request return status 200 content-type text/plain lf-string \"Hello World\"\n"))
		})
	})
})
