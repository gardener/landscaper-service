// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package webhook_test

import (
	"context"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/webhook"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

func createLandscaperDeployment(name, namespace string) *lssv1alpha1.LandscaperDeployment {
	deployment := &lssv1alpha1.LandscaperDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       lssv1alpha1.LandscaperDeploymentDefinition.Names.Kind,
			APIVersion: lssv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return deployment
}

var _ = Describe("LandscaperDeployment", func() {
	var (
		validator webhook.GenericValidator
		ctx       context.Context
	)

	BeforeEach(func() {
		var err error
		validator, err = webhook.ValidatorFromResourceType(logr.Discard(), testenv.Client, envtest.LandscaperServiceScheme, webhook.LandscaperDeploymentsResourceType)
		Expect(err).ToNot(HaveOccurred())

		ctx = context.Background()
	})

	It("should allow valid resource", func() {
		testObj := createLandscaperDeployment("test", "lss-system")
		testObj.Spec = lssv1alpha1.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "test",
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
			Region: "eu",
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())
	})

	It("should deny resource with invalid tenant id", func() {
		testObj := createLandscaperDeployment("test", "lss-system")
		testObj.Spec = lssv1alpha1.LandscaperDeploymentSpec{
			TenantId: "test00001",
			Purpose:  "test",
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
			Region: "eu",
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())

		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Reason).ToNot(BeNil())

		Expect(string(response.Result.Reason)).To(ContainSubstring("spec.tenantId"))
	})

	It("should deny resource with invalid purpose", func() {
		testObj := createLandscaperDeployment("test", "lss-system")
		testObj.Spec = lssv1alpha1.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "",
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
			Region: "eu",
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())

		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Reason).ToNot(BeNil())

		Expect(string(response.Result.Reason)).To(ContainSubstring("spec.purpose"))
	})
})
