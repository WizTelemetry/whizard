package tenant

import (
	"fmt"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/util"
)

func (t *Tenant) compactor() error {

	// finalizers check,  when tenant cr is deleted, ObjectMeta.GetDeletionTimestamp() is not nil, remove finalizers field and call removeTenantFromCompactorByName()
	if t.tenant.ObjectMeta.GetDeletionTimestamp().IsZero() {
		if !containsString(t.tenant.ObjectMeta.Finalizers, constants.FinalizerCompactor) {
			t.tenant.ObjectMeta.Finalizers = append(t.tenant.ObjectMeta.Finalizers, constants.FinalizerCompactor)
			return t.Client.Update(t.Context, t.tenant)
		}
	} else {
		if containsString(t.tenant.ObjectMeta.Finalizers, constants.FinalizerCompactor) {
			if t.tenant.Status.Compactor != nil {
				if err := t.removeTenantFromCompactorByName(t.tenant.Status.Compactor.Namespace, t.tenant.Status.Compactor.Name); err != nil {
					return err
				}
				t.tenant.Status.Compactor = nil
			}
			t.tenant.ObjectMeta.Finalizers = removeString(t.tenant.Finalizers, constants.FinalizerCompactor)
			return t.Client.Update(t.Context, t.tenant)
		}
	}

	// Check if compactor needs to be created or reset
	if needResetCompactor, err := t.needResetCompactor(); err != nil {
		return err
	} else if needResetCompactor {
		klog.V(3).Infof("Reset compactor [%s] for tenant [%s]", t.tenant.Status.Compactor.Name, t.tenant.Name)
		if err := t.removeTenantFromCompactorByName(t.tenant.Status.Compactor.Namespace, t.tenant.Status.Compactor.Name); err != nil {
			return err
		}

		return t.Client.Status().Update(t.Context, t.tenant)
	}

	// when tenant.Labels don't contain Service, remove the bindings to compactor
	if v, ok := t.tenant.Labels[constants.ServiceLabelKey]; !ok || v == "" {
		klog.V(3).Infof("Tenant [%s]'s Service is empty. compactor does not need to be created", t.tenant.Name)
		if t.tenant.Status.Compactor != nil {
			err := t.removeTenantFromCompactorByName(t.tenant.Status.Compactor.Namespace, t.tenant.Status.Compactor.Name)
			if err != nil {
				return err
			}
			return t.Client.Status().Update(t.Context, t.tenant)
		}
		return nil
	}

	// when tenant.Labels don't contain Storage, remove the bindings to compactor
	if t.GetStorage(t.tenant.Labels[constants.StorageLabelKey]) == constants.LocalStorage {
		klog.V(3).Infof("Tenant [%s]'s Storage is empty. compactor does not need to be created", t.tenant.Name)
		return nil
	}

	if v, ok := t.tenant.Labels[constants.ExclusiveLabelKey]; ok && v == "true" {
		compactor := t.createExclusiveCompactorInstance(t.tenant)

		if err := util.CreateOrUpdate(t.Context, t.Client, compactor); err != nil {
			return err
		}

		t.tenant.Status.Compactor = &monitoringv1alpha1.ObjectReference{
			Namespace: compactor.Namespace,
			Name:      compactor.Name,
		}
		return t.Client.Status().Update(t.Context, t.tenant)
	}

	var compactorList monitoringv1alpha1.CompactorList
	ls := make(map[string]string, 2)
	ls[constants.ServiceLabelKey] = t.tenant.Labels[constants.ServiceLabelKey]
	ls[constants.StorageLabelKey] = t.GetStorage(t.tenant.Labels[constants.StorageLabelKey])
	serviceNamespacedName := strings.Split(t.tenant.Labels[constants.ServiceLabelKey], ".")
	err := t.Client.List(t.Context, &compactorList, &client.ListOptions{
		Namespace:     serviceNamespacedName[0],
		LabelSelector: labels.SelectorFromSet(ls),
	})
	if err != nil {
		return err
	}

	// Check duplicates
	for _, item := range compactorList.Items {
		if containsString(item.Spec.Tenants, t.tenant.Spec.Tenant) {
			klog.V(3).Infof("Compactor [%s] has tenant [%s] ,update status ", item.Name, t.tenant.Name)

			t.tenant.Status.Compactor = &monitoringv1alpha1.ObjectReference{
				Namespace: item.Namespace,
				Name:      item.Name,
			}

			return t.Client.Status().Update(t.Context, t.tenant)
		}
	}

	needToCreate := true
	compactor := &monitoringv1alpha1.Compactor{}
	for _, item := range compactorList.Items {
		if len(item.Spec.Tenants) < t.Service.Spec.CompactorTemplateSpec.DefaultTenantsPerCompactor {
			compactor = &item
			compactor.Spec.Tenants = append(compactor.Spec.Tenants, t.tenant.Name)
			needToCreate = false
			break
		}
	}

	if needToCreate {
		compactor = t.createCompactorInstance(t.tenant)
	}

	t.tenant.Status.Compactor = &monitoringv1alpha1.ObjectReference{
		Namespace: compactor.Namespace,
		Name:      compactor.Name,
	}

	if err := util.CreateOrUpdate(t.Context, t.Client, compactor); err != nil {
		return err
	}

	klog.V(3).Infof("create new compactor %s for tenant %s", compactor.Name, t.tenant.Name)

	return t.Client.Status().Update(t.Context, t.tenant)
}

func (t *Tenant) needResetCompactor() (bool, error) {
	if t.tenant.Status.Compactor == nil {
		return false, nil
	}

	compactor := &monitoringv1alpha1.Compactor{}
	err := t.Client.Get(t.Context, types.NamespacedName{
		Namespace: t.tenant.Status.Compactor.Namespace,
		Name:      t.tenant.Status.Compactor.Name,
	}, compactor)
	if err != nil {
		if apierrors.IsNotFound(err) {
			t.Log.V(1).Info("Cannot find compactor [%s] for tenant [%s], create one", t.tenant.Status.Compactor.Name, t.tenant.Name)
			return true, nil
		} else {
			return false, err
		}
	}

	if ok := containsString(compactor.Spec.Tenants, t.tenant.Spec.Tenant); !ok {
		klog.V(3).Infof("Tenant [%s] and compactor [%s] mismatch, need to reset compactor", t.tenant.Name, compactor.Name)
		return true, nil
	}

	if v, ok := compactor.Labels[constants.ServiceLabelKey]; !ok || v != t.tenant.Labels[constants.ServiceLabelKey] {
		klog.V(3).Infof("Tenant [%s] and compactor [%s]'s Service mismatch, need to reset compactor", t.tenant.Name, compactor.Name)
		return true, nil
	}

	if v, ok := compactor.Labels[constants.StorageLabelKey]; !ok || v != t.GetStorage(t.tenant.Labels[constants.StorageLabelKey]) {
		klog.V(3).Infof("Tenant [%s] and compactor [%s]'s Storage mismatch, need to reset compactor", t.tenant.Name, compactor.Name)
		return true, nil
	}

	if v, ok := t.tenant.Labels[constants.ExclusiveLabelKey]; ok && v == "true" {
		if v, ok := compactor.Labels[constants.ExclusiveLabelKey]; !ok || v != "true" || compactor.Name != t.tenant.Name {
			klog.V(3).Infof("Tenant [%s] requires its exclusive compactor, the current shared compactor [%s] should be replaced by an exclusive one", t.tenant.Name, compactor.Name)
			return true, nil
		}
	} else {
		if v, ok := compactor.Labels[constants.ExclusiveLabelKey]; ok && v == "true" && compactor.Name == t.tenant.Name {
			klog.V(3).Infof("Tenant [%s] is not allowed to have its exclusive compactor, the current exclusive compactor [%s] will be replaced by a shared compactor", t.tenant.Name, compactor.Name)
			return true, nil
		}
	}
	return false, nil
}

func (t *Tenant) createCompactorInstance(tenant *monitoringv1alpha1.Tenant) *monitoringv1alpha1.Compactor {
	storage := t.GetStorage(tenant.Labels[constants.StorageLabelKey])

	label := make(map[string]string, 2)
	label[constants.ServiceLabelKey] = tenant.Labels[constants.ServiceLabelKey]
	label[constants.StorageLabelKey] = storage

	serviceNamespacedName := strings.Split(tenant.Labels[constants.ServiceLabelKey], ".")
	storageNamespacedName := strings.Split(storage, ".")
	return &monitoringv1alpha1.Compactor{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-%s-", serviceNamespacedName[1], storageNamespacedName[1]),
			Namespace:    serviceNamespacedName[0],
			Labels:       label,
		},
		Spec: monitoringv1alpha1.CompactorSpec{
			Tenants: []string{tenant.Spec.Tenant},
		},
	}
}

func (t *Tenant) createExclusiveCompactorInstance(tenant *monitoringv1alpha1.Tenant) *monitoringv1alpha1.Compactor {
	storage := t.GetStorage(tenant.Labels[constants.StorageLabelKey])

	label := make(map[string]string, 3)
	label[constants.ServiceLabelKey] = tenant.Labels[constants.ServiceLabelKey]
	label[constants.StorageLabelKey] = storage
	label[constants.ExclusiveLabelKey] = "true"

	serviceNamespacedName := strings.Split(tenant.Labels[constants.ServiceLabelKey], ".")
	return &monitoringv1alpha1.Compactor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      t.tenant.Name,
			Namespace: serviceNamespacedName[0],
			Labels:    label,
		},
		Spec: monitoringv1alpha1.CompactorSpec{
			Tenants: []string{tenant.Spec.Tenant},
		},
	}
}

func (t *Tenant) removeTenantFromCompactorByName(namespace, name string) error {
	if t.tenant.Status.Compactor != nil {
		t.tenant.Status.Compactor = nil
	}

	compactor := &monitoringv1alpha1.Compactor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := t.Client.Get(t.Context, client.ObjectKeyFromObject(compactor), compactor); err != nil {
		return util.IgnoreNotFound(err)
	}

	if ok := containsString(compactor.Spec.Tenants, t.tenant.Spec.Tenant); ok {
		klog.V(3).Infof("compactor %s update, remove tenant %s", compactor.Name, t.tenant.Name)
		compactor.Spec.Tenants = removeString(compactor.Spec.Tenants, t.tenant.Spec.Tenant)
	}

	if len(compactor.Spec.Tenants) == 0 {
		return util.IgnoreNotFound(t.Client.Delete(t.Context, compactor))
	} else {
		return util.CreateOrUpdate(t.Context, t.Client, compactor)
	}
}
