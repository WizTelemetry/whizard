package storegateway

import (
	"fmt"
	"path/filepath"

	"github.com/kubesphere/paodin-monitoring/pkg/resources"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *StoreGateway) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name())}

	if r.store == nil || r.Service.Spec.Thanos.ObjectStorageConfig == nil {
		return sts, resources.OperationDelete, nil
	}

	sts.Spec = appsv1.StatefulSetSpec{
		Replicas: r.store.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.store.NodeSelector,
				Tolerations:  r.store.Tolerations,
				Affinity:     r.store.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "store",
		Image:     r.store.Image,
		Args:      []string{"store"},
		Resources: r.store.Resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "grpc",
				ContainerPort: 10901,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          "http",
				ContainerPort: 10902,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		LivenessProbe: &corev1.Probe{
			FailureThreshold: 4,
			PeriodSeconds:    30,
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Scheme: "HTTP",
					Path:   "/-/healthy",
					Port:   intstr.FromString("http"),
				},
			},
		},
		ReadinessProbe: &corev1.Probe{
			FailureThreshold: 20,
			PeriodSeconds:    5,
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Scheme: "HTTP",
					Path:   "/-/ready",
					Port:   intstr.FromString("http"),
				},
			},
		},
	}

	var tsdbVolume = &corev1.Volume{
		Name: "tsdb",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	if v := r.store.DataVolume; v != nil {
		if pvc := v.PersistentVolumeClaim; pvc != nil {
			if pvc.Name == "" {
				pvc.Name = sts.Name + "-tsdb"
			}
			if pvc.Spec.AccessModes == nil {
				pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
			}
			sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, *pvc)
			container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
				Name:      pvc.Name,
				MountPath: storageDir,
			})
			tsdbVolume = nil
		} else if v.EmptyDir != nil {
			tsdbVolume.EmptyDir = v.EmptyDir
		}
	}
	if tsdbVolume != nil {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, *tsdbVolume)
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      tsdbVolume.Name,
			MountPath: storageDir,
		})
	}

	osConfig := r.Service.Spec.Thanos.ObjectStorageConfig
	osVol := corev1.Volume{
		Name: "secret-" + osConfig.Name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: osConfig.Name,
				Items: []corev1.KeyToPath{{
					Key:  osConfig.Key,
					Path: osConfig.Key,
				}},
			},
		},
	}
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, osVol)
	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      osVol.Name,
		MountPath: filepath.Join(secretsDir, osConfig.Name),
	})

	container.Args = append(container.Args, fmt.Sprintf(`--data-dir="%s"`, storageDir))
	if r.store.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.store.LogLevel)
	}
	if r.store.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.store.LogFormat)
	}
	if r.store.MinTime != "" {
		container.Args = append(container.Args, "--min-time="+r.store.MinTime)
	}
	if r.store.MaxTime != "" {
		container.Args = append(container.Args, "--max-time="+r.store.MaxTime)
	}
	container.Args = append(container.Args, "--objstore.config-file="+filepath.Join(secretsDir, osConfig.Name, osConfig.Key))

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container)

	return sts, resources.OperationCreateOrUpdate, nil
}
