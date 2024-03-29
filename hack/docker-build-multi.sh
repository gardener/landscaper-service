#!/bin/bash
#
# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

PROJECT_ROOT="$(realpath $(dirname $0)/..)"
if [[ -z ${EFFECTIVE_VERSION:-} ]]; then
  EFFECTIVE_VERSION=$("$PROJECT_ROOT/hack/get-version.sh")
fi

DOCKER_BUILDER_NAME="laas-multiarch-builder"
if ! docker buildx ls | grep "$DOCKER_BUILDER_NAME" >/dev/null; then
	docker buildx create --name "$DOCKER_BUILDER_NAME"
fi

for pf in ${PLATFORMS//,/ }; do
  echo "> Building docker images for $pf in version $EFFECTIVE_VERSION ..."
	os=${pf%/*}
	arch=${pf#*/}
	docker buildx build --builder ${DOCKER_BUILDER_NAME} --load --build-arg EFFECTIVE_VERSION=${EFFECTIVE_VERSION} --platform ${pf} -t landscaper-service-controller:${EFFECTIVE_VERSION}-${os}-${arch} -f Dockerfile --target landscaper-service-controller "${PROJECT_ROOT}"
	docker buildx build --builder ${DOCKER_BUILDER_NAME} --load --build-arg EFFECTIVE_VERSION=${EFFECTIVE_VERSION} --platform ${pf} -t landscaper-service-webhooks-server:${EFFECTIVE_VERSION}-${os}-${arch} -f Dockerfile --target landscaper-service-webhooks-server "${PROJECT_ROOT}"
	docker buildx build --builder ${DOCKER_BUILDER_NAME} --load --build-arg EFFECTIVE_VERSION=${EFFECTIVE_VERSION} --platform ${pf} -t landscaper-service-target-shoot-sidecar-server:${EFFECTIVE_VERSION}-${os}-${arch} -f Dockerfile --target landscaper-service-target-shoot-sidecar-server "${PROJECT_ROOT}"
done

docker buildx rm "$DOCKER_BUILDER_NAME"