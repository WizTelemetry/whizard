package monitoringgateway

import (
	"net"
	"net/http"
	"os"
	"time"

	extflag "github.com/efficientgo/tools/extkingpin"
	"github.com/pkg/errors"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/thanos-io/thanos/pkg/extkingpin"
	"github.com/thanos-io/thanos/pkg/httpconfig"

	"gopkg.in/yaml.v2"
)

type RemoteWriteConfig struct {
	Name          string            `yaml:"name,omitempty"`
	URL           *config.URL       `yaml:"url"`
	Headers       map[string]string `yaml:"headers,omitempty"`
	RemoteTimeout model.Duration    `yaml:"remote_timeout,omitempty"`
	TLSConfig     config.TLSConfig  `yaml:"tls_config,omitempty"`
}

// LoadRemoteWritesConfig loads remotewrites config, and prefers file to content
func LoadRemoteWritesConfig(file, content string) ([]RemoteWriteConfig, error) {
	var buff []byte
	if file != "" {
		c, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		buff = c
	} else {
		buff = []byte(content)
	}
	if len(buff) == 0 {
		return nil, nil
	}
	var rws []RemoteWriteConfig
	if err := yaml.UnmarshalStrict(buff, &rws); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return rws, nil
}

type QueryConfig struct {
	DownstreamURL string

	DownstreamTripperConfig
}

// RegisterCommonObjStoreFlags register flags commonly used to configure http servers with.
func RegisterHTTPFlags(cmd extkingpin.FlagClause) (httpBindAddr *string, httpGracePeriod *model.Duration, httpTLSConfig *string) {
	httpBindAddr = cmd.Flag("http-address", "Listen host:port for HTTP endpoints.").Default("0.0.0.0:9090").String()
	httpGracePeriod = extkingpin.ModelDuration(cmd.Flag("http-grace-period", "Time to wait after an interrupt received for HTTP Server.").Default("2m")) // by default it's the same as query.timeout.
	httpTLSConfig = cmd.Flag(
		"http.config",
		"[EXPERIMENTAL] Path to the configuration file that can enable TLS or authentication for all HTTP endpoints.",
	).Default("").String()
	return httpBindAddr, httpGracePeriod, httpTLSConfig
}

func (qc *QueryConfig) RegisterFlag(cmd extflag.FlagClause) *QueryConfig {
	cmd.Flag("query.address", "Addresses of statically configured query API servers (repeatable). The scheme may be prefixed with 'dns+' or 'dnssrv+' to detect query API servers through respective DNS lookups.").
		PlaceHolder("<query>").StringVar(&qc.DownstreamURL)
	qc.DownstreamTripperConfig.TripperPathOrContent = *extflag.RegisterPathOrContent(cmd, "query.config", "YAML file that contains downstream tripper configuration. If your downstream URL is localhost or 127.0.0.1 then it is highly recommended to increase max_idle_conns_per_host to at least 100.", extflag.WithEnvSubstitution())

	return qc
}

type RulesQueryConfig struct {
	DownstreamURL string

	DownstreamTripperConfig
}

func (rc *RulesQueryConfig) RegisterFlag(cmd extflag.FlagClause) *RulesQueryConfig {
	cmd.Flag("rulesquery.address", "Addresses of statically configured query API servers (repeatable). The scheme may be prefixed with 'dns+' or 'dnssrv+' to detect query API servers through respective DNS lookups.").
		PlaceHolder("<query>").StringVar(&rc.DownstreamURL)

	rc.DownstreamTripperConfig.TripperPathOrContent = *extflag.RegisterPathOrContent(cmd, "rulesquery.downstream-tripper-config", "YAML file that contains downstream tripper configuration. If your downstream URL is localhost or 127.0.0.1 then it is highly recommended to increase max_idle_conns_per_host to at least 100.", extflag.WithEnvSubstitution())

	return rc
}

// DownstreamTripperConfig stores the http.Transport configuration for query's HTTP downstream tripper.
type DownstreamTripperConfig struct {
	IdleConnTimeout       model.Duration        `yaml:"idle_conn_timeout"`
	ResponseHeaderTimeout model.Duration        `yaml:"response_header_timeout"`
	TLSHandshakeTimeout   model.Duration        `yaml:"tls_handshake_timeout"`
	ExpectContinueTimeout model.Duration        `yaml:"expect_continue_timeout"`
	MaxIdleConns          *int                  `yaml:"max_idle_conns"`
	MaxIdleConnsPerHost   *int                  `yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost       *int                  `yaml:"max_conns_per_host"`
	TLSConfig             *httpconfig.TLSConfig `yaml:"tls_config"`

	TripperPathOrContent extflag.PathOrContent
}

func ParseTransportConfiguration(downstreamTripperConfContentYaml []byte) (*http.Transport, error) {

	downstreamTripper := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if len(downstreamTripperConfContentYaml) > 0 {
		tripperConfig := &DownstreamTripperConfig{}
		if err := yaml.UnmarshalStrict(downstreamTripperConfContentYaml, tripperConfig); err != nil {
			return nil, errors.Wrap(err, "parsing downstream tripper config YAML file")
		}

		if tripperConfig.IdleConnTimeout > 0 {
			downstreamTripper.IdleConnTimeout = time.Duration(tripperConfig.IdleConnTimeout)
		}
		if tripperConfig.ResponseHeaderTimeout > 0 {
			downstreamTripper.ResponseHeaderTimeout = time.Duration(tripperConfig.ResponseHeaderTimeout)
		}
		if tripperConfig.TLSHandshakeTimeout > 0 {
			downstreamTripper.TLSHandshakeTimeout = time.Duration(tripperConfig.TLSHandshakeTimeout)
		}
		if tripperConfig.ExpectContinueTimeout > 0 {
			downstreamTripper.ExpectContinueTimeout = time.Duration(tripperConfig.ExpectContinueTimeout)
		}
		if tripperConfig.MaxIdleConns != nil {
			downstreamTripper.MaxIdleConns = *tripperConfig.MaxIdleConns
		}
		if tripperConfig.MaxIdleConnsPerHost != nil {
			downstreamTripper.MaxIdleConnsPerHost = *tripperConfig.MaxIdleConnsPerHost
		}
		if tripperConfig.MaxConnsPerHost != nil {
			downstreamTripper.MaxConnsPerHost = *tripperConfig.MaxConnsPerHost
		}
		if tripperConfig.TLSConfig != nil {
			tlsConfig, err := config.NewTLSConfig(&config.TLSConfig{
				CAFile:             tripperConfig.TLSConfig.CAFile,
				CertFile:           tripperConfig.TLSConfig.CertFile,
				KeyFile:            tripperConfig.TLSConfig.KeyFile,
				ServerName:         tripperConfig.TLSConfig.ServerName,
				InsecureSkipVerify: tripperConfig.TLSConfig.InsecureSkipVerify,
			})
			if err != nil {
				return nil, err
			}
			downstreamTripper.TLSClientConfig = tlsConfig
		}
	}

	return downstreamTripper, nil
}
