package thanosreceive

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/paodin-monitoring/api/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/config"
)

const (
	componentName         = "thanosreceive"
	grpcPort        int32 = 10901
	httpPort        int32 = 10902
	remoteWritePort int32 = 19291

	thanosConfigDir   = "/etc/thanos"
	hashringsFileName = "hashrings.json"

	mountDirSecrets    string = "/etc/thanos/secrets/"
	mountDirConfigMaps string = "/etc/thanos/configmaps/"

	volumeNamePrefixSecret    string = "secret-"
	volumeNamePrefixConfigMap string = "configmap-"

	defaultReplicas int32 = 1
	storageDir            = "/thanos"

	RouterOnly     ReceiveMode = "RouterOnly"
	IngestorOnly   ReceiveMode = "IngestorOnly"
	RouterIngestor ReceiveMode = "RouterIngestor"
)

// ReceiveMode represents how the Receive should be deployed
type ReceiveMode string

type ThanosReceive struct {
	Cfg      config.Config
	Context  context.Context
	Client   client.Client
	Instance v1alpha1.ThanosReceive
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *ThanosReceive) GetMode() ReceiveMode {
	hasRouting := r.Instance.Spec.Router != nil
	hasIngesting := r.Instance.Spec.Ingestor != nil
	switch {
	case hasRouting && hasIngesting:
		return RouterIngestor
	case hasRouting && !hasIngesting:
		return RouterOnly
	default:
		return IngestorOnly
	}
}

func (r *ThanosReceive) getReceiveDeploymentName() string {
	return fmt.Sprintf("%s-%s", componentName, r.Instance.Name)
}

func (r *ThanosReceive) getReceiveStatefulSetName() string {
	return fmt.Sprintf("%s-%s", componentName, r.Instance.Name)
}

func (r *ThanosReceive) getReceiveOperatedServiceName() string {
	return fmt.Sprintf("%s-%s-operated", componentName, r.Instance.Name)
}

func (r *ThanosReceive) getRemoteWriteIngressName() string {
	return fmt.Sprintf("%s-%s-remote-write", componentName, r.Instance.Name)
}

func (r *ThanosReceive) getTSDBVolumeName() string {
	if ingesting := r.Instance.Spec.Ingestor; ingesting != nil && ingesting.DataVolume != nil {
		if pvc := ingesting.DataVolume.PersistentVolumeClaim; pvc != nil && pvc.Name != "" {
			return pvc.Name
		}
	}
	return fmt.Sprintf("%s-%s-tsdb", componentName, r.Instance.Name)
}

func (r *ThanosReceive) getHashringsConfigMapName() string {
	return fmt.Sprintf("%s-%s-hashrings", componentName, r.Instance.Name)
}

// HashringConfig represents the configuration for a hashring
// a receive node knows about.
type HashringConfig struct {
	Hashring  string   `json:"hashring,omitempty"`
	Tenants   []string `json:"tenants,omitempty"`
	Endpoints []string `json:"endpoints"`
}
