package main

import (
	"flag"
	"os"
	"strings"

	routev1 "github.com/openshift/api/route/v1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"github.com/six-group/haproxy-operator/controllers/config"
	"github.com/six-group/haproxy-operator/controllers/instance"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	crzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const envLeaderElect = "LEADER_ELECT"

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(configv1alpha1.AddToScheme(scheme))
	utilruntime.Must(proxyv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.Parse()

	setupLogging()

	setupLog.Info("starting operator")

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         strings.EqualFold(os.Getenv(envLeaderElect), "true"),
		LeaderElectionID:       "acc50d8e.haproxy.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Inspect cluster to verify availability of extra features and setup schema
	if err := instance.InspectCluster(); err != nil {
		setupLog.Info("unable to inspect cluster")
	}
	if instance.IsPrometheusAPIAvailable() {
		if err := monitoringv1.AddToScheme(mgr.GetScheme()); err != nil {
			setupLog.Error(err, "")
			os.Exit(1)
		}
	}
	if instance.IsRouteAPIAvailable() {
		if err := routev1.Install(mgr.GetScheme()); err != nil {
			setupLog.Error(err, "")
			os.Exit(1)
		}
	}

	if err = (&instance.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Instance")
		os.Exit(1)
	}
	if err = (&config.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Object: &configv1alpha1.Listen{},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Listen")
		os.Exit(1)
	}
	if err = (&config.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Object: &configv1alpha1.Frontend{},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Frontend")
		os.Exit(1)
	}
	if err = (&config.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Object: &configv1alpha1.Backend{},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Backend")
		os.Exit(1)
	}
	if err = (&config.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Object: &configv1alpha1.Resolver{},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Resolver")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupLogging() {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	opts := crzap.Options{
		Encoder: zapcore.NewJSONEncoder(encCfg),
	}
	logger := crzap.New(crzap.UseFlagOptions(&opts))
	ctrl.SetLogger(logger)
	// replace klog logger
	klog.SetLogger(ctrl.Log)
}
