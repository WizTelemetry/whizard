package tenant

import (
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

func (t *Tenant) ruler() (runtime.Object, resources.Operation, error) {

	ruler := &monitoringv1alpha1.ThanosRuler{}
	if t.tenant.Status.ThanosResource != nil && t.tenant.Status.ThanosResource.ThanosRuler != nil {
		err := t.Client.Get(t.Context, types.NamespacedName{
			Namespace: t.tenant.Status.ThanosResource.ThanosRuler.Namespace,
			Name:      t.tenant.Status.ThanosResource.ThanosRuler.Name,
		}, ruler)
		if err != nil {
			if apierrors.IsNotFound(err) {
				klog.V(3).Infof("Tenant %s not found mapping ruler %s, need to create.", t.tenant.Name, t.tenant.Status.ThanosResource.ThanosRuler)
			} else {
				klog.Errorf("Client Get ThanosRuler error: %v", err)
				return nil, "", err
			}
		} else {
			var needResetRuler bool = false
			// todo ruler check

			if !needResetRuler {
				return nil, "", nil
			} else {
				klog.V(3).Infof("Tenant %s mapping ruler %s need reset", t.tenant.Name, ruler.Name)
				t.tenant.Status.ThanosResource.ThanosRuler = nil
				if err := t.Client.Status().Update(t.Context, t.tenant); err != nil {
					klog.Error(err)
				}
				return ruler, resources.OperationCreateOrUpdate, nil
			}
		}
	}

	klog.V(3).Infof("Tenant %s ruler need to create or reset.", t.tenant.Name)
	ruler = t.createOrUpdateRulerinstance()
	if t.tenant.Status.ThanosResource == nil {
		t.tenant.Status.ThanosResource = &monitoringv1alpha1.ThanosResource{}
	}
	t.tenant.Status.ThanosResource.ThanosRuler = &monitoringv1alpha1.ObjectReference{
		Namespace: ruler.Namespace,
		Name:      ruler.Name,
	}
	if err := t.Client.Status().Update(t.Context, t.tenant); err != nil {
		klog.Error(err)
	}

	return ruler, resources.OperationCreateOrUpdate, nil
}

func (t *Tenant) createOrUpdateRulerinstance() *monitoringv1alpha1.ThanosRuler {

	label := make(map[string]string, 2)
	label[monitoringv1alpha1.MonitoringPaodinTenant] = t.tenant.Name
	label[monitoringv1alpha1.MonitoringPaodinService] = t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService]

	serviceNamespacedName := strings.Split(t.tenant.Labels[monitoringv1alpha1.MonitoringPaodinService], ".")

	// todo: thanosruler config
	return &monitoringv1alpha1.ThanosRuler{ObjectMeta: metav1.ObjectMeta{
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
		Spec: monitoringv1alpha1.ThanosRulerSpec{
			Tenant: t.tenant.Name,
		},
	}
}
