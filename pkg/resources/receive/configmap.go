package receive

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin-monitoring/pkg/resources"
)

func (r *receiveRouter) hashringsConfigMap() (runtime.Object, resources.Operation, error) {

	var cm = &corev1.ConfigMap{ObjectMeta: r.meta(r.name("hashrings-config"))}

	if r.del {
		return cm, resources.OperationDelete, nil
	}

	type HashringConfig struct {
		Hashring  string   `json:"hashring,omitempty"`
		Tenants   []string `json:"tenants,omitempty"`
		Endpoints []string `json:"endpoints"`
	}
	var hashrings []HashringConfig
	var softHashrings []HashringConfig
	for _, i := range r.receive.Ingestors {
		in := i
		ingestor := receiveIngestor{Receive: r.Receive, Ingestor: in}
		if len(in.Tenants) == 0 {
			softHashrings = append(softHashrings, HashringConfig{
				Hashring:  in.Name,
				Tenants:   in.Tenants,
				Endpoints: ingestor.grpcAddrs(),
			})
			continue
		}
		hashrings = append(hashrings, HashringConfig{
			Hashring:  in.Name,
			Tenants:   in.Tenants,
			Endpoints: ingestor.grpcAddrs(),
		})
	}
	hashrings = append(hashrings, softHashrings...) // put soft tenants at the end

	hashringBytes, err := json.MarshalIndent(hashrings, "", "\t")
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}
	cm.Data = map[string]string{
		hashringsFile: string(hashringBytes),
	}

	return cm, resources.OperationCreateOrUpdate, nil
}
