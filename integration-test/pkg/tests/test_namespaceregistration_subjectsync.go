// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	cliutil "github.com/gardener/landscapercli/pkg/util"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
	"github.com/gardener/landscaper-service/test/integration/pkg/util"
)

const (
	USER_CLUSTER_ROLE                 = "landscaper-service:namespace-registrator"
	USER_CLUSTER_ROLE_BINDING         = "landscaper-service:namespace-registrator"
	LS_USER_ROLE_IN_NAMESPACE         = "landscaper-service:namespace-registrator"
	LS_USER_ROLE_BINDING_IN_NAMESPACE = "landscaper-service:namespace-registrator"
	USER_ROLE_IN_NAMESPACE            = "landscaper-service:landscaper-user"
	USER_ROLE_BINDING_IN_NAMESPACE    = "landscaper-service:landscaper-user"

	SUBJECT_LIST_NAME = "subjects"
	LS_USER_NAMESPACE = "ls-user"
	CUSTOM_NS_PREFIX  = "cu-"
)

type NamespaceregistrationSubjectSyncRunner struct {
	BaseTestRunner
	resourceClusterAdminClient client.Client
}

func (r *NamespaceregistrationSubjectSyncRunner) Init(
	ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.BaseInit(r.Name(), ctx, config, clusterClients, clusterTargets, testObjects)
}

func (r *NamespaceregistrationSubjectSyncRunner) Name() string {
	return "NamespaceregistrationSubjectSyncRunner"
}

func (r *NamespaceregistrationSubjectSyncRunner) Description() string {
	description := `This test uses an existing landscaper deployment for a tenant and checks the work of the namespaceregistration and subjectsync controller. 
`
	return description
}

func (r *NamespaceregistrationSubjectSyncRunner) String() string {
	return r.Name()
}

func (r *NamespaceregistrationSubjectSyncRunner) Run() error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	namespaceregistrationName := "cu-test-registration"
	user1 := "user1"
	user2 := "user2"

	for _, deployment := range r.testObjects.LandscaperDeployments {
		logger.Info("checking initial setup")
		if err := r.checkInitialSetup(deployment); err != nil {
			return fmt.Errorf("failed checking initial setup: %w", err)
		}
		if err := r.addUserToSubjectListAndCheckChangePropagated(user1); err != nil {
			return fmt.Errorf("failed adding user to subjectlist and check if change propagated: %w", err)
		}
		if err := r.addNamespaceregistrationAndCheckCreation(namespaceregistrationName); err != nil {
			return fmt.Errorf("failed adding namespaceregistration and check setup: %w", err)
		}
		if err := r.addUserToSubjectListAndCheckChangePropagated(user2); err != nil {
			return fmt.Errorf("failed adding user to subjectlist and check if change propagated: %w", err)
		}
		if err := r.deleteUserFromSubjectsAndCheckSync(user2); err != nil {
			return fmt.Errorf("failed deleting user from subjectlist and check if change propagated: %w", err)
		}
		if err := r.deleteNamespaceregistrationAndCheckForNamespaceDeletion(namespaceregistrationName); err != nil {
			return fmt.Errorf("failed deleting namespaceregistration and checking if namespace is deleted: %w", err)
		}
	}
	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) deleteNamespaceregistrationAndCheckForNamespaceDeletion(namespaceregistrationName string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)
	logger.Info("deleting namespaceregistration", "name", namespaceregistrationName)

	namespaceRegistration := &lssv1alpha1.NamespaceRegistration{
		ObjectMeta: v1.ObjectMeta{
			Name:      namespaceregistrationName,
			Namespace: LS_USER_NAMESPACE,
		},
	}
	if err := r.resourceClusterAdminClient.Delete(
		r.ctx,
		namespaceRegistration); err != nil {
		return fmt.Errorf("failed deleting namespaceregistration %q/%q: %w", namespaceRegistration.Namespace, namespaceRegistration.Name, err)
	}

	logger.Info("waiting for namespaceregistration to be deleted", "name", namespaceregistrationName)
	timeout, err := cliutil.CheckAndWaitUntilObjectNotExistAnymore(
		r.clusterClients.TestCluster,
		types.NamespacedName{Name: namespaceregistrationName, Namespace: LS_USER_NAMESPACE}, namespaceRegistration,
		r.config.SleepTime, r.config.MaxRetries*10)

	if err != nil {
		return fmt.Errorf("failed waiting for namespace to be deleted with error: %w", err)
	}
	if timeout {
		return fmt.Errorf("timeout waiting for namespace to be deleted")
	}

	logger.Info("waiting for namespace to be deleted", "name", namespaceregistrationName)
	namespace := &corev1.Namespace{}
	timeout, err = cliutil.CheckAndWaitUntilObjectNotExistAnymore(
		r.clusterClients.TestCluster,
		types.NamespacedName{Name: namespaceregistrationName}, namespace,
		r.config.SleepTime, r.config.MaxRetries*10)

	if err != nil {
		return fmt.Errorf("failed waiting for namespace to be deleted with error: %w", err)
	}
	if timeout {
		return fmt.Errorf("timeout waiting for namespace to be deleted")
	}

	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) deleteUserFromSubjectsAndCheckSync(username string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)
	logger.Info("check subjectlist change synced after deleting user", "name", username)

	subjects := &lssv1alpha1.SubjectList{}
	if err := r.resourceClusterAdminClient.Get(
		r.ctx,
		types.NamespacedName{Name: SUBJECT_LIST_NAME, Namespace: LS_USER_NAMESPACE},
		subjects); err != nil {
		return fmt.Errorf("failed to get subjects: %w", err)
	}

	lengthSubjectsBefore := len(subjects.Spec.Subjects)

	filteredSubjects := []lssv1alpha1.Subject{}

	for _, s := range subjects.Spec.Subjects {
		if s.Name != username {
			filteredSubjects = append(filteredSubjects, s)
		}
	}
	subjects.Spec.Subjects = filteredSubjects

	if err := r.resourceClusterAdminClient.Update(r.ctx, subjects); err != nil {
		return fmt.Errorf("failed updating subjectlist for %q/%q:%w", LS_USER_NAMESPACE, LS_USER_ROLE_BINDING_IN_NAMESPACE, err)
	}

	// check that there is one subject less after update
	subjectsAfterUpdate := &lssv1alpha1.SubjectList{}
	if err := r.resourceClusterAdminClient.Get(
		r.ctx,
		types.NamespacedName{Name: SUBJECT_LIST_NAME, Namespace: LS_USER_NAMESPACE},
		subjectsAfterUpdate); err != nil {
		return fmt.Errorf("failed to get updated subjects: %w", err)
	}
	if len(subjectsAfterUpdate.Spec.Subjects) != lengthSubjectsBefore-1 {
		return fmt.Errorf("deleting %q from subjects has no effect on length of subjects", username)
	}

	if err := r.checkAllNamespacesForSubjectsSynced(subjects); err != nil {
		return fmt.Errorf("failed checking all namespaces for successfull subjects sync: %w", err)
	}

	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) addNamespaceregistrationAndCheckCreation(namespaceregistrationName string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)
	logger.Info("add namespaceregistation and check correct namespace setup", "name", namespaceregistrationName)

	namespaceRegistration := &lssv1alpha1.NamespaceRegistration{
		ObjectMeta: v1.ObjectMeta{
			Name:      namespaceregistrationName,
			Namespace: LS_USER_NAMESPACE,
		},
	}
	if err := r.resourceClusterAdminClient.Create(
		r.ctx,
		namespaceRegistration); err != nil {
		return fmt.Errorf("failed creating namespaceregistration %q/%q: %w", namespaceRegistration.Namespace, namespaceRegistration.Name, err)
	}

	subjects := &lssv1alpha1.SubjectList{}
	if err := r.resourceClusterAdminClient.Get(
		r.ctx,
		types.NamespacedName{Name: SUBJECT_LIST_NAME, Namespace: LS_USER_NAMESPACE},
		subjects); err != nil {
		return fmt.Errorf("failed to get subjects for deployment: %w", err)
	}

	logger.Info("checking created namespace, role and rolebinding")
	timeout, err := cliutil.CheckConditionPeriodically(func() (bool, error) {
		if err := r.resourceClusterAdminClient.Get(
			r.ctx,
			types.NamespacedName{Name: namespaceRegistration.Name, Namespace: namespaceRegistration.Namespace},
			namespaceRegistration); err != nil {
			return false, fmt.Errorf("failed retrieving namespaceregistration %q/%q: %w", namespaceRegistration.Namespace, namespaceRegistration.Name, err)
		}
		if namespaceRegistration.Status.Phase != "Completed" {
			return false, nil
		}

		role := &rbacv1.Role{}
		if err := r.resourceClusterAdminClient.Get(r.ctx, types.NamespacedName{Name: USER_ROLE_IN_NAMESPACE, Namespace: namespaceRegistration.Name}, role); err != nil {
			return false, nil
		}

		rolebinding := &rbacv1.RoleBinding{}
		if err := r.resourceClusterAdminClient.Get(r.ctx, types.NamespacedName{Name: USER_ROLE_BINDING_IN_NAMESPACE, Namespace: namespaceRegistration.Name}, rolebinding); err != nil {
			return false, nil
		}
		if len(rolebinding.Subjects) != len(subjects.Spec.Subjects) {
			return false, nil
		}
		for i := 0; i < len(subjects.Spec.Subjects); i++ {
			if rolebinding.Subjects[i].Kind != subjects.Spec.Subjects[i].Kind ||
				rolebinding.Subjects[i].Name != subjects.Spec.Subjects[i].Name ||
				rolebinding.Subjects[i].Namespace != subjects.Spec.Subjects[i].Namespace {
				return false, nil
			}
		}

		return true, nil

	}, r.config.SleepTime, r.config.MaxRetries)
	if err != nil {
		return fmt.Errorf("failed checking for namespaceregistration creation: %w", err)
	}
	if timeout {
		return fmt.Errorf("timeout waiting for namespaceregistration creation")
	}

	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) addUserToSubjectListAndCheckChangePropagated(username string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)
	logger.Info("check subjectlist change synced after adding user", "name", username)

	subjects := &lssv1alpha1.SubjectList{}
	if err := r.resourceClusterAdminClient.Get(
		r.ctx,
		types.NamespacedName{Name: SUBJECT_LIST_NAME, Namespace: LS_USER_NAMESPACE},
		subjects); err != nil {
		return fmt.Errorf("failed to get subjects for deployment: %w", err)
	}
	subjects.Spec.Subjects = append(subjects.Spec.Subjects, lssv1alpha1.Subject{Kind: "User", Name: username})
	if err := r.resourceClusterAdminClient.Update(r.ctx, subjects); err != nil {
		return fmt.Errorf("failed updating subjectlist for %q/%q:%w", LS_USER_NAMESPACE, LS_USER_ROLE_BINDING_IN_NAMESPACE, err)
	}

	if err := r.checkAllNamespacesForSubjectsSynced(subjects); err != nil {
		return fmt.Errorf("failed checking all namespaces for successfull subjects sync: %w", err)
	}

	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) checkAllNamespacesForSubjectsSynced(subjects *lssv1alpha1.SubjectList) error {
	roleBindingGetter := func(namespaceName types.NamespacedName) ([]rbacv1.Subject, error) {
		rolebinding := &rbacv1.RoleBinding{}
		if err := r.resourceClusterAdminClient.Get(
			r.ctx,
			namespaceName,
			rolebinding); err != nil {
			return nil, fmt.Errorf("failed to get rolebinding %q/%q: %w", namespaceName.Namespace, namespaceName.Name, err)
		}
		return rolebinding.Subjects, nil
	}
	clusterRoleBindingGetter := func(namespaceName types.NamespacedName) ([]rbacv1.Subject, error) {
		clusterRolebinding := &rbacv1.ClusterRoleBinding{}
		if err := r.resourceClusterAdminClient.Get(
			r.ctx,
			namespaceName,
			clusterRolebinding); err != nil {
			return nil, fmt.Errorf("failed to get clusterrolebinding %q/%q: %w", namespaceName.Namespace, namespaceName.Name, err)
		}
		return clusterRolebinding.Subjects, nil
	}

	if err := r.checkRolebindingForSubjectlistSynced(roleBindingGetter, types.NamespacedName{Namespace: LS_USER_NAMESPACE, Name: LS_USER_ROLE_BINDING_IN_NAMESPACE}, subjects); err != nil {
		return fmt.Errorf("failed sycing subjectlist for %q/%q:%w", LS_USER_NAMESPACE, LS_USER_ROLE_BINDING_IN_NAMESPACE, err)
	}

	if err := r.checkRolebindingForSubjectlistSynced(clusterRoleBindingGetter, types.NamespacedName{Name: USER_CLUSTER_ROLE_BINDING}, subjects); err != nil {
		return fmt.Errorf("failed sycing subjectlist for clusterrolebinding %q:%w", USER_CLUSTER_ROLE_BINDING, err)
	}

	namespaceList := &corev1.NamespaceList{}
	if err := r.resourceClusterAdminClient.List(r.ctx, namespaceList); err != nil {
		return fmt.Errorf("failed retrieving namespace list:%w", err)
	}
	for _, v := range namespaceList.Items {
		if strings.HasPrefix(v.Name, "cu-") {
			if err := r.checkRolebindingForSubjectlistSynced(roleBindingGetter, types.NamespacedName{Namespace: v.Name, Name: USER_ROLE_BINDING_IN_NAMESPACE}, subjects); err != nil {
				return fmt.Errorf("failed sycing subjectlist for %q/%q:%w", v.Name, USER_ROLE_BINDING_IN_NAMESPACE, err)
			}
		}
	}
	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) checkRolebindingForSubjectlistSynced(getSubjects func(types.NamespacedName) ([]rbacv1.Subject, error), namespaceName types.NamespacedName, subjects *lssv1alpha1.SubjectList) error {
	timeout, err := cliutil.CheckConditionPeriodically(func() (bool, error) {
		subjectInBinding, err := getSubjects(namespaceName)
		if err != nil {
			return false, fmt.Errorf("failed to get rolebinding %q/%q: %w", namespaceName.Namespace, namespaceName.Name, err)
		}
		if len(subjectInBinding) != len(subjects.Spec.Subjects) {
			return false, nil
		}
		for i := 0; i < len(subjects.Spec.Subjects); i++ {
			if subjectInBinding[i].Kind != subjects.Spec.Subjects[i].Kind ||
				subjectInBinding[i].Name != subjects.Spec.Subjects[i].Name ||
				subjectInBinding[i].Namespace != subjects.Spec.Subjects[i].Namespace {
				return false, nil
			}
		}
		return true, nil
	}, r.config.SleepTime, r.config.MaxRetries)
	if err != nil {
		return fmt.Errorf("failed checking for subjectlist sync: %w", err)
	}
	if timeout {
		return fmt.Errorf("timeout waiting for subjectlist sync")
	}
	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) checkInitialSetup(deployment *lssv1alpha1.LandscaperDeployment) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)
	logger.Info("check initial setup for deployment", "name", deployment.Name)

	// get admin kubeconfig for resource-shoot cluster from landscaperdeployment.status.instanceRef - Instance.Status.AdminKubeconfig
	logger.Info("build kube client from config in instance.Status")
	if deployment.Status.InstanceRef.Name == "" || deployment.Status.InstanceRef.Namespace == "" {
		return fmt.Errorf("deployment %q instance ref empty", deployment.Name)
	}

	instance := &lssv1alpha1.Instance{}
	if err := r.clusterClients.TestCluster.Get(
		r.ctx,
		types.NamespacedName{Name: deployment.Status.InstanceRef.Name, Namespace: deployment.Status.InstanceRef.Namespace},
		instance); err != nil {
		return fmt.Errorf("failed to get instance for deployment %q: %w", deployment.Name, err)
	}

	if instance.Status.AdminKubeconfig == "" {
		return fmt.Errorf("instance %q for deployment %q missing AdminKubeconfig", instance.Name, deployment.Name)
	}

	kubeconfig, err := base64.StdEncoding.DecodeString(instance.Status.AdminKubeconfig)
	if err != nil {
		return fmt.Errorf("failed to decode admin kubeconfig of instance %q/%q: %w", instance.Namespace, instance.Name, err)
	}

	//build client
	kubeClient, err := util.BuildKubeClient(string(kubeconfig), test.Scheme())
	if err != nil {
		return fmt.Errorf("failed building KubeClient for instance %q/%q: %w", instance.Namespace, instance.Name, err)
	}
	r.resourceClusterAdminClient = kubeClient

	logger.Info("check initial namespace existance", "name", LS_USER_NAMESPACE)
	namespace := &corev1.Namespace{}
	err = r.resourceClusterAdminClient.Get(r.ctx, types.NamespacedName{Name: LS_USER_NAMESPACE}, namespace)
	if err != nil {
		return fmt.Errorf("failed retrieving namespace: %w", err)
	}

	logger.Info("check role existance", "name", LS_USER_ROLE_IN_NAMESPACE)
	role := &rbacv1.Role{}
	err = r.resourceClusterAdminClient.Get(r.ctx, types.NamespacedName{Name: LS_USER_ROLE_IN_NAMESPACE, Namespace: namespace.Name}, role)
	if err != nil {
		return fmt.Errorf("failed retrieving role: %w", err)
	}

	logger.Info("check role binding existance and being empty", "name", LS_USER_ROLE_BINDING_IN_NAMESPACE)
	rolebinding := &rbacv1.RoleBinding{}
	err = r.resourceClusterAdminClient.Get(r.ctx, types.NamespacedName{Name: LS_USER_ROLE_BINDING_IN_NAMESPACE, Namespace: namespace.Name}, rolebinding)
	if err != nil {
		return fmt.Errorf("failed retrieving role: %w", err)
	}

	if len(rolebinding.Subjects) != 0 {
		return fmt.Errorf("initial rolebinding should be empty but contains %q subjects", len(rolebinding.Subjects))
	}

	logger.Info("check cluster role existance", "name", USER_CLUSTER_ROLE)
	clusterRole := &rbacv1.ClusterRole{}
	err = r.resourceClusterAdminClient.Get(r.ctx, types.NamespacedName{Name: USER_CLUSTER_ROLE}, clusterRole)
	if err != nil {
		return fmt.Errorf("failed retrieving clusterrole: %w", err)
	}

	logger.Info("check clusterrole binding existance and being empty", "name", USER_CLUSTER_ROLE_BINDING)
	clusterRolebinding := &rbacv1.ClusterRoleBinding{}
	err = r.resourceClusterAdminClient.Get(r.ctx, types.NamespacedName{Name: USER_CLUSTER_ROLE_BINDING}, clusterRolebinding)
	if err != nil {
		return fmt.Errorf("failed retrieving role: %w", err)
	}

	if len(clusterRolebinding.Subjects) != 0 {
		return fmt.Errorf("initial rolebinding should be empty but contains %q subjects", len(clusterRolebinding.Subjects))
	}

	return nil
}
