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

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/webhook"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

func createInstance(name, namespace string) *lssv1alpha1.Instance {
	instance := &lssv1alpha1.Instance{
		TypeMeta: metav1.TypeMeta{
			Kind:       lssv1alpha1.InstanceDefinition.Names.Kind,
			APIVersion: lssv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return instance
}

var _ = Describe("LandscaperDeployment", func() {
	var (
		validator webhook.GenericValidator
		ctx       context.Context
	)

	BeforeEach(func() {
		var err error
		validator, err = webhook.ValidatorFromResourceType(logging.Discard(), testenv.Client, envtest.LandscaperServiceScheme, webhook.InstancesResourceType)
		Expect(err).ToNot(HaveOccurred())

		ctx = context.Background()
	})

	It("should allow valid resource", func() {
		testObj := createInstance("test", "lss-system")

		testObj.Spec = lssv1alpha1.InstanceSpec{
			TenantId: "test0001",
			ID:       "inst0001",
			ServiceTargetConfigRef: lssv1alpha1.ObjectReference{
				Name:      "test",
				Namespace: "lss-system",
			},
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())
	})

	It("should deny resource with invalid tenant id", func() {
		testObj := createInstance("test", "lss-system")

		testObj.Spec = lssv1alpha1.InstanceSpec{
			TenantId: "test00001",
			ID:       "inst0001",
			ServiceTargetConfigRef: lssv1alpha1.ObjectReference{
				Name:      "test",
				Namespace: "lss-system",
			},
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())

		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Reason).ToNot(BeNil())

		Expect(string(response.Result.Reason)).To(ContainSubstring("spec.tenantId"))
	})

	It("should deny resource with invalid instance id", func() {
		testObj := createInstance("test", "lss-system")

		testObj.Spec = lssv1alpha1.InstanceSpec{
			TenantId: "test0001",
			ID:       "inst00001",
			ServiceTargetConfigRef: lssv1alpha1.ObjectReference{
				Name:      "test",
				Namespace: "lss-system",
			},
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())

		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Reason).ToNot(BeNil())

		Expect(string(response.Result.Reason)).To(ContainSubstring("spec.id"))
	})

	It("should deny resource with invalid service target config ref", func() {
		testObj := createInstance("test", "lss-system")

		testObj.Spec = lssv1alpha1.InstanceSpec{
			TenantId: "test0001",
			ID:       "inst0001",
			ServiceTargetConfigRef: lssv1alpha1.ObjectReference{
				Name:      "",
				Namespace: "",
			},
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())

		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Reason).ToNot(BeNil())

		Expect(string(response.Result.Reason)).To(ContainSubstring("spec.serviceTargetConfigRef.name"))
	})

	It("should allow a valid update", func() {
		testObj := createInstance("test", "lss-system")

		testObj.Spec = lssv1alpha1.InstanceSpec{
			TenantId: "test0001",
			ID:       "inst0001",
			ServiceTargetConfigRef: lssv1alpha1.ObjectReference{
				Name:      "test",
				Namespace: "lss-system",
			},
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		oldObject := testObj.DeepCopyObject()
		testObj.Spec.LandscaperConfiguration.Deployers = []string{
			"manifest",
		}

		request = CreateAdmissionRequestUpdate(testObj, oldObject)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())
	})

	It("should deny an update of the tenant id", func() {
		testObj := createInstance("test", "lss-system")

		testObj.Spec = lssv1alpha1.InstanceSpec{
			TenantId: "test0001",
			ID:       "inst0001",
			ServiceTargetConfigRef: lssv1alpha1.ObjectReference{
				Name:      "test",
				Namespace: "lss-system",
			},
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		oldObject := testObj.DeepCopyObject()
		testObj.Spec.TenantId = "test0002"

		request = CreateAdmissionRequestUpdate(testObj, oldObject)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
	})

	It("should deny an update of the instance id", func() {
		testObj := createInstance("test", "lss-system")

		testObj.Spec = lssv1alpha1.InstanceSpec{
			TenantId: "test0001",
			ID:       "inst0001",
			ServiceTargetConfigRef: lssv1alpha1.ObjectReference{
				Name:      "test",
				Namespace: "lss-system",
			},
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		oldObject := testObj.DeepCopyObject()
		testObj.Spec.ID = "inst0002"

		request = CreateAdmissionRequestUpdate(testObj, oldObject)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
	})

	It("should deny an update of the service target config ref", func() {
		testObj := createInstance("test", "lss-system")

		testObj.Spec = lssv1alpha1.InstanceSpec{
			TenantId: "test0001",
			ID:       "inst0001",
			ServiceTargetConfigRef: lssv1alpha1.ObjectReference{
				Name:      "test",
				Namespace: "lss-system",
			},
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())

		oldObject := testObj.DeepCopyObject()
		testObj.Spec.ServiceTargetConfigRef = lssv1alpha1.ObjectReference{
			Name:      "test1",
			Namespace: "lss-system",
		}

		request = CreateAdmissionRequestUpdate(testObj, oldObject)
		response = validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
	})

	It("should validate high availability config", func() {
		testObj := createInstance("test", "lss-system")
		testObj.Spec = lssv1alpha1.InstanceSpec{
			TenantId: "test0001",
			ID:       "inst0001",
			ServiceTargetConfigRef: lssv1alpha1.ObjectReference{
				Name:      "test",
				Namespace: "lss-system",
			},
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
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
})
