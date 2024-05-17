package v1alpha1

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/models"
	parser "github.com/haproxytech/config-parser/v5"
	"github.com/haproxytech/config-parser/v5/options"
	routev1 "github.com/openshift/api/route/v1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	// Replicas is the desired number of replicas of the HAProxy Instance.
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`
	// Network contains the configuration of Route, Services and other network related configuration.
	Network Network `json:"network"`
	// Configuration is used to bootstrap the global and defaults section of the HAProxy configuration.
	Configuration Configuration `json:"configuration"`
	// Image specifies the HaProxy image including th tag.
	// +kubebuilder:default="haproxy:latest"
	Image string `json:"image"`
	// Sidecars additional sidecar containers
	// +optional
	Sidecars []corev1.Container `json:"sidecars,omitempty"`
	// ServiceAccountName is the name of the ServiceAccount to use to run this Instance.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// ImagePullSecrets is an optional list of secret names in the same namespace to use for pulling any of the images used.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// AllowPrivilegedPorts allows to bind sockets with port numbers less than 1024.
	// +optional
	// +nullable
	AllowPrivilegedPorts *bool `json:"allowPrivilegedPorts,omitempty"`
	// Placement define how the instance's pods should be scheduled.
	// +optional
	// +nullable
	Placement *Placement `json:"placement,omitempty"`
	// ImagePullPolicy one of Always, Never, IfNotPresent.
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Metrics defines the metrics endpoint and scraping configuration.
	// +optional
	// +nullable
	Metrics *Metrics `json:"metrics,omitempty"`
	// +optional
	// +nullable
	// Labels additional labels for the ha-proxy pods
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	// +nullable
	// Env additional environment variables
	Env map[string]string `json:"env,omitempty"`
}

type Placement struct {
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// TopologySpreadConstraints describes how a group of pods ought to spread across topology
	// domains. Scheduler will schedule pods in a way which abides by the constraints.
	// +optional
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

type Network struct {
	// HostNetwork will enable the usage of host network.
	HostNetwork bool `json:"hostNetwork,omitempty"`
	// HostIPs defines an environment variable BIND_ADDRESS in the instance based on the provided host to IP mapping
	HostIPs map[string]string `json:"hostIPs,omitempty"`
	// Route defines the desired state for OpenShift Routes.
	Route RouteSpec `json:"route,omitempty"`
	// Service defines the desired state for a Service.
	Service ServiceSpec `json:"service,omitempty"`
}

type RouteSpec struct {
	// Enabled will toggle the creation of OpenShift Routes.
	Enabled bool `json:"enabled"`
	// TLS provides the ability to configure certificates and termination for the route.
	TLS *routev1.TLSConfig `json:"tls,omitempty"`
}

type ServiceSpec struct {
	// Enabled will toggle the creation of a Service.
	Enabled bool `json:"enabled"`
}

type Metrics struct {
	// Enabled will enable metrics globally for Instance.
	Enabled bool `json:"enabled"`
	// Address to bind the metrics endpoint (default: '0.0.0.0').
	// +optional
	// +kubebuilder:default="0.0.0.0"
	Address *string `json:"address,omitempty"`
	// Port specifies the port used for metrics.
	Port int64 `json:"port"`
	// RelabelConfigs to apply to samples before scraping.
	// More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config
	// +optional
	RelabelConfigs []*monitoringv1.RelabelConfig `json:"relabelings,omitempty"`
	// Interval at which metrics should be scraped
	// If not specified Prometheus' global scrape interval is used.
	// +optional
	Interval monitoringv1.Duration `json:"interval,omitempty"`
}

func (m *Metrics) AddToParser(p parser.Parser) error {
	if !m.Enabled {
		return nil
	}

	frontend := models.Frontend{
		Name: "metrics",
		Mode: "http",
		StatsOptions: &models.StatsOptions{
			StatsEnable:       true,
			StatsURIPrefix:    "/stats",
			StatsRefreshDelay: ptr.To((10 * time.Second).Milliseconds()),
		},
	}
	if err := p.SectionsCreate(parser.Frontends, frontend.Name); err != nil {
		return err
	}
	if err := configuration.CreateEditSection(frontend, parser.Frontends, frontend.Name, p); err != nil {
		return err
	}

	bind := models.Bind{
		BindParams: models.BindParams{
			Name: "metrics",
		},
		Port:    ptr.To(m.Port),
		Address: ptr.Deref(m.Address, "0.0.0.0"),
	}
	if err := p.Insert(parser.Frontends, frontend.Name, "bind", configuration.SerializeBind(bind), 0); err != nil {
		return err
	}

	rule := models.HTTPRequestRule{
		Type:        "use-service",
		ServiceName: "prometheus-exporter",
		Cond:        "if",
		CondTest:    "{ path /metrics }",
	}
	data, err := configuration.SerializeHTTPRequestRule(rule)
	if err != nil {
		return err
	}
	err = p.Insert(parser.Frontends, frontend.Name, "http-request", data, 0)
	if err != nil {
		return err
	}

	return nil
}

type Configuration struct {
	// Global contains the global HAProxy configuration settings
	Global GlobalConfiguration `json:"global"`
	// Defaults presets settings for all frontend, backend and listen
	Defaults DefaultsConfiguration `json:"defaults"`
	// LabelSelector to select other configuration objects of the config.haproxy.com API
	LabelSelector metav1.LabelSelector `json:"selector"`
}

type DefaultsLoggingConfiguration struct {
	// Enabled will enable logs for all proxies
	Enabled bool `json:"enabled"`
	// HTTPLog enables HTTP log format which is the most complete and the best suited for HTTP proxies. It provides
	// the same level of information as the TCP format with additional features which
	// are specific to the HTTP protocol.
	// +optional
	HTTPLog *bool `json:"httpLog,omitempty"`
	// TCPLog enables advanced logging of TCP connections with session state and timers. By default, the log output format
	// is very poor, as it only contains the source and destination addresses, and the instance name.
	// +optional
	TCPLog *bool `json:"tcpLog,omitempty"`
}

func (l *DefaultsLoggingConfiguration) Model() (models.LogTarget, error) {
	logTarget := models.LogTarget{
		Global: l.Enabled,
		Index:  ptr.To(int64(0)),
	}

	return logTarget, logTarget.Validate(strfmt.Default)
}

type GlobalConfiguration struct {
	// Reload enables auto-reload of the configuration using sockets. Requires an image that supports this feature.
	// +kubebuilder:default=false
	Reload bool `json:"reload"`
	// StatsTimeout sets the timeout on the stats socket. Default is set to 10 seconds.
	// +optional
	StatsTimeout *metav1.Duration `json:"statsTimeout,omitempty"`
	// Logging is used to enable and configure logging in the global section of the HAProxy configuration.
	// +optional
	Logging *GlobalLoggingConfiguration `json:"logging,omitempty"`
	// AdditionalParameters can be used to specify any further configuration statements which are not covered in this section explicitly.
	// +optional
	AdditionalParameters string `json:"additionalParameters,omitempty"`
	// AdditionalCertificates can be used to include global ssl certificates which can bes used in any listen
	// +optional
	AdditionalCertificates []configv1alpha1.SSLCertificate `json:"additionalCertificates,omitempty"`
	// Maxconn sets the maximum per-process number of concurrent connections. Proxies will stop accepting connections when this limit is reached.
	// +optional
	Maxconn *int64 `json:"maxconn,omitempty"`
	// Nbthread this setting is only available when support for threads was built in. It makes HAProxy run on specified number of threads.
	// +optional
	Nbthread *int64 `json:"nbthread,omitempty"`
	// TuneOptions sets the global tune options.
	// +optional
	TuneOptions *GlobalTuneOptions `json:"tune,omitempty"`
	// GlobalSSL sets the global SSL options.
	// +optional
	SSL *GlobalSSL `json:"ssl,omitempty"`
	// HardStopAfter is the maximum time the instance will remain alive when a soft-stop is received.
	// +optional
	HardStopAfter *time.Duration `json:"hardStopAfter,omitempty"`
}

func (g *GlobalConfiguration) Model() (models.Global, error) {
	global := models.Global{
		Maxconn:  ptr.Deref(g.Maxconn, 0),
		Nbthread: ptr.Deref(g.Nbthread, 0),
	}

	if g.AdditionalParameters != "" {
		str := strings.ReplaceAll(fmt.Sprintf("%s\n%s", parser.Global, g.AdditionalParameters), "\n", "\n  ")
		p, err := parser.New(options.String(str))
		if err != nil {
			return global, err
		}
		ptrr, err := configuration.ParseGlobalSection(p)
		if err != nil {
			return global, err
		}
		if ptrr != nil {
			global = *ptrr
		}
	}

	if g.StatsTimeout != nil {
		global.StatsTimeout = ptr.To(g.StatsTimeout.Milliseconds())
	}

	if g.Reload {
		global.RuntimeAPIs = append(global.RuntimeAPIs, &models.RuntimeAPI{
			Address: ptr.To("/var/lib/haproxy/run/haproxy.sock"),
			BindParams: models.BindParams{
				ExposeFdListeners: true,
				Level:             "admin",
				Mode:              "600",
			},
		})
	}

	if g.TuneOptions != nil {
		opts, err := g.TuneOptions.Model()
		if err != nil {
			return global, err
		}
		global.TuneOptions = &opts
	}

	if g.Logging != nil {
		_, logSendHostname, err := g.Logging.Model()
		if err != nil {
			return global, err
		}

		global.LogSendHostname = &logSendHostname
	}

	if g.SSL != nil {
		global.SslDefaultBindCiphers = strings.Join(g.SSL.DefaultBindCiphers, ":")
		global.SslDefaultBindCiphersuites = strings.Join(g.SSL.DefaultBindCipherSuites, ":")

		if g.SSL.DefaultBindOptions != nil {
			if g.SSL.DefaultBindOptions.MinVersion != nil {
				global.SslDefaultBindOptions += fmt.Sprintf("ssl-min-ver %s ", *g.SSL.DefaultBindOptions.MinVersion)
			}
		}
	}

	if g.HardStopAfter != nil {
		global.HardStopAfter = ptr.To(g.HardStopAfter.Milliseconds())
	}

	return global, global.Validate(strfmt.Default)
}

func (g *GlobalConfiguration) AddToParser(p parser.Parser) error {
	global, err := g.Model()
	if err != nil {
		return err
	}
	if err := configuration.SerializeGlobalSection(p, &global); err != nil {
		return err
	}

	if g.Logging != nil && g.Logging.Enabled {
		logTarget, _, err := g.Logging.Model()
		if err != nil {
			return err
		}
		if err := p.Insert(parser.Global, parser.GlobalSectionName, "log", configuration.SerializeLogTarget(logTarget), int(*logTarget.Index)); err != nil {
			return err
		}
	}

	return nil
}

type GlobalSSL struct {
	// DefaultBindCiphers sets the list of cipher algorithms ("cipher suite") that are negotiated during the SSL/TLS handshake up to TLSv1.2 for all
	// binds which do not explicitly define theirs.
	// +optional
	DefaultBindCiphers []string `json:"defaultBindCiphers,omitempty"`
	// DefaultBindCipherSuites sets the default list of cipher algorithms ("cipher suite") that are negotiated
	// during the TLSv1.3 handshake for all binds which do not explicitly define theirs.
	// +optional
	DefaultBindCipherSuites []string `json:"defaultBindCipherSuites,omitempty"`
	// DefaultBindOptions sets default ssl-options to force on all binds.
	// +optional
	DefaultBindOptions *GlobalSSLDefaultBindOptions `json:"defaultBindOptions,omitempty"`
}

type GlobalSSLDefaultBindOptions struct {
	// MinVersion enforces use of the specified version or upper on SSL connections
	// instantiated from this listener.
	// +kubebuilder:validation:Enum=SSLv3;TLSv1.0;TLSv1.1;TLSv1.2;TLSv1.3
	// +optional
	MinVersion *string `json:"minVersion,omitempty"`
}

type GlobalTuneOptions struct {
	// Maxrewrite sets the reserved buffer space to this size in bytes. The reserved space is
	// used for header rewriting or appending. The first reads on sockets will never
	// fill more than bufsize-maxrewrite.
	// +optional
	Maxrewrite *int64 `json:"maxrewrite,omitempty"`
	// Bufsize sets the buffer size to this size (in bytes). Lower values allow more
	// sessions to coexist in the same amount of RAM, and higher values allow some
	// applications with very large cookies to work.
	// +optional
	Bufsize *int64 `json:"bufsize,omitempty"`
	// SSL sets the SSL tune options.
	// +optional
	SSL *GlobalSSLTuneOptions `json:"ssl,omitempty"`
}

type GlobalSSLTuneOptions struct {
	// CacheSize sets the size of the global SSL session cache, in a number of blocks. A block
	// is large enough to contain an encoded session without peer certificate.  An
	// encoded session with peer certificate is stored in multiple blocks depending
	// on the size of the peer certificate. The default value may be forced
	// at build time, otherwise defaults to 20000.  Setting this value to 0 disables the SSL session cache.
	// +optional
	CacheSize *int64 `json:"cacheSize,omitempty"`
	// Keylog activates the logging of the TLS keys. It should be used with
	// care as it will consume more memory per SSL session and could decrease
	// performances. This is disabled by default.
	// +optional
	Keylog string `json:"keylog,omitempty"`
	// Lifetime sets how long a cached SSL session may remain valid. This time defaults to 5 min. It is important
	// to understand that it does not guarantee that sessions will last that long, because if the cache is
	// full, the longest idle sessions will be purged despite their configured lifetime.
	// +optional
	Lifetime *metav1.Duration `json:"lifetime,omitempty"`
	// ForcePrivateCache disables SSL session cache sharing between all processes. It
	// should normally not be used since it will force many renegotiations due to
	// clients hitting a random process.
	// +optional
	ForcePrivateCache bool `json:"forcePrivateCache,omitempty"`
	// MaxRecord sets the maximum amount of bytes passed to SSL_write() at a time. Default
	// value 0 means there is no limit. Over SSL/TLS, the client can decipher the
	// data only once it has received a full record.
	// +optional
	MaxRecord *int64 `json:"maxRecord,omitempty"`
	// DefaultDHParam sets the maximum size of the Diffie-Hellman parameters used for generating
	// the ephemeral/temporary Diffie-Hellman key in case of DHE key exchange. The
	// final size will try to match the size of the server's RSA (or DSA) key (e.g,
	// a 2048 bits temporary DH key for a 2048 bits RSA key), but will not exceed
	// this maximum value. Default value if 2048.
	// +optional
	DefaultDHParam int64 `json:"defaultDHParam,omitempty"`
	// CtxCacheSize sets the size of the cache used to store generated certificates to <number>
	// entries. This is an LRU cache. Because generating an SSL certificate
	// dynamically is expensive, they are cached. The default cache size is set to 1000 entries.
	// +optional
	CtxCacheSize int64 `json:"ctxCacheSize,omitempty"`
	// CaptureBufferSize sets the maximum size of the buffer used for capturing client hello cipher
	// list, extensions list, elliptic curves list and elliptic curve point
	// formats. If the value is 0 (default value) the capture is disabled,
	// otherwise a buffer is allocated for each SSL/TLS connection.
	// +optional
	CaptureBufferSize *int64 `json:"captureBufferSize,omitempty"`
}

func (t *GlobalTuneOptions) Model() (models.GlobalTuneOptions, error) {
	opts := models.GlobalTuneOptions{
		Maxrewrite: ptr.Deref(t.Maxrewrite, 0),
		Bufsize:    ptr.Deref(t.Bufsize, 0),
	}

	if t.SSL != nil {
		opts.SslCachesize = t.SSL.CacheSize
		opts.SslKeylog = t.SSL.Keylog
		opts.SslForcePrivateCache = t.SSL.ForcePrivateCache
		opts.SslMaxrecord = t.SSL.MaxRecord
		opts.SslDefaultDhParam = t.SSL.DefaultDHParam
		opts.SslCtxCacheSize = t.SSL.CtxCacheSize
		opts.SslCaptureBufferSize = t.SSL.CaptureBufferSize

		if t.SSL.Lifetime != nil {
			opts.SslLifetime = ptr.To(int64(math.Round(t.SSL.Lifetime.Seconds())))
		}
	}

	return opts, opts.Validate(strfmt.Default)
}

type GlobalLoggingConfiguration struct {
	// Enabled will toggle the creation of a global syslog server.
	Enabled bool `json:"enabled"`
	// Address can be a filesystem path to a UNIX domain socket or a remote syslog target (IPv4/IPv6 address optionally followed by a colon and a UDP port).
	// +kubebuilder:validation:Pattern=^[^\s]+$
	// +kubebuilder:default="/var/lib/rsyslog/rsyslog.sock"
	Address string `json:"address"`
	// Facility must be one of the 24 standard syslog facilities.
	// +kubebuilder:validation:Enum=kern;user;mail;daemon;auth;syslog;lpr;news;uucp;cron;auth2;ftp;ntp;audit;alert;cron2;local0;local1;local2;local3;local4;local5;local6;local7
	// +kubebuilder:default=local0
	Facility string `json:"facility,omitempty"`
	// Level can be specified to filter outgoing messages. By default, all messages are sent.
	// +kubebuilder:validation:Enum=emerg;alert;crit;err;warning;notice;info;debug
	// +optional
	Level string `json:"level,omitempty"`
	// Format is the log format used when generating syslog messages.
	// +kubebuilder:validation:Enum=rfc3164;rfc5424;short;raw
	// +optional
	Format string `json:"format,omitempty"`
	// SendHostname sets the hostname field in the syslog header.  Generally used if one is not relaying logs through an
	// intermediate syslog server.
	// +optional
	SendHostname *bool `json:"sendHostname,omitempty"`
	// Hostname specifies a value for the syslog hostname header, otherwise uses the hostname of the system.
	// +optional
	Hostname *string `json:"hostname,omitempty"`
}

func (l *GlobalLoggingConfiguration) Model() (models.LogTarget, models.GlobalLogSendHostname, error) {
	logTarget := models.LogTarget{
		Address:  l.Address,
		Level:    l.Level,
		Facility: l.Facility,
		Format:   l.Format,
		Index:    ptr.To(int64(0)),
	}

	logSendHostname := models.GlobalLogSendHostname{
		Enabled: ptr.To("disabled"),
	}
	if ptr.Deref(l.SendHostname, false) {
		logSendHostname.Enabled = ptr.To(models.GlobalLogSendHostnameEnabledEnabled)
		logSendHostname.Param = ptr.Deref(l.Hostname, "")
	}

	return logTarget, logSendHostname, multierr.Combine(logTarget.Validate(strfmt.Default), logSendHostname.Validate(strfmt.Default))
}

type DefaultsConfiguration struct {
	// Mode can be either 'tcp' or 'http'. In tcp mode it is a layer 4 proxy. In http mode it is a layer 7 proxy.
	// +kubebuilder:default=http
	// +kubebuilder:validation:Enum=http;tcp
	Mode string `json:"mode"`
	// ErrorFiles custom error files to be used
	// +optional
	ErrorFiles []*configv1alpha1.ErrorFile `json:"errorFiles,omitempty"`
	// Timeouts: check, client, client-fin, connect, http-keep-alive, http-request, queue, server, server-fin, tunnel.
	// The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit.
	// More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html
	// +kubebuilder:default={"client": "5s", "connect": "5s", "server": "10s"}
	Timeouts map[string]metav1.Duration `json:"timeouts"`
	// Logging is used to configure default logging for all proxies.
	// +optional
	Logging *DefaultsLoggingConfiguration `json:"logging,omitempty"`
	// AdditionalParameters can be used to specify any further configuration statements which are not covered in this section explicitly.
	// +optional
	AdditionalParameters string `json:"additionalParameters,omitempty"`
}

func (d *DefaultsConfiguration) Model() (models.Defaults, error) {
	defaults := models.Defaults{}

	if d.AdditionalParameters != "" {
		str := strings.ReplaceAll(fmt.Sprintf("%s\n%s", parser.Defaults, d.AdditionalParameters), "\n", "\n  ")
		p, err := parser.New(options.String(str))
		if err != nil {
			return defaults, err
		}
		if err = configuration.ParseSection(&defaults, parser.Defaults, parser.DefaultSectionName, p); err != nil {
			return defaults, err
		}
	}

	defaults.Mode = d.Mode

	for name, timeout := range d.Timeouts {
		switch name {
		case "check":
			defaults.CheckTimeout = ptr.To(timeout.Milliseconds())
		case "client":
			defaults.ClientTimeout = ptr.To(timeout.Milliseconds())
		case "client-fin":
			defaults.ClientFinTimeout = ptr.To(timeout.Milliseconds())
		case "connect":
			defaults.ConnectTimeout = ptr.To(timeout.Milliseconds())
		case "http-keep-alive":
			defaults.HTTPKeepAliveTimeout = ptr.To(timeout.Milliseconds())
		case "http-request":
			defaults.HTTPRequestTimeout = ptr.To(timeout.Milliseconds())
		case "queue":
			defaults.QueueTimeout = ptr.To(timeout.Milliseconds())
		case "server":
			defaults.ServerTimeout = ptr.To(timeout.Milliseconds())
		case "server-fin":
			defaults.ServerFinTimeout = ptr.To(timeout.Milliseconds())
		case "tunnel":
			defaults.TunnelTimeout = ptr.To(timeout.Milliseconds())
		default:
			return defaults, fmt.Errorf("timeout %s unknown", name)
		}
	}

	for _, ef := range d.ErrorFiles {
		model, err := ef.Model()
		if err != nil {
			return defaults, err
		}

		defaults.ErrorFiles = append(defaults.ErrorFiles, &model)
	}

	if d.Logging != nil {
		defaults.Httplog = ptr.Deref(d.Logging.HTTPLog, false)
		defaults.Tcplog = ptr.Deref(d.Logging.TCPLog, false)
	}

	return defaults, defaults.Validate(strfmt.Default)
}

func (d *DefaultsConfiguration) AddToParser(p parser.Parser) error {
	defaults, err := d.Model()
	if err != nil {
		return err
	}

	if err := p.SectionsCreate(parser.Defaults, parser.DefaultSectionName); err != nil {
		return err
	}

	if err := configuration.CreateEditSection(defaults, parser.Defaults, parser.DefaultSectionName, p); err != nil {
		return err
	}

	if d.Logging != nil && d.Logging.Enabled {
		logTarget, err := d.Logging.Model()
		if err != nil {
			return err
		}
		if err := p.Insert(parser.Defaults, parser.DefaultSectionName, "log", configuration.SerializeLogTarget(logTarget), int(*logTarget.Index)); err != nil {
			return err
		}
	}

	return nil
}

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	// Phase is a simple, high-level summary of where the Listen is in its lifecycle.
	Phase InstancePhase `json:"phase"`
	// Error shows the actual error message if Phase is 'Error'.
	// +optional
	Error string `json:"error,omitempty"`
}

// InstancePhase is a label for the phase of a Instance at the current time.
type InstancePhase string

// These are the valid statuses of a listen configuration.
const (
	InstancePhaseRunning       InstancePhase = "Running"
	InstancePhasePending       InstancePhase = "Pending"
	InstancePhaseInternalError InstancePhase = "Error"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Instance is the Schema for the instances API
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

func (i *Instance) AddToParser(p parser.Parser) error {
	if err := i.Spec.Configuration.Global.AddToParser(p); err != nil {
		return err
	}

	return i.Spec.Configuration.Defaults.AddToParser(p)
}

//+kubebuilder:object:root=true

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instance{}, &InstanceList{})
}
