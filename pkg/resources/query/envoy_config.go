package query

import (
	"bytes"
	"path/filepath"
	"text/template"
)

const (
	envoyConfigDir  = "/etc/envoy/config"
	envoyConfigFile = "envoy.yaml"
	envoyLdsFile    = "lds.yaml"
	envoyCdsFile    = "cds.yaml"
	envoySecretsDir = "/etc/envoy/secrets"
)

var envoyConfigTemplate = template.Must(template.New("envoy-config").Parse(`
node:
  cluster: {{ .NodeCluster }}
  id: {{.NodeId}}

dynamic_resources:
  cds_config:
    path: {{ .CdsPath }}
  lds_config:
    path: {{ .LdsPath }}
`))

var envoyLdsTemplate = template.Must(template.New("envoy-lds").Parse(`
resources:
{{- range .ProxyStores }}
- "@type": type.googleapis.com/envoy.config.listener.v3.Listener
  name: {{ printf "%s:%d" .ListenHost .ListenPort }}
  address:
    socket_address:
      address: {{ .ListenHost }}
      port_value: {{ .ListenPort }}
  filter_chains:
  - filters:
    - name: envoy.http_connection_manager
      typed_config:
        "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
        stat_prefix: store_apis
        http_filters:
        - name: envoy.filters.http.router
        route_config:
          name: local_route
          virtual_hosts:
          - name: local_service
            domains:
            - "*"
            routes:
            - match:
                prefix: "/"
              route:
                host_rewrite_literal: {{ .TargetHost }}
                cluster: {{ printf "%s:%d" .TargetHost .TargetPort }}
{{ end }}
`))

var envoyCdsTemplate = template.Must(template.New("envoy-cds").Parse(`
resources:
{{- range .ProxyStores }}
- "@type": type.googleapis.com/envoy.config.cluster.v3.Cluster
  name: {{ printf "%s:%d" .TargetHost .TargetPort }}
  type: STRICT_DNS
  typed_extension_protocol_options:
    envoy.extensions.upstreams.http.v3.HttpProtocolOptions:
      "@type": type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions
      explicit_http_config:
        http2_protocol_options: {}
  load_assignment:
    cluster_name: {{ printf "%s:%d" .TargetHost .TargetPort }}
    endpoints:
    - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              address: {{ .TargetHost }}
              port_value: {{ .TargetPort }}
  transport_socket:
    name: envoy.transport_sockets.tls
    typed_config:
      "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
      common_tls_context: 
        validation_context:
          trusted_ca:
            filename: {{ .TlsCaFile}}
        alpn_protocols:
        - h2
        - http/1.1
      sni: {{ .TargetHost }}
{{ end }}
`))

func envoyConfigFiles(identifier string, stores Stores) (map[string]string, error) {
	var (
		files = map[string]string{}
		err   error
		b     = new(bytes.Buffer)
	)

	var config = &struct {
		NodeId      string
		NodeCluster string
		CdsPath     string
		LdsPath     string
		ProxyStores []ProxyStore
	}{
		NodeId:      identifier,
		NodeCluster: "query-stores-proxy",
		CdsPath:     filepath.Join(envoyConfigDir, envoyCdsFile),
		LdsPath:     filepath.Join(envoyConfigDir, envoyLdsFile),
		ProxyStores: stores.ProxyStores,
	}

	if err = envoyConfigTemplate.Execute(b, config); err != nil {
		return nil, err
	}
	files[envoyConfigFile] = b.String()
	b.Reset()
	if err = envoyLdsTemplate.Execute(b, config); err != nil {
		return nil, err
	}
	files[envoyLdsFile] = b.String()
	b.Reset()
	if err = envoyCdsTemplate.Execute(b, config); err != nil {
		return nil, err
	}
	files[envoyCdsFile] = b.String()
	b.Reset()

	return files, nil
}
