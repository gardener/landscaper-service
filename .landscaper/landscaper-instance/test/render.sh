#!/usr/bin/env bash

TEST_DIR="$(dirname $0)"
BASE_DIR="$(dirname $0)/.."

RENDER_TMP_DIR="$(mktemp -d)"
RESOURCES_FILE="${RENDER_TMP_DIR}/resources.yaml"

cp -R "${BASE_DIR}/." "${RENDER_TMP_DIR}"
export VERSION=v0.1.0
envsubst <"${BASE_DIR}/resources.yaml" >"${RESOURCES_FILE}"

landscaper-cli blueprints render ${BASE_DIR}/blueprint/installation \
    -c "${TEST_DIR}/component-descriptor.yaml" \
    -f "${TEST_DIR}/values.yaml" \
    -e "${TEST_DIR}/export-templates.yaml" \
    -r "${RESOURCES_FILE}"