// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package subjectsync_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"
	"github.com/gardener/landscaper-service/pkg/operation"
	testutils "github.com/gardener/landscaper-service/test/utils"
	"github.com/gardener/landscaper-service/test/utils/envtest"
	kutil "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
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
		ctrl = subjectsync.NewTestActuator(*op, logging.Discard())
	})

	AfterEach(func() {
		defer ctx.Done()
		if state != nil {
			Expect(testenv.CleanupResources(ctx, state)).ToNot(HaveOccurred())
		}
	})

	SetupNamespacesRolesAndBindings := func() (*corev1.Namespace, *corev1.Namespace) {
		// set up namespaces
		lsUserNamespace := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "ls-user-"}}
		userNamespace := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "user1-"}}
		Expect(testenv.Client.Create(ctx, &lsUserNamespace)).To(Succeed())
		Expect(testenv.Client.Create(ctx, &userNamespace)).To(Succeed())

		// setup up roles and role bindings
		lsUserRole := rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subjectsync.LS_USER_ROLE_IN_NAMESPACE,
				Namespace: lsUserNamespace.Name,
			},
		}
		lsUserRoleBinding := rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE,
				Namespace: lsUserNamespace.Name,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     lsUserRole.Name,
			},
		}
		Expect(testenv.Client.Create(ctx, &lsUserRole)).To(Succeed())
		Expect(testenv.Client.Create(ctx, &lsUserRoleBinding)).To(Succeed())

		userRole := rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subjectsync.USER_ROLE_IN_NAMESPACE,
				Namespace: userNamespace.Name,
			},
		}
		userRoleBinding := rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subjectsync.USER_ROLE_BINDING_IN_NAMESPACE,
				Namespace: userNamespace.Name,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     userRole.Name,
			},
		}
		Expect(testenv.Client.Create(ctx, &userRole)).To(Succeed())
		Expect(testenv.Client.Create(ctx, &userRoleBinding)).To(Succeed())

		return &lsUserNamespace, &userNamespace
	}

	It("should sync role binding subjects to subjectlist", func() {
		var err error

		lsUserNamespace, userNamespace := SetupNamespacesRolesAndBindings()

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())

		subjectlist := state.GetSubjectList("subjectlist")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		updatedLsUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())
		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedLsUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedLsUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Kind).To(Equal("Group"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Name).To(Equal("testgroup"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccount"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		updatedUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedUserRoleBinding.Subjects[1].Kind).To(Equal("Group"))
		Expect(updatedUserRoleBinding.Subjects[1].Name).To(Equal("testgroup"))
		Expect(updatedUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccount"))
		Expect(updatedUserRoleBinding.Subjects[2].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))
	})

	It("should skip unknown/erroneous subjects", func() {
		var err error

		lsUserNamespace, userNamespace := SetupNamespacesRolesAndBindings()

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test2")
		Expect(err).ToNot(HaveOccurred())

		subjectlist := state.GetSubjectList("subjectlist")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		updatedLsUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())
		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(0))

		updatedUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(0))
	})

	It("should remove subjects if removed from the subject list", func() {
		var err error

		lsUserNamespace, userNamespace := SetupNamespacesRolesAndBindings()

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())

		subjectlist := state.GetSubjectList("subjectlist")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		updatedLsUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())
		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedLsUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedLsUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Kind).To(Equal("Group"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Name).To(Equal("testgroup"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccount"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		updatedUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedUserRoleBinding.Subjects[1].Kind).To(Equal("Group"))
		Expect(updatedUserRoleBinding.Subjects[1].Name).To(Equal("testgroup"))
		Expect(updatedUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccount"))
		Expect(updatedUserRoleBinding.Subjects[2].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		//delete group
		subjectlist.Spec.Subjects = append(subjectlist.Spec.Subjects[:1], subjectlist.Spec.Subjects[2:]...)
		Expect(testenv.Client.Update(ctx, subjectlist)).To(Succeed())
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())
		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(2))
		Expect(updatedLsUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedLsUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Kind).To(Equal("ServiceAccount"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Name).To(Equal("testserviceaccount"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(2))
		Expect(updatedUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedUserRoleBinding.Subjects[1].Kind).To(Equal("ServiceAccount"))
		Expect(updatedUserRoleBinding.Subjects[1].Name).To(Equal("testserviceaccount"))
		Expect(updatedUserRoleBinding.Subjects[1].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		//delete user
		subjectlist.Spec.Subjects = subjectlist.Spec.Subjects[1:]
		Expect(testenv.Client.Update(ctx, subjectlist)).To(Succeed())
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())
		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(1))
		Expect(updatedLsUserRoleBinding.Subjects[0].Kind).To(Equal("ServiceAccount"))
		Expect(updatedLsUserRoleBinding.Subjects[0].Name).To(Equal("testserviceaccount"))
		Expect(updatedLsUserRoleBinding.Subjects[0].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(1))
		Expect(updatedUserRoleBinding.Subjects[0].Kind).To(Equal("ServiceAccount"))
		Expect(updatedUserRoleBinding.Subjects[0].Name).To(Equal("testserviceaccount"))
		Expect(updatedUserRoleBinding.Subjects[0].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		//delete service account
		subjectlist.Spec.Subjects = []v1alpha1.Subject{}
		Expect(testenv.Client.Update(ctx, subjectlist)).To(Succeed())
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())
		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(0))

		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(0))
	})

	It("should update entries on modify", func() {
		var err error

		lsUserNamespace, userNamespace := SetupNamespacesRolesAndBindings()

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())

		subjectlist := state.GetSubjectList("subjectlist")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		updatedLsUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())
		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedLsUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedLsUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Kind).To(Equal("Group"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Name).To(Equal("testgroup"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccount"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		updatedUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedUserRoleBinding.Subjects[1].Kind).To(Equal("Group"))
		Expect(updatedUserRoleBinding.Subjects[1].Name).To(Equal("testgroup"))
		Expect(updatedUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccount"))
		Expect(updatedUserRoleBinding.Subjects[2].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		//update values
		subjectlist.Spec.Subjects[0].Kind = "Group"
		subjectlist.Spec.Subjects[0].Name = "testuserMODIFIEDToGroup"
		subjectlist.Spec.Subjects[1].Kind = "User"
		subjectlist.Spec.Subjects[1].Name = "testgroupMODIFIEDToUser"
		subjectlist.Spec.Subjects[2].Name = "testserviceaccountmodified"
		subjectlist.Spec.Subjects[2].Namespace = "ls-user-mod"

		Expect(testenv.Client.Update(ctx, subjectlist)).To(Succeed())

		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())
		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedLsUserRoleBinding.Subjects[0].Kind).To(Equal("Group"))
		Expect(updatedLsUserRoleBinding.Subjects[0].Name).To(Equal("testuserMODIFIEDToGroup"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Kind).To(Equal("User"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Name).To(Equal("testgroupMODIFIEDToUser"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccountmodified"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Namespace).To(Equal("ls-user-mod"))

		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedUserRoleBinding.Subjects[0].Kind).To(Equal("Group"))
		Expect(updatedUserRoleBinding.Subjects[0].Name).To(Equal("testuserMODIFIEDToGroup"))
		Expect(updatedUserRoleBinding.Subjects[1].Kind).To(Equal("User"))
		Expect(updatedUserRoleBinding.Subjects[1].Name).To(Equal("testgroupMODIFIEDToUser"))
		Expect(updatedUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccountmodified"))
		Expect(updatedUserRoleBinding.Subjects[2].Namespace).To(Equal("ls-user-mod"))

	})

	It("should empty role bindings subjects if subjectlist is emptied", func() {
		var err error

		lsUserNamespace, userNamespace := SetupNamespacesRolesAndBindings()

		state, err = testenv.InitResources(ctx, "./testdata/reconcile/test1")
		Expect(err).ToNot(HaveOccurred())

		subjectlist := state.GetSubjectList("subjectlist")
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		updatedLsUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())
		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedLsUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedLsUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Kind).To(Equal("Group"))
		Expect(updatedLsUserRoleBinding.Subjects[1].Name).To(Equal("testgroup"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccount"))
		Expect(updatedLsUserRoleBinding.Subjects[2].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		updatedUserRoleBinding := rbacv1.RoleBinding{}
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(3))
		Expect(updatedUserRoleBinding.Subjects[0].Kind).To(Equal("User"))
		Expect(updatedUserRoleBinding.Subjects[0].Name).To(Equal("testuser"))
		Expect(updatedUserRoleBinding.Subjects[1].Kind).To(Equal("Group"))
		Expect(updatedUserRoleBinding.Subjects[1].Name).To(Equal("testgroup"))
		Expect(updatedUserRoleBinding.Subjects[2].Kind).To(Equal("ServiceAccount"))
		Expect(updatedUserRoleBinding.Subjects[2].Name).To(Equal("testserviceaccount"))
		Expect(updatedUserRoleBinding.Subjects[2].Namespace).To(Equal(subjectsync.LS_USER_NAMESPACE))

		subjectlist.Spec.Subjects = []v1alpha1.Subject{}
		Expect(testenv.Client.Update(ctx, subjectlist)).To(Succeed())
		//reconcile for finalizer
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())
		//reconcile for actual run
		testutils.ShouldReconcile(ctx, ctrl, testutils.RequestFromObject(subjectlist))
		Expect(testenv.Client.Get(ctx, kutil.ObjectKeyFromObject(subjectlist), subjectlist)).To(Succeed())

		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: lsUserNamespace.Name}, &updatedLsUserRoleBinding)).To(Succeed())

		Expect(len(updatedLsUserRoleBinding.Subjects)).To(Equal(0))
		Expect(testenv.Client.Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: userNamespace.Name}, &updatedUserRoleBinding)).To(Succeed())
		Expect(len(updatedUserRoleBinding.Subjects)).To(Equal(0))
	})
})
