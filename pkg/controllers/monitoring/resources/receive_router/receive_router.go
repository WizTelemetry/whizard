package receive_router

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/controllers/monitoring/resources"
)

var (
	configDir     = "/etc/thanos"
	hashringsFile = "hashrings.json"
)

type ReceiveRouter struct {
	resources.ServiceBaseReconciler
	router *v1alpha1.ThanosReceiveRouter
}

func New(reconciler resources.ServiceBaseReconciler) *ReceiveRouter {
	return &ReceiveRouter{
		ServiceBaseReconciler: reconciler,
		router:                reconciler.Service.Spec.Thanos.ReceiveRouter,
	}
}

func (r *ReceiveRouter) labels() map[string]string {
	labels := r.BaseLabels()
	labels["app.kubernetes.io/name"] = "thanos-receive-router"
	return labels
}

func (r *ReceiveRouter) name(nameSuffix ...string) string {
	name := "thanos-receive-router-" + r.Service.Name
	if len(nameSuffix) > 0 {
		name += "-" + strings.Join(nameSuffix, "-")
	}
	return name
}

func (r *ReceiveRouter) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.Service.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *ReceiveRouter) HttpAddr() string {
	return fmt.Sprintf("http://%s:10902", r.name("operated"))
}

func (r *ReceiveRouter) RemoteWriteAddr() string {
	return fmt.Sprintf("http://%s:19291", r.name("operated"))
}

func (r *ReceiveRouter) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.hashringsConfigMap,
		r.deployment,
		r.service,
	})
}
