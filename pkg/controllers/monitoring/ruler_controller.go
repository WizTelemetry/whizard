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
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/ruler"
	"github.com/kubesphere/whizard/pkg/util"
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
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// RulerReconciler reconciles a Ruler object
type RulerReconciler struct {
	DefaulterValidator RulerDefaulterValidator
	Option             *options.RulerOptions
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=rulers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=rulers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.whizard.io,resources=rulers/finalizers,verbs=update
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

	rulerReconcile, err := ruler.New(baseReconciler, instance, r.Option)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, rulerReconcile.Reconcile()
}

// SetupWithManager sets up the controller with the Manager.
func (r *RulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Ruler{}).
		Watches(&source.Kind{Type: &promv1.PrometheusRule{}},
			handler.EnqueueRequestsFromMapFunc(r.mapRuleToRulerFunc)).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Service{}},
			handler.EnqueueRequestsFromMapFunc(r.mapToRulerFunc)).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

func (r *RulerReconciler) mapRuleToRulerFunc(o client.Object) []reconcile.Request {
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

func (r *RulerReconciler) mapToRulerFunc(o client.Object) []reconcile.Request {

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

type RulerDefaulterValidator func(ruler *monitoringv1alpha1.Ruler) (*monitoringv1alpha1.Ruler, error)

func CreateRulerDefaulterValidator(opt *options.RulerOptions) RulerDefaulterValidator {

	return func(ruler *monitoringv1alpha1.Ruler) (*monitoringv1alpha1.Ruler, error) {

		opt.Override(&ruler.Spec)

		if ruler.Spec.Shards == nil || *ruler.Spec.Shards < 0 {
			ruler.Spec.Shards = opt.Shards
		}

		return ruler, nil
	}
}
