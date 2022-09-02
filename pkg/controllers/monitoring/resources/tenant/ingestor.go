package tenant

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/util"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (t *Tenant) ingester() error {

	// finalizers check,  when tenant cr is deleted, ObjectMeta.GetDeletionTimestamp() is not nil, remove finalizers field and call removeTenantFromIngesterbyName()
	if t.tenant.ObjectMeta.GetDeletionTimestamp().IsZero() {
		if !containsString(t.tenant.ObjectMeta.Finalizers, constants.FinalizerIngester) {
			t.tenant.ObjectMeta.Finalizers = append(t.tenant.ObjectMeta.Finalizers, constants.FinalizerIngester)
			return t.Client.Update(t.Context, t.tenant)
		}
	} else {
		if containsString(t.tenant.ObjectMeta.Finalizers, constants.FinalizerIngester) {
			if t.tenant.Status.Ingester != nil {
				if err := t.removeTenantFromIngesterbyName(t.tenant.Status.Ingester.Namespace, t.tenant.Status.Ingester.Name); err != nil {
					return err
				}
				t.tenant.Status.Ingester = nil
			}

			t.tenant.ObjectMeta.Finalizers = removeString(t.tenant.Finalizers, constants.FinalizerIngester)
			return t.Client.Update(t.Context, t.tenant)
		}
	}

	// Check if ingester needs to be created or reset
	if need, err := t.needResetIngester(); err != nil {
		return err
	} else if need {
		klog.V(3).Infof("Reset ingester [%s] for tenant [%s]", t.tenant.Status.Ingester.Name, t.tenant.Name)
		err := t.removeTenantFromIngesterbyName(t.tenant.Status.Ingester.Namespace, t.tenant.Status.Ingester.Name)
		if err != nil {
			return err
		}

		return t.Client.Status().Update(t.Context, t.tenant)
	}

	// when tenant.Labels don't contain Service, remove the bindings to ingester and ruler
	if v, ok := t.tenant.Labels[constants.ServiceLabelKey]; !ok || v == "" {
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
	ls[constants.ServiceLabelKey] = t.tenant.Labels[constants.ServiceLabelKey]
	ls[constants.StorageLabelKey] = t.tenant.Labels[constants.StorageLabelKey]
	serviceNamespacedName := strings.Split(t.tenant.Labels[constants.ServiceLabelKey], ".")
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
	ingester := &monitoringv1alpha1.Ingester{}
	for i := 0; i < len(ingesterMapping)+1; i++ {
		name := createIngesterInstanceName(t.tenant, strconv.Itoa(i))
		if ingesterItem, ok := ingesterMapping[name]; ok {
			if len(ingesterItem.Spec.Tenants) < t.Options.Ingester.DefaultTenantsPerIngester {
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

func (t *Tenant) needResetIngester() (bool, error) {
	if t.tenant.Status.Ingester == nil {
		return false, nil
	}

	ingester := &monitoringv1alpha1.Ingester{}
	err := t.Client.Get(t.Context, types.NamespacedName{
		Namespace: t.tenant.Status.Ingester.Namespace,
		Name:      t.tenant.Status.Ingester.Name,
	}, ingester)
	if err != nil && !apierrors.IsNotFound(err) {
		return false, err
	}

	if err != nil && apierrors.IsNotFound(err) {
		klog.V(3).Infof("Cannot find ingester [%s] for tenant [%s], need to reset ingester", t.tenant.Status.Ingester.Name, t.tenant.Name)
		return true, nil
	}

	if ok := containsString(ingester.Spec.Tenants, t.tenant.Spec.Tenant); !ok {
		klog.V(3).Infof("Tenant [%s] and ingester [%s] mismatch, need to reset ingester", t.tenant.Name, ingester.Name)
		return true, nil
	}

	if v, ok := ingester.Labels[constants.ServiceLabelKey]; !ok || v != t.tenant.Labels[constants.ServiceLabelKey] {
		klog.V(3).Infof("Tenant [%s] and ingester [%s]'s Service mismatch, need to reset ingester", t.tenant.Name, ingester.Name)
		return true, nil
	}

	if v, ok := ingester.Labels[constants.StorageLabelKey]; !ok || v != t.tenant.Labels[constants.StorageLabelKey] {
		klog.V(3).Infof("Tenant [%s] and ingester [%s]'s Storage mismatch, need to reset ingester", t.tenant.Name, ingester.Name)
		return true, nil
	}

	return false, nil
}

func (t *Tenant) removeTenantFromIngesterbyName(namespace, name string) error {
	if t.tenant.Status.Ingester != nil {
		t.tenant.Status.Ingester = nil
	}

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

			if len(ingester.Spec.Tenants) == 0 {
				annotation := ingester.GetAnnotations()
				if annotation == nil {
					annotation = make(map[string]string)
				}
				annotation[constants.LabelNameIngesterState] = constants.IngesterStateDeleting
				annotation[constants.LabelNameIngesterDeletingTime] = strconv.Itoa(int(time.Now().Add(t.Options.Ingester.DefaultIngesterRetentionPeriod).Unix()))
				ingester.Annotations = annotation
			}

			return util.CreateOrUpdate(t.Context, t.Client, ingester)
		}
	}
	return nil
}

func createIngesterInstanceName(tenant *monitoringv1alpha1.Tenant, suffix ...string) string {
	serviceNamespacedName := strings.Split(tenant.Labels[constants.ServiceLabelKey], ".")
	storageNamespacedName := strings.Split(tenant.Labels[constants.StorageLabelKey], ".")

	name := fmt.Sprintf("%s-%s-auto", serviceNamespacedName[1], storageNamespacedName[1])
	if len(suffix) > 0 {
		name += "-" + strings.Join(suffix, "-")
	}
	return name
}

func createIngesterInstance(name string, tenant *monitoringv1alpha1.Tenant) *monitoringv1alpha1.Ingester {
	klog.V(3).Infof("create new ingester %s for tenant %s", name, tenant.Name)
	label := make(map[string]string, 2)
	label[constants.ServiceLabelKey] = tenant.Labels[constants.ServiceLabelKey]
	label[constants.StorageLabelKey] = tenant.Labels[constants.StorageLabelKey]

	namespacedName := strings.Split(tenant.Labels[constants.ServiceLabelKey], ".")
	storage := &monitoringv1alpha1.ObjectReference{}
	if tenant.Labels[constants.StorageLabelKey] != constants.LocalStorage {
		storageNamespacedName := strings.Split(tenant.Labels[constants.StorageLabelKey], ".")
		storage.Namespace = storageNamespacedName[0]
		storage.Name = storageNamespacedName[1]
	}

	return &monitoringv1alpha1.Ingester{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespacedName[0],
			Labels:    label,
		},
		Spec: monitoringv1alpha1.IngesterSpec{
			Tenants: []string{tenant.Spec.Tenant},
			Storage: storage,
		},
	}
}

func addTenantToIngesterInstance(tenant *monitoringv1alpha1.Tenant, ingester *monitoringv1alpha1.Ingester) {
	klog.V(3).Infof("Ingester %s update, add tenant %s", ingester.Name, tenant.Name)

	ingester.Spec.Tenants = append(ingester.Spec.Tenants, tenant.Spec.Tenant)

	annotation := ingester.GetAnnotations()
	if v, ok := annotation[constants.LabelNameIngesterState]; ok && v == constants.IngesterStateDeleting {
		annotation[constants.LabelNameIngesterState] = constants.IngesterStateRunning
		annotation[constants.LabelNameIngesterDeletingTime] = ""
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
