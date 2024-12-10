// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	cliutil "github.com/gardener/landscapercli/pkg/util"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lssutils "github.com/gardener/landscaper-service/pkg/utils"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
	"github.com/gardener/landscaper-service/test/integration/pkg/util"
)

const (
	helmTestNamespace         = "helm-test"
	helmTestTargetName        = "default-target"
	helmTestInstallationName  = "helm-test"
	helmTestComponentName     = "github.com/gardener/landscaper-examples/helm-deployer/helm-chart-1"
	helmTestComponentVersion  = "v0.1.0"
	helmTestRepositoryContext = "europe-docker.pkg.dev/sap-gcp-cp-k8s-stable-hub/landscaper-examples/examples"
	helmTestConfigmapName     = "test-configmap"
)

type HelmDeployerTestRunner struct {
	BaseTestRunner
}

func (r *HelmDeployerTestRunner) Init(
	ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.BaseInit(r.Name(), ctx, config, clusterClients, clusterTargets, testObjects)
}

func (r *HelmDeployerTestRunner) Name() string {
	return "HelmDeployer"
}

func (r *HelmDeployerTestRunner) Description() string {
	description := `This test creates an installation on the tenant virtual cluster using the Landscaper Helm Deployer.
The target used by the installation points to the test cluster. The test succeeds when the installation is in phase succeeded
before the timeout expires and the configmap is correctly created in the target cluster.
`
	return description
}

func (r *HelmDeployerTestRunner) String() string {
	return r.Name()
}

func (r *HelmDeployerTestRunner) Run() error {
	for _, deployment := range r.testObjects.LandscaperDeployments {
		virtualClient, err := r.createVirtualClusterClient(deployment)
		if err != nil {
			return err
		}
		virtualClusterNamespace, err := r.prepare(virtualClient)
		if err != nil {
			return err
		}
		if err := r.createTarget(deployment, virtualClient, virtualClusterNamespace); err != nil {
			return err
		}
		if err := r.createInstallation(deployment, virtualClient, virtualClusterNamespace); err != nil {
			return err
		}
		if err := r.verifyInstallation(deployment, virtualClient, virtualClusterNamespace); err != nil {
			return err
		}
	}
	return nil
}

func (r *HelmDeployerTestRunner) createVirtualClusterClient(deployment *lssv1alpha1.LandscaperDeployment) (client.Client, error) {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("creating virtual cluster client for deployment", "deploymentName", deployment.Name)

	instance := &lssv1alpha1.Instance{}
	if err := r.clusterClients.TestCluster.Get(r.ctx, deployment.Status.InstanceRef.NamespacedName(), instance); err != nil {
		return nil, fmt.Errorf("failed to get instance for deployment: %w", err)
	}

	virtualClient, err := util.BuildKubeClientForInstance(instance, test.Scheme())
	if err != nil {
		return nil, err
	}

	return virtualClient, nil
}

func (r *HelmDeployerTestRunner) prepare(virtualClient client.Client) (string, error) {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	namespace := &corev1.Namespace{}
	if err := r.clusterClients.TestCluster.Get(r.ctx, types.NamespacedName{Name: helmTestNamespace}, namespace); err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Error(err, "failed to get test namespace", "namespace", helmTestNamespace)
			return "", err
		}
	} else {
		logger.Info("deleting namespace", "name", helmTestNamespace)
		if err := cliutil.DeleteNamespace(r.clusterClients.TestCluster, helmTestNamespace, r.config.SleepTime, r.config.MaxRetries); err != nil {
			logger.Error(err, "failed to delete test namespace", "namespace", helmTestNamespace)
			return "", err
		}
	}

	namespace = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: helmTestNamespace,
		},
	}
	logger.Info("creating namespace in test cluster", "name", helmTestNamespace)
	if err := r.clusterClients.TestCluster.Create(r.ctx, namespace); err != nil {
		logger.Error(err, "failed to create test namespace", "namespace", helmTestNamespace)
		return "", err
	}

	namespace = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "helm-test-",
		},
	}
	logger.Info("creating namespace in virtual cluster", "generateName", namespace.GenerateName)
	if err := virtualClient.Create(r.ctx, namespace); err != nil {
		logger.Error(err, "failed to create namespace in virtual cluster", "generateName", namespace.GenerateName)
		return "", err
	}

	return namespace.Name, nil
}

func (r *HelmDeployerTestRunner) createTarget(deployment *lssv1alpha1.LandscaperDeployment, virtualClient client.Client, virtualClusterNamespace string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("creating target for deployment", "deploymentName", deployment.Name)
	if _, err := util.BuildKubernetesClusterTargetWithSecretRef(r.ctx, virtualClient, r.config.TestClusterKubeconfig, helmTestTargetName, virtualClusterNamespace); err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}
	return nil
}

func (r *HelmDeployerTestRunner) createInstallation(deployment *lssv1alpha1.LandscaperDeployment, virtualClient client.Client, virtualClusterNamespace string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("creating installation for deployment",
		"deploymentName", deployment.Name,
		"installationName", helmTestInstallationName,
		"installationNamespace", virtualClusterNamespace)

	release := map[string]interface{}{
		"name":      helmTestConfigmapName,
		"namespace": helmTestNamespace,
	}

	releaseBytes, err := json.Marshal(release)
	if err != nil {
		return fmt.Errorf("faield to marshal release import: %w", err)
	}

	installation := &lsv1alpha1.Installation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helmTestInstallationName,
			Namespace: virtualClusterNamespace,
			Annotations: map[string]string{
				lsv1alpha1.OperationAnnotation: string(lsv1alpha1.ReconcileOperation),
			},
		},
		Spec: lsv1alpha1.InstallationSpec{
			Blueprint: lsv1alpha1.BlueprintDefinition{
				Reference: &lsv1alpha1.RemoteBlueprintReference{
					ResourceName: "blueprint",
				},
			},
			ComponentDescriptor: &lsv1alpha1.ComponentDescriptorDefinition{
				Reference: &lsv1alpha1.ComponentDescriptorReference{
					ComponentName: helmTestComponentName,
					Version:       helmTestComponentVersion,
					RepositoryContext: cdv2.NewUnstructuredType(cdv2.OCIRegistryType, map[string]interface{}{
						"baseUrl": helmTestRepositoryContext,
					}),
				},
			},
			Imports: lsv1alpha1.InstallationImports{
				Targets: []lsv1alpha1.TargetImport{
					{
						Name:   "cluster",
						Target: fmt.Sprintf("#%s", helmTestTargetName),
					},
				},
			},
			ImportDataMappings: map[string]lsv1alpha1.AnyJSON{
				"release":    lsv1alpha1.NewAnyJSON(releaseBytes),
				"testDataIn": lssutils.StringToAnyJSON("helmTest"),
			},
			Exports: lsv1alpha1.InstallationExports{
				Data: []lsv1alpha1.DataExport{
					{
						Name:    "testDataOut",
						DataRef: "do-testdata-out",
					},
				},
			},
		},
		Status: lsv1alpha1.InstallationStatus{},
	}

	if err := virtualClient.Create(r.ctx, installation); err != nil {
		return fmt.Errorf("failed to create installation: %w", err)
	}

	return nil
}

func (r *HelmDeployerTestRunner) verifyInstallation(deployment *lssv1alpha1.LandscaperDeployment, virtualClient client.Client, virtualClusterNamespace string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("verifying installation for deployment",
		"deploymentName", deployment.Name,
		"installationName", helmTestInstallationName,
		"installationNamespace", virtualClusterNamespace)

	timeout, err := cliutil.CheckAndWaitUntilLandscaperInstallationSucceeded(
		virtualClient,
		types.NamespacedName{Name: helmTestInstallationName, Namespace: virtualClusterNamespace},
		r.config.SleepTime, r.config.MaxRetries)

	if err != nil || timeout {
		installation := &lsv1alpha1.Installation{}
		if err := virtualClient.Get(r.ctx, types.NamespacedName{Name: helmTestInstallationName, Namespace: virtualClusterNamespace}, installation); err == nil {
			logger.Error(fmt.Errorf("installation failed"), "installation", "last error", installation.Status.LastError)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to wait for installation to be ready: %w", err)
	}
	if timeout {
		return fmt.Errorf("waiting for installation timed out")
	}

	configMap := &corev1.ConfigMap{}
	if err := r.clusterClients.TestCluster.Get(r.ctx, types.NamespacedName{Name: helmTestConfigmapName, Namespace: helmTestNamespace}, configMap); err != nil {
		return fmt.Errorf("failed to get deployed configmap: %w", err)
	}
	return nil
}
