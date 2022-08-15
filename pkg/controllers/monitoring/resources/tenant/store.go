package tenant

import (
	"context"
	"fmt"
	"strings"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
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
	for _, store := range storeList.Items {
		if store.DeletionTimestamp != nil || !store.DeletionTimestamp.IsZero() {
			continue
		}

		if store.Labels == nil {
			continue
		}

		serviceNamespacedName := store.Labels[monitoringv1alpha1.MonitoringPaodinService]
		if serviceNamespacedName == "" {
			klog.V(3).Infof("Store [%s.%s] Service mismatch", store.Namespace, store.Name)
			continue
		}

		storageNamespacedName := store.Labels[monitoringv1alpha1.MonitoringPaodinStorage]
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

	unionMap(currentStores, expectStores)

	for _, v := range currentStores {
		store := v.(*monitoringv1alpha1.Store)
		if err := t.Client.Delete(t.Context, store); err != nil {
			return err
		}

		klog.V(3).Infof("Delete Store[%s.%s]", store.Namespace, store.Name)
	}

	for _, v := range expectStores {
		m := v.(map[string]string)
		serviceNamespacedName := strings.Split(m[monitoringv1alpha1.MonitoringPaodinService], ".")
		storageNamespacedName := strings.Split(m[monitoringv1alpha1.MonitoringPaodinStorage], ".")
		store := &monitoringv1alpha1.Store{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: fmt.Sprintf("store-%s-%s-", serviceNamespacedName[1], storageNamespacedName[1]),
				Namespace:    serviceNamespacedName[0],
				Labels:       m,
			},
		}

		if err := t.Client.Create(t.Context, store); err != nil {
			return err
		}

		klog.V(3).Infof("Create store[%s.%s] for Service[%s] Storage[%s]",
			store.Namespace, store.Name,
			m[monitoringv1alpha1.MonitoringPaodinService],
			m[monitoringv1alpha1.MonitoringPaodinStorage])
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

		if tenant.Labels[monitoringv1alpha1.MonitoringPaodinStorage] == monitoringv1alpha1.MonitoringLocalStorage {
			klog.V(3).Infof("ignore tenant %s with local storage", tenant.Name)
			continue
		}

		serviceNamespacedName := tenant.Labels[monitoringv1alpha1.MonitoringPaodinService]
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
			monitoringv1alpha1.MonitoringPaodinService: serviceNamespacedName,
			monitoringv1alpha1.MonitoringPaodinStorage: storageNamespacedName,
		}
	}

	return storageMap, nil
}

func unionMap(m1, m2 map[string]interface{}) {
	for k := range m1 {
		if _, ok := m2[k]; ok {
			delete(m1, k)
			delete(m2, k)
		}
	}
}
