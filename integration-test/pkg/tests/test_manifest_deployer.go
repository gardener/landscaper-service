package tests

import (
	"context"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
	"github.com/go-logr/logr"
)

type ManifestDeployerTestRunner struct {
	ctx            context.Context
	log            logr.Logger
	config         *test.TestConfig
	clusterClients *test.ClusterClients
	clusterTargets *test.ClusterTargets
	testObjects    *test.SharedTestObjects
}

func (r *ManifestDeployerTestRunner) Init(
	ctx context.Context, log logr.Logger, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.ctx = ctx
	r.log = log.WithName(r.Name())
	r.config = config
	r.clusterClients = clusterClients
	r.clusterTargets = clusterTargets
	r.testObjects = testObjects
}

func (r *ManifestDeployerTestRunner) Name() string {
	return "ManifestDeployer"
}

func (r *ManifestDeployerTestRunner) Description() string {
	description := `This test creates an installation on the tenant virtual cluster using the Landscaper manifest deployer.
`
	return description
}

func (r *ManifestDeployerTestRunner) String() string {
	return r.Name()
}

func (r *ManifestDeployerTestRunner) Run() error {
	return nil
}
