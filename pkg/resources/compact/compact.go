package compact

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/paodin-monitoring/api/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/resources"
)

var (
	storageDir = "/thanos"
	secretsDir = "/etc/thanos/secrets"
)

type Compact struct {
	resources.ThanosBaseReconciler
	compact *v1alpha1.Compact
}

func New(reconciler resources.ThanosBaseReconciler) *Compact {
	return &Compact{
		ThanosBaseReconciler: reconciler,
		compact:              reconciler.Thanos.Spec.Compact,
	}
}

func (r *Compact) labels() map[string]string {
	labels := r.BaseLabels()
	labels["app.kubernetes.io/name"] = "compact"
	return labels
}

func (r *Compact) name(nameSuffix ...string) string {
	name := r.Thanos.Name + "-compact"
	if len(nameSuffix) > 0 {
		name += "-" + strings.Join(nameSuffix, "-")
	}
	return name
}

func (r *Compact) meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.Thanos.Namespace,
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
