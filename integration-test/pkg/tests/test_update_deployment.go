// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"

	lssutils "github.com/gardener/landscaper-service/pkg/utils"

	"k8s.io/apimachinery/pkg/types"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	cliutil "github.com/gardener/landscapercli/pkg/util"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type UpdateDeploymentRunner struct {
	BaseTestRunner
}

func (r *UpdateDeploymentRunner) Init(
	ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.BaseInit(r.Name(), ctx, config, clusterClients, clusterTargets, testObjects)
}

func (r *UpdateDeploymentRunner) Name() string {
	return "UpdateDeployment"
}

func (r *UpdateDeploymentRunner) Description() string {
	description := `This test updates the specification for an existing tenant Landscaper deployment.
The test succeeds when the corresponding installation is in phase succeeded before the timeout expires.
Otherwise the test fails.
`
	return description
}

func (r *UpdateDeploymentRunner) String() string {
	return r.Name()
}

func (r *UpdateDeploymentRunner) Run() error {
	for _, deployment := range r.testObjects.LandscaperDeployments {
		if err := r.updateDeployment(deployment); err != nil {
			return err
		}
	}

	return nil
}

func (r *UpdateDeploymentRunner) updateDeployment(deployment *lssv1alpha1.LandscaperDeployment) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("updating deployment", "name", deployment.Name)
	if err := r.clusterClients.TestCluster.Get(r.ctx, client.ObjectKeyFromObject(deployment), deployment); err != nil {
		return fmt.Errorf("failed to get deployment %q: %w", deployment.Name, err)
	}

	lssutils.SetOperationAnnotation(deployment, lssv1alpha1.LandscaperServiceOperationIgnore)
	deployment.Spec.LandscaperConfiguration.Deployers = append(deployment.Spec.LandscaperConfiguration.Deployers, "container")

	if err := r.clusterClients.TestCluster.Update(r.ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment %q: %w", deployment.Name, err)
	}

	// wait for the controller to reconcile the instance
	time.Sleep(10 * time.Second)

	instance := &lssv1alpha1.Instance{}
	if err := r.clusterClients.TestCluster.Get(
		r.ctx,
		types.NamespacedName{Name: deployment.Status.InstanceRef.Name, Namespace: deployment.Status.InstanceRef.Namespace},
		instance); err != nil {
		return fmt.Errorf("failed to retrieve instance for deployment %q: %w", deployment.Name, err)
	}

	if !lssutils.HasOperationAnnotation(instance, lssv1alpha1.LandscaperServiceOperationIgnore) {
		return fmt.Errorf("ignore operation annotation of deployment %s was not inherited to instance", deployment.Name)
	}

	installation := &lsv1alpha1.Installation{}
	if err := r.clusterClients.TestCluster.Get(
		r.ctx,
		types.NamespacedName{Name: instance.Status.InstallationRef.Name, Namespace: instance.Status.InstallationRef.Namespace}, installation); err != nil {
		return fmt.Errorf("could not get installation for landscaper deployment %s: %w", deployment.Name, err)
	}

	if installation.Status.InstallationPhase != lsv1alpha1.InstallationPhaseSucceeded {
		return fmt.Errorf("ignore annotation of deployment %s was ignored", deployment.Name)
	}

	if err := r.clusterClients.TestCluster.Get(r.ctx, client.ObjectKeyFromObject(deployment), deployment); err != nil {
		return fmt.Errorf("failed to get deployment %q: %w", deployment.Name, err)
	}
	lssutils.RemoveOperationAnnotation(deployment)
	if err := r.clusterClients.TestCluster.Update(r.ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment %q: %w", deployment.Name, err)
	}

	// waiting for a state change, because installations are already succeeded
	time.Sleep(10 * time.Second)

	logger.Info("waiting for installation being succeeded")

	timeout, err := cliutil.CheckAndWaitUntilLandscaperInstallationSucceeded(
		r.clusterClients.TestCluster,
		types.NamespacedName{Name: instance.Status.InstallationRef.Name, Namespace: instance.Status.InstallationRef.Namespace},
		r.config.SleepTime,
		r.config.MaxRetries*10)

	if err != nil {
		return fmt.Errorf("installation for landscaper deployment %s failed: %w", deployment.Name, err)
	}
	if timeout {
		return fmt.Errorf("waiting for installation of landscaper deployment %s timed out", deployment.Name)
	}

	r.testObjects.LandscaperDeployments[types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}.String()] = deployment
	return nil
}
