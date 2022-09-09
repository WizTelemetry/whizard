package tenant

import (
	"fmt"
	"strings"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type elem struct {
	service string
	storage string
	store   *monitoringv1alpha1.Store
}

func (t *Tenant) store() error {

	storeList := &monitoringv1alpha1.StoreList{}
	if err := t.Client.List(t.Context, storeList); err != nil {
		return err
	}

	currentStores := make(map[string]*elem)
	for _, item := range storeList.Items {
		store := item
		if store.DeletionTimestamp != nil || !store.DeletionTimestamp.IsZero() {
			continue
		}

		if store.Labels == nil {
			store.Labels = make(map[string]string)
		}

		e := &elem{
			service: store.Labels[constants.ServiceLabelKey],
			storage: store.Labels[constants.StorageLabelKey],
			store:   &store,
		}

		key := e.service + "|" + e.storage
		if currentStores[key] != nil {
			if err := t.Client.Delete(t.Context, &store); err != nil {
				return err
			}
		} else {
			currentStores[key] = e
		}
	}

	expectStores, err := t.sortTenantsByStorageAndService()
	if err != nil {
		return err
	}

	for k := range expectStores {
		if e, ok := currentStores[k]; ok {
			if e.service != "" && e.storage != "" {
				store := e.store
				if store.Annotations == nil {
					store.Annotations = make(map[string]string)
				}

				tenantHash, err := t.GetTenantHash(map[string]string{
					constants.StorageLabelKey: e.storage,
					constants.ServiceLabelKey: e.service,
				})
				if err != nil {
					return err
				}

				storageHash, err := t.GetStorageHash(e.storage)
				if err != nil {
					return err
				}

				needUpdate := false
				if store.Annotations[constants.LabelNameTenantHash] != tenantHash {
					store.Annotations[constants.LabelNameTenantHash] = tenantHash
					needUpdate = true
				}

				if store.Annotations[constants.LabelNameStorageHash] != storageHash {
					store.Annotations[constants.LabelNameStorageHash] = storageHash
					needUpdate = true
				}

				if needUpdate {
					if err := t.Client.Update(t.Context, store); err != nil {
						return err
					}
				}
			}

			delete(currentStores, k)
			delete(expectStores, k)
		}
	}

	for _, e := range currentStores {
		if err := t.Client.Delete(t.Context, e.store); err != nil {
			return err
		}

		klog.V(3).Infof("Delete Store[%s.%s]", e.store.Namespace, e.store.Name)
	}

	for _, e := range expectStores {
		serviceNamespacedName := strings.Split(e.service, ".")
		storageNamespacedName := strings.Split(e.storage, ".")
		store := &monitoringv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: fmt.Sprintf("%s-%s-", serviceNamespacedName[1], storageNamespacedName[1]),
				Namespace:    serviceNamespacedName[0],
				Labels: map[string]string{
					constants.ServiceLabelKey: e.service,
					constants.StorageLabelKey: e.storage,
				},
				Annotations: map[string]string{},
			},
		}

		tenantHash, err := t.GetTenantHash(store.Labels)
		if err != nil {
			return err
		}

		storageHash, err := t.GetStorageHash(e.storage)
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
			e.service,
			e.storage)
	}

	return nil
}

func (t *Tenant) sortTenantsByStorageAndService() (map[string]*elem, error) {
	tenantList := &monitoringv1alpha1.TenantList{}
	if err := t.Client.List(t.Context, tenantList); err != nil {
		return nil, err
	}

	storageMap := make(map[string]*elem)
	for _, tenant := range tenantList.Items {

		if tenant.DeletionTimestamp != nil || !tenant.DeletionTimestamp.IsZero() {
			klog.V(3).Infof("ignore tenant %s is deleting", tenant.Name)
			continue
		}

		if tenant.Labels == nil {
			klog.V(3).Infof("ignore tenant %s without service and storage", tenant.Name)
			continue
		}

		serviceNamespacedName := tenant.Labels[constants.ServiceLabelKey]
		if serviceNamespacedName == "" {
			klog.V(3).Infof("ignore tenant %s without service", tenant.Name)
			continue
		}

		storageNamespacedName := t.GetStorage(tenant.Labels[constants.StorageLabelKey])
		if storageNamespacedName == constants.LocalStorage {
			klog.V(3).Infof("ignore tenant %s with local storage", tenant.Name)
			continue
		}

		storageMap[serviceNamespacedName+"|"+storageNamespacedName] = &elem{
			service: serviceNamespacedName,
			storage: storageNamespacedName,
		}
	}

	return storageMap, nil
}
