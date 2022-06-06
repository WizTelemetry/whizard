package query_frontend

import (
	"time"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (q *QueryFrontend) cacheConfigConfigMap() (runtime.Object, resources.Operation, error) {
	var cm = &corev1.ConfigMap{ObjectMeta: q.meta(q.name("cache-config"))}

	if q.queryFrontend == nil {
		return cm, resources.OperationDelete, nil
	}

	type CacheProviderConfig struct {
		Type   string      `yaml:"type"`
		Config interface{} `yaml:"config"`
	}
	var cacheConfig CacheProviderConfig

	var defaultINMEMORYCacheConfig = CacheProviderConfig{
		Type: string(v1alpha1.INMEMORY),
		Config: v1alpha1.InMemoryResponseCacheConfig{
			MaxSize:      "",
			MaxSizeItems: 0,
			Validity:     time.Duration(0),
		},
	}

	if q.queryFrontend.CacheConfig != nil {
		switch q.queryFrontend.CacheConfig.Type {
		case v1alpha1.INMEMORY:
			if q.queryFrontend.CacheConfig.InMemoryResponseCacheConfig == nil {
				cacheConfig = defaultINMEMORYCacheConfig
			} else {
				cacheConfig = CacheProviderConfig{
					Type:   string(v1alpha1.INMEMORY),
					Config: *q.queryFrontend.CacheConfig.InMemoryResponseCacheConfig,
				}
			}

		// todo: support other cache.
		// case v1alpha1.MEMCACHED:
		// case v1alpha1.REDIS:
		default:
			cacheConfig = defaultINMEMORYCacheConfig
		}

	} else {
		cacheConfig = defaultINMEMORYCacheConfig
	}

	cacheConfigBytes, err := yaml.Marshal(cacheConfig)
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}

	cm.Data = map[string]string{
		cacheConfigFile: string(cacheConfigBytes),
	}

	return cm, resources.OperationCreateOrUpdate, nil

}
