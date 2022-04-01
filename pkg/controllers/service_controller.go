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
	"fmt"
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	monitoringv1alpha1 "github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/config"
	"github.com/kubesphere/paodin-monitoring/pkg/resources"
	"github.com/kubesphere/paodin-monitoring/pkg/resources/compact"
	"github.com/kubesphere/paodin-monitoring/pkg/resources/gateway"
	"github.com/kubesphere/paodin-monitoring/pkg/resources/query"
	"github.com/kubesphere/paodin-monitoring/pkg/resources/receive"
	"github.com/kubesphere/paodin-monitoring/pkg/resources/storegateway"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	DefaulterValidator ServiceDefaulterValidator
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.paodin.io,resources=services/finalizers,verbs=update
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

	l.Info("sync service")
	_ = sync.Once{}
	instance := &monitoringv1alpha1.Service{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	instance.Spec, err = r.DefaulterValidator(instance.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	serviceBaseReconciler := resources.ServiceBaseReconciler{
		Client:  r.Client,
		Log:     l,
		Scheme:  r.Scheme,
		Context: ctx,
		Service: instance,
	}

	var reconciles []func() error
	if instance.Spec.Thanos == nil {
		serviceBaseReconciler.Service.Spec.Thanos = &monitoringv1alpha1.Thanos{}
	}
	reconciles = append(reconciles, compact.New(serviceBaseReconciler).Reconcile)
	reconciles = append(reconciles, storegateway.New(serviceBaseReconciler).Reconcile)
	reconciles = append(reconciles, receive.New(serviceBaseReconciler).Reconcile)
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
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

type ServiceDefaulterValidator func(spec monitoringv1alpha1.ServiceSpec) (monitoringv1alpha1.ServiceSpec, error)

func CreateServiceDefaulterValidator(cfg config.Config) ServiceDefaulterValidator {
	var replicas int32 = 1
	var applyDefaultFields = func(defaultFields,
		fields monitoringv1alpha1.CommonThanosFields) monitoringv1alpha1.CommonThanosFields {
		if fields.Image == "" {
			fields.Image = defaultFields.Image
		}
		if fields.LogLevel == "" {
			fields.LogLevel = defaultFields.LogLevel
		}
		if fields.LogFormat == "" {
			fields.LogFormat = defaultFields.LogFormat
		}
		return fields
	}

	return func(spec monitoringv1alpha1.ServiceSpec) (monitoringv1alpha1.ServiceSpec, error) {
		if spec.Thanos == nil {
			return spec, nil
		}

		var thanos = spec.Thanos

		if thanos.DefaultFields.Image == "" {
			thanos.DefaultFields.Image = cfg.ThanosDefaultImage
		}

		if thanos.Query != nil {
			thanos.Query.CommonThanosFields = applyDefaultFields(thanos.DefaultFields, thanos.Query.CommonThanosFields)
			if thanos.Query.Replicas == nil || *thanos.Query.Replicas < 0 {
				thanos.Query.Replicas = &replicas
			}
			if thanos.Query.Envoy.Image == "" {
				thanos.Query.Envoy.Image = cfg.EnvoyDefaultImage
			}

		}
		if thanos.Receive != nil {
			thanos.Receive.Router.CommonThanosFields = applyDefaultFields(thanos.DefaultFields, thanos.Receive.Router.CommonThanosFields)
			if thanos.Receive.Router.Replicas == nil || *thanos.Receive.Router.Replicas < 0 {
				thanos.Receive.Router.Replicas = &replicas
			}
			var ingestors []monitoringv1alpha1.ReceiveIngestor
			for _, i := range thanos.Receive.Ingestors {
				ingestor := i
				if ingestor.Name == "" {
					return spec, fmt.Errorf("ingestor->name can not empty")
				}
				ingestor.CommonThanosFields = applyDefaultFields(thanos.DefaultFields, ingestor.CommonThanosFields)
				if ingestor.Replicas == nil || *ingestor.Replicas < 0 {
					ingestor.Replicas = &replicas
				}
				ingestors = append(ingestors, ingestor)
			}
			thanos.Receive.Ingestors = ingestors
		}
		if thanos.StoreGateway != nil {
			thanos.StoreGateway.CommonThanosFields = applyDefaultFields(thanos.DefaultFields, thanos.StoreGateway.CommonThanosFields)
			if thanos.StoreGateway.Replicas == nil || *thanos.StoreGateway.Replicas < 0 {
				thanos.StoreGateway.Replicas = &replicas
			}
		}
		if thanos.Compact != nil {
			thanos.Compact.CommonThanosFields = applyDefaultFields(thanos.DefaultFields, thanos.Compact.CommonThanosFields)
			if thanos.Compact.Replicas == nil || *thanos.Compact.Replicas < 0 {
				thanos.Compact.Replicas = &replicas
			}
		}

		spec.Thanos = thanos

		return spec, nil
	}
}
