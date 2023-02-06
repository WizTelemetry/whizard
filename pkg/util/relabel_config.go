package util

import (
	"fmt"
	"strings"
)

type relabelConfig struct {
	Action        string   `yaml:"action"`
	SourceLablels []string `yaml:"source_labels"`
	Regex         string   `yaml:"regex"`
}

func CreateKeepTenantsRelabelConfig(tenantLabelName string, tenants []string) (string, error) {

	regex := ""
	for _, tenant := range tenants {
		regex = fmt.Sprintf("%s|^%s$", regex, tenant)
	}

	return YamlMarshal([]relabelConfig{
		{Action: "keep",
			SourceLablels: []string{tenantLabelName},
			Regex:         strings.TrimPrefix(regex, "|"),
		},
	})
}
