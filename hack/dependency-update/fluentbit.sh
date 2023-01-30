#!/bin/bash

SCRIPT_PATH="$(dirname $0)"
SOURCE_PATH="${SCRIPT_PATH}/../.."
TMP_PATH="$(mktemp -d)"
FLUENTBIT_VERSION=$1

docker pull --platform amd64 cr.fluentbit.io/fluent/fluent-bit:${FLUENTBIT_VERSION}
docker tag cr.fluentbit.io/fluent/fluent-bit:${FLUENTBIT_VERSION} eu.gcr.io/gardener-project/landscaper-service/fluent-bit:${FLUENTBIT_VERSION}
docker push eu.gcr.io/gardener-project/landscaper-service/fluent-bit:${FLUENTBIT_VERSION}

export FLUENTBIT_VERSION
export FLUENTBIT_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/fluent-bit ${FLUENTBIT_VERSION})

cat ${SCRIPT_PATH}/template/resources-fluentbit.yaml | envsubst > ${SOURCE_PATH}/.landscaper/logging-stack/resources-fluentbit.yaml
