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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/tenant"
	"github.com/kubesphere/whizard/pkg/util"
)

// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=tenants,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=tenants/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=tenants/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=service,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=storage,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=compactors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=ingesters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=rulers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=stores,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile

func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("tenant", req.NamespacedName)
	l.Info("sync")

	instance := &monitoringv1alpha1.Tenant{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	instance = r.tenantValidator(instance)
	if v, ok := instance.Labels[constants.ServiceLabelKey]; !ok || v == "" {
		l.V(1).Info("tenant does not belong to a service, no need to reconcile")
		return ctrl.Result{}, nil
	}

	baseReconciler := resources.BaseReconciler{
		Client:  r.Client,
		Log:     l,
		Scheme:  r.Scheme,
		Context: ctx,
	}
	t, err := tenant.New(baseReconciler, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, t.Reconcile()
}

// SetupWithManager sets up the controller with the Manager.
func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Tenant{}).
		Watches(&monitoringv1alpha1.Ingester{},
			handler.EnqueueRequestsFromMapFunc(r.mapToTenantbyObjectSpecFunc)).
		Watches(&monitoringv1alpha1.Compactor{},
			handler.EnqueueRequestsFromMapFunc(r.mapToTenantbyObjectSpecFunc)).
		Watches(&monitoringv1alpha1.Store{},
			handler.EnqueueRequestsFromMapFunc(r.mapToTenantByStore)).
		Watches(&monitoringv1alpha1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.mapToTenantByService)).
		Watches(&monitoringv1alpha1.Ruler{},
			handler.EnqueueRequestsFromMapFunc(r.mapToTenantbyObjectSpecFunc)).
		Owns(&monitoringv1alpha1.Ruler{}).
		Complete(r)
}

func (r *TenantReconciler) mapToTenantbyObjectSpecFunc(ctx context.Context, o client.Object) []reconcile.Request {

	var tenants []string
	switch o := o.(type) {
	case *monitoringv1alpha1.Compactor:
		tenants = o.Spec.Tenants
	case *monitoringv1alpha1.Ingester:
		tenants = o.Spec.Tenants
	case *monitoringv1alpha1.Ruler:
		tenants = []string{o.Spec.Tenant}
	}

	if len(tenants) == 0 {
		return nil
	}

	var tenantList monitoringv1alpha1.TenantList
	if err := r.Client.List(r.Context, &tenantList); err != nil {
		log.FromContext(r.Context).WithValues("tenantlist", "").Error(err, "")
		return nil
	}
	var reqs []reconcile.Request
	for _, item := range tenantList.Items {
		//	Avoid the difference between Tenant.Name and Tenant.Spec.Tenant
		if util.Contains(tenants, item.Spec.Tenant) {
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: item.Name,
				},
			})
		}
	}

	return reqs
}

func (r *TenantReconciler) mapToTenantByService(ctx context.Context, o client.Object) []reconcile.Request {

	var tenantList monitoringv1alpha1.TenantList
	if err := r.Client.List(r.Context, &tenantList, client.MatchingLabels(util.ManagedLabelByService(o))); err != nil {
		log.FromContext(r.Context).WithValues("tenantlist", "").Error(err, "")
		return nil
	}

	var reqs []reconcile.Request
	for _, item := range tenantList.Items {
		// Only reconcile the tenant that use the default storage.
		if item.Labels != nil && item.Labels[constants.StorageLabelKey] == constants.DefaultStorage {
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: item.Name,
				},
			})
		}
	}

	return reqs
}

func (r *TenantReconciler) mapToTenantByStore(ctx context.Context, _ client.Object) []reconcile.Request {

	var tenantList monitoringv1alpha1.TenantList
	if err := r.Client.List(r.Context, &tenantList); err != nil {
		log.FromContext(r.Context).WithValues("tenantlist", "").Error(err, "")
		return nil
	}

	var reqs []reconcile.Request
	if len(tenantList.Items) > 0 {
		reqs = append(reqs, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: tenantList.Items[0].Name,
			},
		})
	}

	return reqs
}

func (r *TenantReconciler) tenantValidator(tenant *monitoringv1alpha1.Tenant) *monitoringv1alpha1.Tenant {
	if tenant.Labels == nil {
		tenant.Labels = make(map[string]string)
	}

	if tenant.Spec.Tenant == "" {
		tenant.Spec.Tenant = tenant.Name
	}

	return tenant
}
