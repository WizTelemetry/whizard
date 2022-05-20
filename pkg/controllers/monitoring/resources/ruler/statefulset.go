package ruler

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	promoperator "github.com/prometheus-operator/prometheus-operator/pkg/operator"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/query"
)

func (r *Ruler) statefulSet(ruleConfigMapNames []string) (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name())}

	sts.Spec = appsv1.StatefulSetSpec{
		Replicas: r.ruler.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		ServiceName: r.name(resources.ServiceNameSuffixOperated),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.ruler.Spec.NodeSelector,
				Tolerations:  r.ruler.Spec.Tolerations,
				Affinity:     r.ruler.Spec.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "rule",
		Image:     r.ruler.Spec.Image,
		Args:      []string{"rule"},
		Resources: r.ruler.Spec.Resources,
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

	var watchedDirectories []string
	var configReloaderVolumeMounts []corev1.VolumeMount

	for _, cmName := range ruleConfigMapNames {
		vol := corev1.Volume{
			Name: "configmap-" + cmName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cmName,
					},
				},
			},
		}
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, vol)
		volMount := corev1.VolumeMount{
			Name:      vol.Name,
			MountPath: filepath.Join(rulesDir, cmName),
		}
		container.VolumeMounts = append(container.VolumeMounts, volMount)
		configReloaderVolumeMounts = append(configReloaderVolumeMounts, volMount)
		watchedDirectories = append(watchedDirectories, volMount.MountPath)
	}

	var tsdbVolume = &corev1.Volume{
		Name: "tsdb",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	if v := r.ruler.Spec.DataVolume; v != nil {
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

	if r.ruler.Spec.AlertManagersConfig != nil {
		container.Args = append(container.Args, "--alertmanagers.config=$(ALERTMANAGERS_CONFIG)")
		container.Env = append(container.Env, corev1.EnvVar{
			Name: "ALERTMANAGERS_CONFIG",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: r.ruler.Spec.AlertManagersConfig,
			},
		})
	} else if len(r.ruler.Spec.AlertManagersURL) > 0 {
		for _, url := range r.ruler.Spec.AlertManagersURL {
			container.Args = append(container.Args, fmt.Sprintf("--alertmanagers.url=%s", url))
		}
	}

	if r.ruler.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.ruler.Spec.LogLevel)
	}
	if r.ruler.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.ruler.Spec.LogFormat)
	}
	container.Args = append(container.Args, fmt.Sprintf(`--label=%s="$(POD_NAME)"`, resources.ReplicaLabelNameThanosRuler))
	for k, v := range r.ruler.Spec.Labels {
		container.Args = append(container.Args, fmt.Sprintf(`--label=%s="%s"`, k, v))
	}
	container.Args = append(container.Args, fmt.Sprintf("--data-dir=%s", storageDir))
	container.Args = append(container.Args, fmt.Sprintf("--rule-file=%s/*/*.yaml", rulesDir))
	container.Args = append(container.Args, fmt.Sprintf("--alert.label-drop=%s", resources.ReplicaLabelNameThanosRuler))
	for _, lb := range r.ruler.Spec.AlertDropLabels {
		container.Args = append(container.Args, fmt.Sprintf("--alert.label-drop=%s", lb))
	}
	if r.ruler.Spec.EvaluationInterval != "" {
		container.Args = append(container.Args, fmt.Sprintf("--eval-interval=%s", r.ruler.Spec.EvaluationInterval))
	}

	namespacedName := monitoringv1alpha1.ServiceNamespacedName(r.ruler)

	if namespacedName != nil {
		var service monitoringv1alpha1.Service
		if err := r.Client.Get(r.Context, *namespacedName, &service); err != nil {
			return nil, resources.OperationCreateOrUpdate, err
		}
		query := query.New(resources.ServiceBaseReconciler{
			BaseReconciler: r.BaseReconciler,
			Service:        &service,
		})
		container.Args = append(container.Args, "--query="+query.HttpAddr())

		container.Args = append(container.Args, "--remote-write.config=$(REMOTE_WRITE_CONFIG)")
		container.Env = append(container.Env, corev1.EnvVar{
			Name: "REMOTE_WRITE_CONFIG",
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: r.name("remote-write-config"),
					},
					Key: remoteWriteFile,
				},
			},
		})
	}

	var reloaderConfig = promoperator.ReloaderConfig{Image: r.reloaderConfig.Image}
	if r.reloaderConfig.CPURequest != "0" {
		reloaderConfig.CPURequest = r.reloaderConfig.CPURequest
	}
	if r.reloaderConfig.MemoryRequest != "0" {
		reloaderConfig.MemoryRequest = r.reloaderConfig.MemoryRequest
	}
	if r.reloaderConfig.CPULimit != "0" {
		reloaderConfig.CPULimit = r.reloaderConfig.CPULimit
	}
	if r.reloaderConfig.MemoryLimit != "0" {
		reloaderConfig.MemoryLimit = r.reloaderConfig.MemoryLimit
	}

	var reloadContainer = promoperator.CreateConfigReloader(
		"config-reloader",
		promoperator.ReloaderResources(reloaderConfig),
		promoperator.ReloaderURL(url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("localhost:%d", resources.ThanosHTTPPort),
			Path:   path.Clean("/-/reload"),
		}),
		promoperator.ListenLocal(true),
		promoperator.LocalHost("localhost"),
		promoperator.LogFormat(r.ruler.Spec.LogFormat),
		promoperator.LogLevel(r.ruler.Spec.LogLevel),
		promoperator.WatchedDirectories(watchedDirectories),
		promoperator.VolumeMounts(configReloaderVolumeMounts),
		promoperator.Shard(-1),
	)

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container, reloadContainer)

	return sts, resources.OperationCreateOrUpdate, nil
}
