package tenant

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/util"
)

func (t *Tenant) ingester() error {

	// finalizers check,  when tenant cr is deleted, ObjectMeta.GetDeletionTimestamp() is not nil, remove finalizers field and call removeTenantFromIngesterbyName()
	if t.tenant.ObjectMeta.GetDeletionTimestamp().IsZero() {
		if !containsString(t.tenant.ObjectMeta.Finalizers, monitoringv1alpha1.FinalizerMonitoringIngester) {
			t.tenant.ObjectMeta.Finalizers = append(t.tenant.ObjectMeta.Finalizers, monitoringv1alpha1.FinalizerMonitoringIngester)
			return t.Client.Update(t.Context, t.tenant)
		}
	} else {
		if containsString(t.tenant.ObjectMeta.Finalizers, monitoringv1alpha1.FinalizerMonitoringIngester) {
			if t.tenant.Status.Ingester != nil {
				if err := t.removeTenantFromIngesterbyName(t.tenant.Status.Ingester.Namespace, t.tenant.Status.Ingester.Name); err != nil {
					return err
				}
				t.tenant.Status.Ingester = nil
			}

			t.tenant.ObjectMeta.Finalizers = removeString(t.tenant.Finalizers, monitoringv1alpha1.FinalizerMonitoringIngester)
			return t.Client.Update(t.Context, t.tenant)
		}
	}

	// Check if ingester needs to be created or reset
	ingester := &monitoringv1alpha1.Ingester{}
	if t.tenant.Status.Ingester != nil {
		err := t.Client.Get(t.Context, types.NamespacedName{
			Namespace: t.tenant.Status.Ingester.Namespace,
			Name:      t.tenant.Status.Ingester.Name,
		}, ingester)
		if err != nil {
			if apierrors.IsNotFound(err) {
				klog.V(3).Infof("Cannot find ingester [%s] for tenant [%s], create one", t.tenant.Status.Ingester.Name, t.tenant.Name)
			} else {
				return err
			}
		} else {
			var needResetIngester bool = false
			if ok := containsString(ingester.Spec.Tenants, t.tenant.Spec.Tenant); !ok {
				klog.V(3).Infof("Tenant [%s] and ingester [%s] mismatch, need to reset ingester", t.tenant.Name, ingester.Name)
				needResetIngester = true
			}

			if v, ok := ingester.Labels[monitoringv1alpha1.MonitoringPaodinService]; !ok || v != t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService] {
				klog.V(3).Infof("Tenant [%s] and ingester [%s]'s Service mismatch, need to reset ingester", t.tenant.Name, ingester.Name)
				needResetIngester = true
			}

			if v, ok := ingester.Labels[monitoringv1alpha1.MonitoringPaodinStorage]; !ok || v != t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinStorage] {
				klog.V(3).Infof("Tenant [%s] and ingester [%s]'s Storage mismatch, need to reset ingester", t.tenant.Name, ingester.Name)
				needResetIngester = true
			}

			if !needResetIngester {
				return nil
			} else {
				klog.V(3).Infof("Reset ingester [%s] for tenant [%s]", ingester.Name, t.tenant.Name)
				err := t.removeTenantFromIngesterbyName(ingester.Namespace, ingester.Name)
				if err != nil {
					return err
				}

				return t.Client.Status().Update(t.Context, t.tenant)
			}
		}
	}

	// when tenant.Labels don't contain Service, remove the bindings to ingester and ruler
	if v, ok := t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService]; !ok || v == "" {
		klog.V(3).Infof("Tenant [%s]'s Service is empty. ingester does not need to be created", t.tenant.Name)
		if t.tenant.Status.Ingester != nil {
			err := t.removeTenantFromIngesterbyName(t.tenant.Status.Ingester.Namespace, t.tenant.Status.Ingester.Name)
			if err != nil {
				return err
			}
			return t.Client.Status().Update(t.Context, t.tenant)
		}
		return nil
	}

	var ingesterList monitoringv1alpha1.IngesterList
	ls := make(map[string]string, 2)
	ls[monitoringv1alpha1.MonitoringPaodinService] = t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService]
	ls[monitoringv1alpha1.MonitoringPaodinStorage] = t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinStorage]
	serviceNamespacedName := strings.Split(t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService], ".")
	err := t.Client.List(t.Context, &ingesterList, &client.ListOptions{
		Namespace:     serviceNamespacedName[0],
		LabelSelector: labels.SelectorFromSet(ls),
	})
	if err != nil {
		return err
	}

	ingesterMapping := make(map[string]*monitoringv1alpha1.Ingester, len(ingesterList.Items))
	for _, ingesterItem := range ingesterList.Items {
		ingester := ingesterItem
		ingesterMapping[ingesterItem.Name] = &ingester
		klog.V(3).Infof("Ingester [%s] have Tenants: %v", ingesterItem.Name, ingesterItem.Spec.Tenants)
	}

	// Check duplicates
	for _, ingesterItem := range ingesterMapping {
		if containsString(ingesterItem.Spec.Tenants, t.tenant.Spec.Tenant) {
			klog.V(3).Infof("Ingester [%s] has tenant [%s] ,update status ", ingesterItem.Name, t.tenant.Name)

			t.tenant.Status.Ingester = &monitoringv1alpha1.ObjectReference{
				Namespace: ingesterItem.Namespace,
				Name:      ingesterItem.Name,
			}

			return t.Client.Status().Update(t.Context, t.tenant)
		}
	}

	// create or update ingester instance.
	// traverse ingesterMapping according to the index, if it is currently empty, create a new instance,
	// otherwise check len(ingesterItem.Spec.Tenants) < t.DefaultTenantsPerIngester，if so, select the instance.
	for i := 0; i < len(ingesterMapping)+1; i++ {
		name := createIngesterInstanceName(t.tenant, strconv.Itoa(i))
		if ingesterItem, ok := ingesterMapping[name]; ok {
			if len(ingesterItem.Spec.Tenants) < t.Options.DefaultTenantsPerIngester {
				ingester = ingesterItem
				addTenantToIngesterInstance(t.tenant, ingester)
				break
			}
		} else {
			ingester = createIngesterInstance(name, t.tenant)
			break
		}
	}

	t.tenant.Status.Ingester = &monitoringv1alpha1.ObjectReference{
		Namespace: ingester.Namespace,
		Name:      ingester.Name,
	}

	if err := util.CreateOrUpdate(t.Context, t.Client, ingester); err != nil {
		return err
	}
	return t.Client.Status().Update(t.Context, t.tenant)
}

func (t *Tenant) removeTenantFromIngesterbyName(namespace, name string) error {
	ingester := &monitoringv1alpha1.Ingester{}

	err := t.Client.Get(t.Context, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, ingester)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	} else {
		if ok := containsString(ingester.Spec.Tenants, t.tenant.Spec.Tenant); ok {
			klog.V(3).Infof("ingester %s update, remove tenant %s", ingester.Name, t.tenant.Name)
			ingester.Spec.Tenants = removeString(ingester.Spec.Tenants, t.tenant.Spec.Tenant)
			ingester.Labels[monitoringv1alpha1.MonitoringPaodinTenant] = strings.Join(ingester.Spec.Tenants, "_")

			if len(ingester.Spec.Tenants) == 0 {
				annotation := ingester.GetAnnotations()
				if annotation == nil {
					annotation = make(map[string]string)
				}
				annotation[resources.LabelNameIngesterState] = "deleting"
				annotation[resources.LabelNameIngesterDeletingTime] = strconv.Itoa(int(time.Now().Add(t.Options.DefaultIngesterRetentionPeriod).Unix()))
				ingester.Annotations = annotation
			}

			if t.tenant.Status.Ingester != nil {
				t.tenant.Status.Ingester = nil
			}

			return util.CreateOrUpdate(t.Context, t.Client, ingester)
		}
	}
	return nil
}

func (t *Tenant) deleteIngesterInstance(namespace, name string) error {
	ingester := &monitoringv1alpha1.Ingester{}
	err := t.Client.Get(t.Context, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, ingester)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	annotations := ingester.GetAnnotations()
	if annotations != nil {
		if v, ok := annotations[resources.LabelNameIngesterState]; ok && v == "deleting" {
			if v, ok := annotations[monitoringv1alpha1.MonitoringPaodinTenant]; !ok || len(v) == 0 {
				klog.V(3).Infof("Ingester %s will be deleted.")
				t.Client.Delete(t.Context, ingester)
			}
		}
	}
	return nil
}

// selectStoragebyMatchLabels randomly get Storage by select label
func (t *Tenant) selectStoragebyMatchLabels(matchLabels map[string]string) (*monitoringv1alpha1.Storage, error) {
	storageList := &monitoringv1alpha1.StorageList{}
	err := t.Client.List(t.Context, storageList, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(matchLabels),
	})
	if err != nil {
		return nil, err
	}
	if len(storageList.Items) == 0 {
		return nil, fmt.Errorf("can't select Storage by matchLabels [%v]", matchLabels)
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(storageList.Items))

	return &storageList.Items[index], err
}

// isStorageContainLabel  if storage contains matchLabels
func (t *Tenant) isStorageContainLabel(namespace, name string, matchLabels map[string]string) (bool, error) {
	storage := &monitoringv1alpha1.Storage{}
	err := t.Client.Get(t.Context, types.NamespacedName{Namespace: namespace, Name: name}, storage)
	if err != nil {
		return false, err
	}
	for key, value := range matchLabels {
		if v, ok := storage.Labels[key]; !ok || v != value {
			return false, nil
		}
	}
	return true, nil
}

func createIngesterInstanceName(tenant *monitoringv1alpha1.Tenant, suffix ...string) string {
	serviceNamespacedName := strings.Split(tenant.Labels[monitoringv1alpha1.MonitoringPaodinService], ".")
	storageNamespacedName := strings.Split(tenant.Labels[monitoringv1alpha1.MonitoringPaodinStorage], ".")

	name := fmt.Sprintf("%s-%s-auto", serviceNamespacedName[1], storageNamespacedName[1])
	if len(suffix) > 0 {
		name += "-" + strings.Join(suffix, "-")
	}
	return name
}

func createIngesterInstance(name string, tenant *monitoringv1alpha1.Tenant) *monitoringv1alpha1.Ingester {
	klog.V(3).Infof("create new ingester %s for tenant %s", name, tenant.Name)
	label := make(map[string]string, 2)
	label[monitoringv1alpha1.MonitoringPaodinService] = tenant.Labels[monitoringv1alpha1.MonitoringPaodinService]
	label[monitoringv1alpha1.MonitoringPaodinStorage] = tenant.Labels[monitoringv1alpha1.MonitoringPaodinStorage]
	label[monitoringv1alpha1.MonitoringPaodinTenant] = tenant.Name

	namespacedName := strings.Split(tenant.Labels[monitoringv1alpha1.MonitoringPaodinService], ".")
	// todo: ingester config
	return &monitoringv1alpha1.Ingester{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespacedName[0],
			Labels:    label,
		},
		Spec: monitoringv1alpha1.IngesterSpec{
			Tenants: []string{tenant.Spec.Tenant},
		},
	}
}

func addTenantToIngesterInstance(tenant *monitoringv1alpha1.Tenant, ingester *monitoringv1alpha1.Ingester) {
	klog.V(3).Infof("Ingester %s update, add tenant %s", ingester.Name, tenant.Name)

	ingester.Spec.Tenants = append(ingester.Spec.Tenants, tenant.Spec.Tenant)

	label := ingester.GetLabels()
	if v, ok := label[monitoringv1alpha1.MonitoringPaodinTenant]; !ok || len(v) == 0 {
		label[monitoringv1alpha1.MonitoringPaodinTenant] = tenant.Name
	} else {
		label[monitoringv1alpha1.MonitoringPaodinTenant] = label[monitoringv1alpha1.MonitoringPaodinTenant] + "." + tenant.Name
	}
	ingester.Labels = label

	annotation := ingester.GetAnnotations()
	if v, ok := annotation[resources.LabelNameIngesterState]; ok && v == "deleting" {
		annotation[resources.LabelNameIngesterState] = "running"
		annotation[resources.LabelNameIngesterDeletingTime] = ""
	}
	ingester.Annotations = annotation
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
