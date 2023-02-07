// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
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
