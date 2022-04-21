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

type Options struct {
	ThanosImage                  string `json:"thanosImage,omitempty" yaml:"thanosImage,omitempty"`
	EnvoyImage                   string `json:"envoyImage,omitempty" yaml:"envoyImage,omitempty"`
	PaodinMonitoringGatewayImage string `json:"paodinMonitoringGatewayImage,omitempty" yaml:"paodinMonitoringGatewayImage,omitempty"`
}

func NewOptions() *Options {
	return &Options{
		ThanosImage:                  ThanosDefaultImage,
		EnvoyImage:                   EnvoyDefaultImage,
		PaodinMonitoringGatewayImage: PaodinMonitoringGatewayDefaultImage,
	}
}

func (o *Options) Validate() []error {
	var errs []error
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
}

func (o *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	flag.StringVar(&c.ThanosImage, "thanos-image", ThanosDefaultImage, "Thanos image with tag/version")
	flag.StringVar(&c.EnvoyImage, "envoy-image", EnvoyDefaultImage, "Envoy image with tag/version")
	flag.StringVar(&c.PaodinMonitoringGatewayImage, "paodin-monitoring-gateway-image", PaodinMonitoringGatewayDefaultImage, "Paodin monitoring gateway image with tag/version")
}
