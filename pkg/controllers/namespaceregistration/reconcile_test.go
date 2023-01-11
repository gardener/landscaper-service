// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package namespaceregistration_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper-service/pkg/controllers/namespaceregistration"
	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"
	"github.com/gardener/landscaper-service/pkg/operation"
	testutils "github.com/gardener/landscaper-service/test/utils"
	"github.com/gardener/landscaper-service/test/utils/envtest"
	kutil "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
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
		ctrl = namespaceregistration.NewTestActuator(*op, logging.Discard())
	})

	AfterEach(func() {
		defer ctx.Done()
		if state != nil {
			Expect(testenv.CleanupResources(ctx, state)).ToNot(HaveOccurred())
		}
	})

	It("should add finalizer on reconcile", func() {
		var err error

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())

		namespaceRegistration := state.GetNamespaceRegistration("test-namespace-1")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(namespaceRegistration), namespaceRegistration)).To(Succeed())
		Expect(len(namespaceRegistration.Finalizers)).To(Equal(1))
		Expect(namespaceRegistration.Finalizers[0]).To(Equal(lssv1alpha1.LandscaperServiceFinalizer))
	})

	It("should create namespace with role/rolebinding on namespaceregistration create", func() {
		var err error

		lsUserNamespace := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ls-user"}}
		Expect(testenv.Client.Create(ctx, &lsUserNamespace)).To(Succeed())
		subjectList := lssv1alpha1.SubjectList{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "subjects",
				Namespace: "ls-user",
			},
			Spec: lssv1alpha1.SubjectListSpec{
				Subjects: []lssv1alpha1.Subject{
					{
						Kind: "User",
						Name: "testuser",
					},
				},
			},
		}
		Expect(testenv.Client.Create(ctx, &subjectList)).To(Succeed())

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())

		namespaceRegistration := state.GetNamespaceRegistration("test-namespace-1")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(namespaceRegistration), namespaceRegistration)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(namespaceRegistration), namespaceRegistration)).To(Succeed())
		Expect(namespaceRegistration.Status.Phase).To(Equal("Completed"))

		// check for namespace being created
		namespace := corev1.Namespace{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: namespaceRegistration.Name}, &namespace)).To(Succeed())

		//check for role being created
		role := rbacv1.Role{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_IN_NAMESPACE, Namespace: namespace.Name}, &role)).To(Succeed())

		//check for rolebinding being created
		rolebinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: namespace.Name}, &rolebinding)).To(Succeed())
		Expect(rolebinding.RoleRef.Name).To(Equal(subjectsync.USER_ROLE_IN_NAMESPACE))
	})
})
