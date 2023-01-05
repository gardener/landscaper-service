// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"

	cliutil "github.com/gardener/landscapercli/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type VerifyDeleteRunner struct {
	BaseTestRunner
}

func (r *VerifyDeleteRunner) Init(
	ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.BaseInit(r.Name(), ctx, config, clusterClients, clusterTargets, testObjects)
}

func (r *VerifyDeleteRunner) Name() string {
	return "VerifyDeploymentDeleted"
}

func (r *VerifyDeleteRunner) Description() string {
	description := `This test verifies that a tenant Landscaper deployment has been uninstalled and deleted correctly.
This test succeeds when the virtual-cluster-namespace of the tenant Landscaper deployment has been deleted before the
timeout expires. Otherwise the test fails.
`
	return description
}

func (r *VerifyDeleteRunner) String() string {
	return r.Name()
}

func (r *VerifyDeleteRunner) Run() error {
	for _, namespace := range r.testObjects.HostingClusterNamespaces {
		if err := r.verifyNamespace(namespace); err != nil {
			return err
		}
	}
	return nil
}

func (r *VerifyDeleteRunner) verifyNamespace(namespaceName string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("verifying namespace being deleted", "name", namespaceName)

	timeout, err := cliutil.CheckConditionPeriodically(func() (bool, error) {
		namespace := &corev1.Namespace{}
		if err := r.clusterClients.HostingCluster.Get(r.ctx, types.NamespacedName{Name: namespaceName}, namespace); err != nil {
			if k8serrors.IsNotFound(err) {
				return true, nil
			} else {
				return false, err
			}
		}

		return false, nil
	}, r.config.SleepTime, r.config.MaxRetries*8)

	if timeout {
		return fmt.Errorf("timeout while waiting for namespace %q being deleted", namespaceName)
	}
	if err != nil {
		return fmt.Errorf("error while waiting for namespace %q being deleted: %w", namespaceName, err)
	}

	return nil
}
