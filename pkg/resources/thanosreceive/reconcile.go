package thanosreceive

import (
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ThanosReceive) reconcileConfigMaps() error {
	configmaps, err := r.configMaps()
	if err != nil {
		return err
	}

	for _, desired := range configmaps {
		if err := ctrl.SetControllerReference(&r.Instance, desired, r.Scheme); err != nil {
			return err
		}
		current := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: desired.Namespace,
				Name:      desired.Name,
			},
		}
		_, err := ctrl.CreateOrUpdate(r.Context, r.Client, current, func() error {
			current.Labels = desired.Labels
			current.Annotations = desired.Annotations
			current.OwnerReferences = desired.OwnerReferences
			current.Data = desired.Data
			current.BinaryData = desired.BinaryData
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ThanosReceive) reconcileServices() error {
	desired := r.Service()
	if err := ctrl.SetControllerReference(&r.Instance, &desired, r.Scheme); err != nil {
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
	return err
}

func (r *ThanosReceive) reconcileIngresses() error {
	ingresses := r.Ingresses()

	for _, desired := range ingresses {
		if err := ctrl.SetControllerReference(&r.Instance, desired, r.Scheme); err != nil {
			return err
		}
		current := &netv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: desired.Namespace,
				Name:      desired.Name,
			},
		}
		_, err := ctrl.CreateOrUpdate(r.Context, r.Client, current, func() error {
			current.Labels = desired.Labels
			current.OwnerReferences = desired.OwnerReferences
			current.Spec = desired.Spec
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil

}

func (r *ThanosReceive) reconcileDeployments() error {
	desired := r.deployment()
	if err := ctrl.SetControllerReference(&r.Instance, &desired, r.Scheme); err != nil {
		return err
	}
	current := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: desired.Namespace,
			Name:      desired.Name,
		},
	}
	_, err := ctrl.CreateOrUpdate(r.Context, r.Client, current, func() error {
		current.Labels = desired.Labels
		current.OwnerReferences = desired.OwnerReferences
		current.Spec = desired.Spec
		return nil
	})
	return err
}

func (r *ThanosReceive) reconcileStatefulSets() error {
	desired := r.statefulSet()

	if err := ctrl.SetControllerReference(&r.Instance, &desired, r.Scheme); err != nil {
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

	return err
}

func (r *ThanosReceive) Reconcile() error {
	if err := r.reconcileConfigMaps(); err != nil {
		return err
	}

	if err := r.reconcileServices(); err != nil {
		return err
	}

	switch r.GetMode() {
	case RouterOnly:
		if err := r.reconcileDeployments(); err != nil {
			return err
		}
	case IngestorOnly, RouterIngestor:
		if err := r.reconcileStatefulSets(); err != nil {
			return err
		}
	default:
		return errors.New("unsupported thanos receive mode")
	}

	if err := r.reconcileIngresses(); err != nil {
		return err
	}

	return nil
}
