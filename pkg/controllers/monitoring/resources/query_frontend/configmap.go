package query_frontend

import (
	"time"

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

	type Config struct {
		MaxSize      string        `yaml:"max_size"`
		MaxSizeItems int           `yaml:"max_size_items"`
		Validity     time.Duration `yaml:"validity"`
	}
	type CacheConfig struct {
		Type   string `yaml:"type"`
		Config Config `yaml:"config"`
	}

	var cacheConfig = CacheConfig{
		Type: "IN-MEMORY",
		Config: Config{
			MaxSize:      "",
			MaxSizeItems: 0,
			Validity:     time.Second * 0,
		},
	}

	if q.queryFrontend.MaxSizeInMemoryCacheConfig != "" {
		cacheConfig.Config.MaxSize = q.queryFrontend.MaxSizeInMemoryCacheConfig
	}

	if q.queryFrontend.MaxSizeItemsInMemoryCacheConfig != 0 {
		cacheConfig.Config.MaxSizeItems = int(q.queryFrontend.MaxSizeItemsInMemoryCacheConfig)
	}

	if q.queryFrontend.ValidityInMemoryCacheConfig != 0 {
		cacheConfig.Config.Validity = time.Duration(q.queryFrontend.ValidityInMemoryCacheConfig) * time.Second
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
