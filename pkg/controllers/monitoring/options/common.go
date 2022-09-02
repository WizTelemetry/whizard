package options

import (
	"fmt"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/util"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ResourceList map[string]string

// type ResourceList map[ResourceName]resource.Quantity

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Limits ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList,castkey=ResourceName"`
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Requests ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList,castkey=ResourceName"`
}

type CommonOptions struct {
	Image           string                      `json:"image,omitempty" yaml:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy           `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	Affinity        *corev1.Affinity            `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	NodeSelector    map[string]string           `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	Tolerations     []corev1.Toleration         `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	InnerResources  ResourceRequirements        `json:"resources,omitempty" yaml:"resources,omitempty" mapstructure:"resources,omitempty"`
	Resources       corev1.ResourceRequirements `json:"-" yaml:"-" mapstructure:"-"`
	//Resources      corev1.ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty" mapstructure:"resources,omitempty"`
	Replicas   *int32                     `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	LogLevel   string                     `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	LogFormat  string                     `json:"logFormat,omitempty" yaml:"logFormat,omitempty"`
	Flags      []string                   `json:"flags,omitempty" yaml:"flags,omitempty"`
	DataVolume *v1alpha1.KubernetesVolume `json:"dataVolume,omitempty" yaml:"dataVolume,omitempty"`
}

func NewCommonOptions() CommonOptions {
	var replicas int32 = 1
	return CommonOptions{
		Image:    DefaultWhizardImage,
		Replicas: &replicas,
	}
}
func (o *CommonOptions) Validate() []error {
	var errs []error

	if o.Replicas != nil && *o.Replicas < 0 {
		errs = append(errs, fmt.Errorf("replicas must be >= 0"))
	}
	if o.Resources.Limits == nil {
		o.Resources.Limits = make(corev1.ResourceList)
	}
	if o.Resources.Requests == nil {
		o.Resources.Requests = make(corev1.ResourceList)
	}

	if v, ok := o.InnerResources.Requests[string(corev1.ResourceCPU)]; ok {
		q, err := resource.ParseQuantity(v)
		o.Resources.Requests[corev1.ResourceCPU] = q
		errs = append(errs, err)
	}
	if v, ok := o.InnerResources.Requests[string(corev1.ResourceMemory)]; ok {
		q, err := resource.ParseQuantity(v)
		o.Resources.Requests[corev1.ResourceMemory] = q
		errs = append(errs, err)
	}

	if v, ok := o.InnerResources.Limits[string(corev1.ResourceCPU)]; ok {
		q, err := resource.ParseQuantity(v)
		o.Resources.Limits[corev1.ResourceCPU] = q
		errs = append(errs, err)
	}
	if v, ok := o.InnerResources.Limits[string(corev1.ResourceMemory)]; ok {
		q, err := resource.ParseQuantity(v)
		o.Resources.Limits[corev1.ResourceMemory] = q
		errs = append(errs, err)
	}
	return errs
}

func (o *CommonOptions) ApplyTo(options *CommonOptions) {
	if o.Image != "" {
		options.Image = o.Image
	}

	if o.ImagePullPolicy != "" {
		options.ImagePullPolicy = o.ImagePullPolicy
	}

	if o.Affinity != nil {
		if options.Affinity == nil {
			options.Affinity = o.Affinity
		}

		util.Override(options.Affinity, o.Affinity)
	}

	if o.Tolerations != nil {
		options.Tolerations = o.Tolerations
	}

	if o.NodeSelector != nil {
		options.NodeSelector = o.NodeSelector
	}

	if o.InnerResources.Limits != nil {
		options.InnerResources.Limits = o.InnerResources.Limits
	}
	if o.InnerResources.Requests != nil {
		options.InnerResources.Requests = o.InnerResources.Requests
	}

	if o.Resources.Limits != nil {
		if options.Resources.Limits == nil {
			options.Resources.Limits = o.Resources.Limits
		}
		for k, v := range o.Resources.Limits {
			options.Resources.Limits[k] = v
		}
	}

	if o.Resources.Requests != nil {
		if options.Resources.Requests == nil {
			options.Resources.Requests = o.Resources.Requests
		}
		for k, v := range o.Resources.Requests {
			options.Resources.Requests[k] = v
		}
	}

	if o.Replicas != nil {
		options.Replicas = o.Replicas
	}

	if o.LogLevel != "" {
		options.LogLevel = o.LogLevel
	}

	if o.LogFormat != "" {
		options.LogFormat = o.LogFormat
	}

	if o.Flags != nil {
		options.Flags = o.Flags
	}

	if o.DataVolume != nil {
		if options.DataVolume == nil {
			options.DataVolume = o.DataVolume
		}

		if o.DataVolume.PersistentVolumeClaim != nil {
			options.DataVolume.PersistentVolumeClaim = o.DataVolume.PersistentVolumeClaim
		}

		if o.DataVolume.EmptyDir != nil {
			options.DataVolume.EmptyDir = o.DataVolume.EmptyDir
		}
	}
}

func (o *CommonOptions) Apply(spec *v1alpha1.CommonSpec) {
	if spec.Image == "" {
		spec.Image = o.Image
	}

	if spec.ImagePullPolicy == "" {
		spec.ImagePullPolicy = o.ImagePullPolicy
	}
	if spec.Replicas == nil || *spec.Replicas < 0 {
		spec.Replicas = o.Replicas
	}

	if spec.Affinity == nil {
		spec.Affinity = o.Affinity
	}

	if spec.Tolerations == nil {
		spec.Tolerations = o.Tolerations
	}

	if spec.NodeSelector == nil {
		spec.NodeSelector = o.NodeSelector
	}

	if spec.Resources.Limits == nil {
		spec.Resources.Limits = o.Resources.Limits
	}

	if spec.Resources.Requests == nil {
		spec.Resources.Requests = o.Resources.Requests
	}

	if spec.Replicas == nil {
		spec.Replicas = o.Replicas
	}

	if spec.LogLevel == "" {
		spec.LogLevel = o.LogLevel
	}

	if spec.LogFormat == "" {
		spec.LogFormat = o.LogFormat
	}

	if spec.Flags == nil {
		spec.Flags = o.Flags
	}
}

func (o *CommonOptions) AddFlags(fs *pflag.FlagSet, c *CommonOptions, prefix string) {
	fs.StringVar(&c.Image, prefix+".image", c.Image, "Image with tag/version.")
	fs.StringArrayVar(&c.Flags, prefix+".flags", c.Flags, "Flags with --flag=value.")
	// fs.StringVar(&c.LogLevel, prefix+".log.level", c.LogLevel, "Log filtering level")
	// fs.StringVar(&c.LogFormat, prefix+".log.format", c.LogLevel, "Log format to use. Possible options: logfmt or json")
}

type ContainerOptions struct {
	// Image is the envoy image with tag/version
	Image string `json:"image,omitempty" yaml:"image,omitempty"`

	// Define resources requests and limits for envoy container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`
}

func (o *ContainerOptions) ApplyTo(options *ContainerOptions) {
	if o.Image != "" {
		options.Image = o.Image
	}
	if o.Resources.Limits != nil {
		if options.Resources.Limits == nil {
			options.Resources.Limits = o.Resources.Limits
		}
		for k, v := range o.Resources.Limits {
			options.Resources.Limits[k] = v
		}
	}

	if o.Resources.Requests != nil {
		if options.Resources.Requests == nil {
			options.Resources.Requests = o.Resources.Requests
		}
		for k, v := range o.Resources.Requests {
			options.Resources.Requests[k] = v
		}
	}
}

func (o *ContainerOptions) Apply(spec *v1alpha1.EnvoySpec) {
	if spec.Image == "" {
		spec.Image = o.Image
	}
	if spec.Resources.Limits == nil {
		spec.Resources.Limits = o.Resources.Limits
	}

	if spec.Resources.Requests == nil {
		spec.Resources.Requests = o.Resources.Requests
	}
}

func (o *ContainerOptions) AddFlags(fs *pflag.FlagSet, c *ContainerOptions, prefix string) {
	fs.StringVar(&c.Image, prefix+".image", c.Image, "Image with tag/version.")
}
