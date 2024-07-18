package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/yaml"

	"github.com/WhizardTelemetry/whizard/pkg/client/k8s"
)

func newTestConfig() (*Config, error) {

	var conf = &Config{
		KubernetesOptions: &k8s.KubernetesOptions{
			KubeConfig: "/Users/frezes/.kube/config",
			Master:     "https://127.0.0.1:6443",
			QPS:        1e6,
			Burst:      1e6,
		},
	}
	return conf, nil
}

func saveTestConfig(t *testing.T, conf *Config) {

	content, err := yaml.Marshal(conf)
	if err != nil {
		t.Fatalf("error marshal config. %v", err)
	}
	err = os.WriteFile(fmt.Sprintf("%s.yaml", defaultConfigurationName), content, 0640)
	if err != nil {
		t.Fatalf("error write configuration file, %v", err)
	}
}

func cleanTestConfig(t *testing.T) {
	file := fmt.Sprintf("%s.yaml", defaultConfigurationName)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Log("file not exists, skipping")
		return
	}

	err := os.Remove(file)
	if err != nil {
		t.Fatalf("remove %s file failed", file)
	}

}

func TestGet(t *testing.T) {
	conf, err := newTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	saveTestConfig(t, conf)
	defer cleanTestConfig(t)

	conf2, err := TryLoadFromDisk()
	if err != nil {
		t.Fatal(err)
	}
	opt := cmp.Comparer(func(x, y resource.Quantity) bool {
		return x.Equal(y)
	})

	if diff := cmp.Diff(conf, conf2, opt); diff != "" {
		t.Fatal(diff)
	}
}
