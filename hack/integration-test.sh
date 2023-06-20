#!/bin/sh

# SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

set +e
set -o pipefail

PROJECT_ROOT="$(dirname $0)/.."
ARG=$1

if [[ $ARG == "pull-request" ]]; then
  TEST_CLUSTER="laas-itest-pr"
  HOSTING_CLUSTER="laas-itest-tgt-pr"
  export TEST_PURPOSE="pull-request"
else
  TEST_CLUSTER="laas-itest-rel"
  HOSTING_CLUSTER="laas-itest-tg-rel"
  export TEST_PURPOSE="head-update-release"
fi

GARDENER_CLUSTER="laas-integration-test-service-account"
# LAAS_REPOSITORY="eu.gcr.io/sap-se-gcr-k8s-private/cnudie/gardener/development"
LAAS_REPOSITORY="eu.gcr.io/gardener-project/development"
LAAS_VERSION="$(${PROJECT_ROOT}/hack/get-version.sh)"
REPO_AUTH_URL="https://eu.gcr.io"
REPO_CTX_BASE_URL="eu.gcr.io/sap-se-gcr-k8s-private"
FULL_INTEGRATION_TEST_PATH="$(realpath $INTEGRATION_TEST_PATH)"

export PROJECT_ROOT
export TEST_CLUSTER
export HOSTING_CLUSTER
export GARDENER_CLUSTER
export TARGET_CLUSTER_PROVIDER
export LAAS_VERSION
export LAAS_REPOSITORY
export REPO_AUTH_URL
export REPO_CTX_BASE_URL
export FULL_INTEGRATION_TEST_PATH

unbuffer "${PROJECT_ROOT}/hack/integration-test.py" 2>&1 | tee "${FULL_INTEGRATION_TEST_PATH}/integration_test.log"
status=$?
sync || true
echo "integration test finished with status ${status}"
[ $status -eq 0 ]  || exit 1
