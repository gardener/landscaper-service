#!/bin/sh

PROJECT_ROOT="$(dirname $0)/.."

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
    apk add --no-cache --no-progress bash
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
