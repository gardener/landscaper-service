// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	cliutil "github.com/gardener/landscapercli/pkg/util"

	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type UninstallLAASTestRunner struct {
	test.BaseTestRunner
}

func (r *UninstallLAASTestRunner) Init(
	ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.BaseInit(r.Name(), ctx, config, clusterClients, clusterTargets, testObjects)
}

func (r *UninstallLAASTestRunner) Name() string {
	return "UninstallLAAS"
}

func (r *UninstallLAASTestRunner) Description() string {
	description := `This test uninstalls the Landscaper-As-A-Service controller.
The test succeeds when the installation has been deleted before the timeout expires.
`
	return description
}

func (r *UninstallLAASTestRunner) String() string {
	return r.Name()
}

func (r *UninstallLAASTestRunner) Run() error {
	logger, _ := logging.FromContextOrNew(r.GetCtx(), nil)

	logger.Info("deleting laas service target config")
	if err := r.deleteServiceTargetConfig(); err != nil {
		return err
	}

	logger.Info("deleting laas installation")
	if err := r.deleteInstallation(); err != nil {
		return err
	}

	return nil
}

func (r *UninstallLAASTestRunner) deleteServiceTargetConfig() error {
	serviceTargetConfig := r.GetTestObjects().ServiceTargetConfigs[types.NamespacedName{Name: r.GetClusterTargets().LaasCluster.Name, Namespace: r.GetConfig().LaasNamespace}.String()]

	if err := r.GetClusterClients().TestCluster.Delete(r.GetCtx(), serviceTargetConfig); err != nil {
		return fmt.Errorf("failed to delete service hostingTarget config: %w", err)
	}

	return nil
}

func (r *UninstallLAASTestRunner) deleteInstallation() error {
	installation := r.GetTestObjects().Installations[types.NamespacedName{Name: "laas", Namespace: r.GetConfig().LaasNamespace}.String()]

	if err := r.GetClusterClients().TestCluster.Delete(r.GetCtx(), installation); err != nil {
		return fmt.Errorf("failed to delete laas installation: %w", err)
	}

	timeout, err := cliutil.CheckAndWaitUntilObjectNotExistAnymore(
		r.GetClusterClients().TestCluster,
		types.NamespacedName{Name: installation.Name, Namespace: installation.Namespace},
		installation,
		r.GetConfig().SleepTime,
		r.GetConfig().MaxRetries)

	if err != nil {
		return fmt.Errorf("failed to wait for laas installation to be deleted: %w", err)
	}
	if timeout {
		return fmt.Errorf("waiting for laas installation to be deleted timed out")
	}

	return nil
}
