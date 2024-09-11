/*
Copyright 2024 the Whizard Authors.

Licensed under Apache License, Version 2.0 with a few additional conditions.

You may obtain a copy of the License at

    https://github.com/WhizardTelemetry/whizard/blob/main/LICENSE
*/

package k8s

import (
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type nullClient struct {
}

func NewNullClient() Client {
	return &nullClient{}
}

func (n nullClient) Kubernetes() kubernetes.Interface {
	return nil
}

func (n nullClient) ApiExtensions() apiextensionsclient.Interface {
	return nil
}

func (n nullClient) Discovery() discovery.DiscoveryInterface {
	return nil
}

func (n nullClient) Master() string {
	return ""
}

func (n nullClient) Config() *rest.Config {
	return nil
}
