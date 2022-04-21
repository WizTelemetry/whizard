package monitoringgateway

import (
	"errors"
	"net/http"
)

var errInvalidCert = errors.New("invalid cert")

type CertAuthenticator struct {
}

func NewCertAuthenticator() *CertAuthenticator {
	return &CertAuthenticator{}
}

func (cauth *CertAuthenticator) AuthenticateRequest(req *http.Request) (tenantId string, ok bool) {

	if len(req.TLS.PeerCertificates) == 0 {
		return "", false
	}

	return req.TLS.PeerCertificates[0].Subject.CommonName, true
}

func withAuthorization(f http.HandlerFunc, certAuthenticator *CertAuthenticator) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		tenantId, ok := certAuthenticator.AuthenticateRequest(req)
		if !ok {
			http.Error(w, errInvalidCert.Error(), http.StatusUnauthorized)
		}

		requestInfo, found := requestInfoFrom(ctx)
		if !found {
			http.Error(w, errInvalidCert.Error(), http.StatusUnauthorized)
		}

		if tenantId != requestInfo.TenantId {
			http.Error(w, errInvalidCert.Error(), http.StatusUnauthorized)
		}

		f.ServeHTTP(w, req)
	})
}
