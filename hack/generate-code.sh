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

echo "> Generating openapi definitions"
go install "${PROJECT_ROOT}"/vendor/k8s.io/kube-openapi/cmd/openapi-gen
${GOPATH}/bin/openapi-gen "$@" \
  --v 1 \
  --logtostderr \
  --input-dirs=github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1 \
  --input-dirs=k8s.io/api/core/v1 \
  --input-dirs=k8s.io/apimachinery/pkg/apis/meta/v1 \
  --input-dirs=k8s.io/apimachinery/pkg/api/resource \
  --input-dirs=k8s.io/apimachinery/pkg/types \
  --input-dirs=k8s.io/apimachinery/pkg/runtime \
  --report-filename=${PROJECT_ROOT}/pkg/apis/openapi/api_violations.report \
  --output-package=github.com/gardener/landscaper-service/pkg/apis/openapi \
  -h "${PROJECT_ROOT}/hack/boilerplate.go.txt"
