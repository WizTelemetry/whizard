package gateway

import (
	"context"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

const (
	TLSVersionTLS12 = "TLS12"
	TLSVersionTLS13 = "TLS13"
)

func (g *Gateway) webConfigSecret() (runtime.Object, resources.Operation, error) {
	var secret = &corev1.Secret{ObjectMeta: g.meta(g.name("web-config"))}

	if g.gateway == nil {
		return secret, resources.OperationDelete, nil
	}

	if g.gateway.Spec.WebConfig == nil {
		return secret, resources.OperationDelete, nil
	}

	// set default value
	webConfig := webConfig{
		TLSConfig: TLSConfig{
			MinVersion:               TLSVersionTLS12,
			MaxVersion:               TLSVersionTLS13,
			PreferServerCipherSuites: true,
		},
		HTTPConfig: HTTPConfig{
			HTTP2: true,
		},
	}

	if g.gateway.Spec.WebConfig.HTTPServerTLSConfig != nil {
		if g.gateway.Spec.WebConfig.HTTPServerTLSConfig.KeySecret.Name != "" {
			webConfig.TLSConfig.TLSKeyPath = constants.WhizardCertsMountPath + g.gateway.Spec.WebConfig.HTTPServerTLSConfig.KeySecret.Key
		}
		if g.gateway.Spec.WebConfig.HTTPServerTLSConfig.CertSecret.Name != "" {
			webConfig.TLSConfig.TLSCertPath = constants.WhizardCertsMountPath + g.gateway.Spec.WebConfig.HTTPServerTLSConfig.CertSecret.Key
		}
		if g.gateway.Spec.WebConfig.HTTPServerTLSConfig.ClientCASecret.Name != "" {
			webConfig.TLSConfig.ClientCAs = constants.WhizardCertsMountPath + g.gateway.Spec.WebConfig.HTTPServerTLSConfig.ClientCASecret.Key
		}
	}
	if g.gateway.Spec.WebConfig.BasicAuthUsers != nil {
		users, err := createBasicAuthUsers(g.Client, g.Context, g.gateway.Namespace, g.gateway.Spec.WebConfig.BasicAuthUsers)
		if err != nil {
			return nil, resources.OperationDelete, err
		}
		if len(users) > 0 {
			webConfig.Users = users
		}
	}

	body, err := yaml.Marshal(webConfig)
	if err != nil {
		return nil, resources.OperationDelete, err
	}

	secret.Data = map[string][]byte{
		"web-config.yaml": body,
	}

	return secret, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(g.gateway, secret, g.Scheme)
}

func createBasicAuthUsers(cli client.Client, ctx context.Context, namespace string, basicAuths []v1alpha1.BasicAuth) (map[string]string, error) {
	users := make(map[string]string)

	for _, basicAuth := range basicAuths {

		userSecret := corev1.Secret{}

		if err := cli.Get(ctx, types.NamespacedName{
			Namespace: namespace,
			Name:      basicAuth.Username.Name,
		}, &userSecret); err != nil {
			return nil, err
		}
		passSecret := corev1.Secret{}

		if err := cli.Get(ctx, types.NamespacedName{
			Namespace: namespace,
			Name:      basicAuth.Password.Name,
		}, &passSecret); err != nil {
			return nil, err
		}
		users[string(userSecret.Data[basicAuth.Username.Key])] = string(passSecret.Data[basicAuth.Password.Key])

	}

	return users, nil
}

// We have to redefine webConfig because fields such as passwords are lost when serialized.
// https://pkg.go.dev/github.com/prometheus/exporter-toolkit@v0.8.2/web#Config
type webConfig struct {
	TLSConfig  TLSConfig         `yaml:"tls_server_config"`
	HTTPConfig HTTPConfig        `yaml:"http_server_config"`
	Users      map[string]string `yaml:"basic_auth_users"`
}

type TLSConfig struct {
	TLSCertPath              string   `yaml:"cert_file"`
	TLSKeyPath               string   `yaml:"key_file"`
	ClientAuth               string   `yaml:"client_auth_type"`
	ClientCAs                string   `yaml:"client_ca_file"`
	CipherSuites             []string `yaml:"cipher_suites"`
	CurvePreferences         []string `yaml:"curve_preferences"`
	MinVersion               string   `yaml:"min_version"`
	MaxVersion               string   `yaml:"max_version"`
	PreferServerCipherSuites bool     `yaml:"prefer_server_cipher_suites"`
}

type HTTPConfig struct {
	HTTP2  bool              `yaml:"http2"`
	Header map[string]string `yaml:"headers,omitempty"`
}
