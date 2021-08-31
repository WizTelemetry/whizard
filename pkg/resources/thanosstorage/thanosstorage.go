package thanosstorage

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/kubesphere/paodin-monitoring/api/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/config"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	componentName       = "thanosstorage"
	grpcPort      int32 = 10901
	httpPort      int32 = 10902

	thanosConfigDir = "/etc/thanos"

	mountDirSecrets    string = "/etc/thanos/secrets/"
	mountDirConfigMaps string = "/etc/thanos/configmaps/"

	volumeNamePrefixSecret    string = "secret-"
	volumeNamePrefixConfigMap string = "configmap-"

	defaultReplicas int32 = 1
	storageDir            = "/thanos"
)

type ThanosStorage struct {
	Cfg      config.Config
	Context  context.Context
	Client   client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Instance v1alpha1.ThanosStorage
}

func (s *ThanosStorage) getGatewayStatefulSetName() string {
	return fmt.Sprintf("%s-%s-gateway", componentName, s.Instance.Name)
}

func (s *ThanosStorage) getCompactStatefulSetName() string {
	return fmt.Sprintf("%s-%s-compact", componentName, s.Instance.Name)
}

func (s *ThanosStorage) getGatewayTSDBVolumeName() string {
	if gateway := s.Instance.Spec.Gateway; gateway != nil && gateway.DataVolume != nil {
		if pvc := gateway.DataVolume.PersistentVolumeClaim; pvc != nil && pvc.Name != "" {
			return pvc.Name
		}
	}
	return fmt.Sprintf("%s-%s-gateway-tsdb", componentName, s.Instance.Name)
}

func (s *ThanosStorage) getCompactTSDBVolumeName() string {
	if gateway := s.Instance.Spec.Gateway; gateway != nil && gateway.DataVolume != nil {
		if pvc := gateway.DataVolume.PersistentVolumeClaim; pvc != nil && pvc.Name != "" {
			return pvc.Name
		}
	}
	return fmt.Sprintf("%s-%s-compact-tsdb", componentName, s.Instance.Name)
}

func (s *ThanosStorage) getGatewayOperatedServiceName() string {
	return fmt.Sprintf("%s-%s-gateway-operated", componentName, s.Instance.Name)
}

func (s *ThanosStorage) getCompactOperatedServiceName() string {
	return fmt.Sprintf("%s-%s-compact-operated", componentName, s.Instance.Name)
}
