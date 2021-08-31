package thanosquery

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (q *ThanosQuery) reconcileConfigMaps() error {
	configmaps := q.configMaps()

	for _, desired := range configmaps {
		if err := ctrl.SetControllerReference(&q.Instance, desired, q.Scheme); err != nil {
			return err
		}
		current := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: desired.Namespace,
				Name:      desired.Name,
			},
		}
		_, err := ctrl.CreateOrUpdate(q.Context, q.Client, current, func() error {
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

func (q *ThanosQuery) reconcileIngresses() error {
	ingresses := q.ingresses()

	for _, desired := range ingresses {
		if err := ctrl.SetControllerReference(&q.Instance, desired, q.Scheme); err != nil {
			return err
		}
		current := &netv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: desired.Namespace,
				Name:      desired.Name,
			},
		}
		_, err := ctrl.CreateOrUpdate(q.Context, q.Client, current, func() error {
			current.Labels = desired.Labels
			current.Annotations = desired.Annotations
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

func (q *ThanosQuery) reconcileService() error {
	desired := q.service()
	if err := ctrl.SetControllerReference(&q.Instance, &desired, q.Scheme); err != nil {
		return err
	}
	current := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: desired.Namespace,
			Name:      desired.Name,
		},
	}
	_, err := ctrl.CreateOrUpdate(q.Context, q.Client, current, func() error {
		current.Labels = desired.Labels
		current.OwnerReferences = desired.OwnerReferences
		current.Spec.Selector = desired.Spec.Selector
		current.Spec.Ports = desired.Spec.Ports
		current.Spec.ClusterIP = desired.Spec.ClusterIP
		return nil
	})
	return err
}

func (q *ThanosQuery) reconcileDeployment() error {
	desired := q.deployment()
	if err := ctrl.SetControllerReference(&q.Instance, &desired, q.Scheme); err != nil {
		return err
	}
	current := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: desired.Namespace,
			Name:      desired.Name,
		},
	}
	_, err := ctrl.CreateOrUpdate(q.Context, q.Client, current, func() error {
		current.Labels = desired.Labels
		current.OwnerReferences = desired.OwnerReferences
		current.Spec = desired.Spec
		return nil
	})
	return err
}

func (q *ThanosQuery) Reconcile() error {

	if err := q.reconcileConfigMaps(); err != nil {
		return err
	}

	if err := q.reconcileDeployment(); err != nil {
		return err
	}

	if err := q.reconcileService(); err != nil {
		return err
	}

	if err := q.reconcileIngresses(); err != nil {
		return err
	}

	return nil
}
