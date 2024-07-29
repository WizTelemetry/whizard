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

	"dario.cat/mergo"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	monitoringv1alpha1 "github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources/ruler"
	"github.com/WhizardTelemetry/whizard/pkg/util"
)

// RulerReconciler reconciles a Ruler object
type RulerReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=rulers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=rulers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=rulers/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=services,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services;configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *RulerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("ruler", req.NamespacedName)

	l.Info("sync")

	instance := &monitoringv1alpha1.Ruler{}
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

	service := &monitoringv1alpha1.Service{}
	if err := r.Get(ctx, *util.ServiceNamespacedName(&instance.ObjectMeta), service); err != nil {
		return ctrl.Result{}, err
	}

	if _, err := r.applyConfigurationFromRulerTemplateSpec(instance, resources.ApplyDefaults(service).Spec.RulerTemplateSpec); err != nil {
		return ctrl.Result{}, err
	}

	baseReconciler := resources.BaseReconciler{
		Client:  r.Client,
		Log:     l,
		Scheme:  r.Scheme,
		Context: ctx,
	}

	rulerReconcile, err := ruler.New(baseReconciler, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, rulerReconcile.Reconcile()
}

// SetupWithManager sets up the controller with the Manager.
func (r *RulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Ruler{}).
		Watches(&promv1.PrometheusRule{},
			handler.EnqueueRequestsFromMapFunc(r.mapRuleToRulerFunc)).
		Watches(&monitoringv1alpha1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.mapToRulerFunc)).
		Watches(&monitoringv1alpha1.Query{},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelBySameService))).
		Watches(&monitoringv1alpha1.QueryFrontend{},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelBySameService))).
		Watches(&monitoringv1alpha1.Router{},
			handler.EnqueueRequestsFromMapFunc(r.mapFuncBySelectorFunc(util.ManagedLabelBySameService))).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

func (r *RulerReconciler) mapFuncBySelectorFunc(fn func(metav1.Object) map[string]string) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		rulerList := &monitoringv1alpha1.RulerList{}
		if err := r.Client.List(r.Context, rulerList, client.MatchingLabels(fn(o))); err != nil {
			log.FromContext(r.Context).WithValues("rulerList", "").Error(err, "")
			return nil
		}

		var reqs []reconcile.Request
		for _, item := range rulerList.Items {
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

func (r *RulerReconciler) mapRuleToRulerFunc(ctx context.Context, o client.Object) []reconcile.Request {
	var ns corev1.Namespace
	if err := r.Client.Get(r.Context, types.NamespacedName{Name: o.GetNamespace()}, &ns); err != nil {
		if apierrors.IsNotFound(err) {
			log.FromContext(r.Context).WithValues("namespace", o.GetNamespace()).Error(err, "")
		}
		return nil
	}

	var rulerList monitoringv1alpha1.RulerList
	if err := r.Client.List(r.Context, &rulerList); err != nil {
		log.FromContext(r.Context).WithValues("rulerlist", "").Error(err, "")
		return nil
	}

	var reqs []reconcile.Request
	for _, item := range rulerList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: item.Namespace,
				Name:      item.Name,
			},
		}
		if item.Namespace == ns.Name {
			reqs = append(reqs, req)
			continue
		}

		ruleNsSelector, err := metav1.LabelSelectorAsSelector(item.Spec.RuleNamespaceSelector)
		if err != nil {
			log.FromContext(r.Context).WithValues("ruler", req.NamespacedName).Error(
				err, "failed to convert RuleNamespaceSelector")
			continue
		}
		if ruleNsSelector.Matches(labels.Set(ns.Labels)) {
			reqs = append(reqs, req)
		}
	}

	return reqs
}

func (r *RulerReconciler) mapToRulerFunc(ctx context.Context, o client.Object) []reconcile.Request {

	var rulerList monitoringv1alpha1.RulerList
	if err := r.Client.List(r.Context, &rulerList,
		client.MatchingLabels(util.ManagedLabelByService(o))); err != nil {
		log.FromContext(r.Context).WithValues("rulerlist", "").Error(err, "")
		return nil
	}

	var reqs []reconcile.Request
	for _, item := range rulerList.Items {
		reqs = append(reqs, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: item.Namespace,
				Name:      item.Name,
			},
		})
	}

	return reqs
}

func (r *RulerReconciler) applyConfigurationFromRulerTemplateSpec(ruler *monitoringv1alpha1.Ruler, rulerTemplateSpec monitoringv1alpha1.RulerTemplateSpec) (*monitoringv1alpha1.Ruler, error) {

	err := mergo.Merge(&ruler.Spec, rulerTemplateSpec.RulerSpec)

	return ruler, err
}
