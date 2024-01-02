package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v4/configuration"
	"github.com/haproxytech/client-native/v4/models"
	parser "github.com/haproxytech/config-parser/v4"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// FrontendSpec defines the desired state of Frontend
type FrontendSpec struct {
	BaseSpec `json:",inline"`
	// Binds defines the frontend listening addresses, ports and its configuration.
	// +kubebuilder:validation:MinItems=1
	Binds []Bind `json:"binds"`
	// BackendSwitching rules specify the specific backend used if/unless an ACL-based condition is matched.
	// +optional
	BackendSwitching []BackendSwitchingRule `json:"backendSwitching,omitempty"`
	// DefaultBackend to use when no 'use_backend' rule has been matched.
	DefaultBackend corev1.LocalObjectReference `json:"defaultBackend"`
}

type BackendSwitchingRule struct {
	Rule `json:",inline"`
	// Backend reference used to resolve the backend name.
	Backend BackendReference `json:"backend,omitempty"`
}

func (b *BackendSwitchingRule) Model() (models.BackendSwitchingRule, error) {
	model := models.BackendSwitchingRule{
		Cond:     b.ConditionType,
		CondTest: b.Condition,
		Index:    ptr.To(int64(0)),
		Name:     b.Backend.String(),
	}

	return model, model.Validate(strfmt.Default)
}

type BackendReference struct {
	// Name of a specific backend
	Name *string `json:"name,omitempty"`
	// Mapping of multiple backends
	RegexMapping *RegexBackendMapping `json:"regexMapping,omitempty"`
}

func (b *BackendReference) String() string {
	if b.RegexMapping != nil {
		return fmt.Sprintf("%%[%s,map_reg(%s)]", b.RegexMapping.Parameter, b.RegexMapping.FilePath())
	}

	return ptr.Deref(b.Name, "")
}

func (r *RegexBackendMapping) FoundCondition() string {
	return fmt.Sprintf("{ %s,map_reg(%s) -m found }", r.Parameter, r.FilePath())
}

type RegexBackendMapping struct {
	// Name to identify the mapping
	Name string `json:"name"`
	// Parameter which will be used for the mapping (default: base)
	// +kubebuilder:default=base
	Parameter string `json:"parameter"`
	// LabelSelector to select multiple backends
	LabelSelector metav1.LabelSelector `json:"selector"`
}

func (r *RegexBackendMapping) FilePath() string {
	return fmt.Sprintf("/usr/local/etc/haproxy/%s.map", strings.TrimSuffix(r.Name, ".map"))
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name=Mode,type=string,JSONPath=`.spec.mode`
//+kubebuilder:printcolumn:name=Phase,type=string,JSONPath=`.status.phase`

// Frontend is the Schema for the frontends API
type Frontend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FrontendSpec `json:"spec,omitempty"`
	Status Status       `json:"status,omitempty"`
}

var _ Object = &Frontend{}

func (f *Frontend) SetStatus(status Status) {
	f.Status = status
}

func (f *Frontend) GetStatus() Status {
	return f.Status
}

func (f *Frontend) Model() (models.Frontend, error) {
	model := models.Frontend{
		Name:           f.Name,
		Mode:           f.Spec.Mode,
		DefaultBackend: f.Spec.DefaultBackend.Name,
	}

	if f.Spec.Forwardfor != nil {
		var enabled *string
		if f.Spec.Forwardfor.Enabled {
			enabled = ptr.To(models.ForwardforEnabledEnabled)
		}
		model.Forwardfor = &models.Forwardfor{
			Enabled: enabled,
			Except:  f.Spec.Forwardfor.Except,
			Header:  f.Spec.Forwardfor.Header,
			Ifnone:  f.Spec.Forwardfor.Ifnone,
		}
	}

	for name, timeout := range f.Spec.Timeouts {
		switch name {
		case "client":
			model.ClientTimeout = ptr.To(timeout.Milliseconds())
		case "http-keep-alive":
			model.HTTPKeepAliveTimeout = ptr.To(timeout.Milliseconds())
		case "http-request":
			model.HTTPRequestTimeout = ptr.To(timeout.Milliseconds())
		default:
			return model, fmt.Errorf("timeout %s unknown", name)
		}
	}

	for _, ef := range f.Spec.ErrorFiles {
		m, err := ef.Model()
		if err == nil {
			model.ErrorFiles = append(model.ErrorFiles, &m)
		}
	}

	return model, model.Validate(strfmt.Default)
}

func (f *Frontend) AddToParser(p parser.Parser) error {
	err := p.SectionsCreate(parser.Frontends, f.Name)
	if err != nil {
		return err
	}

	var frontend models.Frontend
	frontend, err = f.Model()
	if err != nil {
		return err
	}

	if err := configuration.CreateEditSection(&frontend, parser.Frontends, f.Name, p); err != nil {
		return err
	}

	err = f.Spec.BaseSpec.AddToParser(p, parser.Frontends, f.Name)
	if err != nil {
		return err
	}

	for idx, bind := range f.Spec.Binds {
		model, err := bind.Model()
		if err != nil {
			return err
		}

		err = p.Insert(parser.Frontends, f.Name, "bind", configuration.SerializeBind(model), idx)
		if err != nil {
			return err
		}
	}

	for idx, rule := range f.Spec.BackendSwitching {
		model, err := rule.Model()
		if err != nil {
			return err
		}

		err = p.Insert(parser.Frontends, f.Name, "use_backend", configuration.SerializeBackendSwitchingRule(model), idx)
		if err != nil {
			return err
		}
	}

	return nil
}

//+kubebuilder:object:root=true

// FrontendList contains a list of Fronted
type FrontendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Frontend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Frontend{}, &FrontendList{})
}
