package query

import (
	"fmt"
	"path/filepath"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/ingester"
	"github.com/kubesphere/whizard/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (q *Query) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: q.meta(q.name())}

	if q.query == nil {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: q.query.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: q.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: q.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: q.query.NodeSelector,
				Tolerations:  q.query.Tolerations,
				Affinity:     q.query.Affinity,
			},
		},
	}

	proxyConfigVol := corev1.Volume{
		Name: "proxy-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: q.name("proxy-config"),
				},
			},
		},
	}
	d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, proxyConfigVol)
	storesConfigVol := corev1.Volume{
		Name: "stores-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: q.name("stores-config"),
				},
			},
		},
	}
	d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, storesConfigVol)

	var queryContainer = corev1.Container{
		Name:      "query",
		Image:     q.query.Image,
		Args:      []string{"query"},
		Resources: q.query.Resources,
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
		VolumeMounts: []corev1.VolumeMount{{
			Name:      storesConfigVol.Name,
			MountPath: configDir,
			ReadOnly:  true,
		}},
	}

	if q.query.LogLevel != "" {
		queryContainer.Args = append(queryContainer.Args, "--log.level="+q.query.LogLevel)
	}
	if q.query.LogFormat != "" {
		queryContainer.Args = append(queryContainer.Args, "--log.format="+q.query.LogFormat)
	}
	queryContainer.Args = append(queryContainer.Args, "--store.sd-files="+filepath.Join(configDir, storesFile))
	for _, labelName := range q.query.ReplicaLabelNames {
		queryContainer.Args = append(queryContainer.Args, "--query.replica-label="+labelName)
	}

	var ingesterList monitoringv1alpha1.IngesterList
	if err := q.Client.List(q.Context, &ingesterList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("ingesterlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range ingesterList.Items {
		ingesterInstance := ingester.New(q.BaseReconciler, &item)
		for _, endpoint := range ingesterInstance.GrpcAddrs() {
			queryContainer.Args = append(queryContainer.Args, "--endpoint="+endpoint)
		}
	}

	var storeList monitoringv1alpha1.StoreList
	if err := q.Client.List(q.Context, &storeList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("storelist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range storeList.Items {
		storeSvcName := util.Join("-", item.Name, constants.ServiceNameSuffix)
		endpoint := fmt.Sprintf("%s.%s.svc:%d", storeSvcName, item.Namespace, constants.GRPCPort)
		queryContainer.Args = append(queryContainer.Args, "--endpoint="+endpoint)
	}

	var rulerList monitoringv1alpha1.RulerList
	if err := q.Client.List(q.Context, &rulerList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("rulerlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range rulerList.Items {
		rulerSvcName := resources.QualifiedName(constants.AppNameRuler, item.Name, constants.ServiceNameSuffix)
		endpoint := fmt.Sprintf("%s.%s.svc:%d", rulerSvcName, item.Namespace, constants.GRPCPort)
		queryContainer.Args = append(queryContainer.Args, "--endpoint="+endpoint)
	}

	for name, value := range q.query.Flags {
		queryContainer.Args = append(queryContainer.Args, fmt.Sprintf("--%s=%s", name, value))
	}

	var envoyContainer = corev1.Container{
		Name:  "proxy",
		Image: q.query.Envoy.Image,
		Args: []string{
			"-c",
			filepath.Join(envoyConfigDir, envoyConfigFile),
			// "-l",
			// "debug",
		},
		Resources: q.query.Envoy.Resources,
		VolumeMounts: []corev1.VolumeMount{{
			Name:      proxyConfigVol.Name,
			MountPath: envoyConfigDir,
			ReadOnly:  true,
		}},
	}

	for _, store := range q.query.Stores {
		if store.CASecret == nil {
			continue
		}
		secretVol := corev1.Volume{
			Name: "secret-" + store.CASecret.Name,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: store.CASecret.Name,
				},
			},
		}
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, secretVol)
		envoyContainer.VolumeMounts = append(envoyContainer.VolumeMounts, corev1.VolumeMount{
			Name:      secretVol.Name,
			ReadOnly:  true,
			SubPath:   store.CASecret.Key,
			MountPath: filepath.Join(envoySecretsDir, store.CASecret.Name, store.CASecret.Key),
		})
	}

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, queryContainer)
	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, envoyContainer)

	return d, resources.OperationCreateOrUpdate, nil
}
