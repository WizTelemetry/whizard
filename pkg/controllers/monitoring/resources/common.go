package resources

import (
	"strings"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	appsv1 "k8s.io/api/apps/v1"
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

func AddTSDBVolume(sts *appsv1.StatefulSet, container *corev1.Container, dataVolume *v1alpha1.KubernetesVolume) {
	if dataVolume == nil ||
		(dataVolume.PersistentVolumeClaim == nil && dataVolume.EmptyDir == nil) {
		return
	}

	if dataVolume.PersistentVolumeClaim != nil {
		pvc := *dataVolume.PersistentVolumeClaim
		if pvc.Name == "" {
			pvc.Name = constants.TSDBVolumeName
		}
		if pvc.Spec.AccessModes == nil {
			pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		}

		sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{pvc}
	} else if dataVolume.EmptyDir != nil {
		sts.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: constants.TSDBVolumeName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: dataVolume.EmptyDir,
				},
			},
		}
	}

	container.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      constants.TSDBVolumeName,
			MountPath: constants.StorageDir,
		},
	}
}
