package gateway

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/resources"
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
	labels["app.kubernetes.io/name"] = "gateway"
	return labels
}

func (g *Gateway) name(nameSuffix ...string) string {
	name := g.Service.Name + "-gateway"
	if len(nameSuffix) > 0 {
		name += "-" + strings.Join(nameSuffix, "-")
	}
	return name
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
		g.role,
		g.serviceAccount,
		g.roleBinding,
		g.deployment,
		g.service,
	})
}
