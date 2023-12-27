package config

import (
	"context"
	"fmt"
	"reflect"

	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciler reconciles any configv1alpha1.Object
type Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Object configv1alpha1.Object
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	object, ok := (r.Object.DeepCopyObject()).(configv1alpha1.Object)
	if !ok {
		logger.Error(fmt.Errorf("interface conversion: %s is not v1alpha1.Object", reflect.TypeOf(r.Object)), "")
		return ctrl.Result{}, nil
	}

	if err := r.Get(ctx, req.NamespacedName, object); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if len(object.GetOwnerReferences()) > 0 {
		return ctrl.Result{}, nil
	}

	instances := &proxyv1alpha1.InstanceList{}
	if err := r.List(ctx, instances, client.InNamespace(object.GetNamespace())); err != nil {
		return reconcile.Result{}, err
	}

	for idx := range instances.Items {
		instance := instances.Items[idx]

		selector, err := metav1.LabelSelectorAsSelector(&instance.Spec.Configuration.LabelSelector)
		if err != nil {
			continue
		}

		if selector.Matches(labels.Set(object.GetLabels())) {
			if err := controllerutil.SetControllerReference(&instance, object, r.Scheme); err != nil {
				return reconcile.Result{}, err
			}

			return ctrl.Result{}, r.Update(ctx, object)
		}
	}

	object.SetStatus(configv1alpha1.Status{
		Phase: configv1alpha1.StatusPhaseInternalError,
		Error: "No Instance with a matching label selector found",
	})

	return ctrl.Result{}, r.Status().Update(ctx, object)
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(r.Object).
		WithEventFilter(predicate.NewPredicateFuncs(func(object client.Object) bool {
			return len(object.GetOwnerReferences()) == 0
		})).
		Complete(r)
}
