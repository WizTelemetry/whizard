package compact

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/controllers/monitoring/resources"
)

var (
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
	labels["app.kubernetes.io/name"] = "thanos-compact"
	return labels
}

func (r *Compact) name(nameSuffix ...string) string {
	name := "thanos-compact-" + r.Store.Name
	if len(nameSuffix) > 0 {
		name += "-" + strings.Join(nameSuffix, "-")
	}
	return name
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
