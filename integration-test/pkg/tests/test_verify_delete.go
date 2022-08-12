// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"

	cliutil "github.com/gardener/landscapercli/pkg/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"

	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type VerifyDeleteRunner struct {
	ctx            context.Context
	log            logr.Logger
	config         *test.TestConfig
	clusterClients *test.ClusterClients
	clusterTargets *test.ClusterTargets
	testObjects    *test.SharedTestObjects
}

func (r *VerifyDeleteRunner) Init(
	ctx context.Context, log logr.Logger, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.ctx = ctx
	r.log = log.WithName(r.Name())
	r.config = config
	r.clusterClients = clusterClients
	r.clusterTargets = clusterTargets
	r.testObjects = testObjects
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
	r.log.Info("verifying namespace being deleted", "name", namespaceName)

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
	}, r.config.SleepTime, r.config.MaxRetries*5)

	if timeout {
		return fmt.Errorf("timeout while waiting for namespace %q being deleted", namespaceName)
	}
	if err != nil {
		return fmt.Errorf("error while waiting for namespace %q being deleted: %w", namespaceName, err)
	}

	return nil
}
