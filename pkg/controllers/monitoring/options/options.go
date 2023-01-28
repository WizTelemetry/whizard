package options

import (
	"time"

	"github.com/prometheus/common/version"
	"github.com/spf13/pflag"
)

const (
	DefaultWhizardImage                  = "thanosio/thanos:v0.30.1"
	DefaultEnvoyImage                    = "envoyproxy/envoy:v1.20.2"
	DefaultRulerWriteProxyImage          = "kubesphere/cortex-tenant:v1.7.2"
	DefaultPrometheusConfigReloaderImage = "quay.io/prometheus-operator/prometheus-config-reloader:v0.55.1"
	DefaultTSDBCleanupImage              = "bash:5.1.16"

	DefaultIngesterRetentionPeriod = time.Hour * 3
	DefaultTenantsPerIngester      = 3
	DefaultTenantsPerCompactor     = 10

	DefaultRouterReplicationFactor uint64 = 1
	DefaultRulerShards             int32  = 1
	DefaultRulerEvaluationInterval        = time.Second * 30
	DefaultStoreMinReplicas        int32  = 2
	DefaultStoreMaxReplicas        int32  = 20

	DefaultServiceAccount = "whizard-controller-manager"
)

var (
	DefaultGatewayImage      = "kubesphere/whizard-monitoring-gateway:" + version.Version
	DefaultBlockManagerImage = "kubesphere/whizard-monitoring-block-manager:" + version.Version
)

type Options struct {
	Compactor     *CompactorOptions     `json:"compactor,omitempty" yaml:"compactor,omitempty" mapstructure:"compactor"`
	Gateway       *GatewayOptions       `json:"gateway,omitempty" yaml:"gateway,omitempty" mapstructure:"gateway"`
	Ingester      *IngesterOptions      `json:"ingester,omitempty" yaml:"ingester,omitempty" mapstructure:"ingester"`
	Query         *QueryOptions         `json:"query,omitempty" yaml:"query,omitempty" mapstructure:"query"`
	QueryFrontend *QueryFrontendOptions `json:"queryFrontend,omitempty" yaml:"queryFrontend,omitempty" mapstructure:"queryFrontend"`
	Router        *RouterOptions        `json:"router,omitempty" yaml:"router,omitempty" mapstructure:"router"`
	Ruler         *RulerOptions         `json:"ruler,omitempty" yaml:"ruler,omitempty" mapstructure:"ruler"`
	Store         *StoreOptions         `json:"store,omitempty" yaml:"store,omitempty" mapstructure:"store"`
	Storage       *StorageOptions       `json:"storage,omitempty" yaml:"storage,omitempty" mapstructure:"storage"`
}

func NewOptions() *Options {
	return &Options{

		Compactor:     NewCompactorOptions(),
		Gateway:       NewGatewayOptions(),
		Ingester:      NewIngesterOptions(),
		Query:         NewQueryOptions(),
		QueryFrontend: NewQueryFrontendOptions(),
		Router:        NewRouterOptions(),
		Ruler:         NewRulerOptions(),
		Store:         NewStoreOptions(),
		Storage:       NewStorageOptions(),
	}
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.Compactor.Validate()...)
	errs = append(errs, o.Gateway.Validate()...)
	errs = append(errs, o.Ingester.Validate()...)
	errs = append(errs, o.Query.Validate()...)
	errs = append(errs, o.QueryFrontend.Validate()...)
	errs = append(errs, o.Router.Validate()...)
	errs = append(errs, o.Ruler.Validate()...)
	errs = append(errs, o.Store.Validate()...)
	errs = append(errs, o.Storage.Validate()...)
	return errs
}

func (o *Options) ApplyTo(options *Options) {

	o.Compactor.ApplyTo(options.Compactor)
	o.Gateway.ApplyTo(options.Gateway)
	o.Ingester.ApplyTo(options.Ingester)
	o.Query.ApplyTo(options.Query)
	o.QueryFrontend.ApplyTo(options.QueryFrontend)
	o.Router.ApplyTo(options.Router)
	o.Ruler.ApplyTo(options.Ruler)
	o.Store.ApplyTo(options.Store)
	o.Storage.ApplyTo(options.Storage)
}

func (o *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	o.Compactor.AddFlags(fs, o.Compactor)
	o.Gateway.AddFlags(fs, o.Gateway)
	o.Ingester.AddFlags(fs, o.Ingester)
	o.Query.AddFlags(fs, o.Query)
	o.QueryFrontend.AddFlags(fs, o.QueryFrontend)
	o.Router.AddFlags(fs, o.Router)
	o.Ruler.AddFlags(fs, o.Ruler)
	o.Store.AddFlags(fs, c.Store)
	o.Storage.AddFlags(fs, c.Storage)
}
