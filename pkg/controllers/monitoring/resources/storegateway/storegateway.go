package storegateway

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

var (
	storageDir = "/thanos"
	secretsDir = "/etc/thanos/secrets"
)

type StoreGateway struct {
	resources.StoreBaseReconciler
	store *v1alpha1.ThanosStoreGateway
}

func New(reconciler resources.StoreBaseReconciler) *StoreGateway {
	return &StoreGateway{
		StoreBaseReconciler: reconciler,
		store:               reconciler.Store.Spec.Thanos.StoreGateway,
	}
}

func (r *StoreGateway) labels() map[string]string {
	labels := r.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameThanosStoreGateway
	labels[resources.LabelNameAppManagedBy] = r.Store.Name
	return labels
}

func (r *StoreGateway) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameThanosStoreGateway, r.Store.Name, nameSuffix...)
}

func (r *StoreGateway) meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.Store.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *StoreGateway) GrpcAddrs() []string {
	var addrs []string
	if r.store == nil {
		return addrs
	}
	addrs = append(addrs, fmt.Sprintf("%s.%s.svc:%d",
		r.name(resources.ServiceNameSuffixOperated), r.Store.Namespace, resources.ThanosGRPCPort))
	return addrs
}

func (r *StoreGateway) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}
