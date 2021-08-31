package thanosstorage

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ThanosStorage) reconcileServices() error {
	for _, desired := range r.services() {
		if err := ctrl.SetControllerReference(&r.Instance, desired, r.Scheme); err != nil {
			return err
		}
		current := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: desired.Namespace,
				Name:      desired.Name,
			},
		}
		_, err := ctrl.CreateOrUpdate(r.Context, r.Client, current, func() error {
			current.Labels = desired.Labels
			current.OwnerReferences = desired.OwnerReferences
			current.Spec.Selector = desired.Spec.Selector
			current.Spec.Ports = desired.Spec.Ports
			current.Spec.ClusterIP = desired.Spec.ClusterIP
			current.Spec.Type = desired.Spec.Type
			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ThanosStorage) reconcileStatefulSets() error {

	for _, desired := range r.statefulSets() {

		if err := ctrl.SetControllerReference(&r.Instance, desired, r.Scheme); err != nil {
			return err
		}
		current := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: desired.Namespace,
				Name:      desired.Name,
			},
		}
		result, err := ctrl.CreateOrUpdate(r.Context, r.Client, current, func() error {
			current.Labels = desired.Labels
			current.OwnerReferences = desired.OwnerReferences
			current.Spec.Replicas = desired.Spec.Replicas
			current.Spec.Selector = desired.Spec.Selector
			current.Spec.Template = desired.Spec.Template
			current.Spec.VolumeClaimTemplates = desired.Spec.VolumeClaimTemplates
			current.Spec.ServiceName = desired.Spec.ServiceName
			return nil
		})

		if result == controllerutil.OperationResultNone && err != nil {
			if sErr, ok := err.(*apierrors.StatusError); ok && sErr.ErrStatus.Code == 422 && sErr.ErrStatus.Reason == metav1.StatusReasonInvalid {
				propagationPolicy := metav1.DeletePropagationForeground
				err = r.Client.Delete(r.Context, current, &client.DeleteOptions{PropagationPolicy: &propagationPolicy})
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ThanosStorage) Reconcile() error {
	if err := r.reconcileServices(); err != nil {
		return err
	}

	if err := r.reconcileStatefulSets(); err != nil {
		return err
	}

	return nil
}
