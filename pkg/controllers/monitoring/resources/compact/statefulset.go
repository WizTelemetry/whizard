package compact

import (
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

func (r *Compact) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name())}

	if r.compact == nil || r.Store.Spec.ObjectStorageConfig == nil {
		return sts, resources.OperationDelete, nil
	}

	sts.Spec = appsv1.StatefulSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.compact.NodeSelector,
				Tolerations:  r.compact.Tolerations,
				Affinity:     r.compact.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "compact",
		Image:     r.compact.Image,
		Args:      []string{"compact", "--wait"},
		Resources: r.compact.Resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          resources.ThanosHTTPPortName,
				ContainerPort: resources.ThanosHTTPPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		LivenessProbe:  resources.ThanosDefaultLivenessProbe(),
		ReadinessProbe: resources.ThanosDefaultReadinessProbe(),
	}

	var tsdbVolume = &corev1.Volume{
		Name: "tsdb",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	if v := r.compact.DataVolume; v != nil {
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

	osConfig := r.Store.Spec.ObjectStorageConfig
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

	if r.compact.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.compact.LogLevel)
	}
	if r.compact.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.compact.LogFormat)
	}
	container.Args = append(container.Args, fmt.Sprintf("--data-dir=%s", storageDir))
	if r.compact.DownsamplingDisable != nil {
		container.Args = append(container.Args, fmt.Sprintf("--downsampling.disable=%v", r.compact.DownsamplingDisable))
	}
	container.Args = append(container.Args, "--objstore.config-file="+filepath.Join(secretsDir, osConfig.Name, osConfig.Key))
	if retention := r.compact.Retention; retention != nil {
		if retention.RetentionRaw != "" {
			container.Args = append(container.Args, "--retention.resolution-raw="+retention.RetentionRaw)
		}
		if retention.Retention5m != "" {
			container.Args = append(container.Args, "--retention.resolution-5m="+retention.Retention5m)
		}
		if retention.Retention5m != "" {
			container.Args = append(container.Args, "--retention.resolution-1h="+retention.Retention5m)
		}
	}
	container.Args = append(container.Args, "--deduplication.replica-label=receive_replica")

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container)

	return sts, resources.OperationCreateOrUpdate, nil
}
