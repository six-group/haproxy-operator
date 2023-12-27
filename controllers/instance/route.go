package instance

import (
	"context"
	"fmt"

	routev1 "github.com/openshift/api/route/v1"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"github.com/six-group/haproxy-operator/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var routeAPIFound = false

func (r *Reconciler) reconcileRoute(ctx context.Context, instance *proxyv1alpha1.Instance, listens *configv1alpha1.ListenList, frontends *configv1alpha1.FrontendList) error {
	if IsRouteAPIAvailable() {
		for i := range listens.Items {
			listen := listens.Items[i]
			if err := r.createOrUpdateRouteForFrontend(ctx, instance, listen.ToFrontend()); err != nil {
				return err
			}
		}

		for i := range frontends.Items {
			frontend := frontends.Items[i]
			if err := r.createOrUpdateRouteForFrontend(ctx, instance, &frontend); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Reconciler) createOrUpdateRouteForFrontend(ctx context.Context, instance *proxyv1alpha1.Instance, frontend *configv1alpha1.Frontend) error {
	logger := log.FromContext(ctx)

	for _, bind := range frontend.Spec.Binds {
		if pointer.BoolDeref(bind.Hidden, false) {
			continue
		}

		route := &routev1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Name:      utils.GetRouteName(frontend, bind),
				Namespace: instance.Namespace,
			},
		}

		result, err := controllerutil.CreateOrUpdate(ctx, r.Client, route, func() error {
			if err := controllerutil.SetOwnerReference(frontend, route, r.Scheme); err != nil {
				return err
			}

			route.Spec.To = routev1.RouteTargetReference{
				Kind: "Service",
				Name: utils.GetServiceName(instance),
			}

			route.Spec.Port = &routev1.RoutePort{
				TargetPort: intstr.FromInt(int(bind.Port)),
			}

			route.Spec.TLS = instance.Spec.Network.Route.TLS

			return nil
		})
		if err != nil {
			return err
		}
		if result != controllerutil.OperationResultNone {
			logger.Info(fmt.Sprintf("Object %s", result), "route", route.Name)
		}
	}

	return nil
}

// IsRouteAPIAvailable returns true if the Route API is present.
func IsRouteAPIAvailable() bool {
	return routeAPIFound
}

// VerifyRouteAPI will verify that the Route API is present.
func VerifyRouteAPI() error {
	found, err := utils.VerifyAPI(routev1.GroupVersion.Group, routev1.GroupVersion.Version)
	if err != nil {
		return err
	}
	routeAPIFound = found
	return nil
}
