// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package utils_test

import (
	"encoding/json"

	"github.com/gardener/landscaper-service/pkg/apis/constants"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	provisioningv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/provisioning/v1alpha2"
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
		refList := []provisioningv1alpha2.ObjectReference{
			{
				Name:      "one",
				Namespace: "ns",
			},
			{
				Name:      "two",
				Namespace: "ns",
			},
		}

		refContained := provisioningv1alpha2.ObjectReference{
			Name:      "two",
			Namespace: "ns",
		}
		Expect(utils.ContainsReference(refList, &refContained)).To(BeTrue())

		refNotContained := provisioningv1alpha2.ObjectReference{
			Name:      "two",
			Namespace: "other",
		}
		Expect(utils.ContainsReference(refList, &refNotContained)).To(BeFalse())
	})

	It("should remove a reference from a reference list", func() {
		refList := []provisioningv1alpha2.ObjectReference{
			{
				Name:      "one",
				Namespace: "ns",
			},
			{
				Name:      "two",
				Namespace: "ns",
			},
		}

		toRemove := provisioningv1alpha2.ObjectReference{
			Name:      "one",
			Namespace: "ns",
		}

		newList := utils.RemoveReference(refList, &toRemove)
		Expect(newList).To(HaveLen(1))
		Expect(newList[0].Name).To(Equal("two"))

		newList = utils.RemoveReference(newList, &toRemove)
		Expect(newList).To(HaveLen(1))

		toRemove = provisioningv1alpha2.ObjectReference{
			Name:      "two",
			Namespace: "ns",
		}

		newList = utils.RemoveReference(newList, &toRemove)
		Expect(newList).To(HaveLen(0))
	})

	It("should detect operation annotations", func() {
		secret := &corev1.Secret{}
		Expect(utils.HasOperationAnnotation(secret, constants.LandscaperServiceOperationIgnore)).To(BeFalse())

		secret.ObjectMeta.Annotations = map[string]string{
			"someKey": "someVar",
		}
		Expect(utils.HasOperationAnnotation(secret, constants.LandscaperServiceOperationIgnore)).To(BeFalse())

		secret.ObjectMeta.Annotations[constants.LandscaperServiceOperationAnnotation] = "invalid"
		Expect(utils.HasOperationAnnotation(secret, constants.LandscaperServiceOperationIgnore)).To(BeFalse())

		secret.ObjectMeta.Annotations[constants.LandscaperServiceOperationAnnotation] = constants.LandscaperServiceOperationIgnore
		Expect(utils.HasOperationAnnotation(secret, constants.LandscaperServiceOperationIgnore)).To(BeTrue())
	})

	It("should set operation annotation", func() {
		secret := &corev1.Secret{}
		utils.SetOperationAnnotation(secret, constants.LandscaperServiceOperationIgnore)
		Expect(secret.ObjectMeta.Annotations).To(HaveKeyWithValue(constants.LandscaperServiceOperationAnnotation, constants.LandscaperServiceOperationIgnore))

		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"someKey": "someVar",
				},
			},
		}
		utils.SetOperationAnnotation(secret, constants.LandscaperServiceOperationIgnore)
		Expect(secret.ObjectMeta.Annotations).To(HaveKeyWithValue(constants.LandscaperServiceOperationAnnotation, constants.LandscaperServiceOperationIgnore))
		Expect(secret.ObjectMeta.Annotations).To(HaveKeyWithValue("someKey", "someVar"))
	})

	It("should remove operation annotation", func() {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"someKey": "someVar",
					constants.LandscaperServiceOperationAnnotation: constants.LandscaperServiceOperationIgnore,
				},
			},
		}
		utils.RemoveOperationAnnotation(secret)
		Expect(secret.ObjectMeta.Annotations).ToNot(HaveKeyWithValue(constants.LandscaperServiceOperationAnnotation, constants.LandscaperServiceOperationIgnore))
		Expect(secret.ObjectMeta.Annotations).To(HaveKeyWithValue("someKey", "someVar"))
	})
})
