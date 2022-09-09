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
	"fmt"
	"strconv"
	"time"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/options"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/ingester"
	"github.com/kubesphere/whizard/pkg/util"
	"github.com/prometheus/common/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// IngesterReconciler reconciles a Ingester object
type IngesterReconciler struct {
	DefaulterValidator IngesterDefaulterValidator
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=ingesters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=ingesters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=ingesters/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=services,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *IngesterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("ingester", req.NamespacedName)

	l.Info("sync")

	instance := &monitoringv1alpha1.Ingester{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// recycle ingester by using the RequeueAfter event
	if v, ok := instance.Annotations[constants.LabelNameIngesterState]; ok && v == constants.IngesterStateDeleting && len(instance.Spec.Tenants) == 0 {
		if deletingTime, ok := instance.Annotations[constants.LabelNameIngesterDeletingTime]; ok {
			i, err := strconv.ParseInt(deletingTime, 10, 64)
			if err == nil {
				d := time.Since(time.Unix(i, 0))
				if d < 0 {
					l.Info("recycle", "recycled time", (-d).String())
					return ctrl.Result{Requeue: true, RequeueAfter: -d}, nil
				} else {
					err := r.Delete(r.Context, instance)
					return ctrl.Result{}, err
				}
			}
		}
	}

	instance, err = r.DefaulterValidator(instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	if instance.Labels == nil ||
		instance.Labels[constants.ServiceLabelKey] == "" {
		return ctrl.Result{}, nil
	}

	baseReconciler := resources.BaseReconciler{
		Client:  r.Client,
		Log:     l,
		Scheme:  r.Scheme,
		Context: ctx,
	}

	if err := ingester.New(baseReconciler, instance).Reconcile(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IngesterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Ingester{}).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Service{}},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelByService))).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Storage{}},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelByStorage))).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func (r *IngesterReconciler) mapFuncBySelectorFunc(fn func(metav1.Object) map[string]string) handler.MapFunc {
	return func(o client.Object) []reconcile.Request {
		ingesterList := &monitoringv1alpha1.IngesterList{}
		if err := r.Client.List(r.Context, ingesterList, client.MatchingLabels(fn(o))); err != nil {
			log.FromContext(r.Context).WithValues("ingesterList", "").Error(err, "")
			return nil
		}

		var reqs []reconcile.Request
		for _, item := range ingesterList.Items {
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: item.Namespace,
					Name:      item.Name,
				},
			})
		}

		return reqs
	}
}

type IngesterDefaulterValidator func(ingester *monitoringv1alpha1.Ingester) (*monitoringv1alpha1.Ingester, error)

func CreateIngesterDefaulterValidator(opt *options.IngesterOptions) IngesterDefaulterValidator {

	return func(ingester *monitoringv1alpha1.Ingester) (*monitoringv1alpha1.Ingester, error) {

		opt.Apply(&ingester.Spec.CommonSpec)

		if ingester.Spec.DataVolume != nil {
			ingester.Spec.DataVolume = opt.DataVolume
		}

		if ingester.Spec.LocalTsdbRetention != "" {
			_, err := model.ParseDuration(ingester.Spec.LocalTsdbRetention)
			if err != nil {
				return nil, fmt.Errorf("invalid localTsdbRetention: %v", err)
			}
		}

		return ingester, nil
	}
}
