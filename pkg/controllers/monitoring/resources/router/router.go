package router

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

type Router struct {
	resources.ServiceBaseReconciler
	router *v1alpha1.Router
}

func New(reconciler resources.ServiceBaseReconciler) *Router {
	return &Router{
		ServiceBaseReconciler: reconciler,
		router:                reconciler.Service.Spec.Router,
	}
}

func (r *Router) labels() map[string]string {
	labels := r.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameRouter
	labels[resources.LabelNameAppManagedBy] = r.Service.Name
	return labels
}

func (r *Router) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameRouter, r.Service.Name, nameSuffix...)
}

func (r *Router) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.Service.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *Router) HttpAddr() string {
	return fmt.Sprintf("http://%s.%s.svc:%d",
		r.name(resources.ServiceNameSuffixOperated), r.Service.Namespace, resources.ThanosHTTPPort)
}

func (r *Router) RemoteWriteAddr() string {
	return fmt.Sprintf("http://%s.%s.svc:%d",
		r.name(resources.ServiceNameSuffixOperated), r.Service.Namespace, resources.ThanosRemoteWritePort)
}

func (r *Router) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.hashringsConfigMap,
		r.deployment,
		r.service,
	})
}
