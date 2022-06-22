package receive_ingestor

import (
	"fmt"
	"path/filepath"

	"github.com/prometheus/common/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

func (r *ReceiveIngestor) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name())}

	sts.Spec = appsv1.StatefulSetSpec{
		Replicas: r.ingestor.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		ServiceName: r.name(resources.ServiceNameSuffixOperated),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.ingestor.Spec.NodeSelector,
				Tolerations:  r.ingestor.Spec.Tolerations,
				Affinity:     r.ingestor.Spec.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "receive",
		Image:     r.ingestor.Spec.Image,
		Args:      []string{"receive"},
		Resources: r.ingestor.Spec.Resources,
		Ports: []corev1.ContainerPort{
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          resources.ThanosGRPCPortName,
				ContainerPort: resources.ThanosGRPCPort,
			},
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          resources.ThanosHTTPPortName,
				ContainerPort: resources.ThanosHTTPPort,
			},
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          resources.ThanosRemoteWritePortName,
				ContainerPort: resources.ThanosRemoteWritePort,
			},
		},
		LivenessProbe:  resources.ThanosDefaultLivenessProbe(),
		ReadinessProbe: resources.ThanosDefaultReadinessProbe(),
		Env: []corev1.EnvVar{
			{
				Name: "POD_NAME",
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
	if v := r.ingestor.Spec.DataVolume; v != nil {
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

	var osConfig *corev1.SecretKeySelector
	if lts := r.ingestor.Spec.LongTermStore; lts != nil {
		var store monitoringv1alpha1.Store
		if err := r.Client.Get(r.Context, types.NamespacedName{
			Namespace: lts.Namespace,
			Name:      lts.Name,
		}, &store); err != nil {
			if !apierrors.IsNotFound(err) {
				r.Log.Error(err, "")
			}
		} else {
			osConfig = store.Spec.ObjectStorageConfig
		}
	}

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

	if r.ingestor.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.ingestor.Spec.LogLevel)
	}
	if r.ingestor.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.ingestor.Spec.LogFormat)
	}
	container.Args = append(container.Args, `--label=thanos_receive_replica="$(POD_NAME)"`)
	container.Args = append(container.Args, fmt.Sprintf("--tsdb.path=%s", storageDir))
	container.Args = append(container.Args, fmt.Sprintf("--receive.local-endpoint=$(POD_NAME).%s:%d", r.name(resources.ServiceNameSuffixOperated), resources.ThanosGRPCPort))
	if r.ingestor.Spec.LocalTsdbRetention != "" {
		container.Args = append(container.Args, "--tsdb.retention="+r.ingestor.Spec.LocalTsdbRetention)
	}
	if osConfig != nil && osConfig.Name != "" {
		container.Args = append(container.Args, "--objstore.config-file="+filepath.Join(secretsDir, osConfig.Name, osConfig.Key))
	} else {
		// set tsdb.max-block-duration by localTsdbRetention to enable block compact when using only local storage
		maxBlockDuration, err := model.ParseDuration("31d")
		if err != nil {
			return nil, resources.OperationCreateOrUpdate, err
		}
		retention := r.ingestor.Spec.LocalTsdbRetention
		if retention == "" {
			retention = "15d"
		}
		retentionDuration, err := model.ParseDuration(retention)
		if err != nil {
			return nil, resources.OperationCreateOrUpdate, err
		}
		if retentionDuration != 0 && retentionDuration/10 < maxBlockDuration {
			maxBlockDuration = retentionDuration / 10
		}

		container.Args = append(container.Args, "--tsdb.max-block-duration="+maxBlockDuration.String())
	}

	namespacedName := monitoringv1alpha1.ServiceNamespacedName(r.ingestor)

	if namespacedName != nil {
		var service monitoringv1alpha1.Service
		if err := r.Client.Get(r.Context, *namespacedName, &service); err != nil {
			if !apierrors.IsNotFound(err) {
				return nil, resources.OperationCreateOrUpdate, err
			}
		} else {
			if service.Spec.TenantHeader != "" {
				container.Args = append(container.Args, "--receive.tenant-header="+service.Spec.TenantHeader)
			}
			if service.Spec.TenantLabelName != "" {
				container.Args = append(container.Args, "--receive.tenant-label-name="+service.Spec.TenantLabelName)
			}
			if service.Spec.DefaultTenantId != "" {
				container.Args = append(container.Args, "--receive.default-tenant-id="+service.Spec.DefaultTenantId)
			}
		}
	}

	for name, value := range r.ingestor.Spec.Flags {
		if name == "receive.hashrings" || name == "receive.hashrings-file" {
			// ignoring these flags to make receiver run with ingestor mode
			// refer to https://github.com/thanos-io/thanos/blob/release-0.26/cmd/thanos/receive.go#L816
			continue
		}
		container.Args = append(container.Args, fmt.Sprintf("--%s=%s", name, value))
	}

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container)

	return sts, resources.OperationCreateOrUpdate, nil
}
