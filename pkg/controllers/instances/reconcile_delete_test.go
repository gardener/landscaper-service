// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances_test

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	lsserrors "github.com/gardener/landscaper-service/pkg/apis/provisioning/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	kutil "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	"github.com/gardener/landscaper-service/pkg/apis/constants"
	provisioningv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/provisioning/v1alpha2"
	instancescontroller "github.com/gardener/landscaper-service/pkg/controllers/instances"
	"github.com/gardener/landscaper-service/pkg/operation"
	testutils "github.com/gardener/landscaper-service/test/utils"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

func expectObjectDeleted(err error, obj client.Object) {
	if err == nil {
		Expect(obj.GetDeletionTimestamp().IsZero()).To(BeFalse())
	} else {
		Expect(apierrors.IsNotFound(err)).To(BeTrue())
	}
}

var _ = Describe("Delete", func() {
	var (
		op    *operation.Operation
		ctrl  *instancescontroller.Controller
		ctx   context.Context
		state *envtest.State
	)

	BeforeEach(func() {
		ctx = context.Background()
		op = operation.NewOperation(testenv.Client, envtest.LandscaperServiceScheme, testutils.DefaultControllerConfiguration())
		ctrl = instancescontroller.NewTestActuator(*op, logging.Discard())
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
		Expect(kutil.HasFinalizer(instance, constants.LandscaperServiceFinalizer)).To(BeTrue())

		Expect(testenv.Client.Delete(ctx, instance)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, instance, 5*time.Second)).To(Succeed())
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
		Expect(kutil.HasFinalizer(instance, constants.LandscaperServiceFinalizer)).To(BeTrue())

		Expect(testenv.Client.Delete(ctx, instance)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, instance, 5*time.Second)).To(Succeed())

		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(config), config)).To(Succeed())
		Expect(config.Status.InstanceRefs).To(HaveLen(1))
	})

	It("should remove the associated context, target and installation", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/delete/test3")
		Expect(err).ToNot(HaveOccurred())

		instance := state.GetInstance("test")
		target := state.GetTarget("test")
		gardenerSa := state.GetTarget("test-gardener-sa")
		installation := state.GetInstallation("test")
		context := state.GetContext("test")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(instance), instance)).To(Succeed())
		Expect(kutil.HasFinalizer(instance, constants.LandscaperServiceFinalizer)).To(BeTrue())

		Expect(testenv.Client.Delete(ctx, instance)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		// installation
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		// target
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		// gardener service account
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		// context
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, instance, 5*time.Second)).To(Succeed())
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, target, 5*time.Second)).To(Succeed())
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, gardenerSa, 5*time.Second)).To(Succeed())
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, installation, 5*time.Second)).To(Succeed())
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, context, 5*time.Second)).To(Succeed())
	})

	It("should remove the secrets referenced by the context", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/delete/test4")
		Expect(err).ToNot(HaveOccurred())

		instance := state.GetInstance("test")
		regpullsecret1 := state.GetSecret("regpullsecret1")
		regpullsecret2 := state.GetSecret("regpullsecret2")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(instance), instance)).To(Succeed())
		Expect(kutil.HasFinalizer(instance, constants.LandscaperServiceFinalizer)).To(BeTrue())

		Expect(testenv.Client.Delete(ctx, instance)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		// context
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, regpullsecret1, 5*time.Second)).To(Succeed())
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, regpullsecret2, 5*time.Second)).To(Succeed())
	})

	It("should handle delete errors", func() {
		var (
			err       error
			operation = "Delete"
			reason    = "failed to delete"
			message   = "error message"
		)

		state, err = testenv.InitResources(ctx, "./testdata/delete/test4")
		Expect(err).ToNot(HaveOccurred())

		instance := state.GetInstance("test")

		ctrl.HandleDeleteFunc = func(ctx context.Context, deployment *provisioningv1alpha2.Instance) (reconcile.Result, error) {
			return reconcile.Result{}, lsserrors.NewWrappedError(fmt.Errorf(reason), operation, reason, message)
		}

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(instance), instance)).To(Succeed())

		Expect(testenv.Client.Delete(ctx, instance)).To(Succeed())

		testutils.ShouldNotReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(instance), instance)).To(Succeed())

		Expect(instance.Status.LastError).ToNot(BeNil())
		Expect(instance.Status.LastError.Operation).To(Equal(operation))
		Expect(instance.Status.LastError.Reason).To(Equal(reason))
		Expect(instance.Status.LastError.Message).To(Equal(message))
		Expect(instance.Status.LastError.LastUpdateTime.Time).Should(BeTemporally("==", instance.Status.LastError.LastTransitionTime.Time))

		time.Sleep(2 * time.Second)

		message = "error message updated"

		testutils.ShouldNotReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(instance), instance)).To(Succeed())

		Expect(instance.Status.LastError).ToNot(BeNil())
		Expect(instance.Status.LastError.Operation).To(Equal(operation))
		Expect(instance.Status.LastError.Reason).To(Equal(reason))
		Expect(instance.Status.LastError.Message).To(Equal(message))
		Expect(instance.Status.LastError.LastUpdateTime.Time).Should(BeTemporally(">", instance.Status.LastError.LastTransitionTime.Time))
	})

	It("should delete the target cluster namespace and rbac objects", func() {
		var err error
		state, err = testenv.InitResources(ctx, "./testdata/delete/test5")
		Expect(err).ToNot(HaveOccurred())

		instance := state.GetInstance("test")

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(instance), instance)).To(Succeed())

		targetName := fmt.Sprintf("%s-%s", instance.Spec.TenantId, instance.Spec.ID)

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: targetName,
			},
		}
		Expect(testenv.Client.Create(ctx, namespace)).To(Succeed())

		serviceAccount := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "laas-test-sa",
				Namespace: targetName,
			},
		}

		Expect(testenv.Client.Create(ctx, serviceAccount)).To(Succeed())

		clusterRole := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("helm-%s-tmp", targetName),
			},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"*"},
					Resources: []string{"*"},
					Verbs:     []string{"*"},
				},
			},
		}
		Expect(testenv.Client.Create(ctx, clusterRole)).To(Succeed())

		clusterRoleBinding := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("helm-%s-rb-tmp", targetName),
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     clusterRole.GetName(),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      serviceAccount.GetName(),
					Namespace: serviceAccount.GetNamespace(),
				},
			},
		}

		Expect(testenv.Client.Create(ctx, clusterRoleBinding)).To(Succeed())

		Expect(testenv.Client.Delete(ctx, instance)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(instance))

		err = testenv.Client.Get(ctx, client.ObjectKeyFromObject(clusterRoleBinding), clusterRoleBinding)
		expectObjectDeleted(err, clusterRoleBinding)

		err = testenv.Client.Get(ctx, client.ObjectKeyFromObject(clusterRole), clusterRole)
		expectObjectDeleted(err, clusterRole)

		err = testenv.Client.Get(ctx, client.ObjectKeyFromObject(namespace), namespace)
		expectObjectDeleted(err, namespace)
	})
})
