package storegateway

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

var (
	storageDir = "/thanos"
	secretsDir = "/etc/thanos/secrets"
)

type StoreGateway struct {
	resources.BaseReconciler
	store *v1alpha1.Store
}

func New(reconciler resources.BaseReconciler, instance *v1alpha1.Store) *StoreGateway {
	return &StoreGateway{
		BaseReconciler: reconciler,
		store:          instance,
	}
}

func (r *StoreGateway) labels() map[string]string {
	labels := r.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameThanosStoreGateway
	labels[resources.LabelNameAppManagedBy] = r.store.Name
	return labels
}

func (r *StoreGateway) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameThanosStoreGateway, r.store.Name, nameSuffix...)
}

func (r *StoreGateway) meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.store.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *StoreGateway) OwnerReferences() []metav1.OwnerReference {
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

func (r *StoreGateway) GrpcAddrs() []string {
	var addrs []string
	if r.store == nil {
		return addrs
	}
	addrs = append(addrs, fmt.Sprintf("%s.%s.svc:%d",
		r.name(resources.ServiceNameSuffixOperated), r.store.Namespace, resources.ThanosGRPCPort))
	return addrs
}

func (r *StoreGateway) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}
