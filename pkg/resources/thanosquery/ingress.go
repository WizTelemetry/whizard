package thanosquery

import (
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (q *ThanosQuery) ingresses() []*netv1.Ingress {
	var ingresses []*netv1.Ingress

	if ingr := q.httpIngress(); ingr != nil {
		ingresses = append(ingresses, ingr)
	}

	if ingr := q.grpcIngress(); ingr != nil {
		ingresses = append(ingresses, ingr)
	}

	return ingresses
}

func (q *ThanosQuery) httpIngress() *netv1.Ingress {
	spec := q.Instance.Spec.HttpIngress
	if spec == nil || spec.Host == "" {
		return nil
	}
	pathType := netv1.PathTypeImplementationSpecific
	ingr := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: q.Instance.Namespace,
			Name:      q.getHttpIngressName(),
			Labels:    q.labels(),
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
									Name: q.getServiceName(),
									Port: netv1.ServiceBackendPort{Number: httpPort},
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
	return ingr
}

func (q *ThanosQuery) grpcIngress() *netv1.Ingress {
	spec := q.Instance.Spec.GrpcIngress
	if spec == nil || spec.Host == "" {
		return nil
	}
	pathType := netv1.PathTypeImplementationSpecific
	ingr := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: q.Instance.Namespace,
			Name:      q.getGrpcIngressName(),
			Labels:    q.labels(),
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/backend-protocol":   "GRPC",
				"nginx.ingress.kubernetes.io/grpc-backend":       "true",
				"nginx.ingress.kubernetes.io/protocol":           "h2c",
				"nginx.ingress.kubernetes.io/force-ssl-redirect": "true",
			},
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
									Name: q.getServiceName(),
									Port: netv1.ServiceBackendPort{Number: grpcPort},
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
	return ingr
}
