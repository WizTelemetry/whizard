package resources

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/util"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func QualifiedName(appName, instanceName string, suffix ...string) string {
	name := appName + "-" + instanceName
	if len(suffix) > 0 {
		name += "-" + strings.Join(suffix, "-")
	}
	return name
}

func DefaultLivenessProbe() *corev1.Probe {
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

func DefaultReadinessProbe() *corev1.Probe {
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

func AddTSDBVolume(sts *appsv1.StatefulSet, container *corev1.Container, dataVolume *v1alpha1.KubernetesVolume) {
	if dataVolume == nil ||
		(dataVolume.PersistentVolumeClaim == nil && dataVolume.EmptyDir == nil) {
		return
	}

	if dataVolume.PersistentVolumeClaim != nil {
		pvc := *dataVolume.PersistentVolumeClaim
		if pvc.Name == "" {
			pvc.Name = constants.TSDBVolumeName
		}
		if pvc.Spec.AccessModes == nil {
			pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		}

		sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{pvc}
	} else if dataVolume.EmptyDir != nil {
		sts.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: constants.TSDBVolumeName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: dataVolume.EmptyDir,
				},
			},
		}
	}

	container.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      constants.TSDBVolumeName,
			MountPath: constants.StorageDir,
		},
	}
}

func GetTenantHash(ctx context.Context, c client.Client, selector map[string]string) (string, error) {
	tenantList := &v1alpha1.TenantList{}
	err := c.List(ctx, tenantList, client.MatchingLabels(selector))
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

func GetStorageHash(ctx context.Context, c client.Client, namespaceName string) (string, error) {
	storage := &v1alpha1.Storage{}
	array := strings.Split(namespaceName, ".")
	if err := c.Get(ctx, client.ObjectKey{Name: array[1], Namespace: array[0]}, storage); err != nil {
		return "", err
	}

	storageConfig, err := GetStorageConfig(ctx, c, namespaceName)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	hash.Write(storageConfig)

	if storage.Spec.S3 != nil {
		tls := storage.Spec.S3.HTTPConfig.TLSConfig
		if bs, err := getValueFromSecret(ctx, c, tls.CA, storage.Namespace); err != nil {
			return "", err
		} else {
			hash.Write(bs)
		}

		if bs, err := getValueFromSecret(ctx, c, tls.Key, storage.Namespace); err != nil {
			return "", err
		} else {
			hash.Write(bs)
		}

		if bs, err := getValueFromSecret(ctx, c, tls.Cert, storage.Namespace); err != nil {
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

func GetStorageConfig(ctx context.Context, c client.Client, namespaceName string) ([]byte, error) {
	storage := &v1alpha1.Storage{}
	array := strings.Split(namespaceName, ".")
	if err := c.Get(ctx, client.ObjectKey{Name: array[1], Namespace: array[0]}, storage); err != nil {
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

		bs, err = getValueFromSecret(ctx, c, storage.Spec.S3.AccessKey, storage.Namespace)
		if err != nil {
			return nil, err
		}
		if n := findYamlNodeByKey(root, "access_key"); n != nil {
			n.SetString(string(bs))
		}

		bs, err = getValueFromSecret(ctx, c, storage.Spec.S3.SecretKey, storage.Namespace)
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

func getValueFromSecret(ctx context.Context, c client.Client, ref *corev1.SecretKeySelector, namespace string) ([]byte, error) {

	if ref == nil {
		return nil, nil
	}

	secret := &corev1.Secret{}
	if err := c.Get(ctx, client.ObjectKey{Name: ref.Name, Namespace: namespace}, secret); err != nil {
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

func VolumesAndVolumeMountsForStorage(ctx context.Context, c client.Client, namespaceName string) ([]corev1.Volume, []corev1.VolumeMount, error) {
	storage := &v1alpha1.Storage{}
	array := strings.Split(namespaceName, ".")
	if err := c.Get(ctx, client.ObjectKey{Name: array[1], Namespace: array[0]}, storage); err != nil {
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
