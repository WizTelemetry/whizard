package options

import (
	"time"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/util"
	"github.com/spf13/pflag"
	"k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CompactorOptions struct {
	CommonOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	DataVolume *v1alpha1.KubernetesVolume `json:"dataVolume,omitempty" yaml:"dataVolume,omitempty"`

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

func (o *CompactorOptions) AddFlags(fs *pflag.FlagSet, c *CompactorOptions) {
	o.CommonOptions.AddFlags(fs, &c.CommonOptions, "compactor")
	fs.IntVar(&c.DefaultTenantsPerCompactor, "default-tenants-per-compactor", c.DefaultTenantsPerCompactor, "Number of tenants processed per compactor")
}

func (o *CompactorOptions) ApplyTo(options *CompactorOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)

	if o.DefaultTenantsPerCompactor != 0 {
		options.DefaultTenantsPerCompactor = o.DefaultTenantsPerCompactor
	}

	if o.DataVolume != nil {
		if options.DataVolume == nil {
			options.DataVolume = o.DataVolume
		} else {
			if o.DataVolume.PersistentVolumeClaim != nil {
				options.DataVolume.PersistentVolumeClaim = o.DataVolume.PersistentVolumeClaim
			}

			if o.DataVolume.EmptyDir != nil {
				options.DataVolume.EmptyDir = o.DataVolume.EmptyDir
			}
		}
	}

	if o.DisableDownsampling != nil {
		options.DisableDownsampling = o.DisableDownsampling
	}

	if o.Retention != nil {
		if options.Retention == nil {
			options.Retention = o.Retention
		} else {
			util.Override(options.Retention, o.Retention)
		}
	}
}

// Override the Options overrides the spec field when it is empty
func (o *CompactorOptions) Override(spec *v1alpha1.CompactorSpec) {
	o.CommonOptions.Override(&spec.CommonSpec)

	if spec.DataVolume == nil {
		spec.DataVolume = o.DataVolume
	}
	if spec.DisableDownsampling == nil {
		spec.DisableDownsampling = o.DisableDownsampling
	}
}

func (o *CompactorOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

type IngesterOptions struct {
	CommonOptions `json:",inline" yaml:",inline"  mapstructure:",squash"`

	DataVolume *v1alpha1.KubernetesVolume `json:"dataVolume,omitempty" yaml:"dataVolume,omitempty"`

	DefaultTenantsPerIngester int `json:"defaultTenantsPerIngester,omitempty" yaml:"defaultTenantsPerIngester,omitempty"`

	// DefaultIngesterRetentionPeriod Whizard default ingester retention period when it has no tenant.
	DefaultIngesterRetentionPeriod time.Duration `json:"defaultIngesterRetentionPeriod,omitempty" yaml:"defaultIngesterRetentionPeriod,omitempty"`

	// LocalTsdbRetention configs how long to retain raw samples on local storage.
	LocalTsdbRetention string `json:"localTsdbRetention,omitempty"`

	// Disable the TSDB cleanup of ingester.
	// The cleanup will delete the blocks that belong to deleted tenants
	// in the data directory of ingester TSDB.
	DisableTSDBCleanup *bool  `json:"disableTSDBCleanup,omitempty"`
	TSDBCleanupImage   string `json:"tsdbCleanupImage,omitempty"`
}

func NewIngesterOptions() *IngesterOptions {
	return &IngesterOptions{
		CommonOptions: NewCommonOptions(),

		DefaultTenantsPerIngester:      DefaultTenantsPerIngester,
		DefaultIngesterRetentionPeriod: DefaultIngesterRetentionPeriod,
		TSDBCleanupImage:               DefaultTSDBCleanupImage,
	}
}

func (o *IngesterOptions) AddFlags(fs *pflag.FlagSet, io *IngesterOptions) {
	o.CommonOptions.AddFlags(fs, &io.CommonOptions, "ingester")

	fs.IntVar(&io.DefaultTenantsPerIngester, "defaultTenantsPerIngester", io.DefaultTenantsPerIngester, "Whizard default tenant count per ingester.")
	fs.DurationVar(&io.DefaultIngesterRetentionPeriod, "defaultIngesterRetentionPeriod", io.DefaultIngesterRetentionPeriod, "Whizard default ingester retention period  when it has no tenant.")
}

func (o *IngesterOptions) ApplyTo(options *IngesterOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)

	if o.DataVolume != nil {
		if options.DataVolume == nil {
			options.DataVolume = o.DataVolume
		} else {
			if o.DataVolume.PersistentVolumeClaim != nil {
				options.DataVolume.PersistentVolumeClaim = o.DataVolume.PersistentVolumeClaim
			}

			if o.DataVolume.EmptyDir != nil {
				options.DataVolume.EmptyDir = o.DataVolume.EmptyDir
			}
		}
	}

	if o.DefaultTenantsPerIngester != 0 {
		options.DefaultTenantsPerIngester = o.DefaultTenantsPerIngester
	}
	if o.DefaultIngesterRetentionPeriod != 0 {
		options.DefaultIngesterRetentionPeriod = o.DefaultIngesterRetentionPeriod
	}
	if o.LocalTsdbRetention != "" {
		options.LocalTsdbRetention = o.LocalTsdbRetention
	}
	if o.DisableTSDBCleanup != nil {
		options.DisableTSDBCleanup = o.DisableTSDBCleanup
	}
	if o.TSDBCleanupImage != "" {
		options.TSDBCleanupImage = o.TSDBCleanupImage
	}
}

// Override the Options overrides the spec field when it is empty
func (o *IngesterOptions) Override(spec *v1alpha1.IngesterSpec) {
	o.CommonOptions.Override(&spec.CommonSpec)

	if spec.DataVolume == nil {
		spec.DataVolume = o.DataVolume
	}
	if spec.LocalTsdbRetention == "" {
		spec.LocalTsdbRetention = o.LocalTsdbRetention
	}

}

func (o *IngesterOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

type GatewayOptions struct {
	CommonOptions `json:",inline" yaml:",inline"  mapstructure:",squash"`

	// Secret name for HTTP Server certificate (Kubernetes TLS secret type)
	ServerCertificate string `json:"serverCertificate,omitempty"`
	// Secret name for HTTP Client CA certificate (Kubernetes TLS secret type)
	ClientCACertificate string `json:"clientCaCertificate,omitempty"`

	NodePort int32 `json:"nodePort,omitempty"`
}

func NewGatewayOptions() *GatewayOptions {
	var replicas int32 = 1

	return &GatewayOptions{
		CommonOptions: CommonOptions{
			Image:    DefaultGatewayImage,
			Replicas: &replicas,
		},
	}
}

func (o *GatewayOptions) ApplyTo(options *GatewayOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)

	if o.NodePort == 0 {
		options.NodePort = o.NodePort
	}
	if o.ClientCACertificate != "" {
		options.ClientCACertificate = o.ClientCACertificate
	}
	if o.ServerCertificate != "" {
		options.ServerCertificate = o.ServerCertificate
	}
}

func (o *GatewayOptions) Override(spec *v1alpha1.GatewaySpec) {
	o.CommonOptions.Override(&spec.CommonSpec)
	if spec.NodePort == 0 {
		spec.NodePort = o.NodePort
	}
	if spec.ServerCertificate != "" {
		spec.ServerCertificate = o.ServerCertificate
	}
	if spec.ClientCACertificate != "" {
		spec.ClientCACertificate = o.ClientCACertificate
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
	CommonOptions `json:",inline" yaml:",inline"  mapstructure:",squash"`

	Envoy *SidecarOptions `json:"envoy,omitempty" yaml:"envoy,omitempty"`

	// Additional StoreApi servers from which Query component queries from
	Stores []v1alpha1.QueryStores `json:"stores,omitempty"`
	// Selector labels that will be exposed in info endpoint.
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`
	// Labels to treat as a replica indicator along which data is deduplicated.
	ReplicaLabelNames []string `json:"replicaLabelNames,omitempty"`
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		CommonOptions: NewCommonOptions(),
		Envoy: &SidecarOptions{
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

	if o.Stores != nil {
		options.Stores = o.Stores
	}
	if o.SelectorLabels != nil {
		options.SelectorLabels = o.SelectorLabels
	}
	if o.ReplicaLabelNames != nil {
		options.ReplicaLabelNames = o.ReplicaLabelNames
	}
}

// Override the Options overrides the spec field when it is empty
func (o *QueryOptions) Override(spec *v1alpha1.QuerySpec) {
	o.CommonOptions.Override(&spec.CommonSpec)
	o.Envoy.Override(&spec.Envoy)

	if spec.Stores == nil {
		spec.Stores = o.Stores
	}
	if spec.SelectorLabels == nil {
		spec.SelectorLabels = o.SelectorLabels
	}
	if o.ReplicaLabelNames == nil {
		spec.ReplicaLabelNames = o.ReplicaLabelNames
	}
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
	CommonOptions `json:",inline" yaml:",inline"  mapstructure:",squash"`

	CacheConfig *v1alpha1.ResponseCacheProviderConfig `json:"cacheConfig,omitempty" yaml:"cacheConfig,omitempty"`
}

func NewQueryFrontendOptions() *QueryFrontendOptions {
	return &QueryFrontendOptions{
		CommonOptions: NewCommonOptions(),
	}
}

func (o *QueryFrontendOptions) ApplyTo(options *QueryFrontendOptions) {
	o.CommonOptions.ApplyTo(&options.CommonOptions)
	if o.CacheConfig != nil {
		options.CacheConfig = o.CacheConfig
	}
}

// Override the Options overrides the spec field when it is empty
func (o *QueryFrontendOptions) Override(spec *v1alpha1.QueryFrontendSpec) {
	o.CommonOptions.Override(&spec.CommonSpec)
	if spec.CacheConfig == nil {
		spec.CacheConfig = o.CacheConfig
	}
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
	CommonOptions `json:",inline" yaml:",inline"  mapstructure:",squash"`

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

	if o.ReplicationFactor != nil {
		options.ReplicationFactor = o.ReplicationFactor
	}
}

// Override the Options overrides the spec field when it is empty
func (o *RouterOptions) Override(spec *v1alpha1.RouterSpec) {
	o.CommonOptions.Override(&spec.CommonSpec)

	if spec.ReplicationFactor == nil {
		spec.ReplicationFactor = o.ReplicationFactor
	}
}

func (o *RouterOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *RouterOptions) AddFlags(fs *pflag.FlagSet, ro *RouterOptions) {
	var factor uint64
	o.CommonOptions.AddFlags(fs, &ro.CommonOptions, "router")
	fs.Uint64Var(&factor, "router.replicationFactor", *ro.ReplicationFactor, "Whizard router replication factor.")

	ro.ReplicationFactor = &factor
}

type RulerOptions struct {
	CommonOptions `json:",inline" yaml:",inline"  mapstructure:",squash"`

	PrometheusConfigReloader SidecarOptions `json:"prometheusConfigReloader,omitempty" yaml:"prometheusConfigReloader,omitempty"`
	RulerQueryProxy          SidecarOptions `json:"rulerQueryProxy" yaml:"rulerQueryProxy,omitempty"`

	// Number of shards to take the hash of fully qualified name of the rule group in order to split rules.
	// Each shard of rules will be bound to one separate statefulset.
	Shards *int32 `json:"shards,omitempty"`
	// Label selectors to select which PrometheusRules to mount for alerting and recording.
	// The result of multiple selectors are ORed.
	RuleSelectors []*metav1.LabelSelector `json:"ruleSelector,omitempty"`
	// Namespaces to be selected for PrometheusRules discovery. If unspecified, only
	// the same namespace as the Ruler object is in is used.
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`

	// Labels configure the external label pairs to Ruler. A default replica label
	// `ruler_replica` will be always added  as a label with the value of the pod's name and it will be dropped in the alerts.
	Labels map[string]string `json:"labels,omitempty"`
	// AlertDropLabels configure the label names which should be dropped in Ruler alerts.
	// The replica label `ruler_replica` will always be dropped in alerts.
	AlertDropLabels []string `json:"alertDropLabels,omitempty"`
	// Define URLs to send alerts to Alertmanager.
	// Note: this field will be ignored if AlertmanagersConfig is specified.
	// Maps to the `alertmanagers.url` arg.
	AlertmanagersURL []string `json:"alertmanagersUrl,omitempty"`
	// Define configuration for connecting to alertmanager. Maps to the `alertmanagers.config` arg.
	AlertmanagersConfig *corev1.SecretKeySelector `json:"alertmanagersConfig,omitempty"`

	// Interval between consecutive evaluations.
	EvaluationInterval time.Duration `json:"evaluationInterval,omitempty"`
}

func NewRulerOptions() *RulerOptions {
	var shards int32 = DefaultRulerShards
	return &RulerOptions{
		CommonOptions:      NewCommonOptions(),
		Shards:             &shards,
		EvaluationInterval: DefaultRulerEvaluationInterval,

		PrometheusConfigReloader: SidecarOptions{
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
		RulerQueryProxy: SidecarOptions{
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
	o.PrometheusConfigReloader.ApplyTo(&options.PrometheusConfigReloader)
	o.RulerQueryProxy.ApplyTo(&options.RulerQueryProxy)

	if *o.Shards != 0 {
		options.Shards = o.Shards
	}
	if o.RuleSelectors != nil {
		options.RuleSelectors = o.RuleSelectors
	}
	if o.RuleNamespaceSelector != nil {
		options.RuleNamespaceSelector = o.RuleNamespaceSelector
	}
	if o.Labels != nil {
		options.Labels = o.Labels
	}
	if o.AlertDropLabels != nil {
		options.AlertDropLabels = o.AlertDropLabels
	}
	if o.AlertmanagersURL != nil {
		options.AlertmanagersURL = o.AlertmanagersURL
	}
	if o.AlertmanagersConfig != nil {
		options.AlertmanagersConfig = o.AlertmanagersConfig
	}
	if o.EvaluationInterval != 0 {
		options.EvaluationInterval = o.EvaluationInterval
	}

}

// Override the Options overrides the spec field when it is empty
func (o *RulerOptions) Override(spec *v1alpha1.RulerSpec) {
	o.CommonOptions.Override(&spec.CommonSpec)
	o.PrometheusConfigReloader.Override(&spec.PrometheusConfigReloader)
	o.RulerQueryProxy.Override(&spec.RulerQueryProxy)

	if spec.Shards == nil {
		spec.Shards = o.Shards
	}
	if spec.RuleSelectors == nil {
		spec.RuleSelectors = o.RuleSelectors
	}
	if spec.RuleNamespaceSelector == nil {
		spec.RuleNamespaceSelector = o.RuleNamespaceSelector
	}
	if spec.Labels == nil {
		spec.Labels = o.Labels
	}
	if spec.AlertDropLabels == nil {
		spec.AlertDropLabels = o.AlertDropLabels
	}
	if spec.AlertmanagersConfig == nil {
		spec.AlertmanagersConfig = o.AlertmanagersConfig
	}
	if spec.EvaluationInterval == "" {
		spec.EvaluationInterval = v1alpha1.Duration(o.EvaluationInterval.String())
	}
}

func (o *RulerOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.CommonOptions.Validate()...)

	return errs
}

func (o *RulerOptions) AddFlags(fs *pflag.FlagSet, ro *RulerOptions) {
	o.CommonOptions.AddFlags(fs, &ro.CommonOptions, "ruler")
	o.PrometheusConfigReloader.AddFlags(fs, &ro.PrometheusConfigReloader, "ruler.prometheus-config-reloader")
	o.RulerQueryProxy.AddFlags(fs, &ro.RulerQueryProxy, "ruler.query-proxy")
}

type StoreOptions struct {
	CommonOptions `json:",inline" yaml:",inline"  mapstructure:",squash"`

	// MinTime specifies start of time range limit to serve
	MinTime string `json:"minTime,omitempty" yaml:"minTime,omitempty"`
	// MaxTime specifies end of time range limit to serve
	MaxTime                    string `json:"maxTime,omitempty" yaml:"maxTime,omitempty"`
	*v1alpha1.KubernetesVolume `json:"dataVolume,omitempty" yaml:"dataVolume,omitempty"`
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

	if o.MinTime != "" {
		options.MinTime = o.MinTime
	}
	if o.MaxTime != "" {
		options.MaxTime = o.MaxTime
	}

	if o.KubernetesVolume != nil {
		if options.KubernetesVolume == nil {
			options.KubernetesVolume = o.KubernetesVolume
		} else {
			if o.KubernetesVolume.PersistentVolumeClaim != nil {
				options.KubernetesVolume.PersistentVolumeClaim = o.KubernetesVolume.PersistentVolumeClaim
			}

			if o.KubernetesVolume.EmptyDir != nil {
				options.KubernetesVolume.EmptyDir = o.KubernetesVolume.EmptyDir
			}
		}
	}

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

// Override the Options overrides the spec field when it is empty
func (o *StoreOptions) Override(spec *v1alpha1.StoreSpec) {
	o.CommonOptions.Override(&spec.CommonSpec)

	if spec.MinTime == "" {
		spec.MinTime = o.MinTime
	}
	if spec.MaxTime != "" {
		spec.MaxTime = o.MaxTime
	}

	if spec.DataVolume == nil {
		spec.DataVolume = o.KubernetesVolume
	}
	if spec.IndexCacheConfig == nil {
		spec.IndexCacheConfig = o.IndexCacheConfig
	} else {
		if spec.IndexCacheConfig.InMemoryIndexCacheConfig == nil {
			spec.IndexCacheConfig.InMemoryIndexCacheConfig = o.IndexCacheConfig.InMemoryIndexCacheConfig
		} else {
			if spec.IndexCacheConfig.MaxSize == "" {
				spec.IndexCacheConfig.MaxSize = o.MaxSize
			}
		}
	}

	if spec.Scaler == nil {
		spec.Scaler = o.AutoScaler
	} else {
		if spec.Scaler.MaxReplicas == 0 {
			spec.Scaler.MaxReplicas = o.MaxReplicas
		}

		if spec.Scaler.MinReplicas == nil || *spec.Scaler.MinReplicas == 0 {
			min := *o.MinReplicas
			spec.Scaler.MinReplicas = &min
		}

		if spec.Scaler.Metrics == nil {
			spec.Scaler.Metrics = o.Metrics
		}

		if spec.Scaler.Behavior == nil {
			spec.Scaler.Behavior = o.Behavior
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
	BlockManager *BlockManagerOptions `json:"blockManager,omitempty"`
}

type BlockManagerOptions struct {
	Enable             *bool `json:"enable,omitempty"`
	CommonOptions      `json:",inline"`
	ServiceAccountName string           `json:"serviceAccountName,omitempty"`
	BlockSyncInterval  *metav1.Duration `json:"blockSyncInterval,omitempty"`
	GC                 *BlockGCOptions  `json:"gc,omitempty"`
}

type BlockGCOptions struct {
	Enable          *bool                       `json:"enable,omitempty"`
	Resources       corev1.ResourceRequirements `json:"resources,omitempty"`
	Image           string                      `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	GCInterval      *metav1.Duration            `json:"gcInterval,omitempty"`
	CleanupTimeout  *metav1.Duration            `json:"cleanupTimeout,omitempty"`
}

func NewStorageOptions() *StorageOptions {
	enable := true
	blockSyncInterval := metav1.Duration{time.Minute}
	return &StorageOptions{
		BlockManager: &BlockManagerOptions{
			Enable:            &enable,
			CommonOptions:     NewCommonOptions(),
			BlockSyncInterval: &blockSyncInterval,
			GC: &BlockGCOptions{
				Enable: &enable,
				Image:  DefaultBlockManagerImage,
			},
			ServiceAccountName: DefaultServiceAccount,
		},
	}
}

func (o *StorageOptions) ApplyTo(options *StorageOptions) {
	if o.BlockManager != nil {
		if options.BlockManager == nil {
			options.BlockManager = o.BlockManager
		} else {
			o.BlockManager.CommonOptions.ApplyTo(&options.BlockManager.CommonOptions)

			if options.BlockManager.Enable == nil {
				options.BlockManager.Enable = o.BlockManager.Enable
			}

			if options.BlockManager.BlockSyncInterval == nil || options.BlockManager.BlockSyncInterval.Duration == 0 {
				options.BlockManager.BlockSyncInterval = o.BlockManager.BlockSyncInterval
			}
			if options.BlockManager.ServiceAccountName == "" {
				options.BlockManager.ServiceAccountName = o.BlockManager.ServiceAccountName
			}

			if o.BlockManager.GC != nil {
				if options.BlockManager.GC == nil {
					options.BlockManager.GC = o.BlockManager.GC
				} else {
					if options.BlockManager.GC.Image == "" {
						options.BlockManager.GC.Image = o.BlockManager.GC.Image
					}

					if options.BlockManager.GC.ImagePullPolicy == "" {
						options.BlockManager.GC.ImagePullPolicy = o.BlockManager.GC.ImagePullPolicy
					}

					if options.BlockManager.GC.Resources.Limits == nil {
						options.BlockManager.GC.Resources.Limits = o.BlockManager.GC.Resources.Limits
					}

					if options.BlockManager.GC.Resources.Requests == nil {
						options.BlockManager.GC.Resources.Requests = o.BlockManager.GC.Resources.Requests
					}

					if options.BlockManager.GC.GCInterval == nil ||
						options.BlockManager.GC.GCInterval.Duration == 0 {
						options.BlockManager.GC.GCInterval = o.BlockManager.GC.GCInterval
					}

					if options.BlockManager.GC.CleanupTimeout == nil ||
						options.BlockManager.GC.CleanupTimeout.Duration == 0 {
						options.BlockManager.GC.CleanupTimeout = o.BlockManager.GC.CleanupTimeout
					}

					if options.BlockManager.GC.Enable == nil {
						options.BlockManager.GC.Enable = o.BlockManager.GC.Enable
					}
				}
			}
		}
	}
}

func (o *StorageOptions) Override(spec *v1alpha1.StorageSpec) {
	o.BlockManager.CommonOptions.Override(&spec.BlockManager.CommonSpec)

	if spec.BlockManager.BlockSyncInterval == nil || spec.BlockManager.BlockSyncInterval.Duration == 0 {
		spec.BlockManager.BlockSyncInterval = o.BlockManager.BlockSyncInterval
	}

	if spec.BlockManager.ServiceAccountName == "" {
		spec.BlockManager.ServiceAccountName = o.BlockManager.ServiceAccountName
	}

	if spec.BlockManager.GC != nil &&
		spec.BlockManager.GC.Enable != nil &&
		*spec.BlockManager.GC.Enable {
		if spec.BlockManager.GC.Image == "" {
			spec.BlockManager.GC.Image = o.BlockManager.GC.Image
		}
		if spec.BlockManager.GC.ImagePullPolicy == "" {
			spec.BlockManager.GC.ImagePullPolicy = o.BlockManager.GC.ImagePullPolicy
		}
		if spec.BlockManager.GC.Resources.Limits == nil {
			spec.BlockManager.GC.Resources.Limits = o.BlockManager.GC.Resources.Limits
		}
		if spec.BlockManager.GC.Resources.Requests == nil {
			spec.BlockManager.GC.Resources.Requests = o.BlockManager.GC.Resources.Requests
		}
		if spec.BlockManager.GC.GCInterval == nil ||
			spec.BlockManager.GC.GCInterval.Duration == 0 {
			spec.BlockManager.GC.GCInterval = o.BlockManager.GC.GCInterval
		}
		if spec.BlockManager.GC.CleanupTimeout == nil ||
			spec.BlockManager.GC.GCInterval.Duration == 0 {
			spec.BlockManager.GC.CleanupTimeout = o.BlockManager.GC.CleanupTimeout
		}
	}
}

func (o *StorageOptions) Validate() []error {
	var errs []error
	if o.BlockManager != nil {
		errs = append(errs, o.BlockManager.CommonOptions.Validate()...)
	}

	return errs
}

func (o *StorageOptions) AddFlags(fs *pflag.FlagSet, s *StorageOptions) {
	if o.BlockManager != nil && s.BlockManager != nil {
		o.BlockManager.CommonOptions.AddFlags(fs, &s.BlockManager.CommonOptions, "storage")
	}
}
