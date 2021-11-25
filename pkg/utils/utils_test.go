// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package utils_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/utils"
)

var _ = Describe("Utils", func() {
	It("should convert a string to AnyJSON", func() {
		input := "teststring"
		asAnyJSON := utils.StringToAnyJSON(input)
		Expect(asAnyJSON.RawMessage).To(Equal(json.RawMessage("\"teststring\"")))
	})

	It("should convert a bool to AnyJSON", func() {
		input := true
		asAnyJSON := utils.BoolToAnyJSON(input)
		Expect(asAnyJSON.RawMessage).To(Equal(json.RawMessage("true")))
	})

	It("should find a reference in a reference list", func() {
		refList := []lssv1alpha1.ObjectReference{
			{
				Name:      "one",
				Namespace: "ns",
			},
			{
				Name:      "two",
				Namespace: "ns",
			},
		}

		refContained := lssv1alpha1.ObjectReference{
			Name:      "two",
			Namespace: "ns",
		}
		Expect(utils.ContainsReference(refList, &refContained)).To(BeTrue())

		refNotContained := lssv1alpha1.ObjectReference{
			Name:      "two",
			Namespace: "other",
		}
		Expect(utils.ContainsReference(refList, &refNotContained)).To(BeFalse())
	})

	It("should remove a reference from a reference list", func() {
		refList := []lssv1alpha1.ObjectReference{
			{
				Name:      "one",
				Namespace: "ns",
			},
			{
				Name:      "two",
				Namespace: "ns",
			},
		}

		toRemove := lssv1alpha1.ObjectReference{
			Name:      "one",
			Namespace: "ns",
		}

		newList := utils.RemoveReference(refList, &toRemove)
		Expect(newList).To(HaveLen(1))
		Expect(newList[0].Name).To(Equal("two"))

		newList = utils.RemoveReference(newList, &toRemove)
		Expect(newList).To(HaveLen(1))

		toRemove = lssv1alpha1.ObjectReference{
			Name:      "two",
			Namespace: "ns",
		}

		newList = utils.RemoveReference(newList, &toRemove)
		Expect(newList).To(HaveLen(0))
	})
})
