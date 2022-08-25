package ruler

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	promoperator "github.com/prometheus-operator/prometheus-operator/pkg/operator"
	promcommonconfig "github.com/prometheus/common/config"
	promconfig "github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/model/relabel"
	yamlv3 "gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/query"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/router"
)

func (r *Ruler) statefulSets() (retResources []resources.Resource) {
	// for target statefulsets
	var targetNames = make(map[string]struct{}, *r.ruler.Spec.Shards)
	for i := 0; i < int(*r.ruler.Spec.Shards); i++ {
		shardSn := i
		targetNames[r.name(strconv.Itoa(shardSn))] = struct{}{}
		retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
			return r.statefulSet(shardSn)
		})
	}

	var stsList appsv1.StatefulSetList
	err := r.Client.List(r.Context, &stsList, client.InNamespace(r.ruler.Namespace))
	if err != nil {
		return errResourcesFunc(err)
	}
	// check statefulsets to be deleted.
	// the statefulsets owned by the ruler have a same name prefix and a shard sequence number suffix
	var namePrefix = r.name() + "-"
	for i := range stsList.Items {
		sts := stsList.Items[i]
		if !strings.HasPrefix(sts.Name, namePrefix) {
			continue
		}
		sn := strings.TrimPrefix(sts.Name, namePrefix)
		if sequenceNumberRegexp.MatchString(sn) {
			if _, ok := targetNames[sts.Name]; !ok {
				retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
					return &sts, resources.OperationDelete, nil
				})
			}
		}
	}
	return
}

func (r *Ruler) statefulSet(shardSn int) (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name(strconv.Itoa(shardSn)))}

	ls := r.labels()
	ls[constants.LabelNameRulerShardSn] = strconv.Itoa(shardSn)

	sts.Spec = appsv1.StatefulSetSpec{
		Replicas: r.ruler.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: ls,
		},
		ServiceName: r.name(strconv.Itoa(shardSn), constants.ServiceNameSuffix),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: ls,
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
				Name:          constants.GRPCPortName,
				ContainerPort: constants.GRPCPort,
			},
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          constants.HTTPPortName,
				ContainerPort: constants.HTTPPort,
			},
		},
		LivenessProbe:  resources.DefaultLivenessProbe(),
		ReadinessProbe: resources.DefaultReadinessProbe(),
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

	for cmName := range r.shardsRuleConfigMapNames[shardSn] {
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
	}

	if r.ruler.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.ruler.Spec.LogLevel)
	}
	if r.ruler.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.ruler.Spec.LogFormat)
	}
	container.Args = append(container.Args, fmt.Sprintf(`--label=%s="$(POD_NAME)"`, constants.RulerReplicaLabelName))
	for k, v := range r.ruler.Spec.Labels {
		container.Args = append(container.Args, fmt.Sprintf(`--label=%s="%s"`, k, v))
	}
	container.Args = append(container.Args, fmt.Sprintf("--data-dir=%s", storageDir))
	container.Args = append(container.Args, fmt.Sprintf("--rule-file=%s/*/*.yaml", rulesDir))
	container.Args = append(container.Args, fmt.Sprintf("--alert.label-drop=%s", constants.RulerReplicaLabelName))
	for _, lb := range r.ruler.Spec.AlertDropLabels {
		container.Args = append(container.Args, fmt.Sprintf("--alert.label-drop=%s", lb))
	}
	if r.ruler.Spec.EvaluationInterval != "" {
		container.Args = append(container.Args, fmt.Sprintf("--eval-interval=%s", r.ruler.Spec.EvaluationInterval))
	}

	for name, value := range r.ruler.Spec.Flags {
		container.Args = append(container.Args, fmt.Sprintf("--%s=%s", name, value))
	}

	var queryProxyContainer *corev1.Container

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

		if service.Spec.Router != nil {
			// remote write config
			receiveRouter := router.New(resources.ServiceBaseReconciler{
				BaseReconciler: r.BaseReconciler,
				Service:        &service,
			})
			writeUrl, err := url.Parse(receiveRouter.RemoteWriteAddr() + "/api/v1/receive")
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			var rwCfg = &promconfig.RemoteWriteConfig{}
			*rwCfg = promconfig.DefaultRemoteWriteConfig
			if rwCfg.Headers == nil {
				rwCfg.Headers = make(map[string]string)
			}
			rwCfg.Headers[service.Spec.TenantHeader] = r.ruler.Spec.Tenant
			rwCfg.URL = &promcommonconfig.URL{URL: writeUrl}
			var cfgs struct {
				RemoteWriteConfigs []*promconfig.RemoteWriteConfig `yaml:"remote_write,omitempty"`
			}
			rwCfg.WriteRelabelConfigs = append(rwCfg.WriteRelabelConfigs, &relabel.Config{
				Regex:  relabel.MustNewRegexp(service.Spec.TenantLabelName),
				Action: relabel.LabelDrop,
			})
			cfgs.RemoteWriteConfigs = append(cfgs.RemoteWriteConfigs, rwCfg)
			content, err := yamlv3.Marshal(&cfgs)
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			container.Args = append(container.Args, "--remote-write.config="+string(content))
		}

		// query config
		if r.ruler.Spec.Tenant == "" {
			container.Args = append(container.Args, "--query="+query.HttpAddr())
		} else {
			proxyResources := corev1.ResourceRequirements{
				Limits:   corev1.ResourceList{},
				Requests: corev1.ResourceList{},
			}
			if r.rulerQueryProxyConfig.CPURequest != "0" {
				proxyResources.Requests[corev1.ResourceCPU] = resource.MustParse(r.rulerQueryProxyConfig.CPURequest)
			}
			if r.rulerQueryProxyConfig.MemoryRequest != "0" {
				proxyResources.Requests[corev1.ResourceMemory] = resource.MustParse(r.rulerQueryProxyConfig.MemoryRequest)
			}
			if r.rulerQueryProxyConfig.CPULimit != "0" {
				proxyResources.Requests[corev1.ResourceCPU] = resource.MustParse(r.rulerQueryProxyConfig.CPULimit)
			}
			if r.rulerQueryProxyConfig.MemoryLimit != "0" {
				proxyResources.Requests[corev1.ResourceMemory] = resource.MustParse(r.rulerQueryProxyConfig.CPURequest)
			}
			queryProxyContainer = &corev1.Container{
				Name:  "query-proxy",
				Image: r.rulerQueryProxyConfig.Image,
				Args: []string{
					"--http-address=127.0.0.1:9080",
				},
				Resources: proxyResources,
			}
			queryProxyContainer.Args = append(queryProxyContainer.Args, "--tenant.label-name="+service.Spec.TenantLabelName)
			queryProxyContainer.Args = append(queryProxyContainer.Args, "--query.address="+query.HttpAddr())

			container.Args = append(container.Args,
				fmt.Sprintf("--query=http://127.0.0.1:9080/%s", r.ruler.Spec.Tenant))
		}
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
			Host:   fmt.Sprintf("localhost:%d", constants.HTTPPort),
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

	if queryProxyContainer != nil {
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *queryProxyContainer)
	}
	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container, reloadContainer)

	return sts, resources.OperationCreateOrUpdate, nil
}
