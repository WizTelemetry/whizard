package options

import (
	"time"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/util"
	"github.com/spf13/pflag"
	"k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	DefaultWhizardImage            = "thanosio/thanos:v0.26.0"
	DefaultEnvoyImage              = "envoyproxy/envoy:v1.20.2"
	DefaultGatewayImage            = "kubesphere/whizard-monitoring-gateway:latest"
	DefaultService                 = "kubesphere-monitoring-system.central"
	DefaultTenantsPerIngester      = 3
	DefaultIngesterRetentionPeriod = time.Hour * 3
	DefaultTenantsPerCompactor     = 10

	DefaultStoreMinReplicas = 2
	DefaultStoreMaxReplicas = 20
)

var PrometheusConfigReloaderDefaultConfig = PrometheusConfigReloaderConfig{
	Image:         "quay.io/prometheus-operator/prometheus-config-reloader:v0.55.1",
	CPURequest:    "100m",
	MemoryRequest: "50Mi",
	CPULimit:      "100m",
	MemoryLimit:   "50Mi",
}

var RulerQueryProxyDefaultConfig = RulerQueryProxyConfig{
	Image:         DefaultGatewayImage,
	CPURequest:    "100m",
	MemoryRequest: "50Mi",
	CPULimit:      "100m",
	MemoryLimit:   "50Mi",
}

type CommonOptions struct {
	Image           string                      `json:"image,omitempty" yaml:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy           `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	Affinity        *corev1.Affinity            `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	NodeSelector    map[string]string           `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	Tolerations     []corev1.Toleration         `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	Resources       corev1.ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`
	Replicas        *int32                      `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	LogLevel        string                      `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	LogFormat       string                      `json:"logFormat,omitempty" yaml:"logFormat,omitempty"`
	Flags           map[string]string           `json:"flags,omitempty" yaml:"flags,omitempty"`
	DataVolume      *v1alpha1.KubernetesVolume  `json:"dataVolume,omitempty" yaml:"dataVolume,omitempty"`
}

func (o *CommonOptions) ApplyTo(options *CommonOptions) {
	if o.Image != "" {
		options.Image = o.Image
	}

	if o.ImagePullPolicy != "" {
		options.ImagePullPolicy = o.ImagePullPolicy
	}

	if o.Affinity != nil {
		if options.Affinity == nil {
			options.Affinity = o.Affinity
		}

		util.Override(options.Affinity, o.Affinity)
	}

	if o.Tolerations != nil {
		options.Tolerations = o.Tolerations
	}

	if o.NodeSelector != nil {
		options.NodeSelector = o.NodeSelector
	}

	if o.Resources.Limits != nil {
		if options.Resources.Limits == nil {
			options.Resources.Limits = o.Resources.Limits
		}
		for k, v := range o.Resources.Limits {
			options.Resources.Limits[k] = v
		}
	}

	if o.Resources.Requests != nil {
		if options.Resources.Requests == nil {
			options.Resources.Requests = o.Resources.Requests
		}
		for k, v := range o.Resources.Requests {
			options.Resources.Requests[k] = v
		}
	}

	if o.Replicas != nil {
		options.Replicas = o.Replicas
	}

	if o.LogLevel != "" {
		options.LogLevel = o.LogLevel
	}

	if o.LogFormat != "" {
		options.LogFormat = o.LogFormat
	}

	if o.Flags != nil {
		if options.Flags == nil {
			options.Flags = make(map[string]string)
		}

		for k, v := range o.Flags {
			options.Flags[k] = v
		}
	}

	if o.DataVolume != nil {
		if options.DataVolume == nil {
			options.DataVolume = o.DataVolume
		}

		if o.DataVolume.PersistentVolumeClaim != nil {
			options.DataVolume.PersistentVolumeClaim = o.DataVolume.PersistentVolumeClaim
		}

		if o.DataVolume.EmptyDir != nil {
			options.DataVolume.EmptyDir = o.DataVolume.EmptyDir
		}
	}
}

func (o *CommonOptions) AddFlags(fs *pflag.FlagSet, c *CommonOptions, prefix string) {
	fs.StringVar(&c.Image, prefix+".image", c.Image, "Image with tag/version")
	fs.StringVar(&c.LogLevel, prefix+".log.level", c.LogLevel, "Log filtering level")
	fs.StringVar(&c.LogFormat, prefix+".log.format", c.LogLevel, "Log format to use. Possible options: logfmt or json")
}

type PrometheusConfigReloaderConfig struct {
	Image         string `json:"image,omitempty" yaml:"image,omitempty"`
	CPURequest    string `json:"cpuRequest,omitempty" yaml:"cpuRequest,omitempty"`
	MemoryRequest string `json:"memoryRequest,omitempty" yaml:"memoryRequest,omitempty"`
	CPULimit      string `json:"cpuLimit,omitempty" yaml:"cpuRequest,omitempty"`
	MemoryLimit   string `json:"memoryLimit,omitempty" yaml:"memoryLimit,omitempty"`
}

func (o *PrometheusConfigReloaderConfig) Validate() []error {
	return nil
}

func (o *PrometheusConfigReloaderConfig) ApplyTo(options *PrometheusConfigReloaderConfig) {
	if o.Image != "" {
		options.Image = o.Image
	}
	if o.CPURequest != "" {
		options.CPURequest = o.CPURequest
	}
	if o.MemoryRequest != "" {
		options.MemoryRequest = o.MemoryRequest
	}
	if o.CPULimit != "" {
		options.CPULimit = o.CPULimit
	}
	if o.MemoryLimit != "" {
		options.MemoryLimit = o.MemoryLimit
	}
}

type RulerQueryProxyConfig struct {
	Image         string `json:"image,omitempty" yaml:"image,omitempty"`
	CPURequest    string `json:"cpuRequest,omitempty" yaml:"cpuRequest,omitempty"`
	MemoryRequest string `json:"memoryRequest,omitempty" yaml:"memoryRequest,omitempty"`
	CPULimit      string `json:"cpuLimit,omitempty" yaml:"cpuRequest,omitempty"`
	MemoryLimit   string `json:"memoryLimit,omitempty" yaml:"memoryLimit,omitempty"`
}

func (o *RulerQueryProxyConfig) Validate() []error {
	var errs []error
	return errs
}

func (o *RulerQueryProxyConfig) ApplyTo(options *PrometheusConfigReloaderConfig) {
	if o.Image != "" {
		options.Image = o.Image
	}
	if o.CPURequest != "" {
		options.CPURequest = o.CPURequest
	}
	if o.MemoryRequest != "" {
		options.MemoryRequest = o.MemoryRequest
	}
	if o.CPULimit != "" {
		options.CPULimit = o.CPULimit
	}
	if o.MemoryLimit != "" {
		options.MemoryLimit = o.MemoryLimit
	}
}

type CompactorOptions struct {
	CommonOptions              `json:"_inline,omitempty" yaml:"_inline,omitempty"`
	DefaultTenantsPerCompactor int `json:"defaultTenantsPerCompactor,omitempty" yaml:"defaultTenantsPerCompactor,omitempty"`
	// DownsamplingDisable specifies whether to disable downsampling
	DownsamplingDisable *bool `json:"downsamplingDisable,omitempty" yaml:"downsamplingDisable,omitempty"`
	// Retention configs how long to retain samples
	Retention *v1alpha1.Retention `json:"retention,omitempty" yaml:"retention,omitempty"`
}

func (o *CompactorOptions) ApplyTo(options *CompactorOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)

	if o.DefaultTenantsPerCompactor != 0 {
		options.DefaultTenantsPerCompactor = o.DefaultTenantsPerCompactor
	}

	if o.DownsamplingDisable != nil {
		options.DownsamplingDisable = o.DownsamplingDisable
	}

	if o.Retention != nil {
		if options.Retention == nil {
			options.Retention = o.Retention
		} else {
			if o.Retention.Retention1h != "" {
				options.Retention.Retention1h = o.Retention.Retention1h
			}

			if o.Retention.RetentionRaw != "" {
				options.Retention.RetentionRaw = o.Retention.RetentionRaw
			}

			if o.Retention.Retention5m != "" {
				options.Retention.Retention5m = o.Retention.Retention5m
			}
		}
	}
}

func (o *CompactorOptions) Validate() []error {
	return nil
}

func (o *CompactorOptions) AddFlags(fs *pflag.FlagSet, c *CompactorOptions) {
	o.CommonOptions.AddFlags(fs, &c.CommonOptions, "compactor")
	fs.IntVar(&c.DefaultTenantsPerCompactor, "default-tenants-per-compactor", c.DefaultTenantsPerCompactor, "Number of tenants processed per compactor")
}

type StoreOptions struct {
	CommonOptions              `json:"_inline,omitempty" yaml:"_inline,omitempty"`
	*v1alpha1.IndexCacheConfig `json:"indexCacheConfig,omitempty" yaml:"indexCacheConfig,omitempty"`
	*v1alpha1.AutoScaler       `json:"scaler,omitempty" yaml:"scaler,omitempty"`
}

func defaultStoreOptions() StoreOptions {
	var min int32 = DefaultStoreMinReplicas
	var stabilizationWindowSeconds int32 = 300
	var cpuAverageUtilization int32 = 80
	var memAverageUtilization int32 = 80

	return StoreOptions{
		CommonOptions: CommonOptions{
			Image: DefaultWhizardImage,
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
			MinReplicas: &min,
			MaxReplicas: DefaultStoreMaxReplicas,
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
	}
}

func (o *StoreOptions) ApplyTo(options *StoreOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)

	if o.IndexCacheConfig != nil {
		if options.IndexCacheConfig == nil {
			options.IndexCacheConfig = o.IndexCacheConfig
		}

		if o.IndexCacheConfig.InMemoryIndexCacheConfig == nil {
			if options.IndexCacheConfig.InMemoryIndexCacheConfig == nil {
				options.IndexCacheConfig.InMemoryIndexCacheConfig = o.InMemoryIndexCacheConfig
			}

			if o.MaxSize != "" {
				options.MaxSize = o.MaxSize
			}
		}
	}

	if o.AutoScaler != nil {
		if options.AutoScaler == nil {
			options.AutoScaler = o.AutoScaler
		}

		if o.AutoScaler.MinReplicas != nil && *o.AutoScaler.MinReplicas > 0 {
			options.MinReplicas = o.MinReplicas
		}

		if o.MaxReplicas > 0 {
			options.MaxReplicas = o.MaxReplicas
		}

		if o.Behavior != nil {
			options.Behavior = o.Behavior
		}
	}
}

func (o *StoreOptions) Validate() []error {
	return nil
}

func (o *StoreOptions) AddFlags(_ *pflag.FlagSet, _ *StoreOptions) {
}

type Options struct {
	WhizardImage string `json:"whizardImage,omitempty" yaml:"whizardImage,omitempty"`
	EnvoyImage   string `json:"envoyImage,omitempty" yaml:"envoyImage,omitempty"`
	GatewayImage string `json:"gatewayImage,omitempty" yaml:"gGatewayImage,omitempty"`

	EnableKubeSphereAdapter        bool          `json:"enableKubeSphereAdapter,omitempty" yaml:"enableKubeSphereAdapter,omitempty"`
	KubeSphereAdapterService       string        `json:"kubeSphereAdapterService,omitempty" yaml:"kubeSphereAdapterService,omitempty"`
	DefaultTenantsPerIngester      int           `json:"defaultTenantsPerIngester,omitempty" yaml:"defaultTenantsPerIngester,omitempty"`
	DefaultIngesterRetentionPeriod time.Duration `json:"defaultIngesterRetentionPeriod,omitempty" yaml:"defaultIngesterRetentionPeriod,omitempty"`

	PrometheusConfigReloader PrometheusConfigReloaderConfig `json:"prometheusConfigReloader,omitempty" yaml:"prometheusConfigReloader,omitempty"`
	RulerQueryProxy          RulerQueryProxyConfig          `json:"rulerQueryProxy,omitempty" yaml:"rulerQueryProxy,omitempty"`

	Compactor CompactorOptions `json:"compactor,omitempty" yaml:"compactor,omitempty"`
	Store     StoreOptions     `json:"store,omitempty" yaml:"store,omitempty"`
}

func NewOptions() *Options {
	return &Options{
		WhizardImage:                   DefaultWhizardImage,
		EnvoyImage:                     DefaultEnvoyImage,
		GatewayImage:                   DefaultGatewayImage,
		DefaultTenantsPerIngester:      DefaultTenantsPerIngester,
		DefaultIngesterRetentionPeriod: DefaultIngesterRetentionPeriod,
		EnableKubeSphereAdapter:        true,
		KubeSphereAdapterService:       DefaultService,
		PrometheusConfigReloader:       PrometheusConfigReloaderDefaultConfig,
		RulerQueryProxy:                RulerQueryProxyDefaultConfig,
		Compactor: CompactorOptions{
			DefaultTenantsPerCompactor: DefaultTenantsPerCompactor,
		},
		Store: defaultStoreOptions(),
	}
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.PrometheusConfigReloader.Validate()...)
	errs = append(errs, o.Store.Validate()...)
	errs = append(errs, o.Compactor.Validate()...)
	return errs
}

func (o *Options) ApplyTo(options *Options) {
	if o.WhizardImage != "" {
		options.WhizardImage = o.WhizardImage
	}
	if o.EnvoyImage != "" {
		options.EnvoyImage = o.EnvoyImage
	}
	if o.GatewayImage != "" {
		options.GatewayImage = o.GatewayImage
	}
	if o.DefaultTenantsPerIngester != 0 {
		options.DefaultTenantsPerIngester = o.DefaultTenantsPerIngester
	}
	if o.DefaultIngesterRetentionPeriod != 0 {
		options.DefaultIngesterRetentionPeriod = o.DefaultIngesterRetentionPeriod
	}
	if o.KubeSphereAdapterService != "" {
		options.KubeSphereAdapterService = o.KubeSphereAdapterService
	}
	options.EnableKubeSphereAdapter = o.EnableKubeSphereAdapter

	o.PrometheusConfigReloader.ApplyTo(&options.PrometheusConfigReloader)
	o.Store.ApplyTo(&options.Store)
	o.Compactor.ApplyTo(&options.Compactor)
}

func (o *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	fs.StringVar(&c.WhizardImage, "whizard-image", DefaultWhizardImage, "Whizard image with tag/version")
	fs.StringVar(&c.EnvoyImage, "envoy-image", DefaultEnvoyImage, "Envoy image with tag/version")
	fs.StringVar(&c.GatewayImage, "gateway-image", DefaultGatewayImage, "Whizard monitoring gateway image with tag/version")
	fs.IntVar(&c.DefaultTenantsPerIngester, "defaultTenantsPerIngester", DefaultTenantsPerIngester, "Whizard default tenant count per ingester. (default 3)")
	fs.DurationVar(&c.DefaultIngesterRetentionPeriod, "defaultIngesterRetentionPeriod", DefaultIngesterRetentionPeriod, "Whizard default ingester retention period. (default 2h)")
	fs.BoolVar(&c.EnableKubeSphereAdapter, "enableKubeSphereAdapter", true, "Enable KubeSphere adapter. (default true)")
	fs.StringVar(&c.KubeSphereAdapterService, "kubeSphereAdapterService", DefaultService, "Default service for tenants generated by kubesphere adapter, format is namespace.name")

	fs.StringVar(&c.PrometheusConfigReloader.Image, "prometheus-config-reloader-image",
		PrometheusConfigReloaderDefaultConfig.Image, "Prometheus Config Reloader image with tag/version")
	fs.StringVar(&c.PrometheusConfigReloader.CPURequest, "prometheus-config-reloader-cpu-request",
		PrometheusConfigReloaderDefaultConfig.CPURequest, "Prometheus Config Reloader CPU request. Value \"0\" disables it and causes no request to be configured.")
	fs.StringVar(&c.PrometheusConfigReloader.CPULimit, "prometheus-config-reloader-cpu-limit",
		PrometheusConfigReloaderDefaultConfig.CPULimit, "Prometheus Config Reloader CPU limit. Value \"0\" disables it and causes no limit to be configured.")
	fs.StringVar(&c.PrometheusConfigReloader.MemoryRequest, "prometheus-config-reloader-memory-request",
		PrometheusConfigReloaderDefaultConfig.MemoryRequest, "Prometheus Config Reloader Memory request. Value \"0\" disables it and causes no request to be configured.")
	fs.StringVar(&c.PrometheusConfigReloader.MemoryLimit, "prometheus-config-reloader-memory-limit",
		PrometheusConfigReloaderDefaultConfig.MemoryLimit, "Prometheus Config Reloader Memory limit. Value \"0\" disables it and causes no limit to be configured.")

	fs.StringVar(&c.RulerQueryProxy.Image, "ruler-query-proxy-image",
		RulerQueryProxyDefaultConfig.Image, "Ruler Query Proxy image with tag/version")
	fs.StringVar(&c.RulerQueryProxy.CPURequest, "ruler-query-proxy-cpu-request",
		RulerQueryProxyDefaultConfig.CPURequest, "Ruler Query Proxy CPU request. Value \"0\" disables it and causes no request to be configured.")
	fs.StringVar(&c.RulerQueryProxy.CPULimit, "ruler-query-proxy-cpu-limit",
		RulerQueryProxyDefaultConfig.CPULimit, "Ruler Query Proxy CPU limit. Value \"0\" disables it and causes no limit to be configured.")
	fs.StringVar(&c.RulerQueryProxy.MemoryRequest, "ruler-query-proxy-memory-request",
		RulerQueryProxyDefaultConfig.MemoryRequest, "Ruler Query Proxy Memory request. Value \"0\" disables it and causes no request to be configured.")
	fs.StringVar(&c.RulerQueryProxy.MemoryLimit, "ruler-query-proxy-memory-limit",
		RulerQueryProxyDefaultConfig.MemoryLimit, "Ruler Query Proxy Memory limit. Value \"0\" disables it and causes no limit to be configured.")

	o.Store.AddFlags(fs, &c.Store)
	o.Compactor.AddFlags(fs, &o.Compactor)
}
