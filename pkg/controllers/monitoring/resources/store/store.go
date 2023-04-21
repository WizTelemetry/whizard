package store

import (
	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/options"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type Store struct {
	resources.BaseReconciler
	store *v1alpha1.Store
	*options.StoreOptions
}

func New(reconciler resources.BaseReconciler, instance *v1alpha1.Store, o *options.StoreOptions) *Store {
	return &Store{
		BaseReconciler: reconciler,
		store:          instance,
		StoreOptions:   o,
	}
}

func (r *Store) labels() map[string]string {
	labels := r.BaseLabels()
	labels[constants.LabelNameAppName] = constants.AppNameStore
	labels[constants.LabelNameAppManagedBy] = r.store.Name

	// Do not copy all labels of the custom resource to the managed workload.
	// util.AppendLabel(labels, r.store.Labels)

	// TODO handle metadata.labels and labelSelector separately in the managed workload,
	//		because labelSelector is an immutable field to be carefully treated.

	return labels
}

func (r *Store) name(nameSuffix ...string) string {
	return r.QualifiedName(constants.AppNameStore, r.store.Name, nameSuffix...)
}

func (r *Store) meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.store.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *Store) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: r.store.APIVersion,
			Kind:       r.store.Kind,
			Name:       r.store.Name,
			UID:        r.store.UID,
			Controller: pointer.BoolPtr(true),
		},
	}
}

func (r *Store) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.horizontalPodAutoscaler,
		r.service,
	})
}
