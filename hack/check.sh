#!/bin/bash

# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

PROJECT_ROOT="$(realpath $(dirname $0)/..)"

if [[ -z ${LOCALBIN:-} ]]; then
  LOCALBIN="$PROJECT_ROOT/bin"
fi
if [[ -z ${LINTER:-} ]]; then
  LINTER="$LOCALBIN/golangci-lint"
fi

GOLANGCI_LINT_CONFIG_FILE=""

laas_module_path=()
int_test_module_path=()
for arg in "$@"; do
  case $arg in
    --golangci-lint-config=*)
      GOLANGCI_LINT_CONFIG_FILE="-c ${arg#*=}"
      shift
      ;;
    $PROJECT_ROOT/integration-test/*)
      int_test_module_path+=("./$(realpath "--relative-base=$PROJECT_ROOT/integration-test" "$arg")")
      ;;
    *)
      laas_module_path+=("./$(realpath "--relative-base=$PROJECT_ROOT" "$arg")")
      ;;
  esac
done

echo "> Check"

echo "integration-test module: ${int_test_module_path[@]}"
(
  cd "$PROJECT_ROOT/integration-test"
  echo "  Executing golangci-lint"
  "$LINTER" run $GOLANGCI_LINT_CONFIG_FILE --timeout 10m "${int_test_module_path[@]}"
  echo "  Executing go vet"
  go vet "${int_test_module_path[@]}"
)

echo "root module: ${laas_module_path[@]}"
echo "  Executing golangci-lint"
"$LINTER" run $GOLANGCI_LINT_CONFIG_FILE --timeout 10m "${laas_module_path[@]}"
echo "  Executing go vet"
go vet "${laas_module_path[@]}"

if [[ ${SKIP_FORMATTING_CHECK:-"false"} == "false" ]]; then
  echo "Checking formatting"
  "$PROJECT_ROOT/hack/format.sh" --verify "$@"
fi

echo "All checks successful"