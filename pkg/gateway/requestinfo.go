package gateway

import (
	"context"
	"net/http"

	"github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/prometheus/common/route"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

type requestInfoKeyType int

const requestInfoKey requestInfoKeyType = iota

type RequestInfo struct {
	Agent *v1alpha1.Agent
}

func requestInfoFrom(ctx context.Context) (*RequestInfo, bool) {
	info, ok := ctx.Value(requestInfoKey).(*RequestInfo)
	return info, ok
}

func withRequestInfo(f http.HandlerFunc, getAgent GetAgentFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		agent, err := getAgent(&types.NamespacedName{
			Namespace: route.Param(ctx, "namespace"),
			Name:      route.Param(ctx, "name"),
		})

		if err != nil {
			if errors.IsNotFound(err) {
				http.NotFound(w, req)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if agent == nil {
			http.NotFound(w, req)
			return
		}

		req = req.WithContext(context.WithValue(ctx, requestInfoKey, &RequestInfo{
			Agent: agent,
		}))

		f.ServeHTTP(w, req)
	})
}

type GetAgentFunc func(name *types.NamespacedName) (*v1alpha1.Agent, error)
