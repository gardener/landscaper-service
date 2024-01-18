// SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package scheduling_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lssscheduling "github.com/gardener/landscaper-service/pkg/controllers/landscaperdeployments/scheduling"
)

var _ = Describe("SortServiceTargetConfigs", func() {

	It("should sort descending by priority", func() {
		configs := []*lssv1alpha1.ServiceTargetConfig{
			{
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 20,
				},
			},
			{
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 10,
				},
			},
			{
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 30,
				},
			},
		}

		lssscheduling.SortServiceTargetConfigs(configs)
		Expect(configs).To(HaveLen(3))
		Expect(configs[0].Spec.Priority).To(Equal(int64(30)))
		Expect(configs[1].Spec.Priority).To(Equal(int64(20)))
		Expect(configs[2].Spec.Priority).To(Equal(int64(10)))
	})

	It("should sort ascending by usage", func() {
		configs := []*lssv1alpha1.ServiceTargetConfig{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "first",
				},
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 10,
				},
				Status: lssv1alpha1.ServiceTargetConfigStatus{
					InstanceRefs: []lssv1alpha1.ObjectReference{
						{
							Name:      "foo",
							Namespace: "bar",
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "second",
				},
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 10,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "third",
				},
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 10,
				},
				Status: lssv1alpha1.ServiceTargetConfigStatus{
					InstanceRefs: []lssv1alpha1.ObjectReference{
						{
							Name:      "foo",
							Namespace: "bar",
						},
						{
							Name:      "foo",
							Namespace: "bar",
						},
					},
				},
			},
		}

		lssscheduling.SortServiceTargetConfigs(configs)
		Expect(configs).To(HaveLen(3))
		Expect(configs[0].GetName()).To(Equal("second"))
		Expect(configs[1].GetName()).To(Equal("first"))
		Expect(configs[2].GetName()).To(Equal("third"))
	})

	It("should sort by priority and usage", func() {
		configs := []*lssv1alpha1.ServiceTargetConfig{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "first",
				},
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 30,
				},
				Status: lssv1alpha1.ServiceTargetConfigStatus{
					InstanceRefs: []lssv1alpha1.ObjectReference{
						{
							Name:      "foo",
							Namespace: "bar",
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "second",
				},
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 20,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "third",
				},
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 40,
				},
				Status: lssv1alpha1.ServiceTargetConfigStatus{
					InstanceRefs: []lssv1alpha1.ObjectReference{
						{
							Name:      "foo",
							Namespace: "bar",
						},
						{
							Name:      "foo",
							Namespace: "bar",
						},
					},
				},
			},
		}

		lssscheduling.SortServiceTargetConfigs(configs)
		Expect(configs).To(HaveLen(3))
		Expect(configs[0].GetName()).To(Equal("second"))
		Expect(configs[1].GetName()).To(Equal("first"))
		Expect(configs[2].GetName()).To(Equal("third"))
	})
})

var _ = Describe("PickServiceTargetConfig", func() {

	It("should pick the config with highest priority", func() {
		configs := []*lssv1alpha1.ServiceTargetConfig{
			{
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 20,
				},
			},
			{
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 10,
				},
			},
			{
				Spec: lssv1alpha1.ServiceTargetConfigSpec{
					Priority: 30,
				},
			},
		}

		conf, err := lssscheduling.PickServiceTargetConfig(configs)
		Expect(err).NotTo(HaveOccurred())
		Expect(conf.Spec.Priority).To(Equal(int64(30)))
	})

	It("should fail if no configs are available", func() {
		configs := make([]*lssv1alpha1.ServiceTargetConfig, 0)
		_, err := lssscheduling.PickServiceTargetConfig(configs)
		Expect(err).To(HaveOccurred())
	})
})
