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

go install "${PROJECT_ROOT}"/vendor/sigs.k8s.io/controller-tools/cmd/controller-gen

echo "Generating OpenAPI specification fore core ..."
controller-gen crd paths=${PROJECT_ROOT}/pkg/apis/core/... output:crd:dir=${PROJECT_ROOT}/pkg/crdmanager/crdresources
