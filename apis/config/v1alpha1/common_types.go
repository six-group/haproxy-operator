package v1alpha1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/models"
	parser "github.com/haproxytech/config-parser/v5"
	"github.com/six-group/haproxy-operator/pkg/defaults"
	"github.com/six-group/haproxy-operator/pkg/hash"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +k8s:deepcopy-gen=false

type Object interface {
	client.Object
	SetStatus(status Status)
	GetStatus() Status
	AddToParser(p parser.Parser) error
}

type BaseSpec struct {
	// Mode can be either 'tcp' or 'http'. In TCP mode it is a layer 4 proxy. In HTTP mode it is a layer 7 proxy.
	// +kubebuilder:default=http
	// +kubebuilder:validation:Enum=http;tcp
	Mode string `json:"mode"`
	// HTTPRequest rules define a set of rules which apply to layer 7 processing.
	// +optional
	HTTPRequest *HTTPRequestRules `json:"httpRequest,omitempty"`
	// TCPRequest rules perform an action on an incoming connection depending on a layer 4 condition.
	// +optional
	TCPRequest []TCPRequestRule `json:"tcpRequest,omitempty"`
	// ACL (Access Control Lists) provides a flexible solution to perform
	// content switching and generally to take decisions based on content extracted
	// from the request, the response or any environmental status
	// +optional
	ACL []ACL `json:"acl,omitempty"`
	// Timeouts: check, connect, http-keep-alive, http-request, queue, server, tunnel.
	// The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit.
	// More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html
	// +optional
	Timeouts map[string]metav1.Duration `json:"timeouts"`
	// ErrorFiles custom error files to be used
	// +optional
	ErrorFiles []*ErrorFile `json:"errorFiles,omitempty"`
	// Forwardfor enable insertion of the X-Forwarded-For header to requests sent to servers
	// +optional
	Forwardfor *Forwardfor `json:"forwardFor,omitempty"`
	// HTTPPretendKeepalive will keep the connection alive. It is recommended not to enable this option by default.
	// +optional
	HTTPPretendKeepalive *bool `json:"httpPretendKeepalive,omitempty"`
}

func (b *BaseSpec) AddToParser(p parser.Parser, sectionType parser.Section, sectionName string) error {
	for idx, acl := range b.ACL {
		model, err := acl.Model()
		if err != nil {
			return err
		}

		err = p.Insert(sectionType, sectionName, "acl", configuration.SerializeACL(model), idx)
		if err != nil {
			return err
		}
	}

	for idx, rule := range b.TCPRequest {
		model, err := rule.Model()
		if err != nil {
			return err
		}

		data, err := configuration.SerializeTCPRequestRule(model)
		if err != nil {
			return err
		}

		err = p.Insert(sectionType, sectionName, "tcp-request", data, idx)
		if err != nil {
			return err
		}
	}

	if b.HTTPRequest != nil {
		rules, err := b.HTTPRequest.Model()
		if err != nil {
			return err
		}
		for idx, rule := range rules {
			if rule != nil {
				data, err := configuration.SerializeHTTPRequestRule(*rule)
				if err != nil {
					return err
				}
				err = p.Insert(sectionType, sectionName, "http-request", data, idx)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

type HashType struct {
	// +kubebuilder:validation:Enum=map-based;consistent
	// +optional
	Method string `json:"method,omitempty"`
	// +kubebuilder:validation:Enum=sdbm;djb2;wt6;crc32
	// +optional
	Function string `json:"function,omitempty"`
	// +kubebuilder:validation:Enum=avalanche
	// +optional
	Modifier string `json:"modifier,omitempty"`
}

func (h *HashType) Model() (*models.HashType, error) {
	model := &models.HashType{
		Function: h.Function,
		Method:   h.Method,
		Modifier: h.Modifier,
	}

	return model, model.Validate(strfmt.Default)
}

type Rule struct {
	// ConditionType specifies the type of the condition matching ('if' or 'unless')
	// +kubebuilder:validation:Enum=if;unless
	// +optional
	ConditionType string `json:"conditionType,omitempty"`
	// Condition is a condition composed of ACLs.
	// +optional
	Condition string `json:"condition,omitempty"`
}

type TCPRequestRule struct {
	Rule `json:",inline"`
	// Type specifies the type of the tcp-request rule.
	// +kubebuilder:validation:Enum=connection;content;inspect-delay;session
	Type string `json:"type"`
	// Action defines the action to perform if the condition applies.
	// +kubebuilder:validation:Enum=accept;capture;do-resolve;expect-netscaler-cip;expect-proxy;reject;sc-inc-gpc0;sc-inc-gpc1;sc-set-gpt0;send-spoe-group;set-dst-port;set-dst;set-priority;set-src;set-var;silent-drop;track-sc0;track-sc1;track-sc2;unset-var;use-service;lua
	// +optional
	Action *string `json:"action"`
	// Timeout sets timeout for the action
	// +optional
	Timeout *metav1.Duration `json:"timeout"`
}

func (t *TCPRequestRule) Model() (models.TCPRequestRule, error) {
	model := models.TCPRequestRule{
		Cond:     t.ConditionType,
		CondTest: t.Condition,
		Type:     t.Type,
		Index:    ptr.To(int64(0)),
	}

	if t.Action != nil {
		model.Action = *t.Action
	}

	if t.Timeout != nil {
		model.Timeout = ptr.To(t.Timeout.Milliseconds())
	}

	return model, model.Validate(strfmt.Default)
}

type ACL struct {
	// Name
	// +kubebuilder:validation:Pattern=^[^\s]+$
	Name string `json:"name"`
	// Criterion is the name of a sample fetch method, or one of its ACL
	// specific declinations.
	// +kubebuilder:validation:Pattern=^[^\s]+$
	Criterion string `json:"criterion"`
	// Values are of the type supported by the criterion.
	Values []string `json:"values"`
}

func (a *ACL) Model() (models.ACL, error) {
	values := strings.Join(a.Values, " ")

	if len(a.Values) > defaults.MaxLineArgs-3 {
		values = fmt.Sprintf("-f %s", a.FilePath())
	}

	model := models.ACL{
		ACLName:   a.Name,
		Criterion: a.Criterion,
		Value:     values,
		Index:     ptr.To(int64(0)),
	}

	return model, model.Validate(strfmt.Default)
}

func (a *ACL) FilePath() string {
	return fmt.Sprintf("/usr/local/etc/haproxy/acl-%s-%s.txt", a.Name, hash.GetMD5Hash(strings.Join(a.Values, "\n")))
}

type ErrorFile struct {
	// Code is the HTTP status code.
	// +kubebuilder:validation:Enum=200;400;401;403;404;405;407;408;410;413;425;429;500;501;502;503;504
	Code int64 `json:"code"`
	// File designates a file containing the full HTTP response.
	File StaticHTTPFile `json:"file"`
}

type StaticHTTPFile struct {
	Name      string             `json:"name"`
	Value     *string            `json:"value,omitempty"`
	ValueFrom ErrorFileValueFrom `json:"valueFrom,omitempty"`
}

func (s *StaticHTTPFile) FilePath() string {
	return fmt.Sprintf("/usr/local/etc/haproxy/%s.http", strings.TrimSuffix(s.Name, ".http"))
}

type ErrorFileValueFrom struct {
	// ConfigMapKeyRef selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
}

func (e *ErrorFile) Model() (models.Errorfile, error) {
	model := models.Errorfile{
		Code: e.Code,
		File: e.File.FilePath(),
	}

	return model, model.Validate(strfmt.Default)
}

type Bind struct {
	// Name for these sockets, which will be reported on the stats page.
	Name string `json:"name"`
	// Address can be a host name, an IPv4 address, an IPv6 address, or '*' (is equal to the special address "0.0.0.0").
	// +kubebuilder:validation:Pattern=^[^\s]+$
	// +optional
	Address string `json:"address,omitempty"`
	// Port
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=1
	Port int64 `json:"port"`
	// PortRangeEnd if set it must be greater than Port
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=1
	// +optional
	PortRangeEnd *int64 `json:"portRangeEnd,omitempty"`
	// Transparent is an optional keyword which is supported only on certain Linux kernels. It
	// indicates that the addresses will be bound even if they do not belong to the
	// local machine, and that packets targeting any of these addresses will be
	// intercepted just as if the addresses were locally configured. This normally
	// requires that IP forwarding is enabled. Caution! do not use this with the
	// default address '*', as it would redirect any traffic for the specified port.
	// +optional
	Transparent bool `json:"transparent,omitempty"`
	// SSL configures OpenSSL
	// +optional
	SSL *SSL `json:"ssl,omitempty"`
	// This setting is only available when support for OpenSSL was built in. It
	// designates a list of PEM file with an optional ssl configuration and a SNI
	// filter per certificate.
	// +optional
	SSLCertificateList *CertificateList `json:"sslCertificateList,omitempty"`
	// Hidden hides the bind and prevent exposing the Bind in services or routes
	// +optional
	Hidden *bool `json:"hidden,omitempty"`
	// AcceptProxy enforces the use of the PROXY protocol over any connection accepted by any of
	// the sockets declared on the same line.
	// +optional
	AcceptProxy *bool `json:"acceptProxy,omitempty"`
}

func (b *Bind) Model() (models.Bind, error) {
	model := models.Bind{
		Address:      b.Address,
		Port:         ptr.To(b.Port),
		PortRangeEnd: b.PortRangeEnd,
		BindParams: models.BindParams{
			Name:        b.Name,
			Transparent: b.Transparent,
			AcceptProxy: ptr.Deref(b.AcceptProxy, false),
		},
	}

	if b.SSL != nil && b.SSL.Enabled {
		model.Ssl = b.SSL.Enabled
		model.Verify = b.SSL.Verify

		if b.SSLCertificateList != nil {
			model.CrtList = b.SSLCertificateList.FilePath()
		}

		if b.SSL.Certificate != nil {
			model.SslCertificate = b.SSL.Certificate.FilePath()
		}

		if b.SSL.CACertificate != nil {
			model.SslCafile = b.SSL.CACertificate.FilePath()
		}

		if b.SSL.MinVersion != "" {
			model.SslMinVer = b.SSL.MinVersion
		}
	}

	return model, model.Validate(strfmt.Default)
}

type ServerParams struct {
	// SSL configures OpenSSL
	// +optional
	SSL *SSL `json:"ssl,omitempty"`
	// Weight parameter is used to adjust the server weight relative to
	// other servers. All servers will receive a load proportional to their weight
	// relative to the sum of all weights.
	// +kubebuilder:validation:Maximum=256
	// +kubebuilder:validation:Minimum=0
	Weight *int64 `json:"weight,omitempty"`
	// Check configures the health checks of the server.
	// +optional
	Check *Check `json:"check,omitempty"`
	// InitAddr indicates in what order the server address should be resolved upon startup if it uses an FQDN.
	// Attempts are made to resolve the address by applying in turn each of the methods mentioned in the comma-delimited
	// list. The first method which succeeds is used.
	// +optional
	InitAddr *string `json:"initAddr,omitempty"`
	// Resolvers points to an existing resolvers to resolve current server hostname.
	// +optional
	Resolvers *corev1.LocalObjectReference `json:"resolvers,omitempty"`
	// SendProxy enforces use of the PROXY protocol over any
	// connection established to this server. The PROXY protocol informs the other
	// end about the layer 3/4 addresses of the incoming connection, so that it can
	// know the client address or the public address it accessed to, whatever the
	// upper layer protocol.
	// +optional
	SendProxy *bool `json:"sendProxy,omitempty"`
	// SendProxyV2 preparing new update.
	SendProxyV2 *ProxyProtocol `json:"SendProxyV2,omitempty"`
	// VerifyHost is only available when support for OpenSSL was built in, and
	// only takes effect if pec.ssl.verify' is set to 'required'. This directive sets
	// a default static hostname to check the server certificate against when no
	// SNI was used to connect to the server.
	// +optional
	VerifyHost string `json:"verifyHost,omitempty"`
	// Cookie sets the cookie value assigned to the server.
	// +optional
	Cookie bool `json:"cookie,omitempty"`
}

type ServerTemplate struct {
	ServerParams `json:",inline"`
	// Prefix for the server names to be built.
	// +kubebuilder:validation:Pattern=^[^\s]+$
	Prefix string `json:"prefix"`
	// NumMin is the min number of servers as server name suffixes this template initializes.
	// +optional
	NumMin *int64 `json:"numMin,omitempty"`
	// Num is the max number of servers as server name suffixes this template initializes.
	Num int64 `json:"num"`
	// FQDN for all the servers this template initializes.
	FQDN string `json:"fqdn"`
	// Port
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=1
	Port int64 `json:"port"`
}

func (s *ServerTemplate) Model() (models.ServerTemplate, error) {
	model := models.ServerTemplate{
		ServerParams: models.ServerParams{
			Weight:   s.Weight,
			InitAddr: s.InitAddr,
		},
		Fqdn:   s.FQDN,
		Port:   ptr.To(s.Port),
		Prefix: s.Prefix,
	}

	if s.NumMin != nil {
		model.NumOrRange = fmt.Sprintf("%d-%d", *s.NumMin, s.Num)
	} else {
		model.NumOrRange = strconv.Itoa(int(s.Num))
	}

	if ptr.Deref(s.SendProxy, false) {
		model.SendProxy = models.ServerParamsSendProxyEnabled
	}

	if s.SendProxyV2 != nil {
		if s.SendProxyV2.V1 && s.SendProxyV2.V2 == nil && !s.SendProxyV2.V2SSL && !s.SendProxyV2.V2SSLCN {
			model.SendProxy = models.ServerParamsSendProxyEnabled
		} else if !s.SendProxyV2.V1 && s.SendProxyV2.V2 != nil && s.SendProxyV2.V2.Enabled && !s.SendProxyV2.V2SSL && !s.SendProxyV2.V2SSLCN {
			model.SendProxyV2 = models.ServerParamsSendProxyV2Enabled
			setProxyProtocolV2ServerTemplate(&model, s)
		} else if !s.SendProxyV2.V1 && s.SendProxyV2.V2 == nil && s.SendProxyV2.V2SSL && !s.SendProxyV2.V2SSLCN {
			model.SendProxyV2Ssl = models.ServerParamsSendProxyV2SslEnabled
		} else if !s.SendProxyV2.V1 && s.SendProxyV2.V2 == nil && !s.SendProxyV2.V2SSL && s.SendProxyV2.V2SSLCN {
			model.SendProxyV2SslCn = models.ServerParamsSendProxyV2SslCnEnabled
		} else {
			return model, fmt.Errorf("you can only select one proxy protocol")
		}
	}

	if s.SSL != nil && s.SSL.Enabled {
		model.Ssl = models.ServerParamsSslEnabled
		model.Verify = s.SSL.Verify

		if s.SSL.Certificate != nil {
			model.SslCertificate = s.SSL.Certificate.FilePath()
		}

		if s.SSL.CACertificate == nil {
			model.Verify = "none"
		} else {
			model.SslCafile = s.SSL.CACertificate.FilePath()
		}

		if s.SSL.MinVersion != "" {
			model.SslMinVer = s.SSL.MinVersion
		}

		if s.SSL.SNI != "" {
			model.Sni = s.SSL.SNI
		}
	}

	if s.Check != nil && s.Check.Enabled {
		model.Check = models.ServerParamsCheckEnabled

		if s.Check.Inter != nil {
			model.Inter = ptr.To(s.Check.Inter.Milliseconds())
		}

		model.Rise = s.Check.Rise
		model.Fall = s.Check.Fall
	}

	if s.Resolvers != nil {
		model.Resolvers = s.Resolvers.Name
	}

	return model, model.Validate(strfmt.Default)
}

func setProxyProtocolV2ServerTemplate(model *models.ServerTemplate, s *ServerTemplate) {
	if s.SendProxyV2.V2.Options != nil {
		var v2Options []string
		if s.SendProxyV2.V2.Options.Ssl {
			v2Options = append(v2Options, "ssl")
		}
		if s.SendProxyV2.V2.Options.CertCn {
			v2Options = append(v2Options, "cert-cn")
		}
		if s.SendProxyV2.V2.Options.SslCipher {
			v2Options = append(v2Options, "ssl-cipher")
		}
		if s.SendProxyV2.V2.Options.CertSig {
			v2Options = append(v2Options, "cert-sig")
		}
		if s.SendProxyV2.V2.Options.CertKey {
			v2Options = append(v2Options, "cert-key")
		}
		if s.SendProxyV2.V2.Options.Authority {
			v2Options = append(v2Options, "authority")
		}
		if s.SendProxyV2.V2.Options.Crc32c {
			v2Options = append(v2Options, "crc32c")
		}
		if s.SendProxyV2.V2.Options.UniqueID {
			v2Options = append(v2Options, "unique-id")
		}
		model.ProxyV2Options = v2Options
	}
}

type Server struct {
	ServerParams `json:",inline"`
	// Name of the server.
	Name string `json:"name"`
	// Address can be a host name, an IPv4 address, an IPv6 address.
	// +kubebuilder:validation:Pattern=^[^\s]+$
	Address string `json:"address"`
	// Port
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=1
	Port int64 `json:"port"`
}

func (s *Server) Model() (models.Server, error) {
	model := models.Server{
		ServerParams: models.ServerParams{
			Weight:     s.Weight,
			InitAddr:   s.InitAddr,
			Verifyhost: s.VerifyHost,
		},
		Name:    s.Name,
		Address: s.Address,
		Port:    ptr.To(s.Port),
	}

	if s.Cookie {
		model.ServerParams.Cookie = hash.GetMD5Hash(model.Address + ":" + strconv.Itoa(int(*model.Port)))
	}

	if ptr.Deref(s.SendProxy, false) {
		model.SendProxy = models.ServerParamsSendProxyEnabled
	}

	if s.SendProxyV2 != nil {
		if s.SendProxyV2.V1 && s.SendProxyV2.V2 == nil && !s.SendProxyV2.V2SSL && !s.SendProxyV2.V2SSLCN {
			model.SendProxy = models.ServerParamsSendProxyEnabled
		} else if !s.SendProxyV2.V1 && s.SendProxyV2.V2 != nil && s.SendProxyV2.V2.Enabled && !s.SendProxyV2.V2SSL && !s.SendProxyV2.V2SSLCN {
			model.SendProxyV2 = models.ServerParamsSendProxyV2Enabled
			setProxyProtocolV2Server(&model, s)
		} else if !s.SendProxyV2.V1 && s.SendProxyV2.V2 == nil && s.SendProxyV2.V2SSL && !s.SendProxyV2.V2SSLCN {
			model.SendProxyV2Ssl = models.ServerParamsSendProxyV2SslEnabled
		} else if !s.SendProxyV2.V1 && s.SendProxyV2.V2 == nil && !s.SendProxyV2.V2SSL && s.SendProxyV2.V2SSLCN {
			model.SendProxyV2SslCn = models.ServerParamsSendProxyV2SslCnEnabled
		} else {
			return model, fmt.Errorf("you can only select one proxy protocol")
		}
	}

	if s.SSL != nil && s.SSL.Enabled {
		model.Ssl = models.ServerParamsSslEnabled

		if s.SSL.Certificate != nil {
			model.SslCertificate = s.SSL.Certificate.FilePath()
		}

		if s.SSL.CACertificate == nil {
			model.Verify = "none"
		} else {
			model.SslCafile = s.SSL.CACertificate.FilePath()
		}

		if s.SSL.MinVersion != "" {
			model.SslMinVer = s.SSL.MinVersion
		}

		if s.SSL.SNI != "" {
			model.Sni = s.SSL.SNI
		}
	}

	if s.Check != nil && s.Check.Enabled {
		model.Check = models.ServerParamsCheckEnabled

		if s.Check.Inter != nil {
			model.Inter = ptr.To(s.Check.Inter.Milliseconds())
		}

		model.Rise = s.Check.Rise
		model.Fall = s.Check.Fall
	}

	if s.Resolvers != nil {
		model.Resolvers = s.Resolvers.Name
	}

	return model, model.Validate(strfmt.Default)
}

func setProxyProtocolV2Server(model *models.Server, s *Server) {
	if s.SendProxyV2.V2.Options != nil {
		var v2Options []string
		if s.SendProxyV2.V2.Options.Ssl {
			v2Options = append(v2Options, "ssl")
		}
		if s.SendProxyV2.V2.Options.CertCn {
			v2Options = append(v2Options, "cert-cn")
		}
		if s.SendProxyV2.V2.Options.SslCipher {
			v2Options = append(v2Options, "ssl-cipher")
		}
		if s.SendProxyV2.V2.Options.CertSig {
			v2Options = append(v2Options, "cert-sig")
		}
		if s.SendProxyV2.V2.Options.CertKey {
			v2Options = append(v2Options, "cert-key")
		}
		if s.SendProxyV2.V2.Options.Authority {
			v2Options = append(v2Options, "authority")
		}
		if s.SendProxyV2.V2.Options.Crc32c {
			v2Options = append(v2Options, "crc32c")
		}
		if s.SendProxyV2.V2.Options.UniqueID {
			v2Options = append(v2Options, "unique-id")
		}
		model.ProxyV2Options = v2Options
	}
}

type Check struct {
	// Enable enables health checks on a server. If not set, no health checking is performed, and the server is always
	// considered available.
	Enabled bool `json:"enabled"`
	// Inter sets the interval between two consecutive health checks. If left unspecified, the delay defaults to 2000 ms.
	// +optional
	Inter *metav1.Duration `json:"inter,omitempty"`
	// Rise specifies the number of consecutive successful health checks after a server will be considered as operational.
	// This value defaults to 2 if unspecified.
	// +optional
	Rise *int64 `json:"rise,omitempty"`
	// Fall specifies the number of consecutive unsuccessful health checks after a server will be considered as dead.
	// This value defaults to 3 if unspecified.
	// +optional
	Fall *int64 `json:"fall,omitempty"`
}

type Balance struct {
	// Algorithm is the algorithm used to select a server when doing load balancing. This only applies when no persistence information is available, or when a connection is redispatched to another server.
	// +kubebuilder:validation:Enum=roundrobin;static-rr;leastconn;first;source;uri;hdr;random;rdp-cookie
	Algorithm string `json:"algorithm"`
}

func (b *Balance) Model() (models.Balance, error) {
	model := models.Balance{
		Algorithm: ptr.To(b.Algorithm),
	}

	return model, model.Validate(strfmt.Default)
}

type SSL struct {
	// Enabled enables SSL deciphering on connections instantiated from this listener. A
	// certificate is necessary. All contents in the buffers will
	// appear in clear text, so that ACLs and HTTP processing will only have access
	// to deciphered contents. SSLv3 is disabled per default, set MinVersion to SSLv3
	// to enable it.
	Enabled bool `json:"enabled"`
	// MinVersion enforces use of the specified version or upper on SSL connections
	// instantiated from this listener.
	// +kubebuilder:validation:Enum=SSLv3;TLSv1.0;TLSv1.1;TLSv1.2;TLSv1.3
	// +optional
	MinVersion string `json:"minVersion,omitempty"`
	// Verify is only available when support for OpenSSL was built in. If set
	// to 'none', client certificate is not requested. This is the default. In other
	// cases, a client certificate is requested. If the client does not provide a
	// certificate after the request and if 'Verify' is set to 'required', then the
	// handshake is aborted, while it would have succeeded if set to 'optional'. The verification
	// of the certificate provided by the client using CAs from CACertificate.
	// On verify failure the handshake abortes, regardless of the 'verify' option.
	// +kubebuilder:validation:Enum=none;optional;required
	// +optional
	Verify string `json:"verify,omitempty"`
	// CACertificate configures the CACertificate used for the Server or Bind client certificate
	// +optional
	CACertificate *SSLCertificate `json:"caCertificate,omitempty"`
	// Certificate configures a PEM based Certificate file containing both the required certificates and any
	// associated private keys.
	// +optional
	Certificate *SSLCertificate `json:"certificate,omitempty"`
	// SNI parameter evaluates the sample fetch expression, converts it to a
	// string and uses the result as the host name sent in the SNI TLS extension to
	// the server.
	// +optional
	SNI string `json:"sni,omitempty"`
	// Alpn enables the TLS ALPN extension and advertises the specified protocol
	// list as supported on top of ALPN.
	// +optional
	Alpn []string `json:"alpn,omitempty"`
}

type SSLCertificate struct {
	Name      string                    `json:"name"`
	Value     *string                   `json:"value,omitempty"`
	ValueFrom []SSLCertificateValueFrom `json:"valueFrom,omitempty"`
}

func (s *SSLCertificate) FilePath() string {
	return fmt.Sprintf("/usr/local/etc/haproxy/%s.crt", strings.TrimSuffix(s.Name, ".crt"))
}

type SSLCertificateValueFrom struct {
	// ConfigMapKeyRef selects a key of a ConfigMap
	// +optional
	ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	// SecretKeyRef selects a key of a secret in the pod namespace
	// +optional
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

type CertificateListElement struct {
	// Certificate that will be presented to clients who provide a valid
	// TLSServerNameIndication field matching the SNIFilter.
	Certificate SSLCertificate `json:"certificate"`
	// SNIFilter specifies the filter for the SSL Certificate.  Wildcards are supported in the SNIFilter. Negative filter are also supported.
	SNIFilter string `json:"sniFilter"`
	// Alpn enables the TLS ALPN extension and advertises the specified protocol
	// list as supported on top of ALPN.
	// +optional
	Alpn []string `json:"alpn,omitempty"`
}

type CertificateList struct {
	// Name is the name of the certificate list
	Name string `json:"name"`
	// Elements is a list of SSL configuration and a SNI filter per certificate. If backend switching based on regex is used the host certificate
	Elements []CertificateListElement `json:"elements,omitempty"`
	// LabelSelector to select multiple backend certificates
	LabelSelector *metav1.LabelSelector `json:"selector,omitempty"`
}

func (r *CertificateList) FilePath() string {
	return fmt.Sprintf("/usr/local/etc/haproxy/%s.map", strings.TrimSuffix(r.Name, ".map"))
}

type HTTPRequestRules struct {
	// SetHeader sets HTTP header fields
	SetHeader []HTTPHeaderRule `json:"setHeader,omitempty"`
	// SetPath sets request path
	SetPath []HTTPPathRule `json:"setPath,omitempty"`
	// AddHeader appends HTTP header fields
	AddHeader []HTTPHeaderRule `json:"addHeader,omitempty"`
	// Redirect performs an HTTP redirection based on a redirect rule.
	// +optional
	Redirect []Redirect `json:"redirect,omitempty"`
	// ReplacePath matches the value of the path using a regex and completely replaces it with the specified format.
	// The replacement does not modify the scheme, the authority and the query-string.
	// +optional
	ReplacePath []ReplacePath `json:"replacePath,omitempty"`
	// Deny stops the evaluation of the rules and immediately rejects the request and emits an HTTP 403 error.
	// Optionally the status code specified as an argument to deny_status.
	// +optional
	Deny *Deny `json:"deny,omitempty"`
	// DenyStatus is the HTTP status code.
	// +kubebuilder:validation:Minimum=200
	// +kubebuilder:validation:Maximum=599
	// +optional
	DenyStatus *int64 `json:"denyStatus,omitempty"`
	// Return stops the evaluation of the rules and immediately returns a response.
	Return *HTTPReturn `json:"return,omitempty"`
}

func (h *HTTPRequestRules) Model() (models.HTTPRequestRules, error) {
	model := models.HTTPRequestRules{}

	for idx, header := range h.SetHeader {
		model = append(model, &models.HTTPRequestRule{
			Type:      "set-header",
			Index:     ptr.To(int64(idx)),
			HdrName:   header.Name,
			HdrFormat: header.Value.String(),
			Cond:      header.ConditionType,
			CondTest:  header.Condition,
		})
	}

	for idx, path := range h.SetPath {
		model = append(model, &models.HTTPRequestRule{
			Type:     "set-path",
			Index:    ptr.To(int64(idx)),
			PathFmt:  path.Value,
			Cond:     path.ConditionType,
			CondTest: path.Condition,
		})
	}

	for idx, header := range h.AddHeader {
		model = append(model, &models.HTTPRequestRule{
			Type:      "add-header",
			Index:     ptr.To(int64(idx)),
			HdrName:   header.Name,
			HdrFormat: header.Value.String(),
			Cond:      header.ConditionType,
			CondTest:  header.Condition,
		})
	}

	for idx, header := range h.ReplacePath {
		model = append(model, &models.HTTPRequestRule{
			Type:       "replace-path",
			Index:      ptr.To(int64(idx)),
			PathMatch:  header.MatchRegex,
			PathFmt:    header.ReplaceFmt,
			Cond:       header.ConditionType,
			CondTest:   header.Condition,
		})
	}

	if h.Deny != nil && h.Deny.Enabled {
		model = append(model, &models.HTTPRequestRule{
			DenyStatus: h.DenyStatus,
			Index:      ptr.To(int64(0)),
			Type:       "deny",
			Cond:       h.Deny.ConditionType,
			CondTest:   h.Deny.Condition,
		})
	}

	for idx, redirect := range h.Redirect {
		redirectRule := &models.HTTPRequestRule{
			Cond:       redirect.ConditionType,
			CondTest:   redirect.Condition,
			Index:      ptr.To(int64(idx)),
			RedirCode:  redirect.Code,
			RedirValue: redirect.Value,
			Type:       "redirect",
		}

		switch redirect.Type {
		case RedirectType{Location: true}:
			redirectRule.RedirType = models.HTTPRequestRuleRedirTypeLocation
		case RedirectType{Prefix: true}:
			redirectRule.RedirType = models.HTTPRequestRuleRedirTypePrefix
		case RedirectType{Scheme: true}:
			redirectRule.RedirType = models.HTTPRequestRuleRedirTypeScheme
		case RedirectType{}:
			redirectRule.RedirType = ""
		default:
			return models.HTTPRequestRules{}, fmt.Errorf("you can only select one redirect type")
		}

		if redirect.Option != nil {
			option := ""
			if redirect.Option.DropQuery {
				option = option + HTTPRequestRuleRedirectOptionDropQuery + " "
			}
			if redirect.Option.AppendSlash {
				option = option + HTTPRequestRuleRedirectOptionAppendSlash + " "
			}
			if redirect.Option.SetCookie != nil {
				option = option + HTTPRequestRuleRedirectOptionSetCookie + " " + strings.ToUpper(redirect.Option.SetCookie.Name)
				if redirect.Option.SetCookie.Value != "" {
					option = option + redirect.Option.SetCookie.Value + " "
				} else {
					option = option + " "
				}
			}
			if redirect.Option.ClearCookie != nil {
				option = option + HTTPRequestRuleRedirectOptionClearCookie + " " + strings.ToUpper(redirect.Option.ClearCookie.Name)
				if redirect.Option.ClearCookie.Value == "=" {
					option = option + redirect.Option.ClearCookie.Value
				}
			}
			option = strings.TrimSuffix(option, " ")
			redirectRule.RedirOption = option
		}
		model = append(model, redirectRule)
	}

	if h.Return != nil {
		value := h.Return.Content.Value
		if strings.Contains(h.Return.Content.Format, "string") {
			value = fmt.Sprintf("\"%s\"", h.Return.Content.Value)
		}

		model = append(model, &models.HTTPRequestRule{
			Type:                "return",
			ReturnStatusCode:    h.Return.Status,
			ReturnContentType:   ptr.To(h.Return.Content.Type),
			ReturnContentFormat: h.Return.Content.Format,
			ReturnContent:       value,
		})
	}

	for i := 0; i < len(model); i++ {
		model[i].Index = ptr.To(int64(i))
	}

	return model, model.Validate(strfmt.Default)
}

type HTTPReturn struct {
	// Status can be optionally specified, the default status code used for the response is 200.
	// +kubebuilder:default=200
	Status *int64 `json:"status,omitempty"`
	// Content is a full HTTP response specifying the errorfile to use, or the response payload specifying the file or the string to use.
	Content HTTPReturnContent `json:"content"`
}

type HTTPReturnContent struct {
	// Type specifies the content-type of the HTTP response.
	Type string `json:"type"`
	// ContentFormat defines the format of the Content. Can be one an errorfile or a string.
	// +kubebuilder:validation:Enum=default-errorfile;errorfile;errorfiles;file;lf-file;string;lf-string
	Format string `json:"format"`
	// Value specifying the file or the string to use.
	Value string `json:"value"`
}

type HTTPHeaderRule struct {
	Rule `json:",inline"`
	// Name specifies the header name
	Name string `json:"name"`
	// Value specifies the header value
	Value HTTPHeaderValue `json:"value"`
}

type HTTPPathRule struct {
	Rule `json:",inline"`
	// Value specifies the path value
	Value string `json:"format,omitempty"`
}

type HTTPHeaderValue struct {
	// Env variable with the header value
	Env *corev1.EnvVar `json:"env,omitempty"`
	// Str with the header value
	Str *string `json:"str,omitempty"`
	// Format specifies the format of the header value (implicit default is '%s')
	Format *string `json:"format,omitempty"`
}

type ReplacePath struct {
	Rule `json:",inline"`
	// MatchRegex is a string pattern used to identify the paths that need to be replaced.
	MatchRegex string `json:"matchRegex"`
	// ReplaceFmt defines the format string used to replace the values that match the pattern.
	ReplaceFmt string `json:"replaceFmt"`
}

func (h *HTTPHeaderValue) String() string {
	str := ptr.Deref(h.Str, "")
	if h.Env != nil {
		str = fmt.Sprintf("${%s}", h.Env.Name)
	}

	if h.Format != nil {
		return fmt.Sprintf(*h.Format, str)
	}

	return str
}

// Status defines the observed state of an object
type Status struct {
	// Phase is a simple, high-level summary of where the object is in its lifecycle.
	Phase StatusPhase `json:"phase"`
	// ObservedGeneration the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Error shows the actual error message if Phase is 'Error'.
	// +optional
	Error string `json:"error,omitempty"`
}

// StatusPhase is a label for the phase of an object at the current time.
type StatusPhase string

// These are the valid statuses of a listen configuration.
const (
	StatusPhaseActive        StatusPhase = "Active"
	StatusPhaseInternalError StatusPhase = "Error"
)

type Forwardfor struct {
	Enabled bool `json:"enabled"`
	// Pattern: ^[^\s]+$
	Except string `json:"except,omitempty"`
	// Pattern: ^[^\s]+$
	Header string `json:"header,omitempty"`
	Ifnone bool   `json:"ifnone,omitempty"`
}

type Deny struct {
	Rule `json:",inline"`
	// Enabled enables deny http request
	Enabled bool `json:"enabled"`
}

type Redirect struct {
	// +optional
	Rule `json:",inline"`
	// Code indicates which type of HTTP redirection is desired.
	// +kubebuilder:validation:Enum=301;302;303;307;308
	// +optional
	Code *int64 `json:"code,omitempty"`
	// Type selects a mode and value to redirect
	// +optional
	Type RedirectType `json:"type,omitempty"`
	// Value to redirect
	// +optional
	Value string `json:"value,omitempty"`
	// Value to redirect
	// +optional
	Option *RedirectOption `json:"option,omitempty"`
}

type RedirectType struct {
	// Location replaces the entire location of a URL.
	// +optional
	Location bool `json:"location"`
	// Prefix adds a prefix to the URL's location.
	// +optional
	Prefix bool `json:"insert"`
	// Scheme redirects to a different scheme.
	// +optional
	Scheme bool `json:"prefix"`
}

type RedirectOption struct {
	// DropQuery removes the query string from the original URL when performing the concatenation.
	// +optional
	DropQuery bool `json:"dropQuery,omitempty"`
	// AppendSlash adds a / character at the end of the URL.
	// +optional
	AppendSlash bool `json:"appendSlash,omitempty"`
	// SetCookie adds header to the redirection. It will be added with NAME (and optionally "=value")
	// +optional
	SetCookie *RedirectCookie `json:"SetCookie,omitempty"`
	// ClearCookie is to instruct the browser to delete the cookie. It will be added with NAME (and optionally "=").
	// To add "=" type any string in the value field
	// +optional
	ClearCookie *RedirectCookie `json:"ClearCookie,omitempty"`
}

type RedirectCookie struct {
	// Name
	// +optional
	Name string `json:"name,omitempty"`
	// Value
	// +optional
	Value string `json:"value,omitempty"`
}

const (
	HTTPRequestRuleRedirectOptionDropQuery   = "drop-query"
	HTTPRequestRuleRedirectOptionAppendSlash = "append-slash"
	HTTPRequestRuleRedirectOptionSetCookie   = "set-cookie"
	HTTPRequestRuleRedirectOptionClearCookie = "clear-cookie"
)

type HTTPPretendKeepalive struct {
	Enabled bool `json:"enabled"`
}

type Cookie struct {
	// Name of the cookie which will be monitored, modified or inserted in order to bring persistence.
	Name string `json:"name,omitempty"`
	// Mode could be 'rewrite', 'insert', 'prefix'. Select one.
	// +optional
	Mode CookieMode `json:"mode,omitempty"`
	// Indirect no cookie will be emitted to a client which already has a valid one
	// for the server which has processed the request.
	// +optional
	Indirect *bool `json:"indirect,omitempty"`
	// NoCache recommended in conjunction with the insert mode when there is a cache
	// between the client and HAProx
	// +optional
	NoCache *bool `json:"noCache,omitempty"`
	// PostOnly ensures that cookie insertion will only be performed on responses to POST requests.
	// +optional
	PostOnly *bool `json:"postOnly,omitempty"`
	// Preserve only be used with "insert" and/or "indirect". It allows the server
	// to emit the persistence cookie itself.
	// +optional
	Preserve *bool `json:"preserve,omitempty"`
	// HTTPOnly add an "HttpOnly" cookie attribute when a cookie is inserted.
	// It doesn't share the cookie with non-HTTP components.
	// +optional
	HTTPOnly *bool `json:"httpOnly,omitempty"`
	// Secure add a "Secure" cookie attribute when a cookie is inserted. The user agent
	// never emits this cookie over non-secure channels. The cookie will be presented
	// only over SSL/TLS connections.
	// +optional
	Secure *bool `json:"secure,omitempty"`
	// Dynamic activates dynamic cookies, when used, a session cookie is dynamically created for each server,
	// based on the IP and port of the server, and a secret key.
	// +optional
	Dynamic *bool `json:"dynamic,omitempty"`
	// Domain specify the domain at which a cookie is inserted. You can specify
	// several domain names by invoking this option multiple times.
	// +optional
	Domain []string `json:"domain,omitempty"`
	// MaxIdle cookies are ignored after some idle time.
	// +optional
	MaxIdle int64 `json:"maxIdle,omitempty"`
	// MaxLife cookies are ignored after some life time.
	// +optional
	MaxLife int64 `json:"maxLife,omitempty"`
	// Attribute add an extra attribute when a cookie is inserted.
	// +optional
	Attribute []string `json:"attribute,omitempty"`
}

type CookieMode struct {
	// Rewrite the cookie will be provided by the server.
	Rewrite bool `json:"rewrite"`
	// Insert cookie will have to be inserted by haproxy in server responses.
	Insert bool `json:"insert"`
	// Prefix is needed in some specific environments where the client does not support
	// more than one single cookie and the application already needs it.
	Prefix bool `json:"prefix"`
}

type ProxyProtocol struct {
	// V1 parameter enforces use of the PROXY protocol version 1.
	// +optional
	V1 bool `json:"v1"`
	// V2 parameter enforces use of the PROXY protocol version 2.
	// +optional
	V2 *ProxyProtocolV2 `json:"v2"`
	// V2SSL parameter add the SSL information extension of the PROXY protocol to the PROXY protocol header.
	// +optional
	V2SSL bool `json:"v2SSL"`
	// V2SSLCN parameter add the SSL information extension of the PROXY protocol to the PROXY protocol header and he SSL information extension
	// along with the Common Name from the subject of the client certificate (if any), is added to the PROXY protocol header.
	// +optional
	V2SSLCN bool `json:"v2SSLCN"`
}

type ProxyProtocolV2 struct {
	// Enabled enables the PROXY protocol version 2.
	// +optional
	Enabled bool `json:"enabled"`
	// Options is a list of options to add to the PROXY protocol header.
	// +optional
	Options *ProxyProtocolV2Options `json:"options,omitempty"`
}

type ProxyProtocolV2Options struct {
	// Ssl is equivalent to use V2SSL.
	// +optional
	Ssl bool `json:"ssl"`
	// CertCn is equivalent to use V2SSLCN.
	// +optional
	CertCn bool `json:"certCn"`
	// SslCipher is the name of the used cipher.
	// +optional
	SslCipher bool `json:"sslCipher"`
	// CertSig is the signature algorithm of the used certificate.
	// +optional
	CertSig bool `json:"certSig"`
	// CertKey is the key algorithm of the used certificate.
	// +optional
	CertKey bool `json:"certKey"`
	// Authority is the host name value passed by the client (only SNI from a TLS)
	// +optional
	Authority bool `json:"authority"`
	// Crc32c is the checksum of the PROXYv2 header.
	// +optional
	Crc32c bool `json:"crc32C"`
	// UniqueId sends a unique ID generated using the frontend's "unique-id-format" within the PROXYv2 header.
	// This unique-id is primarily meant for "mode tcp". It can lead to unexpected results in "mode http".
	// +optional
	UniqueID bool `json:"uniqueID"`
}
