#!/bin/bash
#
# Copyright (c) 2024 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# SPDX-License-Identifier: Apache-2.0

set -e

if [ -z $1 ]; then
  echo "provider argument is required"
  exit 1
fi

SOURCE_PATH="$(realpath $(dirname $0)/..)"
EFFECTIVE_VERSION="$(${SOURCE_PATH}/hack/get-version.sh)"

echo -n "> Updating helm chart version"
${SOURCE_PATH}/hack/update-helm-chart-version.sh ${EFFECTIVE_VERSION}

echo "> Create Component Version ${EFFECTIVE_VERSION}"

PROVIDER=$1
TMP_PATH="$(mktemp -d)"
COMPONENT_ARCHIVE_PATH="${TMP_PATH}/ctf"
COMMIT_SHA=$(git rev-parse HEAD)

LANDSCAPER_SERVICE_CHART_PATH="${SOURCE_PATH}/charts/landscaper-service"
SHOOT_SIDECAR_CHART_PATH="${SOURCE_PATH}/charts/landscaper-service-target-shoot-sidecar-server"
SHOOT_SIDECAR_RBAC_CHART_PATH="${SOURCE_PATH}/charts/sidecar-rbac"

ocm add componentversions --create --file ${COMPONENT_ARCHIVE_PATH} ${SOURCE_PATH}/.landscaper/components.yaml \
  --settings ${SOURCE_PATH}/.landscaper/ocm-settings.yaml \
  -- VERSION=${EFFECTIVE_VERSION} \
     COMMIT_SHA=${COMMIT_SHA} \
     PROVIDER=${PROVIDER} \
     LANDSCAPER_SERVICE_CHART_PATH=${LANDSCAPER_SERVICE_CHART_PATH} \
     SHOOT_SIDECAR_CHART_PATH=${SHOOT_SIDECAR_CHART_PATH} \
     SHOOT_SIDECAR_RBAC_CHART_PATH=${SHOOT_SIDECAR_RBAC_CHART_PATH}

echo "> Transfer Component version ${EFFECTIVE_VERSION} to ${PROVIDER}"
ocm ${OCM_CONFIG} transfer ctf --copy-resources --recursive --overwrite --lookup ${PROVIDER} ${COMPONENT_ARCHIVE_PATH} ${PROVIDER}

echo "> Remote Component Version Landscaper Service"
ocm get componentversion --repo OCIRegistry::${PROVIDER} "github.com/gardener/landscaper-service:${EFFECTIVE_VERSION}" -o yaml

echo "> Remote Component Version Landscaper Instance"
ocm get componentversion --repo OCIRegistry::${PROVIDER} "github.com/gardener/landscaper-service/landscaper-instance:${EFFECTIVE_VERSION}" -o yaml