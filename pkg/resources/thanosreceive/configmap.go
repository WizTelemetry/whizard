package thanosreceive

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/paodin-monitoring/api/v1alpha1"
)

func (r *ThanosReceive) configMaps() ([]*corev1.ConfigMap, error) {
	var cms []*corev1.ConfigMap
	if m := r.GetMode(); m == RouterOnly || m == RouterIngestor {
		cm, err := r.hashringsConfigMap()
		if err != nil {
			return nil, err
		}
		cms = append(cms, cm)
	}

	return cms, nil
}

func (r *ThanosReceive) hashringsConfigMap() (*corev1.ConfigMap, error) {

	var hashrings []HashringConfig
	if routingSpec := r.Instance.Spec.Router; routingSpec != nil {

		f := func(h *v1alpha1.RouterHashringConfig) error {
			if h == nil {
				return nil
			}
			var endpoints []string
			if len(h.Endpoints) > 0 {
				endpoints = h.Endpoints
			} else {
				var err error
				endpoints, err = r.selectEndpoints(h.EndpointsNamespaceSelector, h.EndpointsSelector)
				if err != nil {
					return err
				}
			}
			hashrings = append(hashrings, HashringConfig{
				Hashring:  h.Name,
				Tenants:   h.Tenants,
				Endpoints: endpoints,
			})
			return nil
		}

		for _, h := range routingSpec.HardTenantHashrings {
			if len(h.Tenants) == 0 {
				return nil, fmt.Errorf("HardTenantHashring can not be with empty tenants")
			}
			if err := f(h); err != nil {
				return nil, err
			}
		}

		if routingSpec.SoftTenantHashring != nil {
			if len(routingSpec.SoftTenantHashring.Tenants) > 0 {
				return nil, fmt.Errorf("SoftTenantHashring must be with empty tenants")
			}
			if err := f(routingSpec.SoftTenantHashring); err != nil {
				return nil, err
			}
		}
	}
	hashringBytes, err := json.MarshalIndent(hashrings, "", "\t")
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: r.Instance.Namespace,
			Name:      r.getHashringsConfigMapName(),
			Labels:    r.labels(),
		},
		Data: map[string]string{
			hashringsFileName: string(hashringBytes),
		},
	}, nil
}

func (r *ThanosReceive) selectEndpoints(nsSelector, epsSelector *metav1.LabelSelector) ([]string, error) {
	var (
		namespaces []string
		endpoints  []string
	)

	if nsSelector == nil {
		namespaces = append(namespaces, r.Instance.Namespace)
	} else {
		var namespaceList corev1.NamespaceList
		if nsSelector, err := metav1.LabelSelectorAsSelector(nsSelector); err != nil {
			return nil, err
		} else if err = r.Client.List(r.Context, &namespaceList, &client.ListOptions{LabelSelector: nsSelector}); err != nil {
			return nil, err
		}
		for _, ns := range namespaceList.Items {
			namespaces = append(namespaces, ns.Name)
		}
	}

	if epsSelector == nil {
		return endpoints, nil
	}

	for _, ns := range namespaces {
		var endpointsList corev1.EndpointsList
		if selector, err := metav1.LabelSelectorAsSelector(epsSelector); err != nil {
			return nil, err
		} else if err = r.Client.List(r.Context, &endpointsList, &client.ListOptions{Namespace: ns, LabelSelector: selector}); err != nil {
			return nil, err
		}
		for _, eps := range endpointsList.Items {
			for _, subset := range eps.Subsets {
				var port *int32
				for _, p := range subset.Ports {
					if p.Name == "grpc" {
						port = &p.Port
						break
					}
				}
				if port == nil {
					break
				}
				for _, addr := range subset.Addresses {
					endpoints = append(endpoints, fmt.Sprintf("%s:%d", addr.IP, *port))
				}
			}
		}
	}

	return endpoints, nil
}
