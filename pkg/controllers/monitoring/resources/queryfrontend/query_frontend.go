package queryfrontend

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

const (
	configDir       = "/etc/whizard"
	cacheConfigFile = "cache-config.yaml"
	envoyConfigFile = "envoy.yaml"
)

type QueryFrontend struct {
	resources.BaseReconciler
	queryFrontend *monitoringv1alpha1.QueryFrontend
}

func New(reconciler resources.BaseReconciler, q *monitoringv1alpha1.QueryFrontend) (*QueryFrontend, error) {
	if err := reconciler.SetService(q); err != nil {
		return nil, err
	}
	return &QueryFrontend{
		BaseReconciler: reconciler,
		queryFrontend:  q,
	}, nil
}

func (q *QueryFrontend) labels() map[string]string {
	labels := q.BaseLabels()
	labels[constants.LabelNameAppName] = constants.AppNameQueryFrontend
	labels[constants.LabelNameAppManagedBy] = q.Service.Name
	return labels
}

func (q *QueryFrontend) name(nameSuffix ...string) string {
	return q.QualifiedName(constants.AppNameQueryFrontend, q.queryFrontend.Name, nameSuffix...)
}

func (q *QueryFrontend) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: q.Service.Namespace,
		Labels:    q.labels(),
	}
}

func (q *QueryFrontend) HttpAddr() string {
	return fmt.Sprintf("http://%s.%s.svc:%d",
		q.name(constants.ServiceNameSuffix), q.Service.Namespace, constants.HTTPPort)
}

func (q *QueryFrontend) HttpsAddr() string {
	return fmt.Sprintf("https://%s.%s.svc:%d",
		q.name(constants.ServiceNameSuffix), q.Service.Namespace, constants.HTTPPort)
}

func (q *QueryFrontend) Reconcile() error {
	return q.ReconcileResources([]resources.Resource{
		q.cacheConfigConfigMap,
		q.webConfigSecret,
		q.deployment,
		q.service,
	})
}
