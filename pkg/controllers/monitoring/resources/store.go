package resources

import (
	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type StoreBaseReconciler struct {
	Store *v1alpha1.Store

	BaseReconciler
}

func (r *StoreBaseReconciler) BaseLabels() map[string]string {
	return map[string]string{
		constants.LabelNameAppManagedBy: r.Store.Name,
		constants.LabelNameAppPartOf:    "store",
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
