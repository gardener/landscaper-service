#!/bin/bash

# SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

set -e

K8S_VERSION="1.21.x"

echo "> Setup Test Environment for K8s Version ${K8S_VERSION}"

CURRENT_DIR=$(dirname $0)
PROJECT_ROOT="${CURRENT_DIR}"/..
export KUBEBUILDER_ASSETS=$(setup-envtest use -p path ${K8S_VERSION})
mkdir -p ${PROJECT_ROOT}/tmp/test
rm -f ${PROJECT_ROOT}/tmp/test/bin
ln -s "${KUBEBUILDER_ASSETS}" ${PROJECT_ROOT}/tmp/test/bin
