// SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/apis/validation"
)

func createLandscaperDeployment() *v1alpha1.LandscaperDeployment {
	return &v1alpha1.LandscaperDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-landscaper-deployment",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.LandscaperDeploymentSpec{
			TenantId: "test-id1",
			Purpose:  "test-purpose",
		},
	}
}

var _ = Describe("Validation of LandscaperDeployments", func() {
	It("should accept supported deployers", func() {
		ld := createLandscaperDeployment()
		ld.Spec.LandscaperConfiguration.Deployers = []string{"manifest", "helm", "container"}
		errList := validation.ValidateLandscaperDeployment(ld, nil)
		Expect(errList).To(BeEmpty())
	})

	It("should accept default deployers", func() {
		ld := createLandscaperDeployment()
		ld.Spec.LandscaperConfiguration.Deployers = nil
		errList := validation.ValidateLandscaperDeployment(ld, nil)
		Expect(errList).To(BeEmpty())
	})

	It("should reject unsupported deployers", func() {
		ld := createLandscaperDeployment()
		ld.Spec.LandscaperConfiguration.Deployers = []string{"fantasy-deployer"}
		errList := validation.ValidateLandscaperDeployment(ld, nil)
		Expect(errList).To(HaveLen(1))
		Expect(errList[0].Type).To(Equal(field.ErrorTypeNotSupported))
		Expect(errList[0].Field).To(Equal("spec.landscaperConfiguration.deployers"))
	})
})
