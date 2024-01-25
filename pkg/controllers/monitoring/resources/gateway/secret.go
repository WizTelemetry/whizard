package gateway

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

const (
	TLSVersionTLS12 = "TLS12"
	TLSVersionTLS13 = "TLS13"
)

func (g *Gateway) webConfigSecret() (runtime.Object, resources.Operation, error) {
	var secret = &corev1.Secret{ObjectMeta: g.meta(g.name("web-config"))}

	if g.gateway == nil {
		return secret, resources.OperationDelete, nil
	}

	if g.gateway.Spec.WebConfig == nil {
		return secret, resources.OperationDelete, nil
	}

	if len(g.gateway.Spec.WebConfig.BasicAuthUsers) > 0 {

		var secret = &corev1.Secret{}
		err := g.Client.Get(g.Context, types.NamespacedName{Namespace: g.gateway.Namespace, Name: g.name("built-in-user")}, secret)
		if errors.IsNotFound(err) {
			if err := g.generateBuiltInBasicAuthUserSecret(); err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
		} else if err != nil {
			return nil, resources.OperationCreateOrUpdate, err
		}

		g.gateway.Spec.WebConfig.BasicAuthUsers = append(g.gateway.Spec.WebConfig.BasicAuthUsers, v1alpha1.BasicAuth{
			Username: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: g.name("built-in-user"),
				},
				Key: "username",
			},
			Password: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: g.name("built-in-user"),
				},
				Key: "password",
			},
		})
	}

	body, err := g.BaseReconciler.CreateWebConfig(g.gateway.Namespace, g.gateway.Spec.WebConfig)
	if err != nil {
		return nil, resources.OperationDelete, err
	}

	secret.Data = map[string][]byte{
		constants.WhizardWebConfigFile: body,
	}

	return secret, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(g.gateway, secret, g.Scheme)
}

func (g *Gateway) generateBuiltInBasicAuthUserSecret() error {

	user := randomString(16)
	clearTextPassword := randomString(32)
	password, err := bcrypt.GenerateFromPassword([]byte(clearTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	var secret = &corev1.Secret{ObjectMeta: g.meta(g.name("built-in-user"))}
	secret.StringData = map[string]string{
		"username":  user,
		"password":  string(password),
		"cpassword": clearTextPassword,
	}
	if err := ctrl.SetControllerReference(g.gateway, secret, g.Scheme); err != nil {
		return err
	}

	return g.Client.Create(g.Context, secret)
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// randomString returns a random string with a fixed length
func randomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
