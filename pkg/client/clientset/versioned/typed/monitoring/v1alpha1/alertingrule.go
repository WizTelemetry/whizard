/*
Copyright 2021 The KubeSphere authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	scheme "github.com/kubesphere/paodin/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AlertingRulesGetter has a method to return a AlertingRuleInterface.
// A group's client should implement this interface.
type AlertingRulesGetter interface {
	AlertingRules(namespace string) AlertingRuleInterface
}

// AlertingRuleInterface has methods to work with AlertingRule resources.
type AlertingRuleInterface interface {
	Create(ctx context.Context, alertingRule *v1alpha1.AlertingRule, opts v1.CreateOptions) (*v1alpha1.AlertingRule, error)
	Update(ctx context.Context, alertingRule *v1alpha1.AlertingRule, opts v1.UpdateOptions) (*v1alpha1.AlertingRule, error)
	UpdateStatus(ctx context.Context, alertingRule *v1alpha1.AlertingRule, opts v1.UpdateOptions) (*v1alpha1.AlertingRule, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.AlertingRule, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.AlertingRuleList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AlertingRule, err error)
	AlertingRuleExpansion
}

// alertingRules implements AlertingRuleInterface
type alertingRules struct {
	client rest.Interface
	ns     string
}

// newAlertingRules returns a AlertingRules
func newAlertingRules(c *MonitoringV1alpha1Client, namespace string) *alertingRules {
	return &alertingRules{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the alertingRule, and returns the corresponding alertingRule object, and an error if there is any.
func (c *alertingRules) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AlertingRule, err error) {
	result = &v1alpha1.AlertingRule{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("alertingrules").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AlertingRules that match those selectors.
func (c *alertingRules) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AlertingRuleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.AlertingRuleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("alertingrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested alertingRules.
func (c *alertingRules) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("alertingrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a alertingRule and creates it.  Returns the server's representation of the alertingRule, and an error, if there is any.
func (c *alertingRules) Create(ctx context.Context, alertingRule *v1alpha1.AlertingRule, opts v1.CreateOptions) (result *v1alpha1.AlertingRule, err error) {
	result = &v1alpha1.AlertingRule{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("alertingrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(alertingRule).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a alertingRule and updates it. Returns the server's representation of the alertingRule, and an error, if there is any.
func (c *alertingRules) Update(ctx context.Context, alertingRule *v1alpha1.AlertingRule, opts v1.UpdateOptions) (result *v1alpha1.AlertingRule, err error) {
	result = &v1alpha1.AlertingRule{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("alertingrules").
		Name(alertingRule.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(alertingRule).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *alertingRules) UpdateStatus(ctx context.Context, alertingRule *v1alpha1.AlertingRule, opts v1.UpdateOptions) (result *v1alpha1.AlertingRule, err error) {
	result = &v1alpha1.AlertingRule{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("alertingrules").
		Name(alertingRule.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(alertingRule).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the alertingRule and deletes it. Returns an error if one occurs.
func (c *alertingRules) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("alertingrules").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *alertingRules) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("alertingrules").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched alertingRule.
func (c *alertingRules) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AlertingRule, err error) {
	result = &v1alpha1.AlertingRule{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("alertingrules").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
