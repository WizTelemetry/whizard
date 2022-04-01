package storegateway

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/paodin-monitoring/pkg/api/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/resources"
)

var (
	storageDir = "/thanos"
	secretsDir = "/etc/thanos/secrets"
)

type StoreGateway struct {
	resources.ThanosBaseReconciler
	store *v1alpha1.StoreGateway
}

func New(reconciler resources.ThanosBaseReconciler) *StoreGateway {
	return &StoreGateway{
		ThanosBaseReconciler: reconciler,
		store:                reconciler.Thanos.Spec.StoreGateway,
	}
}

func (r *StoreGateway) labels() map[string]string {
	labels := r.BaseLabels()
	labels["app.kubernetes.io/name"] = "storegateway"
	return labels
}

func (r *StoreGateway) name(nameSuffix ...string) string {
	name := r.Thanos.Name + "-storegateway"
	if len(nameSuffix) > 0 {
		name += "-" + strings.Join(nameSuffix, "-")
	}
	return name
}

func (r *StoreGateway) meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.Thanos.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *StoreGateway) GrpcAddrs() []string {
	var addrs []string
	if r.store == nil {
		return addrs
	}
	addrs = append(addrs, fmt.Sprintf("%s:%d", r.name("operated"), 10901))
	return addrs
}

func (r *StoreGateway) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}
