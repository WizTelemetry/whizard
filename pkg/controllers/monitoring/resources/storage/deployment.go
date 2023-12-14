package storage

import (
	"fmt"

	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
	"github.com/prometheus-operator/prometheus-operator/pkg/k8sutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	webContainerName = "web"
	gcContainerName  = "gc"
)

func (s *Storage) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: s.meta(s.name())}

	if s.storage.Spec.BlockManager == nil ||
		s.storage.Spec.BlockManager.Enable == nil ||
		*s.storage.Spec.BlockManager.Enable == false {
		return d, resources.OperationDelete, nil
	}

	if err := s.Client.Get(s.Context, client.ObjectKeyFromObject(d), d); err != nil {
		if !util.IsNotFound(err) {
			return nil, "", err
		}
	}

	d.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: s.labels(),
	}

	if d.Spec.Replicas == nil || *d.Spec.Replicas == 0 {
		d.Spec.Replicas = s.storage.Spec.BlockManager.Replicas
	}

	d.Spec.Template.Labels = s.labels()

	d.Spec.Template.Spec.Affinity = s.storage.Spec.BlockManager.Affinity
	d.Spec.Template.Spec.NodeSelector = s.storage.Spec.BlockManager.NodeSelector
	d.Spec.Template.Spec.Tolerations = s.storage.Spec.BlockManager.Tolerations
	d.Spec.Template.Spec.ServiceAccountName = s.storage.Spec.BlockManager.ServiceAccountName
	d.Spec.Template.Spec.SecurityContext = s.storage.Spec.BlockManager.SecurityContext

	if s.storage.Spec.BlockManager.ImagePullSecrets != nil && len(s.storage.Spec.BlockManager.ImagePullSecrets) > 0 {
		d.Spec.Template.Spec.ImagePullSecrets = s.storage.Spec.BlockManager.ImagePullSecrets
	}

	volumes, volumeMounts, err := s.VolumesAndVolumeMountsForStorage(fmt.Sprintf("%s.%s", s.storage.Namespace, s.storage.Name))
	if err != nil {
		return nil, "", err
	}
	d.Spec.Template.Spec.Volumes = volumes

	var webContainer *corev1.Container
	for i := 0; i < len(d.Spec.Template.Spec.Containers); i++ {
		if d.Spec.Template.Spec.Containers[i].Name == webContainerName {
			webContainer = &d.Spec.Template.Spec.Containers[i]
		}
	}

	needToAppend := false
	if webContainer == nil {
		webContainer = &corev1.Container{
			Name: webContainerName,
			Ports: []corev1.ContainerPort{
				{
					Protocol:      corev1.ProtocolTCP,
					Name:          constants.HTTPPortName,
					ContainerPort: constants.HTTPPort,
				},
			},
		}

		needToAppend = true
	}

	webContainer.VolumeMounts = volumeMounts

	webContainer.Image = s.storage.Spec.BlockManager.Image
	webContainer.ImagePullPolicy = s.storage.Spec.BlockManager.ImagePullPolicy

	if webContainer.LivenessProbe == nil {
		webContainer.LivenessProbe = s.DefaultLivenessProbe()
	}

	if webContainer.ReadinessProbe == nil {
		webContainer.ReadinessProbe = s.DefaultReadinessProbe()
	}

	webContainer.Resources = s.storage.Spec.BlockManager.Resources

	storageConfig, err := s.GetStorageConfig(util.Join(".", s.storage.Namespace, s.storage.Name))
	if err != nil {
		return nil, "", err
	}

	webContainer.Args = []string{
		"tools",
		"bucket",
		"web",
		"--objstore.config=" + string(storageConfig),
	}

	if s.storage.Spec.BlockManager.LogLevel != "" {
		webContainer.Args = append(webContainer.Args, "--log.level="+s.storage.Spec.BlockManager.LogLevel)
	}
	if s.storage.Spec.BlockManager.LogFormat != "" {
		webContainer.Args = append(webContainer.Args, "--log.format="+s.storage.Spec.BlockManager.LogFormat)
	}

	if s.storage.Spec.BlockManager.BlockSyncInterval != nil {
		webContainer.Args = append(webContainer.Args, "--refresh="+s.storage.Spec.BlockManager.BlockSyncInterval.Duration.String())
	}

	if needToAppend {
		d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, *webContainer)
	}

	var gcContainer *corev1.Container
	for i := 0; i < len(d.Spec.Template.Spec.Containers); i++ {
		if d.Spec.Template.Spec.Containers[i].Name == gcContainerName {
			if s.isGCEnabled() {
				gcContainer = &d.Spec.Template.Spec.Containers[i]
			} else {
				//delete gc container if gc is disable
				d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers[:i], d.Spec.Template.Spec.Containers[i+1:]...)
			}

			break
		}
	}

	if s.isGCEnabled() {
		needToAppend := false
		if gcContainer == nil {
			gcContainer = &corev1.Container{
				Name:            gcContainerName,
				Image:           s.storage.Spec.BlockManager.GC.Image,
				ImagePullPolicy: s.storage.Spec.BlockManager.GC.ImagePullPolicy,
			}
			needToAppend = true
		}

		gcContainer.VolumeMounts = volumeMounts

		gcContainer.Resources = s.storage.Spec.BlockManager.GC.Resources

		args := []string{"--objstore.config=" + string(storageConfig)}

		if s.storage.Spec.BlockManager.GC.GCInterval != nil &&
			s.storage.Spec.BlockManager.GC.GCInterval.Duration != 0 {
			args = append(args, "--gc.interval="+s.storage.Spec.BlockManager.GC.GCInterval.Duration.String())
		}

		if s.storage.Spec.BlockManager.GC.CleanupTimeout != nil &&
			s.storage.Spec.BlockManager.GC.CleanupTimeout.Duration != 0 {
			args = append(args, "--gc.cleanup-timeout="+s.storage.Spec.BlockManager.GC.CleanupTimeout.Duration.String())
		}

		if s.storage.Spec.BlockManager.GC.DefaultTenantId != "" {
			args = append(args, "--tenant.default-id="+s.storage.Spec.BlockManager.GC.DefaultTenantId)
		}

		if s.storage.Spec.BlockManager.GC.TenantLabelName != "" {
			args = append(args, "--tenant.label-name="+s.storage.Spec.BlockManager.GC.TenantLabelName)
		}

		gcContainer.Args = args

		if needToAppend {
			d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, *gcContainer)
		}
	}

	if len(s.storage.Spec.BlockManager.Containers) > 0 {
		containers, err := k8sutil.MergePatchContainers(d.Spec.Template.Spec.Containers, s.storage.Spec.BlockManager.Containers)
		if err != nil {
			return nil, "", fmt.Errorf("failed to merge containers spec: %w", err)
		}
		d.Spec.Template.Spec.Containers = containers
	}

	return d, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(s.storage, d, s.Scheme)
}

func (s *Storage) isGCEnabled() bool {
	if s.storage.Spec.BlockManager.GC == nil ||
		s.storage.Spec.BlockManager.GC.Enable == nil ||
		!*s.storage.Spec.BlockManager.GC.Enable {
		return false
	}

	return true
}
