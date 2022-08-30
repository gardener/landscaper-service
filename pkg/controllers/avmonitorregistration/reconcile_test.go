// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package avmonitorregistration_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper-service/pkg/operation"
	"github.com/gardener/landscaper-service/test/utils/envtest"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"

	avmonitorregistration "github.com/gardener/landscaper-service/pkg/controllers/avmonitorregistration"
	testutils "github.com/gardener/landscaper-service/test/utils"
)

var _ = Describe("Reconcile", func() {
	var (
		op    *operation.Operation
		ctrl  reconcile.Reconciler
		ctx   context.Context
		state *envtest.State
	)

	BeforeEach(func() {
		ctx = context.Background()
		op = operation.NewOperation(testenv.Client, envtest.LandscaperServiceScheme, testutils.DefaultControllerConfiguration())
		ctrl = avmonitorregistration.NewTestActuator(*op, logging.Discard())
	})

	AfterEach(func() {
		defer ctx.Done()
		if state != nil {
			Expect(testenv.CleanupResources(ctx, state)).ToNot(HaveOccurred())
		}
	})

	It("should create avmonitoringcollection with empty spec for missing installation reference in instance", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())
		op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace = state.Namespace

		instance := state.GetInstance("test")
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		availabilitycollection := &lssv1alpha1.AvailabilityCollection{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Namespace: op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace, Name: op.Config().AvailabilityMonitoring.AvailabilityCollectionName}, availabilitycollection)).To(Succeed())
		Expect(len(availabilitycollection.Spec.InstanceRefs)).To(Equal(0))
	})

	It("should create avmonitoringcollection with one spec for instance", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test2")
		Expect(err).ToNot(HaveOccurred())
		op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace = state.Namespace

		instance := state.GetInstance("test")
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		availabilitycollection := &lssv1alpha1.AvailabilityCollection{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Namespace: op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace, Name: op.Config().AvailabilityMonitoring.AvailabilityCollectionName}, availabilitycollection)).To(Succeed())
		Expect(len(availabilitycollection.Spec.InstanceRefs)).To(Equal(1))
	})
})
