// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"reflect"

	"k8s.io/apimachinery/pkg/util/json"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
)

// InstallationSpecDeepEquals test whether two installation specs are deeply equal.
func InstallationSpecDeepEquals(specA, specB *lsv1alpha1.InstallationSpec) bool {
	testImportDataMappings := func(specA, specB *lsv1alpha1.InstallationSpec) bool {
		for k, v := range specA.ImportDataMappings {
			if _, ok := specB.ImportDataMappings[k]; !ok {
				return false
			}

			var importDataMappingA interface{}
			if err := json.Unmarshal(v.RawMessage, &importDataMappingA); err != nil {
				return false
			}

			var importDataMappingB interface{}
			if err := json.Unmarshal(specB.ImportDataMappings[k].RawMessage, &importDataMappingB); err != nil {
				return false
			}

			if !reflect.DeepEqual(importDataMappingA, importDataMappingB) {
				return false
			}

			delete(specA.ImportDataMappings, k)
			delete(specB.ImportDataMappings, k)
		}

		return true
	}

	if len(specA.ImportDataMappings) != len(specB.ImportDataMappings) {
		return false
	}

	if !testImportDataMappings(specA, specB) {
		return false
	}
	if !testImportDataMappings(specB, specA) {
		return false
	}

	return reflect.DeepEqual(specA, specB)
}
