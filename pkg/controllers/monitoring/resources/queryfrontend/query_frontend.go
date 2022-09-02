package queryfrontend

import (
	"fmt"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	configDir       = "/etc/whizard"
	cacheConfigFile = "cache-config.yaml"
)

type QueryFrontend struct {
	resources.ServiceBaseReconciler
	queryFrontend *monitoringv1alpha1.QueryFrontend
}

func New(reconciler resources.ServiceBaseReconciler) *QueryFrontend {
	return &QueryFrontend{
		ServiceBaseReconciler: reconciler,
		queryFrontend:         reconciler.Service.Spec.QueryFrontend,
	}
}

func (q *QueryFrontend) labels() map[string]string {
	labels := q.BaseLabels()
	labels[constants.LabelNameAppName] = constants.AppNameQueryFrontend
	labels[constants.LabelNameAppManagedBy] = q.Service.Name
	return labels
}

func (q *QueryFrontend) name(nameSuffix ...string) string {
	return resources.QualifiedName(constants.AppNameQueryFrontend, q.Service.Name, nameSuffix...)
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
		q.name(constants.ServiceNameSuffix), q.Service.Namespace, constants.HTTPPort)
}

func (q *QueryFrontend) Reconcile() error {
	return q.ReconcileResources([]resources.Resource{
		q.cacheConfigConfigMap,
		q.deployment,
		q.service,
	})
}
