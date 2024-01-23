package resources

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/prometheus-operator/prometheus-operator/pkg/k8sutil"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/util"
)

type Operation string

const (
	OperationCreateOrUpdate Operation = "CreateOrUpdate"
	OperationDelete         Operation = "Delete"
)

type Resource func() (runtime.Object, Operation, error)

type BaseReconciler struct {
	Context context.Context
	Client  client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Service *v1alpha1.Service
}

func (r *BaseReconciler) SetService(o runtime.Object) error {
	accessor, err := meta.Accessor(o)
	if err != nil {
		return err
	}

	lables := accessor.GetLabels()
	if lables == nil || lables[constants.ServiceLabelKey] == "" {
		return fmt.Errorf("service missing")
	}

	namespacedName := strings.Split(lables[constants.ServiceLabelKey], ".")
	s := &v1alpha1.Service{}
	if err := r.Client.Get(r.Context, client.ObjectKey{Name: namespacedName[1], Namespace: namespacedName[0]}, s); err != nil {
		return err
	}

	if s.Spec.TenantHeader == "" {
		s.Spec.TenantHeader = constants.DefaultTenantHeader
	}
	if s.Spec.TenantLabelName == "" {
		s.Spec.TenantLabelName = constants.DefaultTenantLabelName
	}
	if s.Spec.DefaultTenantId == "" {
		s.Spec.DefaultTenantId = constants.DefaultTenantId
	}

	r.Service = s

	return nil
}

func (r *BaseReconciler) GetStorage(storage string) string {
	if storage == "" {
		storage = constants.DefaultStorage
	}

	if storage == constants.DefaultStorage {
		if r.Service != nil && r.Service.Spec.Storage != nil {
			storage = fmt.Sprintf("%s.%s", r.Service.Spec.Storage.Namespace, r.Service.Spec.Storage.Name)
		} else {
			storage = constants.LocalStorage
		}
	}

	return storage
}

func (r *BaseReconciler) BaseLabels() map[string]string {
	return map[string]string{}
}

func (r *BaseReconciler) ReconcileResources(resources []Resource) error {
	for _, resource := range resources {
		obj, operation, err := resource()
		if err != nil {
			return err
		}

		switch operation {
		case OperationDelete:
			err := r.Client.Delete(r.Context, obj.(client.Object))
			if err != nil && !apierrors.IsNotFound(err) {
				return err
			}
		case OperationCreateOrUpdate:
			if err = func() error {
				switch desired := obj.(type) {
				case *appsv1.Deployment:
					return util.CreateOrUpdateDeployment(r.Context, r.Client, desired)
				case *corev1.Service:
					return util.CreateOrUpdateService(r.Context, r.Client, desired)
				case *corev1.ServiceAccount:
					return util.CreateOrUpdateServiceAccount(r.Context, r.Client, desired)
				case *appsv1.StatefulSet:
					err = util.CreateOrUpdateStatefulSet(r.Context, r.Client, desired)
					sErr, ok := err.(*apierrors.StatusError)
					if ok && sErr.ErrStatus.Code == 422 && sErr.ErrStatus.Reason == metav1.StatusReasonInvalid {
						// Gather only reason for failed update
						failMsg := make([]string, len(sErr.ErrStatus.Details.Causes))
						for i, cause := range sErr.ErrStatus.Details.Causes {
							failMsg[i] = cause.Message
						}
						r.Log.Info("recreating StatefulSet because the update operation wasn't possible", "reason", strings.Join(failMsg, ", "))
						if err = r.Client.Delete(r.Context, obj.(client.Object), client.PropagationPolicy(metav1.DeletePropagationForeground)); err != nil {
							return errors.Wrap(err, "failed to delete StatefulSet to avoid forbidden action")
						}
						return nil
					}
					if err != nil {
						return errors.Wrap(err, "updating StatefulSet failed")
					}
					return nil
				default:
					return util.CreateOrUpdate(r.Context, r.Client, desired.(client.Object))
				}
			}(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *BaseReconciler) QualifiedName(appName, instanceName string, suffix ...string) string {
	name := appName + "-" + instanceName
	if len(suffix) > 0 {
		name += "-" + strings.Join(suffix, "-")
	}
	return name
}

func (r *BaseReconciler) DefaultLivenessProbe() *corev1.Probe {
	return &corev1.Probe{
		FailureThreshold: 4,
		PeriodSeconds:    30,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTP",
				Path:   "/-/healthy",
				Port:   intstr.FromString(constants.HTTPPortName),
			},
		},
	}
}

func (r *BaseReconciler) DefaultReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		FailureThreshold: 20,
		PeriodSeconds:    5,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTP",
				Path:   "/-/ready",
				Port:   intstr.FromString(constants.HTTPPortName),
			},
		},
	}
}

func (r *BaseReconciler) DefaultLivenessProbeWithTLS() *corev1.Probe {
	return &corev1.Probe{
		FailureThreshold: 4,
		PeriodSeconds:    30,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTPS",
				Path:   "/-/healthy",
				Port:   intstr.FromString(constants.HTTPPortName),
			},
		},
	}
}

func (r *BaseReconciler) DefaultReadinessProbeWithTLS() *corev1.Probe {
	return &corev1.Probe{
		FailureThreshold: 20,
		PeriodSeconds:    5,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTPS",
				Path:   "/-/ready",
				Port:   intstr.FromString(constants.HTTPPortName),
			},
		},
	}
}

func (r *BaseReconciler) AddTSDBVolume(sts *appsv1.StatefulSet, container *corev1.Container, dataVolume *v1alpha1.KubernetesVolume) {
	var volumeName string
	if dataVolume == nil || // If dataVolume is not specified, default to a new EmptyDirVolumeSource
		(dataVolume.PersistentVolumeClaim == nil && dataVolume.EmptyDir == nil) {
		volumeName = constants.TSDBVolumeName
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: constants.TSDBVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	} else if dataVolume.PersistentVolumeClaim != nil {
		pvc := *dataVolume.PersistentVolumeClaim
		if pvc.Name == "" {
			pvc.Name = constants.TSDBVolumeName
		}
		volumeName = pvc.Name
		if pvc.Spec.AccessModes == nil {
			pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		}

		sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{pvc}
		policy := dataVolume.PersistentVolumeClaimRetentionPolicy
		if policy != nil {
			if policy.WhenDeleted != appsv1.RetainPersistentVolumeClaimRetentionPolicyType &&
				policy.WhenDeleted != appsv1.DeletePersistentVolumeClaimRetentionPolicyType {
				policy.WhenDeleted = appsv1.RetainPersistentVolumeClaimRetentionPolicyType
			}
			if policy.WhenScaled != appsv1.RetainPersistentVolumeClaimRetentionPolicyType &&
				policy.WhenScaled != appsv1.DeletePersistentVolumeClaimRetentionPolicyType {
				policy.WhenScaled = appsv1.RetainPersistentVolumeClaimRetentionPolicyType
			}
			sts.Spec.PersistentVolumeClaimRetentionPolicy = policy
		}
	} else if dataVolume.EmptyDir != nil {
		volumeName = constants.TSDBVolumeName
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: constants.TSDBVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: dataVolume.EmptyDir,
			},
		})
	}

	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      volumeName,
		MountPath: constants.StorageDir,
	})
}

func (r *BaseReconciler) GetTenantHash(selector map[string]string) (string, error) {
	tenantList := &v1alpha1.TenantList{}
	err := r.Client.List(r.Context, tenantList, client.MatchingLabels(selector))
	if err != nil {
		return "", err
	}

	var tenants []string
	for _, item := range tenantList.Items {
		if item.DeletionTimestamp != nil || !item.DeletionTimestamp.IsZero() {
			continue
		}

		tenants = append(tenants, item.Name)
	}
	sort.Strings(tenants)

	hash := md5.New()
	for _, tenant := range tenants {
		hash.Write([]byte(tenant))
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (r *BaseReconciler) GetStorageHash(namespaceName string) (string, error) {
	storage := &v1alpha1.Storage{}
	array := strings.Split(namespaceName, ".")
	if err := r.Client.Get(r.Context, client.ObjectKey{Name: array[1], Namespace: array[0]}, storage); err != nil {
		return "", err
	}

	storageConfig, err := r.GetStorageConfig(namespaceName)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	hash.Write(storageConfig)

	if storage.Spec.S3 != nil {
		tls := storage.Spec.S3.HTTPConfig.TLSConfig
		if bs, err := r.getValueFromSecret(tls.CA, storage.Namespace); err != nil {
			return "", err
		} else {
			hash.Write(bs)
		}

		if bs, err := r.getValueFromSecret(tls.Key, storage.Namespace); err != nil {
			return "", err
		} else {
			hash.Write(bs)
		}

		if bs, err := r.getValueFromSecret(tls.Cert, storage.Namespace); err != nil {
			return "", err
		} else {
			hash.Write(bs)
		}
	}

	hashStr := hex.EncodeToString(hash.Sum(nil))
	return hashStr, nil
}

type bucketConfig struct {
	Type   string      `yaml:"type"`
	Config interface{} `yaml:"config"`
}

func (r *BaseReconciler) GetStorageConfig(namespaceName string) ([]byte, error) {
	if namespaceName == constants.LocalStorage {
		return nil, nil
	}

	storage := &v1alpha1.Storage{}
	array := strings.Split(namespaceName, ".")
	if err := r.Client.Get(r.Context, client.ObjectKey{Name: array[1], Namespace: array[0]}, storage); err != nil {
		return nil, err
	}

	if storage.Spec.S3 != nil {
		b := &bucketConfig{
			constants.StorageProviderS3,
			*storage.Spec.S3,
		}

		root := &yaml.Node{}
		bs, err := yaml.Marshal(b)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(bs, root); err != nil {
			return nil, err
		}

		bs, err = r.getValueFromSecret(storage.Spec.S3.AccessKey, storage.Namespace)
		if err != nil {
			return nil, err
		}
		if n := findYamlNodeByKey(root, "access_key"); n != nil {
			n.SetString(string(bs))
		}

		bs, err = r.getValueFromSecret(storage.Spec.S3.SecretKey, storage.Namespace)
		if err != nil {
			return nil, err
		}
		if n := findYamlNodeByKey(root, "secret_key"); n != nil {
			n.SetString(string(bs))
		}

		if ref := storage.Spec.S3.HTTPConfig.TLSConfig.CA; ref != nil {
			if n := findYamlNodeByKey(root, "ca_file"); n != nil {
				n.SetString(fmt.Sprintf("%s%s/%s", constants.ConfigPath, ref.Name, ref.Key))
			}
		}

		if ref := storage.Spec.S3.HTTPConfig.TLSConfig.Cert; ref != nil {
			if n := findYamlNodeByKey(root, "cert_file"); n != nil {
				n.SetString(fmt.Sprintf("%s%s/%s", constants.ConfigPath, ref.Name, ref.Key))
			}
		}

		if ref := storage.Spec.S3.HTTPConfig.TLSConfig.Key; ref != nil {
			if n := findYamlNodeByKey(root, "key_file"); n != nil {
				n.SetString(fmt.Sprintf("%s%s/%s", constants.ConfigPath, ref.Name, ref.Key))
			}
		}

		return yaml.Marshal(root)
	}

	return nil, nil
}

func (r *BaseReconciler) getValueFromSecret(ref *corev1.SecretKeySelector, namespace string) ([]byte, error) {

	if ref == nil {
		return nil, nil
	}

	secret := &corev1.Secret{}
	if err := r.Client.Get(r.Context, client.ObjectKey{Name: ref.Name, Namespace: namespace}, secret); err != nil {
		return nil, err
	}

	return secret.Data[ref.Key], nil
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

func (r *BaseReconciler) VolumesAndVolumeMountsForStorage(namespaceName string) ([]corev1.Volume, []corev1.VolumeMount, error) {
	storage := &v1alpha1.Storage{}
	array := strings.Split(namespaceName, ".")
	if err := r.Client.Get(r.Context, client.ObjectKey{Name: array[1], Namespace: array[0]}, storage); err != nil {
		return nil, nil, err
	}

	if storage.Spec.S3 == nil {
		return nil, nil, nil
	}

	tls := storage.Spec.S3.HTTPConfig.TLSConfig
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	volumes, volumeMounts = appendVolumesAndVolumeMounts(volumes, volumeMounts, tls.CA)
	volumes, volumeMounts = appendVolumesAndVolumeMounts(volumes, volumeMounts, tls.Cert)
	volumes, volumeMounts = appendVolumesAndVolumeMounts(volumes, volumeMounts, tls.Key)

	return volumes, volumeMounts, nil
}

func appendVolumesAndVolumeMounts(volumes []corev1.Volume, volumeMounts []corev1.VolumeMount, ref *corev1.SecretKeySelector) ([]corev1.Volume, []corev1.VolumeMount) {
	if ref == nil {
		return volumes, volumeMounts
	}

	mode := corev1.SecretVolumeSourceDefaultMode
	volume := corev1.Volume{
		Name: ref.Name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  ref.Name,
				DefaultMode: &mode,
			},
		},
	}

	replaced := util.ReplaceInSlice(volumes, func(v interface{}) bool {
		return v.(corev1.Volume).Name == volume.Name
	}, volume)

	if !replaced {
		volumes = append(volumes, volume)
	}

	volumeMount := corev1.VolumeMount{
		Name:      ref.Name,
		ReadOnly:  true,
		MountPath: fmt.Sprintf("%s%s/", constants.ConfigPath, ref.Name),
	}

	replaced = util.ReplaceInSlice(volumeMounts, func(v interface{}) bool {
		return v.(corev1.VolumeMount).Name == volumeMount.Name
	}, volumeMount)

	if !replaced {
		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumes, volumeMounts
}

func BuildEnvoySidecarContainer(spec v1alpha1.SidecarSpec, volumeMounts []corev1.VolumeMount) corev1.Container {
	envoyContainer := corev1.Container{
		Name:         "envoy-sidecar",
		Args:         []string{"-c", constants.EnvoyConfigMountPath + "envoy.yaml"},
		Image:        spec.Image,
		Resources:    spec.Resources,
		VolumeMounts: volumeMounts,
	}

	return envoyContainer
}

// BuildCommonVolumes returns a set of volumes to be mounted on statefulset spec that are common components
func BuildCommonVolumes(tlsAssetSecrets []string, config string, configmaps []string, secrets []string) ([]corev1.Volume, []corev1.VolumeMount, error) {

	assetsVolume := corev1.Volume{
		Name: "tls-assets",
		VolumeSource: v1.VolumeSource{
			Projected: &corev1.ProjectedVolumeSource{
				Sources: []corev1.VolumeProjection{},
			},
		},
	}
	for _, assetShard := range tlsAssetSecrets {
		assetsVolume.Projected.Sources = append(assetsVolume.Projected.Sources,
			v1.VolumeProjection{
				Secret: &v1.SecretProjection{
					LocalObjectReference: v1.LocalObjectReference{Name: assetShard},
				},
			})
	}

	volumes := []v1.Volume{
		{
			Name: "envoy-config",
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: config,
					},
				},
			},
		},
	}

	volumeMounts := []v1.VolumeMount{
		{
			Name:      "envoy-config",
			ReadOnly:  true,
			MountPath: constants.EnvoyConfigMountPath,
		},
	}

	if len(tlsAssetSecrets) > 0 {
		volumes = append(volumes, assetsVolume)
		volumeMounts = append(volumeMounts, v1.VolumeMount{

			Name:      "tls-assets",
			ReadOnly:  true,
			MountPath: constants.EnvoyCertsMountPath,
		})
	}

	// Mount related secrets
	rn := k8sutil.NewResourceNamerWithPrefix("secret")
	for _, s := range secrets {
		name, err := rn.DNS1123Label(s)
		if err != nil {
			return nil, nil, err
		}

		volumes = append(volumes, v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: s,
				},
			},
		})
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      name,
			ReadOnly:  true,
			MountPath: constants.EnvoySecretMountPath + s,
		})
	}

	rn = k8sutil.NewResourceNamerWithPrefix("configmap")
	for _, c := range configmaps {
		name, err := rn.DNS1123Label(c)
		if err != nil {
			return nil, nil, err
		}

		volumes = append(volumes, v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: c,
					},
				},
			},
		})
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      name,
			ReadOnly:  true,
			MountPath: constants.EnvoyConfigMapMountPath + c,
		})
	}

	return volumes, volumeMounts, nil
}

func (r *BaseReconciler) GetValueFromSecret(ref *corev1.SecretKeySelector, namespace string) ([]byte, error) {

	if ref == nil {
		return nil, nil
	}

	secret := &corev1.Secret{}
	if err := r.Client.Get(r.Context, client.ObjectKey{Name: ref.Name, Namespace: namespace}, secret); err != nil {
		return nil, err
	}

	return secret.Data[ref.Key], nil
}
