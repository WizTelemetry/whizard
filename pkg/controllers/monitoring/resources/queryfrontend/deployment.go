package queryfrontend

import (
	"fmt"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/query"
	"github.com/kubesphere/whizard/pkg/util"
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
		"--query-frontend.forward-header",
		"--query-frontend.org-id-header",
	}
	// unsupportedArgs is the args that are not allowed to be set by the user.
	unsupportedArgs = []string{
		// Deprecation
		"--log.request.decision",
		"--http-address",
	}
)

func (q *QueryFrontend) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: q.meta(q.name())}

	if q.queryFrontend == nil {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: q.queryFrontend.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: q.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: q.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector:    q.queryFrontend.Spec.NodeSelector,
				Tolerations:     q.queryFrontend.Spec.Tolerations,
				Affinity:        q.queryFrontend.Spec.Affinity,
				SecurityContext: q.queryFrontend.Spec.SecurityContext,
			},
		},
	}

	cacheConfigVol := corev1.Volume{
		Name: "cache-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: q.name("cache-config"),
				},
			},
		},
	}

	hashCode, err := q.GetTenantHash(map[string]string{
		constants.ServiceLabelKey: q.queryFrontend.Labels[constants.ServiceLabelKey],
	})
	if err != nil {
		return nil, "", err
	}

	var container = corev1.Container{
		Name:      "query-frontend",
		Image:     q.queryFrontend.Spec.Image,
		Args:      []string{"query-frontend"},
		Resources: q.queryFrontend.Spec.Resources,
		Ports: []corev1.ContainerPort{
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          constants.HTTPPortName,
				ContainerPort: constants.HTTPPort,
			},
		},
		LivenessProbe:  q.DefaultLivenessProbe(),
		ReadinessProbe: q.DefaultReadinessProbe(),
		VolumeMounts: []corev1.VolumeMount{{
			Name:      cacheConfigVol.Name,
			MountPath: configDir,
			ReadOnly:  true,
		}},
		Env: []corev1.EnvVar{
			{
				Name:  constants.TenantHash,
				Value: hashCode,
			},
		},
	}

	data := make(map[string]string, 8)

	// If there is remote-query configured in the related Service,
	// QueryFrontend will preferentially query from configured remote-query target,
	// else it'll query from the Query directly.
	var addr string
	if q.Service != nil && q.Service.Spec.RemoteQuery != nil {
		addr = q.Service.Spec.RemoteQuery.URL
	} else {
		var err error
		addr, err = q.queryAddress()
		if err != nil {
			return nil, "", err
		}
	}

	if url, err := url.Parse(addr); err == nil && url.Scheme == "https" {

		addr = "http://127.0.0.1:" + constants.CustomProxyPort

		data["ProxyServiceEnabled"] = "true"
		data["ProxyLocalListenPort"] = constants.CustomProxyPort
		data["ProxyServiceAddress"] = url.Hostname()
		data["ProxyServicePort"] = url.Port()

	} else if err != nil {
		return nil, "", err
	}

	if q.queryFrontend.Spec.HTTPServerTLSConfig != nil {
		data["LocalServiceEnabled"] = "true"
		data["ServiceMappingPort"] = strconv.Itoa(constants.HTTPPort)
		data["ServiceListenPort"] = constants.QueryFrontendHTTPPort
		data["ServiceTLSCertFile"] = constants.EnvoyCertsMountPath + q.queryFrontend.Spec.HTTPServerTLSConfig.CertSecret.Key
		data["ServiceTLSKeyFile"] = constants.EnvoyCertsMountPath + q.queryFrontend.Spec.HTTPServerTLSConfig.KeySecret.Key

		container.Args = append(container.Args, "--http-address=127.0.0.1:"+constants.QueryFrontendHTTPPort)
		container.LivenessProbe.HTTPGet.Scheme = "HTTPS"
		container.ReadinessProbe.HTTPGet.Scheme = "HTTPS"
	}

	if len(data) > 0 {
		if err := q.envoyConfigMap(data); err != nil {
			return nil, "", err
		}
		var volumeMounts = []corev1.VolumeMount{}
		var volumes = []corev1.Volume{}

		// Mount tls volume
		if q.queryFrontend.Spec.HTTPServerTLSConfig != nil {

			tlsAsset := []string{q.queryFrontend.Spec.HTTPServerTLSConfig.KeySecret.Name, q.queryFrontend.Spec.HTTPServerTLSConfig.CertSecret.Name}
			volumes, volumeMounts, _ = resources.BuildCommonVolumes(tlsAsset, q.name("envoy-config"), nil, nil)
		} else {
			volumes, volumeMounts, _ = resources.BuildCommonVolumes(nil, q.name("envoy-config"), nil, nil)
		}

		envoyContainer := resources.BuildEnvoySidecarContainer(q.queryFrontend.Spec.Envoy, volumeMounts)
		d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, envoyContainer)
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, volumes...)
	}

	container.Args = append(container.Args, "--query-frontend.downstream-url="+addr)
	container.Args = append(container.Args, "--labels.response-cache-config-file="+filepath.Join(configDir, cacheConfigFile))
	container.Args = append(container.Args, "--query-range.response-cache-config-file="+filepath.Join(configDir, cacheConfigFile))

	if q.queryFrontend.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+q.queryFrontend.Spec.LogLevel)
	}
	if q.queryFrontend.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+q.queryFrontend.Spec.LogFormat)
	}

	for _, flag := range q.queryFrontend.Spec.Flags {
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

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)
	d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, cacheConfigVol)

	if q.queryFrontend.Spec.ImagePullSecrets != nil && len(q.queryFrontend.Spec.ImagePullSecrets) > 0 {
		d.Spec.Template.Spec.ImagePullSecrets = q.queryFrontend.Spec.ImagePullSecrets
	}

	return d, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(q.queryFrontend, d, q.Scheme)
}

func (q *QueryFrontend) queryAddress() (string, error) {
	queryList := &v1alpha1.QueryList{}
	if err := q.Client.List(q.Context, queryList, client.MatchingLabels(util.ManagedLabelBySameService(q.queryFrontend))); err != nil {
		return "", err
	}

	if len(queryList.Items) > 0 {
		if len(queryList.Items) > 1 {
			return "", fmt.Errorf("more than one query defined for service %s/%s", q.Service.Name, q.Service.Namespace)
		}

		o := queryList.Items[0]
		r, err := query.New(q.BaseReconciler, &o, nil)
		if err != nil {
			return "", err
		}

		if o.Spec.HTTPServerTLSConfig != nil {
			return r.HttpsAddr(), nil
		}

		return r.HttpAddr(), nil
	}

	return "", fmt.Errorf("no query frontend or query exist for service %s/%s", q.Service.Name, q.Service.Namespace)
}
