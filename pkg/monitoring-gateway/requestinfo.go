package monitoringgateway

import (
	"context"
	"net/http"

	"github.com/prometheus/common/route"
)

type requestInfoKeyType int

const requestInfoKey requestInfoKeyType = iota

type RequestInfo struct {
	TenantId string
}

func requestInfoFrom(ctx context.Context) (*RequestInfo, bool) {
	info, ok := ctx.Value(requestInfoKey).(*RequestInfo)
	return info, ok
}

func withRequestInfo(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		req = req.WithContext(context.WithValue(ctx, requestInfoKey, &RequestInfo{
			TenantId: route.Param(ctx, "tenant_id"),
		}))

		f.ServeHTTP(w, req)
	})
}
