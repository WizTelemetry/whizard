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
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/router"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// RouterReconciler reconciles a Service object
type RouterReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
	Options *options.RouterOptions
}

//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=services/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=ingesters,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=stores,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=rulers,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services;configmaps;serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *RouterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("service", req.NamespacedName)

	l.Info("sync")

	instance := &monitoringv1alpha1.Service{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	instance, err = r.validator(instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	routerReconciler := router.New(
		resources.ServiceBaseReconciler{
			Service: instance,
			BaseReconciler: resources.BaseReconciler{
				Client:  r.Client,
				Log:     l,
				Scheme:  r.Scheme,
				Context: ctx,
			},
		})

	return ctrl.Result{}, routerReconciler.Reconcile()
}

// SetupWithManager sets up the controller with the Manager.
func (r *RouterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Service{}).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Ingester{}},
			handler.EnqueueRequestsFromMapFunc(r.mapToServiceFunc)).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

func (r *RouterReconciler) mapToServiceFunc(o client.Object) []reconcile.Request {

	namespacedName := monitoringv1alpha1.ServiceNamespacedName(o)

	if namespacedName == nil {
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: *namespacedName,
	}}
}

func (r *RouterReconciler) validator(service *monitoringv1alpha1.Service) (*monitoringv1alpha1.Service, error) {
	if service.Spec.TenantHeader == "" {
		service.Spec.TenantHeader = constants.DefaultTenantHeader
	}
	if service.Spec.TenantLabelName == "" {
		service.Spec.TenantLabelName = constants.DefaultTenantLabelName
	}
	if service.Spec.DefaultTenantId == "" {
		service.Spec.DefaultTenantId = constants.DefaultTenantId
	}

	if service.Spec.Router != nil {
		r.Options.Apply(&service.Spec.Router.CommonSpec)
	}

	return service, nil

}
