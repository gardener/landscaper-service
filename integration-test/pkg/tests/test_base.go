// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

// BaseTestRunner is the base struct for test runners.
type BaseTestRunner struct {
	ctx            context.Context
	config         *test.TestConfig
	clusterClients *test.ClusterClients
	clusterTargets *test.ClusterTargets
	testObjects    *test.SharedTestObjects
}

// BaseInit initializes the test runner.
func (r *BaseTestRunner) BaseInit(name string, ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)
	logger = logger.WithName(name)
	r.ctx = logging.NewContext(ctx, logger)
	r.config = config
	r.clusterClients = clusterClients
	r.clusterTargets = clusterTargets
	r.testObjects = testObjects
}
