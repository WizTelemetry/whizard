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

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/options"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/ruler"
)

// ThanosRulerReconciler reconciles a ThanosRuler object
type ThanosRulerReconciler struct {
	DefaulterValidator   ThanosRulerDefaulterValidator
	ReloaderConfig       options.PrometheusConfigReloaderConfig
	RulerQueryProxyConfig options.RulerQueryProxyConfig
	client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosrulers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosrulers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosrulers/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=alertingrules,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=rulegroups,verbs=get;list;watch
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services;configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ThanosRulerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("thanosruler", req.NamespacedName)

	l.Info("sync")

	instance := &monitoringv1alpha1.ThanosRuler{}
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

	if err := ruler.New(baseReconciler, instance, r.ReloaderConfig, r.RulerQueryProxyConfig).Reconcile(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ThanosRulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.ThanosRuler{}).
		Watches(&source.Kind{Type: &promv1.PrometheusRule{}},
			handler.EnqueueRequestsFromMapFunc(r.mapRuleToThanosRulerFunc)).
		Watches(&source.Kind{Type: &monitoringv1alpha1.AlertingRule{}},
			handler.EnqueueRequestsFromMapFunc(r.mapRuleToThanosRulerFunc)).
		Watches(&source.Kind{Type: &monitoringv1alpha1.RuleGroup{}},
			handler.EnqueueRequestsFromMapFunc(r.mapRuleGroupToThanosRulerFunc)).
		Watches(&source.Kind{Type: &monitoringv1alpha1.Service{}},
			handler.EnqueueRequestsFromMapFunc(r.mapToThanosRulerFunc)).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func (r *ThanosRulerReconciler) mapRuleToThanosRulerFunc(o client.Object) []reconcile.Request {
	var ns corev1.Namespace
	if err := r.Client.Get(r.Context, types.NamespacedName{Name: o.GetNamespace()}, &ns); err != nil {
		if apierrors.IsNotFound(err) {
			log.FromContext(r.Context).WithValues("namespace", o.GetNamespace()).Error(err, "")
		}
		return nil
	}

	var thanosRulerList monitoringv1alpha1.ThanosRulerList
	if err := r.Client.List(r.Context, &thanosRulerList); err != nil {
		log.FromContext(r.Context).WithValues("thanosrulerlist", "").Error(err, "")
		return nil
	}

	var reqs []reconcile.Request
	for _, ruler := range thanosRulerList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: ruler.Namespace,
				Name:      ruler.Name,
			},
		}
		if ruler.Namespace == ns.Name {
			reqs = append(reqs, req)
			continue
		}

		ruleNsSelector, err := metav1.LabelSelectorAsSelector(ruler.Spec.RuleNamespaceSelector)
		if err != nil {
			log.FromContext(r.Context).WithValues("thanosruler", req.NamespacedName).Error(
				err, "failed to convert RuleNamespaceSelector")
			continue
		}
		if ruleNsSelector.Matches(labels.Set(ns.Labels)) {
			reqs = append(reqs, req)
			continue
		}

		alertingRuleNsSelector, err := metav1.LabelSelectorAsSelector(ruler.Spec.AlertingRuleNamespaceSelector)
		if err != nil {
			log.FromContext(r.Context).WithValues("thanosruler", req.NamespacedName).Error(
				err, "failed to convert AlertingRuleNamespaceSelector")
			continue
		}
		if alertingRuleNsSelector.Matches(labels.Set(ns.Labels)) {
			reqs = append(reqs, req)
		}
	}

	return reqs
}

func (r *ThanosRulerReconciler) mapRuleGroupToThanosRulerFunc(o client.Object) []reconcile.Request {
	var alertingRuleList monitoringv1alpha1.AlertingRuleList
	if err := r.Client.List(r.Context, &alertingRuleList,
		client.InNamespace(o.GetNamespace()),
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByRuleGroup(o))); err != nil {
		log.FromContext(r.Context).WithValues("thanosrulerlist", "").Error(err, "")
		return nil
	}

	var reqs []reconcile.Request
	for _, rule := range alertingRuleList.Items {
		reqs = append(reqs, r.mapRuleToThanosRulerFunc(&rule)...)
	}
	return reqs
}

func (r *ThanosRulerReconciler) mapToThanosRulerFunc(o client.Object) []reconcile.Request {

	var thanosRulerList monitoringv1alpha1.ThanosRulerList
	if err := r.Client.List(r.Context, &thanosRulerList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(o))); err != nil {
		log.FromContext(r.Context).WithValues("thanosrulerlist", "").Error(err, "")
		return nil
	}

	var reqs []reconcile.Request
	for _, ruler := range thanosRulerList.Items {
		reqs = append(reqs, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: ruler.Namespace,
				Name:      ruler.Name,
			},
		})
	}

	return reqs
}

type ThanosRulerDefaulterValidator func(ruler *monitoringv1alpha1.ThanosRuler) (*monitoringv1alpha1.ThanosRuler, error)

func CreateThanosRulerDefaulterValidator(opt options.Options) ThanosRulerDefaulterValidator {
	var replicas int32 = 1

	return func(ruler *monitoringv1alpha1.ThanosRuler) (*monitoringv1alpha1.ThanosRuler, error) {

		if ruler.Spec.Image == "" {
			ruler.Spec.Image = opt.ThanosImage
		}
		if ruler.Spec.Replicas == nil || *ruler.Spec.Replicas < 0 {
			ruler.Spec.Replicas = &replicas
		}

		return ruler, nil
	}
}
