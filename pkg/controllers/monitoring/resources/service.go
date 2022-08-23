package resources

import (
	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type ServiceBaseReconciler struct {
	Service *v1alpha1.Service

	BaseReconciler
}

func (r *ServiceBaseReconciler) BaseLabels() map[string]string {
	return map[string]string{
		constants.LabelNameAppManagedBy: r.Service.Name,
		constants.LabelNameAppPartOf:    "service",
	}
}

func (r *ServiceBaseReconciler) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: r.Service.APIVersion,
			Kind:       r.Service.Kind,
			Name:       r.Service.Name,
			UID:        r.Service.UID,
			Controller: pointer.BoolPtr(true),
		},
	}
}
