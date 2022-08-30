// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployments_test

import (
	"context"
	"fmt"
	"time"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lsserrors "github.com/gardener/landscaper-service/pkg/apis/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	kutil "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	deploymentscontroller "github.com/gardener/landscaper-service/pkg/controllers/landscaperdeployments"
	"github.com/gardener/landscaper-service/pkg/operation"
	testutils "github.com/gardener/landscaper-service/test/utils"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

var _ = Describe("SortServiceTargetConfigs", func() {
	It("should sort descending by priority", func() {
		configs := &lssv1alpha1.ServiceTargetConfigList{
			Items: []lssv1alpha1.ServiceTargetConfig{
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
			},
		}

		deploymentscontroller.SortServiceTargetConfigs(configs)
		Expect(configs.Items).To(HaveLen(3))
		Expect(configs.Items[0].Spec.Priority).To(Equal(int64(30)))
		Expect(configs.Items[1].Spec.Priority).To(Equal(int64(20)))
		Expect(configs.Items[2].Spec.Priority).To(Equal(int64(10)))
	})

	It("should sort ascending by usage", func() {
		configs := &lssv1alpha1.ServiceTargetConfigList{
			Items: []lssv1alpha1.ServiceTargetConfig{
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
			},
		}

		deploymentscontroller.SortServiceTargetConfigs(configs)
		Expect(configs.Items).To(HaveLen(3))
		Expect(configs.Items[0].GetName()).To(Equal("second"))
		Expect(configs.Items[1].GetName()).To(Equal("first"))
		Expect(configs.Items[2].GetName()).To(Equal("third"))
	})

	It("should sort by priority and usage", func() {
		configs := &lssv1alpha1.ServiceTargetConfigList{
			Items: []lssv1alpha1.ServiceTargetConfig{
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
			},
		}

		deploymentscontroller.SortServiceTargetConfigs(configs)
		Expect(configs.Items).To(HaveLen(3))
		Expect(configs.Items[0].GetName()).To(Equal("second"))
		Expect(configs.Items[1].GetName()).To(Equal("first"))
		Expect(configs.Items[2].GetName()).To(Equal("third"))
	})
})

var _ = Describe("Reconcile", func() {
	var (
		op    *operation.Operation
		ctrl  *deploymentscontroller.Controller
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

	It("should set finalizer and update observed generation", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())

		deployment := state.GetDeployment("test")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		Expect(kutil.HasFinalizer(deployment, lssv1alpha1.LandscaperServiceFinalizer)).To(BeTrue())
		Expect(deployment.Status.ObservedGeneration).To(Equal(int64(1)))
	})

	It("should select target configuration and create instance", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test2")
		Expect(err).ToNot(HaveOccurred())

		deployment := state.GetDeployment("test")
		config := state.GetConfig("config2")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		Expect(deployment.Status.InstanceRef).ToNot(BeNil())

		instance := &lssv1alpha1.Instance{}
		err = testenv.Client.Get(ctx, types.NamespacedName{Name: deployment.Status.InstanceRef.Name, Namespace: deployment.Status.InstanceRef.Namespace}, instance)
		Expect(err).ToNot(HaveOccurred())
		Expect(instance.Spec.ServiceTargetConfigRef.Name).To(Equal("config2"))
		Expect(instance.Spec.LandscaperConfiguration).To(Equal(deployment.Spec.LandscaperConfiguration))
		Expect(instance.Spec.TenantId).To(Equal(deployment.Spec.TenantId))
		Expect(instance.Spec.ID).To(MatchRegexp("[a-f0-9]+"))
		Expect(instance.Spec.ID).To(HaveLen(8))

		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(config), config)).To(Succeed())
		Expect(config.Status.InstanceRefs).To(HaveLen(1))
		Expect(config.Status.InstanceRefs[0].Name).To(Equal(instance.Name))
		Expect(config.Status.InstanceRefs[0].Namespace).To(Equal(instance.Namespace))
	})

	It("should not create an instance when no target configuration is available", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test3")
		Expect(err).ToNot(HaveOccurred())

		deployment := state.GetDeployment("test")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		testutils.ShouldNotReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(deployment.Status.InstanceRef).To(BeNil())
	})

	It("should mutate an existing instance", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test2")
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
		uid := instance.Spec.ID

		deployment.Spec.LandscaperConfiguration.Deployers = []string{
			"foo",
		}
		Expect(testenv.Client.Update(ctx, deployment)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		err = testenv.Client.Get(ctx, types.NamespacedName{Name: deployment.Status.InstanceRef.Name, Namespace: deployment.Status.InstanceRef.Namespace}, instance)
		Expect(err).ToNot(HaveOccurred())
		Expect(instance.Spec.LandscaperConfiguration).To(Equal(deployment.Spec.LandscaperConfiguration))
		Expect(instance.Spec.ID).To(Equal(uid))
	})

	It("should not create instances with duplicated ids", func() {
		var err error

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test5")
		Expect(err).ToNot(HaveOccurred())

		deployment := state.GetDeployment("test")
		existingInstance := state.GetInstance("existing")

		callCount := 0
		uniqueId := "eb08fabb"
		ctrl.UniqueIDFunc = func() string {
			var id string
			if callCount == 0 {
				id = existingInstance.Spec.ID
			} else {
				id = uniqueId
			}
			callCount += 1
			return id
		}

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		Expect(deployment.Status.InstanceRef).ToNot(BeNil())

		instance := &lssv1alpha1.Instance{}
		err = testenv.Client.Get(ctx, types.NamespacedName{Name: deployment.Status.InstanceRef.Name, Namespace: deployment.Status.InstanceRef.Namespace}, instance)
		Expect(err).ToNot(HaveOccurred())

		Expect(instance.Spec.ID).To(Equal(uniqueId))
	})

	It("should handle reconcile errors", func() {
		var (
			err       error
			operation = "Reconcile"
			reason    = "failed to reconcile"
			message   = "error message"
		)

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test5")
		Expect(err).ToNot(HaveOccurred())

		deployment := state.GetDeployment("test")

		ctrl.ReconcileFunc = func(ctx context.Context, deployment *lssv1alpha1.LandscaperDeployment) error {
			return lsserrors.NewWrappedError(fmt.Errorf(reason), operation, reason, message)
		}

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
		testutils.ShouldNotReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())

		Expect(deployment.Status.LastError).ToNot(BeNil())
		Expect(deployment.Status.LastError.Operation).To(Equal(operation))
		Expect(deployment.Status.LastError.Reason).To(Equal(reason))
		Expect(deployment.Status.LastError.Message).To(Equal(message))
		Expect(deployment.Status.LastError.LastUpdateTime.Time).Should(BeTemporally("==", deployment.Status.LastError.LastTransitionTime.Time))

		time.Sleep(2 * time.Second)

		message = "error message updated"

		testutils.ShouldNotReconcile(ctx, ctrl, testutils.RequestFromObject(deployment))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(deployment), deployment)).To(Succeed())

		Expect(deployment.Status.LastError).ToNot(BeNil())
		Expect(deployment.Status.LastError.Operation).To(Equal(operation))
		Expect(deployment.Status.LastError.Reason).To(Equal(reason))
		Expect(deployment.Status.LastError.Message).To(Equal(message))
		Expect(deployment.Status.LastError.LastUpdateTime.Time).Should(BeTemporally(">", deployment.Status.LastError.LastTransitionTime.Time))
	})
})
