#!/bin/bash
#
# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

PROJECT_ROOT="$(realpath $(dirname $0)/..)"
if [[ -z ${LOCALBIN:-} ]]; then
  LOCALBIN="$PROJECT_ROOT/bin"
fi
if [[ -z ${CODE_GEN_SCRIPT:-} ]]; then
  CODE_GEN_SCRIPT="$LOCALBIN/kube_codegen.sh"
fi
if [[ -z ${CONTROLLER_GEN:-} ]]; then
  CONTROLLER_GEN="$LOCALBIN/controller-gen"
fi
LAAS_MODULE_PATH="github.com/gardener/landscaper-service"

# Code generation expects this repo to lie under <whatever>/github.com/gardener/landscaper-service, so let's verify that this is the case.
src_path="$(realpath "$PROJECT_ROOT")"
for parent in $(tr '/' '\n' <<< $LAAS_MODULE_PATH | tac); do
  if [[ "$src_path" != */$parent ]]; then
    echo "error: code generation expects the landscaper-service repository to be contained into a folder structure matching its module path '$LAAS_MODULE_PATH'"
    echo "expected path element: $parent"
    echo "actual path element: ${src_path##*/}"
    exit 1
  fi
  src_path="${src_path%/$parent}"
done

rm -f ${GOPATH}/bin/deepcopy-gen
rm -f ${GOPATH}/bin/defaulter-gen
rm -f ${GOPATH}/bin/conversion-gen

source "$CODE_GEN_SCRIPT"

echo "> Generating deepcopy/conversion/defaulter functions"
kube::codegen::gen_helpers \
  --input-pkg-root "$LAAS_MODULE_PATH" \
  --output-base "$src_path" \
  --boilerplate "${PROJECT_ROOT}/hack/boilerplate.go.txt"

echo
echo "> Generating CRDs"
"$CONTROLLER_GEN" crd paths="$PROJECT_ROOT/pkg/apis/..." output:crd:artifacts:config="$PROJECT_ROOT/pkg/crdmanager/crdresources"
