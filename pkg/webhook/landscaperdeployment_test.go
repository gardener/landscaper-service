// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package webhook_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/webhook"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

func createLandscaperDeployment(name, namespace string) *lssv1alpha1.LandscaperDeployment {
	deployment := &lssv1alpha1.LandscaperDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LandscaperDeployment",
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
		validator, err = webhook.ValidatorFromResourceType(logging.Discard(), testenv.Client, envtest.LandscaperServiceScheme, webhook.LandscaperDeploymentsResourceType)
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
		testObj.Spec = lssv1alpha1.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "",
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
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
		testObj.Spec = lssv1alpha1.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "test",
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
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
		testObj.Spec = lssv1alpha1.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "test",
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
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

	It("should validate high availability config", func() {
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
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		oldObject := testObj.DeepCopyObject()
		testObj.Spec.HighAvailabilityConfig = &lssv1alpha1.HighAvailabilityConfig{
			ControlPlaneFailureTolerance: "node",
		}

		request = CreateAdmissionRequestUpdate(testObj, oldObject)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		oldObject = testObj.DeepCopyObject()
		testObj.Spec.HighAvailabilityConfig.ControlPlaneFailureTolerance = "zone"

		request = CreateAdmissionRequestUpdate(testObj, oldObject)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
	})

	It("shall validate an external data plane", func() {
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
			DataPlane: &lssv1alpha1.DataPlane{
				Kubeconfig: "{}",
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		testObj.Spec.DataPlane = &lssv1alpha1.DataPlane{
			SecretRef: &lssv1alpha1.SecretReference{
				Key: "kubeconfig",
				ObjectReference: lssv1alpha1.ObjectReference{
					Name:      "dataplane",
					Namespace: "default",
				},
			},
		}

		request = CreateAdmissionRequest(testObj)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())
	})

	It("shall deny an invalid external data plane", func() {
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
			DataPlane: &lssv1alpha1.DataPlane{
				Kubeconfig: "{}",
				SecretRef: &lssv1alpha1.SecretReference{
					Key: "kubeconfig",
					ObjectReference: lssv1alpha1.ObjectReference{
						Name:      "dataplane",
						Namespace: "default",
					},
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
	})

	It("shall deny switching to and from external data plane", func() {
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
			DataPlane: &lssv1alpha1.DataPlane{
				Kubeconfig: "{}",
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		oldObject := testObj.DeepCopyObject()
		testObj.Spec.DataPlane = nil

		request = CreateAdmissionRequestUpdate(testObj, oldObject)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())

		oldObject = testObj.DeepCopyObject()
		testObj.Spec.DataPlane = &lssv1alpha1.DataPlane{
			Kubeconfig: "{}",
		}

		request = CreateAdmissionRequestUpdate(testObj, oldObject)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
	})

	It("shall deny an invalid combination of internal and external data plane", func() {
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
			DataPlane: &lssv1alpha1.DataPlane{
				Kubeconfig: "{}",
			},
			OIDCConfig: &lssv1alpha1.OIDCConfig{
				ClientID: "test",
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())

		testObj.Spec = lssv1alpha1.LandscaperDeploymentSpec{
			TenantId: "test0001",
			Purpose:  "test",
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
			DataPlane: &lssv1alpha1.DataPlane{
				Kubeconfig: "{}",
			},
			HighAvailabilityConfig: &lssv1alpha1.HighAvailabilityConfig{
				ControlPlaneFailureTolerance: "zone",
			},
		}

		request = CreateAdmissionRequest(testObj)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
	})
})
