#!/bin/bash

# SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

set -e

if [ -n "$EFFECTIVE_VERSION" ] ; then
  # running in the pipeline use the provided EFFECTIVE_VERSION
  echo "$EFFECTIVE_VERSION"
  exit 0
fi

SOURCE_PATH="$(dirname $0)/.."
VERSION="$(cat "${SOURCE_PATH}/VERSION")"

pushd ${SOURCE_PATH} > /dev/null 2>&1

if [[ "$VERSION" = *-dev ]] ; then
  VERSION="$VERSION-$(git rev-parse HEAD)"
fi

popd > /dev/null 2>&1

echo "$VERSION"
