#!/bin/bash

# SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

set -e

CURRENT_DIR=$(dirname $0)
PROJECT_ROOT="${CURRENT_DIR}"/..

# On MacOS there is a strange race condition
# between port allocation of envtest suites when go test
# runs all the tests in parallel without any limits (spins up around 10+ environments).
#
# To avoid flakes, set we're setting the go-test parallel flag
# to limit the number of parallel executions.
#
# TODO: check the controller-runtime for root-cause and real mitigation
# https://github.com/kubernetes-sigs/controller-runtime/pull/1567
if [[ "${OSTYPE}" == "darwin"* ]]; then
  P_FLAG="-p=1"
fi

go test -mod=vendor -race ${P_FLAG} ${PROJECT_ROOT}/cmd/... \
                    ${PROJECT_ROOT}/pkg/...
