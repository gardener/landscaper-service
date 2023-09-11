// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package namespaceregistration_test

import (
	"context"
	"time"

	"github.com/gardener/landscaper/apis/core/v1alpha1/helper"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kutil "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper-service/pkg/controllers/namespaceregistration"
	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"
	"github.com/gardener/landscaper-service/pkg/operation"
	testutils "github.com/gardener/landscaper-service/test/utils"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

var _ = Describe("Reconcile", func() {
	var (
		op    *operation.TargetShootSidecarOperation
		ctrl  reconcile.Reconciler
		ctx   context.Context
		state *envtest.State
	)

	BeforeEach(func() {
		ctx = context.Background()
		op = operation.NewTargetShootSidecarOperation(testenv.Client, envtest.LandscaperServiceScheme, testutils.DefaultTargetShootConfiguration())
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

		namespaceRegistration := state.GetNamespaceRegistration(subjectsync.CUSTOM_NS_PREFIX + "test-namespace-1")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(namespaceRegistration), namespaceRegistration)).To(Succeed())
		Expect(len(namespaceRegistration.Finalizers)).To(Equal(1))
		Expect(namespaceRegistration.Finalizers[0]).To(Equal(lssv1alpha1.LandscaperServiceFinalizer))
	})

	It("should create namespace with role/rolebinding on namespaceregistration create", func() {
		var err error

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())

		namespaceRegistration := state.GetNamespaceRegistration(subjectsync.CUSTOM_NS_PREFIX + "test-namespace-1")
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
		Expect(role.Rules[0].APIGroups).To(ContainElement("landscaper.gardener.cloud"))
		Expect(role.Rules[0].Resources).To(ContainElement("*"))
		Expect(role.Rules[0].Verbs).To(ContainElement("*"))
		Expect(role.Rules[1].APIGroups).To(ContainElement(""))
		Expect(role.Rules[1].Resources).To(ContainElements("secrets", "configmaps"))
		Expect(role.Rules[1].Verbs).To(ContainElement("*"))

		//check for rolebinding being created
		rolebinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: namespace.Name}, &rolebinding)).To(Succeed())
		Expect(rolebinding.RoleRef.Name).To(Equal(subjectsync.USER_ROLE_IN_NAMESPACE))
		Expect(len(rolebinding.Subjects)).To(Equal(1))
		Expect(rolebinding.Subjects[0].Kind).To(Equal("User"))
		Expect(rolebinding.Subjects[0].Name).To(Equal("testuser"))
	})

	It("should delete namespace with role/rolebinding on namespaceregistration deletion", func() {
		var err error

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test2")
		Expect(err).ToNot(HaveOccurred())

		namespaceRegistration := state.GetNamespaceRegistration(subjectsync.CUSTOM_NS_PREFIX + "test-namespace-2")
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
		Expect(len(rolebinding.Subjects)).To(Equal(1))
		Expect(rolebinding.Subjects[0].Kind).To(Equal("User"))
		Expect(rolebinding.Subjects[0].Name).To(Equal("testuser"))

		// deletion
		Expect(testenv.Client.Delete(ctx, namespaceRegistration)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, namespaceRegistration, 5*time.Second)).To(Succeed())
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(&namespace), &namespace)).To(Succeed())

		// check for namespace being deleted
		Expect(namespace.Status.Phase).To(Equal(corev1.NamespaceTerminating))
	})

	It("should delete namespace with installation", func() {
		var err error

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test3")
		Expect(err).ToNot(HaveOccurred())

		namespaceRegistration := state.GetNamespaceRegistration(subjectsync.CUSTOM_NS_PREFIX + "test-namespace-3")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(namespaceRegistration), namespaceRegistration)).To(Succeed())
		Expect(namespaceRegistration.Status.Phase).To(Equal("Completed"))

		// check for namespace being created
		namespace := corev1.Namespace{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: namespaceRegistration.Name}, &namespace)).To(Succeed())

		// add Installation to prevent delete
		compRef := &lsv1alpha1.ComponentDescriptorDefinition{
			Reference: &lsv1alpha1.ComponentDescriptorReference{
				ComponentName: "component",
				Version:       "v0.1.0",
			},
		}

		installation := &lsv1alpha1.Installation{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: namespace.Name, Namespace: namespace.Name},
			Spec: lsv1alpha1.InstallationSpec{
				ComponentDescriptor: compRef,
			},
		}

		Expect(testenv.Client.Create(ctx, installation)).To(Succeed())

		// delete
		Expect(testenv.Client.Delete(ctx, namespaceRegistration)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))

		// failed deletion
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: namespaceRegistration.Name}, &namespace)).To(Succeed())
		Expect(namespace.DeletionTimestamp.IsZero()).To(BeTrue())

		// successful deletion
		Expect(testenv.Client.Delete(ctx, installation)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, namespaceRegistration, 5*time.Second)).To(Succeed())
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(&namespace), &namespace)).To(Succeed())

		// check for namespace being deleted
		Expect(namespace.Status.Phase).To(Equal(corev1.NamespaceTerminating))
	})

	It("should delete namespace with execution", func() {
		var err error

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test4")
		Expect(err).ToNot(HaveOccurred())

		namespaceRegistration := state.GetNamespaceRegistration(subjectsync.CUSTOM_NS_PREFIX + "test-namespace-4")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(namespaceRegistration), namespaceRegistration)).To(Succeed())
		Expect(namespaceRegistration.Status.Phase).To(Equal("Completed"))

		// check for namespace being created
		namespace := corev1.Namespace{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: namespaceRegistration.Name}, &namespace)).To(Succeed())

		// add execution to prevent delete
		execution := &lsv1alpha1.Execution{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: namespace.Name, Namespace: namespace.Name},
		}

		Expect(testenv.Client.Create(ctx, execution)).To(Succeed())

		// delete
		Expect(testenv.Client.Delete(ctx, namespaceRegistration)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))

		// failed deletion
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: namespaceRegistration.Name}, &namespace)).To(Succeed())
		Expect(namespace.DeletionTimestamp.IsZero()).To(BeTrue())

		// successful deletion
		Expect(testenv.Client.Delete(ctx, execution)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, namespaceRegistration, 5*time.Second)).To(Succeed())
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(&namespace), &namespace)).To(Succeed())

		// check for namespace being deleted
		Expect(namespace.Status.Phase).To(Equal(corev1.NamespaceTerminating))
	})

	It("should delete namespace with deploy item", func() {
		var err error

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test5")
		Expect(err).ToNot(HaveOccurred())

		namespaceRegistration := state.GetNamespaceRegistration(subjectsync.CUSTOM_NS_PREFIX + "test-namespace-5")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(namespaceRegistration), namespaceRegistration)).To(Succeed())
		Expect(namespaceRegistration.Status.Phase).To(Equal("Completed"))

		// check for namespace being created
		namespace := corev1.Namespace{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: namespaceRegistration.Name}, &namespace)).To(Succeed())

		// add execution to prevent delete
		di := &lsv1alpha1.DeployItem{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: namespace.Name, Namespace: namespace.Name},
		}

		Expect(testenv.Client.Create(ctx, di)).To(Succeed())

		// delete
		Expect(testenv.Client.Delete(ctx, namespaceRegistration)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))

		// failed deletion
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: namespaceRegistration.Name}, &namespace)).To(Succeed())
		Expect(namespace.DeletionTimestamp.IsZero()).To(BeTrue())

		// successful deletion
		Expect(testenv.Client.Delete(ctx, di)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, namespaceRegistration, 5*time.Second)).To(Succeed())
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(&namespace), &namespace)).To(Succeed())

		// check for namespace being deleted
		Expect(namespace.Status.Phase).To(Equal(corev1.NamespaceTerminating))
	})

	It("should delete namespace with target sync", func() {
		var err error

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test6")
		Expect(err).ToNot(HaveOccurred())

		namespaceRegistration := state.GetNamespaceRegistration(subjectsync.CUSTOM_NS_PREFIX + "test-namespace-6")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(namespaceRegistration), namespaceRegistration)).To(Succeed())
		Expect(namespaceRegistration.Status.Phase).To(Equal("Completed"))

		// check for namespace being created
		namespace := corev1.Namespace{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: namespaceRegistration.Name}, &namespace)).To(Succeed())

		// add execution to prevent delete
		targetSync := &lsv1alpha1.TargetSync{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: namespace.Name, Namespace: namespace.Name},
		}

		Expect(testenv.Client.Create(ctx, targetSync)).To(Succeed())

		targetSync = &lsv1alpha1.TargetSync{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: namespace.Name + "-2", Namespace: namespace.Name},
		}

		Expect(testenv.Client.Create(ctx, targetSync)).To(Succeed())

		// delete
		Expect(testenv.Client.Delete(ctx, namespaceRegistration)).To(Succeed())

		counter := 1
		result := testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		for (result.Requeue || result.RequeueAfter > 0) && counter < 5 {
			counter++
			result = testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
			time.Sleep(1 * time.Second)
		}

		// successful deletion
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, namespaceRegistration, 5*time.Second)).To(Succeed())
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(&namespace), &namespace)).To(Succeed())

		// check for namespace being deleted
		Expect(namespace.Status.Phase).To(Equal(corev1.NamespaceTerminating))
	})

	It("should delete namespace with installation without uninstall", func() {
		var err error

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test7")
		Expect(err).ToNot(HaveOccurred())

		namespaceRegistration := state.GetNamespaceRegistration(subjectsync.CUSTOM_NS_PREFIX + "test-namespace-7")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(namespaceRegistration), namespaceRegistration)).To(Succeed())
		Expect(namespaceRegistration.Status.Phase).To(Equal("Completed"))

		// check for namespace being created
		namespace := corev1.Namespace{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: namespaceRegistration.Name}, &namespace)).To(Succeed())

		// add execution to prevent delete
		inst := &lsv1alpha1.Installation{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: namespace.Name, Namespace: namespace.Name},
		}

		metav1.SetMetaDataAnnotation(&inst.ObjectMeta, lsv1alpha1.DeleteWithoutUninstallAnnotation, "true")
		controllerutil.AddFinalizer(inst, lsv1alpha1.LandscaperFinalizer)
		Expect(testenv.Client.Create(ctx, inst)).To(Succeed())

		// delete
		Expect(testenv.Client.Delete(ctx, namespaceRegistration)).To(Succeed())

		// delete installations
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(inst), inst)).To(Succeed())
		Expect(inst.GetDeletionTimestamp().IsZero()).ToNot(BeTrue())
		Expect(helper.HasOperation(inst.ObjectMeta, lsv1alpha1.ReconcileOperation)).ToNot(BeTrue())

		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(inst), inst)).To(Succeed())
		Expect(helper.HasOperation(inst.ObjectMeta, lsv1alpha1.ReconcileOperation)).To(BeTrue())

		controllerutil.RemoveFinalizer(inst, lsv1alpha1.LandscaperFinalizer)
		Expect(testenv.Client.Update(ctx, inst)).To(Succeed())
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(namespaceRegistration))

		// successful deletion
		Expect(testenv.WaitForObjectToBeDeleted(ctx, testenv.Client, namespaceRegistration, 5*time.Second)).To(Succeed())
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(&namespace), &namespace)).To(Succeed())

		// check for namespace being deleted
		Expect(namespace.Status.Phase).To(Equal(corev1.NamespaceTerminating))
	})
})
