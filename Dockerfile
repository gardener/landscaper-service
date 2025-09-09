# SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

#### BASE ####
FROM gcr.io/distroless/static-debian12:nonroot@sha256:e8a4044e0b4ae4257efa45fc026c0bc30ad320d43bd4c1a7d5271bd241e386d0 as base

#### Landscaper Service controller ####
FROM base AS landscaper-service-controller

ARG TARGETOS
ARG TARGETARCH
WORKDIR /
COPY bin/landscaper-service-controller-$TARGETOS.$TARGETARCH /landscaper-service-controller
USER 65532:65532

WORKDIR /

ENTRYPOINT ["/landscaper-service-controller"]

#### Landscaper Service webhooks server ####
FROM base AS landscaper-service-webhooks-server

ARG TARGETOS
ARG TARGETARCH
WORKDIR /
COPY bin/landscaper-service-webhooks-server-$TARGETOS.$TARGETARCH /landscaper-service-webhooks-server
USER 65532:65532

WORKDIR /

ENTRYPOINT ["/landscaper-service-webhooks-server"]

#### Landscaper Target-shoot Sidecar server ####
FROM base AS landscaper-service-target-shoot-sidecar-server

ARG TARGETOS
ARG TARGETARCH
WORKDIR /
COPY bin/landscaper-service-target-shoot-sidecar-server-$TARGETOS.$TARGETARCH /landscaper-service-target-shoot-sidecar-server
USER 65532:65532

WORKDIR /

ENTRYPOINT ["/landscaper-service-target-shoot-sidecar-server"]
