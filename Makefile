REPO ?= kubesphere
TAG ?= $(shell cat VERSION | tr -d " \t\n\r")
VERSION=$(shell cat VERSION | tr -d " \t\n\r")

# Image URL to use all building/pushing image targets
CONTROLLER_MANAGER_IMG=${REPO}/whizard-controller-manager:${TAG}
MONITORING_GATEWAY_IMG=${REPO}/whizard-monitoring-gateway:${TAG}
MONITORING_AGENT_PROXY_IMG=${REPO}/whizard-monitoring-agent-proxy:${TAG}
MONITORING_BLOCK_MANAGER_IMG=${REPO}/whizard-monitoring-block-manager:${TAG}

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.30.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)


# CONTAINER_CLI defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_CLI ?= docker
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
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

bundle: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(CONTROLLER_MANAGER_IMG)
	cd config/default && $(KUSTOMIZE) edit set namespace kubesphere-monitoring-system
	$(KUSTOMIZE) build config/default > config/bundle.yaml

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test $$(go list ./... | grep -v /e2e) -coverprofile cover.out

# Utilize Kind or modify the e2e tests to load the image locally, enabling compatibility with other vendors.
.PHONY: test-e2e  # Run the e2e tests against a Kind k8s instance that is spun up.
test-e2e:
	go test ./test/e2e/ -v -ginkgo.v

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix



##@ Build

.PHONY: build
build: manifests generate fmt vet controller-manager monitoring-gateway monitoring-block-manager ## Build manager binary.

controller-manager: 
	$(GO_BUILD_RECIPE)  -o bin/controller-manager cmd/controller-manager/controller-manager.go

monitoring-gateway: 
	$(GO_BUILD_RECIPE) -o bin/monitoring-gateway cmd/monitoring-gateway/*

monitoring-block-manager:
	go build -o bin/monitoring-block-manager cmd/monitoring-block-manager/block-manager.go


# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: docker-build-controller-manager docker-build-monitoring-gateway docker-build-monitoring-agent-proxy docker-build-monitoring-block-manager ## Build docker image with the manager.

docker-build-controller-manager: 
	${CONTAINER_CLI} ${CONTAINER_BUILDER} ${CONTAINER_BUILD_EXTRA_ARGS} --build-arg GOLDFLAGS="$(GO_BUILD_LDFLAGS)" -t $(CONTROLLER_MANAGER_IMG) -f build/controller-manager/Dockerfile .

docker-build-monitoring-gateway:
	${CONTAINER_CLI} ${CONTAINER_BUILDER} ${CONTAINER_BUILD_EXTRA_ARGS} --build-arg GOLDFLAGS="$(GO_BUILD_LDFLAGS)"  -t $(MONITORING_GATEWAY_IMG) -f build/monitoring-gateway/Dockerfile .

docker-build-monitoring-agent-proxy:
	${CONTAINER_CLI} ${CONTAINER_BUILDER} ${CONTAINER_BUILD_EXTRA_ARGS} --build-arg GOLDFLAGS="$(GO_BUILD_LDFLAGS)"  -t $(MONITORING_AGENT_PROXY_IMG) -f build/monitoring-agent-proxy/Dockerfile .

docker-build-monitoring-block-manager:
	${CONTAINER_CLI} ${CONTAINER_BUILDER} ${CONTAINER_BUILD_EXTRA_ARGS}  -t $(MONITORING_BLOCK_MANAGER_IMG) -f build/monitoring-block-manager/Dockerfile .

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) apply --server-side -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${CONTROLLER_MANAGER_IMG}
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply --server-side -f -

.PHONY: undeploy
undeploy: kustomize ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -


##@ Others

generate-api-docs: gen-crd-api-reference-docs manifests ## Generate API documentation
	$(API_DOC_GEN) -api-dir "github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1" -config "$(PWD)/tools/docs/config.json" -template-dir "$(PWD)/tools/docs/templates" -out-file "$(PWD)/docs/monitoring/api.md"

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

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize-$(KUSTOMIZE_VERSION)
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen-$(CONTROLLER_TOOLS_VERSION)
ENVTEST ?= $(LOCALBIN)/setup-envtest-$(ENVTEST_VERSION)
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
API_DOC_GEN ?= $(LOCALBIN)/gen-crd-api-reference-docs-$(API_DOC_GEN_VERSION)

## Tool Versions
KUSTOMIZE_VERSION ?= v5.4.1
CONTROLLER_TOOLS_VERSION ?= v0.17.2
ENVTEST_VERSION ?= release-0.18
GOLANGCI_LINT_VERSION ?= v1.57.2
API_DOC_GEN_VERSION ?= v0.3.0

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

.PHONY: envtest
envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,${GOLANGCI_LINT_VERSION})

.PHONY: gen-crd-api-reference-docs
gen-crd-api-reference-docs: $(API_DOC_GEN) ## Download gen-crd-api-reference-docs locally if necessary.
$(API_DOC_GEN): $(LOCALBIN)
	$(call go-install-tool,$(API_DOC_GEN),github.com/ahmetb/gen-crd-api-reference-docs,$(API_DOC_GEN_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef
