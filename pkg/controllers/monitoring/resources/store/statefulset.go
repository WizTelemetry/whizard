package store

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	mainContainerName = "store"
)

var (
	// unsupportedArgs is the args that are not allowed to be set by the user.
	unsupportedArgs = []string{
		"--http-address",
		"--grpc-address",
	}
)

func (r *Store) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name())}
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
			Name: mainContainerName,
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

	container.Image = r.store.Spec.Image
	container.ImagePullPolicy = r.store.Spec.ImagePullPolicy

	container.VolumeMounts = []corev1.VolumeMount{}
	r.AddTSDBVolume(sts, container, r.store.Spec.DataVolume)

	volumes, volumeMounts, err := r.VolumesAndVolumeMountsForStorage(r.store.Labels[constants.StorageLabelKey])
	if err != nil {
		return nil, "", err
	}
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, volumes...)
	container.VolumeMounts = append(container.VolumeMounts, volumeMounts...)

	if container.LivenessProbe == nil {
		container.LivenessProbe = r.DefaultLivenessProbe()
	}

	if container.ReadinessProbe == nil {
		container.ReadinessProbe = r.DefaultReadinessProbe()
	}

	container.Resources = r.store.Spec.Resources

	storageHash, err := r.GetStorageHash(r.store.Labels[constants.StorageLabelKey])
	if err != nil {
		return nil, "", err
	}

	env := corev1.EnvVar{
		Name:  constants.StorageHash,
		Value: storageHash,
	}
	replaced := util.ReplaceInSlice(container.Env, func(v interface{}) bool {
		return v.(corev1.EnvVar).Name == constants.StorageHash
	}, env)
	if !replaced {
		container.Env = append(container.Env, env)
	}

	tenantHash, err := r.GetTenantHash(map[string]string{
		constants.StorageLabelKey: r.store.Labels[constants.StorageLabelKey],
		constants.ServiceLabelKey: r.store.Labels[constants.ServiceLabelKey],
	})
	if err != nil {
		return nil, "", err
	}

	env = corev1.EnvVar{
		Name:  constants.TenantHash,
		Value: tenantHash,
	}
	replaced = util.ReplaceInSlice(container.Env, func(v interface{}) bool {
		return v.(corev1.EnvVar).Name == constants.TenantHash
	}, env)
	if !replaced {
		container.Env = append(container.Env, env)
	}

	if args, err := r.megerArgs(); err != nil {
		return nil, "", err
	} else {
		container.Args = args
	}

	if needToAppend {
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *container)
	}

	return sts, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.store, sts, r.Scheme)
}

func (r *Store) megerArgs() ([]string, error) {
	storageConfig, err := r.GetStorageConfig(r.store.Labels[constants.StorageLabelKey])
	if err != nil {
		return nil, err
	}
	relabelConfig, err := r.createRelabelConfig()
	if err != nil {
		return nil, err
	}
	defaultArgs := []string{
		"store",
		fmt.Sprintf("--data-dir=%s", constants.StorageDir),
		"--objstore.config=" + string(storageConfig),
		"--selector.relabel-config=" + relabelConfig,
	}

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

	for _, flag := range r.store.Spec.Flags {
		arg := util.GetArgName(flag)
		if util.Contains(unsupportedArgs, arg) {
			klog.V(3).Infof("ignore the unsupported flag %s", arg)
			continue
		}

		replaced := util.ReplaceInSlice(defaultArgs, func(v interface{}) bool {
			return util.GetArgName(v.(string)) == util.GetArgName(flag)
		}, flag)
		if !replaced {
			defaultArgs = append(defaultArgs, flag)
		}
	}

	sort.Strings(defaultArgs[1:])
	return defaultArgs, nil
}

func (r *Store) createRelabelConfig() (string, error) {
	namespacedName := strings.Split(r.store.Labels[constants.ServiceLabelKey], ".")
	svc := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName[1],
			Namespace: namespacedName[0],
		},
	}
	if err := r.Client.Get(r.Context, client.ObjectKeyFromObject(svc), svc); err != nil {
		return "", err
	}

	label := svc.Spec.TenantLabelName
	if len(label) == 0 {
		label = constants.DefaultTenantLabelName
	}
	var tenants []string
	tenantList := &v1alpha1.TenantList{}
	err := r.Client.List(r.Context, tenantList, client.MatchingLabels(map[string]string{
		constants.StorageLabelKey: r.store.Labels[constants.StorageLabelKey],
		constants.ServiceLabelKey: r.store.Labels[constants.ServiceLabelKey],
	}))
	if err != nil {
		return "", err
	}

	for _, item := range tenantList.Items {
		if item.DeletionTimestamp != nil || !item.DeletionTimestamp.IsZero() {
			continue
		}
		tenants = append(tenants, item.Spec.Tenant)
	}
	sort.Strings(tenants)
	return util.CreateKeepTenantsRelabelConfig(label, tenants)

}
