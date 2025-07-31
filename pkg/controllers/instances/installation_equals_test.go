// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances_test

import (
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gardener/landscaper-service/pkg/controllers/instances"
	"github.com/gardener/landscaper-service/pkg/utils"
)

func createDummyInstallationSpec() *lsv1alpha1.InstallationSpec {
	return &lsv1alpha1.InstallationSpec{
		ComponentDescriptor: &lsv1alpha1.ComponentDescriptorDefinition{
			Reference: &lsv1alpha1.ComponentDescriptorReference{
				ComponentName: "compa",
				Version:       "v0.0.1",
			},
		},
		Context: "test",
		Blueprint: lsv1alpha1.BlueprintDefinition{
			Reference: &lsv1alpha1.RemoteBlueprintReference{
				ResourceName: "blueprint",
			},
		},
		Imports: lsv1alpha1.InstallationImports{
			Data: []lsv1alpha1.DataImport{
				{
					Name:    "testimport",
					DataRef: "testimport",
				},
			},
			Targets: []lsv1alpha1.TargetImport{
				{
					Name:   "testtarget",
					Target: "testtarget",
				},
			},
		},
		Exports: lsv1alpha1.InstallationExports{
			Data: []lsv1alpha1.DataExport{
				{
					Name:    "dataexport",
					DataRef: "dataexport",
				},
			},
			Targets: []lsv1alpha1.TargetExport{
				{
					Name:   "targetexport",
					Target: "targetexport",
				},
			},
		},
	}
}

var _ = Describe("Installation Equals", func() {
	It("should correctly test equality of empty specs", func() {
		specA := &lsv1alpha1.InstallationSpec{}
		specB := &lsv1alpha1.InstallationSpec{}

		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeTrue())
	})

	It("should correctly test equality of identical specs", func() {
		specA := createDummyInstallationSpec()
		specB := createDummyInstallationSpec()

		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeTrue())
	})

	It("should correctly test inequality of different specs", func() {
		specA := createDummyInstallationSpec()
		specB := createDummyInstallationSpec()

		specB.Imports.Data = append(specB.Imports.Data, lsv1alpha1.DataImport{Name: "foo"})

		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeFalse())

		specA = createDummyInstallationSpec()
		specB = createDummyInstallationSpec()

		specA.Exports.Targets = []lsv1alpha1.TargetExport{}

		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeFalse())

		specA = createDummyInstallationSpec()
		specB = createDummyInstallationSpec()

		specB.ComponentDescriptor.Reference.ComponentName = "invalid"
		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeFalse())
	})

	It("should correctly test equality of equal specs with import data mappings", func() {
		specA := createDummyInstallationSpec()
		specB := createDummyInstallationSpec()

		importDataMappings := map[string]lsv1alpha1.AnyJSON{
			"a": utils.StringToAnyJSON("a-string"),
			"b": utils.BoolToAnyJSON(true),
			"c": utils.IntToAnyJSON(42),
		}

		specA.ImportDataMappings = importDataMappings
		specB.ImportDataMappings = importDataMappings

		specA.ImportDataMappings["d"] = lsv1alpha1.NewAnyJSON([]byte("{ \"keyA\": \"val\", \"keyB\": \"val\" }"))
		specB.ImportDataMappings["d"] = lsv1alpha1.NewAnyJSON([]byte("{ \"keyB\": \"val\", \"keyA\": \"val\" }"))

		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeTrue())
	})

	It("should correctly test inequality of different specs with import data mappings", func() {
		specA := createDummyInstallationSpec()
		specB := createDummyInstallationSpec()

		specA.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"a": lsv1alpha1.NewAnyJSON([]byte("{ \"keyA\": \"val\", \"keyB\": \"val\" }")),
		}
		specB.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"a": lsv1alpha1.NewAnyJSON([]byte("{ \"keyB\": \"val\", \"keyA\": \"val2\" }")),
		}
		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeFalse())

		specA.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"a": lsv1alpha1.NewAnyJSON([]byte("{ \"keyA\": \"val\", \"keyB\": \"val\" }")),
		}
		specB.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"a": utils.StringToAnyJSON("test"),
		}
		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeFalse())

		specA.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"a": lsv1alpha1.NewAnyJSON([]byte("{ \"keyA\": \"val\", \"keyB\": \"val\" }")),
		}
		specB.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"b": utils.StringToAnyJSON("test"),
		}
		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeFalse())

		specA.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"a": lsv1alpha1.NewAnyJSON([]byte("{ \"keyA\": \"val\", \"keyB\": \"val\" }")),
		}
		specB.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"a": lsv1alpha1.NewAnyJSON([]byte("{ \"keyA\": \"val\", \"keyB\": \"val\" }")),
			"b": utils.StringToAnyJSON("test"),
		}
		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeFalse())

		specA.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"a": lsv1alpha1.NewAnyJSON([]byte("{ \"keyA\": \"val\", \"keyB\": \"val\" }")),
			"b": utils.StringToAnyJSON("test"),
		}
		specB.ImportDataMappings = map[string]lsv1alpha1.AnyJSON{
			"a": lsv1alpha1.NewAnyJSON([]byte("{ \"keyA\": \"val\", \"keyB\": \"val\" }")),
		}
		Expect(instances.InstallationSpecDeepEquals(specA, specB)).To(BeFalse())
	})
})
