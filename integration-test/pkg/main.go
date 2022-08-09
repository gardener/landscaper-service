// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/controller-utils/pkg/logger"
	cliquickstart "github.com/gardener/landscapercli/cmd/quickstart"
	cliutil "github.com/gardener/landscapercli/pkg/util"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
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
		new(tests.UpdateDeploymentRunner),
		new(tests.VerifyDeploymentRunner),
		new(tests.ManifestDeployerTestRunner),
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

	ctx := context.Background()
	defer ctx.Done()

	log, err := logger.NewCliLogger()
	if err != nil {
		return err
	}

	logger.SetLogger(log)

	log.Info("running integration test with flags",
		"LAAS Component", config.LaasComponent,
		"LAAS Version", config.LaasVersion,
		"LAAS Repository", config.LaasRepository,
		"Landscaper Version", config.LandscaperVersion,
		"Landscaper Namespace", config.LandscaperNamespace,
		"LAAS Namespace", config.LaasNamespace,
		"Test Namespace", config.TestNamespace,
		"Provider Type", config.ProviderType,
	)

	clusterClients, err := test.NewClusterClients(config)
	if err != nil {
		return err
	}

	log.Info("========== Uninstalling Landscaper ==========")
	if err := uninstallLandscaper(ctx, log, clusterClients.TestCluster, config); err != nil {
		return err
	}

	log.Info("========== Cleaning up before test ==========")
	if err := cleanupResources(ctx, log, clusterClients.TestCluster, clusterClients.HostingCluster, config); err != nil {
		return err
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

	return runTestSuite(ctx, log, clusterClients, clusterTargets, config)
}

// runTestSuite runs the tests defined in integrationTests
func runTestSuite(ctx context.Context, log logr.Logger, clusterClients *test.ClusterClients, clusterTarget *test.ClusterTargets, config *test.TestConfig) error {
	log.Info("========== Running test suite ==========")
	testObjects := &test.SharedTestObjects{
		Installations:            make(map[string]*lsv1alpha1.Installation),
		ServiceTargetConfigs:     make(map[string]*lssv1alpha1.ServiceTargetConfig),
		LandscaperDeployments:    make(map[string]*lssv1alpha1.LandscaperDeployment),
		HostingClusterNamespaces: make([]string, 0),
	}

	succeededTests := make([]test.TestRunner, 0, len(integrationTests))
	testsNotRun := make([]test.TestRunner, len(integrationTests))
	copy(testsNotRun, integrationTests)

	log.Info("following tests will be run")
	for _, runner := range integrationTests {
		log.Info("test", "name", runner.Name(), "description", runner.Description())
	}

	for _, runner := range integrationTests {
		testsNotRun = testsNotRun[1:]
		log.Info("********** running test", "name", runner.Name())
		runner.Init(ctx, log, config, clusterClients, clusterTarget, testObjects)
		if err := runner.Run(); err != nil {
			return logTestSummary(log, succeededTests, testsNotRun, err, runner)
		}
		succeededTests = append(succeededTests, runner)
	}

	return logTestSummary(log, succeededTests, testsNotRun, nil, nil)
}

func logTestSummary(log logr.Logger, succeededTests, testsNotRun []test.TestRunner, err error, failedTest test.TestRunner) error {
	log.Info("==========  Test summary ==========")
	log.Info("successful tests", "tests", fmt.Sprintf("%v", succeededTests), "total", fmt.Sprintf("%d/%d", len(succeededTests), len(integrationTests)))
	log.Info("tests not run", "tests", fmt.Sprintf("%v", testsNotRun), "total", fmt.Sprintf("%d/%d", len(testsNotRun), len(integrationTests)))
	if err != nil {
		log.Error(err, "error while running test", "name", failedTest.Name())
	}
	return err
}

// uninstallLandscaper uninstalls the landscaper including the namespace it is installed in
func uninstallLandscaper(ctx context.Context, log logr.Logger, kclient client.Client, config *test.TestConfig) error {
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

	if err := util.ForceDeleteInstallations(ctx, log, kclient, config.TestClusterKubeconfig, config.LandscaperNamespace); err != nil {
		return err
	}

	if err := util.RemoveFinalizerLandscaperResources(ctx, kclient, config.LandscaperNamespace); err != nil {
		return err
	}

	log.Info("Waiting for resources to be deleted on the K8s cluster...")
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
    verbosity: 10
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
            verbosityLevel: 10
        helm:
          deployer:
            verbosityLevel: 10
        manifest:
          deployer:
            verbosityLevel: 10
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

	tmpFile, err := ioutil.TempFile(".", "landscaper-values-")
	if err != nil {
		return fmt.Errorf("cannot create temporary file: %w", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			fmt.Printf("Cannot remove temporary file %s: %s", tmpFile.Name(), err.Error())
		}
	}()

	if err := ioutil.WriteFile(tmpFile.Name(), []byte(landscaperValues), os.ModePerm); err != nil {
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
func cleanupResources(ctx context.Context, log logr.Logger, hostingClient, laasClient client.Client, config *test.TestConfig) error {
	// LAAS Namespace
	if err := util.DeleteValidatingWebhookConfiguration(ctx, hostingClient, "landscaper-service-validation-webhook", config.LaasNamespace); err != nil {
		return err
	}

	if err := util.RemoveFinalizerLandscaperResources(ctx, hostingClient, config.LaasNamespace); err != nil {
		return err
	}

	if err := util.RemoveFinalizerLaasResources(ctx, hostingClient, config.LaasNamespace); err != nil {
		return err
	}

	if err := util.CleanupLaasResources(ctx, log, hostingClient, config.LaasNamespace, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	if err := cliutil.DeleteNamespace(hostingClient, config.LaasNamespace, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	// Test Namespace
	if err := util.RemoveFinalizerLaasResources(ctx, hostingClient, config.TestNamespace); err != nil {
		return err
	}

	if err := util.CleanupLaasResources(ctx, log, hostingClient, config.TestNamespace, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	if err := util.RemoveFinalizerLandscaperResources(ctx, hostingClient, config.TestNamespace); err != nil {
		return err
	}

	if err := cliutil.DeleteNamespace(hostingClient, config.TestNamespace, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	if err := util.DeleteVirtualClusterNamespaces(ctx, log, laasClient, config.SleepTime, config.MaxRetries); err != nil {
		return err
	}

	return nil
}
