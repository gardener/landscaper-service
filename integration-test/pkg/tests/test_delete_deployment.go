// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"

	cliutil "github.com/gardener/landscapercli/pkg/util"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type DeleteDeploymentRunner struct {
	BaseTestRunner
}

func (r *DeleteDeploymentRunner) Init(
	ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.BaseInit(r.Name(), ctx, config, clusterClients, clusterTargets, testObjects)
}

func (r *DeleteDeploymentRunner) Name() string {
	return "DeleteDeployment"
}

func (r *DeleteDeploymentRunner) Description() string {
	description := `This test deletes a Landscaper deployment for a tenant.
The test waits until the deployment crd has been deleted, which also means, that the
corresponding installation has been deleted.
The test fails when the deployment crd hasn't been deleted before the timeout has expired.'
`
	return description
}

func (r *DeleteDeploymentRunner) String() string {
	return r.Name()
}

func (r *DeleteDeploymentRunner) Run() error {
	for _, deployment := range r.testObjects.LandscaperDeployments {
		if err := r.deleteDeployment(deployment); err != nil {
			return err
		}
	}
	return nil
}

func (r *DeleteDeploymentRunner) deleteDeployment(deployment *lssv1alpha1.LandscaperDeployment) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)
	logger.Info("deleting deployment", "name", deployment.Name)

	if err := r.clusterClients.TestCluster.Delete(r.ctx, deployment); err != nil {
		return fmt.Errorf("failed to delete deployment %q: %w", deployment.Name, err)
	}

	logger.Info("waiting for deployment to be deleted", "name", deployment.Name)
	timeout, err := cliutil.CheckAndWaitUntilObjectNotExistAnymore(
		r.clusterClients.TestCluster,
		types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, deployment,
		r.config.SleepTime, r.config.MaxRetries*5)

	if err != nil {
		return fmt.Errorf("failed to wait for deployment %q to be deleted: %w", deployment.Name, err)
	}

	if timeout {
		return fmt.Errorf("waiting for deployment %q to be deleted timed out", deployment.Name)
	}

	return nil
}
