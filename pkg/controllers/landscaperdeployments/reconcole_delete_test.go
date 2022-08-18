// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployments_test

import (
	"context"
	"time"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	kutil "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	deploymentscontroller "github.com/gardener/landscaper-service/pkg/controllers/landscaperdeployments"
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
		op = operation.NewOperation(testenv.Client, envtest.LandscaperServiceScheme, testutils.DefaultControllerConfiguration())
		ctrl = deploymentscontroller.NewTestActuator(*op, logging.Discard())
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

		deployment := state.GetDeployment("test")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		Expect(kutil.HasFinalizer(deployment, lssv1alpha1.LandscaperServiceFinalizer)).To(BeTrue())

		Expect(testenv.Client.Delete(ctx, deployment)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))

		Expect(testenv.WaitForObjectToBeDeleted(ctx, deployment, 5*time.Second)).To(Succeed())
	})

	It("should remove the referenced instance", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/delete/test2")
		Expect(err).ToNot(HaveOccurred())

		deployment := state.GetDeployment("test")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		Expect(deployment.Status.InstanceRef).ToNot(BeNil())

		instance := &lssv1alpha1.Instance{}
		err = testenv.Client.Get(ctx, types.NamespacedName{Name: deployment.Status.InstanceRef.Name, Namespace: deployment.Status.InstanceRef.Namespace}, instance)
		Expect(err).ToNot(HaveOccurred())

		Expect(testenv.Client.Delete(ctx, deployment)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))

		Expect(testenv.WaitForObjectToBeDeleted(ctx, deployment, 5*time.Second)).To(Succeed())
		Expect(testenv.WaitForObjectToBeDeleted(ctx, instance, 5*time.Second)).To(Succeed())
	})
})
