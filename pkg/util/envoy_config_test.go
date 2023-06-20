package util

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var staticConfig = `
static_resources:
  listeners:
    - name: self_listener
      address:
        socket_address:
          address: 0.0.0.0
          port_value: 10902
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
                      filename: /etc/whizard/certs/tls.crt
                    private_key:
                      filename: /etc/whizard/certs/tls.key
    - name: proxy_listener
      address:
        socket_address:
          address: 0.0.0.0
          port_value: 10000
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
  clusters:
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
                    port_value: 10904
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
                    address: query-whizard-operated.kubesphere-monitoring-system.svc
                    port_value: 10902
      transport_socket:
        name: envoy.transport_sockets.tls
        typed_config:
          '@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
          common_tls_context:
            validation_context:
              trust_chain_verification: ACCEPT_UNTRUSTED
            alpn_protocols: 'h2,http/1.1'
          sni: query-whizard-operated.kubesphere-monitoring-system.svc
`

func TestTempExec(t *testing.T) {
	variables := map[string]string{
		"ServiceMappingPort":   "10902",
		"ServiceListenPort":    "10904",
		"ServiceTLSCertFile":   "/etc/whizard/certs/tls.crt",
		"ServiceTLSKeyFile":    "/etc/whizard/certs/tls.key",
		"ProxyServiceEnabled":  "true",
		"ProxyLocalListenPort": "10000",
		"ProxyServiceAddress":  "query-whizard-operated.kubesphere-monitoring-system.svc",
		"ProxyServicePort":     "10902",
	}
	var buff strings.Builder

	tmpl := EnvoyStaticConfigTemplate
	if err := tmpl.Execute(&buff, variables); err != nil {
		t.Error(err)
	}
	t.Log(buff.String())

	if diff := cmp.Diff(buff.String(), staticConfig); diff != "" {
		t.Error(diff)
	}

}
