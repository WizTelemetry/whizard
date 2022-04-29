package gateway

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

const (
	secretsDir = "/etc/gateway/secrets"
)

type Gateway struct {
	resources.ServiceBaseReconciler
	gateway *v1alpha1.Gateway
}

func New(reconciler resources.ServiceBaseReconciler) *Gateway {
	return &Gateway{
		ServiceBaseReconciler: reconciler,
		gateway:               reconciler.Service.Spec.Gateway,
	}
}

func (g *Gateway) labels() map[string]string {
	labels := g.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameGateway
	return labels
}

func (g *Gateway) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameGateway, g.Service.Name, nameSuffix...)
}

func (g *Gateway) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       g.Service.Namespace,
		Labels:          g.labels(),
		OwnerReferences: g.OwnerReferences(),
	}
}

func (g *Gateway) Reconcile() error {
	return g.ReconcileResources([]resources.Resource{
		g.deployment,
		g.service,
	})
}
