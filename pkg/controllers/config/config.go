package config

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var (
	// singleton instance of config package
	_config = defaultConfig()

	// decodeHook is used to configure mapstructure.DecoderConfig options
	decodeHook = mapstructure.ComposeDecodeHookFunc(
		StringToResourceQuantityHookFunc(),
		StringToMetav1TimeHookFunc(),
		StringToMetav1DurationHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)
)

const (
	// defaultConfigurationName is the default name of configuration
	defaultConfigurationName = "whizard"
	// defaultConfigurationPath the default location of the configuration file
	defaultConfigurationPath = "/etc/whizard"
)

type config struct {
	cfg         *Config
	cfgChangeCh chan Config
	watchOnce   sync.Once
	loadOnce    sync.Once
}

func (c *config) watchConfig() <-chan Config {
	c.watchOnce.Do(func() {
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			cfg := New()
			if err := viper.Unmarshal(cfg, viper.DecodeHook(decodeHook)); err != nil {
				klog.Warning("config reload error", err)
			} else {
				c.cfgChangeCh <- *cfg
			}
		})
	})
	return c.cfgChangeCh
}

func (c *config) loadFromDisk() (*Config, error) {
	var err error
	c.loadOnce.Do(func() {
		if err = viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				err = fmt.Errorf("error parsing configuration file %s", err)
			}
		}
		err = viper.Unmarshal(c.cfg, viper.DecodeHook(decodeHook))
	})
	return c.cfg, err
}

func defaultConfig() *config {
	viper.SetConfigName(defaultConfigurationName)
	viper.AddConfigPath(defaultConfigurationPath)

	// Load from current working directory, only used for debugging
	viper.AddConfigPath(".")

	// Load from Environment variables
	viper.SetEnvPrefix("whizard")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &config{
		cfg:         New(),
		cfgChangeCh: make(chan Config),
		watchOnce:   sync.Once{},
		loadOnce:    sync.Once{},
	}
}

type Config struct {
}

func New() *Config {
	return &Config{}
}

// TryLoadFromDisk loads configuration from default location after server startup
// return nil error if configuration file not exists
func TryLoadFromDisk() (*Config, error) {
	return _config.loadFromDisk()
}

// WatchConfigChange return config change channel
func WatchConfigChange() <-chan Config {
	return _config.watchConfig()
}

// ToMap simply converts config to map[string]bool
// to hide sensitive information
func (conf *Config) ToMap() map[string]bool {
	conf.stripEmptyOptions()
	result := make(map[string]bool, 0)

	if conf == nil {
		return result
	}

	c := reflect.Indirect(reflect.ValueOf(conf))

	for i := 0; i < c.NumField(); i++ {
		name := strings.Split(c.Type().Field(i).Tag.Get("json"), ",")[0]
		if strings.HasPrefix(name, "-") {
			continue
		}

		if c.Field(i).IsNil() {
			result[name] = false
		} else {
			result[name] = true
		}
	}

	return result
}

// Remove invalid options before serializing to json or yaml
func (conf *Config) stripEmptyOptions() {
}

// StringToResourceQuantityHookFunc returns a DecodeHookFunc that converts
// strings to resource.Quantity.
func StringToResourceQuantityHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(resource.MustParse("100Mi")) {
			return data, nil
		}

		// Convert it by parsing
		return resource.ParseQuantity(data.(string))
	}
}

func StringToMetav1TimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(metav1.Now()) {
			return data, nil
		}

		// Convert it by parsing
		pt := metav1.Time{}
		pt.UnmarshalQueryParameter(data.(string))

		return pt, nil
	}

}

func StringToMetav1DurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(metav1.Duration{Duration: time.Minute}) {
			return data, nil
		}

		d, err := time.ParseDuration(data.(string))
		if err != nil {
			return nil, err
		}
		return metav1.Duration{Duration: d}, nil
	}
}
