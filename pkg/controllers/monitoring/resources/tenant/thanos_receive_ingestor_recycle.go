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

type IngesterRecycleController struct {
	DeleteIngesterEventChan        chan DeleteIngesterEvent
	deleteChan                     chan types.NamespacedName
	DefaultIngesterRetentionPeriod time.Duration

	Client  client.Client
	Scheme  *runtime.Scheme
	Context context.Context
}

type DeleteIngesterEvent struct {
	NamespacedName types.NamespacedName
	DeleteDuration time.Duration
}

const (
	retryDuration = time.Minute * 10
)

func NewIngesterRecycleController(client client.Client, scheme *runtime.Scheme, context context.Context, deleteIngesterEventChan chan DeleteIngesterEvent) *IngesterRecycleController {
	r := &IngesterRecycleController{
		DeleteIngesterEventChan: deleteIngesterEventChan,
		deleteChan:              make(chan types.NamespacedName, 100),
		Client:                  client,
		Scheme:                  scheme,
		Context:                 context,
	}
	return r
}

func (r *IngesterRecycleController) Recycle() error {

	go func() {
		if err := r.traverse(); err != nil {
			klog.Errorf("traverse ingester error: %v", err)
		}
		for {
			select {
			case e := <-r.DeleteIngesterEventChan:
				klog.V(3).Infof("ThanosReceiveIngester %s/%s will be deleted after %s.", e.NamespacedName.Namespace, e.NamespacedName.Name, e.DeleteDuration.String())
				time.AfterFunc(e.DeleteDuration, func() {
					r.deleteChan <- e.NamespacedName
				})

			case namespacedName := <-r.deleteChan:
				err := r.deleteIngesterInstance(namespacedName.Namespace, namespacedName.Name)
				if err != nil {
					klog.Errorf("%v", err)
					event := DeleteIngesterEvent{
						NamespacedName: namespacedName,
						DeleteDuration: retryDuration,
					}
					r.DeleteIngesterEventChan <- event
				}

			}
		}
	}()

	return nil
}

// traverse ingesters looking for ingesters to be deleted
func (r *IngesterRecycleController) traverse() error {

	ingesterList := &monitoringv1alpha1.IngesterList{}

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
		ingesterList = &monitoringv1alpha1.IngesterList{}
		return r.Client.List(context.Background(), ingesterList)
	})
	if err != nil {
		return err
	}

	for _, ingester := range ingesterList.Items {
		annotations := ingester.GetAnnotations()
		if annotations != nil {
			if v, ok := annotations[resources.LabelNameReceiveIngesterState]; ok && v == "deleting" {
				if len(ingester.Spec.Tenants) == 0 {
					event := DeleteIngesterEvent{
						NamespacedName: types.NamespacedName{
							Namespace: ingester.Namespace,
							Name:      ingester.Name,
						},
					}
					if deletingTime, ok := annotations[resources.LabelNameReceiveIngesterDeletingTime]; ok {
						i, err := strconv.ParseInt(deletingTime, 10, 64)
						if err != nil {
							event.DeleteDuration = retryDuration
						}
						d := time.Since(time.Unix(i, 0))
						// now - DeletingTime >  DefaultIngesterRetentionPeriod-retryretryDuration
						if d > r.DefaultIngesterRetentionPeriod-retryDuration {
							event.DeleteDuration = retryDuration
						} else {
							event.DeleteDuration = r.DefaultIngesterRetentionPeriod - d
						}
					} else {
						event.DeleteDuration = retryDuration
					}
					r.DeleteIngesterEventChan <- event
				}
			}
		}
	}
	return nil
}

func (r *IngesterRecycleController) deleteIngesterInstance(namespace, name string) error {
	ingester := &monitoringv1alpha1.Ingester{}
	err := r.Client.Get(context.Background(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, ingester)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	annotations := ingester.GetAnnotations()
	if annotations != nil {
		if v, ok := annotations[resources.LabelNameReceiveIngesterState]; ok && v == "deleting" {
			if len(ingester.Spec.Tenants) == 0 {
				klog.V(3).Infof("ThanosReceiveIngester %s/%s will be deleted.", namespace, name)
				return r.Client.Delete(context.Background(), ingester)
			}
		}
	}
	return nil
}
