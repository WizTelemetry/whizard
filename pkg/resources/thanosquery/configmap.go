package thanosquery

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"time"

	envoy_config_accesslog "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	envoy_config_bootstrap "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_config_cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_extentions_accessloggers "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/file/v3"
	envoy_extensions_filters_network_http "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoy_extensions_transport_sockets_tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	envoy_service_discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	ghodss_yaml "github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/kubesphere/paodin-monitoring/api/v1alpha1"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (q *ThanosQuery) configMaps() ([]*corev1.ConfigMap, error) {
	storeSDCm, err := q.storeSDConfigMap()
	if err != nil {
		return nil, err
	}
	envoyCm, err := q.envoyConfigMap()
	if err != nil {
		return nil, err
	}
	return []*corev1.ConfigMap{
		storeSDCm,
		envoyCm,
	}, nil
}

func (q *ThanosQuery) storeSDConfigMap() (*corev1.ConfigMap, error) {
	var (
		targets, proxyTargets []model.LabelSet
		listenerPort          uint32 = envoyListenerStartPort
	)
	for _, store := range q.Instance.Spec.Stores {
		if store.Address == "" {
			continue
		}
		if storeRequireProxy(store) {
			proxyTargets = append(proxyTargets, model.LabelSet{
				model.AddressLabel: model.LabelValue(fmt.Sprintf("%s:%d", envoyListenerAddress, listenerPort)),
			})
			listenerPort += 1
		} else {
			targets = append(targets, model.LabelSet{
				model.AddressLabel: model.LabelValue(store.Address),
			})
		}
	}
	var groups []targetgroup.Group
	if len(targets) > 0 {
		groups = append(groups, targetgroup.Group{Targets: targets})
	}
	if len(proxyTargets) > 0 {
		groups = append(groups, targetgroup.Group{Targets: proxyTargets})
	}

	out, err := yaml.Marshal(groups)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: q.Instance.Namespace,
			Name:      q.getStoreSDConfigMapName(),
			Labels:    q.labels(),
		},
		Data: map[string]string{
			storeSDFileName: string(out),
		},
	}, err
}

func (q *ThanosQuery) envoyConfigMap() (*corev1.ConfigMap, error) {
	m := jsonpb.Marshaler{OrigName: true}

	bootstrap, err := getEnvoyBootstrap(q.Instance)
	if err != nil {
		return nil, err
	}
	configString, err := m.MarshalToString(bootstrap)
	if err != nil {
		return nil, err
	}
	configBytes, err := ghodss_yaml.JSONToYAML([]byte(configString))
	if err != nil {
		return nil, err
	}

	lds, err := getEnvoyLDS(q.Instance)
	if err != nil {
		return nil, err
	}
	ldsString, err := m.MarshalToString(lds)
	if err != nil {
		return nil, err
	}
	ldsBytes, err := ghodss_yaml.JSONToYAML([]byte(ldsString))
	if err != nil {
		return nil, err
	}

	cds, err := getEnvoyCDS(q.Instance)
	if err != nil {
		return nil, err
	}
	cdsString, err := m.MarshalToString(cds)
	if err != nil {
		return nil, err
	}
	cdsBytes, err := ghodss_yaml.JSONToYAML([]byte(cdsString))
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: q.Instance.Namespace,
			Name:      q.getEnvoyConfigMapName(),
			Labels:    q.labels(),
		},
		Data: map[string]string{
			envoyConfigFileName: string(configBytes),
			envoyLDSFileName:    string(ldsBytes),
			envoyCDSFileName:    string(cdsBytes),
		},
	}, nil
}

func getEnvoyBootstrap(instance v1alpha1.ThanosQuery) (*envoy_config_bootstrap.Bootstrap, error) {
	// TODO secret reload

	var bootstrap = &envoy_config_bootstrap.Bootstrap{}
	adminAccessLog, err := anypb.New(&envoy_extentions_accessloggers.FileAccessLog{Path: "/tmp/admin_access.log"})
	if err != nil {
		return nil, err
	}
	bootstrap.Admin = &envoy_config_bootstrap.Admin{
		AccessLog: []*envoy_config_accesslog.AccessLog{{
			Name: "query-stores-proxy-access-log",
			ConfigType: &envoy_config_accesslog.AccessLog_TypedConfig{
				TypedConfig: adminAccessLog,
			},
		}},
	}
	bootstrap.Node = &envoy_config_core.Node{
		Cluster: "query-stores-proxy",
		Id:      instance.Namespace + "/" + instance.Name,
	}
	bootstrap.DynamicResources = &envoy_config_bootstrap.Bootstrap_DynamicResources{
		LdsConfig: &envoy_config_core.ConfigSource{
			ConfigSourceSpecifier: &envoy_config_core.ConfigSource_Path{
				Path: filepath.Join(envoyConfigDir, envoyLDSFileName),
			},
		},
		CdsConfig: &envoy_config_core.ConfigSource{
			ConfigSourceSpecifier: &envoy_config_core.ConfigSource_Path{
				Path: filepath.Join(envoyConfigDir, envoyCDSFileName),
			},
		},
	}
	return bootstrap, nil
}

func getEnvoyLDS(instance v1alpha1.ThanosQuery) (*envoy_service_discovery.DiscoveryResponse, error) {
	var (
		listeners    []*envoy_config_listener.Listener
		listenerPort uint32 = envoyListenerStartPort
	)

	for _, store := range instance.Spec.Stores {
		if !storeRequireProxy(store) {
			continue
		}
		host, _, err := net.SplitHostPort(store.Address)
		if err != nil {
			return nil, err
		}

		var (
			listenerName = store.Address
			clusterName  = store.Address
		)
		var filter = &envoy_extensions_filters_network_http.HttpConnectionManager{
			CodecType:  envoy_extensions_filters_network_http.HttpConnectionManager_AUTO,
			StatPrefix: "store_apis",
			HttpFilters: []*envoy_extensions_filters_network_http.HttpFilter{{
				Name: "envoy.filters.http.router",
			}},
			RouteSpecifier: &envoy_extensions_filters_network_http.HttpConnectionManager_RouteConfig{
				RouteConfig: &envoy_config_route.RouteConfiguration{
					Name: "local_route",
					VirtualHosts: []*envoy_config_route.VirtualHost{{
						Name:    "local_service",
						Domains: []string{"*"},
						Routes: []*envoy_config_route.Route{{
							Match: &envoy_config_route.RouteMatch{
								PathSpecifier: &envoy_config_route.RouteMatch_Prefix{Prefix: "/"},
							},
							Action: &envoy_config_route.Route_Route{
								Route: &envoy_config_route.RouteAction{
									ClusterSpecifier: &envoy_config_route.RouteAction_Cluster{
										Cluster: clusterName,
									},
									HostRewriteSpecifier: &envoy_config_route.RouteAction_HostRewriteLiteral{
										HostRewriteLiteral: host,
									},
								},
							},
						}},
					}},
				},
			},
		}

		filterAny, err := anypb.New(filter)
		if err != nil {
			return nil, err
		}

		var listener = &envoy_config_listener.Listener{
			Name: listenerName,
			Address: &envoy_config_core.Address{
				Address: &envoy_config_core.Address_SocketAddress{
					SocketAddress: &envoy_config_core.SocketAddress{
						Address: envoyListenerAddress,
						PortSpecifier: &envoy_config_core.SocketAddress_PortValue{
							PortValue: listenerPort,
						},
					},
				},
			},
			FilterChains: []*envoy_config_listener.FilterChain{{
				Filters: []*envoy_config_listener.Filter{{
					Name: "envoy.http_connection_manager",
					ConfigType: &envoy_config_listener.Filter_TypedConfig{
						TypedConfig: filterAny,
					},
				}},
			}},
		}

		listeners = append(listeners, listener)

		listenerPort += 1
	}

	var resources []*anypb.Any
	for _, listener := range listeners {
		listenerAny, err := anypb.New(listener)
		if err != nil {
			return nil, err
		}
		resources = append(resources, listenerAny)
	}

	return &envoy_service_discovery.DiscoveryResponse{
		Resources: resources,
	}, nil
}

func getEnvoyCDS(instance v1alpha1.ThanosQuery) (*envoy_service_discovery.DiscoveryResponse, error) {
	var clusters []*envoy_config_cluster.Cluster

	for _, store := range instance.Spec.Stores {
		if !storeRequireProxy(store) {
			continue
		}

		host, portString, err := net.SplitHostPort(store.Address)
		if err != nil {
			return nil, err
		}
		port, err := strconv.ParseUint(portString, 10, 32)
		if err != nil {
			return nil, err
		}

		var (
			clusterName = store.Address
		)

		var tlsCtx = &envoy_extensions_transport_sockets_tls.UpstreamTlsContext{
			CommonTlsContext: &envoy_extensions_transport_sockets_tls.CommonTlsContext{
				ValidationContextType: &envoy_extensions_transport_sockets_tls.CommonTlsContext_ValidationContext{
					ValidationContext: &envoy_extensions_transport_sockets_tls.CertificateValidationContext{
						TrustedCa: &envoy_config_core.DataSource{
							Specifier: &envoy_config_core.DataSource_Filename{
								// TODO check ca and client cert
								Filename: filepath.Join(envoySecretsDir, store.SecretName, "ca.crt"),
							},
						},
					},
				},
				AlpnProtocols: []string{"h2", "http/1.1"},
			},
			Sni: host,
		}

		tlsCtxAny, err := anypb.New(tlsCtx)
		if err != nil {
			return nil, err
		}

		var cluster = &envoy_config_cluster.Cluster{
			Name:           clusterName,
			ConnectTimeout: durationpb.New(time.Second * 30),
			ClusterDiscoveryType: &envoy_config_cluster.Cluster_Type{
				Type: envoy_config_cluster.Cluster_STRICT_DNS,
			},
			Http2ProtocolOptions: &envoy_config_core.Http2ProtocolOptions{},
			DnsLookupFamily:      envoy_config_cluster.Cluster_V4_ONLY,
			LbPolicy:             envoy_config_cluster.Cluster_ROUND_ROBIN,
			TransportSocket: &envoy_config_core.TransportSocket{
				Name: "envoy.transport_sockets.tls",
				ConfigType: &envoy_config_core.TransportSocket_TypedConfig{
					TypedConfig: tlsCtxAny,
				},
			},
			LoadAssignment: &envoy_config_endpoint.ClusterLoadAssignment{
				ClusterName: clusterName,
				Endpoints: []*envoy_config_endpoint.LocalityLbEndpoints{{
					LbEndpoints: []*envoy_config_endpoint.LbEndpoint{{
						HostIdentifier: &envoy_config_endpoint.LbEndpoint_Endpoint{
							Endpoint: &envoy_config_endpoint.Endpoint{
								Address: &envoy_config_core.Address{
									Address: &envoy_config_core.Address_SocketAddress{
										SocketAddress: &envoy_config_core.SocketAddress{
											Address: host,
											PortSpecifier: &envoy_config_core.SocketAddress_PortValue{
												PortValue: uint32(port),
											},
										},
									},
								},
							},
						},
					}},
				}},
			},
		}

		clusters = append(clusters, cluster)
	}

	var resources []*anypb.Any
	for _, cluster := range clusters {
		clusterAny, err := anypb.New(cluster)
		if err != nil {
			return nil, err
		}
		resources = append(resources, clusterAny)
	}

	return &envoy_service_discovery.DiscoveryResponse{
		Resources: resources,
	}, nil
}
