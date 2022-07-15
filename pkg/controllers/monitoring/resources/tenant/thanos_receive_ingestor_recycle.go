package tenant

import (
	"context"
	"errors"
	"strconv"
	"time"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IngestorRecycleController struct {
	DeleteIngestorEventChan        chan DeleteIngestorEvent
	deleteChan                     chan types.NamespacedName
	DefaultIngestorRetentionPeriod time.Duration

	Client  client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

type DeleteIngestorEvent struct {
	NamespacedName types.NamespacedName
	DeleteDuration time.Duration
}

const (
	retryDuration = time.Minute * 10
)

func NewIngestorRecycleController(client client.Client, scheme *runtime.Scheme, context context.Context, deleteIngestorEventChan chan DeleteIngestorEvent) *IngestorRecycleController {
	r := &IngestorRecycleController{
		DeleteIngestorEventChan: deleteIngestorEventChan,
		deleteChan:              make(chan types.NamespacedName, 100),
		Client:                  client,
		Scheme:                  scheme,
		Context:                 context,
	}
	return r
}

func (r *IngestorRecycleController) Recycle() error {

	go func() {
		if err := r.traverse(); err != nil {
			klog.Errorf("traverse ingestor error: %v", err)
		}
		for {
			select {
			case e := <-r.DeleteIngestorEventChan:
				klog.V(3).Infof("ThanosReceiveIngestor %s/%s will be deleted after %s.", e.NamespacedName.Namespace, e.NamespacedName.Name, e.DeleteDuration.String())
				time.AfterFunc(e.DeleteDuration, func() {
					r.deleteChan <- e.NamespacedName
				})

			case namespacedName := <-r.deleteChan:
				err := r.deleteIngestorInstance(namespacedName.Namespace, namespacedName.Name)
				if err != nil {
					klog.Errorf("%v", err)
					event := DeleteIngestorEvent{
						NamespacedName: namespacedName,
						DeleteDuration: retryDuration,
					}
					r.DeleteIngestorEventChan <- event
				}

			}
		}
	}()

	return nil
}

// traverse ingestors looking for ingestors to be deleted
func (r *IngestorRecycleController) traverse() error {

	ingestorList := &monitoringv1alpha1.ThanosReceiveIngestorList{}

	backoff := wait.Backoff{
		Steps:    5,
		Duration: 10 * time.Second,
		Cap:      10 * time.Minute,
	}
	isErrCacheNotStartedFunc := func(err error) bool {
		return errors.Is(err, &cache.ErrCacheNotStarted{})
	}

	// retry when InformerCache not ready
	err := retry.OnError(backoff, isErrCacheNotStartedFunc, func() error {
		ingestorList = &monitoringv1alpha1.ThanosReceiveIngestorList{}
		return r.Client.List(context.Background(), ingestorList)
	})
	if err != nil {
		return err
	}

	for _, ingestor := range ingestorList.Items {
		annotations := ingestor.GetAnnotations()
		if annotations != nil {
			if v, ok := annotations[resources.LabelNameReceiveIngestorState]; ok && v == "deleting" {
				if len(ingestor.Spec.Tenants) == 0 {
					event := DeleteIngestorEvent{
						NamespacedName: types.NamespacedName{
							Namespace: ingestor.Namespace,
							Name:      ingestor.Name,
						},
					}
					if deletingTime, ok := annotations[resources.LabelNameReceiveIngestorDeletingTime]; ok {
						i, err := strconv.ParseInt(deletingTime, 10, 64)
						if err != nil {
							event.DeleteDuration = retryDuration
						}
						d := time.Since(time.Unix(i, 0))
						// now - DeletingTime >  DefaultIngestorRetentionPeriod-retryretryDuration
						if d > r.DefaultIngestorRetentionPeriod-retryDuration {
							event.DeleteDuration = retryDuration
						} else {
							event.DeleteDuration = r.DefaultIngestorRetentionPeriod - d
						}
					} else {
						event.DeleteDuration = retryDuration
					}
					r.DeleteIngestorEventChan <- event
				}
			}
		}
	}
	return nil
}

func (r *IngestorRecycleController) deleteIngestorInstance(namespace, name string) error {
	ingestor := &monitoringv1alpha1.ThanosReceiveIngestor{}
	err := r.Client.Get(context.Background(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, ingestor)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	annotations := ingestor.GetAnnotations()
	if annotations != nil {
		if v, ok := annotations[resources.LabelNameReceiveIngestorState]; ok && v == "deleting" {
			if len(ingestor.Spec.Tenants) == 0 {
				klog.V(3).Infof("ThanosReceiveIngestor %s/%s will be deleted.", namespace, name)
				return r.Client.Delete(context.Background(), ingestor)
			}
		}
	}
	return nil
}
