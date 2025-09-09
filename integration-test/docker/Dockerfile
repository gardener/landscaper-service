# SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.25.1-alpine3.22@sha256:b6ed3fd0452c0e9bcdef5597f29cc1418f61672e9d3a2f55bf02e7222c014abd

RUN apk add --no-cache --no-progress \
    bash \
    git \
    curl \
    openssl openssl-dev \
    expect \
    gcc \
    build-base \
    cargo \
    python3 python3-dev py3-pip  \
    libc-dev \
    libffi-dev

RUN curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash && \
    ln -s "$(which helm)" /bin/helm3

RUN curl -LO https://dl.k8s.io/release/v1.25.14/bin/linux/amd64/kubectl && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

RUN pip3 install --upgrade pip --break-system-packages && \
    pip3 install gardener-cicd-libs --break-system-packages