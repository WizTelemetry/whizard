package tenant

import (
	"context"
	"fmt"
	"strings"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (t *Tenant) store() error {

	storeList := &monitoringv1alpha1.StoreList{}
	if err := t.Client.List(t.Context, storeList); err != nil {
		return err
	}

	currentStores := make(map[string]interface{})
	for _, item := range storeList.Items {
		store := item
		if store.DeletionTimestamp != nil || !store.DeletionTimestamp.IsZero() {
			continue
		}

		if store.Labels == nil {
			continue
		}

		serviceNamespacedName := store.Labels[constants.ServiceLabelKey]
		if serviceNamespacedName == "" {
			klog.V(3).Infof("Store [%s.%s] Service mismatch", store.Namespace, store.Name)
			continue
		}

		storageNamespacedName := store.Labels[constants.StorageLabelKey]
		if storageNamespacedName == "" {
			klog.V(3).Infof("Store [%s.%s] Storage mismatch", store.Namespace, store.Name)
			continue
		}

		if currentStores[serviceNamespacedName+storageNamespacedName] != nil {
			if err := t.Client.Delete(t.Context, &store); err != nil {
				return err
			}
		} else {
			currentStores[serviceNamespacedName+storageNamespacedName] = &store
		}
	}

	expectStores, err := sortTenantsByStorageAndService(t.Context, t.Client)
	if err != nil {
		return err
	}

	for k, v := range currentStores {
		if _, ok := expectStores[k]; ok {
			store := v.(*monitoringv1alpha1.Store)

			if store.Annotations == nil {
				store.Annotations = make(map[string]string)
			}

			tenantHash, err := resources.GetTenantHash(t.Context, t.Client, map[string]string{
				constants.StorageLabelKey: store.Labels[constants.StorageLabelKey],
				constants.ServiceLabelKey: store.Labels[constants.ServiceLabelKey],
			})
			if err != nil {
				return err
			}

			storageHash, err := resources.GetStorageHash(t.Context, t.Client, store.Labels[constants.StorageLabelKey])
			if err != nil {
				return err
			}

			needUpdate := false
			if store.Annotations[constants.LabelNameTenantHash] != tenantHash {
				store.Annotations[constants.LabelNameTenantHash] = tenantHash
				needUpdate = true
			}

			if store.Annotations[constants.LabelNameStorageHash] != tenantHash {
				store.Annotations[constants.LabelNameStorageHash] = storageHash
				needUpdate = true
			}

			if needUpdate {
				if err := t.Client.Update(t.Context, store); err != nil {
					return err
				}
			}

			delete(currentStores, k)
			delete(expectStores, k)
		}
	}

	for _, v := range currentStores {
		store := v.(*monitoringv1alpha1.Store)
		if err := t.Client.Delete(t.Context, store); err != nil {
			return err
		}

		klog.V(3).Infof("Delete Store[%s.%s]", store.Namespace, store.Name)
	}

	for _, v := range expectStores {
		m := v.(map[string]string)
		serviceNamespacedName := strings.Split(m[constants.ServiceLabelKey], ".")
		storageNamespacedName := strings.Split(m[constants.StorageLabelKey], ".")
		store := &monitoringv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: fmt.Sprintf("%s-%s-%s-", constants.AppNameStore, serviceNamespacedName[1], storageNamespacedName[1]),
				Namespace:    serviceNamespacedName[0],
				Labels:       m,
				Annotations:  map[string]string{},
			},
		}

		tenantHash, err := resources.GetTenantHash(t.Context, t.Client, map[string]string{
			constants.StorageLabelKey: store.Labels[constants.StorageLabelKey],
			constants.ServiceLabelKey: store.Labels[constants.ServiceLabelKey],
		})
		if err != nil {
			return err
		}

		storageHash, err := resources.GetStorageHash(t.Context, t.Client, store.Labels[constants.StorageLabelKey])
		if err != nil {
			return err
		}

		store.Annotations[constants.LabelNameTenantHash] = tenantHash
		store.Annotations[constants.LabelNameStorageHash] = storageHash

		if err := t.Client.Create(t.Context, store); err != nil {
			return err
		}

		klog.V(3).Infof("Create store[%s.%s] for Service[%s] Storage[%s]",
			store.Namespace, store.Name,
			m[constants.ServiceLabelKey],
			m[constants.StorageLabelKey])
	}

	return nil
}

func sortTenantsByStorageAndService(ctx context.Context, c client.Client) (map[string]interface{}, error) {
	tenantList := &monitoringv1alpha1.TenantList{}
	if err := c.List(ctx, tenantList); err != nil {
		return nil, err
	}

	storageMap := make(map[string]interface{})
	for _, tenant := range tenantList.Items {

		if tenant.DeletionTimestamp != nil || !tenant.DeletionTimestamp.IsZero() {
			klog.V(3).Infof("ignore tenant %s is deleting", tenant.Name)
			continue
		}

		if tenant.Labels[constants.StorageLabelKey] == constants.LocalStorage {
			klog.V(3).Infof("ignore tenant %s with local storage", tenant.Name)
			continue
		}

		serviceNamespacedName := tenant.Labels[constants.ServiceLabelKey]
		if serviceNamespacedName == "" {
			klog.V(3).Infof("ignore tenant %s without service", tenant.Name)
			continue
		}

		storageNamespacedName := ""
		if tenant.Spec.Storage != nil {
			storageNamespacedName = fmt.Sprintf("%s.%s", tenant.Spec.Storage.Namespace, tenant.Spec.Storage.Name)
		} else {
			service := &monitoringv1alpha1.Service{}
			serviceNamespacedName := strings.Split(serviceNamespacedName, ".")
			if err := c.Get(ctx, types.NamespacedName{
				Namespace: serviceNamespacedName[0],
				Name:      serviceNamespacedName[1],
			}, service); err != nil {
				klog.V(3).Infof("get service %s failed, %s", serviceNamespacedName, err)
				return nil, err
			}
			if service.Spec.Storage != nil {
				storageNamespacedName = fmt.Sprintf("%s.%s", service.Spec.Storage.Namespace, service.Spec.Storage.Name)
			}
		}

		if storageNamespacedName == "" {
			klog.V(3).Infof("Tenant [%s] Storage mismatch", tenant.Name)
			continue
		}

		storageMap[serviceNamespacedName+storageNamespacedName] = map[string]string{
			constants.ServiceLabelKey: serviceNamespacedName,
			constants.StorageLabelKey: storageNamespacedName,
		}
	}

	return storageMap, nil
}
