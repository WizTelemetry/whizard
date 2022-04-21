package resources

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
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

type BaseReconciler struct {
	Context  context.Context
	Client   client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *BaseReconciler) BaseLabels() map[string]string {
	return map[string]string{}
}

func (r *BaseReconciler) ReconcileResources(resources []Resource) error {
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
			if err = func() error {
				switch desired := obj.(type) {
				case *appsv1.Deployment:
					return util.CreateOrUpdateDeployment(r.Context, r.Client, desired)
				case *corev1.Service:
					return util.CreateOrUpdateService(r.Context, r.Client, desired)
				case *corev1.ServiceAccount:
					return util.CreateOrUpdateServiceAccount(r.Context, r.Client, desired)
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
