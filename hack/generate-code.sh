#!/bin/bash
#
# Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

rm -f ${GOPATH}/bin/deepcopy-gen
rm -f ${GOPATH}/bin/defaulter-gen
rm -f ${GOPATH}/bin/conversion-gen

PROJECT_MOD_ROOT="github.com/gardener/landscaper-service"

CURRENT_DIR=$(dirname $0)
PROJECT_ROOT="${CURRENT_DIR}"/..

chmod +x ${PROJECT_ROOT}/vendor/k8s.io/code-generator/*

export GOFLAGS=-mod=vendor

echo "> Generating groups for Landscaper Service"
bash "${PROJECT_ROOT}"/vendor/k8s.io/code-generator/generate-internal-groups.sh \
  deepcopy,defaulter,conversion \
  $PROJECT_MOD_ROOT/pkg/generated \
  $PROJECT_MOD_ROOT/pkg/apis \
  $PROJECT_MOD_ROOT/pkg/apis \
  "core:v1alpha1" \
  --go-header-file "${PROJECT_ROOT}/hack/boilerplate.go.txt"
