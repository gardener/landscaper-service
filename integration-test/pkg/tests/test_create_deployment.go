// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	cliutil "github.com/gardener/landscapercli/pkg/util"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

const (
	// LandscaperDeploymentInstallationTimeout is the timeout at which the landscaper deployment installation is considered to be failed.
	LandscaperDeploymentInstallationTimeout = 35 * time.Minute
)

type CreateDeploymentRunner struct {
	BaseTestRunner
}

func (r *CreateDeploymentRunner) Init(
	ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.BaseInit(r.Name(), ctx, config, clusterClients, clusterTargets, testObjects)
}

func (r *CreateDeploymentRunner) Name() string {
	return "CreateDeployment"
}

func (r *CreateDeploymentRunner) Description() string {
	description := `This test creates a Landscaper deployment for a tenant and waits until
the corresponding installation has succeeded. If the installation is in phase succeeded before the test timeout has expired,
the test succeeds, otherwise the test has failed.
`
	return description
}

func (r *CreateDeploymentRunner) String() string {
	return r.Name()
}

func (r *CreateDeploymentRunner) Run() error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)
	logger.Info("creating landscaper deployment")
	if err := r.createDeployment(); err != nil {
		return err
	}
	return nil
}

func (r *CreateDeploymentRunner) createDeployment() error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	deployment := &lssv1alpha1.LandscaperDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: r.config.TestNamespace,
		},
		Spec: lssv1alpha1.LandscaperDeploymentSpec{
			TenantId: createTenantId(),
			Purpose:  "integration-test",
			LandscaperConfiguration: lssv1alpha1.LandscaperConfiguration{
				Deployers: []string{
					"helm",
					"manifest",
				},
			},
			HighAvailabilityConfig: &lssv1alpha1.HighAvailabilityConfig{
				ControlPlaneFailureTolerance: "node",
			},
		},
	}

	if len(r.config.TestPurpose) > 0 {
		deployment.Name = fmt.Sprintf("test-%s", r.config.TestPurpose)
	}

	if err := r.clusterClients.TestCluster.Create(r.ctx, deployment); err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	logger.Info("waiting for instance being created")

	timeout, err := cliutil.CheckConditionPeriodically(func() (bool, error) {
		if err := r.clusterClients.TestCluster.Get(r.ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, deployment); err != nil {
			return false, err
		}

		if deployment.Status.InstanceRef != nil {
			return true, nil
		}

		return false, nil
	}, r.config.SleepTime, r.config.MaxRetries)

	if err != nil {
		return fmt.Errorf("failed to wait for instance being created: %w", err)
	}
	if timeout {
		return fmt.Errorf("timeout while wating for instance being created")
	}

	instance := &lssv1alpha1.Instance{}

	logger.Info("waiting for installation being created")

	timeout, err = cliutil.CheckConditionPeriodically(func() (bool, error) {
		if err := r.clusterClients.TestCluster.Get(
			r.ctx,
			types.NamespacedName{Name: deployment.Status.InstanceRef.Name, Namespace: deployment.Status.InstanceRef.Namespace},
			instance); err != nil {

			return false, err
		}

		if instance.Status.InstallationRef != nil {
			return true, nil
		}

		return false, nil
	}, r.config.SleepTime, r.config.MaxRetries)

	if err != nil {
		return fmt.Errorf("failed to wait for installation being created: %w", err)
	}
	if timeout {
		return fmt.Errorf("timeout while wating for installation being created")
	}

	logger.Info("waiting for installation being succeeded")
	maxRetries := int(LandscaperDeploymentInstallationTimeout.Seconds() / r.config.SleepTime.Seconds())

	timeout, err = cliutil.CheckAndWaitUntilLandscaperInstallationSucceeded(
		r.clusterClients.TestCluster,
		types.NamespacedName{Name: instance.Status.InstallationRef.Name, Namespace: instance.Status.InstallationRef.Namespace},
		r.config.SleepTime,
		maxRetries)

	if err != nil || timeout {
		installation := &lsv1alpha1.Installation{}
		if err := r.clusterClients.TestCluster.Get(r.ctx, types.NamespacedName{Name: instance.Status.InstallationRef.Name, Namespace: instance.Status.InstallationRef.Namespace}, installation); err == nil {
			logger.Error(fmt.Errorf("installation failed"), "installation", "last error", installation.Status.LastError)
		}
	}

	if err != nil {
		return fmt.Errorf("installation for landscaper deployment %s failed: %w", deployment.Name, err)
	}
	if timeout {
		return fmt.Errorf("waiting for installation of landscaper deployment %s timed out", deployment.Name)
	}

	r.testObjects.LandscaperDeployments[types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}.String()] = deployment

	return nil
}

func createTenantId() string {
	low := 10000
	high := 99999
	return fmt.Sprintf("tc-%d", low+rand.Intn(high-low))
}
