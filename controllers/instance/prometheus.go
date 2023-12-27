package instance

import (
	"context"
	"fmt"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"github.com/six-group/haproxy-operator/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var prometheusAPIFound = false

func (r *Reconciler) reconcilePrometheusConfiguration(ctx context.Context, instance *proxyv1alpha1.Instance) error {
	if IsPrometheusAPIAvailable() {
		return r.reconcilePodMonitor(ctx, instance)
	}

	return nil
}

func (r *Reconciler) reconcilePodMonitor(ctx context.Context, instance *proxyv1alpha1.Instance) error {
	logger := log.FromContext(ctx)

	monitor := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetServiceName(instance),
			Namespace: instance.Namespace,
		},
	}

	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, monitor, func() error {
		if err := controllerutil.SetOwnerReference(instance, monitor, r.Scheme); err != nil {
			return err
		}

		monitor.Spec.Selector = metav1.LabelSelector{MatchLabels: utils.GetAppSelectorLabels(instance)}

		monitor.Spec.Endpoints = []monitoringv1.Endpoint{
			{
				Port:           "metrics",
				Path:           "/metrics",
				RelabelConfigs: instance.Spec.Metrics.RelabelConfigs,
				Interval:       instance.Spec.Metrics.Interval,
				Scheme:         "http",
			},
		}

		return nil
	})
	if err != nil {
		return err
	}
	if result != controllerutil.OperationResultNone {
		logger.Info(fmt.Sprintf("Object %s", result), "servicemonitor", monitor.Name)
	}

	return nil
}

// IsPrometheusAPIAvailable returns true if the Prometheus API is present.
func IsPrometheusAPIAvailable() bool {
	return prometheusAPIFound
}

// VerifyPrometheusAPI will verify that the Prometheus API is present.
func VerifyPrometheusAPI() error {
	found, err := utils.VerifyAPI(monitoringv1.SchemeGroupVersion.Group, monitoringv1.SchemeGroupVersion.Version)
	if err != nil {
		return err
	}
	prometheusAPIFound = found
	return nil
}
