package receive_ingestor

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

const (
	secretsDir = "/etc/thanos/secrets"
	storageDir = "/thanos"
)

type ReceiveIngestor struct {
	resources.BaseReconciler
	ingestor *monitoringv1alpha1.ThanosReceiveIngestor
}

func New(reconciler resources.BaseReconciler, ingestor *monitoringv1alpha1.ThanosReceiveIngestor) *ReceiveIngestor {
	return &ReceiveIngestor{
		BaseReconciler: reconciler,
		ingestor:       ingestor,
	}
}

func (r *ReceiveIngestor) labels() map[string]string {
	labels := r.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameThanosReceiveIngestor
	labels[resources.LabelNameAppManagedBy] = r.ingestor.Name
	return labels
}

func (r *ReceiveIngestor) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameThanosReceiveIngestor, r.ingestor.Name, nameSuffix...)
}

func (r *ReceiveIngestor) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.ingestor.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *ReceiveIngestor) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: r.ingestor.APIVersion,
			Kind:       r.ingestor.Kind,
			Name:       r.ingestor.Name,
			UID:        r.ingestor.UID,
			Controller: pointer.BoolPtr(true),
		},
	}
}

func (r *ReceiveIngestor) GrpcAddrs() []string {
	var addrs []string
	if r.ingestor.Spec.Replicas == nil {
		addrs = make([]string, 1)
	} else {
		addrs = make([]string, *r.ingestor.Spec.Replicas)
	}
	for i := range addrs {
		addrs[i] = fmt.Sprintf("%s-%d.%s.%s.svc:%d",
			r.name(), i, r.name(resources.ServiceNameSuffixOperated), r.ingestor.Namespace, resources.ThanosGRPCPort)
	}
	return addrs
}

func (r *ReceiveIngestor) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}
