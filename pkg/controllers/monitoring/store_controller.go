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

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/options"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/store"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
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
)

// StoreReconciler reconciles a Store object
type StoreReconciler struct {
	DefaulterValidator StoreDefaulterValidator
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context

	Options *options.StoreOptions
}

//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=stores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=stores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=stores/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete

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

	baseReconciler := resources.BaseReconciler{
		Client:  r.Client,
		Log:     l,
		Scheme:  r.Scheme,
		Context: ctx,
	}

	if err := store.New(baseReconciler, instance, r.Options).Reconcile(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Store{}).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Storage{}},
			handler.EnqueueRequestsFromMapFunc(r.reconcileRequestFromStorage)).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&autoscalingv2beta2.HorizontalPodAutoscaler{}).
		Complete(r)
}

func (r *StoreReconciler) reconcileRequestFromStorage(o client.Object) []reconcile.Request {
	storeList := &monitoringv1alpha1.StoreList{}
	if err := r.Client.List(r.Context, storeList, client.MatchingLabels(monitoringv1alpha1.ManagedLabelByStorage(o))); err != nil {
		log.FromContext(r.Context).WithValues("storeList", "").Error(err, "")
		return nil
	}

	var reqs []reconcile.Request
	for _, item := range storeList.Items {
		reqs = append(reqs, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: item.Namespace,
				Name:      item.Name,
			},
		})
	}

	return reqs
}

type StoreDefaulterValidator func(store *monitoringv1alpha1.Store) (*monitoringv1alpha1.Store, error)

func CreateStoreDefaulterValidator(opt *options.StoreOptions) StoreDefaulterValidator {

	return func(store *monitoringv1alpha1.Store) (*monitoringv1alpha1.Store, error) {

		opt.Apply(&store.Spec.CommonSpec)

		if store.Spec.IndexCacheConfig == nil {
			store.Spec.IndexCacheConfig = opt.IndexCacheConfig
		} else {
			if store.Spec.IndexCacheConfig.InMemoryIndexCacheConfig == nil {
				store.Spec.IndexCacheConfig.InMemoryIndexCacheConfig = opt.IndexCacheConfig.InMemoryIndexCacheConfig
			} else {
				if store.Spec.IndexCacheConfig.MaxSize == "" {
					store.Spec.IndexCacheConfig.MaxSize = opt.MaxSize
				}
			}
		}

		if store.Spec.Scaler == nil {
			store.Spec.Scaler = opt.AutoScaler
		} else {
			if store.Spec.Scaler.MaxReplicas == 0 {
				store.Spec.Scaler.MaxReplicas = opt.MaxReplicas
			}

			if store.Spec.Scaler.MinReplicas == nil || *store.Spec.Scaler.MinReplicas == 0 {
				min := *opt.MinReplicas
				store.Spec.Scaler.MinReplicas = &min
			}

			if store.Spec.Scaler.Metrics == nil {
				store.Spec.Scaler.Metrics = opt.Metrics
			}

			if store.Spec.Scaler.Behavior == nil {
				store.Spec.Scaler.Behavior = opt.Behavior
			}
		}

		if store.Spec.DataVolume == nil {
			store.Spec.DataVolume = opt.DataVolume
		}

		return store, nil
	}
}
