package util

import (
	"text/template"

	"github.com/lithammer/dedent"
)

var EnvoyStaticConfigTemplate = template.Must(template.New("envoy.yaml").Parse(dedent.Dedent(`
static_resources:
  listeners:
{{ if .LocalServiceEnabled }}
    - name: self_listener
      address:
        socket_address:
          address: 0.0.0.0
          port_value: {{ .ServiceMappingPort }}
      filter_chains:
        - filters:
            - name: envoy.filters.network.http_connection_manager
              typed_config:
                '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                scheme_header_transformation:
                  scheme_to_overwrite: https
                stat_prefix: ingress_http
                codec_type: auto
                route_config:
                  virtual_hosts:
                    - name: default
                      domains:
                        - '*'
                      routes:
                        - match:
                            prefix: /
                          route:
                            cluster: local_serivce
                http_filters:
                  - name: envoy.filters.http.router
          transport_socket:
            name: envoy.transport_sockets.tls
            typed_config:
              '@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext
              common_tls_context:
                tls_certificates:
                  - certificate_chain:
                      filename: {{ .ServiceTLSCertFile }}
                    private_key:
                      filename: {{ .ServiceTLSKeyFile }}
{{ end }}
{{ if .ProxyServiceEnabled }}
    - name: proxy_listener
      address:
        socket_address:
          address: 0.0.0.0
          port_value: {{ .ProxyLocalListenPort }}
      filter_chains:
        - filters:
            - name: envoy.filters.network.http_connection_manager
              typed_config:
                '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                scheme_header_transformation:
                  scheme_to_overwrite: https
                stat_prefix: ingress_http
                codec_type: auto
                route_config:
                  virtual_hosts:
                    - name: default
                      domains:
                        - '*'
                      routes:
                        - match:
                            prefix: /
                          route:
                            cluster: upstrean_serivce
                http_filters:
                  - name: envoy.filters.http.router
{{ end }}
  clusters:
{{ if .LocalServiceEnabled }}  
  - name: local_serivce
    type: STATIC
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: local_serivce
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: 127.0.0.1
                    port_value: {{ .ServiceListenPort }}
{{ end }} 
{{ if .ProxyServiceEnabled }}
  - name: upstrean_serivce
    type: STRICT_DNS
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: upstrean_serivce
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: {{ .ProxyServiceAddress }}
                    port_value: {{ .ProxyServicePort }}
    transport_socket:
      name: envoy.transport_sockets.tls
      typed_config:
        '@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
        common_tls_context:
          validation_context:
            trust_chain_verification: ACCEPT_UNTRUSTED
          alpn_protocols: 'h2,http/1.1'
        sni: {{ .ProxyServiceAddress }}
{{ end }} 
`)))
