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
	logger.Info("checking initial setup")

	for _, deployment := range r.testObjects.LandscaperDeployments {
		if err := r.checkInitialSetup(deployment); err != nil {
			return fmt.Errorf("failed checking initial setup: %w", err)
		}
		if err := r.addUserToSubjectListAndCheckChangePropagated(deployment, "user1"); err != nil {
			return fmt.Errorf("failed adding user to subjectlist and check if change propagated: %w", err)
		}
		if err := r.addNamespaceregistrationAndCheckCreation("cu-test-registration"); err != nil {
			return fmt.Errorf("failed adding namespaceregistration and check setup: %w", err)
		}
		if err := r.addUserToSubjectListAndCheckChangePropagated(deployment, "user2"); err != nil {
			return fmt.Errorf("failed adding user to subjectlist and check if change propagated: %w", err)
		}
	}
	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) addNamespaceregistrationAndCheckCreation(namespaceregistrationName string) error {
	namespaceRegistration := &lssv1alpha1.NamespaceRegistration{
		ObjectMeta: v1.ObjectMeta{
			Name:      namespaceregistrationName,
			Namespace: "ls-user",
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
		types.NamespacedName{Name: "subjects", Namespace: "ls-user"},
		subjects); err != nil {
		return fmt.Errorf("failed to get subjects for deployment: %w", err)
	}

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
		if err := r.resourceClusterAdminClient.Get(r.BaseTestRunner.ctx, types.NamespacedName{Name: "user-role", Namespace: namespaceRegistration.Name}, role); err != nil {
			return false, nil
		}

		rolebinding := &rbacv1.RoleBinding{}
		if err := r.resourceClusterAdminClient.Get(r.BaseTestRunner.ctx, types.NamespacedName{Name: "user-role-binding", Namespace: namespaceRegistration.Name}, rolebinding); err != nil {
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

func (r *NamespaceregistrationSubjectSyncRunner) addUserToSubjectListAndCheckChangePropagated(deployment *lssv1alpha1.LandscaperDeployment, username string) error {
	// add user to subjectlist
	subjects := &lssv1alpha1.SubjectList{}
	if err := r.resourceClusterAdminClient.Get(
		r.ctx,
		types.NamespacedName{Name: "subjects", Namespace: "ls-user"},
		subjects); err != nil {
		return fmt.Errorf("failed to get subjects for deployment: %w", err)
	}
	subjects.Spec.Subjects = append(subjects.Spec.Subjects, lssv1alpha1.Subject{Kind: "User", Name: username})
	r.resourceClusterAdminClient.Update(r.ctx, subjects)

	if err := r.checkRolebindingForSubjectlistSynced(types.NamespacedName{Namespace: "ls-user", Name: "ls-user-role-binding"}, subjects); err != nil {
		return fmt.Errorf("failed sycing subjectlist for ls-user/ls-user-role-binding:%w", err)
	}

	namespaceList := &corev1.NamespaceList{}
	if err := r.resourceClusterAdminClient.List(r.ctx, namespaceList); err != nil {
		return fmt.Errorf("failed retrieving namespace list:%w", err)
	}
	for _, v := range namespaceList.Items {
		if strings.HasPrefix(v.Name, "cu-") {
			if err := r.checkRolebindingForSubjectlistSynced(types.NamespacedName{Namespace: v.Name, Name: "user-role-binding"}, subjects); err != nil {
				return fmt.Errorf("failed sycing subjectlist for %q/user-role-binding:%w", v.Name, err)
			}
		}
	}
	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) checkRolebindingForSubjectlistSynced(namespaceName types.NamespacedName, subjects *lssv1alpha1.SubjectList) error {
	timeout, err := cliutil.CheckConditionPeriodically(func() (bool, error) {
		rolebinding := &rbacv1.RoleBinding{}
		if err := r.resourceClusterAdminClient.Get(
			r.ctx,
			namespaceName,
			rolebinding); err != nil {
			return false, fmt.Errorf("failed to get rolebinding %q/%q: %w", namespaceName.Namespace, namespaceName.Name, err)
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
		return fmt.Errorf("failed checking for subjectlist sync: %w", err)
	}
	if timeout {
		return fmt.Errorf("timeout waiting for subjectlist sync")
	}
	return nil
}

func (r *NamespaceregistrationSubjectSyncRunner) checkInitialSetup(deployment *lssv1alpha1.LandscaperDeployment) error {
	// get admin kubeconfig for resource-shoot cluster from landscaperdeployment.status.instanceRef - Instance.Status.AdminKubeconfig
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

	//check for namespace ls-user, role, rolebinding, subjectsynclist
	namespace := &corev1.Namespace{}
	err = r.resourceClusterAdminClient.Get(r.BaseTestRunner.ctx, types.NamespacedName{Name: "ls-user"}, namespace)
	if err != nil {
		return fmt.Errorf("failed retrieving namespace: %w", err)
	}

	role := &rbacv1.Role{}
	err = r.resourceClusterAdminClient.Get(r.BaseTestRunner.ctx, types.NamespacedName{Name: "ls-user-role", Namespace: namespace.Name}, role)
	if err != nil {
		return fmt.Errorf("failed retrieving role: %w", err)
	}

	rolebinding := &rbacv1.RoleBinding{}
	err = r.resourceClusterAdminClient.Get(r.BaseTestRunner.ctx, types.NamespacedName{Name: "ls-user-role-binding", Namespace: namespace.Name}, rolebinding)
	if err != nil {
		return fmt.Errorf("failed retrieving role: %w", err)
	}

	if len(rolebinding.Subjects) != 0 {
		return fmt.Errorf("initial rolebinding should be empty but contains %q subjects", len(rolebinding.Subjects))
	}

	return nil
}
