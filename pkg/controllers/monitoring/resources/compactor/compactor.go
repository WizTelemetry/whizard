package compactor

import (
	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type Compactor struct {
	resources.BaseReconciler
	compactor *v1alpha1.Compactor
}

func New(reconciler resources.BaseReconciler, compactor *v1alpha1.Compactor) (*Compactor, error) {
	if err := reconciler.SetService(compactor); err != nil {
		return nil, err
	}
	return &Compactor{
		BaseReconciler: reconciler,
		compactor:      compactor,
	}, nil
}

func (r *Compactor) labels() map[string]string {
	labels := r.BaseLabels()
	labels[constants.LabelNameAppName] = constants.AppNameCompactor
	labels[constants.LabelNameAppManagedBy] = r.compactor.Name
	util.AppendLabel(labels, r.compactor.Labels)
	return labels
}

func (r *Compactor) name(nameSuffix ...string) string {
	return r.QualifiedName(constants.AppNameCompactor, r.compactor.Name, nameSuffix...)
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
	})
}
