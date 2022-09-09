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
	"github.com/kubesphere/whizard/pkg/util"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	clusterv1alpha1 "kubesphere.io/api/cluster/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ClusterReconciler reconciles a Service object
type ClusterReconciler struct {
	client.Client
	Scheme                          *runtime.Scheme
	Context                         context.Context
	KubesphereAdapterDefaultService string
}

//+kubebuilder:rbac:groups=cluster.kubesphere.io,resources=clusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=tenants,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("cluster", req.NamespacedName)

	l.Info("sync")

	instance := &clusterv1alpha1.Cluster{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	tenant := &monitoringv1alpha1.Tenant{}
	err = r.Get(ctx, types.NamespacedName{Name: req.Name}, tenant)
	if err != nil {
		if apierrors.IsNotFound(err) {
			tenant = r.createTenantInstance(instance)
			util.CreateOrUpdate(ctx, r.Client, tenant)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1alpha1.Cluster{}).WithEventFilter(&ResourceCustomPredicate{}).
		Owns(&monitoringv1alpha1.Tenant{}).
		Complete(r)
}

func (r *ClusterReconciler) mapToTenantFunc(o client.Object) []reconcile.Request {
	req := types.NamespacedName{
		Name: o.GetName(),
	}
	return []ctrl.Request{{
		NamespacedName: req,
	},
	}
}

func (r *ClusterReconciler) createTenantInstance(cluster *clusterv1alpha1.Cluster) *monitoringv1alpha1.Tenant {

	label := make(map[string]string, 1)
	label[constants.ServiceLabelKey] = r.KubesphereAdapterDefaultService
	label[constants.StorageLabelKey] = constants.DefaultStorage
	return &monitoringv1alpha1.Tenant{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cluster.Name,
			Labels: label,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: cluster.APIVersion,
					Kind:       cluster.Kind,
					Name:       cluster.Name,
					UID:        cluster.UID,
					Controller: pointer.BoolPtr(true),
				},
			},
		},
		Spec: monitoringv1alpha1.TenantSpec{
			Tenant: cluster.Name,
		},
	}
}

type ResourceCustomPredicate struct {
	predicate.Funcs
}

// Update ignore cluster update event
func (r *ResourceCustomPredicate) Update(e event.UpdateEvent) bool {
	return false
}
