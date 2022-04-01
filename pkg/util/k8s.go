package util

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateOrUpdateDeployment(ctx context.Context, cli client.Client, desired *appsv1.Deployment) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var current = &appsv1.Deployment{}
		err := cli.Get(ctx, client.ObjectKeyFromObject(desired), current)
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			return cli.Create(ctx, desired)
		}

		mergeMetadata(&desired.ObjectMeta, &current.ObjectMeta)

		return cli.Update(ctx, desired)
	})
}

func CreateOrUpdateConfigMap(ctx context.Context, cli client.Client, desired *corev1.ConfigMap) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var current = &appsv1.Deployment{}
		err := cli.Get(ctx, client.ObjectKeyFromObject(desired), current)
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			return cli.Create(ctx, desired)
		}

		mergeMetadata(&desired.ObjectMeta, &current.ObjectMeta)

		if apiequality.Semantic.DeepEqual(current, desired) {
			return nil
		}

		return cli.Update(ctx, desired)
	})
}

func CreateOrUpdate(ctx context.Context, cli client.Client, desired client.Object) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var current = desired.DeepCopyObject().(client.Object)
		err := cli.Get(ctx, client.ObjectKeyFromObject(desired), current)
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			return cli.Create(ctx, desired)
		}

		return cli.Update(ctx, desired)
	})
}

func mergeMetadata(new, old *metav1.ObjectMeta) {
	new.ResourceVersion = old.ResourceVersion
	new.Labels = labels.Merge(old.Labels, new.Labels)
	new.Annotations = labels.Merge(old.Annotations, new.Annotations)
}

// IndexOwnerRef returns the index of the owner reference in the slice if found, or -1.
func IndexOwnerRef(ownerReferences []metav1.OwnerReference, ref metav1.OwnerReference) int {
	for index, r := range ownerReferences {
		if referSameObject(r, ref) {
			return index
		}
	}
	return -1
}

// Returns true if a and b point to the same object.
func referSameObject(a, b metav1.OwnerReference) bool {
	aGV, err := schema.ParseGroupVersion(a.APIVersion)
	if err != nil {
		return false
	}

	bGV, err := schema.ParseGroupVersion(b.APIVersion)
	if err != nil {
		return false
	}

	return aGV.Group == bGV.Group && a.Kind == b.Kind && a.Name == b.Name
}
