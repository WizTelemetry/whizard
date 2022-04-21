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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	monitoringv1alpha1 "github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/controllers/monitoring/options"
	"github.com/kubesphere/paodin-monitoring/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin-monitoring/pkg/controllers/monitoring/resources/compact"
	"github.com/kubesphere/paodin-monitoring/pkg/controllers/monitoring/resources/storegateway"
)

// StoreReconciler reconciles a Store object
type StoreReconciler struct {
	DefaulterValidator StoreDefaulterValidator
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=stores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=stores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=stores/finalizers,verbs=update
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
func (r *StoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("store", req.NamespacedName)

	l.Info("sync")

	instance := &monitoringv1alpha1.Store{}
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

	storeBaseReconciler := resources.StoreBaseReconciler{
		BaseReconciler: resources.BaseReconciler{
			Client:  r.Client,
			Log:     l,
			Scheme:  r.Scheme,
			Context: ctx,
		},
		Store: instance,
	}

	if instance.Spec.Thanos == nil {
		// to clean up resources defined by spec.Thanos
		storeBaseReconciler.Store.Spec.Thanos = &monitoringv1alpha1.ThanosStore{}
	}

	var reconciles []func() error
	reconciles = append(reconciles, compact.New(storeBaseReconciler).Reconcile)
	reconciles = append(reconciles, storegateway.New(storeBaseReconciler).Reconcile)
	for _, reconcile := range reconciles {
		if err := reconcile(); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Store{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

type StoreDefaulterValidator func(store *monitoringv1alpha1.Store) (*monitoringv1alpha1.Store, error)

func CreateStoreDefaulterValidator(opt options.Options) StoreDefaulterValidator {
	var replicas int32 = 1

	return func(store *monitoringv1alpha1.Store) (*monitoringv1alpha1.Store, error) {

		if store.Spec.Thanos == nil {
			return store, nil
		}

		var thanos = store.Spec.Thanos

		if thanos.StoreGateway != nil {
			if thanos.StoreGateway.Image == "" {
				thanos.StoreGateway.Image = opt.ThanosImage
			}
			if thanos.StoreGateway.Replicas == nil || *thanos.StoreGateway.Replicas < 0 {
				thanos.StoreGateway.Replicas = &replicas
			}
		}
		if thanos.Compact != nil {
			if thanos.Compact.Image == "" {
				thanos.Compact.Image = opt.ThanosImage
			}
			if thanos.Compact.Replicas == nil || *thanos.Compact.Replicas < 0 {
				thanos.Compact.Replicas = &replicas
			}
		}

		store.Spec.Thanos = thanos

		return store, nil
	}
}
