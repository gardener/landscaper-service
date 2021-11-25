// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	kutil "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	instancescontroller "github.com/gardener/landscaper-service/pkg/controllers/instances"
	"github.com/gardener/landscaper-service/pkg/operation"
	testutils "github.com/gardener/landscaper-service/test/utils"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

var _ = Describe("Delete", func() {
	var (
		op    *operation.Operation
		ctrl  reconcile.Reconciler
		ctx   context.Context
		state *envtest.State
	)

	BeforeEach(func() {
		ctx = context.Background()
		op = operation.NewOperation(logr.Discard(), testenv.Client, envtest.LandscaperServiceScheme)
		ctrl = instancescontroller.NewTestActuator(*op)
	})

	AfterEach(func() {
		defer ctx.Done()
		if state != nil {
			Expect(testenv.CleanupResources(ctx, state)).ToNot(HaveOccurred())
		}
	})

	It("should remove the finalizer", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/delete/test1")
		Expect(err).ToNot(HaveOccurred())

		instance := state.GetInstance("test")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(instance), instance)).To(Succeed())
		Expect(kutil.HasFinalizer(instance, lssv1alpha1.LandscaperServiceFinalizer)).To(BeTrue())

		Expect(testenv.Client.Delete(ctx, instance)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		Expect(testenv.WaitForObjectToBeDeleted(ctx, instance, 5*time.Second)).To(Succeed())
	})

	It("should remove the instance reference from the service target config", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/delete/test2")
		Expect(err).ToNot(HaveOccurred())

		instance := state.GetInstance("test")
		config := state.GetConfig("default")

		Expect(config.Status.InstanceRefs).To(HaveLen(2))

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(instance), instance)).To(Succeed())
		Expect(kutil.HasFinalizer(instance, lssv1alpha1.LandscaperServiceFinalizer)).To(BeTrue())

		Expect(testenv.Client.Delete(ctx, instance)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		Expect(testenv.WaitForObjectToBeDeleted(ctx, instance, 5*time.Second)).To(Succeed())

		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(config), config)).To(Succeed())
		Expect(config.Status.InstanceRefs).To(HaveLen(1))
	})

	It("should remove the associated target and installation", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/delete/test3")
		Expect(err).ToNot(HaveOccurred())

		instance := state.GetInstance("test")
		target := state.GetTarget("test")
		installation := state.GetInstallation("test")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(instance), instance)).To(Succeed())
		Expect(kutil.HasFinalizer(instance, lssv1alpha1.LandscaperServiceFinalizer)).To(BeTrue())

		Expect(testenv.Client.Delete(ctx, instance)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		// installation
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		// target
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		Expect(testenv.WaitForObjectToBeDeleted(ctx, instance, 5*time.Second)).To(Succeed())
		Expect(testenv.WaitForObjectToBeDeleted(ctx, target, 5*time.Second)).To(Succeed())
		Expect(testenv.WaitForObjectToBeDeleted(ctx, installation, 5*time.Second)).To(Succeed())
	})
})
