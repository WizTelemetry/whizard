package storage

import (
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
	webContainerName = "web"
	gcContainerName  = "gc"
)

func (s *Storage) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: s.meta(s.name())}

	if s.storage.Spec.Bucket == nil ||
		s.storage.Spec.Bucket.Enable == nil ||
		*s.storage.Spec.Bucket.Enable == false {
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
		d.Spec.Replicas = s.storage.Spec.Bucket.Replicas
	}

	d.Spec.Template.Labels = s.labels()

	d.Spec.Template.Spec.Affinity = s.storage.Spec.Bucket.Affinity
	d.Spec.Template.Spec.NodeSelector = s.storage.Spec.Bucket.NodeSelector
	d.Spec.Template.Spec.Tolerations = s.storage.Spec.Bucket.Tolerations
	d.Spec.Template.Spec.ServiceAccountName = s.storage.Spec.Bucket.ServiceAccountName

	var webContainer *corev1.Container
	for i := 0; i < len(d.Spec.Template.Spec.Containers); i++ {
		if d.Spec.Template.Spec.Containers[i].Name == webContainerName {
			webContainer = &d.Spec.Template.Spec.Containers[i]
		}
	}

	needToAppend := false
	if webContainer == nil {
		webContainer = &corev1.Container{
			Name:            webContainerName,
			Image:           s.storage.Spec.Bucket.Image,
			ImagePullPolicy: s.storage.Spec.Bucket.ImagePullPolicy,
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

	if webContainer.LivenessProbe == nil {
		webContainer.LivenessProbe = s.DefaultLivenessProbe()
	}

	if webContainer.ReadinessProbe == nil {
		webContainer.ReadinessProbe = s.DefaultReadinessProbe()
	}

	webContainer.Resources = s.storage.Spec.Bucket.Resources

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

	if s.storage.Spec.Bucket.LogLevel != "" {
		webContainer.Args = append(webContainer.Args, "--log.level="+s.storage.Spec.Bucket.LogLevel)
	}
	if s.storage.Spec.Bucket.LogFormat != "" {
		webContainer.Args = append(webContainer.Args, "--log.format="+s.storage.Spec.Bucket.LogFormat)
	}

	if s.storage.Spec.Bucket.Refresh != nil {
		webContainer.Args = append(webContainer.Args, "--refresh="+s.storage.Spec.Bucket.Refresh.Duration.String())
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
				Image:           s.storage.Spec.Bucket.GC.Image,
				ImagePullPolicy: s.storage.Spec.Bucket.GC.ImagePullPolicy,
			}
			needToAppend = true
		}

		gcContainer.Resources = s.storage.Spec.Bucket.Resources

		gcContainer.Args = []string{"--objstore.config=" + string(storageConfig)}

		if s.storage.Spec.Bucket.GC.Interval != nil &&
			s.storage.Spec.Bucket.GC.Interval.Duration != 0 {
			gcContainer.Args = append(gcContainer.Args, "--interval="+s.storage.Spec.Bucket.GC.Interval.String())
		}

		if s.storage.Spec.Bucket.GC.CleanupTimeout != nil &&
			s.storage.Spec.Bucket.GC.CleanupTimeout.Duration != 0 {
			gcContainer.Args = append(gcContainer.Args, "--cleanup-timeout="+s.storage.Spec.Bucket.GC.CleanupTimeout.String())
		}

		if needToAppend {
			d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, *gcContainer)
		}
	}

	return d, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(s.storage, d, s.Scheme)
}

func (s *Storage) isGCEnabled() bool {
	if s.storage.Spec.Bucket.GC == nil ||
		s.storage.Spec.Bucket.GC.Enable == nil ||
		!*s.storage.Spec.Bucket.GC.Enable {
		return false
	}

	return true
}
