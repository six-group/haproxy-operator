package instance

import (
	"context"
	"fmt"
	"sort"

	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"github.com/six-group/haproxy-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *Reconciler) reconcileService(ctx context.Context, instance *proxyv1alpha1.Instance, listens *configv1alpha1.ListenList, frontends *configv1alpha1.FrontendList) error {
	logger := log.FromContext(ctx)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetServiceName(instance),
			Namespace: instance.Namespace,
		},
	}

	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		if err := controllerutil.SetOwnerReference(instance, service, r.Scheme); err != nil {
			return err
		}

		service.Labels = utils.GetAppSelectorLabels(instance)

		if len(instance.Spec.Network.HostIPs) == 0 {
			service.Spec.Selector = utils.GetAppSelectorLabels(instance)
		}

		service.Spec.Ports = []corev1.ServicePort{}
		for _, listen := range listens.Items {
			for _, bind := range listen.Spec.Binds {
				if ptr.Deref(bind.Hidden, false) {
					continue
				}

				service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
					Name:       fmt.Sprintf("tcp-%d", bind.Port),
					Port:       int32(bind.Port),
					TargetPort: intstr.FromInt32(int32(bind.Port)),
					Protocol:   corev1.ProtocolTCP,
				})
			}
		}

		for _, frontend := range frontends.Items {
			for _, bind := range frontend.Spec.Binds {
				if ptr.Deref(bind.Hidden, false) {
					continue
				}

				service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
					Name:       fmt.Sprintf("tcp-%d", bind.Port),
					Port:       int32(bind.Port),
					TargetPort: intstr.FromInt32(int32(bind.Port)),
					Protocol:   corev1.ProtocolTCP,
				})
			}
		}

		if instance.Spec.Metrics != nil && instance.Spec.Metrics.Enabled {
			service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
				Name:       "metrics",
				Port:       int32(instance.Spec.Metrics.Port),
				TargetPort: intstr.FromInt32(int32(instance.Spec.Metrics.Port)),
				Protocol:   corev1.ProtocolTCP,
			})
		}

		sort.Slice(service.Spec.Ports, func(i, j int) bool {
			return service.Spec.Ports[i].Name < service.Spec.Ports[j].Name
		})

		return nil
	})
	if err != nil {
		return err
	}
	if result != controllerutil.OperationResultNone {
		logger.Info(fmt.Sprintf("Object %s", result), "service", service.Name)
	}

	if len(instance.Spec.Network.HostIPs) > 0 {
		if err := r.reconcileServiceEndpoints(ctx, instance, service); err != nil {
			return err
		}
	}

	return nil
}

func (r *Reconciler) reconcileServiceEndpoints(ctx context.Context, instance *proxyv1alpha1.Instance, service *corev1.Service) error {
	logger := log.FromContext(ctx)

	endpoints := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetServiceName(instance),
			Namespace: instance.Namespace,
		},
	}

	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, endpoints, func() error {
		if err := controllerutil.SetOwnerReference(instance, endpoints, r.Scheme); err != nil {
			return err
		}

		var addresses []corev1.EndpointAddress
		for host, ip := range instance.Spec.Network.HostIPs {
			addresses = append(addresses, corev1.EndpointAddress{
				IP:       ip,
				NodeName: ptr.To(host),
			})
		}
		sort.Slice(addresses, func(i, j int) bool {
			return addresses[i].IP < addresses[j].IP
		})

		var ports []corev1.EndpointPort
		for _, port := range service.Spec.Ports {
			ports = append(ports, corev1.EndpointPort{
				Name:     port.Name,
				Port:     port.Port,
				Protocol: port.Protocol,
			})
		}
		sort.Slice(ports, func(i, j int) bool {
			return ports[i].Name < ports[j].Name
		})

		endpoints.Subsets = []corev1.EndpointSubset{
			{
				Addresses: addresses,
				Ports:     ports,
			},
		}

		return nil
	})
	if err != nil {
		return err
	}
	if result != controllerutil.OperationResultNone {
		logger.Info(fmt.Sprintf("Object %s", result), "endpoints", service.Name)
	}

	return nil
}
