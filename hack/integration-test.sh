#!/bin/sh

# SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

set -e

apk add --no-cache --no-progress bash

if ! command -v git &> /dev/null
then
    apk add --no-cache --no-progress git
fi

PROJECT_ROOT="$(dirname $0)/.."
TARGET_CLUSTER="laas-integration-test"
TARGET_CLUSTER_PROVIDER="gcp"
LAAS_REPOSITORY="eu.gcr.io/sap-se-gcr-k8s-private/cnudie/gardener/development"
REPO_AUTH_URL="https://eu.gcr.io"
REPO_CTX_BASE_URL="eu.gcr.io/sap-se-gcr-k8s-private"

unset EFFECTIVE_VERSION
LAAS_VERSION="$(${PROJECT_ROOT}/hack/get-version.sh)"

export PROJECT_ROOT
export TARGET_CLUSTER
export TARGET_CLUSTER_PROVIDER
export LAAS_VERSION
export LAAS_REPOSITORY
export REPO_AUTH_URL
export REPO_CTX_BASE_URL


if ! command -v curl &> /dev/null
then
    apk add --no-cache --no-progress curl openssl
fi

if ! command -v python3 &> /dev/null
then
    echo "Python3 could not be found"
    echo "Try installing it..."
    apk add --no-cache --no-progress python3 python3-dev py3-pip gcc libc-dev libffi-dev openssl-dev cargo build-base
fi

if ! command -v helm &> /dev/null
then
    echo "Helm could not be found"
    echo "Try installing it..."
    export DESIRED_VERSION="v3.7.1"
    curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
    # symlink to /bin/helm3 as it is required by the integration test script
    ln -s "$(which helm)" /bin/helm3
fi

if ! command -v kubectl &> /dev/null
then
    echo "Kubectl could not be found"
    echo "Try installing it..."
    curl -LO https://dl.k8s.io/release/v1.21.0/bin/linux/amd64/kubectl
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
fi

echo "Running pip3 install --upgrade pip"
pip3 install --upgrade pip

echo "Running pip3 install gardener-cicd-libs"
pip3 install gardener-cicd-libs

"${PROJECT_ROOT}/hack/integration-test.py"
