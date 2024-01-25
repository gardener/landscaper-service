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

var _ = Describe("Scheduling", func() {

	const (
		namespace1 = "test-namespace-1"

		config1 = "test-config-1"
		config2 = "test-config-2"
		config3 = "test-config-3"

		tenant1 = "test-tenant-1"
		tenant2 = "test-tenant-2"

		key1   = "test-key-1"
		value1 = "test-value-1"
		value2 = "test-value-2"
	)

	buildLandscaperDeployment := func(tenantID string, labels map[string]string) *lssv1alpha1.LandscaperDeployment {
		return &lssv1alpha1.LandscaperDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: lssv1alpha1.LandscaperDeploymentSpec{
				TenantId: tenantID,
			},
		}
	}

	buildServiceTargetConfig := func(name string, prio int64, restricted bool) *lssv1alpha1.ServiceTargetConfig {
		return &lssv1alpha1.ServiceTargetConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Namespace:   namespace1,
				Annotations: map[string]string{lssv1alpha1.ServiceTargetConfigVisibleLabelName: "true"},
			},
			Spec: lssv1alpha1.ServiceTargetConfigSpec{
				Priority:   prio,
				Restricted: restricted,
			},
		}
	}

	It("should apply the rule with highest prio", func() {
		// The scheduling contains three rules, all of which match with the LandscaperDeployment.
		// Therefore, the rule with the highest prio should win.

		serviceTargetConfigs := []lssv1alpha1.ServiceTargetConfig{
			*buildServiceTargetConfig(config1, 1, false),
			*buildServiceTargetConfig(config2, 1, false),
			*buildServiceTargetConfig(config3, 1, false),
		}

		deployment := buildLandscaperDeployment(tenant1, nil)

		scheduling := &lssv1alpha1.TargetScheduling{
			Spec: lssv1alpha1.TargetSchedulingSpec{
				Rules: []lssv1alpha1.SchedulingRule{
					{
						Priority: 4,
						ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
							{Name: config1, Namespace: namespace1},
						},
						Selector: []lssv1alpha1.Selector{
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenant1}},
						},
					},
					{
						Priority: 8, // highest prio
						ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
							{Name: config2, Namespace: namespace1},
						},
						Selector: []lssv1alpha1.Selector{
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenant1}},
						},
					},
					{
						Priority: 6,
						ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
							{Name: config3, Namespace: namespace1},
						},
						Selector: []lssv1alpha1.Selector{
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenant1}},
						},
					},
				},
			},
		}

		config, err := lssscheduling.FindServiceTargetConfig(scheduling, deployment, serviceTargetConfigs)
		Expect(err).NotTo(HaveOccurred())
		Expect(config.Name).To(Equal(config2))
	})

	It("should only apply matching rules", func() {
		// The scheduling contains three rules, but only one of them matches with the LandscaperDeployment.

		serviceTargetConfigs := []lssv1alpha1.ServiceTargetConfig{
			*buildServiceTargetConfig(config1, 1, false),
			*buildServiceTargetConfig(config2, 1, false),
			*buildServiceTargetConfig(config3, 1, false),
		}

		deployment := buildLandscaperDeployment(tenant1, map[string]string{key1: value1})

		// Only the second rule matches
		scheduling := &lssv1alpha1.TargetScheduling{
			Spec: lssv1alpha1.TargetSchedulingSpec{
				Rules: []lssv1alpha1.SchedulingRule{
					{
						Priority: 4,
						ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
							{Name: config1, Namespace: namespace1},
						},
						Selector: []lssv1alpha1.Selector{
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenant2}}, // other tenant
							{MatchLabel: &lssv1alpha1.LabelSelector{Name: key1, Value: value1}},
						},
					},
					{
						// matching rule
						Priority: 4,
						ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
							{Name: config2, Namespace: namespace1},
						},
						Selector: []lssv1alpha1.Selector{
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenant1}},
							{MatchLabel: &lssv1alpha1.LabelSelector{Name: key1, Value: value1}},
						},
					},
					{
						Priority: 4,
						ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
							{Name: config3, Namespace: namespace1},
						},
						Selector: []lssv1alpha1.Selector{
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenant1}},
							{MatchLabel: &lssv1alpha1.LabelSelector{Name: key1, Value: value2}}, // other value
						},
					},
				},
			},
		}

		config, err := lssscheduling.FindServiceTargetConfig(scheduling, deployment, serviceTargetConfigs)
		Expect(err).NotTo(HaveOccurred())
		Expect(config.Name).To(Equal(config2))
	})

	It("should only select an existing service target config", func() {
		// The scheduling contains one rule, and this rule matches with the LandscaperDeployment.
		// The rule has three ServiceTargetConfigs, but only one them exists.

		serviceTargetConfigs := []lssv1alpha1.ServiceTargetConfig{
			*buildServiceTargetConfig(config2, 1, false), // the only existing config
		}

		deployment := buildLandscaperDeployment(tenant1, map[string]string{key1: value1})

		scheduling := &lssv1alpha1.TargetScheduling{
			Spec: lssv1alpha1.TargetSchedulingSpec{
				Rules: []lssv1alpha1.SchedulingRule{
					{
						Priority: 4,
						ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
							{Name: config1, Namespace: namespace1},
							{Name: config2, Namespace: namespace1},
							{Name: config3, Namespace: namespace1},
						},
						Selector: []lssv1alpha1.Selector{
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenant1}},
						},
					},
				},
			},
		}

		config, err := lssscheduling.FindServiceTargetConfig(scheduling, deployment, serviceTargetConfigs)
		Expect(err).NotTo(HaveOccurred())
		Expect(config.Name).To(Equal(config2))
	})

	It("should pick the service target config with highest prio", func() {
		// The scheduling contains one rule, and this rule matches with the LandscaperDeployment.
		// The rule has three ServiceTargetConfigs, which all exist.
		// Therefore, the ServiceTargetConfig with the highest prio should win.

		serviceTargetConfigs := []lssv1alpha1.ServiceTargetConfig{
			*buildServiceTargetConfig(config1, 10, false),
			*buildServiceTargetConfig(config2, 30, false), // highest prio
			*buildServiceTargetConfig(config3, 20, false),
		}

		deployment := buildLandscaperDeployment(tenant1, nil)

		scheduling := &lssv1alpha1.TargetScheduling{
			Spec: lssv1alpha1.TargetSchedulingSpec{
				Rules: []lssv1alpha1.SchedulingRule{
					{
						Priority: 4,
						ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
							{Name: config1, Namespace: namespace1},
							{Name: config2, Namespace: namespace1},
							{Name: config3, Namespace: namespace1},
						},
						Selector: []lssv1alpha1.Selector{
							{MatchTenant: &lssv1alpha1.TenantSelector{ID: tenant1}},
						},
					},
				},
			},
		}

		config, err := lssscheduling.FindServiceTargetConfig(scheduling, deployment, serviceTargetConfigs)
		Expect(err).NotTo(HaveOccurred())
		Expect(config.Name).To(Equal(config2))
	})

	It("should pick an unrestricted service target config if scheduling is not configured", func() {
		// There is no scheduling defined. Therefore, the only unrestricted ServiceTargetConfig should be selected.

		serviceTargetConfigs := []lssv1alpha1.ServiceTargetConfig{
			*buildServiceTargetConfig(config1, 10, true),
			*buildServiceTargetConfig(config2, 10, false), // unrestricted
			*buildServiceTargetConfig(config3, 10, true),
		}

		deployment := buildLandscaperDeployment(tenant1, nil)

		config, err := lssscheduling.FindServiceTargetConfig(nil, deployment, serviceTargetConfigs)
		Expect(err).NotTo(HaveOccurred())
		Expect(config.Name).To(Equal(config2))
	})

	It("should pick an unrestricted service target config if no scheduling rule matches", func() {
		// The scheduling has one rule which does not match with the LandscaperDeployment.
		// Therefore, the only unrestricted ServiceTargetConfig should be selected.

		serviceTargetConfigs := []lssv1alpha1.ServiceTargetConfig{
			*buildServiceTargetConfig(config1, 10, true),
			*buildServiceTargetConfig(config2, 10, false), // unrestricted
			*buildServiceTargetConfig(config3, 10, true),
		}

		deployment := buildLandscaperDeployment(tenant1, map[string]string{key1: value1})

		scheduling := &lssv1alpha1.TargetScheduling{
			Spec: lssv1alpha1.TargetSchedulingSpec{
				Rules: []lssv1alpha1.SchedulingRule{
					{
						Priority: 4,
						ServiceTargetConfigs: []lssv1alpha1.ObjectReference{
							{Name: config3, Namespace: namespace1},
						},
						Selector: []lssv1alpha1.Selector{
							{MatchLabel: &lssv1alpha1.LabelSelector{Name: key1, Value: value2}}, // other value
						},
					},
				},
			},
		}

		config, err := lssscheduling.FindServiceTargetConfig(scheduling, deployment, serviceTargetConfigs)
		Expect(err).NotTo(HaveOccurred())
		Expect(config.Name).To(Equal(config2))
	})

	It("should pick an unrestricted service target config if the scheduling resource has no rules", func() {
		// The scheduling has no rules.
		// Therefore, the only unrestricted ServiceTargetConfig should be selected.

		serviceTargetConfigs := []lssv1alpha1.ServiceTargetConfig{
			*buildServiceTargetConfig(config1, 10, true),
			*buildServiceTargetConfig(config2, 10, false), // unrestricted
			*buildServiceTargetConfig(config3, 10, true),
		}

		deployment := buildLandscaperDeployment(tenant1, map[string]string{key1: value1})

		scheduling := &lssv1alpha1.TargetScheduling{}

		config, err := lssscheduling.FindServiceTargetConfig(scheduling, deployment, serviceTargetConfigs)
		Expect(err).NotTo(HaveOccurred())
		Expect(config.Name).To(Equal(config2))
	})
})
