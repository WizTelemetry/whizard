package options

import (
	"flag"

	"github.com/spf13/pflag"
)

const (
	ThanosDefaultImage                  = "thanosio/thanos:v0.25.2"
	EnvoyDefaultImage                   = "envoyproxy/envoy:v1.20.2"
	PaodinMonitoringGatewayDefaultImage = "junotx/paodin-monitoring-gateway:latest"
)

var PrometheusConfigReloaderDefaultConfig = PrometheusConfigReloaderConfig{
	Image:         "quay.io/prometheus-operator/prometheus-config-reloader:v0.55.1",
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

type Options struct {
	ThanosImage                  string `json:"thanosImage,omitempty" yaml:"thanosImage,omitempty"`
	EnvoyImage                   string `json:"envoyImage,omitempty" yaml:"envoyImage,omitempty"`
	PaodinMonitoringGatewayImage string `json:"paodinMonitoringGatewayImage,omitempty" yaml:"paodinMonitoringGatewayImage,omitempty"`

	PrometheusConfigReloader PrometheusConfigReloaderConfig `json:"prometheusConfigReloader,omitempty" yaml:"prometheusConfigReloader,omitempty"`
}

func NewOptions() *Options {
	return &Options{
		ThanosImage:                  ThanosDefaultImage,
		EnvoyImage:                   EnvoyDefaultImage,
		PaodinMonitoringGatewayImage: PaodinMonitoringGatewayDefaultImage,
		PrometheusConfigReloader:     PrometheusConfigReloaderDefaultConfig,
	}
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.PrometheusConfigReloader.Validate()...)
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
	o.PrometheusConfigReloader.ApplyTo(&options.PrometheusConfigReloader)
}

func (o *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	flag.StringVar(&c.ThanosImage, "thanos-image", ThanosDefaultImage, "Thanos image with tag/version")
	flag.StringVar(&c.EnvoyImage, "envoy-image", EnvoyDefaultImage, "Envoy image with tag/version")
	flag.StringVar(&c.PaodinMonitoringGatewayImage, "paodin-monitoring-gateway-image", PaodinMonitoringGatewayDefaultImage, "Paodin monitoring gateway image with tag/version")

	flag.StringVar(&c.PrometheusConfigReloader.Image, "prometheus-config-reloader-image",
		PrometheusConfigReloaderDefaultConfig.Image, "Prometheus Config Reloader image with tag/version")
	flag.StringVar(&c.PrometheusConfigReloader.CPURequest, "prometheus-config-reloader-cpu-request",
		PrometheusConfigReloaderDefaultConfig.CPURequest, "Prometheus Config Reloader CPU request. Value \"0\" disables it and causes no request to be configured.")
	flag.StringVar(&c.PrometheusConfigReloader.CPULimit, "prometheus-config-reloader-cpu-limit",
		PrometheusConfigReloaderDefaultConfig.CPULimit, "Prometheus Config Reloader CPU limit. Value \"0\" disables it and causes no limit to be configured.")
	flag.StringVar(&c.PrometheusConfigReloader.MemoryRequest, "prometheus-config-reloader-memory-request",
		PrometheusConfigReloaderDefaultConfig.MemoryRequest, "Prometheus Config Reloader Memory request. Value \"0\" disables it and causes no request to be configured.")
	flag.StringVar(&c.PrometheusConfigReloader.MemoryLimit, "prometheus-config-reloader-memory-limit",
		PrometheusConfigReloaderDefaultConfig.MemoryLimit, "Prometheus Config Reloader Memory limit. Value \"0\" disables it and causes no limit to be configured.")
}