package compact

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

type Compact struct {
	resources.BaseReconciler
	compact *v1alpha1.Compact
}

func New(reconciler resources.BaseReconciler, compact *v1alpha1.Compact) *Compact {
	return &Compact{
		BaseReconciler: reconciler,
		compact:        compact,
	}
}

func (r *Compact) labels() map[string]string {
	labels := r.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameThanosCompact
	labels[resources.LabelNameAppManagedBy] = r.compact.Name
	return labels
}

func (r *Compact) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameThanosCompact, r.compact.Name, nameSuffix...)
}

func (r *Compact) meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.compact.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *Compact) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: r.compact.APIVersion,
			Kind:       r.compact.Kind,
			Name:       r.compact.Name,
			UID:        r.compact.UID,
			Controller: pointer.BoolPtr(true),
		},
	}
}

func (r *Compact) Reconcile() error {

	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}
