package ingester

import (
	"fmt"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/options"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Ingester struct {
	resources.BaseReconciler
	ingester *monitoringv1alpha1.Ingester
	options  *options.IngesterOptions
}

func New(reconciler resources.BaseReconciler, ingester *monitoringv1alpha1.Ingester, options *options.IngesterOptions) (*Ingester, error) {
	if err := reconciler.SetService(ingester); err != nil {
		return nil, err
	}
	return &Ingester{
		BaseReconciler: reconciler,
		ingester:       ingester,
		options:        options,
	}, nil
}

func (r *Ingester) labels() map[string]string {
	labels := r.BaseLabels()
	labels[constants.LabelNameAppName] = constants.AppNameIngester
	labels[constants.LabelNameAppManagedBy] = r.ingester.Name
	return labels
}

func (r *Ingester) name(nameSuffix ...string) string {
	return r.QualifiedName(constants.AppNameIngester, r.ingester.Name, nameSuffix...)
}

func (r *Ingester) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: r.ingester.Namespace,
		Labels:    r.labels(),
	}
}

func (r *Ingester) GrpcAddrs() []string {
	var addrs []string
	if r.ingester.Spec.Replicas == nil {
		addrs = make([]string, 1)
	} else {
		addrs = make([]string, *r.ingester.Spec.Replicas)
	}
	for i := range addrs {
		addrs[i] = fmt.Sprintf("%s-%d.%s.%s.svc:%d",
			r.name(), i, r.name(constants.ServiceNameSuffix), r.ingester.Namespace, constants.GRPCPort)
	}
	return addrs
}

func (r *Ingester) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}
