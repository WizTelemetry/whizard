package store

import (
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
	"k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Store) horizontalPodAutoscaler() (runtime.Object, resources.Operation, error) {
	var hpa = &v2beta2.HorizontalPodAutoscaler{ObjectMeta: r.meta(r.name())}

	if err := r.Client.Get(r.Context, client.ObjectKeyFromObject(hpa), hpa); err != nil {
		if !util.IsNotFound(err) {
			return nil, "", err
		}
	}

	hpa.Spec.ScaleTargetRef = v2beta2.CrossVersionObjectReference{
		Kind:       "StatefulSet",
		APIVersion: "apps/v1",
		Name:       r.name(),
	}

	if hpa.Labels == nil {
		hpa.Labels = r.labels()
	}

	hpa.Spec.MinReplicas = r.store.Spec.Scaler.MinReplicas
	hpa.Spec.MaxReplicas = r.store.Spec.Scaler.MaxReplicas
	hpa.Spec.Behavior = r.store.Spec.Scaler.Behavior
	hpa.Spec.Metrics = r.store.Spec.Scaler.Metrics

	return hpa, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.store, hpa, r.Scheme)
}
