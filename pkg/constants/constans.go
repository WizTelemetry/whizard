package constants

import (
	"github.com/prometheus/common/version"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	DefaultTenantHeader    = "WHIZARD-TENANT"
	DefaultTenantId        = "default-tenant"
	DefaultTenantLabelName = "tenant_id"

	ServiceLabelKey = "monitoring.whizard.io/service"
	StorageLabelKey = "monitoring.whizard.io/storage"
	TenantLabelKey  = "monitoring.whizard.io/tenant"

	ExclusiveLabelKey  = "monitoring.whizard.io/exclusive"
	SoftTenantLabelKey = "monitoring.whizard.io/soft-tenant"

	FinalizerIngester  = "finalizers.monitoring.whizard.io/ingester"
	FinalizerCompactor = "finalizers.monitoring.whizard.io/compactor"
	FinalizerDeletePVC = "finalizers.monitoring.whizard.io/deletePVC"

	DefaultStorage = "default"
	LocalStorage   = "local"

	ReceiveReplicaLabelName = "receive_replica"
	RulerReplicaLabelName   = "ruler_replica"

	AppNameGateway       = "gateway"
	AppNameQuery         = "query"
	AppNameQueryFrontend = "query-frontend"
	AppNameRouter        = "router"
	AppNameIngester      = "ingester"
	AppNameRuler         = "ruler"
	AppNameStore         = "store"
	AppNameCompactor     = "compactor"
	AppNameStorage       = "storage"
	AppNameBlockManager  = "block-manager"

	ServiceNameSuffix = "operated"

	LabelNameAppName      = "app.kubernetes.io/name"
	LabelNameAppManagedBy = "app.kubernetes.io/managed-by"
	LabelNameAppPartOf    = "app.kubernetes.io/part-of"

	LabelNameIngesterState        = "monitoring.whizard.io/ingester-state"
	LabelNameIngesterDeletingTime = "monitoring.whizard.io/ingester-deleting-time"

	LabelNameStorePartition = "monitoring.whizard.io/store-partition"

	LabelNameRulerShardSn = "monitoring.whizard.io/ruler-shard-sn"

	LabelNameStorageHash = "monitoring.whizard.io/storage-hash"
	LabelNameTenantHash  = "monitoring.whizard.io/tenant-hash"
	LabelNameConfigHash  = "monitoring.whizard.io/config-hash"

	TenantHash  = "TENANT_HASH"
	StorageHash = "STORAGE_HASH"

	IngesterStateDeleting = "deleting"
	IngesterStateRunning  = "running"
)

// Mount path of config files in containers.
const (
	ConfigPath     = "/etc/whizard/"
	StorageDir     = "/whizard"
	TSDBVolumeName = "tsdb"

	WhizardConfigMountPath     = "/etc/whizard/config/"
	WhizardWebConfigMountPath  = "/etc/whizard/web_config/"
	WhizardWebConfigFile       = "web-config.yaml"
	WhizardCertsMountPath      = "/etc/whizard/certs/"
	WhizardConfigMapsMountPath = "/etc/whizard/configmaps/"
	WhizardSecretsMountPath    = "/etc/whizard/secrets/"

	EnvoyConfigMountPath    = "/etc/envoy/config/"
	EnvoyCertsMountPath     = "/etc/envoy/certs/"
	EnvoyConfigMapMountPath = "/etc/envoy/configmap/"
	EnvoySecretMountPath    = "/etc/envoy/secret/"

	// common
	WhizardObjStoreConfigFile       = "/etc/whizard/config/objstore.yaml"
	WhizardTracingConfigFile        = "/etc/whizard/config/tracing.yaml"
	WhizardRequestLoggingConfigFile = "/etc/whizard/config/logging-config.yaml"
	// query

	// compactor
	WhizardSelectorRelabelConfigFile = "/etc/whizard/config/selector-relabel-config.yaml"
	// query-frontend
	WhizardLabelsResponseCacheConfigFile            = "/etc/whizard/config/labels-response-cache-config.yaml"
	WhizardQueryFrontendDownstreamTripperConfigFile = "/etc/whizard/config/query-frontend-downstream-tripper-config.yaml"
	WhizardQueryRangeResponseCacheConfigFile        = "/etc/whizard/config/query-range-response-cache-config.yaml"
	// receive
	WhizardReceiveRelabelConfigFile = "/etc/whizard/config/receive-relabel-config.yaml"
	WhizardReceiveHashringsFile     = "/etc/whizard/config/hashrings.yaml"
	// ruler
	WhizardAlertRelabelConfigFile  = "/etc/whizard/config/alert-relabel-config.yaml"
	WhizardAlertmanagersConfigFile = "/etc/whizard/config/alertmanagers-config.yaml"
	WhizardQueryConfigFile         = "/etc/whizard/config/query-config.yaml"
	WhizardRemoteWriteConfigFile   = "/etc/whizard/config/remote-write-config.yaml"
	// sidecar
	WhizardReloaderConfigFile = "/etc/whizard/config/reloader-config.yaml"
	// store
	WhizardIndexCacheConfigFile = "/etc/whizard/config/index-cache.config"

	// tls certs
	WhizardGRPCServerTLSPath = "/etc/whizard/certs/grpc-server-tls"
	WhizardGRPCClientTLSPath = "/etc/whizard/certs/grpc-client-tls"
	WhizardHTTPServerTLSPath = "/etc/whizard/certs/http-server-tls"
	WhizardHTTPClientTLSPath = "/etc/whizard/certs/http-server-tls"
)

const (
	StorageProviderFILESYSTEM string = "FILESYSTEM"
	StorageProviderGCS        string = "GCS"
	StorageProviderS3         string = "S3"
	StorageProviderAZURE      string = "AZURE"
	StorageProviderSWIFT      string = "SWIFT"
	StorageProviderCOS        string = "COS"
	StorageProviderALIYUNOSS  string = "ALIYUNOSS"
	StorageProviderBOS        string = "BOS"
)

const (
	// The version is the same as thanos mod version
	DefaultWhizardBaseImage = "docker.io/thanosio/thanos:v0.38.0"
	// The version is the same as prometheus-operator mod version
	DefaultPrometheusConfigReloaderImage = "quay.io/prometheus-operator/prometheus-config-reloader:v0.81.0"

	DefaultEnvoyImage               = "docker.io/envoyproxy/envoy:v1.20.2"
	DefaultRulerWriteProxyImage     = "docker.io/kubesphere/cortex-tenant:v1.12.5"
	DefaultIngesterTSDBCleanupImage = "bash:5.1.16"
)

var DefaultWhizardMonitoringGatewayImage = "docker.io/kubesphere/whizard-monitoring-gateway:" + version.Version
var DefaultWhizardBlockManagerImage = "docker.io/kubesphere/whizard-monitoring-block-manager:" + version.Version

const (
	GRPCPortName = "grpc"
	GRPCPort     = 10901
	HTTPPortName = "http"
	HTTPPort     = 10902

	// receive
	RemoteWritePortName = "remote-write"
	RemoteWritePort     = 19291
	CapnprotoPortName   = "capnproto"
	CapnprotoPort       = 19391
)

// ConponentProbePreset defines standard probe presets for components.
var (
	ComponentProbePresetHTTPLivenessProbe = corev1.Probe{
		InitialDelaySeconds: 30,
		TimeoutSeconds:      30,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    6,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: corev1.URISchemeHTTP,
				Path:   "/-/healthy",
				Port:   intstr.FromString(HTTPPortName),
			},
		},
	}
	ComponentProbePresetHTTPReadinessProbe = corev1.Probe{
		InitialDelaySeconds: 30,
		TimeoutSeconds:      30,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    6,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: corev1.URISchemeHTTP,
				Path:   "/-/ready",
				Port:   intstr.FromString(HTTPPortName),
			},
		},
	}
	ComponentProbePresetHTTPSLivenessProbe = corev1.Probe{
		InitialDelaySeconds: 30,
		TimeoutSeconds:      30,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    6,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: corev1.URISchemeHTTPS,
				Path:   "/-/healthy",
				Port:   intstr.FromString(HTTPPortName),
			},
		},
	}
	ComponentProbePresetHTTPSReadinessProbe = corev1.Probe{
		InitialDelaySeconds: 30,
		TimeoutSeconds:      30,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    6,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: corev1.URISchemeHTTPS,
				Path:   "/-/ready",
				Port:   intstr.FromString(HTTPPortName),
			},
		},
	}

	// Use a TCP Socket probe when the service requires authentication.
	ComponentProbePresetTCPSocketProbe = corev1.Probe{
		InitialDelaySeconds: 30,
		TimeoutSeconds:      30,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    6,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromString(HTTPPortName),
			},
		},
	}
)

// ComponentResourcePreset defines standard resource presets for different components.
var (
	// ComponentResourcePresetMedium defines medium resource presets for query, query-frontend, ruler, and compactor.
	ComponentResourcePresetMedium = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("50m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("2"),
			corev1.ResourceMemory: resource.MustParse("4Gi"),
		},
	}

	// ComponentResourcePresetLarge defines large resource presets for ingester and store.
	ComponentResourcePresetLarge = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("50m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("4"),
			corev1.ResourceMemory: resource.MustParse("16Gi"),
		},
	}

	// ComponentResourcePresetSmall defines small resource presets for router and sidecar.
	ComponentResourcePresetSmall = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("50m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("512Mi"),
		},
	}
)
