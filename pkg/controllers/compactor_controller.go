/*
Copyright 2021 The WhizardTelemetry Authors.

This program is free software: you can redistribute it and/or modify
it under the terms of the Server Side Public License, version 1,
as published by MongoDB, Inc.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
Server Side Public License for more details.

You should have received a copy of the Server Side Public License
along with this program. If not, see
<http://www.mongodb.com/licensing/server-side-public-license>.

As a special exception, the copyright holders give permission to link the
code of portions of this program with the OpenSSL library under certain
conditions as described in each individual source file and distribute
linked combinations including the program with the OpenSSL library. You
must comply with the Server Side Public License in all respects for
all of the code used other than as permitted herein. If you modify file(s)
with this exception, you may extend this exception to your version of the
file(s), but you are not obligated to do so. If you do not wish to do so,
delete this exception statement from your version. If you delete this
exception statement from all source files in the program, then also delete
it in the license file.
*/

package controllers

import (
	"context"

	"github.com/imdario/mergo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	monitoringv1alpha1 "github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources/compactor"
	"github.com/WhizardTelemetry/whizard/pkg/util"
)

// CompactorReconciler reconciles a compactor object
type CompactorReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=compactors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=compactors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=compactors/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=services,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=storages,verbs=get;list;watch
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
func (r *CompactorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("compactor", req.NamespacedName)

	l.Info("sync")

	instance := &monitoringv1alpha1.Compactor{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if instance.Labels == nil ||
		instance.Labels[constants.ServiceLabelKey] == "" ||
		instance.Labels[constants.StorageLabelKey] == "" {
		return ctrl.Result{}, nil
	}

	service := &monitoringv1alpha1.Service{}
	if err := r.Get(ctx, *util.ServiceNamespacedName(&instance.ObjectMeta), service); err != nil {
		return ctrl.Result{}, err
	}

	if _, err := r.applyConfigurationFromCompactorTemplateSpec(instance, resources.ApplyDefaults(service).Spec.CompactorTemplateSpec); err != nil {
		return ctrl.Result{}, err
	}

	if instance.GetDeletionTimestamp().IsZero() {
		if len(instance.Finalizers) == 0 {
			instance.Finalizers = append(instance.Finalizers, constants.FinalizerDeletePVC)
		}
	} else {
		if err := util.DeletePVC(r.Context, r.Client, instance); err != nil {
			return ctrl.Result{}, err
		}

		instance.Finalizers = nil
		return ctrl.Result{}, r.Client.Update(r.Context, instance)
	}

	if len(instance.Spec.Tenants) == 0 {
		klog.V(3).Infof("ignore compactor %s/%s because of empty tenants", instance.Name, instance.Namespace)
		return ctrl.Result{}, nil
	}

	baseReconciler := resources.BaseReconciler{
		Client:  r.Client,
		Log:     l,
		Scheme:  r.Scheme,
		Context: ctx,
	}

	compactorReconciler, err := compactor.New(baseReconciler, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, compactorReconciler.Reconcile()
}

// SetupWithManager sets up the controller with the Manager.
func (r *CompactorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Compactor{}).
		Watches(&monitoringv1alpha1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelByService))).
		Watches(&monitoringv1alpha1.Storage{},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelByStorage))).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func (r *CompactorReconciler) mapFuncBySelectorFunc(fn func(metav1.Object) map[string]string) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		compactorList := &monitoringv1alpha1.CompactorList{}
		if err := r.Client.List(r.Context, compactorList, client.MatchingLabels(fn(o))); err != nil {
			log.FromContext(r.Context).WithValues("compactorList", "").Error(err, "")
			return nil
		}

		var reqs []reconcile.Request
		for _, item := range compactorList.Items {
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

func (r *CompactorReconciler) applyConfigurationFromCompactorTemplateSpec(compactor *monitoringv1alpha1.Compactor, compactorTemplateSpec monitoringv1alpha1.CompactorTemplateSpec) (*monitoringv1alpha1.Compactor, error) {

	err := mergo.Merge(&compactor.Spec, compactorTemplateSpec.CompactorSpec)

	return compactor, err
}
