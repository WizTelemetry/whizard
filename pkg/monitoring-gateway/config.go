package monitoringgateway

import (
	"net"
	"net/http"
	"reflect"
	"time"

	extflag "github.com/efficientgo/tools/extkingpin"
	"github.com/pkg/errors"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/thanos-io/thanos/pkg/extkingpin"

	"gopkg.in/yaml.v2"
)

type ExternalRemoteWriteConfig struct {
	Name          string            `yaml:"name,omitempty"`
	URL           *config.URL       `yaml:"url"`
	Headers       map[string]string `yaml:"headers,omitempty"`
	RemoteTimeout model.Duration    `yaml:"remote_timeout,omitempty"`

	// The HTTP basic authentication credentials for the targets.
	BasicAuth *BasicAuth `yaml:"basic_auth,omitempty" json:"basic_auth,omitempty"`
	// The bearer token for the targets. Deprecated in favour of
	// Authorization.Credentials.
	BearerToken string `yaml:"bearer_token,omitempty" json:"bearer_token,omitempty"`
	// TLSConfig to use to connect to the targets.
	TLSConfig config.TLSConfig `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
}

// BasicAuth contains basic HTTP authentication credentials.
type BasicAuth struct {
	Username     string `yaml:"username" json:"username"`
	Password     string `yaml:"password,omitempty" json:"password,omitempty"`
	PasswordFile string `yaml:"password_file,omitempty" json:"password_file,omitempty"`
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

type QueryConfig struct {
	DownstreamURL string

	DownstreamTripperConfig
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
	cmd.Flag("rules-query.address", "Addresses of statically configured query API servers (repeatable). The scheme may be prefixed with 'dns+' or 'dnssrv+' to detect query API servers through respective DNS lookups.").
		PlaceHolder("<query>").StringVar(&rc.DownstreamURL)

	rc.DownstreamTripperConfig.TripperPathOrContent = *extflag.RegisterPathOrContent(cmd, "rules-query.config", "YAML file that contains downstream tripper configuration. If your downstream URL is localhost or 127.0.0.1 then it is highly recommended to increase max_idle_conns_per_host to at least 100.", extflag.WithEnvSubstitution())

	return rc
}

type RemoteWriteConfig struct {
	DownstreamURL string
	DownstreamTripperConfig
}

func (rwc *RemoteWriteConfig) RegisterFlag(cmd extflag.FlagClause) *RemoteWriteConfig {
	cmd.Flag("remote-write.address", "Address to send remote write requests.").
		PlaceHolder("<query>").StringVar(&rwc.DownstreamURL)

	rwc.TripperPathOrContent = *extflag.RegisterPathOrContent(cmd, "remote-write.config", "YAML file that contains downstream tripper configuration. If your downstream URL is localhost or 127.0.0.1 then it is highly recommended to increase max_idle_conns_per_host to at least 100.", extflag.WithEnvSubstitution())

	return rwc
}

// DownstreamTripperConfig stores the http.Transport configuration for query's HTTP downstream tripper.
type DownstreamTripperConfig struct {
	IdleConnTimeout       model.Duration          `yaml:"idle_conn_timeout"`
	ResponseHeaderTimeout model.Duration          `yaml:"response_header_timeout"`
	TLSHandshakeTimeout   model.Duration          `yaml:"tls_handshake_timeout"`
	ExpectContinueTimeout model.Duration          `yaml:"expect_continue_timeout"`
	MaxIdleConns          *int                    `yaml:"max_idle_conns"`
	MaxIdleConnsPerHost   *int                    `yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost       *int                    `yaml:"max_conns_per_host"`
	HTTPClientConfig      config.HTTPClientConfig `yaml:",inline"`

	TripperPathOrContent extflag.PathOrContent
}

func ParseTransportConfiguration(downstreamTripperConfContentYaml []byte) (http.RoundTripper, error) {

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
		if !reflect.DeepEqual(tripperConfig.HTTPClientConfig, config.HTTPClientConfig{}) {
			// todo: load DownstreamTripperConfig
			rt, err := config.NewRoundTripperFromConfig(tripperConfig.HTTPClientConfig, "")
			if err != nil {
				return nil, err
			}
			return rt, nil
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

	}

	return downstreamTripper, nil
}
