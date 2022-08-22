package resources

import (
	"strings"

	"github.com/kubesphere/whizard/pkg/constants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func QualifiedName(appName, instanceName string, suffix ...string) string {
	name := appName + "-" + instanceName
	if len(suffix) > 0 {
		name += "-" + strings.Join(suffix, "-")
	}
	return name
}

func DefaultLivenessProbe() *corev1.Probe {
	return &corev1.Probe{
		FailureThreshold: 4,
		PeriodSeconds:    30,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTP",
				Path:   "/-/healthy",
				Port:   intstr.FromString(constants.HTTPPortName),
			},
		},
	}
}

func DefaultReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		FailureThreshold: 20,
		PeriodSeconds:    5,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTP",
				Path:   "/-/ready",
				Port:   intstr.FromString(constants.HTTPPortName),
			},
		},
	}
}
