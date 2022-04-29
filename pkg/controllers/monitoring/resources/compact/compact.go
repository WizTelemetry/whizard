package compact

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

const (
	storageDir = "/thanos"
	secretsDir = "/etc/thanos/secrets"
)

type Compact struct {
	resources.StoreBaseReconciler
	compact *v1alpha1.ThanosCompact
}

func New(reconciler resources.StoreBaseReconciler) *Compact {
	return &Compact{
		StoreBaseReconciler: reconciler,
		compact:             reconciler.Store.Spec.Thanos.Compact,
	}
}

func (r *Compact) labels() map[string]string {
	labels := r.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameThanosCompact
	labels[resources.LabelNameAppManagedBy] = r.Store.Name
	return labels
}

func (r *Compact) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameThanosCompact, r.Store.Name, nameSuffix...)
}

func (r *Compact) meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.Store.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *Compact) Reconcile() error {

	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}
