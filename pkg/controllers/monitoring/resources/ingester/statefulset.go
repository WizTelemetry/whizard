package ingester

import (
	"fmt"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
	"github.com/prometheus/common/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	// sliceArgs is the args that can be set repeatedly.
	// An error will occur if a non-slice arg is set repeatedly.
	sliceArgs = []string{
		"--label",
	}
	// unsupportedArgs is the args that are not allowed to be set by the user.
	unsupportedArgs = []string{
		"--receive.hashrings",
		"--receive.hashrings-file",
	}
)

func (r *Ingester) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name())}

	sts.Spec = appsv1.StatefulSetSpec{
		Replicas: r.ingester.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		ServiceName: r.name(constants.ServiceNameSuffix),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.ingester.Spec.NodeSelector,
				Tolerations:  r.ingester.Spec.Tolerations,
				Affinity:     r.ingester.Spec.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "receive",
		Image:     r.ingester.Spec.Image,
		Args:      []string{"receive"},
		Resources: r.ingester.Spec.Resources,
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
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          constants.RemoteWritePortName,
				ContainerPort: constants.RemoteWritePort,
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

	resources.AddTSDBVolume(sts, &container, r.ingester.Spec.DataVolume)

	var storageConfig []byte
	if r.ingester.Spec.Storage != nil {
		var err error
		storageConfig, err = resources.GetStorageConfig(r.Context, r.Client,
			fmt.Sprintf("%s.%s", r.ingester.Spec.Storage.Namespace, r.ingester.Spec.Storage.Name))
		if err != nil {
			return nil, "", err
		}
	}

	if r.ingester.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.ingester.Spec.LogLevel)
	}
	if r.ingester.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.ingester.Spec.LogFormat)
	}
	container.Args = append(container.Args, fmt.Sprintf("--label=%s=\"$(POD_NAME)\"", constants.ReceiveReplicaLabelName))
	container.Args = append(container.Args, fmt.Sprintf("--tsdb.path=%s", constants.StorageDir))
	container.Args = append(container.Args, fmt.Sprintf("--receive.local-endpoint=$(POD_NAME).%s:%d", r.name(constants.ServiceNameSuffix), constants.GRPCPort))
	if r.ingester.Spec.LocalTsdbRetention != "" {
		container.Args = append(container.Args, "--tsdb.retention="+r.ingester.Spec.LocalTsdbRetention)
	}
	if storageConfig != nil {
		container.Args = append(container.Args, "--objstore.config="+string(storageConfig))
		volumes, volumeMounts, err := resources.VolumesAndVolumeMountsForStorage(r.Context, r.Client, r.ingester.Labels[constants.StorageLabelKey])
		if err != nil {
			return nil, "", err
		}
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, volumes...)
		container.VolumeMounts = append(container.VolumeMounts, volumeMounts...)
	} else {
		// set tsdb.max-block-duration by localTsdbRetention to enable block compact when using only local storage
		maxBlockDuration, err := model.ParseDuration("31d")
		if err != nil {
			return nil, resources.OperationCreateOrUpdate, err
		}
		retention := r.ingester.Spec.LocalTsdbRetention
		if retention == "" {
			retention = "15d"
		}
		retentionDuration, err := model.ParseDuration(retention)
		if err != nil {
			return nil, resources.OperationCreateOrUpdate, err
		}
		if retentionDuration != 0 && retentionDuration/10 < maxBlockDuration {
			maxBlockDuration = retentionDuration / 10
		}

		container.Args = append(container.Args, "--tsdb.max-block-duration="+maxBlockDuration.String())
	}

	namespacedName := monitoringv1alpha1.ServiceNamespacedName(r.ingester)

	if namespacedName != nil {
		var service monitoringv1alpha1.Service
		if err := r.Client.Get(r.Context, *namespacedName, &service); err != nil {
			if !apierrors.IsNotFound(err) {
				return nil, resources.OperationCreateOrUpdate, err
			}
		} else {
			if service.Spec.TenantHeader != "" {
				container.Args = append(container.Args, "--receive.tenant-header="+service.Spec.TenantHeader)
			}
			if service.Spec.TenantLabelName != "" {
				container.Args = append(container.Args, "--receive.tenant-label-name="+service.Spec.TenantLabelName)
			}
			if service.Spec.DefaultTenantId != "" {
				container.Args = append(container.Args, "--receive.default-tenant-id="+service.Spec.DefaultTenantId)
			}
		}
	}

	for _, flag := range r.ingester.Spec.Flags {
		arg := util.GetArgName(flag)
		if util.Contains(unsupportedArgs, arg) {
			continue
		}

		if util.Contains(sliceArgs, arg) {
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

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container)

	return sts, resources.OperationCreateOrUpdate, nil
}
