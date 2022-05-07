/*
Copyright 2021 The KubeSphere authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package monitoring

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/options"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/receive_ingestor"
)

// ThanosReceiveIngestorReconciler reconciles a ThanosReceiveIngestor object
type ThanosReceiveIngestorReconciler struct {
	DefaulterValidator ThanosReceiveIngestorDefaulterValidator
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosreceiveingestors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosreceiveingestors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosreceiveingestors/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=services,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ThanosReceiveIngestorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("thanosreceiveingestor", req.NamespacedName)

	l.Info("sync")

	instance := &monitoringv1alpha1.ThanosReceiveIngestor{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	instance, err = r.DefaulterValidator(instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	baseReconciler := resources.BaseReconciler{
		Client:  r.Client,
		Log:     l,
		Scheme:  r.Scheme,
		Context: ctx,
	}

	if err := receive_ingestor.New(baseReconciler, instance).Reconcile(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ThanosReceiveIngestorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.ThanosReceiveIngestor{}).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Service{}},
			handler.EnqueueRequestsFromMapFunc(r.mapToIngestorFunc)).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func (r *ThanosReceiveIngestorReconciler) mapToIngestorFunc(o client.Object) []reconcile.Request {

	var ingestorList monitoringv1alpha1.ThanosReceiveIngestorList
	if err := r.Client.List(r.Context, &ingestorList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(o))); err != nil {
		log.FromContext(r.Context).WithValues("thanosreceiveingestorlist", "").Error(err, "")
		return nil
	}

	var reqs []reconcile.Request
	for _, ingestor := range ingestorList.Items {
		reqs = append(reqs, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: ingestor.Namespace,
				Name:      ingestor.Name,
			},
		})
	}

	return reqs
}

type ThanosReceiveIngestorDefaulterValidator func(ingestor *monitoringv1alpha1.ThanosReceiveIngestor) (*monitoringv1alpha1.ThanosReceiveIngestor, error)

func CreateThanosReceiveIngestorDefaulterValidator(opt options.Options) ThanosReceiveIngestorDefaulterValidator {
	var replicas int32 = 1

	return func(ingestor *monitoringv1alpha1.ThanosReceiveIngestor) (*monitoringv1alpha1.ThanosReceiveIngestor, error) {

		if ingestor.Spec.Image == "" {
			ingestor.Spec.Image = opt.ThanosImage
		}
		if ingestor.Spec.Replicas == nil || *ingestor.Spec.Replicas < 0 {
			ingestor.Spec.Replicas = &replicas
		}

		return ingestor, nil
	}
}
