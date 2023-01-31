#!/bin/bash

# SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

SCRIPT_PATH="$(dirname $0)"
SOURCE_PATH="${SCRIPT_PATH}/../.."
TMP_PATH="$(mktemp -d)"
INGRESS_NGINX_CHART_VERSION=$1
INGRESS_NGINX_IMAGE_VERSION=$2

export HELM_EXPERIMENTAL_OCI=1

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm pull ingress-nginx/ingress-nginx --version ${INGRESS_NGINX_CHART_VERSION} --untar --destination ${TMP_PATH}
helm package ${TMP_PATH}/ingress-nginx -d ${TMP_PATH}
helm push ${TMP_PATH}/ingress-nginx-${INGRESS_NGINX_CHART_VERSION}.tgz oci://eu.gcr.io/gardener-project/landscaper-service/charts

docker pull --platform amd64 registry.k8s.io/ingress-nginx:/controller${INGRESS_NGINX_IMAGE_VERSION}
docker tag registry.k8s.io/ingress-nginx/controller:${INGRESS_NGINX_IMAGE_VERSION} eu.gcr.io/gardener-project/landscaper-service/ingress-nginx/controller:${INGRESS_NGINX_IMAGE_VERSION}
docker push eu.gcr.io/gardener-project/landscaper-service/ingress-nginx/controller:${INGRESS_NGINX_IMAGE_VERSION}

export INGRESS_NGINX_CHART_VERSION
export INGRESS_NGINX_IMAGE_VERSION
export INGRESS_NGINX_CHART_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/charts/ingress-nginx ${INGRESS_NGINX_CHART_VERSION})
export INGRESS_NGINX_IMAGE_DIGEST=$(go run ${SCRIPT_PATH}/get-digest.go eu.gcr.io gardener-project/landscaper-service/ingress-nginx/controller ${INGRESS_NGINX_IMAGE_VERSION})

cat ${SCRIPT_PATH}/template/resource-ingress-controller.yaml | envsubst > ${SOURCE_PATH}/.landscaper/ingress-controller/resources.yaml
