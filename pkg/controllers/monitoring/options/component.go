package options

import (
	"time"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/spf13/pflag"
	"k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CompactorOptions struct {
	CommonOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	DefaultTenantsPerCompactor int `json:"defaultTenantsPerCompactor,omitempty" yaml:"defaultTenantsPerCompactor,omitempty"`
	// DisableDownsampling specifies whether to disable downsampling
	DisableDownsampling *bool `json:"disableDownsampling,omitempty" yaml:"disableDownsampling,omitempty"`
	// Retention configs how long to retain samples
	Retention *v1alpha1.Retention `json:"retention,omitempty" yaml:"retention,omitempty"`
}

func NewCompactorOptions() *CompactorOptions {
	return &CompactorOptions{
		CommonOptions:              NewCommonOptions(),
		DefaultTenantsPerCompactor: DefaultTenantsPerCompactor,
	}
}

func (o *CompactorOptions) ApplyTo(options *CompactorOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)

	if o.DefaultTenantsPerCompactor != 0 {
		options.DefaultTenantsPerCompactor = o.DefaultTenantsPerCompactor
	}

	if o.DisableDownsampling != nil {
		options.DisableDownsampling = o.DisableDownsampling
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
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *CompactorOptions) AddFlags(fs *pflag.FlagSet, c *CompactorOptions) {
	o.CommonOptions.AddFlags(fs, &c.CommonOptions, "compactor")
	fs.IntVar(&c.DefaultTenantsPerCompactor, "default-tenants-per-compactor", DefaultTenantsPerCompactor, "Number of tenants processed per compactor")
}

type IngesterOptions struct {
	CommonOptions `json:",inline" yaml:",inline"`

	DefaultTenantsPerIngester int `json:"defaultTenantsPerIngester,omitempty" yaml:"defaultTenantsPerIngester,omitempty"`
	// DefaultIngesterRetentionPeriod ... Maybe it can be clearer?
	DefaultIngesterRetentionPeriod time.Duration `json:"defaultIngesterRetentionPeriod,omitempty" yaml:"defaultIngesterRetentionPeriod,omitempty"`

	// LocalTsdbRetention configs how long to retain raw samples on local storage.
	LocalTsdbRetention string `json:"localTsdbRetention,omitempty"`
}

func NewIngesterOptions() *IngesterOptions {
	return &IngesterOptions{
		CommonOptions: NewCommonOptions(),

		DefaultTenantsPerIngester:      DefaultTenantsPerIngester,
		DefaultIngesterRetentionPeriod: DefaultIngesterRetentionPeriod,
	}
}

func (o *IngesterOptions) ApplyTo(options *IngesterOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)
}

func (o *IngesterOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *IngesterOptions) AddFlags(fs *pflag.FlagSet, io *IngesterOptions) {
	o.CommonOptions.AddFlags(fs, &io.CommonOptions, "ingester")

	fs.IntVar(&io.DefaultTenantsPerIngester, "defaultTenantsPerIngester", DefaultTenantsPerIngester, "Whizard default tenant count per ingester.")
	fs.DurationVar(&io.DefaultIngesterRetentionPeriod, "defaultIngesterRetentionPeriod", DefaultIngesterRetentionPeriod, "Whizard default ingester retention period.")
}

type GatewayOptions struct {
	CommonOptions `json:",inline" yaml:",inline"`

	// Secret name for HTTP Server certificate (Kubernetes TLS secret type)
	ServerCertificate string `json:"serverCertificate,omitempty"`
	// Secret name for HTTP Client CA certificate (Kubernetes TLS secret type)
	ClientCACertificate string `json:"clientCaCertificate,omitempty"`

	NodePort int32 `json:"nodePort,omitempty"`
}

func NewGatewayOptions() *GatewayOptions {
	o := &GatewayOptions{
		CommonOptions: NewCommonOptions(),
	}

	o.Image = DefaultGatewayImage
	return o
}

func (o *GatewayOptions) ApplyTo(options *GatewayOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)

	if options.NodePort == 0 {
		options.NodePort = o.NodePort
	}
}

func (o *GatewayOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *GatewayOptions) AddFlags(fs *pflag.FlagSet, g *GatewayOptions) {
	o.CommonOptions.AddFlags(fs, &g.CommonOptions, "gateway")
}

type QueryOptions struct {
	CommonOptions `json:",inline" yaml:",inline"`

	Envoy *ContainerOptions `json:"envoy,omitempty" yaml:"envoy,omitempty"`
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		CommonOptions: NewCommonOptions(),
		Envoy: &ContainerOptions{
			Image: DefaultEnvoyImage,
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
	}
}

func (o *QueryOptions) ApplyTo(options *QueryOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)
	o.Envoy.ApplyTo(options.Envoy)

}

func (o *QueryOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *QueryOptions) AddFlags(fs *pflag.FlagSet, qo *QueryOptions) {
	o.CommonOptions.AddFlags(fs, &qo.CommonOptions, "query")
	o.Envoy.AddFlags(fs, qo.Envoy, "query.envoy")
}

type QueryFrontendOptions struct {
	CommonOptions `json:",inline" yaml:",inline"`

	CacheConfig *v1alpha1.ResponseCacheProviderConfig `json:"cacheConfig,omitempty" yaml:"cacheConfig,omitempty"`
}

func NewQueryFrontendOptions() *QueryFrontendOptions {
	return &QueryFrontendOptions{
		CommonOptions: NewCommonOptions(),
	}
}

func (o *QueryFrontendOptions) ApplyTo(options *QueryFrontendOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)

}

func (o *QueryFrontendOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *QueryFrontendOptions) AddFlags(fs *pflag.FlagSet, qfo *QueryFrontendOptions) {
	o.CommonOptions.AddFlags(fs, &qfo.CommonOptions, "query-frontend")
}

type RouterOptions struct {
	CommonOptions `json:",inline" yaml:",inline"`

	// How many times to replicate incoming write requests
	ReplicationFactor *uint64 `json:"replicationFactor,omitempty"`
}

func NewRouterOptions() *RouterOptions {
	var factor uint64 = DefaultRouterReplicationFactor
	return &RouterOptions{
		CommonOptions: NewCommonOptions(),

		ReplicationFactor: &factor,
	}
}

func (o *RouterOptions) ApplyTo(options *RouterOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)

}

func (o *RouterOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *RouterOptions) AddFlags(fs *pflag.FlagSet, ro *RouterOptions) {
	var factor uint64
	o.CommonOptions.AddFlags(fs, &ro.CommonOptions, "router")
	fs.Uint64Var(&factor, "router.replicationFactor", DefaultRouterReplicationFactor, "Whizard router replication factor.")

	ro.ReplicationFactor = &factor
}

type RulerOptions struct {
	CommonOptions `json:",inline" yaml:",inline"`

	PrometheusConfigReloader *ContainerOptions `json:"prometheusConfigReloader,omitempty" yaml:"prometheusConfigReloader,omitempty"`
	RulerQueryProxy          *ContainerOptions `json:"rulerQueryProxy" yaml:"rulerQueryProxy,omitempty"`

	// Number of shards to take the hash of fully qualified name of the rule group in order to split rules.
	// Each shard of rules will be bound to one separate statefulset.
	Shards *int32 `json:"shards,omitempty"`
	// A label selector to select which PrometheusRules to mount for alerting and
	// recording.
	RuleSelector *metav1.LabelSelector `json:"ruleSelector,omitempty"`
	// Namespaces to be selected for PrometheusRules discovery. If unspecified, only
	// the same namespace as the Ruler object is in is used.
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`

	// Labels configure the external label pairs to Ruler. A default replica label
	// `ruler_replica` will be always added  as a label with the value of the pod's name and it will be dropped in the alerts.
	Labels map[string]string `json:"labels,omitempty"`
	// AlertDropLabels configure the label names which should be dropped in Ruler alerts.
	// The replica label `ruler_replica` will always be dropped in alerts.
	AlertDropLabels []string `json:"alertDropLabels,omitempty"`
	// Define configuration for connecting to alertmanager. Maps to the `alertmanagers.config` arg.
	AlertManagersConfig *corev1.SecretKeySelector `json:"alertmanagersConfig,omitempty"`

	// Interval between consecutive evaluations.
	EvaluationInterval time.Duration `json:"evaluationInterval,omitempty"`
}

func NewRulerOptions() *RulerOptions {
	var shards int32 = DefaultRulerShards
	return &RulerOptions{
		CommonOptions:      NewCommonOptions(),
		Shards:             &shards,
		EvaluationInterval: DefaultRulerEvaluationInterval,

		PrometheusConfigReloader: &ContainerOptions{
			Image: DefaultPrometheusConfigReloaderImage,
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
		RulerQueryProxy: &ContainerOptions{
			Image: DefaultGatewayImage,
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
	}
}

func (o *RulerOptions) ApplyTo(options *RulerOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)
	o.PrometheusConfigReloader.ApplyTo(options.PrometheusConfigReloader)

}

func (o *RulerOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *RulerOptions) AddFlags(fs *pflag.FlagSet, ro *RulerOptions) {
	o.CommonOptions.AddFlags(fs, &ro.CommonOptions, "ruler")
	o.PrometheusConfigReloader.AddFlags(fs, ro.PrometheusConfigReloader, "ruler.prometheus-config-reloader")
	o.RulerQueryProxy.AddFlags(fs, ro.RulerQueryProxy, "ruler.query-proxy")
}

type StoreOptions struct {
	CommonOptions              `json:",inline" yaml:",inline"`
	*v1alpha1.IndexCacheConfig `json:"indexCacheConfig,omitempty" yaml:"indexCacheConfig,omitempty"`
	*v1alpha1.AutoScaler       `json:"scaler,omitempty" yaml:"scaler,omitempty"`
}

func NewStoreOptions() *StoreOptions {
	var replicas int32 = DefaultStoreMinReplicas
	var min int32 = DefaultStoreMinReplicas
	var stabilizationWindowSeconds int32 = 300
	var cpuAverageUtilization int32 = 80
	var memAverageUtilization int32 = 80

	return &StoreOptions{
		CommonOptions: CommonOptions{
			Image:    DefaultWhizardImage,
			Replicas: &replicas,
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
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *StoreOptions) AddFlags(fs *pflag.FlagSet, s *StoreOptions) {
	o.CommonOptions.AddFlags(fs, &s.CommonOptions, "store")
}

type StorageOptions struct {
	Bucket *BucketOptions `json:"bucket,omitempty"`
}

type BucketOptions struct {
	Enable             *bool `json:"enable,omitempty"`
	CommonOptions      `json:",inline"`
	ServiceAccountName string           `json:"serviceAccountName,omitempty"`
	Refresh            *metav1.Duration `json:"refresh,omitempty"`
	GC                 *BucketGCOptions `json:"gc,omitempty"`
}

type BucketGCOptions struct {
	Enable          *bool                       `json:"enable,omitempty"`
	Resources       corev1.ResourceRequirements `json:"resources,omitempty"`
	Image           string                      `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	Interval        *metav1.Duration            `json:"interval,omitempty"`
	CleanupTimeout  *metav1.Duration            `json:"cleanupTimeout,omitempty"`
}

func NewStorageOptions() *StorageOptions {
	enable := true
	refresh := metav1.Duration{time.Minute}
	return &StorageOptions{
		Bucket: &BucketOptions{
			Enable:        &enable,
			CommonOptions: NewCommonOptions(),
			Refresh:       &refresh,
			GC: &BucketGCOptions{
				Enable: &enable,
				Image:  DefaultBucketImage,
			},
			ServiceAccountName: DefaultServiceAccount,
		},
	}
}

func (o *StorageOptions) ApplyTo(options *StorageOptions) {
	if o.Bucket != nil {
		if options.Bucket == nil {
			options.Bucket = o.Bucket
		} else {
			o.Bucket.CommonOptions.ApplyTo(&options.Bucket.CommonOptions)

			if options.Bucket.Enable == nil {
				options.Bucket.Enable = o.Bucket.Enable
			}

			if options.Bucket.Refresh == nil || options.Bucket.Refresh.Duration == 0 {
				options.Bucket.Refresh = o.Bucket.Refresh
			}
			if options.Bucket.ServiceAccountName == "" {
				options.Bucket.ServiceAccountName = o.Bucket.ServiceAccountName
			}

			if o.Bucket.GC != nil {
				if options.Bucket.GC == nil {
					options.Bucket.GC = o.Bucket.GC
				} else {
					if options.Bucket.GC.Image == "" {
						options.Bucket.GC.Image = o.Bucket.GC.Image
					}

					if options.Bucket.GC.ImagePullPolicy == "" {
						options.Bucket.GC.ImagePullPolicy = o.Bucket.GC.ImagePullPolicy
					}

					if options.Bucket.GC.Interval == nil ||
						options.Bucket.GC.Interval.Duration == 0 {
						options.Bucket.GC.Interval = o.Bucket.GC.Interval
					}

					if options.Bucket.GC.CleanupTimeout == nil ||
						options.Bucket.GC.CleanupTimeout.Duration == 0 {
						options.Bucket.GC.CleanupTimeout = o.Bucket.GC.CleanupTimeout
					}

					if options.Bucket.GC.Enable == nil {
						options.Bucket.GC.Enable = o.Bucket.GC.Enable
					}
				}
			}
		}
	}
}

func (o *StorageOptions) Validate() []error {
	var errs []error
	if o.Bucket != nil {
		errs = append(errs, o.Bucket.CommonOptions.Validate()...)
	}

	return errs
}

func (o *StorageOptions) AddFlags(fs *pflag.FlagSet, s *StorageOptions) {
	if o.Bucket != nil && s.Bucket != nil {
		o.Bucket.CommonOptions.AddFlags(fs, &s.Bucket.CommonOptions, "storage")
	}
}
