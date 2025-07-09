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

func createServiceTargetConfig(name, namespace string) *lssv1alpha1.ServiceTargetConfig {
	config := &lssv1alpha1.ServiceTargetConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceTargetConfig",
			APIVersion: lssv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return config
}

var _ = Describe("ServiceTargetConfig", func() {
	var (
		validator webhook.GenericValidator
		ctx       context.Context
	)

	BeforeEach(func() {
		var err error
		validator, err = webhook.ValidatorFromResourceType(logging.Discard(), testenv.Client, envtest.LandscaperServiceScheme, webhook.ServiceTargetConfigsResourceType)
		Expect(err).ToNot(HaveOccurred())

		ctx = context.Background()
	})

	It("should allow valid resource", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.Labels = map[string]string{
			lssv1alpha1.ServiceTargetConfigVisibleLabelName: "true",
		}
		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			Priority: 10,
			SecretRef: lssv1alpha1.SecretReference{
				ObjectReference: lssv1alpha1.ObjectReference{
					Name:      "target",
					Namespace: "lss-system",
				},
				Key: "kubeconfig",
			},
			IngressDomain: "ingress.external",
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())
	})

	It("should deny resource with missing labels", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			Priority: 10,
			SecretRef: lssv1alpha1.SecretReference{
				ObjectReference: lssv1alpha1.ObjectReference{
					Name:      "target",
					Namespace: "lss-system",
				},
				Key: "kubeconfig",
			},
			IngressDomain: "ingress.external",
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Message).ToNot(BeNil())

		Expect(response.Result.Message).To(ContainSubstring("metadata.labels.config.landscaper-service.gardener.cloud/visible"))
	})

	It("should deny resource with invalid visible label", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.Labels = map[string]string{
			lssv1alpha1.ServiceTargetConfigVisibleLabelName: "abc",
		}
		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			Priority: 10,
			SecretRef: lssv1alpha1.SecretReference{
				ObjectReference: lssv1alpha1.ObjectReference{
					Name:      "target",
					Namespace: "lss-system",
				},
				Key: "kubeconfig",
			},
			IngressDomain: "ingress.external",
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Message).ToNot(BeNil())

		Expect(response.Result.Message).To(ContainSubstring("metadata.labels.config.landscaper-service.gardener.cloud/visible"))
	})

	It("should deny resource with invalid secret reference", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.Labels = map[string]string{
			lssv1alpha1.ServiceTargetConfigVisibleLabelName: "true",
		}
		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			Priority: 10,
			SecretRef: lssv1alpha1.SecretReference{
				ObjectReference: lssv1alpha1.ObjectReference{
					Name:      "",
					Namespace: "",
				},
				Key: "",
			},
			IngressDomain: "ingress.external",
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Message).ToNot(BeNil())

		Expect(response.Result.Message).To(ContainSubstring("spec.secretRef.key"))
		Expect(response.Result.Message).To(ContainSubstring("spec.secretRef.name"))
		Expect(response.Result.Message).ToNot(ContainSubstring("spec.secretRef.namespace"))
	})

	It("should deny resource with missing ingress domain", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.Labels = map[string]string{
			lssv1alpha1.ServiceTargetConfigVisibleLabelName: "true",
		}
		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			Priority: 10,
			SecretRef: lssv1alpha1.SecretReference{
				ObjectReference: lssv1alpha1.ObjectReference{
					Name:      "target",
					Namespace: "lss-system",
				},
				Key: "kubeconfig",
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Message).ToNot(BeNil())

		Expect(response.Result.Message).To(ContainSubstring("spec.ingressDomain"))
	})
})
