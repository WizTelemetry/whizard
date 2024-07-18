package query

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
)

const (
	configDir  = "/etc/whizard"
	storesFile = "store-sd.yaml"
)

type Query struct {
	resources.BaseReconciler
	query *v1alpha1.Query
}

func New(reconciler resources.BaseReconciler, q *v1alpha1.Query) (*Query, error) {
	if err := reconciler.SetService(q); err != nil {
		return nil, err
	}
	return &Query{
		BaseReconciler: reconciler,
		query:          q,
	}, nil
}

type Stores struct {
	DirectStores []DirectStore
	ProxyStores  []ProxyStore
}

type DirectStore struct {
	Address string
}

type ProxyStore struct {
	ListenHost string
	ListenPort uint32
	TargetHost string
	TargetPort uint32
	TlsCaFile  string
}

func (q *Query) stores() (*Stores, error) {
	var stores = &Stores{}
	var listenPort uint32 = 11000

	for _, store := range q.query.Spec.Stores {
		for _, address := range store.Addresses {
			if store.CASecret == nil {
				stores.DirectStores = append(stores.DirectStores, DirectStore{Address: address})
				continue
			}
			host, portString, err := net.SplitHostPort(address)
			if err != nil {
				return nil, err
			}
			port, err := strconv.ParseUint(portString, 10, 32)
			if err != nil {
				return nil, err
			}
			stores.ProxyStores = append(stores.ProxyStores, ProxyStore{
				ListenHost: "127.0.0.1",
				ListenPort: listenPort,
				TargetHost: host,
				TargetPort: uint32(port),
				TlsCaFile:  filepath.Join(envoySecretsDir, store.CASecret.Name, store.CASecret.Key),
			})
			listenPort++
		}
	}

	return stores, nil
}

func (q *Query) labels() map[string]string {
	labels := q.BaseLabels()
	labels[constants.LabelNameAppName] = constants.AppNameQuery
	labels[constants.LabelNameAppManagedBy] = q.Service.Name
	return labels
}

func (q *Query) name(nameSuffix ...string) string {
	return q.QualifiedName(constants.AppNameQuery, q.query.Name, nameSuffix...)
}

func (q *Query) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: q.Service.Namespace,
		Labels:    q.labels(),
	}
}

func (q *Query) HttpAddr() string {
	return fmt.Sprintf("http://%s.%s.svc:%d",
		q.name(constants.ServiceNameSuffix), q.Service.Namespace, constants.HTTPPort)
}

func (q *Query) HttpsAddr() string {
	return fmt.Sprintf("https://%s.%s.svc:%d",
		q.name(constants.ServiceNameSuffix), q.Service.Namespace, constants.HTTPPort)
}

func (q *Query) Reconcile() error {
	return q.ReconcileResources([]resources.Resource{
		q.proxyConfigMap,
		q.storesConfigMap,
		q.webConfigSecret,
		q.deployment,
		q.service,
	})
}
