package thanosreceive

import (
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *ThanosReceive) Ingresses() []*netv1.Ingress {
	var ingresses []*netv1.Ingress

	if routingSpec := r.Instance.Spec.Router; routingSpec != nil {
		if spec := routingSpec.RemoteWriteIngress; spec != nil && spec.Host != "" {
			pathType := netv1.PathTypeImplementationSpecific
			ingr := &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: r.Instance.Namespace,
					Name:      r.getRemoteWriteIngressName(),
					Labels:    r.labels(),
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{{
						Host: spec.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{{
									PathType: &pathType,
									Path:     spec.Path,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: r.getReceiveOperatedServiceName(),
											Port: netv1.ServiceBackendPort{Number: remoteWritePort},
										},
									},
								}},
							},
						},
					}},
				},
			}
			if spec.SecretName != "" {
				ingr.Spec.TLS = append(ingr.Spec.TLS, netv1.IngressTLS{
					Hosts:      []string{spec.Host},
					SecretName: spec.SecretName,
				})
			}
			ingresses = append(ingresses, ingr)
		}
	}

	return ingresses
}
