package v1alpha1_test

import (
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDefaultsConfiguration_AddToParser_H1CaseAdjustBogusServerEnabled(t *testing.T) {
	p, err := parser.New()
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	d := &proxyv1alpha1.DefaultsConfiguration{
		Mode:                    "http",
		Timeouts:                map[string]metav1.Duration{},
		H1CaseAdjustBogusServer: true,
	}

	if err := d.AddToParser(p); err != nil {
		t.Fatalf("failed to add defaults to parser: %v", err)
	}

	cfg := p.String()
	if !strings.Contains(cfg, "option h1-case-adjust-bogus-server") {
		t.Fatalf("expected h1-case-adjust-bogus-server option in defaults section, got:\n%s", cfg)
	}
}

func TestDefaultsConfiguration_AddToParser_H1CaseAdjustBogusServerDisabledByDefault(t *testing.T) {
	p, err := parser.New()
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	d := &proxyv1alpha1.DefaultsConfiguration{
		Mode:     "http",
		Timeouts: map[string]metav1.Duration{},
	}

	if err := d.AddToParser(p); err != nil {
		t.Fatalf("failed to add defaults to parser: %v", err)
	}

	cfg := p.String()
	if strings.Contains(cfg, "option h1-case-adjust-bogus-server") {
		t.Fatalf("did not expect h1-case-adjust-bogus-server option in defaults section, got:\n%s", cfg)
	}
}
