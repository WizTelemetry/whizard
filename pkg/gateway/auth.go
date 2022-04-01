package gateway

import (
	"errors"
	"net/http"
)

var errInvalidCert = errors.New("invalid cert")

type CertAuthenticator struct {
	serviceGetter struct{}
	agentGetter   struct{}
}

func NewCertAuthenticator() *CertAuthenticator {
	return &CertAuthenticator{}
}

func (cauth *CertAuthenticator) AuthenticateRequest(req *http.Request) (agent string, ok bool) {

	if len(req.TLS.PeerCertificates) == 0 {
		return "", false
	}

	cert := req.TLS.PeerCertificates[0]

	if len(cert.Subject.Organization) == 0 {
		return "", false
	}

	return cert.Subject.Organization[0] + "/" + cert.Subject.CommonName, true
}

func withAuthorization(f http.HandlerFunc, certAuthenticator *CertAuthenticator) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		agent, ok := certAuthenticator.AuthenticateRequest(req)

		if !ok {
			http.Error(w, errInvalidCert.Error(), http.StatusUnauthorized)
		}

		requestInfo, found := requestInfoFrom(ctx)
		if !found || requestInfo.Agent == nil {
			http.NotFound(w, req)
		}

		if agent != requestInfo.Agent.Namespace+"/"+requestInfo.Agent.Name {
			http.Error(w, errInvalidCert.Error(), http.StatusUnauthorized)
		}

		f.ServeHTTP(w, req)
	})
}
