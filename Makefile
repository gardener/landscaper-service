# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

REPO_ROOT                                      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION                                        := $(shell cat $(REPO_ROOT)/VERSION)
EFFECTIVE_VERSION                              := $(VERSION)-$(shell git rev-parse HEAD)

REGISTRY                                       := eu.gcr.io/gardener-project/landscaper-service
LANDSCAPER_SERVICE_CONTROLLER_IMAGE_REPOSITORY         := $(REGISTRY)/landscaper-service-controller
LANDSCAPER_SERVICE_WEBHOOKS_SERVER_IMAGE_REPOSITORY    := $(REGISTRY)/landscaper-service-webhooks-server
LANDSCAPER_SERVICE_TARGET_SHOOT_SIDECAR_SERVER_IMAGE_REPOSITORY    := $(REGISTRY)/landscaper-service-target-shoot-sidecar-server

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
	@echo "Building docker images for version $(EFFECTIVE_VERSION)"
	@docker build --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) -t $(LANDSCAPER_SERVICE_CONTROLLER_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION) -f Dockerfile --target landscaper-service-controller .
	@docker build --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) -t $(LANDSCAPER_SERVICE_WEBHOOKS_SERVER_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION) -f Dockerfile --target landscaper-service-webhooks-server .
	@docker build --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) -t $(LANDSCAPER_SERVICE_TARGET_SHOOT_SIDECAR_SERVER_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION) -f Dockerfile --target landscaper-service-target-shoot-sidecar-server .

.PHONY: docker-push
docker-push:
	@echo "Pushing docker images for version $(EFFECTIVE_VERSION) to registry $(REGISTRY)"
	@if ! docker images $(LANDSCAPER_SERVICE_CONTROLLER_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(EFFECTIVE_VERSION); then echo "$(LANDSCAPER_SERVICE_CONTROLLER_IMAGE_REPOSITORY) version $(EFFECTIVE_VERSION) is not yet built. Please run 'make docker-images'"; false; fi
	@if ! docker images $(LANDSCAPER_SERVICE_WEBHOOKS_SERVER_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(EFFECTIVE_VERSION); then echo "$(LANDSCAPER_SERVICE_WEBHOOKS_SERVER_IMAGE_REPOSITORY) version $(EFFECTIVE_VERSION) is not yet built. Please run 'make docker-images'"; false; fi
	@if ! docker images $(LANDSCAPER_SERVICE_TARGET_SHOOT_SIDECAR_SERVER_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(EFFECTIVE_VERSION); then echo "$(LANDSCAPER_SERVICE_TARGET_SHOOT_SIDECAR_SERVER_IMAGE_REPOSITORY) version $(EFFECTIVE_VERSION) is not yet built. Please run 'make docker-images'"; false; fi
	@docker push $(LANDSCAPER_SERVICE_CONTROLLER_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION)
	@docker push $(LANDSCAPER_SERVICE_WEBHOOKS_SERVER_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION)
	@docker push $(LANDSCAPER_SERVICE_TARGET_SHOOT_SIDECAR_SERVER_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION)

.PHONY: docker-all
docker-all: docker-images docker-push

.PHONY: cnudie
cnudie:
	@$(REPO_ROOT)/hack/generate-cd.sh

.PHONY: helm-charts
helm-charts:
	@$(REPO_ROOT)/.ci/publish-helm-charts

.PHONY: build-resources
build-resources: docker-all helm-charts cnudie

.PHONY: build-int-test-image
build-int-test-image:
	@docker buildx build --platform linux/amd64 integration-test/docker -t eu.gcr.io/gardener-project/landscaper-service/integration-test:1.20.10-alpine3.18 --push
