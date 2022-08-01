package compactor

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

const (
	storageDir = "/thanos"
	secretsDir = "/etc/thanos/secrets"
)

type Compactor struct {
	resources.BaseReconciler
	compactor *v1alpha1.Compactor
}

func New(reconciler resources.BaseReconciler, compactor *v1alpha1.Compactor) *Compactor {
	return &Compactor{
		BaseReconciler: reconciler,
		compactor:      compactor,
	}
}

func (r *Compactor) labels() map[string]string {
	labels := r.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameCompactor
	labels[resources.LabelNameAppManagedBy] = r.compactor.Name
	return labels
}

func (r *Compactor) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameCompactor, r.compactor.Name, nameSuffix...)
}

func (r *Compactor) meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.compactor.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *Compactor) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: r.compactor.APIVersion,
			Kind:       r.compactor.Kind,
			Name:       r.compactor.Name,
			UID:        r.compactor.UID,
			Controller: pointer.BoolPtr(true),
		},
	}
}

func (r *Compactor) Reconcile() error {

	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}