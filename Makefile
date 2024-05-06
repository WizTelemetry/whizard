REPO ?= kubesphere
TAG ?= $(shell cat VERSION | tr -d " \t\n\r")
VERSION=$(shell cat VERSION | tr -d " \t\n\r")

CONTROLLER_MANAGER_IMG=${REPO}/whizard-controller-manager:${TAG}
MONITORING_GATEWAY_IMG=${REPO}/whizard-monitoring-gateway:${TAG}
MONITORING_AGENT_PROXY_IMG=${REPO}/whizard-monitoring-agent-proxy:${TAG}
MONITORING_BLOCK_MANAGER_IMG=${REPO}/whizard-monitoring-block-manager:${TAG}


# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:allowDangerousTypes=true"

GV="monitoring:v1alpha1"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)


CONTAINER_CLI ?=docker
# buildx build
CONTAINER_BUILDER ?= build
# --platform=linux/amd64,linux/arm64 --push
CONTAINER_BUILD_EXTRA_ARGS ?= 


BUILD_DATE=$(shell date +"%Y%m%d-%T")
# source: https://docs.github.com/en/free-pro-team@latest/actions/reference/environment-variables#default-environment-variables
ifndef GITHUB_ACTIONS
	BUILD_USER?=$(USER)
	BUILD_BRANCH?=$(shell git branch --show-current)
	BUILD_REVISION?=$(shell git rev-parse --short HEAD)
else
	BUILD_USER=Action-Run-ID-$(GITHUB_RUN_ID)
	BUILD_BRANCH=$(GITHUB_REF:refs/heads/%=%)
	BUILD_REVISION=$(GITHUB_SHA)
endif


# The Prometheus common library import path
PROMETHEUS_COMMON_PKG=github.com/prometheus/common

# The ldflags for the go build process to set the version related data.
GO_BUILD_LDFLAGS= \
	-X $(PROMETHEUS_COMMON_PKG)/version.Revision=$(BUILD_REVISION)  \
	-X $(PROMETHEUS_COMMON_PKG)/version.BuildUser=$(BUILD_USER) \
	-X $(PROMETHEUS_COMMON_PKG)/version.BuildDate=$(BUILD_DATE) \
	-X $(PROMETHEUS_COMMON_PKG)/version.Branch=$(BUILD_BRANCH) \
	-X $(PROMETHEUS_COMMON_PKG)/version.Version=$(VERSION)


GO_BUILD_RECIPE=\
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	CGO_ENABLED=0 \
	go build -ldflags="$(GO_BUILD_LDFLAGS)"

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

clientset:
	./hack/generate_client.sh $(GV)

bundle: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(CONTROLLER_MANAGER_IMG)
	cd config/default && $(KUSTOMIZE) edit set namespace kubesphere-monitoring-system
	$(KUSTOMIZE) build config/default > config/bundle.yaml

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

.PHONY: build
build: controller-manager monitoring-gateway monitoring-block-manager

controller-manager: 
	$(GO_BUILD_RECIPE)  -o bin/controller-manager cmd/controller-manager/controller-manager.go

monitoring-gateway: 
	$(GO_BUILD_RECIPE) -o bin/monitoring-gateway cmd/monitoring-gateway/*

#monitoring-agent-proxy:
#	go build -o bin/monitoring-agent-proxy cmd/monitoring-agent-proxy/monitoring-agent-proxy.go

monitoring-block-manager:
	go build -o bin/monitoring-block-manager cmd/monitoring-block-manager/block-manager.go

docker-build: docker-build-controller-manager docker-build-monitoring-gateway docker-build-monitoring-agent-proxy docker-build-monitoring-block-manager

docker-build-controller-manager: 
	${CONTAINER_CLI} ${CONTAINER_BUILDER} ${CONTAINER_BUILD_EXTRA_ARGS} --build-arg GOLDFLAGS="$(GO_BUILD_LDFLAGS)" -t $(CONTROLLER_MANAGER_IMG) -f build/controller-manager/Dockerfile .

docker-build-monitoring-gateway:
	${CONTAINER_CLI} ${CONTAINER_BUILDER} ${CONTAINER_BUILD_EXTRA_ARGS} --build-arg GOLDFLAGS="$(GO_BUILD_LDFLAGS)"  -t $(MONITORING_GATEWAY_IMG) -f build/monitoring-gateway/Dockerfile .

docker-build-monitoring-agent-proxy:
	${CONTAINER_CLI} ${CONTAINER_BUILDER} ${CONTAINER_BUILD_EXTRA_ARGS} --build-arg GOLDFLAGS="$(GO_BUILD_LDFLAGS)"  -t $(MONITORING_AGENT_PROXY_IMG) -f build/monitoring-agent-proxy/Dockerfile .

docker-build-monitoring-block-manager:
	${CONTAINER_CLI} ${CONTAINER_BUILDER} ${CONTAINER_BUILD_EXTRA_ARGS}  -t $(MONITORING_BLOCK_MANAGER_IMG) -f build/monitoring-block-manager/Dockerfile .

##@ Deployment

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply --server-side  -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${CONTROLLER_MANAGER_IMG}
	$(KUSTOMIZE) build config/default | kubectl apply --server-side -f -

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete -f -


CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

CRD-REF-DOCS =$(shell pwd)/bin/crd-ref-docs
crd-ref-docs: ## Download crd-ref-docs locally if necessary.
	$(call go-get-tool,$(CRD-REF-DOCS),github.com/elastic/crd-ref-docs@v0.0.12)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef


docgen: crd-ref-docs
	$(CRD-REF-DOCS) \
	    --source-path=./pkg/api/monitoring/v1alpha1 \
	    --config=./tools/docgen/config.yaml \
		--output-path=./docs/monitoring/api.md \
	    --renderer=markdown

# stripped-down-crds is a version of the whizard CRDs with all
# description fields being removed. It is meant as a workaround for the issue
# that `kubectl apply -f ...` might fail with the full version of the CRDs
# because of too long annotations field.
# See https://github.com/prometheus-operator/prometheus-operator/issues/4355
stripped-down-crds: manifests
	cd config/crd/bases && \
	for f in *.yaml; do \
		echo "---" > ../../../charts/whizard/crds/$$f; \
		gojsontoyaml -yamltojson < $$f | jq 'walk(if type == "object" then with_entries(if .value|type=="object" then . else select(.key | test("description") | not) end) else . end)' | gojsontoyaml >> ../../../charts/whizard/crds/$$f; \
	done && \
	cp -rf ../../../charts/whizard/crds ../../../charts/whizard-crds/;

cut-new-version: stripped-down-crds ## Update appVersion in helm chart.
	$(shell sed -i '' -e 's/appVersion:.*/appVersion: "'$(VERSION)'"/g'  charts/whizard/Chart.yaml)
	   