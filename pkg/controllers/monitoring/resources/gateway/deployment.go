package gateway

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"time"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/query"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/queryfrontend"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/router"
	monitoringgateway "github.com/kubesphere/whizard/pkg/monitoring-gateway"
	"github.com/kubesphere/whizard/pkg/util"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (g *Gateway) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: g.meta(g.name())}

	if g.gateway == nil {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: g.gateway.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: g.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: g.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: g.gateway.Spec.NodeSelector,
				Tolerations:  g.gateway.Spec.Tolerations,
				Affinity:     g.gateway.Spec.Affinity,
			},
		},
	}

	if g.gateway.Spec.ImagePullSecrets != nil && len(g.gateway.Spec.ImagePullSecrets) > 0 {
		d.Spec.Template.Spec.ImagePullSecrets = g.gateway.Spec.ImagePullSecrets
	}

	var container = corev1.Container{
		Name:      "gateway",
		Image:     g.gateway.Spec.Image,
		Args:      []string{},
		Resources: g.gateway.Spec.Resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "http",
				ContainerPort: 9090,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		// LivenessProbe: &corev1.Probe{
		// 	FailureThreshold: 4,
		// 	PeriodSeconds:    30,
		// 	ProbeHandler: corev1.ProbeHandler{
		// 		HTTPGet: &corev1.HTTPGetAction{
		// 			Scheme: "HTTP",
		// 			Path:   "/-/healthy",
		// 			Port:   intstr.FromString("http"),
		// 		},
		// 	},
		// },
		// ReadinessProbe: &corev1.Probe{
		// 	FailureThreshold: 20,
		// 	PeriodSeconds:    5,
		// 	ProbeHandler: corev1.ProbeHandler{
		// 		HTTPGet: &corev1.HTTPGetAction{
		// 			Scheme: "HTTP",
		// 			Path:   "/-/ready",
		// 			Port:   intstr.FromString("http"),
		// 		},
		// 	},
		// },
	}

	if g.gateway.Spec.ServerCertificate != "" {
		serverCertVol := corev1.Volume{
			Name: "server-certificate",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: g.gateway.Spec.ServerCertificate,
				},
			},
		}
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, serverCertVol)
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      serverCertVol.Name,
			MountPath: filepath.Join(secretsDir, g.gateway.Spec.ServerCertificate),
		})
		container.Args = append(container.Args, "--server-tls-key="+filepath.Join(secretsDir, g.gateway.Spec.ServerCertificate, "tls.key"))
		container.Args = append(container.Args, "--server-tls-crt="+filepath.Join(secretsDir, g.gateway.Spec.ServerCertificate, "tls.crt"))
	}

	if g.gateway.Spec.ClientCACertificate != "" {
		clientCaCertVol := corev1.Volume{
			Name: "client-ca-certificate",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: g.gateway.Spec.ClientCACertificate,
				},
			},
		}
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, clientCaCertVol)
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      clientCaCertVol.Name,
			MountPath: filepath.Join(secretsDir, g.gateway.Spec.ClientCACertificate),
		})
		container.Args = append(container.Args, "--server-tls-client-ca="+filepath.Join(secretsDir, g.gateway.Spec.ClientCACertificate, "tls.crt"))
	}

	if g.gateway.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+g.gateway.Spec.LogLevel)
	}
	if g.gateway.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+g.gateway.Spec.LogFormat)
	}

	if g.Service.Spec.TenantHeader != "" {
		container.Args = append(container.Args, "--tenant.header="+g.Service.Spec.TenantHeader)
	}
	if g.Service.Spec.TenantLabelName != "" {
		container.Args = append(container.Args, "--tenant.label-name="+g.Service.Spec.TenantLabelName)
	}

	queryFrontendAddr, err := g.queryfrontendAddress()
	if err != nil {
		return nil, "", err
	}
	queryAddr, err := g.queryAddress()
	if err != nil {
		return nil, "", err
	}
	if g.Service != nil && g.Service.Spec.RemoteQuery != nil {
		// If there is remote query config in service,
		// Gateway will query metrics from QueryFrontend (which is put in front of remote-query),
		// while query rules from Query (which aggregates rules from all rulers).
		if queryFrontendAddr == "" {
			return nil, "", fmt.Errorf("no query frontend exist for service %s/%s", g.Service.Name, g.Service.Namespace)
		}
		queryFrontendUrl, err := url.Parse(queryFrontendAddr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid query frontend address: %s", queryFrontendAddr)
		}
		container.Args = append(container.Args, fmt.Sprintf("--query.address=%s", queryFrontendAddr))
		if queryFrontendUrl.Scheme == "https" {
			cfg := config{TLSConfig: &config_util.TLSConfig{InsecureSkipVerify: true}}
			buff, _ := yaml.Marshal(cfg)
			container.Args = append(container.Args, fmt.Sprintf("--query.config=%s", buff))
		}
		if queryAddr == "" {
			return nil, "", fmt.Errorf("no query exist for service %s/%s", g.Service.Name, g.Service.Namespace)
		}
		queryUrl, err := url.Parse(queryAddr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid query address: %s", queryAddr)
		}
		container.Args = append(container.Args, fmt.Sprintf("--rules-query.address=%s", queryAddr))
		if queryUrl.Scheme == "https" {
			cfg := config{TLSConfig: &config_util.TLSConfig{InsecureSkipVerify: true}}
			buff, _ := yaml.Marshal(cfg)
			container.Args = append(container.Args, fmt.Sprintf("--rules-query.config=%s", buff))
		}
	} else {
		// If there is no remote query config, the Gateway will preferentially query all from QueryFrontend
		var addr = queryFrontendAddr
		if addr == "" {
			addr = queryAddr
		}
		if addr == "" {
			return nil, "", fmt.Errorf("no query frontend and query exist for service %s/%s", g.Service.Name, g.Service.Namespace)
		}
		queryUrl, err := url.Parse(addr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid query address: %s", addr)
		}
		container.Args = append(container.Args, fmt.Sprintf("--query.address=%s", addr))
		if queryUrl.Scheme == "https" {
			cfg := config{TLSConfig: &config_util.TLSConfig{InsecureSkipVerify: true}}
			buff, _ := yaml.Marshal(cfg)
			container.Args = append(container.Args, fmt.Sprintf("--query.config=%s", buff))
		}
	}

	var rwsCfg []*monitoringgateway.RemoteWriteConfig
	// write to router
	routerAddr, err := g.remoteWriteAddress()
	if err != nil {
		return nil, "", err
	}
	url, err := url.Parse(routerAddr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid router address: %s", queryAddr)
	}
	url.Path = path.Join(url.Path, "/api/v1/receive")
	rwRouter := &monitoringgateway.RemoteWriteConfig{URL: &config_util.URL{URL: url}}
	if url.Scheme == "https" {
		rwRouter.TLSConfig = config_util.TLSConfig{InsecureSkipVerify: true}
	}
	rwsCfg = append(rwsCfg, rwRouter)
	// write to configured remote-writes targets
	if g.Service != nil {
		for _, rw := range g.Service.Spec.RemoteWrites {
			url, err := url.Parse(rw.URL)
			if err != nil {
				return nil, "", fmt.Errorf("invalid remote write url: %s", rw.URL)
			}
			rwCfg := &monitoringgateway.RemoteWriteConfig{
				Name:    rw.Name,
				URL:     &config_util.URL{URL: url},
				Headers: rw.Headers,
			}
			if rw.RemoteTimeout != "" {
				timeout, err := time.ParseDuration(string(rw.RemoteTimeout))
				if err != nil {
					return nil, "", fmt.Errorf("invalid remoteTimeout: %s", rw.RemoteTimeout)
				}
				rwCfg.RemoteTimeout = model.Duration(timeout)
			}
			if url.Scheme == "https" {
				rwCfg.TLSConfig = config_util.TLSConfig{InsecureSkipVerify: true}
			}
			rwsCfg = append(rwsCfg, rwCfg)
		}
	}
	// add remote-writes.config flag to gateway
	buff, err := yaml.Marshal(rwsCfg)
	if err != nil {
		return nil, "", err
	}
	container.Args = append(container.Args, fmt.Sprintf("--remote-writes.config=%s", buff))

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)

	return d, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(g.gateway, d, g.Scheme)
}

func (g *Gateway) queryfrontendAddress() (string, error) {
	queryFrontendList := &v1alpha1.QueryFrontendList{}
	if err := g.Client.List(g.Context, queryFrontendList, client.MatchingLabels(util.ManagedLabelBySameService(g.gateway))); err != nil {
		return "", err
	}

	if len(queryFrontendList.Items) > 0 {
		if len(queryFrontendList.Items) > 1 {
			return "", fmt.Errorf("more than one query frontend defined for service %s/%s", g.Service.Name, g.Service.Namespace)
		}

		q := queryFrontendList.Items[0]
		r, err := queryfrontend.New(g.BaseReconciler, &q)
		if err != nil {
			return "", err
		}
		if q.Spec.HTTPServerTLSConfig != nil {
			return r.HttpsAddr(), nil
		}

		return r.HttpAddr(), nil
	}

	return "", nil
}

func (g *Gateway) queryAddress() (string, error) {

	queryList := &v1alpha1.QueryList{}
	if err := g.Client.List(g.Context, queryList, client.MatchingLabels(util.ManagedLabelBySameService(g.gateway))); err != nil {
		return "", err
	}

	if len(queryList.Items) > 0 {
		if len(queryList.Items) > 1 {
			return "", fmt.Errorf("more than one query defined for service %s/%s", g.Service.Name, g.Service.Namespace)
		}

		q := queryList.Items[0]
		r, err := query.New(g.BaseReconciler, &q, nil)
		if err != nil {
			return "", err
		}
		if q.Spec.HTTPServerTLSConfig != nil {
			return r.HttpsAddr(), nil
		}

		return r.HttpAddr(), nil
	}

	return "", nil
}

func (g *Gateway) remoteWriteAddress() (string, error) {
	routerList := &v1alpha1.RouterList{}
	if err := g.Client.List(g.Context, routerList, client.MatchingLabels(util.ManagedLabelBySameService(g.gateway))); err != nil {
		return "", err
	}

	if len(routerList.Items) > 0 {
		if len(routerList.Items) > 1 {
			return "", fmt.Errorf("more than one router defined for service %s/%s", g.Service.Name, g.Service.Namespace)
		}

		o := routerList.Items[0]
		r, err := router.New(g.BaseReconciler, &o, nil)
		if err != nil {
			return "", err
		}
		if o.Spec.HTTPServerTLSConfig != nil {
			return r.RemoteWriteHTTPSAddr(), nil
		}

		return r.RemoteWriteAddr(), nil
	}

	return "", fmt.Errorf("no router defined for service %s/%s", g.Service.Name, g.Service.Namespace)
}

type config struct {
	TLSConfig *config_util.TLSConfig `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
}
