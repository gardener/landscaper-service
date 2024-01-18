// SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package scheduling_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gardener/landscaper-service/pkg/controllers/landscaperdeployments/scheduling"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

var _ = Describe("Evaluation of selectors", func() {

	newLandscaperDeployment := func(tenantID string, labels map[string]string) *lssv1alpha1.LandscaperDeployment {
		return &lssv1alpha1.LandscaperDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: lssv1alpha1.LandscaperDeploymentSpec{
				TenantId: tenantID,
			},
		}
	}

	It("should evaluate a tenant selector", func() {
		const tenantID = "test-tenant"

		selector := &lssv1alpha1.Selector{
			MatchTenant: &lssv1alpha1.TenantSelector{ID: tenantID},
		}

		deployment := newLandscaperDeployment(tenantID, nil)
		match, err := scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeTrue())

		deployment = newLandscaperDeployment("other-tenant", nil)
		match, err = scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeFalse())
	})

	It("should evaluate a label selector", func() {
		const (
			tenantID   = "test-tenant"
			labelName  = "test-label"
			labelValue = "test-value"
		)

		selector := &lssv1alpha1.Selector{
			MatchLabel: &lssv1alpha1.LabelSelector{
				Name:  labelName,
				Value: labelValue,
			},
		}

		deployment := newLandscaperDeployment(tenantID, map[string]string{labelName: labelValue})
		match, err := scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeTrue())

		deployment = newLandscaperDeployment(tenantID, nil)
		match, err = scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeFalse())
	})

	It("should evaluate an or selector", func() {
		const tenantID = "test-tenant"

		selector := &lssv1alpha1.Selector{
			Or: []lssv1alpha1.Selector{
				{MatchTenant: &lssv1alpha1.TenantSelector{ID: "other-tenant"}},
				{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenantID}},
				{MatchTenant: &lssv1alpha1.TenantSelector{ID: "yet-another-tenant"}},
			},
		}

		deployment := newLandscaperDeployment(tenantID, nil)
		match, err := scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeTrue())

		deployment = newLandscaperDeployment("still-another-tenant", nil)
		match, err = scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeFalse())
	})

	It("should evaluate an and selector", func() {
		const (
			tenantID   = "test-tenant"
			labelName  = "test-label"
			labelValue = "test-value"
		)

		selector := &lssv1alpha1.Selector{
			And: []lssv1alpha1.Selector{
				{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenantID}},
				{MatchLabel: &lssv1alpha1.LabelSelector{Name: labelName, Value: labelValue}},
			},
		}

		deployment := newLandscaperDeployment(tenantID, map[string]string{labelName: labelValue})
		match, err := scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeTrue())

		deployment = newLandscaperDeployment(tenantID, map[string]string{labelName: "another-value"})
		match, err = scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeFalse())
	})

	It("should evaluate a not selector", func() {
		const tenantID = "test-tenant"

		selector := &lssv1alpha1.Selector{
			Not: &lssv1alpha1.Selector{
				MatchTenant: &lssv1alpha1.TenantSelector{ID: tenantID},
			},
		}

		deployment := newLandscaperDeployment(tenantID, nil)
		match, err := scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeFalse())

		deployment = newLandscaperDeployment("other-tenant", nil)
		match, err = scheduling.EvaluateSelector(selector, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(match).To(BeTrue())
	})
})
