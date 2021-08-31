package thanosquery

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/kubesphere/paodin-monitoring/api/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/config"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	componentName       = "thanosquery"
	grpcPort      int32 = 10901
	httpPort      int32 = 10902

	thanosConfigDir = "/etc/thanos"
	storeSDFileName = "store-sd.yaml"

	envoyConfigDir         = "/etc/envoy"
	envoyConfigFileName    = "envoy.yaml"
	envoyLDSFileName       = "lds.yaml"
	envoyCDSFileName       = "cds.yaml"
	envoySecretsDir        = "/etc/envoy/secrets"
	envoyListenerAddress   = "127.0.0.1"
	envoyListenerStartPort = 11000

	volumeNamePrefixSecret    string = "secret-"
	volumeNamePrefixConfigMap string = "configmap-"

	defaultReplicas int32 = 1
)

type ThanosQuery struct {
	Cfg      config.Config
	Context  context.Context
	Client   client.Client
	Instance v1alpha1.ThanosQuery
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (q *ThanosQuery) getDeploymentName() string {
	return componentName + "-" + q.Instance.Name
}

func (q *ThanosQuery) getServiceName() string {
	return componentName + "-" + q.Instance.Name + "-operated"
}

func (q *ThanosQuery) getStoreSDConfigMapName() string {
	return componentName + "-" + q.Instance.Name + "-store-sd"
}

func (q *ThanosQuery) getEnvoyConfigMapName() string {
	return componentName + "-" + q.Instance.Name + "-envoy"
}

func (q *ThanosQuery) getHttpIngressName() string {
	return componentName + "-" + q.Instance.Name + "-http"
}

func (q *ThanosQuery) getGrpcIngressName() string {
	return componentName + "-" + q.Instance.Name + "-grpc"
}

func storeRequireProxy(store v1alpha1.QueryStore) bool {
	return store.Address != "" && store.SecretName != ""
}
