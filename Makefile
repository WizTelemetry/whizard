REPO ?= kubesphere
TAG ?= latest

CONTROLLER_MANAGER_IMG=${REPO}/whizard-controller-manager:${TAG}
MONITORING_GATEWAY_IMG=${REPO}/whizard-monitoring-gateway:${TAG}
MONITORING_AGENT_PROXY_IMG=${REPO}/whizard-monitoring-agent-proxy:${TAG}
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

GV="monitoring:v1alpha1"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

all: build

##@ Development

manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: manifests generate fmt vet ## Run tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out

##@ Build

build: controller-manager monitoring-gateway monitoring-agent-proxy

controller-manager: 
	go build -o bin/controller-manager cmd/controller-manager/controller-manager.go

monitoring-gateway: 
	go build -o bin/monitoring-gateway cmd/monitoring-gateway/monitoring-gateway.go

monitoring-agent-proxy:
	go build -o bin/monitoring-agent-proxy cmd/monitoring-agent-proxy/monitoring-agent-proxy.go

docker-build: docker-build-controller-manager docker-build-monitoring-gateway docker-build-monitoring-agent-proxy

docker-build-controller-manager: 
	docker build -t $(CONTROLLER_MANAGER_IMG) -f build/controller-manager/Dockerfile .

docker-build-monitoring-gateway:
	docker build -t $(MONITORING_GATEWAY_IMG) -f build/monitoring-gateway/Dockerfile .

docker-build-monitoring-agent-proxy:
	docker build -t $(MONITORING_AGENT_PROXY_IMG) -f build/monitoring-agent-proxy/Dockerfile .

##@ Deployment

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete -f -


CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

clientset:
	./hack/generate_client.sh $(GV)

bundle: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(CONTROLLER_MANAGER_IMG)
	cd config/default && $(KUSTOMIZE) edit set namespace kubesphere-monitoring-system
	$(KUSTOMIZE) build config/default > config/bundle.yaml

MONITORING_TYPE_GOES=$(shell find pkg/api/monitoring -name *_types.go | tr '\n' ' ')
docs/monitoring/api.md: tools/docgen/docgen.go $(TYPE_GOES)
	go run github.com/kubesphere/whizard/tools/docgen $(MONITORING_TYPE_GOES) > docs/monitoring/api.md

whizard.key:
	openssl genrsa -des3 -passout pass:x -out whizard.pass.key 2048
	openssl rsa -passin pass:x -in whizard.pass.key -out whizard.key
	rm whizard.pass.key

whizard.crt: whizard.key
	openssl req -new -key whizard.key -out whizard.csr
	openssl x509 -req -sha256 -days 365 -in whizard.csr -signkey whizard.key -out whizard.crt
	rm whizard.csr