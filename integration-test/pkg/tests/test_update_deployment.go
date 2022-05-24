// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	cliutil "github.com/gardener/landscapercli/pkg/util"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type UpdateDeploymentRunner struct {
	ctx         context.Context
	log         logr.Logger
	kclient     client.Client
	config      *test.TestConfig
	target      *lsv1alpha1.Target
	testObjects *test.SharedTestObjects
}

func (r *UpdateDeploymentRunner) Init(
	ctx context.Context, log logr.Logger,
	kclient client.Client, config *test.TestConfig,
	target *lsv1alpha1.Target, testObjects *test.SharedTestObjects) {
	r.ctx = ctx
	r.log = log.WithName(r.Name())
	r.kclient = kclient
	r.config = config
	r.target = target
	r.testObjects = testObjects
}

func (r *UpdateDeploymentRunner) Name() string {
	return "UpdateDeployment"
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
	r.log.Info("updating deployment", "name", deployment.Name)
	deployment.Spec.LandscaperConfiguration.Deployers = append(deployment.Spec.LandscaperConfiguration.Deployers, "container")

	if err := r.kclient.Update(r.ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment %q: %w", deployment.Name, err)
	}

	instance := &lssv1alpha1.Instance{}
	if err := r.kclient.Get(
		r.ctx,
		types.NamespacedName{Name: deployment.Status.InstanceRef.Name, Namespace: deployment.Status.InstanceRef.Namespace},
		instance); err != nil {
		return fmt.Errorf("failed to retrieve instance for deployment %q: %w", deployment.Name, err)
	}

	// waiting for a state change, because installations are already succeeded
	time.Sleep(10 * time.Second)

	r.log.Info("waiting for installation being succeeded")

	timeout, err := cliutil.CheckAndWaitUntilLandscaperInstallationSucceeded(
		r.kclient,
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
