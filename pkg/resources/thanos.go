package resources

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/util"
)

var (
	ManagedByLabel = "app.kubernetes.io/managed-by"
)

type Operation string

const (
	OperationCreateOrUpdate Operation = "CreateOrUpdate"
	OperationDelete         Operation = "Delete"
)

type Resource func() (runtime.Object, Operation, error)

type ThanosBaseReconciler struct {
	Thanos *v1alpha1.Thanos

	Context  context.Context
	Client   client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *ThanosBaseReconciler) BaseLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/managed-by": r.Thanos.Name,
		"app.kubernetes.io/part-of":    "thanos",
	}
}

func (r *ThanosBaseReconciler) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: r.Thanos.APIVersion,
			Kind:       r.Thanos.Kind,
			Name:       r.Thanos.Name,
			UID:        r.Thanos.UID,
			Controller: pointer.BoolPtr(true),
		},
	}
}

func (r *ThanosBaseReconciler) ReconcileResources(resources []Resource) error {
	for _, resource := range resources {
		obj, operation, err := resource()
		if err != nil {
			return err
		}
		switch operation {
		case OperationDelete:
			err := r.Client.Delete(r.Context, obj.(client.Object))
			if !apierrors.IsNotFound(err) {
				return err
			}
		case OperationCreateOrUpdate:
			if func() error {
				switch desired := obj.(type) {
				case *appsv1.Deployment:
					return util.CreateOrUpdateDeployment(r.Context, r.Client, desired)
				default:
					return util.CreateOrUpdate(r.Context, r.Client, desired.(client.Object))
				}
			}(); err != nil {
				return err
			}
		}
	}

	return nil
}
