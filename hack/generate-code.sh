#!/bin/bash
#
# Copyright (c) 2023 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# SPDX-License-Identifier: Apache-2.0


set -o errexit
set -o nounset
set -o pipefail

CURRENT_DIR=$(dirname $0)
PROJECT_ROOT="${CURRENT_DIR}"/..
PROJECT_MOD_ROOT="github.com/gardener/landscaper-service"

CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${PROJECT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

source "${PROJECT_ROOT}/${CODEGEN_PKG}/kube_codegen.sh"

echo "> Generating groups for Landscaper Service Core"
kube::codegen::gen_helpers \
  --input-pkg-root "${PROJECT_MOD_ROOT}/pkg/apis/core" \
  --output-base "${PROJECT_ROOT}/../../.." \
  --boilerplate "${PROJECT_ROOT}/hack/boilerplate.go.txt"

echo "> Generating groups for Landscaper Service Config"
kube::codegen::gen_helpers \
  --input-pkg-root "${PROJECT_MOD_ROOT}/pkg/apis/config" \
  --output-base "${PROJECT_ROOT}/../../.." \
  --boilerplate "${PROJECT_ROOT}/hack/boilerplate.go.txt"
