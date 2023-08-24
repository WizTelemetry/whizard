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
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/options"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/gateway"
	"github.com/kubesphere/whizard/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// GatewayReconciler reconciles a Service object
type GatewayReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
	Options *options.GatewayOptions
}

//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=gateways/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=gateways/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=services,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=queries,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=queryfrontends,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=routers,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services;configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("gateway", req.NamespacedName)

	l.Info("sync")

	instance := &monitoringv1alpha1.Gateway{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if instance.Labels == nil ||
		instance.Labels[constants.ServiceLabelKey] == "" {
		return ctrl.Result{}, nil
	}

	instance = r.validator(instance)
	gatewayReconciler, err := gateway.New(
		resources.BaseReconciler{
			Client:  r.Client,
			Log:     l,
			Scheme:  r.Scheme,
			Context: ctx,
		},
		instance,
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, gatewayReconciler.Reconcile()
}

var p predicate.Predicate = predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		return false
	},
}

type ResourceCustomPredicate struct {
	predicate.Funcs
}

// Update ignore cluster update event
func (r *ResourceCustomPredicate) Update(e event.UpdateEvent) bool {
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Gateway{}).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Service{}},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelByService))).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Query{}},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelBySameService))).
		Watches(&source.Kind{Type: &monitoringv1alpha1.QueryFrontend{}},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelBySameService))).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Router{}},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelBySameService))).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Tenant{}},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelBySameService)), builder.WithPredicates(p)).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

func (r *GatewayReconciler) mapFuncBySelectorFunc(fn func(metav1.Object) map[string]string) handler.MapFunc {
	return func(o client.Object) []reconcile.Request {
		gatewayList := &monitoringv1alpha1.GatewayList{}
		if err := r.Client.List(r.Context, gatewayList, client.MatchingLabels(fn(o))); err != nil {
			log.FromContext(r.Context).WithValues("gatewayList", "").Error(err, "")
			return nil
		}

		var reqs []reconcile.Request
		for _, item := range gatewayList.Items {
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

func (r *GatewayReconciler) validator(g *monitoringv1alpha1.Gateway) *monitoringv1alpha1.Gateway {
	r.Options.Override(&g.Spec)
	return g
}
