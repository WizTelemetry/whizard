package ruler

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/query"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/queryfrontend"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/router"
	monitoringgateway "github.com/kubesphere/whizard/pkg/monitoring-gateway"
	"github.com/kubesphere/whizard/pkg/util"
	"github.com/prometheus-operator/prometheus-operator/pkg/k8sutil"
	promoperator "github.com/prometheus-operator/prometheus-operator/pkg/operator"
	promcommonconfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	promconfig "github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/model/relabel"
	"github.com/thanos-io/thanos/pkg/httpconfig"
	"gopkg.in/yaml.v3"
	yamlv3 "gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// repeatableArgs is the args that can be set repeatedly.
	// An error will occur if a non-repeatable arg is set repeatedly.
	repeatableArgs = []string{
		"--query",
		"--query.sd-files",
		"--rule-file",
		"--alertmanagers.url",
		"--alert.label-drop",
	}
	// unsupportedArgs is the args that are not allowed to be set by the user.
	unsupportedArgs = []string{
		"--receive.local-endpoint",
		"--http-address",
		"--grpc-address",
	}
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
				NodeSelector:    r.ruler.Spec.NodeSelector,
				Tolerations:     r.ruler.Spec.Tolerations,
				Affinity:        r.ruler.Spec.Affinity,
				SecurityContext: r.ruler.Spec.SecurityContext,
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
		LivenessProbe:  r.DefaultLivenessProbe(),
		ReadinessProbe: r.DefaultReadinessProbe(),
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
	rn := k8sutil.NewResourceNamerWithPrefix("configmap")
	for cmName := range r.shardsRuleConfigMapNames[shardSn] {
		name, err := rn.DNS1123Label(cmName)
		if err != nil {
			return nil, "", err
		}
		vol := corev1.Volume{
			Name: name,
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

	r.AddTSDBVolume(sts, &container, r.ruler.Spec.DataVolume)

	if r.ruler.Spec.AlertmanagersConfig != nil {
		container.Args = append(container.Args, "--alertmanagers.config=$(ALERTMANAGERS_CONFIG)")
		container.Env = append(container.Env, corev1.EnvVar{
			Name: "ALERTMANAGERS_CONFIG",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: r.ruler.Spec.AlertmanagersConfig,
			},
		})
	} else if len(r.ruler.Spec.AlertmanagersURL) > 0 {
		for _, url := range r.ruler.Spec.AlertmanagersURL {
			container.Args = append(container.Args, fmt.Sprintf("--alertmanagers.url=%s", url))
		}
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

	namespacedName := util.ServiceNamespacedName(r.ruler)

	if namespacedName != nil {
		var service monitoringv1alpha1.Service
		if err := r.Client.Get(r.Context, *namespacedName, &service); err != nil {
			return nil, resources.OperationCreateOrUpdate, err
		}

		writeAddr, err := r.remoteWriteAddress()
		if err != nil {
			return nil, "", err
		}
		queryAddr, err := r.queryAddress()
		if err != nil {
			return nil, "", err
		}
		var rwCfg = &promconfig.RemoteWriteConfig{}
		*rwCfg = promconfig.DefaultRemoteWriteConfig
		var cfgs struct {
			RemoteWriteConfigs []*promconfig.RemoteWriteConfig `yaml:"remote_write,omitempty"`
		}

		// If there is remote-writes configured in the related service instance,
		// the rulers should also write calculated metrics to them.
		//
		// TODO currently global rulers write directly to remote targets without tenant header.
		// 		If the tenant header is required, do it later.
		for _, rw := range service.Spec.RemoteWrites {
			var rwCfg = &promconfig.RemoteWriteConfig{}
			*rwCfg = promconfig.DefaultRemoteWriteConfig
			rwCfg.Name = rw.Name
			writeUrl, err := url.Parse(rw.URL)
			if err != nil {
				return nil, "", fmt.Errorf("invalid remote write url: %s", rw.URL)
			}
			rwCfg.URL = &promcommonconfig.URL{URL: writeUrl}
			if writeUrl.Scheme == "https" {
				// TODO support certificate validation
				rwCfg.HTTPClientConfig.TLSConfig.InsecureSkipVerify = true
			}
			if !reflect.DeepEqual(rw.HTTPClientConfig.BasicAuth, v1alpha1.BasicAuth{}) {
				secret := &corev1.Secret{}
				if err := r.Client.Get(r.Context, client.ObjectKey{Name: rw.HTTPClientConfig.BasicAuth.Username.Name, Namespace: r.Service.Namespace}, secret); err != nil {
					return nil, "", err
				}
				username := string(secret.Data[rw.HTTPClientConfig.BasicAuth.Username.Key])
				if err := r.Client.Get(r.Context, client.ObjectKey{Name: rw.HTTPClientConfig.BasicAuth.Password.Name, Namespace: r.Service.Namespace}, secret); err != nil {
					return nil, "", err
				}
				password := promcommonconfig.Secret(secret.Data[rw.HTTPClientConfig.BasicAuth.Password.Key])
				rwCfg.HTTPClientConfig.BasicAuth = &promcommonconfig.BasicAuth{
					Username: username,
					Password: password,
				}

				//	basicAuthEnc := func(username, password string) string {
				//		auth := username + ":" + password
				//		return base64.StdEncoding.EncodeToString([]byte(auth))
				//	}(username, strings.TrimSpace(string(password)))
				//	if len(rwCfg.Headers) == 0 {
				//		rwCfg.Headers = make(map[string]string, 1)
				//	}
				//	rwCfg.Headers["Authorization"] = "Basic " + basicAuthEnc
			}
			if rw.HTTPClientConfig.BearerToken != "" {
				rwCfg.HTTPClientConfig.BearerToken = promcommonconfig.Secret(rw.HTTPClientConfig.BearerToken)
				//	if len(rwCfg.Headers) == 0 {
				//		rwCfg.Headers = make(map[string]string, 1)
				//	}
				//	bearerEnc := fmt.Sprintf("%s %s", "Bearer", string(rw.HTTPClientConfig.BearerToken))
				//	rwCfg.Headers["Authorization"] = bearerEnc
			}
			rwCfg.Headers = rw.Headers
			if rw.RemoteTimeout != "" {
				timeout, err := time.ParseDuration(string(rw.RemoteTimeout))
				if err != nil {
					return nil, "", fmt.Errorf("invalid remote timeout: %s", rw.RemoteTimeout)
				}
				rwCfg.RemoteTimeout = model.Duration(timeout)
			}

			cfgs.RemoteWriteConfigs = append(cfgs.RemoteWriteConfigs, rwCfg)
		}

		// proxy config
		// if the tenant exists, append QueryProxy
		// otherwise, append WriteProxy
		if r.ruler.Spec.Tenant == "" { // as global ruler if no tenant specified
			var hasQueryFlag bool
			for _, flag := range r.ruler.Spec.Flags {
				if strings.HasPrefix(flag, "--query=") {
					hasQueryFlag = true
					break
				}
			}
			// If --query flag in spec.Flags is not specified, the global ruler will query from the queryAddr
			// which points to the QueryFrontend/Query under the same whizard service.
			// If and only if the ruler needs to query external data out of whizard service, the --query flag can be specified.
			if !hasQueryFlag {
				url, err := url.Parse(queryAddr)
				if err != nil {
					return nil, "", err
				}
				queryconfig := httpconfig.Config{
					EndpointsConfig: httpconfig.EndpointsConfig{
						Scheme:          url.Scheme,
						StaticAddresses: []string{url.Host},
					},
				}
				if url.Scheme == "https" {
					queryconfig.HTTPClientConfig = httpconfig.ClientConfig{
						TLSConfig: httpconfig.TLSConfig{
							InsecureSkipVerify: true,
						},
					}
				}
				if r.Service.Spec.RemoteQuery != nil {
					if !reflect.DeepEqual(r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth, v1alpha1.BasicAuth{}) {
						secret := &corev1.Secret{}

						if err := r.Client.Get(r.Context, client.ObjectKey{Name: r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth.Username.Name, Namespace: r.Service.Namespace}, secret); err != nil {
							return nil, "", err
						}
						queryconfig.HTTPClientConfig.BasicAuth.Username = string(secret.Data[r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth.Username.Key])
						if err := r.Client.Get(r.Context, client.ObjectKey{Name: r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth.Password.Name, Namespace: r.Service.Namespace}, secret); err != nil {
							return nil, "", err
						}
						queryconfig.HTTPClientConfig.BasicAuth.Password = string(secret.Data[r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth.Password.Key])
					}
					if r.Service.Spec.RemoteQuery.HTTPClientConfig.BearerToken != "" {
						queryconfig.HTTPClientConfig.BearerToken = string(r.Service.Spec.RemoteQuery.HTTPClientConfig.BearerToken)
					}
				}
				queryConfigs := []httpconfig.Config{}
				queryConfigs = append(queryConfigs, queryconfig)
				buff, _ := yamlv3.Marshal(queryConfigs)
				container.Args = append(container.Args, "--query.config="+string(buff))
			}

			writeUrl, err := url.Parse("http://127.0.0.1:8081/push")
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			rwCfg.URL = &promcommonconfig.URL{URL: writeUrl}

			cfgs.RemoteWriteConfigs = append(cfgs.RemoteWriteConfigs, rwCfg)
			content, err := yamlv3.Marshal(&cfgs)
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			root := &yaml.Node{}
			if err := yaml.Unmarshal(content, root); err != nil {
				return nil, "", err
			}
			for _, rwCfg := range cfgs.RemoteWriteConfigs {
				if rwCfg.HTTPClientConfig.BasicAuth != nil {
					if n := findYamlNodeByKey(root, "password"); n != nil {
						n.SetString(string(rwCfg.HTTPClientConfig.BasicAuth.Password))
					}
				}
				if rwCfg.HTTPClientConfig.BearerToken != "" {
					if n := findYamlNodeByKey(root, "bearer_token"); n != nil {
						n.SetString(string(rwCfg.HTTPClientConfig.BearerToken))
					}
				}
			}

			body, err := yaml.Marshal(root)
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			container.Args = append(container.Args, "--remote-write.config="+string(body))

			if url, err := url.Parse(writeAddr); err == nil && url.Scheme == "https" {
				writeAddr = "http://127.0.0.1:" + constants.CustomProxyPort

				data := make(map[string]string, 4)

				data["ProxyServiceEnabled"] = "true"
				data["ProxyLocalListenPort"] = constants.CustomProxyPort
				data["ProxyServiceAddress"] = url.Hostname()
				data["ProxyServicePort"] = url.Port()

				if err := r.envoyConfigMap(data); err != nil {
					return nil, "", err
				}
				var volumeMounts = []corev1.VolumeMount{}
				var volumes = []corev1.Volume{}

				volumes, volumeMounts, _ = resources.BuildCommonVolumes(nil, r.name("envoy-config"), nil, nil)

				envoyContainer := resources.BuildEnvoySidecarContainer(r.ruler.Spec.Envoy, volumeMounts)
				sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, envoyContainer)
				sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, volumes...)

			}

			writeProxyContainer, err := r.addWriteProxyContainer(&service.Spec, writeAddr)
			if err != nil {
				return nil, "", err
			}
			sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *writeProxyContainer)
		} else {
			container.Args = append(container.Args,
				fmt.Sprintf("--query=http://127.0.0.1:9080/%s", r.ruler.Spec.Tenant))

			//	rewrite proxy container
			// writeUrl, err := url.Parse(writeAddr + "/api/v1/receive")
			writeUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:9080/%s/api/v1/receive", r.ruler.Spec.Tenant))

			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			rwCfg.URL = &promcommonconfig.URL{URL: writeUrl}
			cfgs.RemoteWriteConfigs = append(cfgs.RemoteWriteConfigs, rwCfg)

			for i := range cfgs.RemoteWriteConfigs {
				rwCfg := cfgs.RemoteWriteConfigs[i]
				if rwCfg.Headers == nil {
					rwCfg.Headers = make(map[string]string)
				}
				rwCfg.Headers[service.Spec.TenantHeader] = r.ruler.Spec.Tenant
				rwCfg.WriteRelabelConfigs = append(rwCfg.WriteRelabelConfigs,
					&relabel.Config{
						Regex:  relabel.MustNewRegexp(service.Spec.TenantLabelName),
						Action: relabel.LabelDrop,
					})
				cfgs.RemoteWriteConfigs[i] = rwCfg
			}
			content, err := yamlv3.Marshal(&cfgs)
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			root := &yaml.Node{}
			if err := yaml.Unmarshal(content, root); err != nil {
				return nil, "", err
			}
			for _, rwCfg := range cfgs.RemoteWriteConfigs {
				if rwCfg.HTTPClientConfig.BasicAuth != nil {
					if n := findYamlNodeByKey(root, "password"); n != nil {
						n.SetString(string(rwCfg.HTTPClientConfig.BasicAuth.Password))
					}
				}
				if rwCfg.HTTPClientConfig.BearerToken != "" {
					if n := findYamlNodeByKey(root, "bearer_token"); n != nil {
						n.SetString(string(rwCfg.HTTPClientConfig.BearerToken))
					}
				}
			}

			body, err := yaml.Marshal(root)
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			container.Args = append(container.Args, "--remote-write.config="+string(body))

			queryProxyContainer, _ := r.addQueryProxyContainer(&service.Spec, queryAddr, writeAddr)
			sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *queryProxyContainer)
		}
	}

	for _, flag := range r.ruler.Spec.Flags {
		arg := util.GetArgName(flag)
		if util.Contains(unsupportedArgs, arg) {
			klog.V(3).Infof("ignore the unsupported flag %s", arg)
			continue
		}

		if util.Contains(repeatableArgs, arg) {
			container.Args = append(container.Args, flag)
			continue
		}

		replaced := util.ReplaceInSlice(container.Args, func(v interface{}) bool {
			return util.GetArgName(v.(string)) == util.GetArgName(flag)
		}, flag)
		if !replaced {
			container.Args = append(container.Args, flag)
		}
	}

	sort.Strings(container.Args[1:])

	var reloaderConfig = promoperator.ContainerConfig{
		Image:         r.Options.PrometheusConfigReloader.Image,
		CPURequest:    r.Options.PrometheusConfigReloader.Resources.Requests.Cpu().String(),
		MemoryRequest: r.Options.PrometheusConfigReloader.Resources.Requests.Memory().String(),
		CPULimit:      r.Options.PrometheusConfigReloader.Resources.Limits.Cpu().String(),
		MemoryLimit:   r.Options.PrometheusConfigReloader.Resources.Limits.Memory().String(),
	}

	var reloadContainer = promoperator.CreateConfigReloader(
		"config-reloader",
		promoperator.ReloaderConfig(reloaderConfig),
		promoperator.ReloaderURL(url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("127.0.0.1:%d", constants.HTTPPort),
			Path:   path.Clean("/-/reload"),
		}),
		promoperator.ListenLocal(true),
		promoperator.LocalHost("127.0.0.1"),
		promoperator.LogFormat(r.ruler.Spec.LogFormat),
		promoperator.LogLevel(r.ruler.Spec.LogLevel),
		promoperator.WatchedDirectories(watchedDirectories),
		promoperator.VolumeMounts(configReloaderVolumeMounts),
		promoperator.Shard(-1),
	)

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container, reloadContainer)

	if r.ruler.Spec.ImagePullSecrets != nil && len(r.ruler.Spec.ImagePullSecrets) > 0 {
		sts.Spec.Template.Spec.ImagePullSecrets = r.ruler.Spec.ImagePullSecrets
	}

	if len(r.ruler.Spec.Containers) > 0 {
		containers, err := k8sutil.MergePatchContainers(sts.Spec.Template.Spec.Containers, r.ruler.Spec.Containers)
		if err != nil {
			return nil, "", fmt.Errorf("failed to merge containers spec: %w", err)
		}
		sts.Spec.Template.Spec.Containers = containers
	}

	return sts, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.ruler, sts, r.Scheme)
}

func (r *Ruler) remoteWriteAddress() (string, error) {
	routerList := &monitoringv1alpha1.RouterList{}
	if err := r.Client.List(r.Context, routerList, client.MatchingLabels(util.ManagedLabelBySameService(r.ruler))); err != nil {
		return "", err
	}

	if len(routerList.Items) > 0 {
		if len(routerList.Items) > 1 {
			return "", fmt.Errorf("more than one router defined for service %s/%s", r.Service.Name, r.Service.Namespace)
		}

		o := routerList.Items[0]
		r, err := router.New(r.BaseReconciler, &o, nil)
		if err != nil {
			return "", err
		}
		if o.Spec.HTTPServerTLSConfig != nil {
			return r.RemoteWriteHTTPSAddr(), nil
		}

		return r.RemoteWriteAddr(), nil
	}

	return "", fmt.Errorf("no router defined for service %s/%s", r.Service.Name, r.Service.Namespace)
}

// queryAddress returns the address from which the ruler should query.
// The ruler should query from the QueryFrontend with a higher performance the Query
// because the feature https://thanos.io/v0.31/proposals-accepted/202205-vertical-query-sharding/
func (r *Ruler) queryAddress() (string, error) {
	queryFrontendList := &monitoringv1alpha1.QueryFrontendList{}
	if err := r.Client.List(r.Context, queryFrontendList, client.MatchingLabels(util.ManagedLabelBySameService(r.ruler))); err != nil {
		return "", err
	}

	if len(queryFrontendList.Items) > 0 {
		if len(queryFrontendList.Items) > 1 {
			return "", fmt.Errorf("more than one query frontend defined for service %s/%s", r.Service.Name, r.Service.Namespace)
		}

		q := queryFrontendList.Items[0]
		r, err := queryfrontend.New(r.BaseReconciler, &q)
		if err != nil {
			return "", err
		}
		if q.Spec.HTTPServerTLSConfig != nil {
			return r.HttpsAddr(), nil
		}

		return r.HttpAddr(), nil
	}

	queryList := &monitoringv1alpha1.QueryList{}
	if err := r.Client.List(r.Context, queryList, client.MatchingLabels(util.ManagedLabelBySameService(r.ruler))); err != nil {
		return "", err
	}

	if len(queryList.Items) > 0 {
		if len(queryList.Items) > 1 {
			return "", fmt.Errorf("more than one query defined for service %s/%s", r.Service.Name, r.Service.Namespace)
		}

		o := queryList.Items[0]
		r, err := query.New(r.BaseReconciler, &o, nil)
		if err != nil {
			return "", err
		}
		if o.Spec.HTTPServerTLSConfig != nil {
			return r.HttpsAddr(), nil
		}
		return r.HttpAddr(), nil
	}

	return "", fmt.Errorf("no query frontend or query exist for service %s/%s", r.Service.Name, r.Service.Namespace)
}

type gatewayConfig struct {
	// The HTTP basic authentication credentials for the targets.
	BasicAuth *monitoringgateway.BasicAuth `yaml:"basic_auth,omitempty" json:"basic_auth,omitempty"`
	// The bearer token for the targets.
	BearerToken string `yaml:"bearer_token,omitempty" json:"bearer_token,omitempty"`

	TLSConfig *promcommonconfig.TLSConfig `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
}

type BasicAuth struct {
}

func (r *Ruler) addQueryProxyContainer(serviceSpec *monitoringv1alpha1.ServiceSpec, queryAddr, remoteWriteAddr string) (*corev1.Container, error) {

	queryProxyContainer := &corev1.Container{
		Name:  "query-proxy",
		Image: r.Options.RulerQueryProxy.Image,
		Args: []string{
			"--http-address=127.0.0.1:9080",
		},
		Resources: r.Options.RulerQueryProxy.Resources,
	}
	queryProxyContainer.Args = append(queryProxyContainer.Args, "--tenant.label-name="+serviceSpec.TenantLabelName)
	queryProxyContainer.Args = append(queryProxyContainer.Args, "--tenant.header="+serviceSpec.TenantHeader)

	var cfg = gatewayConfig{}
	if url, err := url.Parse(queryAddr); err == nil && url.Scheme == "https" {
		cfg.TLSConfig = &promcommonconfig.TLSConfig{
			InsecureSkipVerify: true,
		}
	}

	if r.Service.Spec.RemoteQuery != nil {
		if !reflect.DeepEqual(r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth, v1alpha1.BasicAuth{}) {
			secret := &corev1.Secret{}
			cfg.BasicAuth = &monitoringgateway.BasicAuth{}
			if err := r.Client.Get(r.Context, client.ObjectKey{Name: r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth.Username.Name, Namespace: r.Service.Namespace}, secret); err != nil {
				return nil, err
			}
			cfg.BasicAuth.Username = string(secret.Data[r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth.Username.Key])
			if err := r.Client.Get(r.Context, client.ObjectKey{Name: r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth.Password.Name, Namespace: r.Service.Namespace}, secret); err != nil {
				return nil, err
			}
			cfg.BasicAuth.Password = string(secret.Data[r.Service.Spec.RemoteQuery.HTTPClientConfig.BasicAuth.Password.Key])
		}
		if r.Service.Spec.RemoteQuery.HTTPClientConfig.BearerToken != "" {
			cfg.BearerToken = string(r.Service.Spec.RemoteQuery.HTTPClientConfig.BearerToken)
		}
	}
	if !reflect.DeepEqual(cfg, config{}) {
		buff, _ := yamlv3.Marshal(cfg)
		queryProxyContainer.Args = append(queryProxyContainer.Args, fmt.Sprintf("--query.config=%s", buff))
	}
	queryProxyContainer.Args = append(queryProxyContainer.Args, "--query.address="+queryAddr)

	remoteWritesConigs := []monitoringgateway.RemoteWriteConfig{}
	rwcfg := monitoringgateway.RemoteWriteConfig{}
	if url, err := url.Parse(remoteWriteAddr); err == nil {
		rwcfg.URL = &promcommonconfig.URL{URL: url}
		if url.Scheme == "https" {
			rwcfg.TLSConfig = promcommonconfig.TLSConfig{
				InsecureSkipVerify: true,
			}
		}
		remoteWritesConigs = append(remoteWritesConigs, rwcfg)

		buff, _ := yamlv3.Marshal(remoteWritesConigs)
		queryProxyContainer.Args = append(queryProxyContainer.Args, fmt.Sprintf("--remote-writes.config=%s", buff))
	}

	return queryProxyContainer, nil
}

// cortex-tenant config
// https://github.com/blind-oracle/cortex-tenant/blob/main/config.go#L13
type config struct {
	Listen      string
	ListenPprof string `yaml:"listen_pprof"`

	Target string

	LogLevel        string `yaml:"log_level"`
	Timeout         time.Duration
	TimeoutShutdown time.Duration `yaml:"timeout_shutdown"`
	Concurrency     int
	Metadata        bool

	Auth struct {
		Egress struct {
			Username string
			Password string
		}
	}

	Tenant struct {
		Label       string
		LabelRemove bool `yaml:"label_remove"`
		Header      string
		Default     string
		AcceptAll   bool `yaml:"accept_all"`
	}
}

func (r *Ruler) addWriteProxyContainer(serviceSpec *monitoringv1alpha1.ServiceSpec, writeAddr string) (*corev1.Container, error) {
	var writeProxyContainer *corev1.Container
	cfg := &config{
		Listen:          "127.0.0.1:8081",
		LogLevel:        "warn",
		Timeout:         time.Second * 10,
		TimeoutShutdown: time.Second * 10,
		Concurrency:     1000,
		Metadata:        false,
	}

	writeUrl, err := url.Parse(writeAddr + "/api/v1/receive")
	if err != nil {
		return writeProxyContainer, err
	}
	cfg.Target = writeUrl.String()

	cfg.Tenant.Label = serviceSpec.TenantLabelName
	cfg.Tenant.LabelRemove = true
	cfg.Tenant.Header = serviceSpec.TenantHeader
	cfg.Tenant.Default = serviceSpec.DefaultTenantId

	cfgContent, err := yamlv3.Marshal(cfg)
	if err != nil {
		return writeProxyContainer, err
	}

	writeProxyContainer = &corev1.Container{
		Name:  "write-proxy",
		Image: r.Options.RulerWriteProxy.Image,
		Args: []string{
			"--config-content=" + string(cfgContent),
		},
		Resources: r.Options.RulerWriteProxy.Resources,
	}
	return writeProxyContainer, nil
}

func findYamlNodeByKey(root *yaml.Node, key string) *yaml.Node {

	for i := 0; i < len(root.Content); i++ {
		if root.Content[i].Value == key && i+1 < len(root.Content) {
			return root.Content[i+1]
		}

		if n := findYamlNodeByKey(root.Content[i], key); n != nil {
			return n
		}
	}
	return nil
}
