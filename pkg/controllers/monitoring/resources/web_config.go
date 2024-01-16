package resources

import (
	"context"

	"github.com/kubesphere/whizard/pkg/constants"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
)

func (r *BaseReconciler) CreateWebConfigVolumeMount(secretName string, webConfig *v1alpha1.WebConfig) (volumes []corev1.Volume, volumeMounts []corev1.VolumeMount) {
	volumes = append(volumes, corev1.Volume{
		Name: "web-config",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	})
	volumeMounts = append(volumeMounts, corev1.VolumeMount{
		Name:      "web-config",
		MountPath: constants.WhizardWebConfigMountPath,
		ReadOnly:  true,
	})

	tlsAssets := []string{}
	if webConfig.HTTPServerTLSConfig != nil {
		if webConfig.HTTPServerTLSConfig.KeySecret.Name != "" {
			tlsAssets = append(tlsAssets, webConfig.HTTPServerTLSConfig.KeySecret.Name)
		}
		if webConfig.HTTPServerTLSConfig.CertSecret.Name != "" {
			tlsAssets = append(tlsAssets, webConfig.HTTPServerTLSConfig.CertSecret.Name)
		}
		if webConfig.HTTPServerTLSConfig.ClientCASecret.Name != "" {
			tlsAssets = append(tlsAssets, webConfig.HTTPServerTLSConfig.ClientCASecret.Name)
		}
		if len(tlsAssets) > 0 {
			assetsVolume := corev1.Volume{
				Name: "tls-assets",
				VolumeSource: corev1.VolumeSource{
					Projected: &corev1.ProjectedVolumeSource{
						Sources: []corev1.VolumeProjection{},
					},
				},
			}
			for _, assetShard := range tlsAssets {
				assetsVolume.Projected.Sources = append(assetsVolume.Projected.Sources,
					corev1.VolumeProjection{
						Secret: &corev1.SecretProjection{
							LocalObjectReference: corev1.LocalObjectReference{Name: assetShard},
						},
					})
			}
			volumes = append(volumes, assetsVolume)
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      "tls-assets",
				ReadOnly:  true,
				MountPath: constants.WhizardCertsMountPath,
			})
		}
	}

	return
}

func (r *BaseReconciler) CreateWebConfig(namespace string, webConfig *v1alpha1.WebConfig) ([]byte, error) {

	// set default value
	wc := WebConfig{
		TLSConfig: TLSConfig{
			MinVersion:               TLSVersionTLS12,
			MaxVersion:               TLSVersionTLS13,
			PreferServerCipherSuites: true,
		},
		HTTPConfig: HTTPConfig{
			HTTP2: true,
		},
	}

	if webConfig.HTTPServerTLSConfig != nil {
		if webConfig.HTTPServerTLSConfig.KeySecret.Name != "" {
			wc.TLSConfig.TLSKeyPath = constants.WhizardCertsMountPath + webConfig.HTTPServerTLSConfig.KeySecret.Key
		}
		if webConfig.HTTPServerTLSConfig.CertSecret.Name != "" {
			wc.TLSConfig.TLSCertPath = constants.WhizardCertsMountPath + webConfig.HTTPServerTLSConfig.CertSecret.Key
		}
		if webConfig.HTTPServerTLSConfig.ClientCASecret.Name != "" {
			wc.TLSConfig.ClientCAs = constants.WhizardCertsMountPath + webConfig.HTTPServerTLSConfig.ClientCASecret.Key
		}
	}

	if webConfig.BasicAuthUsers != nil {
		users, err := createBasicAuthUsers(r.Client, r.Context, namespace, webConfig.BasicAuthUsers)
		if err != nil {
			return nil, err
		}
		if len(users) > 0 {
			wc.Users = users
		}
	}

	body, err := yaml.Marshal(wc)
	if err != nil {
		return nil, err
	}

	return body, nil
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

const (
	TLSVersionTLS12 = "TLS12"
	TLSVersionTLS13 = "TLS13"
)

// We have to redefine webConfig because fields such as passwords are lost when serialized.
// https://pkg.go.dev/github.com/prometheus/exporter-toolkit@v0.8.2/web#Config
type WebConfig struct {
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
