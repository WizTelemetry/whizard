package monitoringgateway

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

type RemoteWriteConfig struct {
	Name          string            `yaml:"name,omitempty"`
	URL           *config.URL       `yaml:"url"`
	Headers       map[string]string `yaml:"headers,omitempty"`
	RemoteTimeout model.Duration    `yaml:"remote_timeout,omitempty"`
	TLSConfig     config.TLSConfig  `yaml:"tls_config,omitempty"`
}

// LoadRemoteWritesConfig loads remotewrites config, and prefers file to content
func LoadRemoteWritesConfig(file, content string) ([]RemoteWriteConfig, error) {
	var buff []byte
	if file != "" {
		c, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		buff = c
	} else {
		buff = []byte(content)
	}
	if len(buff) == 0 {
		return nil, nil
	}
	var rws []RemoteWriteConfig
	if err := yaml.UnmarshalStrict(buff, &rws); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return rws, nil
}
