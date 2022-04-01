package receive

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

func (r *receiveIngestor) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name())}

	if r.del {
		return sts, resources.OperationDelete, nil
	}

	sts.Spec = appsv1.StatefulSetSpec{
		Replicas: r.Ingestor.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		ServiceName: r.name("operated"),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.Ingestor.NodeSelector,
				Tolerations:  r.Ingestor.Tolerations,
				Affinity:     r.Ingestor.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "receive",
		Image:     r.Ingestor.Image,
		Args:      []string{"receive"},
		Resources: r.Ingestor.Resources,
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
			{
				Name:          "remote-write",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: 19291,
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
		Env: []corev1.EnvVar{
			{
				Name: "NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
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
	if v := r.Ingestor.DataVolume; v != nil {
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
	if osConfig != nil && osConfig.Name != "" {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
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
		})
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      "secret-" + osConfig.Name,
			MountPath: filepath.Join(secretsDir, osConfig.Name),
		})
	}

	if r.Ingestor.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.Ingestor.LogLevel)
	}
	if r.Ingestor.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.Ingestor.LogFormat)
	}
	container.Args = append(container.Args, `--label=receive_replica="$(NAME)"`)
	container.Args = append(container.Args, fmt.Sprintf(`--tsdb.path="%s"`, storageDir))
	container.Args = append(container.Args, fmt.Sprintf("--receive.local-endpoint=$(NAME).%s:%d", r.name("operated"), 10901))
	if r.Ingestor.LocalTsdbRetention != "" {
		container.Args = append(container.Args, "--tsdb.retention="+r.Ingestor.LocalTsdbRetention)
	}
	if osConfig != nil && osConfig.Name != "" {
		container.Args = append(container.Args, "--objstore.config-file="+filepath.Join(secretsDir, osConfig.Name, osConfig.Key))
	} else {
		// TODO enable block compact when using only local storage
	}

	if r.Service.Spec.TenantHeader != "" {
		container.Args = append(container.Args, "--receive.tenant-header="+r.Service.Spec.TenantHeader)
	}
	if r.Service.Spec.TenantLabelName != "" {
		container.Args = append(container.Args, "--receive.tenant-label-name="+r.Service.Spec.TenantLabelName)
	}
	if r.Service.Spec.DefaultTenantId != "" {
		container.Args = append(container.Args, "--receive.default-tenant-id="+r.Service.Spec.DefaultTenantId)
	}

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container)

	return sts, resources.OperationCreateOrUpdate, nil
}
