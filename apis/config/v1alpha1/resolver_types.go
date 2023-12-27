package v1alpha1

import (
	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v4/configuration"
	"github.com/haproxytech/client-native/v4/models"
	parser "github.com/haproxytech/config-parser/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// ResolverSpec defines the desired state of Resolver
type ResolverSpec struct {
	// Nameservers used to configure a nameservers.
	Nameservers []Nameserver `json:"nameservers,omitempty"`
	// AcceptedPayloadSize defines the maximum payload size accepted by HAProxy and announced to all the  name servers
	// configured in this resolver.
	// +kubebuilder:validation:Maximum=8192
	// +kubebuilder:validation:Minimum=512
	// +optional
	AcceptedPayloadSize *int64 `json:"acceptedPayloadSize,omitempty"`
	// ParseResolvConf if true, adds all nameservers found in /etc/resolv.conf to this resolvers nameservers list.
	// +optional
	ParseResolvConf *bool `json:"parseResolvConf,omitempty"`
	// ResolveRetries defines the number <nb> of queries to send to resolve a server name before giving up. Default value: 3
	// +kubebuilder:validation:Minimum=1
	// +optional
	ResolveRetries *int64 `json:"resolveRetries,omitempty"`
	// Hold defines the period during which the last name resolution should be kept based on the last resolution status.
	// +optional
	Hold *Hold `json:"hold,omitempty"`
	// Timeouts defines timeouts related to name resolution.
	// +optional
	Timeouts *Timeouts `json:"timeouts,omitempty"`
}

type Nameserver struct {
	// Name specifies a unique name of the nameserver.
	// +kubebuilder:validation:Pattern="^[A-Za-z0-9-_.:]+$"
	Name string `json:"name"`
	// Address
	// +kubebuilder:validation:Pattern=^[^\s]+$
	Address string `json:"address"`
	// Port
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=1
	Port int64 `json:"port,omitempty"`
}

func (n *Nameserver) Model() (models.Nameserver, error) {
	model := models.Nameserver{
		Name:    n.Name,
		Address: pointer.String(n.Address),
		Port:    pointer.Int64(n.Port),
	}

	return model, model.Validate(strfmt.Default)
}

type Hold struct {
	// Nx defines interval between two successive name resolution when the last answer was nx.
	Nx *metav1.Duration `json:"nx,omitempty"`
	// Obsolete defines interval between two successive name resolution when the last answer was obsolete.
	Obsolete *metav1.Duration `json:"obsolete,omitempty"`
	// Other defines interval between two successive name resolution when the last answer was other.
	Other *metav1.Duration `json:"other,omitempty"`
	// Refused defines interval between two successive name resolution when the last answer was nx.
	Refused *metav1.Duration `json:"refused,omitempty"`
	// Timeout defines interval between two successive name resolution when the last answer was timeout.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// Valid defines interval between two successive name resolution when the last answer was valid.
	Valid *metav1.Duration `json:"valid,omitempty"`
}

type Timeouts struct {
	// Resolve time to trigger name resolutions when no other time applied. Default value: 1s
	// +optional
	Resolve *metav1.Duration `json:"resolve,omitempty"`
	// Retry time between two DNS queries, when no valid response have been received. Default value: 1s
	// +optional
	Retry *metav1.Duration `json:"retry,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name=Mode,type=string,JSONPath=`.spec.mode`
//+kubebuilder:printcolumn:name=Phase,type=string,JSONPath=`.status.phase`

// Resolver is the Schema for the Resolver API
type Resolver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResolverSpec `json:"spec,omitempty"`
	Status Status       `json:"status,omitempty"`
}

var _ Object = &Resolver{}

func (r *Resolver) SetStatus(status Status) {
	r.Status = status
}

func (r *Resolver) GetStatus() Status {
	return r.Status
}

func (r *Resolver) Model() (models.Resolver, error) {
	model := models.Resolver{
		Name:            r.Name,
		ParseResolvConf: pointer.BoolDeref(r.Spec.ParseResolvConf, false),
		ResolveRetries:  pointer.Int64Deref(r.Spec.ResolveRetries, 3),
		TimeoutResolve:  1000,
		TimeoutRetry:    1000,
	}

	if r.Spec.AcceptedPayloadSize != nil {
		model.AcceptedPayloadSize = *r.Spec.AcceptedPayloadSize
	}

	if r.Spec.Hold != nil {
		if r.Spec.Hold.Nx != nil {
			model.HoldNx = pointer.Int64(r.Spec.Hold.Nx.Milliseconds())
		}
		if r.Spec.Hold.Obsolete != nil {
			model.HoldObsolete = pointer.Int64(r.Spec.Hold.Obsolete.Milliseconds())
		}
		if r.Spec.Hold.Other != nil {
			model.HoldOther = pointer.Int64(r.Spec.Hold.Other.Milliseconds())
		}
		if r.Spec.Hold.Refused != nil {
			model.HoldRefused = pointer.Int64(r.Spec.Hold.Refused.Milliseconds())
		}
		if r.Spec.Hold.Timeout != nil {
			model.HoldTimeout = pointer.Int64(r.Spec.Hold.Timeout.Milliseconds())
		}
		if r.Spec.Hold.Valid != nil {
			model.HoldValid = pointer.Int64(r.Spec.Hold.Valid.Milliseconds())
		}
	}

	if r.Spec.Timeouts != nil {
		if r.Spec.Timeouts.Resolve != nil {
			model.TimeoutResolve = r.Spec.Timeouts.Resolve.Milliseconds()
		}
		if r.Spec.Timeouts.Retry != nil {
			model.TimeoutRetry = r.Spec.Timeouts.Retry.Milliseconds()
		}
	}

	return model, model.Validate(strfmt.Default)
}

func (r *Resolver) AddToParser(p parser.Parser) error {
	err := p.SectionsCreate(parser.Resolvers, r.Name)
	if err != nil {
		return err
	}

	var resolver models.Resolver
	resolver, err = r.Model()
	if err != nil {
		return err
	}

	if err := configuration.SerializeResolverSection(p, &resolver); err != nil {
		return err
	}

	for idx, nameserver := range r.Spec.Nameservers {
		model, err := nameserver.Model()
		if err != nil {
			return err
		}

		err = p.Insert(parser.Resolvers, r.Name, "nameserver", configuration.SerializeNameserver(model), idx)
		if err != nil {
			return err
		}
	}

	return nil
}

//+kubebuilder:object:root=true

// ResolverList contains a list of Resolver
type ResolverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Resolver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Resolver{}, &ResolverList{})
}
