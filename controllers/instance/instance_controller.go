package instance

import (
	"context"

	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciler reconciles a Instance object
type Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=proxy.haproxy.com,resources=instances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=proxy.haproxy.com,resources=instances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=proxy.haproxy.com,resources=instances/finalizers,verbs=update

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance := &proxyv1alpha1.Instance{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	selector, err := metav1.LabelSelectorAsSelector(&instance.Spec.Configuration.LabelSelector)
	if err != nil {
		return reconcile.Result{}, err
	}

	listens := &configv1alpha1.ListenList{}
	if err := r.List(ctx, listens, client.InNamespace(instance.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return reconcile.Result{}, err
	}

	frontends := &configv1alpha1.FrontendList{}
	if err := r.List(ctx, frontends, client.InNamespace(instance.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return reconcile.Result{}, err
	}

	backends := &configv1alpha1.BackendList{}
	if err := r.List(ctx, backends, client.InNamespace(instance.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return reconcile.Result{}, err
	}

	resolvers := &configv1alpha1.ResolverList{}
	if err := r.List(ctx, resolvers, client.InNamespace(instance.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return reconcile.Result{}, err
	}

	if len(listens.Items) == 0 && len(frontends.Items) == 0 {
		instance.Status = proxyv1alpha1.InstanceStatus{
			Phase: proxyv1alpha1.InstancePhasePending,
			Error: "at least one listen or frontend must exist with the instance as owner",
		}

		return reconcile.Result{}, r.Status().Update(ctx, instance)
	}

	if err := r.reconcileConfig(ctx, instance, listens, frontends, backends, resolvers); err != nil {
		return reconcile.Result{}, r.handleError(ctx, instance, err)
	}

	if instance.Spec.Network.Service.Enabled {
		if err := r.reconcileService(ctx, instance, listens, frontends); err != nil {
			return reconcile.Result{}, r.handleError(ctx, instance, err)
		}
	}

	if instance.Spec.Network.Route.Enabled {
		if err := r.reconcileRoute(ctx, instance, listens, frontends); err != nil {
			return reconcile.Result{}, r.handleError(ctx, instance, err)
		}
	}

	if instance.Spec.Metrics != nil && instance.Spec.Metrics.Enabled {
		if err := r.reconcilePrometheusConfiguration(ctx, instance); err != nil {
			return reconcile.Result{}, r.handleError(ctx, instance, err)
		}
	}

	if err := r.reconcileStatefulSet(ctx, instance); err != nil {
		return reconcile.Result{}, r.handleError(ctx, instance, err)
	}

	if err := r.reconcilePDB(ctx, instance); err != nil {
		return reconcile.Result{}, r.handleError(ctx, instance, err)
	}

	instance.Status = proxyv1alpha1.InstanceStatus{
		Phase: proxyv1alpha1.InstancePhaseRunning,
	}
	if err := r.Status().Update(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	r.updateConfig(ctx, instance, listens, frontends, backends, resolvers)

	return ctrl.Result{}, nil
}

func (r *Reconciler) handleError(ctx context.Context, instance *proxyv1alpha1.Instance, err error) error {
	instance.Status = proxyv1alpha1.InstanceStatus{
		Phase: proxyv1alpha1.InstancePhaseInternalError,
		Error: err.Error(),
	}

	return r.Status().Update(ctx, instance)
}

func (r *Reconciler) updateConfig(ctx context.Context, instance *proxyv1alpha1.Instance, listens *configv1alpha1.ListenList, frontends *configv1alpha1.FrontendList, backends *configv1alpha1.BackendList, resolvers *configv1alpha1.ResolverList) {
	for i := range listens.Items {
		listen := listens.Items[i]
		_ = r.updateConfigObject(ctx, instance, &listen)
	}

	for i := range frontends.Items {
		frontend := frontends.Items[i]
		_ = r.updateConfigObject(ctx, instance, &frontend)
	}

	for i := range backends.Items {
		backend := backends.Items[i]
		_ = r.updateConfigObject(ctx, instance, &backend)
	}

	for i := range resolvers.Items {
		resolvers := resolvers.Items[i]
		_ = r.updateConfigObject(ctx, instance, &resolvers)
	}
}

func (r *Reconciler) updateConfigObject(ctx context.Context, instance *proxyv1alpha1.Instance, object configv1alpha1.Object) error {
	logger := log.FromContext(ctx)

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, object, func() error {
		return controllerutil.SetControllerReference(instance, object, r.Scheme)
	})
	if err != nil {
		logger.Error(err, "Unable to set controller reference", object.GetObjectKind().GroupVersionKind().Kind, object.GetName())
		return err
	}

	object.SetStatus(configv1alpha1.Status{
		Phase:              configv1alpha1.StatusPhaseActive,
		ObservedGeneration: object.GetGeneration(),
	})
	if err := r.Status().Update(ctx, object); err != nil {
		logger.Error(err, "Unable to update status", object.GetObjectKind().GroupVersionKind().Kind, object.GetName())
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&proxyv1alpha1.Instance{}).
		Owns(&configv1alpha1.Listen{}).
		Owns(&configv1alpha1.Frontend{}).
		Owns(&configv1alpha1.Backend{}).
		Owns(&configv1alpha1.Resolver{}).
		Complete(r)
}
