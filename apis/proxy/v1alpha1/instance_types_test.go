package v1alpha1

import (
	"strings"
	"testing"
	"time"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/models"
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

func TestDefaultsConfigurationModelWithMaxconn(t *testing.T) {
	d := &DefaultsConfiguration{
		Mode:    "http",
		Maxconn: ptr.To(int64(2000)),
		Timeouts: map[string]metav1.Duration{
			"client":  {Duration: 5 * time.Second},
			"connect": {Duration: 5 * time.Second},
			"server":  {Duration: 10 * time.Second},
		},
	}

	model, err := d.Model()
	if err != nil {
		t.Fatalf("Model() returned error: %v", err)
	}

	if model.Maxconn == nil || *model.Maxconn != 2000 {
		t.Fatalf("expected maxconn 2000, got %#v", model.Maxconn)
	}
}

func TestDefaultsConfigurationAddToParserWithMaxconn(t *testing.T) {
	d := &DefaultsConfiguration{
		Mode:    "http",
		Maxconn: ptr.To(int64(2000)),
		Timeouts: map[string]metav1.Duration{
			"client":  {Duration: 5 * time.Second},
			"connect": {Duration: 5 * time.Second},
			"server":  {Duration: 10 * time.Second},
		},
	}

	p, err := parser.New()
	if err != nil {
		t.Fatalf("parser.New() returned error: %v", err)
	}

	if err := d.AddToParser(p); err != nil {
		t.Fatalf("AddToParser() returned error: %v", err)
	}

	if !strings.Contains(p.String(), "maxconn 2000") {
		t.Fatalf("expected generated config to contain maxconn 2000, got:\n%s", p.String())
	}
}

func TestDefaultsConfiguration_AddToParser_H1CaseAdjustBogusServerEnabled(t *testing.T) {
	p, err := parser.New()
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	d := &DefaultsConfiguration{
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

	d := &DefaultsConfiguration{
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

func TestDefaultsConfigurationModelWithOptions(t *testing.T) {
	d := &DefaultsConfiguration{
		Mode:     "http",
		Timeouts: map[string]metav1.Duration{},
		Options: &DefaultsOptions{
			LogSeparateErrors: ptr.To(true),
			LogHealthChecks:   ptr.To(true),
			Dontlognull:       ptr.To(false),
			DontlogNormal:     ptr.To(false),
			HTTPLogCLF:        ptr.To(true),
			Redispatch:        ptr.To(true),
		},
	}

	model, err := d.Model()
	if err != nil {
		t.Fatalf("Model() returned error: %v", err)
	}
	if model.LogSeparateErrors != models.DefaultsBaseLogSeparateErrorsEnabled {
		t.Fatalf("unexpected log-separate-errors value: %s", model.LogSeparateErrors)
	}
	if model.LogHealthChecks != models.DefaultsBaseLogHealthChecksEnabled {
		t.Fatalf("unexpected log-health-checks value: %s", model.LogHealthChecks)
	}
	if model.Dontlognull != models.DefaultsBaseDontlognullDisabled {
		t.Fatalf("unexpected dontlognull value: %s", model.Dontlognull)
	}
	if model.DontlogNormal != models.DefaultsBaseDontlogNormalDisabled {
		t.Fatalf("unexpected dontlog-normal value: %s", model.DontlogNormal)
	}
	if !model.Httplog {
		t.Fatalf("expected httplog to be enabled")
	}
	if !model.Clflog {
		t.Fatalf("expected clflog to be enabled")
	}
	if model.Redispatch == nil || model.Redispatch.Enabled == nil || *model.Redispatch.Enabled != models.RedispatchEnabledEnabled {
		t.Fatalf("expected redispatch to be enabled")
	}
}

func TestDefaultsConfigurationAddToParserWithOptions(t *testing.T) {
	d := &DefaultsConfiguration{
		Mode:     "http",
		Timeouts: map[string]metav1.Duration{},
		Options: &DefaultsOptions{
			LogSeparateErrors: ptr.To(true),
			LogHealthChecks:   ptr.To(true),
			Dontlognull:       ptr.To(false),
			DontlogNormal:     ptr.To(false),
			HTTPLogCLF:        ptr.To(true),
		},
	}

	p, err := parser.New()
	if err != nil {
		t.Fatalf("parser.New() returned error: %v", err)
	}

	if err := d.AddToParser(p); err != nil {
		t.Fatalf("AddToParser() returned error: %v", err)
	}

	cfg := p.String()
	checks := []string{
		"option httplog",
		"option log-separate-errors",
		"option log-health-checks",
		"no option dontlognull",
		"no option dontlog-normal",
	}

	for _, check := range checks {
		if !strings.Contains(cfg, check) {
			t.Fatalf("expected generated config to contain %q, got:\n%s", check, cfg)
		}
	}
}

func TestDefaultsConfigurationAddToParserWithHTTPLogCLF(t *testing.T) {
	d := &DefaultsConfiguration{
		Mode:     "http",
		Timeouts: map[string]metav1.Duration{},
		Options: &DefaultsOptions{
			HTTPLogCLF: ptr.To(true),
		},
	}

	p, err := parser.New()
	if err != nil {
		t.Fatalf("parser.New() returned error: %v", err)
	}

	if err := d.AddToParser(p); err != nil {
		t.Fatalf("AddToParser() returned error: %v", err)
	}

	cfg := p.String()
	if !strings.Contains(cfg, "option httplog clf") {
		t.Fatalf("expected generated config to contain %q, got:\n%s", "option httplog clf", cfg)
	}
}
