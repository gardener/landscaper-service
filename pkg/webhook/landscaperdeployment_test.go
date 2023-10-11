// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package webhook_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	provisioningv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/provisioning/v1alpha2"
	"github.com/gardener/landscaper-service/pkg/webhook"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

func createLandscaperDeployment(name, namespace string) *provisioningv1alpha2.LandscaperDeployment {
	deployment := &provisioningv1alpha2.LandscaperDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       provisioningv1alpha2.LandscaperDeploymentDefinition.Names.Kind,
			APIVersion: provisioningv1alpha2.SchemeGroupVersion.String(),
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
		validator, err = webhook.ValidatorFromResourceType(logging.Discard(), testenv.Client, envtest.LandscaperServiceScheme, webhook.LandscaperDeploymentsResourceType)
		Expect(err).ToNot(HaveOccurred())

		ctx = context.Background()
	})

	It("should allow valid resource", func() {
		testObj := createLandscaperDeployment("test", "lss-system")
		testObj.Spec = provisioningv1alpha2.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "test",
			LandscaperConfiguration: provisioningv1alpha2.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())
	})

	It("should deny resource with invalid tenant id", func() {
		testObj := createLandscaperDeployment("test", "lss-system")
		testObj.Spec = provisioningv1alpha2.LandscaperDeploymentSpec{
			TenantId: "test00001",
			Purpose:  "test",
			LandscaperConfiguration: provisioningv1alpha2.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())

		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Message).ToNot(BeNil())

		Expect(response.Result.Message).To(ContainSubstring("spec.tenantId"))
	})

	It("should deny resource with invalid purpose", func() {
		testObj := createLandscaperDeployment("test", "lss-system")
		testObj.Spec = provisioningv1alpha2.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "",
			LandscaperConfiguration: provisioningv1alpha2.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())

		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Message).ToNot(BeNil())

		Expect(response.Result.Message).To(ContainSubstring("spec.purpose"))
	})

	It("should allow a valid update", func() {
		testObj := createLandscaperDeployment("test", "lss-system")
		testObj.Spec = provisioningv1alpha2.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "test",
			LandscaperConfiguration: provisioningv1alpha2.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		oldObj := testObj.DeepCopyObject()
		testObj.Spec.LandscaperConfiguration.Deployers = []string{
			"manifest",
		}

		request = CreateAdmissionRequestUpdate(testObj, oldObj)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())
	})

	It("should deny an update of the tenant id", func() {
		testObj := createLandscaperDeployment("test", "lss-system")
		testObj.Spec = provisioningv1alpha2.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "test",
			LandscaperConfiguration: provisioningv1alpha2.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		oldObj := testObj.DeepCopyObject()
		testObj.Spec.TenantId = "test0002"

		request = CreateAdmissionRequestUpdate(testObj, oldObj)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
	})
})
