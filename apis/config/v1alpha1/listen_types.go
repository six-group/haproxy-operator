package v1alpha1

import (
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListenSpec defines the desired state of Listen
type ListenSpec struct {
	BaseSpec `json:",inline"`
	// Binds defines the frontend listening addresses, ports and its configuration.
	// +kubebuilder:validation:MinItems=1
	Binds []Bind `json:"binds"`
	// Servers defines the backend servers and its configuration.
	// +optional
	Servers []Server `json:"servers,omitempty"`
	// ServerTemplates defines the backend server templates and its configuration.
	// +optional
	ServerTemplates []ServerTemplate `json:"serverTemplates,omitempty"`
	// CheckTimeout sets an additional check timeout, but only after a connection has been already
	// established.
	// +optional
	CheckTimeout *metav1.Duration `json:"checkTimeout,omitempty"`
	// Balance defines the load balancing algorithm to be used in a backend.
	// +optional
	Balance *Balance `json:"balance,omitempty"`
	// Redispatch enable or disable session redistribution in case of connection failure
	// +optional
	Redispatch *bool `json:"redispatch,omitempty"`
	// HashType Specify a method to use for mapping hashes to servers
	// +optional
	HashType *HashType `json:"hashType,omitempty"`
	// Cookie enables cookie-based persistence in a backend.
	// +optional
	Cookie *Cookie `json:"cookie,omitempty"`
	// HostCertificate specifies a certificate for that host used in the crt-list of a frontend
	// +optional
	HostCertificate *CertificateListElement `json:"hostCertificate,omitempty"`
	// HTTPCheck Enables HTTP protocol to check on the servers health
	// +optional
	HTTPCheck *HTTPChk `json:"httpCheck,omitempty"`
	// TCPCheck Perform health checks using tcp-check send/expect sequences
	// +optional
	TCPCheck *bool `json:"tcpCheck,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name=Mode,type=string,JSONPath=`.spec.mode`
//+kubebuilder:printcolumn:name=Phase,type=string,JSONPath=`.status.phase`

// Listen is the Schema for the frontends API
type Listen struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ListenSpec `json:"spec,omitempty"`
	Status Status     `json:"status,omitempty"`
}

var _ Object = &Listen{}

func (l *Listen) SetStatus(status Status) {
	l.Status = status
}

func (l *Listen) GetStatus() Status {
	return l.Status
}

func (l *Listen) ToFrontend() *Frontend {
	frontend := Frontend{
		TypeMeta:   l.TypeMeta,
		ObjectMeta: l.ObjectMeta,
		Spec: FrontendSpec{
			BaseSpec: l.Spec.BaseSpec,
			Binds:    l.Spec.Binds,
			DefaultBackend: corev1.LocalObjectReference{
				Name: l.Name,
			},
		},
	}

	delete(frontend.Spec.Timeouts, "check")
	delete(frontend.Spec.Timeouts, "connect")
	delete(frontend.Spec.Timeouts, "queue")
	delete(frontend.Spec.Timeouts, "server")
	delete(frontend.Spec.Timeouts, "tunnel")

	return &frontend
}

func (l *Listen) ToBackend() *Backend {
	backend := Backend{
		TypeMeta:   l.TypeMeta,
		ObjectMeta: l.ObjectMeta,
		Spec: BackendSpec{
			BaseSpec:        l.Spec.BaseSpec,
			CheckTimeout:    l.Spec.CheckTimeout,
			Servers:         l.Spec.Servers,
			ServerTemplates: l.Spec.ServerTemplates,
			Balance:         l.Spec.Balance,
			Redispatch:      l.Spec.Redispatch,
			HashType:        l.Spec.HashType,
			Cookie:          l.Spec.Cookie,
			HostCertificate: l.Spec.HostCertificate,
			HTTPChk:         l.Spec.HTTPCheck,
			TCPCheck:        l.Spec.TCPCheck,
		},
	}

	delete(backend.Spec.Timeouts, "client")

	return &backend
}

func (l *Listen) AddToParser(p parser.Parser) error {
	return multierr.Combine(l.DeepCopy().ToFrontend().AddToParser(p), l.DeepCopy().ToBackend().AddToParser(p))
}

//+kubebuilder:object:root=true

// ListenList contains a list of Listen
type ListenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Listen `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Listen{}, &ListenList{})
}
