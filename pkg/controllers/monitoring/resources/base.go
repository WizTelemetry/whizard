package resources

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kubesphere/whizard/pkg/util"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
				case *appsv1.StatefulSet:
					err = util.CreateOrUpdateStatefulSet(r.Context, r.Client, desired)
					sErr, ok := err.(*apierrors.StatusError)
					if ok && sErr.ErrStatus.Code == 422 && sErr.ErrStatus.Reason == metav1.StatusReasonInvalid {
						// Gather only reason for failed update
						failMsg := make([]string, len(sErr.ErrStatus.Details.Causes))
						for i, cause := range sErr.ErrStatus.Details.Causes {
							failMsg[i] = cause.Message
						}
						r.Log.Info("recreating StatefulSet because the update operation wasn't possible", "reason", strings.Join(failMsg, ", "))
						if err = r.Client.Delete(r.Context, obj.(client.Object), client.PropagationPolicy(metav1.DeletePropagationForeground)); err != nil {
							return errors.Wrap(err, "failed to delete StatefulSet to avoid forbidden action")
						}
						return nil
					}
					if err != nil {
						return errors.Wrap(err, "updating StatefulSet failed")
					}
					return nil
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
