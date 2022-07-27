package options

import (
	"time"

	"github.com/spf13/pflag"
)

const (
	ThanosDefaultImage                   = "thanosio/thanos:v0.26.0"
	EnvoyDefaultImage                    = "envoyproxy/envoy:v1.20.2"
	PaodinMonitoringGatewayDefaultImage  = "kubesphere/paodin-monitoring-gateway:latest"
	PaodinDefaultService                 = "kubesphere-monitoring-system.central"
	PaodinDefaultTenantsPerIngestor      = 3
	PaodinDefaultIngestorRetentionPeriod = "3h"
)

var PrometheusConfigReloaderDefaultConfig = PrometheusConfigReloaderConfig{
	Image:         "quay.io/prometheus-operator/prometheus-config-reloader:v0.55.1",
	CPURequest:    "100m",
	MemoryRequest: "50Mi",
	CPULimit:      "100m",
	MemoryLimit:   "50Mi",
}

var RulerQueryProxyDefaultConfig = RulerQueryProxyConfig{
	Image:         PaodinMonitoringGatewayDefaultImage,
	CPURequest:    "100m",
	MemoryRequest: "50Mi",
	CPULimit:      "100m",
	MemoryLimit:   "50Mi",
}

type PrometheusConfigReloaderConfig struct {
	Image         string `json:"image,omitempty" yaml:"image,omitempty"`
	CPURequest    string `json:"cpuRequest,omitempty" yaml:"cpuRequest,omitempty"`
	MemoryRequest string `json:"memoryRequest,omitempty" yaml:"memoryRequest,omitempty"`
	CPULimit      string `json:"cpuLimit,omitempty" yaml:"cpuRequest,omitempty"`
	MemoryLimit   string `json:"memoryLimit,omitempty" yaml:"memoryLimit,omitempty"`
}

func (o *PrometheusConfigReloaderConfig) Validate() []error {
	var errs []error
	return errs
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

type Options struct {
	ThanosImage                  string `json:"thanosImage,omitempty" yaml:"thanosImage,omitempty"`
	EnvoyImage                   string `json:"envoyImage,omitempty" yaml:"envoyImage,omitempty"`
	PaodinMonitoringGatewayImage string `json:"paodinMonitoringGatewayImage,omitempty" yaml:"paodinMonitoringGatewayImage,omitempty"`

	EnableKubeSphereAdapter        bool   `json:"enableKubeSphereAdapter,omitempty" yaml:"enableKubeSphereAdapter,omitempty"`
	KubeSphereAdapterService       string `json:"kubeSphereAdapterService,omitempty" yaml:"kubeSphereAdapterService,omitempty"`
	DefaultTenantsPerIngestor      int    `json:"defaultTenantsPerIngestor,omitempty" yaml:"defaultTenantsPerIngestor,omitempty"`
	DefaultIngestorRetentionPeriod string `json:"defaultIngestorRetentionPeriod,omitempty" yaml:"defaultIngestorRetentionPeriod,omitempty"`

	PrometheusConfigReloader PrometheusConfigReloaderConfig `json:"prometheusConfigReloader,omitempty" yaml:"prometheusConfigReloader,omitempty"`
	RulerQueryProxy          RulerQueryProxyConfig          `json:"rulerQueryProxy,omitempty" yaml:"rulerQueryProxy,omitempty"`
}

func NewOptions() *Options {
	return &Options{
		ThanosImage:                    ThanosDefaultImage,
		EnvoyImage:                     EnvoyDefaultImage,
		PaodinMonitoringGatewayImage:   PaodinMonitoringGatewayDefaultImage,
		DefaultTenantsPerIngestor:      PaodinDefaultTenantsPerIngestor,
		DefaultIngestorRetentionPeriod: PaodinDefaultIngestorRetentionPeriod,
		EnableKubeSphereAdapter:        true,
		KubeSphereAdapterService:       PaodinDefaultService,
		PrometheusConfigReloader:       PrometheusConfigReloaderDefaultConfig,
		RulerQueryProxy:                RulerQueryProxyDefaultConfig,
	}
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.PrometheusConfigReloader.Validate()...)
	if _, err := time.ParseDuration(o.DefaultIngestorRetentionPeriod); err != nil {
		errs = append(errs, err)
	}
	return errs
}

func (o *Options) ApplyTo(options *Options) {
	if o.ThanosImage != "" {
		options.ThanosImage = o.ThanosImage
	}
	if o.EnvoyImage != "" {
		options.EnvoyImage = o.EnvoyImage
	}
	if o.PaodinMonitoringGatewayImage != "" {
		options.PaodinMonitoringGatewayImage = o.PaodinMonitoringGatewayImage
	}
	if o.DefaultTenantsPerIngestor != 0 {
		options.DefaultTenantsPerIngestor = o.DefaultTenantsPerIngestor
	}
	if o.DefaultIngestorRetentionPeriod != "" {
		options.DefaultIngestorRetentionPeriod = o.DefaultIngestorRetentionPeriod
	}
	if o.KubeSphereAdapterService != "" {
		options.KubeSphereAdapterService = o.KubeSphereAdapterService
	}
	options.EnableKubeSphereAdapter = o.EnableKubeSphereAdapter

	o.PrometheusConfigReloader.ApplyTo(&options.PrometheusConfigReloader)
}

func (o *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	fs.StringVar(&c.ThanosImage, "thanos-image", ThanosDefaultImage, "Thanos image with tag/version")
	fs.StringVar(&c.EnvoyImage, "envoy-image", EnvoyDefaultImage, "Envoy image with tag/version")
	fs.StringVar(&c.PaodinMonitoringGatewayImage, "paodin-monitoring-gateway-image", PaodinMonitoringGatewayDefaultImage, "Paodin monitoring gateway image with tag/version")
	fs.IntVar(&c.DefaultTenantsPerIngestor, "defaultTenantsPerIngestor", PaodinDefaultTenantsPerIngestor, "Paodin default tenant count per ingestor. (default 3)")
	fs.StringVar(&c.DefaultIngestorRetentionPeriod, "defaultIngestorRetentionPeriod", PaodinDefaultIngestorRetentionPeriod, "Paodin default ingestor retention period. (default 2h)")
	fs.BoolVar(&c.EnableKubeSphereAdapter, "enableKubeSphereAdapter", true, "Enable KubeSphere adapter. (default true)")
	fs.StringVar(&c.KubeSphereAdapterService, "kubeSphereAdapterService", PaodinDefaultService, "Default service for tenants generated by kubesphere adapter, format is namespace.name")

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
}
