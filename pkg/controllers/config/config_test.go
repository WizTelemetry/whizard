package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/client/k8s"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/options"
	"k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func newTestConfig() (*Config, error) {
	var replicas1 int32 = 1
	var replicas2 int32 = 2
	var replicas uint64 = 1
	var pflase bool = false
	var ptrue bool = true
	var stabilizationWindowSeconds int32 = 300
	var cpuAverageUtilization int32 = 80
	var memAverageUtilization int32 = 80
	var conf = &Config{
		KubernetesOptions: &k8s.KubernetesOptions{
			KubeConfig: "/Users/frezes/.kube/config",
			Master:     "https://127.0.0.1:6443",
			QPS:        1e6,
			Burst:      1e6,
		},
		MonitoringOptions: &options.Options{
			Compactor: &options.CompactorOptions{
				DefaultTenantsPerCompactor: 10,
				CommonOptions: options.CommonOptions{
					Replicas:  &replicas1,
					Image:     "kubesphere/thanos:v0.29.5",
					LogLevel:  "info",
					LogFormat: "logfmt",
					Flags:     []string{"--block-files-concurrency=20", "--compact.blocks-fetch-concurrency=5", "--web.disable"},
					Resources: corev1.ResourceRequirements{},
				},
				DataVolume: &v1alpha1.KubernetesVolume{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaim{
						Spec: corev1.PersistentVolumeClaimSpec{
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceStorage: resource.MustParse("10Gi"),
								},
							},
						},
					},
				},
			},
			Gateway: &options.GatewayOptions{
				CommonOptions: options.CommonOptions{
					Image:     "kubesphere/whizard-options-gateway:v0.5.0",
					Replicas:  &replicas2,
					LogLevel:  "info",
					LogFormat: "logfmt",
				},
			},
			Ingester: &options.IngesterOptions{
				LocalTsdbRetention:             "7d",
				DisableTSDBCleanup:             &pflase,
				DefaultTenantsPerIngester:      3,
				DefaultIngesterRetentionPeriod: time.Hour * 3,
				TSDBCleanupImage:               "bash:5.1.16",
				CommonOptions: options.CommonOptions{
					Image:     "kubesphere/thanos:v0.29.5",
					Replicas:  &replicas2,
					LogLevel:  "info",
					LogFormat: "logfmt",
					Resources: corev1.ResourceRequirements{},
				},
				DataVolume: &v1alpha1.KubernetesVolume{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaim{
						Spec: corev1.PersistentVolumeClaimSpec{
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceStorage: resource.MustParse("20Gi"),
								},
							},
						},
					},
				},
			},
			Query: &options.QueryOptions{
				CommonOptions: options.CommonOptions{
					Image:     "kubesphere/thanos:v0.28.5",
					Replicas:  &replicas2,
					LogLevel:  "info",
					LogFormat: "logfmt",
					Flags:     []string{"--query.max-concurrent=200"},
				},
				Envoy: &options.SidecarOptions{
					Image: "envoyproxy/envoy:corev1.20.2",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("500Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("500Mi"),
						},
					},
				},
			},
			QueryFrontend: &options.QueryFrontendOptions{
				CommonOptions: options.CommonOptions{
					Image:     "kubesphere/thanos:v0.29.5",
					Replicas:  &replicas2,
					LogLevel:  "info",
					LogFormat: "logfmt",
				},
			},
			Router: &options.RouterOptions{
				ReplicationFactor: &replicas,
				CommonOptions: options.CommonOptions{
					Image:     "kubesphere/thanos:v0.29.5",
					Replicas:  &replicas2,
					LogLevel:  "info",
					LogFormat: "logfmt",
				},
			},
			Ruler: &options.RulerOptions{
				Shards:             &replicas1,
				EvaluationInterval: "1m",
				RuleSelectors: []*metav1.LabelSelector{
					&metav1.LabelSelector{
						MatchLabels: map[string]string{"role": "alert-rules"},
					},
				},
				AlertmanagersURL: []string{"dnssrv+http://alertmanager-operated.kubesphere-monitoring-system.svc:9093"},
				CommonOptions: options.CommonOptions{
					Image:     "kubesphere/thanos:v0.29.5",
					Replicas:  &replicas1,
					LogLevel:  "info",
					LogFormat: "logfmt",
				},
				PrometheusConfigReloader: options.SidecarOptions{
					Image: "kubesphere/prometheus-config-reloader:v0.55.1",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("50Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("50Mi"),
						},
					},
				},
				RulerQueryProxy: options.SidecarOptions{
					Image: "kubesphere/whizard-monitoring-gateway:v0.5.0",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("50Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("50Mi"),
						},
					},
				},
				RulerWriteProxy: options.SidecarOptions{
					Image: "kubesphere/cortex-tenant:v1.7.2",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("50Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("200m"),
							corev1.ResourceMemory: resource.MustParse("400Mi"),
						},
					},
				},
			},
			Store: &options.StoreOptions{
				CommonOptions: options.CommonOptions{
					Replicas:  &replicas1,
					Image:     "kubesphere/thanos:v0.29.5",
					LogLevel:  "info",
					LogFormat: "logfmt",
					Flags:     []string{"--max-time=-36h", "--web.disable"},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("500Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("500Mi"),
						},
					},
				},
				AutoScaler: &v1alpha1.AutoScaler{
					MinReplicas: &replicas2,
					MaxReplicas: 20,
					Behavior: &v2beta2.HorizontalPodAutoscalerBehavior{
						ScaleUp: &v2beta2.HPAScalingRules{
							StabilizationWindowSeconds: &stabilizationWindowSeconds,
						},
					},
					Metrics: []v2beta2.MetricSpec{
						{
							Type: v2beta2.ResourceMetricSourceType,
							Resource: &v2beta2.ResourceMetricSource{
								Name: corev1.ResourceCPU,
								Target: v2beta2.MetricTarget{
									Type:               v2beta2.UtilizationMetricType,
									AverageUtilization: &cpuAverageUtilization,
								},
							},
						},
						{
							Type: v2beta2.ResourceMetricSourceType,
							Resource: &v2beta2.ResourceMetricSource{
								Name: corev1.ResourceMemory,
								Target: v2beta2.MetricTarget{
									Type:               v2beta2.UtilizationMetricType,
									AverageUtilization: &memAverageUtilization,
								},
							},
						},
					},
				},
			},
			Storage: &options.StorageOptions{
				BlockManager: &options.BlockManagerOptions{
					Enable:             &ptrue,
					ServiceAccountName: "whizard-controller-manager",
					BlockSyncInterval:  &metav1.Duration{Duration: time.Minute},
					CommonOptions: options.CommonOptions{
						Image:    "kubesphere/thanos:v0.29.5",
						Replicas: &replicas1,
					},
					GC: &options.BlockGCOptions{
						Enable: &ptrue,
						Image:  "kubesphere/whizard-monitoring-block-manager:v0.5.0",
					},
				},
			},
		},
	}
	return conf, nil
}

func saveTestConfig(t *testing.T, conf *Config) {

	content, err := yaml.Marshal(conf)
	if err != nil {
		t.Fatalf("error marshal config. %v", err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s.yaml", defaultConfigurationName), content, 0640)
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
