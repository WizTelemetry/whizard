package query_frontend

import (
	"fmt"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	configDir       = "/etc/thanos"
	cacheConfigFile = "cache-config.yaml"
)

type QueryFrontend struct {
	resources.ServiceBaseReconciler
	queryFrontend *monitoringv1alpha1.ThanosQueryFrontend
}

func New(reconciler resources.ServiceBaseReconciler) *QueryFrontend {
	return &QueryFrontend{
		ServiceBaseReconciler: reconciler,
		queryFrontend:         reconciler.Service.Spec.Thanos.QueryFrontend,
	}
}

func (q *QueryFrontend) labels() map[string]string {
	labels := q.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameThanosQuery
	labels[resources.LabelNameAppManagedBy] = q.Service.Name
	return labels
}

func (q *QueryFrontend) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameThanosQueryFrontend, q.Service.Name, nameSuffix...)
}

func (q *QueryFrontend) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       q.Service.Namespace,
		Labels:          q.labels(),
		OwnerReferences: q.OwnerReferences(),
	}
}

func (q *QueryFrontend) HttpAddr() string {
	return fmt.Sprintf("http://%s.%s.svc:%d",
		q.name(resources.ServiceNameSuffixOperated), q.Service.Namespace, resources.ThanosHTTPPort)
}

func (q *QueryFrontend) Reconcile() error {
	return q.ReconcileResources([]resources.Resource{
		q.cacheConfigConfigMap,
		q.deployment,
		q.service,
	})
}
