# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

REPO_ROOT                                      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION                                        := $(shell cat $(REPO_ROOT)/VERSION)
EFFECTIVE_VERSION                              := $(shell $(REPO_ROOT)/hack/get-version.sh)

REGISTRY                                       := europe-docker.pkg.dev/sap-gcp-cp-k8s-stable-hub/landscaper

DOCKER_BUILDER_NAME := "laas-multiarch"
DOCKER_PLATFORM     := "linux/amd64"

.PHONY: install-requirements
install-requirements:
	@go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
	@$(REPO_ROOT)/hack/install-requirements.sh

.PHONY: revendor
revendor:
	@$(REPO_ROOT)/hack/revendor.sh
	@cd $(REPO_ROOT)/integration-test && $(REPO_ROOT)/hack/revendor.sh

.PHONY: format
format:
	@$(REPO_ROOT)/hack/format.sh $(REPO_ROOT)/pkg $(REPO_ROOT)/cmd $(REPO_ROOT)/hack $(REPO_ROOT)/test $(REPO_ROOT)/integration-test/pkg

.PHONY: check
check: revendor check-fast

.PHONY: check-fast
check-fast:
	@$(REPO_ROOT)/hack/check.sh --golangci-lint-config=./.golangci.yaml $(REPO_ROOT)/cmd/... $(REPO_ROOT)/pkg/... $(REPO_ROOT)/hack/... $(REPO_ROOT)/test/...
	@cd $(REPO_ROOT)/integration-test && $(REPO_ROOT)/hack/check.sh --golangci-lint-config=$(REPO_ROOT)/.golangci.yaml ./pkg/...

.PHONY: verify
verify: check

.PHONY: setup-testenv
setup-testenv:
	@$(REPO_ROOT)/hack/setup-testenv.sh

.PHONY: test
test: setup-testenv
	@$(REPO_ROOT)/hack/test.sh

.PHONY: integration-test
integration-test:
	@cd $(REPO_ROOT)/integration-test && go run ./pkg --kubeconfig $(KUBECONFIG) --laas-version $(EFFECTIVE_VERSION) --laas-repository $(REGISTRY) --provider-type $(CLUSTER_PROVIDER_TYPE)

.PHONY: generate-code
generate-code:
	$(REPO_ROOT)/hack/generate-code.sh ./... && $(REPO_ROOT)/hack/generate-crd.sh

.PHONY: generate
generate: generate-code format revendor

#################################################################
# Rules related to binary build, docker image build and release #
#################################################################

.PHONY: install
install:
	@EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) ./hack/install.sh

.PHONY: docker-images
docker-images:
	@$(REPO_ROOT)/hack/prepare-docker-builder.sh
	@echo "Building docker images for version $(EFFECTIVE_VERSION)"
	@docker buildx build --builder $(DOCKER_BUILDER_NAME) --load --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) --platform $(DOCKER_PLATFORM) -t landscaper-service-controller:$(EFFECTIVE_VERSION) -f Dockerfile --target landscaper-service-controller .
	@docker buildx build --builder $(DOCKER_BUILDER_NAME) --load --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) --platform $(DOCKER_PLATFORM) -t landscaper-service-webhooks-server:$(EFFECTIVE_VERSION) -f Dockerfile --target landscaper-service-webhooks-server .
	@docker buildx build --builder $(DOCKER_BUILDER_NAME) --load --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) --platform $(DOCKER_PLATFORM) -t landscaper-service-target-shoot-sidecar-server:$(EFFECTIVE_VERSION) -f Dockerfile --target landscaper-service-target-shoot-sidecar-server .

.PHONY: component
component:
	@$(REPO_ROOT)/hack/generate-cd.sh $(REGISTRY)

.PHONY: build-resources
build-resources: docker-images component

.PHONY: build-int-test-image
build-int-test-image:
	@docker buildx build --platform linux/amd64 integration-test/docker -t europe-docker.pkg.dev/sap-gcp-cp-k8s-stable-hub/landscaper/integration-test:1.21.7-alpine3.18 --push
