package resources

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
)

type ServiceBaseReconciler struct {
	Service *v1alpha1.Service

	BaseReconciler
}

func (r *ServiceBaseReconciler) BaseLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/managed-by": r.Service.Name,
		"app.kubernetes.io/part-of":    "service",
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
