package store

import (
	"fmt"

	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	mainContainerName = "store"
)

func (r *Store) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.store.Name)}
	if err := r.Client.Get(r.Context, client.ObjectKeyFromObject(sts), sts); err != nil {
		if !util.IsNotFound(err) {
			return nil, "", err
		}
	}

	sts.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: r.labels(),
	}

	if sts.Spec.Replicas == nil || *sts.Spec.Replicas == 0 {
		sts.Spec.Replicas = r.store.Spec.Replicas
	}

	sts.Spec.Template.Labels = r.labels()

	sts.Spec.Template.Spec.Affinity = r.store.Spec.Affinity
	sts.Spec.Template.Spec.NodeSelector = r.store.Spec.NodeSelector
	sts.Spec.Template.Spec.Tolerations = r.store.Spec.Tolerations
	sts.Spec.Template.Spec.Volumes = []corev1.Volume{}
	sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}

	var container *corev1.Container
	for i := 0; i < len(sts.Spec.Template.Spec.Containers); i++ {
		if sts.Spec.Template.Spec.Containers[i].Name == mainContainerName {
			container = &sts.Spec.Template.Spec.Containers[i]
		}
	}

	needToAppend := false
	if container == nil {
		container = &corev1.Container{
			Name:  mainContainerName,
			Image: r.store.Spec.Image,
			Ports: []corev1.ContainerPort{
				{
					Protocol:      corev1.ProtocolTCP,
					Name:          constants.GRPCPortName,
					ContainerPort: constants.GRPCPort,
				},
				{
					Protocol:      corev1.ProtocolTCP,
					Name:          constants.HTTPPortName,
					ContainerPort: constants.HTTPPort,
				},
			},
		}

		needToAppend = true
	}

	container.VolumeMounts = []corev1.VolumeMount{}
	resources.AddTSDBVolume(sts, container, r.store.Spec.DataVolume)

	volumes, volumeMounts, err := resources.VolumesAndVolumeMountsForStorage(r.Context, r.Client, r.store.Labels[constants.StorageLabelKey])
	if err != nil {
		return nil, "", err
	}
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, volumes...)
	container.VolumeMounts = append(container.VolumeMounts, volumeMounts...)

	if container.LivenessProbe == nil {
		container.LivenessProbe = resources.DefaultLivenessProbe()
	}

	if container.ReadinessProbe == nil {
		container.ReadinessProbe = resources.DefaultReadinessProbe()
	}

	container.Resources = r.store.Spec.Resources

	env := corev1.EnvVar{
		Name:  constants.StorageHash,
		Value: r.store.Annotations[constants.LabelNameStorageHash],
	}
	replaced := util.ReplaceInSlice(container.Env, func(v interface{}) bool {
		return v.(corev1.EnvVar).Name == constants.StorageHash
	}, env)
	if !replaced {
		container.Env = append(container.Env, env)
	}

	env = corev1.EnvVar{
		Name:  constants.TenantHash,
		Value: r.store.Annotations[constants.LabelNameTenantHash],
	}
	replaced = util.ReplaceInSlice(container.Env, func(v interface{}) bool {
		return v.(corev1.EnvVar).Name == constants.TenantHash
	}, env)
	if !replaced {
		container.Env = append(container.Env, env)
	}

	if err := r.megerArgs(container); err != nil {
		return nil, "", err
	}

	if needToAppend {
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *container)
	}

	return sts, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.store, sts, r.Scheme)
}

func (r *Store) megerArgs(container *corev1.Container) error {
	defaultArgs := []string{"store", fmt.Sprintf("--data-dir=%s", constants.StorageDir)}

	storageConfig, err := resources.GetStorageConfig(r.Context, r.Client, r.store.Labels[constants.StorageLabelKey])
	if err != nil {
		return err
	}
	defaultArgs = append(defaultArgs, "--objstore.config="+string(storageConfig))

	if r.store.Spec.IndexCacheConfig != nil &&
		r.store.Spec.IndexCacheConfig.InMemoryIndexCacheConfig != nil &&
		r.store.Spec.IndexCacheConfig.MaxSize != "" {
		defaultArgs = append(defaultArgs, "--index-cache-size="+r.store.Spec.IndexCacheConfig.MaxSize)
	}

	if r.store.Spec.LogLevel != "" {
		defaultArgs = append(defaultArgs, "--log.level="+r.store.Spec.LogLevel)
	}
	if r.store.Spec.LogFormat != "" {
		defaultArgs = append(defaultArgs, "--log.format="+r.store.Spec.LogFormat)
	}
	if r.store.Spec.MinTime != "" {
		defaultArgs = append(defaultArgs, "--min-time="+r.store.Spec.MinTime)
	}
	if r.store.Spec.MaxTime != "" {
		defaultArgs = append(defaultArgs, "--max-time="+r.store.Spec.MaxTime)
	}

	if len(container.Args) > 0 && container.Args[0] != "store" {
		container.Args = append([]string{"store"}, container.Args...)
	}

	for name, value := range r.store.Spec.Flags {
		arg := fmt.Sprintf("--%s=%s", name, value)
		replaced := util.ReplaceInSlice(defaultArgs, func(v interface{}) bool {
			return util.GetArgName(v.(string)) == name
		}, arg)

		if !replaced {
			defaultArgs = append(defaultArgs, arg)
		}
	}

	for _, arg := range defaultArgs {

		replaced := util.ReplaceInSlice(container.Args, func(v interface{}) bool {
			return util.GetArgName(v.(string)) == util.GetArgName(arg)
		}, arg)

		if !replaced {
			container.Args = append(container.Args, arg)
		}
	}

	return nil
}
