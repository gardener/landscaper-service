// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package healthwatcher_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	kutil "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/operation"
	"github.com/gardener/landscaper-service/test/utils/envtest"

	healthwatcher "github.com/gardener/landscaper-service/pkg/controllers/healthwatcher"
	testutils "github.com/gardener/landscaper-service/test/utils"
)

type TestServiceTargetKubeClientExtractor struct{}

func (e *TestServiceTargetKubeClientExtractor) GetKubeClientFromServiceTargetConfig(ctx context.Context, name string, namespace string, client client.Client) (client.Client, error) {
	// return the original kubeclient to fake a target cluster being the core cluster
	return client, nil
}

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
		ctrl = healthwatcher.NewTestActuator(*op, &TestServiceTargetKubeClientExtractor{}, logging.Discard())
	})

	AfterEach(func() {
		defer ctx.Done()
		if state != nil {
			Expect(testenv.CleanupResources(ctx, state)).ToNot(HaveOccurred())
		}
	})

	It("should add self monitoring status as successful", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())
		op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace = state.Namespace
		op.Config().AvailabilityMonitoring.SelfLandscaperNamespace = state.Namespace

		//set lastUpdateTime of LsHealthCheck to recent
		lsHealthObject := state.GetLsHealthCheck("default")
		lsHealthObject.LastUpdateTime = v1.Now()
		Expect(testenv.Client.Update(ctx, lsHealthObject)).To(Succeed())

		availabilityCollection := state.GetAvailabilityCollection("availability")
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(availabilityCollection))

		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(availabilityCollection), availabilityCollection)).To(Succeed())
		Expect(availabilityCollection.Status.Self.Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusOk)))
	})

	It("should add self monitoring status as failed due to timeout", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())
		op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace = state.Namespace
		op.Config().AvailabilityMonitoring.SelfLandscaperNamespace = state.Namespace

		//set lastUpdateTime of LsHealthCheck to a timed-out value
		lsHealthObject := state.GetLsHealthCheck("default")
		lsHealthObject.LastUpdateTime = v1.Time{Time: v1.Now().Add(time.Minute * -6)}
		Expect(testenv.Client.Update(ctx, lsHealthObject)).To(Succeed())

		availabilityCollection := state.GetAvailabilityCollection("availability")
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(availabilityCollection))

		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(availabilityCollection), availabilityCollection)).To(Succeed())
		Expect(availabilityCollection.Status.Self.Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusFailed)))
		Expect(availabilityCollection.Status.Self.FailedReason).To(ContainSubstring("timeout"))
	})

	It("should add self monitoring status as success with lshealthcheck failed state if it is not in timoueout", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test2")
		Expect(err).ToNot(HaveOccurred())
		op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace = state.Namespace
		op.Config().AvailabilityMonitoring.SelfLandscaperNamespace = state.Namespace

		//set lastUpdateTime of LsHealthCheck to recent
		lsHealthObject := state.GetLsHealthCheck("default")
		lsHealthObject.LastUpdateTime = v1.Now()
		Expect(testenv.Client.Update(ctx, lsHealthObject)).To(Succeed())

		availabilityCollection := state.GetAvailabilityCollection("availability")
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(availabilityCollection))

		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(availabilityCollection), availabilityCollection)).To(Succeed())
		Expect(availabilityCollection.Status.Self.Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusOk)))
		Expect(availabilityCollection.Status.Self.FailedReason).To(ContainSubstring("to transition to status=Failed"))
	})
	It("should add self monitoring status as failed due to lshealthcheck failed state and in timeout", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test2")
		Expect(err).ToNot(HaveOccurred())
		op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace = state.Namespace
		op.Config().AvailabilityMonitoring.SelfLandscaperNamespace = state.Namespace

		//set lastUpdateTime of LsHealthCheck to recent
		lsHealthObject := state.GetLsHealthCheck("default")
		lsHealthObject.LastUpdateTime = v1.Time{Time: v1.Now().Add(time.Minute * -6)}
		Expect(testenv.Client.Update(ctx, lsHealthObject)).To(Succeed())

		availabilityCollection := state.GetAvailabilityCollection("availability")
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(availabilityCollection))

		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(availabilityCollection), availabilityCollection)).To(Succeed())
		Expect(availabilityCollection.Status.Self.Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusFailed)))
		Expect(availabilityCollection.Status.Self.FailedReason).To(ContainSubstring("problems"))
		Expect(availabilityCollection.Status.Self.FailedReason).To(ContainSubstring("timeout - failed recovering from failed state within time"))
	})

	It("should collect lshealthcheck from two successful instances", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test3")
		Expect(err).ToNot(HaveOccurred())
		op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace = state.Namespace
		op.Config().AvailabilityMonitoring.SelfLandscaperNamespace = state.Namespace

		//set lastUpdateTime of LsHealthCheck to recent
		lsHealthObject := state.GetLsHealthCheck("default")
		lsHealthObject.LastUpdateTime = v1.Now()
		Expect(testenv.Client.Update(ctx, lsHealthObject)).To(Succeed())

		lshealthcheck1 := state.GetLsHealthCheckInNamespace("default", fmt.Sprintf("instance1namespace-%s", state.Namespace))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(lshealthcheck1), lshealthcheck1)).To(Succeed())
		lshealthcheck1.LastUpdateTime = v1.Now()
		Expect(testenv.Client.Update(ctx, lshealthcheck1)).To(Succeed())

		lshealthcheck2 := state.GetLsHealthCheckInNamespace("default", fmt.Sprintf("instance2namespace-%s", state.Namespace))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(lshealthcheck2), lshealthcheck2)).To(Succeed())
		lshealthcheck2.LastUpdateTime = v1.Now()
		Expect(testenv.Client.Update(ctx, lshealthcheck2)).To(Succeed())

		availabilityCollection := state.GetAvailabilityCollection("availability3")
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(availabilityCollection))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(availabilityCollection), availabilityCollection)).To(Succeed())
		Expect(time.Until(availabilityCollection.Status.LastRun.Time) > time.Minute*-1).To(Equal(true)) //last reported is up-to-data
		Expect(len(availabilityCollection.Status.Instances)).To(Equal(2))
		Expect(availabilityCollection.Status.Instances[0].Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusOk)))
		Expect(availabilityCollection.Status.Instances[1].Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusOk)))
	})

	It("should collect lshealthcheck from one successful and one failed (and timeouted) instances", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test3")
		Expect(err).ToNot(HaveOccurred())
		op.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace = state.Namespace
		op.Config().AvailabilityMonitoring.SelfLandscaperNamespace = state.Namespace

		//set lastUpdateTime of LsHealthCheck to recent
		lsHealthObject := state.GetLsHealthCheck("default")
		lsHealthObject.LastUpdateTime = v1.Now()
		Expect(testenv.Client.Update(ctx, lsHealthObject)).To(Succeed())

		lshealthcheck1 := state.GetLsHealthCheckInNamespace("default", fmt.Sprintf("instance1namespace-%s", state.Namespace))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(lshealthcheck1), lshealthcheck1)).To(Succeed())
		lshealthcheck1.LastUpdateTime = v1.Now()
		Expect(testenv.Client.Update(ctx, lshealthcheck1)).To(Succeed())

		lshealthcheck2 := state.GetLsHealthCheckInNamespace("default", fmt.Sprintf("instance2namespace-%s", state.Namespace))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(lshealthcheck2), lshealthcheck2)).To(Succeed())
		lshealthcheck2.LastUpdateTime = v1.Time{Time: v1.Now().Add(time.Minute * -6)}
		lshealthcheck2.Status = lsv1alpha1.LsHealthCheckStatusFailed
		lshealthcheck2.Description = "problems"
		Expect(testenv.Client.Update(ctx, lshealthcheck2)).To(Succeed())

		availabilityCollection := state.GetAvailabilityCollection("availability3")
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(availabilityCollection))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(availabilityCollection), availabilityCollection)).To(Succeed())
		Expect(time.Until(availabilityCollection.Status.LastRun.Time) > time.Minute*-1).To(Equal(true)) //last reported is up-to-data
		Expect(len(availabilityCollection.Status.Instances)).To(Equal(2))
		Expect(availabilityCollection.Status.Instances[0].Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusOk)))
		Expect(availabilityCollection.Status.Instances[1].Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusFailed)))
		Expect(availabilityCollection.Status.Instances[1].FailedReason).To(ContainSubstring("problems"))
	})

})
var _ = Describe("failed/succeded state handling", func() {
	It("should set status to successful if not timeout and lshealthcheck is ok", func() {
		avInstance := &lssv1alpha1.AvailabilityInstance{}
		lsHealthChecks := &lsv1alpha1.LsHealthCheckList{
			Items: []lsv1alpha1.LsHealthCheck{
				{
					Status:         lsv1alpha1.LsHealthCheckStatusOk,
					LastUpdateTime: v1.Now(),
				},
			},
		}
		timeout := time.Minute * 5
		healthwatcher.ExportTransferLsHealthCheckStatusToAvailabilityInstance(avInstance, lsHealthChecks, timeout)
		Expect(avInstance.Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusOk)))
	})
	It("should set status to successful if lshealthcheck is failed but not timeouted", func() {
		avInstance := &lssv1alpha1.AvailabilityInstance{}
		lsHealthChecks := &lsv1alpha1.LsHealthCheckList{
			Items: []lsv1alpha1.LsHealthCheck{
				{
					Status:         lsv1alpha1.LsHealthCheckStatusFailed,
					Description:    "Problem Description",
					LastUpdateTime: v1.Now(),
				},
			},
		}
		timeout := time.Minute * 5
		healthwatcher.ExportTransferLsHealthCheckStatusToAvailabilityInstance(avInstance, lsHealthChecks, timeout)
		Expect(avInstance.Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusOk)))
		Expect(avInstance.FailedReason).To(ContainSubstring("to transition to status=Failed"))
	})
	It("should set status to failed if timeout but lshealthcheck is ok", func() {
		avInstance := &lssv1alpha1.AvailabilityInstance{}
		lsHealthChecks := &lsv1alpha1.LsHealthCheckList{
			Items: []lsv1alpha1.LsHealthCheck{
				{
					Status:         lsv1alpha1.LsHealthCheckStatusOk,
					LastUpdateTime: v1.Time{Time: v1.Now().Add(time.Minute * -6)},
				},
			},
		}
		timeout := time.Minute * 5
		healthwatcher.ExportTransferLsHealthCheckStatusToAvailabilityInstance(avInstance, lsHealthChecks, timeout)
		Expect(avInstance.Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusFailed)))
		Expect(avInstance.FailedReason).To(ContainSubstring("timeout"))
	})
	It("should set status to failed if timeout and lshealthcheck is failed", func() {
		avInstance := &lssv1alpha1.AvailabilityInstance{}
		lsHealthChecks := &lsv1alpha1.LsHealthCheckList{
			Items: []lsv1alpha1.LsHealthCheck{
				{
					Status:         lsv1alpha1.LsHealthCheckStatusFailed,
					Description:    "Problem Description",
					LastUpdateTime: v1.Time{Time: v1.Now().Add(time.Minute * -6)},
				},
			},
		}
		timeout := time.Minute * 5
		healthwatcher.ExportTransferLsHealthCheckStatusToAvailabilityInstance(avInstance, lsHealthChecks, timeout)
		Expect(avInstance.Status).To(Equal(string(lsv1alpha1.LsHealthCheckStatusFailed)))
		Expect(avInstance.FailedReason).To(ContainSubstring("timeout"))
		Expect(avInstance.FailedReason).To(ContainSubstring("Problem Description"))
	})
})
