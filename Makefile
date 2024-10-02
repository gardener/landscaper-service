# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

REPO_ROOT                                      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION                                        := $(shell cat $(REPO_ROOT)/VERSION)
EFFECTIVE_VERSION                              := $(shell $(REPO_ROOT)/hack/get-version.sh)

REGISTRY                                       := europe-docker.pkg.dev/sap-gcp-cp-k8s-stable-hub/landscaper

CODE_DIRS := $(REPO_ROOT)/cmd/... $(REPO_ROOT)/pkg/... $(REPO_ROOT)/test/... $(REPO_ROOT)/integration-test/...

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Development

.PHONY: revendor
revendor: ## Runs 'go mod tidy' for all go modules in this repo.
	@$(REPO_ROOT)/hack/revendor.sh

.PHONY: format
format: goimports ## Runs the formatter.
	@@FORMATTER=$(FORMATTER) $(REPO_ROOT)/hack/format.sh $(CODE_DIRS)

.PHONY: check
check: golangci-lint goimports ## Runs linter, 'go vet', and checks if the formatter has been run.
	@LINTER=$(LINTER) FORMATTER=$(FORMATTER) $(REPO_ROOT)/hack/check.sh --golangci-lint-config="$(REPO_ROOT)/.golangci.yaml" $(CODE_DIRS)

.PHONY: verify
verify: check ## Alias for 'make check'.

.PHONY: generate-code
generate-code: code-gen controller-gen ## Runs code generation (deepcopy/conversion/defaulter functions, CRDs).
	@CODE_GEN_SCRIPT=$(CODE_GEN_SCRIPT) CONTROLLER_GEN=$(CONTROLLER_GEN) $(REPO_ROOT)/hack/generate-code.sh

.PHONY: generate # Runs code and docs generation and the formatter.
generate: generate-code format


##@ Build

PLATFORMS ?= linux/arm64,linux/amd64

.PHONY: build
build: ## Build binaries for all os/arch combinations specified in PLATFORMS.
	@PLATFORMS=$(PLATFORMS) COMPONENT=landscaper-service-controller $(REPO_ROOT)/hack/build.sh
	@PLATFORMS=$(PLATFORMS) COMPONENT=landscaper-service-webhooks-server $(REPO_ROOT)/hack/build.sh
	@PLATFORMS=$(PLATFORMS) COMPONENT=landscaper-service-target-shoot-sidecar-server $(REPO_ROOT)/hack/build.sh
	
.PHONY: docker-images
docker-images: build ## Builds images for all controllers locally. The images are suffixed with -$OS-$ARCH
	@PLATFORMS=$(PLATFORMS) $(REPO_ROOT)/hack/docker-build-multi.sh

.PHONY: component
component: ocm ## Builds and pushes the Component Descriptor. Also pushes the images and combines them into multi-platform images. Requires the docker images to have been built before.
	@OCM=$(OCM) $(REPO_ROOT)/hack/generate-cd.sh $(REGISTRY)

.PHONY: build-resources ## Wrapper for 'make docker-images component'.
build-resources: docker-images component

.PHONY: build-int-test-image
build-int-test-image:
	- docker buildx create --name project-v3-builder
	docker buildx use project-v3-builder
	@docker buildx build --platform linux/amd64 integration-test/docker -t europe-docker.pkg.dev/sap-gcp-cp-k8s-stable-hub/landscaper/integration-test:1.22.4-alpine3.19 --push
	- docker buildx rm project-v3-builder


##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(REPO_ROOT)/bin

## Tool Binaries
CODE_GEN_SCRIPT ?= $(LOCALBIN)/kube_codegen.sh
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
FORMATTER ?= $(LOCALBIN)/goimports
LINTER ?= $(LOCALBIN)/golangci-lint
OCM ?= $(LOCALBIN)/ocm

## Tool Versions
CODE_GEN_VERSION ?= $(shell  $(REPO_ROOT)/hack/extract-module-version.sh k8s.io/code-generator)
# renovate: datasource=github-releases depName=kubernetes-sigs/controller-tools
CONTROLLER_TOOLS_VERSION ?= v0.16.3
# renovate: datasource=github-tags depName=golang/tools
FORMATTER_VERSION ?= v0.25.0
# renovate: datasource=github-releases depName=golangci/golangci-lint
LINTER_VERSION ?= v1.61.0
# renovate: datasource=github-releases depName=open-component-model/ocm
OCM_VERSION ?= v0.13.0

.PHONY: localbin
localbin: ## Creates the local bin folder, if it doesn't exist. Not meant to be called manually, used as requirement for the other tool commands.
	@test -d $(LOCALBIN) || mkdir -p $(LOCALBIN)

.PHONY: code-gen
code-gen: localbin ## Download the code-gen script locally.
	@test -s $(CODE_GEN_SCRIPT) && test -s $(LOCALBIN)/kube_codegen_version && cat $(LOCALBIN)/kube_codegen_version | grep -q $(CODE_GEN_VERSION) || \
	( echo "Downloading code generator script $(CODE_GEN_VERSION) ..."; \
	curl -sfL "https://raw.githubusercontent.com/kubernetes/code-generator/$(CODE_GEN_VERSION)/kube_codegen.sh" --output "$(CODE_GEN_SCRIPT)" && chmod +x "$(CODE_GEN_SCRIPT)" && \
	echo $(CODE_GEN_VERSION) > $(LOCALBIN)/kube_codegen_version )

.PHONY: controller-gen
controller-gen: localbin ## Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
	@test -s $(CONTROLLER_GEN) && $(CONTROLLER_GEN) --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	( echo "Installing controller-gen $(CONTROLLER_TOOLS_VERSION) ..."; \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION) )

.PHONY: goimports
goimports: localbin ## Download goimports locally if necessary. If wrong version is installed, it will be overwritten.
	@test -s $(FORMATTER) && test -s $(LOCALBIN)/goimports_version && cat $(LOCALBIN)/goimports_version | grep -q $(FORMATTER_VERSION) || \
	( echo "Installing goimports $(FORMATTER_VERSION) ..."; \
	GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@$(FORMATTER_VERSION) && \
	echo $(FORMATTER_VERSION) > $(LOCALBIN)/goimports_version )

.PHONY: golangci-lint
golangci-lint: localbin ## Download golangci-lint locally if necessary. If wrong version is installed, it will be overwritten.
	@test -s $(LINTER) && $(LINTER) --version | grep -q $(LINTER_VERSION) || \
	( echo "Installing golangci-lint $(LINTER_VERSION) ..."; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) $(LINTER_VERSION) )

.PHONY: ocm
ocm: localbin ## Install OCM CLI if necessary. If wrong version is installed, it will be overwritten.
	@test -s $(OCM) && $(OCM) --version | grep -q $(subst v,,$(OCM_VERSION)) || \
	( echo "Installing OCM tooling $(OCM_VERSION) ..."; \
	curl -sSfL https://ocm.software/install.sh | OCM_VERSION=$(subst v,,$(OCM_VERSION)) bash -s $(LOCALBIN) )
