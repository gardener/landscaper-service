// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"
	"time"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	cliquickstart "github.com/gardener/landscapercli/cmd/quickstart"
	cliutil "github.com/gardener/landscapercli/pkg/util"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"

	"github.com/gardener/landscaper-service/test/integration/pkg/test"
	"github.com/gardener/landscaper-service/test/integration/pkg/tests"
	"github.com/gardener/landscaper-service/test/integration/pkg/util"
)

var (
	// tests are run in the order they are defined here
	integrationTests = []test.TestRunner{
		new(tests.InstallLAASTestRunner),
		new(tests.CreateDeploymentRunner),
		new(tests.VerifyDeploymentRunner),
		new(tests.NamespaceregistrationSubjectSyncRunner),
		new(tests.UpdateDeploymentRunner),
		new(tests.VerifyDeploymentRunner),
		new(tests.ManifestDeployerTestRunner),
		new(tests.HelmDeployerTestRunner),
		new(tests.ContainerDeployerTestRunner),
		new(tests.DeleteDeploymentRunner),
		new(tests.VerifyDeleteRunner),
		new(tests.UninstallLAASTestRunner),
	}
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while running integration test: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	config := test.ParseConfig()
	if err := test.VerifyConfig(config); err != nil {
		return err
	}

	landscaperVersion, err := util.GetLandscaperVersion(test.RepoRootDir)
	if err != nil {
		return err
	}

	config.LandscaperVersion = landscaperVersion
	config.GardenerProject = "laas"
	config.ShootSecretBindingName = "laas-canary"

	ctx := context.Background()
	defer ctx.Done()

	log, err := logging.NewCliLogger()
	if err != nil {
		return err
	}

	ctx = logging.NewContext(ctx, log)

	log.Info("running integration test with flags",
		"LAAS Component", config.LaasComponent,
		"LAAS Version", config.LaasVersion,
		"LAAS Repository", config.LaasRepository,
		"Landscaper Version", config.LandscaperVersion,
		"Landscaper Namespace", config.LandscaperNamespace,
		"LAAS Namespace", config.LaasNamespace,
		"Test Namespace", config.TestNamespace,
	)

	clusterClients, err := test.NewClusterClients(config)
	if err != nil {
		return err
	}

	log.Info("========== Uninstalling Landscaper ==========")
	if err := uninstallLandscaper(ctx, clusterClients.TestCluster, config); err != nil {
		log.Error(err, "uninstall landscaper")
	}

	log.Info("========== Cleaning up before test ==========")
	if err := cleanupResources(ctx, clusterClients.TestCluster, clusterClients.HostingCluster, config, log); err != nil {
		log.Error(err, "cleanup resources")
	}

	log.Info("========== Installing Landscaper ==========")
	if err := installLandscaper(ctx, config); err != nil {
		return err
	}

	log.Info("========== Creating prerequisites ==========")
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.LaasNamespace,
		},
	}

	if err := clusterClients.TestCluster.Create(ctx, namespace); err != nil {
		return fmt.Errorf("failed to create laas namespace: %w", err)
	}

	namespace = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.TestNamespace,
		},
	}

	if err := clusterClients.TestCluster.Create(ctx, namespace); err != nil {
		return fmt.Errorf("failed to create test namespace: %w", err)
	}

	clusterTargets, err := test.NewClusterTargets(ctx, clusterClients.TestCluster, config)
	if err != nil {
		return err
	}

	if err := util.BuildLandscaperContext(ctx, clusterClients.TestCluster, config.RegistryPullSecrets, "laas", config.LaasNamespace); err != nil {
		return fmt.Errorf("cannot build landscaper context: %w", err)
	}

	return runTestSuite(ctx, clusterClients, clusterTargets, config)
}

// runTestSuite runs the tests defined in integrationTests
func runTestSuite(ctx context.Context, clusterClients *test.ClusterClients, clusterTarget *test.ClusterTargets, config *test.TestConfig) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	logger.Info("========== Running test suite ==========")
	testObjects := &test.SharedTestObjects{
		Installations:            make(map[string]*lsv1alpha1.Installation),
		ServiceTargetConfigs:     make(map[string]*lssv1alpha1.ServiceTargetConfig),
		LandscaperDeployments:    make(map[string]*lssv1alpha1.LandscaperDeployment),
		HostingClusterNamespaces: make([]string, 0),
	}

	succeededTests := make([]test.TestRunner, 0, len(integrationTests))
	testsNotRun := make([]test.TestRunner, len(integrationTests))
	copy(testsNotRun, integrationTests)

	logger.Info("following tests will be run")
	for i, runner := range integrationTests {
		logger.Info(fmt.Sprintf("[%d] name: %s, description:\n%s", i, runner.Name(), runner.Description()))
	}

	for _, runner := range integrationTests {
		testsNotRun = testsNotRun[1:]
		logger.Info("********** running test", "name", runner.Name())
		runner.Init(ctx, config, clusterClients, clusterTarget, testObjects)
		if err := runner.Run(); err != nil {
			return logTestSummary(ctx, succeededTests, testsNotRun, err, runner)
		}
		succeededTests = append(succeededTests, runner)
	}

	return logTestSummary(ctx, succeededTests, testsNotRun, nil, nil)
}

func logTestSummary(ctx context.Context, succeededTests, testsNotRun []test.TestRunner, err error, failedTest test.TestRunner) error {
	logger, _ := logging.FromContextOrNew(ctx, nil)

	logger.Info("==========  Test summary ==========")
	logger.Info("successful tests", "tests", fmt.Sprintf("%v", succeededTests), "total", fmt.Sprintf("%d/%d", len(succeededTests), len(integrationTests)))
	logger.Info("tests not run", "tests", fmt.Sprintf("%v", testsNotRun), "total", fmt.Sprintf("%d/%d", len(testsNotRun), len(integrationTests)))
	if err != nil {
		logger.Error(err, "error while running test", "name", failedTest.Name())
	}
	return err
}

// uninstallLandscaper uninstalls the landscaper including the namespace it is installed in
func uninstallLandscaper(ctx context.Context, kclient client.Client, config *test.TestConfig) error {
	logger, _ := logging.FromContextOrNew(ctx, nil)

	landscaperNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.LandscaperNamespace,
		},
	}

	if err := kclient.Get(ctx, types.NamespacedName{Name: config.LandscaperNamespace}, landscaperNamespace); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		} else {
			return fmt.Errorf("failed to get landscaper namespace %s: %w", config.LandscaperNamespace, err)
		}
	}

	uninstallArgs := []string{
		"--kubeconfig",
		config.TestClusterKubeconfig,
		"--namespace",
		config.LandscaperNamespace,
		"--delete-namespace",
	}

	uninstallCmd := cliquickstart.NewUninstallCommand(ctx)
	uninstallCmd.SetArgs(uninstallArgs)

	if err := uninstallCmd.Execute(); err != nil {
		return fmt.Errorf("failed to uninstall landscaper: %w", err)
	}

	if err := util.DeleteValidatingWebhookConfiguration(ctx, kclient, "landscaper-validation-webhook", config.LandscaperNamespace); err != nil {
		return err
	}

	if err := util.ForceDeleteInstallations(ctx, kclient, config.TestClusterKubeconfig, config.LandscaperNamespace); err != nil {
		return err
	}

	if err := util.RemoveFinalizerLandscaperResources(ctx, kclient, config.LandscaperNamespace); err != nil {
		return err
	}

	logger.Info("Waiting for resources to be deleted on the K8s cluster...")
	namespace := &corev1.Namespace{}
	timeout, err := cliutil.CheckAndWaitUntilObjectNotExistAnymore(kclient, types.NamespacedName{Name: config.LandscaperNamespace}, namespace, config.SleepTime, config.MaxRetries)

	if err != nil {
		return err
	}

	if timeout {
		return fmt.Errorf("waiting for landscaper namespace deletion timed out")
	}

	return nil
}

// buildLandscaperValues builds the landscaper chart values
func buildLandscaperValues(namespace string) ([]byte, error) {
	const valuesTemplate = `
landscaper:
  landscaper:
    verbosity: debug
    registryConfig: # contains optional oci secrets
      allowPlainHttpRegistries: true
      secrets: {}
    deployers:
    - container
    - helm
    - manifest
    deployersConfig:
      Deployers:
        container:
          deployer:
            verbosityLevel: debug
        helm:
          deployer:
            verbosityLevel: debug
        manifest:
          deployer:
            verbosityLevel: debug
    deployerManagement:
      namespace: {{ .Namespace }}
      agent:
        namespace: {{ .Namespace }}
`

	t, err := template.New("valuesTemplate").Parse(valuesTemplate)
	if err != nil {
		return nil, err
	}

	data := struct {
		Namespace string
	}{
		Namespace: namespace,
	}

	b := &bytes.Buffer{}
	if err := t.Execute(b, data); err != nil {
		return nil, fmt.Errorf("could not template helm values: %w", err)
	}

	return b.Bytes(), nil
}

// installLandscaper installs the landscaper
func installLandscaper(ctx context.Context, config *test.TestConfig) error {
	landscaperValues, err := buildLandscaperValues(config.LandscaperNamespace)
	if err != nil {
		return fmt.Errorf("cannot template landscaper values: %w", err)
	}

	tmpFile, err := os.CreateTemp(".", "landscaper-values-")
	if err != nil {
		return fmt.Errorf("cannot create temporary file: %w", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			fmt.Printf("Cannot remove temporary file %s: %s", tmpFile.Name(), err.Error())
		}
	}()

	if err := os.WriteFile(tmpFile.Name(), []byte(landscaperValues), os.ModePerm); err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	}

	installArgs := []string{
		"--kubeconfig",
		config.TestClusterKubeconfig,
		"--landscaper-values",
		tmpFile.Name(),
		"--namespace",
		config.LandscaperNamespace,
		"--landscaper-chart-version",
		config.LandscaperVersion,
	}
	installCmd := cliquickstart.NewInstallCommand(ctx)
	installCmd.SetArgs(installArgs)

	if err := installCmd.Execute(); err != nil {
		return fmt.Errorf("install command failed: %w", err)
	}

	return nil
}

// cleanupResources removes all landscaper and laas resource in the laas and test namespace.
// It also tries to remove all virtual cluster namespaces that are still present in the cluster.
func cleanupResources(ctx context.Context, hostingClient, laasClient client.Client, config *test.TestConfig, log logging.Logger) error {
	// LAAS Namespace
	log.Info("execute DeleteValidatingWebhookConfiguration")
	if err := util.DeleteValidatingWebhookConfiguration(ctx, hostingClient, "landscaper-service-validation-webhook", config.LaasNamespace); err != nil {
		return err
	}

	log.Info("execute RemoveFinalizerLandscaperResources")
	if err := util.RemoveFinalizerLandscaperResources(ctx, hostingClient, config.LaasNamespace); err != nil {
		return err
	}

	log.Info("execute RemoveFinalizerLaasResources")
	if err := util.RemoveFinalizerLaasResources(ctx, hostingClient, config.LaasNamespace); err != nil {
		return err
	}

	log.Info("execute CleanupLaasResources")
	if err := util.CleanupLaasResources(ctx, hostingClient, config.LaasNamespace, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	// this should help prevent race conditions
	time.Sleep(time.Second * 10)

	log.Info("execute DeleteNamespace")
	if err := cliutil.DeleteNamespace(hostingClient, config.LaasNamespace, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	// Test Namespace
	log.Info("execute RemoveFinalizerLaasResources")
	if err := util.RemoveFinalizerLaasResources(ctx, hostingClient, config.TestNamespace); err != nil {
		return err
	}

	log.Info("execute CleanupLaasResources")
	if err := util.CleanupLaasResources(ctx, hostingClient, config.TestNamespace, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	log.Info("execute RemoveFinalizerLandscaperResources")
	if err := util.RemoveFinalizerLandscaperResources(ctx, hostingClient, config.TestNamespace); err != nil {
		return err
	}

	// this should help prevent race conditions
	time.Sleep(time.Second * 10)

	log.Info("execute DeleteNamespace")
	if err := cliutil.DeleteNamespace(hostingClient, config.TestNamespace, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	log.Info("execute DeleteTargetClusterNamespaces")
	if err := util.DeleteTargetClusterNamespaces(ctx, laasClient, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	log.Info("execute DeleteTestShootClusters")
	if err := util.DeleteTestShootClusters(ctx, config.GardenerServiceAccountKubeconfig, config.GardenerProject, config.TestPurpose, test.Scheme()); err != nil {
		return err
	}

	log.Info("cleanupResources finished")

	return nil
}
