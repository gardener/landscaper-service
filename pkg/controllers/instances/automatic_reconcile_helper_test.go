// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/landscaper-service/pkg/controllers/instances"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

var _ = Describe("AutomaticReconcileHelper", func() {
	var (
		state *envtest.State
		ctx   context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	AfterEach(func() {
		defer ctx.Done()
		if state != nil {
			Expect(testenv.CleanupResources(ctx, state)).ToNot(HaveOccurred())
		}
	})

	It("handle the automatic reconciling functionality correctly", func() {
		const automaticReconcileInterval = 2 * time.Minute

		state, err := testenv.InitResources(ctx, "./testdata/automatic_reconcile_helper/test1")
		Expect(err).ToNot(HaveOccurred())

		instance := state.GetInstance("test")
		Expect(instance).ToNot(BeNil())

		rh := instances.NewAutomaticReconcileHelper(testenv.Client, clock.RealClock{})

		/// 1. Instance has a changed spec.
		// Expect the last reconcile time to be set and the requeue interval to equal the duration set in the Instance spec.
		By("instance has changed", func() {
			oldInstance := instance.DeepCopy()
			instance.Spec.LandscaperConfiguration.Deployers = []string{
				"customDeployer",
			}

			Expect(testenv.Client.Update(ctx, instance)).To(Succeed())

			result, err := rh.ComputeAutomaticReconcile(ctx, instance, oldInstance, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).ToNot(BeNil())
			Expect(result.Requeue).To(BeTrue())
			Expect(result.RequeueAfter).To(Equal(automaticReconcileInterval))

			now := time.Now()

			Expect(testenv.Client.Get(ctx, client.ObjectKeyFromObject(instance), instance)).To(Succeed())
			Expect(instance.Status.AutomaticReconcileStatus).ToNot(BeNil())
			Expect(now.Compare(instance.Status.AutomaticReconcileStatus.LastReconcileTime.Time) >= 0).To(BeTrue())
		})

		time.Sleep(1 * time.Second)

		/// 2. Instance has no changed spec.
		/// Expect the last reconcile time to be not updated and the requeue interval to be less than the duration set
		/// in the instance spec.
		By("instance has not changed", func() {
			Expect(testenv.Client.Get(ctx, client.ObjectKeyFromObject(instance), instance)).To(Succeed())
			lastReconcileTime := instance.Status.AutomaticReconcileStatus.LastReconcileTime
			oldInstance := instance.DeepCopy()

			result, err := rh.ComputeAutomaticReconcile(ctx, instance, oldInstance, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).ToNot(BeNil())
			Expect(result.Requeue).To(BeTrue())
			Expect(result.RequeueAfter < automaticReconcileInterval).To(BeTrue())

			Expect(testenv.Client.Get(ctx, client.ObjectKeyFromObject(instance), instance)).To(Succeed())
			Expect(lastReconcileTime.Time.Equal(instance.Status.AutomaticReconcileStatus.LastReconcileTime.Time)).To(BeTrue())
		})

		time.Sleep(1 * time.Second)

		/// 3. Instance has a changed spec.
		// Expect the last reconcile time to be set and the requeue interval to be less than the duration set in the instance spec.
		By("instance has changed again", func() {
			Expect(testenv.Client.Get(ctx, client.ObjectKeyFromObject(instance), instance)).To(Succeed())
			lastReconcileTime := instance.Status.AutomaticReconcileStatus.LastReconcileTime
			oldInstance := instance.DeepCopy()
			instance.Spec.LandscaperConfiguration.Deployers = []string{
				"newDeployer",
			}

			Expect(testenv.Client.Update(ctx, instance)).To(Succeed())

			result, err := rh.ComputeAutomaticReconcile(ctx, instance, oldInstance, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).ToNot(BeNil())
			Expect(result.Requeue).To(BeTrue())
			Expect(result.RequeueAfter < automaticReconcileInterval).To(BeTrue())

			Expect(testenv.Client.Get(ctx, client.ObjectKeyFromObject(instance), instance)).To(Succeed())
			Expect(lastReconcileTime.Time.Equal(instance.Status.AutomaticReconcileStatus.LastReconcileTime.Time)).To(BeFalse())
		})
	})
})
