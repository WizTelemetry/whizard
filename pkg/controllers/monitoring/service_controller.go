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
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/options"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/gateway"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/query"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/query_frontend"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/receive_router"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	DefaulterValidator ServiceDefaulterValidator
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=services/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosreceiveingestors,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=stores,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosrulers,verbs=get;list;watch
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
func (r *ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	instance, err = r.DefaulterValidator(instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	serviceBaseReconciler := resources.ServiceBaseReconciler{
		BaseReconciler: resources.BaseReconciler{
			Client:  r.Client,
			Log:     l,
			Scheme:  r.Scheme,
			Context: ctx,
		},
		Service: instance,
	}

	var reconciles []func() error
	if instance.Spec.Thanos == nil {
		// to clean up resources defined by spec.Thanos
		serviceBaseReconciler.Service.Spec.Thanos = &monitoringv1alpha1.Thanos{}
	}
	reconciles = append(reconciles, receive_router.New(serviceBaseReconciler).Reconcile)
	reconciles = append(reconciles, query_frontend.New(serviceBaseReconciler).Reconcile)
	reconciles = append(reconciles, query.New(serviceBaseReconciler).Reconcile)
	reconciles = append(reconciles, gateway.New(serviceBaseReconciler).Reconcile)
	for _, reconcile := range reconciles {
		if err := reconcile(); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Service{}).
		Watches(&source.Kind{Type: &monitoringv1alpha1.ThanosReceiveIngestor{}},
			handler.EnqueueRequestsFromMapFunc(r.mapToServiceFunc)).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Store{}},
			handler.EnqueueRequestsFromMapFunc(r.mapToServiceFunc)).
		Watches(&source.Kind{Type: &monitoringv1alpha1.ThanosRuler{}},
			handler.EnqueueRequestsFromMapFunc(r.mapToServiceFunc)).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

func (r *ServiceReconciler) mapToServiceFunc(o client.Object) []reconcile.Request {

	namespacedName := monitoringv1alpha1.ServiceNamespacedName(o)

	if namespacedName == nil {
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: *namespacedName,
	}}
}

type ServiceDefaulterValidator func(service *monitoringv1alpha1.Service) (*monitoringv1alpha1.Service, error)

func CreateServiceDefaulterValidator(opt options.Options) ServiceDefaulterValidator {
	var replicas int32 = 1

	return func(service *monitoringv1alpha1.Service) (*monitoringv1alpha1.Service, error) {

		if service.Spec.TenantHeader == "" {
			service.Spec.TenantHeader = "THANOS-TENANT"
		}
		if service.Spec.TenantLabelName == "" {
			service.Spec.TenantLabelName = "tenant_id"
		}
		if service.Spec.DefaultTenantId == "" {
			service.Spec.DefaultTenantId = "default-tenant"
		}

		if service.Spec.Gateway != nil && service.Spec.Gateway.Image == "" {
			service.Spec.Gateway.Image = opt.PaodinMonitoringGatewayImage
		}

		if service.Spec.Thanos == nil {
			return service, nil
		}

		var thanos = service.Spec.Thanos

		if thanos.Query != nil {
			if thanos.Query.Replicas == nil || *thanos.Query.Replicas < 0 {
				thanos.Query.Replicas = &replicas
			}
			if thanos.Query.Image == "" {
				thanos.Query.Image = opt.ThanosImage
			}
			if thanos.Query.Envoy.Image == "" {
				thanos.Query.Envoy.Image = opt.EnvoyImage
			}

		}
		if thanos.ReceiveRouter != nil {
			if thanos.ReceiveRouter.Replicas == nil || *thanos.ReceiveRouter.Replicas < 0 {
				thanos.ReceiveRouter.Replicas = &replicas
			}
			if thanos.ReceiveRouter.Image == "" {
				thanos.ReceiveRouter.Image = opt.ThanosImage
			}
		}
		if thanos.QueryFrontend != nil {
			if thanos.QueryFrontend.Replicas == nil || *thanos.QueryFrontend.Replicas < 0 {
				thanos.QueryFrontend.Replicas = &replicas
			}
			if thanos.QueryFrontend.Image == "" {
				thanos.QueryFrontend.Image = opt.ThanosImage
			}
		}

		service.Spec.Thanos = thanos

		return service, nil
	}
}
