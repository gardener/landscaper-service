#!/bin/bash

SCRIPT_PATH="$(dirname $0)"
SOURCE_PATH="${SCRIPT_PATH}/../.."
TMP_PATH="$(mktemp -d)"
SAP_BTP_OPERATOR_VERSION=$1
KUBE_RBAC_PROXY_VERSION=$2

export HELM_EXPERIMENTAL_OCI=1

helm repo add sap-btp-operator https://sap.github.io/sap-btp-service-operator
helm repo update
helm pull sap-btp-operator/sap-btp-operator --version ${SAP_BTP_OPERATOR_VERSION} --untar --destination ${TMP_PATH}
helm package ${TMP_PATH}/sap-btp-operator -d ${TMP_PATH}
helm push ${TMP_PATH}/sap-btp-operator-${SAP_BTP_OPERATOR_VERSION}.tgz oci://eu.gcr.io/gardener-project/landscaper-service/charts

docker pull --platform amd64 ghcr.io/sap/sap-btp-service-operator/controller:${SAP_BTP_OPERATOR_VERSION}
docker tag ghcr.io/sap/sap-btp-service-operator/controller:${SAP_BTP_OPERATOR_VERSION} eu.gcr.io/gardener-project/landscaper-service/sap-btp-service-operator-controller:${SAP_BTP_OPERATOR_VERSION}
docker push eu.gcr.io/gardener-project/landscaper-service/sap-btp-service-operator-controller:${SAP_BTP_OPERATOR_VERSION}

docker pull --platform amd64 quay.io/brancz/kube-rbac-proxy:${KUBE_RBAC_PROXY_VERSION}
docker tag quay.io/brancz/kube-rbac-proxy:${KUBE_RBAC_PROXY_VERSION} eu.gcr.io/gardener-project/landscaper-service/kube-rbac-proxy:${KUBE_RBAC_PROXY_VERSION}
docker push eu.gcr.io/gardener-project/landscaper-service/kube-rbac-proxy:${KUBE_RBAC_PROXY_VERSION}

export SAP_BTP_OPERATOR_VERSION
export KUBE_RBAC_PROXY_VERSION
export SAP_BTP_OPERATOR_CHART_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/charts/sap-btp-operator ${SAP_BTP_OPERATOR_VERSION})
export SAP_BTP_OPERATOR_CONTROLLER_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/sap-btp-service-operator-controller ${SAP_BTP_OPERATOR_VERSION})
export KUBE_RBAC_PROXY_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/kube-rbac-proxy ${KUBE_RBAC_PROXY_VERSION})

cat ${SCRIPT_PATH}/template/resources-sap-btp-service-operator.yaml | envsubst > ${SOURCE_PATH}/.landscaper/logging-stack/resources-sap-btp-service-operator.yaml
