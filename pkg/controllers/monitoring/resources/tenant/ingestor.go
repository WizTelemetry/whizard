package tenant

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/common/model"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/util"
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

			if t.Service != nil {
				// If ingester.status.Tenants contain the deleted tenant, remove that
				var ingesterList monitoringv1alpha1.IngesterList
				err := t.Client.List(t.Context, &ingesterList, &client.ListOptions{
					Namespace:     t.Service.Namespace,
					LabelSelector: labels.SelectorFromSet(util.ManagedLabelByService(t.Service)),
				})
				if err != nil {
					return err
				}
				for _, ingester := range ingesterList.Items {
					var tenantsStatus []v1alpha1.IngesterTenantStatus
					for _, tenant := range ingester.Status.Tenants {
						if tenant.Name != t.tenant.Name {
							tenantsStatus = append(tenantsStatus, tenant)
						}
					}
					if !reflect.DeepEqual(tenantsStatus, ingester.Status.Tenants) {
						ingester.Status.Tenants = tenantsStatus
						if err := t.Client.Status().Update(t.Context, &ingester); err != nil {
							return err
						}
					}
				}
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

	if v, ok := t.tenant.Labels[constants.ExclusiveLabelKey]; ok && v == "true" {
		ingester := t.createIngesterInstance(t.tenant.Name, t.tenant, true)

		t.tenant.Status.Ingester = &monitoringv1alpha1.ObjectReference{
			Namespace: ingester.Namespace,
			Name:      ingester.Name,
		}

		if err := util.CreateOrUpdate(t.Context, t.Client, ingester); err != nil {
			return err
		}
		return t.Client.Status().Update(t.Context, t.tenant)
	}

	var ingesterList monitoringv1alpha1.IngesterList
	ls := make(map[string]string, 2)
	ls[constants.ServiceLabelKey] = t.tenant.Labels[constants.ServiceLabelKey]
	ls[constants.StorageLabelKey] = t.GetStorage(t.tenant.Labels[constants.StorageLabelKey])
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
		name := t.createIngesterInstanceName(strconv.Itoa(i))
		if ingesterItem, ok := ingesterMapping[name]; ok {
			if len(ingesterItem.Spec.Tenants) < t.Service.Spec.IngesterTemplateSpec.DefaultTenantsPerIngester {
				ingester = ingesterItem
				addTenantToIngesterInstance(t.tenant, ingester)
				break
			}
		} else {
			ingester = t.createIngesterInstance(name, t.tenant, false)
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

	if v, ok := ingester.Labels[constants.StorageLabelKey]; !ok || v != t.GetStorage(t.tenant.Labels[constants.StorageLabelKey]) {
		klog.V(3).Infof("Tenant [%s] and ingester [%s]'s Storage mismatch, need to reset ingester", t.tenant.Name, ingester.Name)
		return true, nil
	}

	if v, ok := t.tenant.Labels[constants.ExclusiveLabelKey]; ok && v == "true" {
		if v, ok := ingester.Labels[constants.ExclusiveLabelKey]; !ok || v != "true" || ingester.Name != t.tenant.Name {
			klog.V(3).Infof("Tenant [%s] requires its exclusive ingester, the current shared ingester [%s] should be replaced by an exclusive one", t.tenant.Name, ingester.Name)
			return true, nil
		}
	} else {
		if v, ok := ingester.Labels[constants.ExclusiveLabelKey]; ok && v == "true" && ingester.Name == t.tenant.Name {
			klog.V(3).Infof("Tenant [%s] is not allowed to have its exclusive ingester, the current exclusive ingester [%s] will be replaced by a shared ingester", t.tenant.Name, ingester.Name)
			return true, nil
		}
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

				var retentionPeriod time.Duration
				// When ingester uses object storage, the ingester retention period uses the DefaultIngesterRetentionPeriod.
				// When it uses local storage, its retention period is the same as LocalTsdbRetention
				if v, ok := ingester.Labels[constants.StorageLabelKey]; ok && v != constants.LocalStorage {
					period, _ := model.ParseDuration(string(t.Service.Spec.IngesterTemplateSpec.DefaultIngesterRetentionPeriod))
					retentionPeriod = time.Duration(period)
				} else {
					if ingester.Spec.LocalTsdbRetention != "" {
						period, _ := model.ParseDuration(ingester.Spec.LocalTsdbRetention)
						retentionPeriod = time.Duration(period)
					} else if t.Service.Spec.IngesterTemplateSpec.DefaultIngesterRetentionPeriod != "" {
						period, _ := model.ParseDuration(ingester.Spec.LocalTsdbRetention)
						retentionPeriod = time.Duration(period)
					}
					if retentionPeriod <= 0 {
						period, _ := model.ParseDuration(string(t.Service.Spec.IngesterTemplateSpec.DefaultIngesterRetentionPeriod))
						retentionPeriod = time.Duration(period)
					}
				}

				annotation[constants.LabelNameIngesterState] = constants.IngesterStateDeleting
				annotation[constants.LabelNameIngesterDeletingTime] = strconv.Itoa(int(time.Now().Add(retentionPeriod).Unix()))
				ingester.Annotations = annotation
			}

			return util.CreateOrUpdate(t.Context, t.Client, ingester)
		}
	}
	return nil
}

func (t *Tenant) createIngesterInstanceName(suffix ...string) string {
	storage := t.GetStorage(t.tenant.Labels[constants.StorageLabelKey])

	serviceNamespacedName := strings.Split(t.tenant.Labels[constants.ServiceLabelKey], ".")
	storageNamespacedName := strings.Split(storage, ".")
	storageName := constants.LocalStorage
	if len(storageNamespacedName) >= 2 {
		storageName = storageNamespacedName[1]
	}

	name := fmt.Sprintf("%s-%s-auto", serviceNamespacedName[1], storageName)
	if len(suffix) > 0 {
		name += "-" + strings.Join(suffix, "-")
	}
	return name
}

func (t *Tenant) createIngesterInstance(name string, tenant *monitoringv1alpha1.Tenant, isExclusive bool) *monitoringv1alpha1.Ingester {
	klog.V(3).Infof("create new ingester %s for tenant %s", name, tenant.Name)
	storage := t.GetStorage(tenant.Labels[constants.StorageLabelKey])

	label := make(map[string]string, 3)
	label[constants.ServiceLabelKey] = tenant.Labels[constants.ServiceLabelKey]
	label[constants.StorageLabelKey] = storage
	if isExclusive {
		label[constants.ExclusiveLabelKey] = "true"
	}

	namespacedName := strings.Split(tenant.Labels[constants.ServiceLabelKey], ".")

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
