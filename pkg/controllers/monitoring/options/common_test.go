package options

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestCommonOptionsApplyTo(t *testing.T) {
	var replicas1 int32 = 1
	var replicas2 int32 = 2

	testCases := []struct {
		name    string
		options CommonOptions
		conf    CommonOptions
		want    *CommonOptions
	}{
		{
			"good case 1",
			NewCommonOptions(),
			CommonOptions{},
			&CommonOptions{
				Image:    DefaultWhizardImage,
				Replicas: &replicas1,
			},
		},
		{
			"good case 2",
			NewCommonOptions(),
			CommonOptions{
				Image:    "thanos/thanos:v0.28.0",
				Replicas: &replicas2,
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup: &[]int64{0}[0],
				},
			},
			&CommonOptions{
				Image:    "thanos/thanos:v0.28.0",
				Replicas: &replicas2,
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup: &[]int64{0}[0],
				},
			},
		},
		{
			"good case 3",
			CommonOptions{
				Image:    DefaultWhizardImage,
				Replicas: &replicas2,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("500Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("1"),
						corev1.ResourceMemory: resource.MustParse("4Gi"),
					},
				},
			},
			CommonOptions{
				Image:    "thanos/thanos:v0.28.0",
				Replicas: &replicas1,
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("2"),
						corev1.ResourceMemory: resource.MustParse("8Gi"),
					},
				},
			},
			&CommonOptions{
				Image:    "thanos/thanos:v0.28.0",
				Replicas: &replicas1,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("500Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("2"),
						corev1.ResourceMemory: resource.MustParse("8Gi"),
					},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.conf.ApplyTo(&tt.options)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nget = %+v, \nwant %+v", got, tt.want)
			}
		})
	}
}

func TestSidecarOptionsApplyTo(t *testing.T) {

	testCases := []struct {
		name    string
		options SidecarOptions
		conf    SidecarOptions
		want    *SidecarOptions
	}{
		{
			"good case 1",
			SidecarOptions{
				Image: DefaultEnvoyImage,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("500Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("500Mi"),
					},
				}},
			SidecarOptions{},
			&SidecarOptions{
				Image: DefaultEnvoyImage,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("500Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("500Mi"),
					},
				},
			},
		},
		{
			"good case 2",
			SidecarOptions{
				Image: DefaultRulerWriteProxyImage,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("50Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("200m"),
						corev1.ResourceMemory: resource.MustParse("400Mi"),
					},
				},
			},
			SidecarOptions{
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("200m"),
						corev1.ResourceMemory: resource.MustParse("2Gi"),
					},
				},
			},
			&SidecarOptions{
				Image: DefaultRulerWriteProxyImage,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("50Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("200m"),
						corev1.ResourceMemory: resource.MustParse("2Gi"),
					},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.conf.ApplyTo(&tt.options)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nget = %+v, \nwant %+v", got, tt.want)
			}
		})
	}
}
