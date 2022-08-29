#!/bin/bash

# SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

set -e

K8S_VERSION="1.21.x"

echo "> Setup Test Environment for K8s Version ${K8S_VERSION}"

CURRENT_DIR=$(dirname $0)
PROJECT_ROOT="${CURRENT_DIR}"/..

# TODO: setup-envtest currently doesnt support darwin/arm64 / force amd64
ARCH_ARG=""
if [[ $(go env GOOS) == "darwin" && $(go env GOARCH) == "arm64" ]]; then
  ARCH_ARG="--arch amd64"
fi

export KUBEBUILDER_ASSETS=$(setup-envtest use -p path ${K8S_VERSION} ${ARCH_ARG})

mkdir -p ${PROJECT_ROOT}/tmp/test
rm -f ${PROJECT_ROOT}/tmp/test/bin
ln -s "${KUBEBUILDER_ASSETS}" ${PROJECT_ROOT}/tmp/test/bin

# TODO: The landscaper crd files used for testing are currently not exported via landscaper api module.
#       To avoid adding the landscaper module, download the needed crd files directly.
LANDSCAPER_APIS_VERSION=$(go list -m -mod=readonly -f {{.Version}}  github.com/gardener/landscaper/apis)
LANDSCAPER_CRD_URL="https://raw.githubusercontent.com/gardener/landscaper/${LANDSCAPER_APIS_VERSION}/pkg/landscaper/crdmanager/crdresources"
LANDSCAPER_CRD_DIR="${PROJECT_ROOT}/tmp/landscapercrd"
LANDSCAPER_CRDS="landscaper.gardener.cloud_installations.yaml landscaper.gardener.cloud_targets.yaml landscaper.gardener.cloud_dataobjects.yaml landscaper.gardener.cloud_contexts.yaml landscaper.gardener.cloud_lshealthchecks.yaml"
mkdir -p ${PROJECT_ROOT}/tmp/landscapercrd

for crd in $LANDSCAPER_CRDS; do
  (cd ${LANDSCAPER_CRD_DIR} && curl -s -O "$LANDSCAPER_CRD_URL/$crd")
done
