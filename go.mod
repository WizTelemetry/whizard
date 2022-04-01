module github.com/kubesphere/paodin-monitoring

go 1.16

require (
	github.com/alecthomas/kong v0.5.0
	github.com/go-kit/log v0.2.0
	github.com/go-logr/logr v1.0.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/prometheus-community/prom-label-proxy v0.4.1-0.20220310103857-b961d28b26ab
	github.com/prometheus/common v0.32.1
	github.com/prometheus/prometheus v1.8.2-0.20211214150951-52c693a63be1
	github.com/thanos-io/thanos v0.25.2
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.22.4
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.3
	k8s.io/code-generator v0.21.3
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a
	sigs.k8s.io/controller-runtime v0.9.5
)

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.8.2-0.20211119115433-692a54649ed7
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.8.0
)
