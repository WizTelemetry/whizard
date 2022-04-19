package receive

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/resources"
	"github.com/kubesphere/paodin-monitoring/pkg/util"
)

var (
	configDir     = "/etc/thanos"
	hashringsFile = "hashrings.json"
	secretsDir    = "/etc/thanos/secrets"
	storageDir    = "/thanos"

	LabelReceiveTypeKey           = "receive.type"
	LabelReceiveTypeValueRouter   = "router"
	LabelReceiveTypeValueIngestor = "ingestor"
	LabelReceiveIngestorNameKey   = "ingestor.name"
)

type Receive struct {
	resources.ServiceBaseReconciler
	receive *v1alpha1.Receive
}

func New(reconciler resources.ServiceBaseReconciler) *Receive {
	return &Receive{
		ServiceBaseReconciler: reconciler,
		receive:               reconciler.Service.Spec.Thanos.Receive,
	}
}

func (r *Receive) labels() map[string]string {
	labels := r.BaseLabels()
	labels["app.kubernetes.io/name"] = "receive"
	return labels
}

func (r *Receive) name(nameSuffix ...string) string {
	name := r.Service.Name + "-receive"
	if len(nameSuffix) > 0 {
		name += "-" + strings.Join(nameSuffix, "-")
	}
	return name
}

func (r *Receive) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.Service.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *Receive) HttpAddr() string {
	router := receiveRouter{Receive: *r, Router: r.receive.Router}
	return fmt.Sprintf("http://%s:10902", router.name("operated"))
}

func (r *Receive) RemoteWriteAddr() string {
	router := receiveRouter{Receive: *r, Router: r.receive.Router}
	return fmt.Sprintf("http://%s:19291", router.name("operated"))
}

func (r *Receive) GrpcAddrs() []string {
	var svcs []string
	if r.receive != nil {
		for _, i := range r.receive.Ingestors {
			in := i
			ingestor := receiveIngestor{Receive: *r, Ingestor: in}
			svcs = append(svcs, ingestor.grpcAddrs()...)
		}
	}
	return svcs
}

func (r *Receive) Reconcile() error {

	var ress []resources.Resource

	var stss = &appsv1.StatefulSetList{}
	var ingestorCommonLabels = r.labels()
	ingestorCommonLabels[LabelReceiveTypeKey] = LabelReceiveTypeValueIngestor
	if err := r.Client.List(r.Context, stss, &client.ListOptions{
		Namespace:     r.Service.Namespace,
		LabelSelector: labels.SelectorFromSet(ingestorCommonLabels),
	}); err != nil {
		return err
	}
	var oldIngestors []string
	for _, item := range stss.Items {
		sts := item
		for _, ref := range r.OwnerReferences() {
			if util.IndexOwnerRef(sts.ObjectMeta.OwnerReferences, ref) > 0 {
				if name, ok := sts.Labels[LabelReceiveIngestorNameKey]; ok {
					oldIngestors = append(oldIngestors, name)
				}
			}
		}
	}

	if r.receive == nil {
		router := receiveRouter{
			Receive: *r,
			Router:  v1alpha1.ReceiveRouter{},
			del:     true,
		}
		ress = append(ress, router.resources()...)
		for _, name := range oldIngestors {
			ingestor := receiveIngestor{
				Receive:  *r,
				Ingestor: v1alpha1.ReceiveIngestor{Name: name},
				del:      true,
			}
			ress = append(ress, ingestor.resources()...)
		}
		return r.ReconcileResources(ress)
	}

	var newIngestors = map[string]struct{}{}
	for _, i := range r.receive.Ingestors {
		in := i
		ingestor := receiveIngestor{Receive: *r, Ingestor: in}
		ress = append(ress, ingestor.resources()...)
		newIngestors[r.name("ingestor", i.Name)] = struct{}{}
	}
	router := receiveRouter{Receive: *r, Router: r.receive.Router}
	ress = append(ress, router.resources()...)

	for _, name := range oldIngestors {
		if _, ok := newIngestors[name]; ok {
			continue
		}
		ingestor := receiveIngestor{
			Receive:  *r,
			Ingestor: v1alpha1.ReceiveIngestor{Name: name},
			del:      true,
		}
		ress = append(ress, ingestor.resources()...)
	}

	return r.ReconcileResources(ress)
}

type receiveRouter struct {
	Receive
	del    bool
	Router v1alpha1.ReceiveRouter
}

func (r *receiveRouter) labels() map[string]string {
	labels := r.Receive.labels()
	labels[LabelReceiveTypeKey] = LabelReceiveTypeValueRouter
	return labels
}

func (r *receiveRouter) name(nameSuffix ...string) string {
	return r.Receive.name(append([]string{"router"}, nameSuffix...)...)
}

func (r *receiveRouter) meta(name string) metav1.ObjectMeta {
	meta := r.Receive.meta(name)
	meta.Labels = r.labels()
	return meta
}

func (r *receiveRouter) resources() []resources.Resource {
	var ress []resources.Resource
	ress = append(ress, r.hashringsConfigMap)
	ress = append(ress, r.deployment)
	ress = append(ress, r.service)
	return ress
}

type receiveIngestor struct {
	Receive
	del      bool
	Ingestor v1alpha1.ReceiveIngestor
}

func (r *receiveIngestor) labels() map[string]string {
	labels := r.Receive.labels()
	labels[LabelReceiveTypeKey] = LabelReceiveTypeValueIngestor
	labels[LabelReceiveIngestorNameKey] = r.Ingestor.Name

	return labels
}

func (r *receiveIngestor) name(nameSuffix ...string) string {
	return r.Receive.name(append([]string{"ingestor", r.Ingestor.Name}, nameSuffix...)...)
}

func (r *receiveIngestor) meta(name string) metav1.ObjectMeta {
	meta := r.Receive.meta(name)
	meta.Labels = r.labels()
	return meta
}

func (r *receiveIngestor) grpcAddrs() []string {
	addrs := make([]string, *r.Ingestor.Replicas)
	for i := range addrs {
		addrs[i] = fmt.Sprintf("%s-%d.%s:10901", r.name(), i, r.name("operated"))
	}
	return addrs
}

func (r *receiveIngestor) resources() []resources.Resource {
	var ress []resources.Resource
	ress = append(ress, r.statefulSet)
	ress = append(ress, r.service)
	return ress
}
