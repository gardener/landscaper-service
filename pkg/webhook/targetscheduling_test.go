// SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package webhook_test

import (
	"context"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/webhook"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

func createTargetScheduling(name, namespace string) *lssv1alpha1.TargetScheduling {
	return &lssv1alpha1.TargetScheduling{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TargetScheduling",
			APIVersion: lssv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

var _ = Describe("TargetScheduling", func() {
	var (
		validator webhook.GenericValidator
		ctx       context.Context
	)

	BeforeEach(func() {
		var err error
		validator, err = webhook.ValidatorFromResourceType(logging.Discard(), testenv.Client, envtest.LandscaperServiceScheme, webhook.TargetSchedulingsResourceType)
		Expect(err).ToNot(HaveOccurred())

		ctx = context.Background()
	})

	expectErrorAtPath := func(scheduling *lssv1alpha1.TargetScheduling, path string) {
		request := CreateAdmissionRequest(scheduling)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeFalse())
		Expect(response.Result).ToNot(BeNil())
		Expect(response.Result.Message).ToNot(BeNil())
		Expect(response.Result.Message).To(ContainSubstring(path))
	}

	It("should allow valid resource", func() {
		testObj := createTargetScheduling("test", "lss-system")
		testObj.Spec.Rules = []lssv1alpha1.SchedulingRule{
			{
				Priority: 10,
				ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
					{Name: "test01", Namespace: "lss-system"},
					{Name: "test02", Namespace: "lss-system"},
				},
				Selector: []lssv1alpha1.Selector{
					{
						MatchTenant: &lssv1alpha1.TenantSelector{ID: "test-tenant-1"},
					},
					{
						MatchLabel: &lssv1alpha1.LabelSelector{Name: "region", Value: "eu"},
					},
					{
						Or: []lssv1alpha1.Selector{
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: "test-tenant-2"}},
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: "test-tenant-3"}},
						},
					},
					{
						And: []lssv1alpha1.Selector{
							{MatchLabel: &lssv1alpha1.LabelSelector{Name: "region", Value: "eu"}},
							{MatchLabel: &lssv1alpha1.LabelSelector{Name: "direction", Value: "north"}},
						},
					},
					{
						Not: &lssv1alpha1.Selector{
							MatchLabel: &lssv1alpha1.LabelSelector{Name: "direction", Value: "north"},
						},
					},
				},
			},
		}

		request := CreateAdmissionRequest(testObj)
		response := validator.Handle(ctx, request)
		Expect(response).ToNot(BeNil())
		Expect(response.Allowed).To(BeTrue())
	})

	It("should deny a rule with negative priority", func() {
		testObj := createTargetScheduling("test", "lss-system")
		testObj.Spec.Rules = []lssv1alpha1.SchedulingRule{
			{
				Priority: -10, // negative!
				ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
					{Name: "test01", Namespace: "lss-system"},
				},
				Selector: []lssv1alpha1.Selector{},
			},
		}
		expectErrorAtPath(testObj, "spec.rules[0].priority")
	})

	It("should deny a rule without servicetargetconfigs", func() {
		testObj := createTargetScheduling("test", "lss-system")
		testObj.Spec.Rules = []lssv1alpha1.SchedulingRule{
			{
				Priority:             10,
				ServiceTargetConfigs: []lssv1alpha1.ObjectReference{}, // empty!
				Selector:             []lssv1alpha1.Selector{},
			},
		}
		expectErrorAtPath(testObj, "spec.rules[0].serviceTargetConfigs")
	})

	It("should deny an selector without term", func() {
		testObj := createTargetScheduling("test", "lss-system")
		testObj.Spec.Rules = []lssv1alpha1.SchedulingRule{
			{
				Priority: 10,
				ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
					{Name: "test01", Namespace: "lss-system"},
				},
				Selector: []lssv1alpha1.Selector{
					{
						MatchTenant: &lssv1alpha1.TenantSelector{ID: "test-tenant-0"},
					},
					{
						// selector without term!
					},
					{
						MatchTenant: &lssv1alpha1.TenantSelector{ID: "test-tenant-2"},
					},
				},
			},
		}
		expectErrorAtPath(testObj, "spec.rules[0].selector[1]")
	})

	It("should deny an selector with more than one term", func() {
		testObj := createTargetScheduling("test", "lss-system")
		testObj.Spec.Rules = []lssv1alpha1.SchedulingRule{
			{
				Priority: 10,
				ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
					{Name: "test01", Namespace: "lss-system"},
				},
				Selector: []lssv1alpha1.Selector{
					{
						MatchTenant: &lssv1alpha1.TenantSelector{ID: "test-tenant-0"},
					},
					{
						MatchTenant: &lssv1alpha1.TenantSelector{ID: "test-tenant-1"},
						MatchLabel:  &lssv1alpha1.LabelSelector{Name: "region", Value: "eu"}, // second term!
					},
					{
						MatchTenant: &lssv1alpha1.TenantSelector{ID: "test-tenant-2"},
					},
				},
			},
		}
		expectErrorAtPath(testObj, "spec.rules[0].selector[1]")
	})

	It("should deny an invalid selector with nested terms", func() {
		testObj := createTargetScheduling("test", "lss-system")
		testObj.Spec.Rules = []lssv1alpha1.SchedulingRule{
			{
				Priority: 10,
				ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
					{Name: "test01", Namespace: "lss-system"},
				},
				Selector: []lssv1alpha1.Selector{
					{
						Or: []lssv1alpha1.Selector{
							{
								And: []lssv1alpha1.Selector{
									{
										Not: &lssv1alpha1.Selector{
											Or: []lssv1alpha1.Selector{
												{
													MatchTenant: &lssv1alpha1.TenantSelector{}, // no tenant id!
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		expectErrorAtPath(testObj, "spec.rules[0].selector[0].or[0].and[0].not.or[0]")
	})

})
