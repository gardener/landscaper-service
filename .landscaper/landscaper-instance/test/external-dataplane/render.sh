#!/usr/bin/env bash

TEST_DIR="$(dirname $0)"
BASE_DIR="$(dirname $0)/../.."

RENDER_TMP_DIR="$(mktemp -d)"
RESOURCES_FILE="${RENDER_TMP_DIR}/resources.yaml"
COMPONENT_DESCRIPTOR_FILE="${RENDER_TMP_DIR}/test/component-descriptor.yaml"

cp -R "${BASE_DIR}/." "${RENDER_TMP_DIR}"
export VERSION=v0.1.0
envsubst <"${BASE_DIR}/resources.yaml" >"${RESOURCES_FILE}"

LANDSCAPER_COMPONENT_REF="$(yq ./.landscaper/landscaper-instance/component-references.yaml -ojson -I=0)"
yq ".component.componentReferences += [${LANDSCAPER_COMPONENT_REF}]" "${TEST_DIR}/component-descriptor.yaml" > "${COMPONENT_DESCRIPTOR_FILE}"

echo "!!! render global installation blueprint !!!"
landscaper-cli blueprints render ${BASE_DIR}/blueprint/installation-ext-dataplane \
    -c "${COMPONENT_DESCRIPTOR_FILE}" \
    -f "${TEST_DIR}/values-global.yaml" \
    -e "${TEST_DIR}/export-templates.yaml" \
    -r "${RESOURCES_FILE}"

echo "!!! render landscaper rbac blueprint !!!"
landscaper-cli blueprints render ${BASE_DIR}/blueprint/rbac \
    -c "${COMPONENT_DESCRIPTOR_FILE}" \
    -f "${TEST_DIR}/values-rbac.yaml" \
    -r "${RESOURCES_FILE}"

echo "!!! render sidecar rbac blueprint !!!"
landscaper-cli blueprints render ${BASE_DIR}/blueprint/rbac \
    -c "${COMPONENT_DESCRIPTOR_FILE}" \
    -f "${TEST_DIR}/values-sidecar-rbac.yaml" \
    -r "${RESOURCES_FILE}"
