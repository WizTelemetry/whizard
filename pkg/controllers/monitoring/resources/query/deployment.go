package query

import (
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
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

	var ingestorList monitoringv1alpha1.ThanosReceiveIngestorList
	if err := q.Client.List(q.Context, &ingestorList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("thanosreceiveingestorlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range ingestorList.Items {
		ingestorSvcName := resources.QualifiedName(resources.AppNameThanosRuler, item.Name, resources.ServiceNameSuffixOperated)
		endpoint := fmt.Sprintf("%s.%s.svc:%d", ingestorSvcName, item.Namespace, resources.ThanosGRPCPort)
		queryContainer.Args = append(queryContainer.Args, "--endpoint="+endpoint)
	}

	var storeList monitoringv1alpha1.StoreList
	if err := q.Client.List(q.Context, &storeList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("storelist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range storeList.Items {
		storeSvcName := resources.QualifiedName(resources.AppNameThanosStoreGateway, item.Name, resources.ServiceNameSuffixOperated)
		endpoint := fmt.Sprintf("%s.%s.svc:%d", storeSvcName, item.Namespace, resources.ThanosGRPCPort)
		queryContainer.Args = append(queryContainer.Args, "--endpoint="+endpoint)
	}

	var rulerList monitoringv1alpha1.ThanosRulerList
	if err := q.Client.List(q.Context, &rulerList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("thanosrulerlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range rulerList.Items {
		rulerSvcName := resources.QualifiedName(resources.AppNameThanosRuler, item.Name, resources.ServiceNameSuffixOperated)
		endpoint := fmt.Sprintf("%s.%s.svc:%d", rulerSvcName, item.Namespace, resources.ThanosGRPCPort)
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
