package compactor

import (
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

func (r *Compactor) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name())}

	if r.compactor.Spec.Storage == nil {
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
				NodeSelector: r.compactor.Spec.NodeSelector,
				Tolerations:  r.compactor.Spec.Tolerations,
				Affinity:     r.compactor.Spec.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "compactor",
		Image:     r.compactor.Spec.Image,
		Args:      []string{"compact", "--wait"},
		Resources: r.compactor.Spec.Resources,
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
	if v := r.compactor.Spec.DataVolume; v != nil {
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

	storage := &v1alpha1.Storage{}
	err := r.Client.Get(r.Context, types.NamespacedName{Namespace: r.compactor.Spec.Storage.Namespace, Name: r.compactor.Spec.Storage.Name}, storage)
	if err != nil {
		return nil, "", err
	}

	osVol := corev1.Volume{
		Name: storage.Name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: "secret-" + storage.Name,
				Items: []corev1.KeyToPath{{
					Key:  resources.SecretBucketKey,
					Path: resources.SecretBucketKey,
				}},
			},
		},
	}
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, osVol)
	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      osVol.Name,
		MountPath: filepath.Join(secretsDir, storage.Name),
	})

	if r.compactor.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.compactor.Spec.LogLevel)
	}
	if r.compactor.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.compactor.Spec.LogFormat)
	}
	container.Args = append(container.Args, fmt.Sprintf("--data-dir=%s", storageDir))
	if r.compactor.Spec.DownsamplingDisable != nil {
		container.Args = append(container.Args, fmt.Sprintf("--downsampling.disable=%v", r.compactor.Spec.DownsamplingDisable))
	}
	container.Args = append(container.Args, "--objstore.config-file="+filepath.Join(secretsDir, storage.Name, resources.SecretBucketKey))
	if retention := r.compactor.Spec.Retention; retention != nil {
		if retention.RetentionRaw != "" {
			container.Args = append(container.Args, fmt.Sprintf("--retention.resolution-raw=%s", retention.RetentionRaw))
		}
		if retention.Retention5m != "" {
			container.Args = append(container.Args, fmt.Sprintf("--retention.resolution-5m=%s", retention.Retention5m))
		}
		if retention.Retention5m != "" {
			container.Args = append(container.Args, fmt.Sprintf("--retention.resolution-1h=%s", retention.Retention5m))
		}
	}
	container.Args = append(container.Args, "--deduplication.replica-label=receive_replica")

	for name, value := range r.compactor.Spec.Flags {
		container.Args = append(container.Args, fmt.Sprintf("--%s=%s", name, value))
	}

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container)

	return sts, resources.OperationCreateOrUpdate, nil
}
