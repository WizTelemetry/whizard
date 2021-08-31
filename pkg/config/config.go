package config

import (
	"flag"
)

const (
	ThanosDefaultImage = "thanosio/thanos:v0.22.0"
	EnvoyDefaultImage  = "envoyproxy/envoy:v1.19.0"
)

type Config struct {
	ThanosDefaultImage string
	EnvoyDefaultImage  string
}

func (c *Config) AddFlags() {
	flag.StringVar(&c.ThanosDefaultImage, "thanos-default-image", ThanosDefaultImage, "Thanos default image with tag/version")
	flag.StringVar(&c.EnvoyDefaultImage, "envoy-default-image", EnvoyDefaultImage, "Envoy default image with tag/version")
}
