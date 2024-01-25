package gateway

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

const (
	secretsDir = "/etc/gateway/secrets"
)

type Gateway struct {
	resources.BaseReconciler
	gateway *v1alpha1.Gateway
}

func New(reconciler resources.BaseReconciler, g *v1alpha1.Gateway) (*Gateway, error) {
	if err := reconciler.SetService(g); err != nil {
		return nil, err
	}
	return &Gateway{
		BaseReconciler: reconciler,
		gateway:        g,
	}, nil
}

func (g *Gateway) labels() map[string]string {
	labels := g.BaseLabels()
	labels[constants.LabelNameAppName] = constants.AppNameGateway
	labels[constants.LabelNameAppManagedBy] = g.gateway.Name

	// Do not copy all labels of the custom resource to the managed workload.
	// util.AppendLabel(labels, g.gateway.Labels)

	// TODO handle metadata.labels and labelSelector separately in the managed workload,
	//		because labelSelector is an immutable field to be carefully treated.

	return labels

}

func (g *Gateway) name(nameSuffix ...string) string {
	return g.QualifiedName(constants.AppNameGateway, g.gateway.Name, nameSuffix...)
}

func (g *Gateway) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: g.Service.Namespace,
		Labels:    g.labels(),
	}
}

func (g *Gateway) HttpAddr() string {
	return fmt.Sprintf("http://%s.%s.svc:%d",
		g.name(constants.ServiceNameSuffix), g.Service.Namespace, 9090)
}

func (g *Gateway) HttpsAddr() string {
	return fmt.Sprintf("https://%s.%s.svc:%d",
		g.name(constants.ServiceNameSuffix), g.Service.Namespace, 9090)
}

func (g *Gateway) Reconcile() error {
	return g.ReconcileResources([]resources.Resource{
		g.deployment,
		g.service,
		g.tenantsAdmissionConfigMap,
		g.webConfigSecret,
	})
}
