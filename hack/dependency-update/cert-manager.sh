#!/bin/bash

SCRIPT_PATH="$(dirname $0)"
SOURCE_PATH="${SCRIPT_PATH}/../.."
TMP_PATH="$(mktemp -d)"
CERT_MANAGER_VERSION=$1

export HELM_EXPERIMENTAL_OCI=1

helm repo add jetstack https://charts.jetstack.io
helm repo update
helm pull jetstack/cert-manager --version ${CERT_MANAGER_VERSION} --untar --destination ${TMP_PATH}
helm package ${TMP_PATH}/cert-manager -d ${TMP_PATH}
helm push ${TMP_PATH}/cert-manager-${CERT_MANAGER_VERSION}.tgz oci://eu.gcr.io/gardener-project/landscaper-service/charts

docker pull --platform amd64 quay.io/jetstack/cert-manager-controller:${CERT_MANAGER_VERSION}
docker tag quay.io/jetstack/cert-manager-controller:${CERT_MANAGER_VERSION} eu.gcr.io/gardener-project/landscaper-service/cert-manager-controller:${CERT_MANAGER_VERSION}
docker push eu.gcr.io/gardener-project/landscaper-service/cert-manager-controller:${CERT_MANAGER_VERSION}

docker pull --platform amd64 quay.io/jetstack/cert-manager-webhook:${CERT_MANAGER_VERSION}
docker tag quay.io/jetstack/cert-manager-webhook:${CERT_MANAGER_VERSION} eu.gcr.io/gardener-project/landscaper-service/cert-manager-webhook:${CERT_MANAGER_VERSION}
docker push eu.gcr.io/gardener-project/landscaper-service/cert-manager-webhook:${CERT_MANAGER_VERSION}

docker pull --platform amd64 quay.io/jetstack/cert-manager-cainjector:${CERT_MANAGER_VERSION}
docker tag quay.io/jetstack/cert-manager-cainjector:${CERT_MANAGER_VERSION} eu.gcr.io/gardener-project/landscaper-service/cert-manager-cainjector:${CERT_MANAGER_VERSION}
docker push eu.gcr.io/gardener-project/landscaper-service/cert-manager-cainjector:${CERT_MANAGER_VERSION}

docker pull --platform amd64 quay.io/jetstack/cert-manager-ctl:${CERT_MANAGER_VERSION}
docker tag quay.io/jetstack/cert-manager-ctl:${CERT_MANAGER_VERSION} eu.gcr.io/gardener-project/landscaper-service/cert-manager-ctl:${CERT_MANAGER_VERSION}
docker push eu.gcr.io/gardener-project/landscaper-service/cert-manager-ctl:${CERT_MANAGER_VERSION}

docker pull --platform amd64 quay.io/jetstack/cert-manager-acmesolver:${CERT_MANAGER_VERSION}
docker tag quay.io/jetstack/cert-manager-acmesolver:${CERT_MANAGER_VERSION} eu.gcr.io/gardener-project/landscaper-service/cert-manager-acmesolver:${CERT_MANAGER_VERSION}
docker push eu.gcr.io/gardener-project/landscaper-service/cert-manager-acmesolver:${CERT_MANAGER_VERSION}


export CERT_MANAGER_VERSION
export CERT_MANAGER_CHART_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/charts/cert-manager ${CERT_MANAGER_VERSION})
export CERT_MANAGER_CAINJECTOR_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/cert-manager-cainjector ${CERT_MANAGER_VERSION})
export CERT_MANAGER_CONTROLLER_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/cert-manager-controller ${CERT_MANAGER_VERSION})
export CERT_MANAGER_CTL_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/cert-manager-ctl ${CERT_MANAGER_VERSION})
export CERT_MANAGER_WEBHOOK_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/cert-manager-webhook ${CERT_MANAGER_VERSION})
export CERT_MANAGER_ACMESOLVER_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/cert-manager-acmesolver ${CERT_MANAGER_VERSION})

cat ${SCRIPT_PATH}/template/resources-cert-manager.yaml | envsubst > ${SOURCE_PATH}/.landscaper/logging-stack/resources-cert-manager.yaml
