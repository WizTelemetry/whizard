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

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
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

	v1alpha1 "github.com/kubesphere/paodin-monitoring/api/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/config"
	"github.com/kubesphere/paodin-monitoring/pkg/resources/thanosreceive"
)

// ThanosReceiveReconciler reconciles a ThanosReceive object
type ThanosReceiveReconciler struct {
	Cfg config.Config
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosreceives,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosreceives/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=thanosreceives/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=namespaces;endpoints,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services;configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ThanosReceive object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ThanosReceiveReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx, "thanosreceive", req.NamespacedName)

	instance := &v1alpha1.ThanosReceive{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	tr := thanosreceive.ThanosReceive{
		Cfg:      r.Cfg,
		Client:   r.Client,
		Instance: *instance,
		Log:      l,
		Scheme:   r.Scheme,
		Context:  ctx,
	}

	if err := tr.Reconcile(); err != nil {
		if apierrors.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ThanosReceiveReconciler) SetupWithManager(mgr ctrl.Manager) error {

	receivesToRequests := func(receives []v1alpha1.ThanosReceive) []reconcile.Request {
		var requests []reconcile.Request
		for _, receive := range receives {
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{Namespace: receive.Namespace, Name: receive.Name}})
		}
		return requests
	}
	endpointsRequestMapper := handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		if eps, ok := object.(*corev1.Endpoints); !ok {
			return nil
		} else {
			var ctx = context.TODO()
			if receives, err := mapEndpointsToReceives(eps, r.Client, ctx); err != nil {
				log.FromContext(ctx).Error(err, "")
				return nil
			} else {
				return receivesToRequests(receives)
			}
		}
	})
	namespaceRequestMapper := handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		if ns, ok := object.(*corev1.Namespace); !ok {
			return nil
		} else {
			var ctx = context.TODO()
			if receives, err := mapNamespaceToReceives(ns, r.Client, ctx); err != nil {
				log.FromContext(ctx).Error(err, "")
				return nil
			} else {
				return receivesToRequests(receives)
			}
		}
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ThanosReceive{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&netv1.Ingress{}).
		Owns(&corev1.ConfigMap{}).
		Watches(&source.Kind{Type: &corev1.Endpoints{}}, endpointsRequestMapper).
		Watches(&source.Kind{Type: &corev1.Namespace{}}, namespaceRequestMapper).
		Complete(r)
}

func mapEndpointsToReceives(eps *corev1.Endpoints, client client.Client, ctx context.Context) (
	[]v1alpha1.ThanosReceive, error) {
	var (
		receiveList v1alpha1.ThanosReceiveList
		epsNs       corev1.Namespace
		receives    []v1alpha1.ThanosReceive
	)
	if err := client.List(ctx, &receiveList); err != nil {
		return nil, err
	}
	if err := client.Get(ctx, types.NamespacedName{Name: eps.Namespace}, &epsNs); err != nil {
		return nil, err
	}
	for _, receive := range receiveList.Items {
		if receive.Spec.Router == nil {
			continue
		}

		var hashrings []*v1alpha1.RouterHashringConfig
		if receive.Spec.Router.SoftTenantHashring != nil {
			hashrings = append(hashrings, receive.Spec.Router.SoftTenantHashring)
		}
		hashrings = append(hashrings, receive.Spec.Router.HardTenantHashrings...)

		for _, hashring := range hashrings {
			if hashring.EndpointsNamespaceSelector == nil {
				if epsNs.Name != receive.Namespace {
					continue
				}
			} else if nsSelector, err := metav1.LabelSelectorAsSelector(hashring.EndpointsNamespaceSelector); err != nil {
				return nil, err
			} else if !nsSelector.Matches(labels.Set(epsNs.Labels)) {
				continue
			}
			if epsSelector, err := metav1.LabelSelectorAsSelector(hashring.EndpointsSelector); err != nil {
				return nil, err
			} else if !epsSelector.Matches(labels.Set(eps.Labels)) {
				continue
			}
			receives = append(receives, receive)
			break
		}
	}
	return receives, nil
}

func mapNamespaceToReceives(ns *corev1.Namespace, client client.Client, ctx context.Context) (
	[]v1alpha1.ThanosReceive, error) {
	var (
		receiveList v1alpha1.ThanosReceiveList
		receives    []v1alpha1.ThanosReceive
	)
	if err := client.List(ctx, &receiveList); err != nil {
		return nil, err
	}
	for _, receive := range receiveList.Items {
		if receive.Spec.Router == nil {
			continue
		}

		var hashrings []*v1alpha1.RouterHashringConfig
		if receive.Spec.Router.SoftTenantHashring != nil {
			hashrings = append(hashrings, receive.Spec.Router.SoftTenantHashring)
		}
		hashrings = append(hashrings, receive.Spec.Router.HardTenantHashrings...)

		for _, hashring := range hashrings {
			if hashring.EndpointsNamespaceSelector == nil {
				if ns.Name != receive.Namespace {
					continue
				}
			} else if nsSelector, err := metav1.LabelSelectorAsSelector(hashring.EndpointsNamespaceSelector); err != nil {
				return nil, err
			} else if !nsSelector.Matches(labels.Set(ns.Labels)) {
				continue
			}
			receives = append(receives, receive)
			break
		}
	}
	return receives, nil
}
