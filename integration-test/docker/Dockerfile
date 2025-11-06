# SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.25.4-alpine3.22@sha256:d3f0cf7723f3429e3f9ed846243970b20a2de7bae6a5b66fc5914e228d831bbb

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