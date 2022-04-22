// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	cliutil "github.com/gardener/landscapercli/pkg/util"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type DeleteDeploymentRunner struct {
	ctx         context.Context
	log         logr.Logger
	kclient     client.Client
	config      *test.TestConfig
	target      *lsv1alpha1.Target
	testObjects *test.SharedTestObjects
}

func (r *DeleteDeploymentRunner) Init(
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

func (r *DeleteDeploymentRunner) Name() string {
	return "DeleteDeployment"
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
	r.log.Info("deleting deployment", "name", deployment.Name)

	if err := r.kclient.Delete(r.ctx, deployment); err != nil {
		return fmt.Errorf("failed to delete deployment %q: %w", deployment.Name, err)
	}

	r.log.Info("waiting for deployment to be deleted", "name", deployment.Name)
	timeout, err := cliutil.CheckAndWaitUntilObjectNotExistAnymore(
		r.kclient,
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
