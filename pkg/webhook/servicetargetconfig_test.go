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

	lsscore "github.com/gardener/landscaper-service/pkg/apis/core"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/webhook"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

func createServiceTargetConfig(name, namespace string) *lssv1alpha1.ServiceTargetConfig {
	config := &lssv1alpha1.ServiceTargetConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       lssv1alpha1.ServiceTargetConfigDefinition.Names.Kind,
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
		validator, err = webhook.ValidatorFromResourceType(logr.Discard(), testenv.Client, envtest.LandscaperServiceScheme, webhook.ServiceTargetConfigsResourceType)
		Expect(err).ToNot(HaveOccurred())

		ctx = context.Background()
	})

	It("should allow valid resource", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.ObjectMeta.Labels = map[string]string{
			lsscore.ServiceTargetConfigRegionLabelName:  "eu",
			lsscore.ServiceTargetConfigVisibleLabelName: "true",
		}
		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			ProviderType: "aws",
			Priority:     10,
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
		Expect(response.Allowed).To(BeTrue())
	})

	It("should deny resource with missing labels", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			ProviderType: "aws",
			Priority:     10,
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
		Expect(response.Result.Reason).ToNot(BeNil())

		Expect(string(response.Result.Reason)).To(ContainSubstring("metadata.labels.config.landscaper-service.gardener.cloud/region"))
		Expect(string(response.Result.Reason)).To(ContainSubstring("metadata.labels.config.landscaper-service.gardener.cloud/visible"))
	})

	It("should deny resource with invalid visible label", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.ObjectMeta.Labels = map[string]string{
			lsscore.ServiceTargetConfigRegionLabelName:  "eu",
			lsscore.ServiceTargetConfigVisibleLabelName: "abc",
		}
		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			ProviderType: "aws",
			Priority:     10,
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
		Expect(response.Result.Reason).ToNot(BeNil())

		Expect(string(response.Result.Reason)).To(ContainSubstring("metadata.labels.config.landscaper-service.gardener.cloud/visible"))
	})

	It("should deny resource with unknown provider type", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.ObjectMeta.Labels = map[string]string{
			lsscore.ServiceTargetConfigRegionLabelName:  "eu",
			lsscore.ServiceTargetConfigVisibleLabelName: "false",
		}
		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			ProviderType: "invalid",
			Priority:     10,
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
		Expect(response.Result.Reason).ToNot(BeNil())

		Expect(string(response.Result.Reason)).To(ContainSubstring("spec.providerType"))
	})

	It("should deny resource with invalid secret reference", func() {
		testObj := createServiceTargetConfig("test", "lss-system")

		testObj.ObjectMeta.Labels = map[string]string{
			lsscore.ServiceTargetConfigRegionLabelName:  "eu",
			lsscore.ServiceTargetConfigVisibleLabelName: "true",
		}
		testObj.Spec = lssv1alpha1.ServiceTargetConfigSpec{
			ProviderType: "aws",
			Priority:     10,
			SecretRef: lssv1alpha1.SecretReference{
				ObjectReference: lssv1alpha1.ObjectReference{
					Name:      "",
					Namespace: "",
				},
				Key: "",
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Reason).ToNot(BeNil())

		Expect(string(response.Result.Reason)).To(ContainSubstring("spec.secretRef.key"))
		Expect(string(response.Result.Reason)).To(ContainSubstring("spec.secretRef.name"))
		Expect(string(response.Result.Reason)).ToNot(ContainSubstring("spec.secretRef.namespace"))
	})
})
