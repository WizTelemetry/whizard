package resources

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
)

type StoreBaseReconciler struct {
	Store *v1alpha1.Store

	BaseReconciler
}

func (r *StoreBaseReconciler) BaseLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/managed-by": r.Store.Name,
		"app.kubernetes.io/part-of":    "store",
	}
}

func (r *StoreBaseReconciler) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: r.Store.APIVersion,
			Kind:       r.Store.Kind,
			Name:       r.Store.Name,
			UID:        r.Store.UID,
			Controller: pointer.BoolPtr(true),
		},
	}
}
