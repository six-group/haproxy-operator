package v1alpha1

import (
	"strings"
	"testing"
	"time"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestDefaultsConfigurationAddToParserWithRedispatchOption(t *testing.T) {
	p, err := parser.New()
	if err != nil {
		t.Fatalf("new parser: %v", err)
	}

	d := &DefaultsConfiguration{
		Mode: "http",
		Timeouts: map[string]metav1.Duration{
			"connect": {Duration: 5 * time.Second},
			"client":  {Duration: 50 * time.Second},
			"server":  {Duration: 50 * time.Second},
		},
		Options: &DefaultsOptions{
			Redispatch: ptr.To(true),
		},
	}

	if err := d.AddToParser(p); err != nil {
		t.Fatalf("add defaults to parser: %v", err)
	}

	cfg := p.String()
	if !strings.Contains(cfg, "option redispatch") {
		t.Fatalf("expected config to contain option redispatch, got:\n%s", cfg)
	}
	if strings.Contains(cfg, "option redispatch 3") {
		t.Fatalf("expected config not to contain redispatch interval, got:\n%s", cfg)
	}
}

func TestDefaultsConfigurationAddToParserWithRetries(t *testing.T) {
	p, err := parser.New()
	if err != nil {
		t.Fatalf("new parser: %v", err)
	}

	d := &DefaultsConfiguration{
		Mode: "http",
		Timeouts: map[string]metav1.Duration{
			"connect": {Duration: 5 * time.Second},
			"client":  {Duration: 10 * time.Second},
			"server":  {Duration: 10 * time.Second},
		},
		Retries: ptr.To(int64(5)),
	}

	if err := d.AddToParser(p); err != nil {
		t.Fatalf("add defaults to parser: %v", err)
	}

	cfg := p.String()
	if !strings.Contains(cfg, "retries 5") {
		t.Fatalf("expected config to contain retries 5, got:\n%s", cfg)
	}
}
