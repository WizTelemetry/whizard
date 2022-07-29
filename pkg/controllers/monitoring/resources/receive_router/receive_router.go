package receive_router

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

const (
	configDir     = "/etc/thanos"
	hashringsFile = "hashrings.json"
)

type ReceiveRouter struct {
	resources.ServiceBaseReconciler
	router *v1alpha1.Router
}

func New(reconciler resources.ServiceBaseReconciler) *ReceiveRouter {
	return &ReceiveRouter{
		ServiceBaseReconciler: reconciler,
		router:                reconciler.Service.Spec.Router,
	}
}

func (r *ReceiveRouter) labels() map[string]string {
	labels := r.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameThanosReceiveRouter
	labels[resources.LabelNameAppManagedBy] = r.Service.Name
	return labels
}

func (r *ReceiveRouter) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameThanosReceiveRouter, r.Service.Name, nameSuffix...)
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
	return fmt.Sprintf("http://%s.%s.svc:%d",
		r.name(resources.ServiceNameSuffixOperated), r.Service.Namespace, resources.ThanosHTTPPort)
}

func (r *ReceiveRouter) RemoteWriteAddr() string {
	return fmt.Sprintf("http://%s.%s.svc:%d",
		r.name(resources.ServiceNameSuffixOperated), r.Service.Namespace, resources.ThanosRemoteWritePort)
}

func (r *ReceiveRouter) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.hashringsConfigMap,
		r.deployment,
		r.service,
	})
}
