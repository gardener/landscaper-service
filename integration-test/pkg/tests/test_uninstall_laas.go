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

	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type UninstallLAASTestRunner struct {
	ctx         context.Context
	log         logr.Logger
	kclient     client.Client
	config      *test.TestConfig
	target      *lsv1alpha1.Target
	testObjects *test.SharedTestObjects
}

func (r *UninstallLAASTestRunner) Init(
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

func (r *UninstallLAASTestRunner) Name() string {
	return "UninstallLAAS"
}

func (r *UninstallLAASTestRunner) Run() error {
	r.log.Info("deleting default service target config")
	if err := r.deleteServiceTargetConfig(); err != nil {
		return err
	}

	r.log.Info("deleting laas installation")
	if err := r.deleteInstallation(); err != nil {
		return err
	}

	return nil
}

func (r *UninstallLAASTestRunner) deleteServiceTargetConfig() error {
	serviceTargetConfig := r.testObjects.ServiceTargetConfigs[types.NamespacedName{Name: "default-target", Namespace: r.config.LaasNamespace}.String()]

	if err := r.kclient.Delete(r.ctx, serviceTargetConfig); err != nil {
		return fmt.Errorf("failed to delete service target config: %w", err)
	}

	return nil
}

func (r *UninstallLAASTestRunner) deleteInstallation() error {
	installation := r.testObjects.Installations[types.NamespacedName{Name: "laas", Namespace: r.config.LaasNamespace}.String()]

	if err := r.kclient.Delete(r.ctx, installation); err != nil {
		return fmt.Errorf("failed to delete laas installation: %w", err)
	}

	timeout, err := cliutil.CheckAndWaitUntilObjectNotExistAnymore(
		r.kclient,
		types.NamespacedName{Name: installation.Name, Namespace: installation.Namespace},
		installation,
		r.config.SleepTime,
		r.config.MaxRetries)

	if err != nil {
		return fmt.Errorf("failed to wait for laas installation to be deleted: %w", err)
	}
	if timeout {
		return fmt.Errorf("waiting for laas installation to be deleted timed out")
	}

	return nil
}
