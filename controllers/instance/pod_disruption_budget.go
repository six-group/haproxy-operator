package instance

import (
	"context"
	"fmt"
	"github.com/six-group/haproxy-operator/pkg/utils"

	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *Reconciler) reconcilePDB(ctx context.Context, instance *proxyv1alpha1.Instance) error {
	logger := log.FromContext(ctx)

	pdb := &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-haproxy", instance.Name),
			Namespace: instance.Namespace,
		},
	}

	if instance.Spec.PodDisruptionBudget.MaxUnavailable != nil || instance.Spec.PodDisruptionBudget.MinAvailable != nil {
		result, err := controllerutil.CreateOrUpdate(ctx, r.Client, pdb, func() error {
			pdb.Spec.Selector = metav1.SetAsLabelSelector(utils.GetPodLabels(instance))
			pdb.Spec.MaxUnavailable = instance.Spec.PodDisruptionBudget.MaxUnavailable
			pdb.Spec.MinAvailable = instance.Spec.PodDisruptionBudget.MinAvailable
			return nil
		})
		if err != nil {
			return err
		}
		if result != controllerutil.OperationResultNone {
			logger.Info(fmt.Sprintf("Object %s", result), "poddisruptionbudget", pdb.Name)
		}

	} else {

		err := r.Get(ctx, client.ObjectKeyFromObject(pdb), pdb)
		if err == nil {
			err = r.Delete(ctx, pdb)
			if err != nil {
				return err
			}
			logger.Info("deleted", "poddisruptionbudget", pdb.Name)
		}

	}

	return nil
}
