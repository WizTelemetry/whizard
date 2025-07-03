package ingester

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/prometheus-operator/prometheus-operator/pkg/k8sutil"
	"github.com/prometheus/common/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
	"github.com/WhizardTelemetry/whizard/pkg/util"
)

const (
	initContainerName = "cleanup"
)

var (
	// repeatableArgs is the args that can be set repeatedly.
	// An error will occur if a non-repeatable arg is set repeatedly.
	repeatableArgs = []string{
		"--label",
	}
	// unsupportedArgs is the args that are not allowed to be set by the user.
	unsupportedArgs = []string{
		"--receive.hashrings",
		"--receive.hashrings-file",
		"--http-address",
		"--grpc-address",
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
				NodeSelector:    r.ingester.Spec.NodeSelector,
				Tolerations:     r.ingester.Spec.Tolerations,
				Affinity:        r.ingester.Spec.Affinity,
				SecurityContext: r.ingester.Spec.SecurityContext,
			},
		},
	}

	// To make sure there is enough time to upload the block when ingester is terminated.
	terminationGracePeriodSeconds := int64(time.Hour)
	sts.Spec.Template.Spec.TerminationGracePeriodSeconds = &terminationGracePeriodSeconds

	if r.ingester.Spec.ImagePullSecrets != nil && len(r.ingester.Spec.ImagePullSecrets) > 0 {
		sts.Spec.Template.Spec.ImagePullSecrets = r.ingester.Spec.ImagePullSecrets
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
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          constants.CapNProtoPortName,
				ContainerPort: constants.CapNProtoPort,
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

	r.AddTSDBVolume(sts, &container, r.ingester.Spec.DataVolume)

	var storageConfig []byte
	if r.ingester.Labels != nil {
		if namespacedName := r.ingester.Labels[constants.StorageLabelKey]; namespacedName != "" {
			var err error
			storageConfig, err = r.GetStorageConfig(namespacedName)
			if err != nil {
				return nil, "", err
			}
		}
	}

	if r.ingester.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.ingester.Spec.LogLevel)
	}
	if r.ingester.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.ingester.Spec.LogFormat)
	}

	if r.ingester.Spec.OtlpEnableTargetInfo != nil && !*r.ingester.Spec.OtlpEnableTargetInfo {
		container.Args = append(container.Args, "--no-receive.otlp-enable-target-info")
	}

	for _, attr := range r.ingester.Spec.OtlpResourceAttributes {
		container.Args = append(container.Args, "--receive.otlp-promote-resource-attributes="+attr)
	}

	container.Args = append(container.Args, fmt.Sprintf("--label=%s=\"$(POD_NAME)\"", constants.ReceiveReplicaLabelName))
	container.Args = append(container.Args, fmt.Sprintf("--tsdb.path=%s", constants.StorageDir))
	container.Args = append(container.Args, fmt.Sprintf("--receive.local-endpoint=$(POD_NAME).%s:%d", r.name(constants.ServiceNameSuffix), constants.GRPCPort))
	if r.ingester.Spec.LocalTsdbRetention != "" {
		container.Args = append(container.Args, "--tsdb.retention="+r.ingester.Spec.LocalTsdbRetention)
	}
	if storageConfig != nil {
		container.Args = append(container.Args, "--objstore.config="+string(storageConfig))
		volumes, volumeMounts, err := r.VolumesAndVolumeMountsForStorage(r.ingester.Labels[constants.StorageLabelKey])
		if err != nil {
			return nil, "", err
		}
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, volumes...)
		container.VolumeMounts = append(container.VolumeMounts, volumeMounts...)
	} else {
		// set tsdb.max-block-duration by localTsdbRetention to enable block compact when using only local storage
		// https://prometheus.io/docs/prometheus/latest/storage/#compaction
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

	if r.Service.Spec.TenantHeader != "" {
		container.Args = append(container.Args, "--receive.tenant-header="+r.Service.Spec.TenantHeader)
	}
	if r.Service.Spec.TenantLabelName != "" {
		container.Args = append(container.Args, "--receive.tenant-label-name="+r.Service.Spec.TenantLabelName)
	}
	if r.Service.Spec.DefaultTenantId != "" {
		container.Args = append(container.Args, "--receive.default-tenant-id="+r.Service.Spec.DefaultTenantId)
	}

	for _, flag := range r.ingester.Spec.Flags {
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

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container)
	sts.Spec.Template.Spec.InitContainers = r.generateInitContainer(getTSDBVolumeMount(container))

	if len(r.ingester.Spec.Containers.Raw) > 0 {
		var err error
		r.ingester.Spec.EmbeddedContainers, err = util.DecodeRawToContainers(r.ingester.Spec.Containers)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode containers: %w", err)
		}
		containers, err := k8sutil.MergePatchContainers(sts.Spec.Template.Spec.Containers, r.ingester.Spec.EmbeddedContainers)
		if err != nil {
			return nil, "", fmt.Errorf("failed to merge containers spec: %w", err)
		}
		sts.Spec.Template.Spec.Containers = containers
	}

	return sts, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.ingester, sts, r.Scheme)
}

var cleanupScript = `
#!/bin/bash

echo [$(date "+%Y-%m-%d %H:%M:%S")] begin to cleanup block
echo TSDB path: ${0}
echo tenants: ${1}

files=$(ls -d $0/*)
tenants=(${1//,/ })

for f in ${files[@]}
do
  if test -d $f; then
    name=$(basename $f)
    if [[ ! "${tenants[@]}" =~ "$name" ]]; then
      echo [$(date "+%Y-%m-%d %H:%M:%S")] tenant $name does not exist, delete data directory $f
      rm -rf $f
    fi
  fi
done

echo [$(date "+%Y-%m-%d %H:%M:%S")] cleanup block end
`

func (r *Ingester) generateInitContainer(tsdbVolumeMount *corev1.VolumeMount) []corev1.Container {

	// The tsdbVolumeMount is nil means ingester uses empty dir as the storage of TSDB, no need to cleanup.
	if (r.Service.Spec.IngesterTemplateSpec.DisableTSDBCleanup != nil && *r.Service.Spec.IngesterTemplateSpec.DisableTSDBCleanup) ||
		tsdbVolumeMount == nil {
		return nil
	}

	// Soft tenant ingesters should not clean up data.
	if v, ok := r.ingester.Labels[constants.SoftTenantLabelKey]; ok && v == "true" {
		return nil
	}

	var tenants []string
	for _, tenant := range r.ingester.Status.Tenants {
		tenants = append(tenants, tenant.Name)
	}
	sort.Strings(tenants)

	return []corev1.Container{
		{
			Name:  initContainerName,
			Image: r.ingester.Spec.IngesterTSDBCleanUp.Image,
			Command: []string{
				"bash",
				"-c",
				cleanupScript,
			},
			Args: []string{
				constants.StorageDir,
				strings.Join(append(tenants, r.Service.Spec.DefaultTenantId), ","),
			},
			Resources:    r.ingester.Spec.Resources,
			VolumeMounts: []corev1.VolumeMount{*tsdbVolumeMount},
		},
	}

}

func getTSDBVolumeMount(container corev1.Container) *corev1.VolumeMount {

	for _, item := range container.VolumeMounts {
		if item.Name == constants.TSDBVolumeName {
			return &item
		}
	}

	return nil
}
