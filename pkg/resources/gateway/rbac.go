package gateway

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/resources"
)

func (g *Gateway) serviceAccount() (runtime.Object, resources.Operation, error) {
	var sa = &corev1.ServiceAccount{ObjectMeta: g.meta(g.name())}
	if g.gateway == nil {
		return sa, resources.OperationDelete, nil
	}
	return sa, resources.OperationCreateOrUpdate, nil
}

func (g *Gateway) role() (runtime.Object, resources.Operation, error) {
	var r = &rbacv1.Role{ObjectMeta: g.meta(g.name())}
	if g.gateway == nil {
		return r, resources.OperationDelete, nil
	}
	r.Rules = append(r.Rules, rbacv1.PolicyRule{
		Verbs:     []string{"get", "list", "watch"},
		APIGroups: []string{v1alpha1.SchemeGroupVersion.Group},
		Resources: []string{"agents"},
	})
	return r, resources.OperationCreateOrUpdate, nil
}

func (g *Gateway) roleBinding() (runtime.Object, resources.Operation, error) {
	var rb = &rbacv1.RoleBinding{ObjectMeta: g.meta(g.name())}
	if g.gateway == nil {
		return rb, resources.OperationDelete, nil
	}
	rb.RoleRef = rbacv1.RoleRef{
		APIGroup: rbacv1.SchemeGroupVersion.Group,
		Kind:     "Role",
		Name:     g.name(),
	}
	rb.Subjects = append(rb.Subjects, rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      g.name(),
		Namespace: rb.Namespace,
	})
	return rb, resources.OperationCreateOrUpdate, nil
}
