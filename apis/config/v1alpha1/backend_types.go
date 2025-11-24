package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/six-group/haproxy-operator/pkg/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// BackendSpec defines the desired state of Backend
type BackendSpec struct {
	BaseSpec `json:",inline"`
	// CheckTimeout sets an additional check timeout, but only after a connection has been already
	// established.
	// +optional
	CheckTimeout *metav1.Duration `json:"checkTimeout,omitempty"`
	// Servers defines the backend servers and its configuration.
	Servers []Server `json:"servers,omitempty"`
	// ServerTemplates defines the backend server templates and its configuration.
	ServerTemplates []ServerTemplate `json:"serverTemplates,omitempty"`
	// Balance defines the load balancing algorithm to be used in a backend.
	// +optional
	Balance *Balance `json:"balance,omitempty"`
	// HostRegex specifies a regular expression used for backend switching rules.
	// +optional
	HostRegex string `json:"hostRegex,omitempty"`
	// HostCertificate specifies a certificate for that host used in the crt-list of a frontend
	// +optional
	HostCertificate *CertificateListElement `json:"hostCertificate,omitempty"`
	// Redispatch enable or disable session redistribution in case of connection failure
	// +optional
	Redispatch *bool `json:"redispatch,omitempty"`
	// HashType specifies a method to use for mapping hashes to servers
	// +optional
	HashType *HashType `json:"hashType,omitempty"`
	// Cookie enables cookie-based persistence in a backend.
	// +optional
	Cookie *Cookie `json:"cookie,omitempty"`
	// HTTPChk Enables HTTP protocol to check on the servers health
	// +optional
	HTTPChk *HTTPChk `json:"httpchk,omitempty"`
	// TCPCheck Perform health checks using tcp-check send/expect sequences
	// +optional
	TCPCheck *bool `json:"tcpCheck,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name=Mode,type=string,JSONPath=`.spec.mode`
//+kubebuilder:printcolumn:name=Phase,type=string,JSONPath=`.status.phase`

// Backend is the Schema for the backend API
type Backend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackendSpec `json:"spec,omitempty"`
	Status Status      `json:"status,omitempty"`
}

var _ Object = &Backend{}

func (b *Backend) SetStatus(status Status) {
	b.Status = status
}

func (b *Backend) GetStatus() Status {
	return b.Status
}

func (b *Backend) Model() (models.Backend, error) {
	model := models.Backend{
		BackendBase: models.BackendBase{
			Name: b.Name,
			Mode: b.Spec.Mode,
		},
	}

	if b.Spec.CheckTimeout != nil {
		model.CheckTimeout = ptr.To(b.Spec.CheckTimeout.Milliseconds())
	}

	if b.Spec.Forwardfor != nil {
		var enabled *string
		if b.Spec.Forwardfor.Enabled {
			enabled = ptr.To(models.ForwardforEnabledEnabled)
		}
		model.Forwardfor = &models.Forwardfor{
			Enabled: enabled,
			Except:  b.Spec.Forwardfor.Except,
			Header:  b.Spec.Forwardfor.Header,
			Ifnone:  b.Spec.Forwardfor.Ifnone,
		}
	}

	if b.Spec.HTTPChk != nil {
		model.AdvCheck = models.BackendBaseAdvCheckHttpchk
		model.HttpchkParams = &models.HttpchkParams{
			URI:    b.Spec.HTTPChk.URI,
			Method: b.Spec.HTTPChk.Method,
		}
	} else if b.Spec.TCPCheck != nil && *b.Spec.TCPCheck {
		model.AdvCheck = models.BackendBaseAdvCheckTCPDashCheck
	}

	if b.Spec.HTTPPretendKeepalive != nil && *b.Spec.HTTPPretendKeepalive {
		model.HTTPPretendKeepalive = models.BackendBaseHTTPPretendKeepaliveEnabled
	}

	if b.Spec.Redispatch != nil && *b.Spec.Redispatch {
		model.Redispatch = &models.Redispatch{
			Enabled:  ptr.To(models.RedispatchEnabledEnabled),
			Interval: ptr.To(int64(3)),
		}
	}

	if b.Spec.HashType != nil {
		ht, err := b.Spec.HashType.Model()
		if err == nil {
			model.HashType = ht
		}
	}

	if b.Spec.Balance != nil {
		model.Balance = &models.Balance{
			Algorithm: ptr.To(strings.ToLower(b.Spec.Balance.Algorithm)),
		}
	}

	if b.Spec.Cookie != nil {
		name := hash.GetMD5Hash(b.Spec.Cookie.Name)

		model.Cookie = &models.Cookie{
			Httponly: ptr.Deref(b.Spec.Cookie.HTTPOnly, false),
			Indirect: ptr.Deref(b.Spec.Cookie.Indirect, false),
			Maxidle:  b.Spec.Cookie.MaxIdle,
			Maxlife:  b.Spec.Cookie.MaxLife,
			Name:     ptr.To(name),
			Nocache:  ptr.Deref(b.Spec.Cookie.NoCache, false),
			Postonly: ptr.Deref(b.Spec.Cookie.PostOnly, false),
			Preserve: ptr.Deref(b.Spec.Cookie.Preserve, false),
			Secure:   ptr.Deref(b.Spec.Cookie.Secure, false),
			Dynamic:  ptr.Deref(b.Spec.Cookie.Dynamic, false),
		}

		for _, attr := range b.Spec.Cookie.Attribute {
			attrs := &models.Attr{Value: attr}
			model.Cookie.Attrs = append(model.Cookie.Attrs, attrs)
		}

		for _, domain := range b.Spec.Cookie.Domain {
			domains := &models.Domain{Value: domain}
			model.Cookie.Domains = append(model.Cookie.Domains, domains)
		}

		switch b.Spec.Cookie.Mode {
		case CookieMode{Rewrite: true}:
			model.Cookie.Type = models.CookieTypeRewrite
		case CookieMode{Insert: true}:
			model.Cookie.Type = models.CookieTypeInsert
		case CookieMode{Prefix: true}:
			model.Cookie.Type = models.CookieTypePrefix
		case CookieMode{}:
			model.Cookie.Type = ""
		default:
			return models.Backend{}, fmt.Errorf("you can only select one cookie mode")
		}

		if ptr.Deref(b.Spec.Cookie.Dynamic, false) {
			model.DynamicCookieKey = name
		}
	}

	for name, timeout := range b.Spec.Timeouts {
		switch name {
		case "check":
			model.CheckTimeout = ptr.To(timeout.Milliseconds())
		case "connect":
			model.ConnectTimeout = ptr.To(timeout.Milliseconds())
		case "http-keep-alive":
			model.HTTPKeepAliveTimeout = ptr.To(timeout.Milliseconds())
		case "http-request":
			model.HTTPRequestTimeout = ptr.To(timeout.Milliseconds())
		case "queue":
			model.QueueTimeout = ptr.To(timeout.Milliseconds())
		case "server":
			model.ServerTimeout = ptr.To(timeout.Milliseconds())
		case "tunnel":
			model.TunnelTimeout = ptr.To(timeout.Milliseconds())
		default:
			return model, fmt.Errorf("timeout %s unknown", name)
		}
	}

	for _, ef := range b.Spec.ErrorFiles {
		m, err := ef.Model()
		if err == nil {
			model.ErrorFiles = append(model.ErrorFiles, &m)
		}
	}

	return model, model.Validate(strfmt.Default)
}

func (b *Backend) AddToParser(p parser.Parser) error {
	err := p.SectionsCreate(parser.Backends, b.Name)
	if err != nil {
		return err
	}

	var backend models.Backend
	backend, err = b.Model()
	if err != nil {
		return err
	}

	configOpts := &options.ConfigurationOptions{}
	if err := configuration.CreateEditSection(&backend.BackendBase, parser.Backends, b.Name, p, configOpts); err != nil {
		return err
	}

	err = b.Spec.AddToParser(p, parser.Backends, b.Name)
	if err != nil {
		return err
	}

	for idx, server := range b.Spec.Servers {
		model, err := server.Model()

		if server.SSL != nil && server.SSL.Verify == "required" {
			model.Verify = server.SSL.Verify
			model.Alpn = strings.Join(server.SSL.Alpn, ",")
		}

		if err != nil {
			return err
		}

		err = p.Insert(parser.Backends, b.Name, "server", configuration.SerializeServer(model, configOpts), idx)
		if err != nil {
			return err
		}
	}

	for idx, template := range b.Spec.ServerTemplates {
		model, err := template.Model()
		if err != nil {
			return err
		}

		err = p.Insert(parser.Backends, b.Name, "server-template", configuration.SerializeServerTemplate(model, configOpts), idx)
		if err != nil {
			return err
		}
	}

	return nil
}

//+kubebuilder:object:root=true

// BackendList contains a list of Backend
type BackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backend{}, &BackendList{})
}
