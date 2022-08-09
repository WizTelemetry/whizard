package store

import (
	"fmt"

	storecache "github.com/thanos-io/thanos/pkg/store/cache"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/storage"
)

func (r *Store) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name())}

	if r.store.Spec.Storage == nil {
		return sts, resources.OperationDelete, nil
	}

	sts.Spec = appsv1.StatefulSetSpec{
		Replicas: r.store.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.store.Spec.NodeSelector,
				Tolerations:  r.store.Spec.Tolerations,
				Affinity:     r.store.Spec.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "store",
		Image:     r.store.Spec.Image,
		Args:      []string{"store"},
		Resources: r.store.Spec.Resources,
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
	if v := r.store.Spec.DataVolume; v != nil {
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

	var objstoreConfig string
	storageInstance := &v1alpha1.Storage{}
	err := r.Client.Get(r.Context, types.NamespacedName{Namespace: r.store.Spec.Storage.Namespace, Name: r.store.Spec.Storage.Name}, storageInstance)
	if err != nil {
		return nil, "", err
	}
	objstoreConfig, err = storage.New(r.BaseReconciler, storageInstance).String()
	if err != nil {
		return nil, "", err
	}

	// index cache config
	if cacheConfig := r.store.Spec.IndexCacheConfig; cacheConfig != nil {
		var c *storecache.IndexCacheConfig

		switch cacheConfig.Type {
		case v1alpha1.INMEMORY:
			if inMemory := cacheConfig.InMemoryIndexCacheConfig; inMemory != nil {
				c = &storecache.IndexCacheConfig{
					Type:   storecache.INMEMORY,
					Config: inMemory,
				}
			}
		case v1alpha1.MEMCACHED:
			// TODO
			fallthrough
		case v1alpha1.REDIS:
			// TODO
			fallthrough
		default:
			return nil, resources.OperationCreateOrUpdate, fmt.Errorf("unsupported cache type: %s", cacheConfig.Type)
		}
		if c != nil {
			content, err := yaml.Marshal(c)
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			container.Args = append(container.Args, "--index-cache.config="+string(content))
		}
	}

	container.Args = append(container.Args, fmt.Sprintf("--data-dir=%s", storageDir))
	if r.store.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.store.Spec.LogLevel)
	}
	if r.store.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.store.Spec.LogFormat)
	}
	if r.store.Spec.MinTime != "" {
		container.Args = append(container.Args, "--min-time="+r.store.Spec.MinTime)
	}
	if r.store.Spec.MaxTime != "" {
		container.Args = append(container.Args, "--max-time="+r.store.Spec.MaxTime)
	}
	container.Args = append(container.Args, "--objstore.config="+objstoreConfig)

	for name, value := range r.store.Spec.Flags {
		container.Args = append(container.Args, fmt.Sprintf("--%s=%s", name, value))
	}

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container)

	return sts, resources.OperationCreateOrUpdate, nil
}
