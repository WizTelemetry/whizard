package query_frontend

import (
	"fmt"
	"path/filepath"

	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/query"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (q *QueryFrontend) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: q.meta(q.name())}

	if q.queryFrontend == nil {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: q.queryFrontend.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: q.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: q.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: q.queryFrontend.NodeSelector,
				Tolerations:  q.queryFrontend.Tolerations,
				Affinity:     q.queryFrontend.Affinity,
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

	var container = corev1.Container{
		Name:      "query-frontend",
		Image:     q.queryFrontend.Image,
		Args:      []string{"query-frontend"},
		Resources: q.queryFrontend.Resources,
		Ports: []corev1.ContainerPort{
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          resources.ThanosHTTPPortName,
				ContainerPort: resources.ThanosHTTPPort,
			},
		},
		LivenessProbe:  resources.ThanosDefaultLivenessProbe(),
		ReadinessProbe: resources.ThanosDefaultReadinessProbe(),
		VolumeMounts: []corev1.VolumeMount{{
			Name:      cacheConfigVol.Name,
			MountPath: configDir,
			ReadOnly:  true,
		}},
	}

	query := query.New(q.ServiceBaseReconciler)
	container.Args = append(container.Args, "--query-frontend.downstream-url="+query.HttpAddr())
	container.Args = append(container.Args, "--labels.response-cache-config-file="+filepath.Join(configDir, cacheConfigFile))
	container.Args = append(container.Args, "--query-range.response-cache-config-file="+filepath.Join(configDir, cacheConfigFile))
	for param, value := range q.queryFrontend.Flags {
		container.Args = append(container.Args, fmt.Sprintf("--%s=%s", param, value))
	}

	if q.queryFrontend.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+q.queryFrontend.LogLevel)
	}
	if q.queryFrontend.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+q.queryFrontend.LogFormat)
	}

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)
	d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, cacheConfigVol)

	return d, resources.OperationCreateOrUpdate, nil
}
