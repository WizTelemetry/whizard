package monitoringgateway

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
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

		// remove the prefix /:tenant_id from path
		if index := indexByteNth(req.URL.Path, '/', 2); index > 0 {
			req.URL.Path = req.URL.Path[index:]
		}
		ctx := req.Context()

		req = req.WithContext(context.WithValue(ctx, requestInfoKey, &RequestInfo{
			TenantId: mux.Vars(req)["tenant_id"],
		}))

		f.ServeHTTP(w, req)
	})
}
