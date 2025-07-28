package ingester

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	monitoringv1alpha1 "github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
)

type Ingester struct {
	resources.BaseReconciler
	ingester *monitoringv1alpha1.Ingester
}

func New(reconciler resources.BaseReconciler, ingester *monitoringv1alpha1.Ingester) (*Ingester, error) {
	if err := reconciler.SetService(ingester); err != nil {
		return nil, err
	}
	return &Ingester{
		BaseReconciler: reconciler,
		ingester:       ingester,
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
	if r.ingester.Spec.Replicas != nil {
		addrs = make([]string, *r.ingester.Spec.Replicas)
	} else if r.Service.Spec.IngesterTemplateSpec.Replicas != nil {
		addrs = make([]string, *r.Service.Spec.IngesterTemplateSpec.Replicas)
	} else {
		// whizard ingester default replicas is 2
		addrs = make([]string, 2)
	}
	for i := range addrs {
		addrs[i] = fmt.Sprintf("%s-%d.%s.%s.svc:%d",
			r.name(), i, r.name(constants.ServiceNameSuffix), r.ingester.Namespace, constants.GRPCPort)
	}
	return addrs
}

func (r *Ingester) Endpoints() []string {
	return []string{fmt.Sprintf("dnssrv+_grpc._tcp.%s.%s.svc", r.name(constants.ServiceNameSuffix), r.ingester.Namespace)}
}

func (r *Ingester) Address() []string {
	var addrs []string
	if r.ingester.Spec.Replicas != nil {
		addrs = make([]string, *r.ingester.Spec.Replicas)
	} else if r.Service.Spec.IngesterTemplateSpec.Replicas != nil {
		addrs = make([]string, *r.Service.Spec.IngesterTemplateSpec.Replicas)
	} else {
		// whizard ingester default replicas is 2
		addrs = make([]string, 2)
	}
	for i := range addrs {
		addrs[i] = fmt.Sprintf("%s-%d.%s.%s.svc", r.name(), i, r.name(constants.ServiceNameSuffix), r.ingester.Namespace)
	}
	return addrs
}

func (r *Ingester) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}
