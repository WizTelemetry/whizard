package router

import (
	"fmt"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	configDir     = "/etc/whizard"
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
	labels[constants.LabelNameAppName] = constants.AppNameRouter
	labels[constants.LabelNameAppManagedBy] = r.Service.Name
	return labels
}

func (r *Router) name(nameSuffix ...string) string {
	return resources.QualifiedName(constants.AppNameRouter, r.Service.Name, nameSuffix...)
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
		r.name(constants.ServiceNameSuffix), r.Service.Namespace, constants.HTTPPort)
}

func (r *Router) RemoteWriteAddr() string {
	return fmt.Sprintf("http://%s.%s.svc:%d",
		r.name(constants.ServiceNameSuffix), r.Service.Namespace, constants.RemoteWritePort)
}

func (r *Router) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.hashringsConfigMap,
		r.deployment,
		r.service,
	})
}
