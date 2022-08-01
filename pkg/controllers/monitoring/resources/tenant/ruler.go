package tenant

import (
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/util"
)

func (t *Tenant) ruler() error {

	ruler := &monitoringv1alpha1.Ruler{}
	if t.tenant.Status.Ruler != nil {
		err := t.Client.Get(t.Context, types.NamespacedName{
			Namespace: t.tenant.Status.Ruler.Namespace,
			Name:      t.tenant.Status.Ruler.Name,
		}, ruler)
		if err != nil {
			if apierrors.IsNotFound(err) {
				klog.V(3).Infof("Cannot find ruler [%s] for tenant [%s], create one", t.tenant.Status.Ruler, t.tenant.Name)
			} else {
				return err
			}
		} else {
			var needResetRuler bool = false
			// todo: more ruler check
			if v, ok := ruler.Labels[monitoringv1alpha1.MonitoringPaodinService]; !ok || v != t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService] {
				klog.V(3).Infof("Tenant [%s] and ruler [%s]'s Service mismatch, need to reset ingester", t.tenant.Name, ruler.Name)
				needResetRuler = true
			}

			if !needResetRuler {
				return nil
			} else {
				klog.V(3).Infof("Ruler [%s] is already assigned to tenant [%s],  reset ruler for this tenant", ruler.Name, t.tenant.Name)
			}
		}
	}

	// when tenant.Labels don't contain Service, remove the bindings to ingester and ruler
	if v, ok := t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService]; !ok || v == "" {
		klog.V(3).Infof("Tenant [%s]'s Service is empty. ruler does not need to be created", t.tenant.Name)
		if t.tenant.Status.Ruler != nil && ruler != nil {
			if err := t.Client.Delete(t.Context, ruler); err != nil {
				return err
			}
			t.tenant.Status.Ruler = nil
			return t.Client.Status().Update(t.Context, t.tenant)
		}
		return nil
	}
	ruler = t.createOrUpdateRulerinstance()
	t.tenant.Status.Ruler = &monitoringv1alpha1.ObjectReference{
		Namespace: ruler.Namespace,
		Name:      ruler.Name,
	}
	err := util.CreateOrUpdate(t.Context, t.Client, ruler)
	if err != nil {
		return err
	}
	return t.Client.Status().Update(t.Context, t.tenant)
}

func (t *Tenant) createOrUpdateRulerinstance() *monitoringv1alpha1.Ruler {

	label := make(map[string]string, 2)
	label[monitoringv1alpha1.MonitoringPaodinTenant] = t.tenant.Name
	label[monitoringv1alpha1.MonitoringPaodinService] = t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService]

	serviceNamespacedName := strings.Split(t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService], ".")

	// todo: thanosruler config
	return &monitoringv1alpha1.Ruler{ObjectMeta: metav1.ObjectMeta{
		Name:      t.tenant.Name,
		Namespace: serviceNamespacedName[0],
		Labels:    label,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: t.tenant.APIVersion,
				Kind:       t.tenant.Kind,
				Name:       t.tenant.Name,
				UID:        t.tenant.UID,
				Controller: pointer.BoolPtr(true),
			},
		},
	},
		Spec: monitoringv1alpha1.RulerSpec{
			Tenant: t.tenant.Name,
		},
	}
}