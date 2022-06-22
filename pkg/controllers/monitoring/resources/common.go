package resources

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	ThanosGRPCPort        = 10901
	ThanosHTTPPort        = 10902
	ThanosRemoteWritePort = 19291

	ThanosGRPCPortName        = "grpc"
	ThanosHTTPPortName        = "http"
	ThanosRemoteWritePortName = "remote-write"

	ReplicaLabelNamePrometheus    = "prometheus_replica"
	ReplicaLabelNameThanosReceive = "thanos_receive_replica"
	ReplicaLabelNameThanosRuler   = "thanos_ruler_replica"

	AppNameGateway               = "paodin-monitoring-gateway"
	AppNameThanosQuery           = "thanos-query"
	AppNameThanosQueryFrontend   = "thanos-query-frontend"
	AppNameThanosReceiveRouter   = "thanos-receive-router"
	AppNameThanosReceiveIngestor = "thanos-receive-ingestor"
	AppNameThanosRuler           = "thanos-ruler"
	AppNameThanosStoreGateway    = "thanos-store-gateway"
	AppNameThanosCompact         = "thanos-compact"

	ServiceNameSuffixOperated = "operated"

	LabelNameAppComponent = "app.kubernetes.io/component"
	LabelNameAppName      = "app.kubernetes.io/name"
	LableNameAppInstance  = "app.kubernetes.io/instance"
	LabelNameAppManagedBy = "app.kubernetes.io/managed-by"
	LabelNameAppPartOf    = "app.kubernetes.io/part-of"
)

func QualifiedName(appName, instanceName string, suffix ...string) string {
	name := appName + "-" + instanceName
	if len(suffix) > 0 {
		name += "-" + strings.Join(suffix, "-")
	}
	return name
}

func ThanosDefaultLivenessProbe() *corev1.Probe {
	return &corev1.Probe{
		FailureThreshold: 4,
		PeriodSeconds:    30,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTP",
				Path:   "/-/healthy",
				Port:   intstr.FromString(ThanosHTTPPortName),
			},
		},
	}
}

func ThanosDefaultReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		FailureThreshold: 20,
		PeriodSeconds:    5,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTP",
				Path:   "/-/ready",
				Port:   intstr.FromString(ThanosHTTPPortName),
			},
		},
	}
}
