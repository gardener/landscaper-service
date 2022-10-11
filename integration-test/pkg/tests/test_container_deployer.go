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
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
	"github.com/gardener/landscaper-service/test/integration/pkg/util"
)

const (
	containerTestNamespace         = "example"
	containerTestTargetName        = "default-target"
	containerTestInstallationName  = "container-test"
	containerTestComponentName     = "github.com/gardener/landscaper-examples/container-deployer/container-1"
	containerTestComponentVersion  = "v0.1.0"
	containerTestRepositoryContext = "eu.gcr.io/gardener-project/landscaper/examples"
	containerTestConfigmapName     = "test-configmap"
)

type ContainerDeployerTestRunner struct {
	BaseTestRunner
}

func (r *ContainerDeployerTestRunner) Init(
	ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.BaseInit(r.Name(), ctx, config, clusterClients, clusterTargets, testObjects)
}

func (r *ContainerDeployerTestRunner) Name() string {
	return "ContainerDeployer"
}

func (r *ContainerDeployerTestRunner) Description() string {
	description := `This test creates an installation on the tenant virtual cluster using the Landscaper Container Deployer.
The target used by the installation points to the test cluster. The test succeeds when the installation is in phase succeeded
before the timeout expires and the configmap is correctly created in the target cluster.
`
	return description
}

func (r *ContainerDeployerTestRunner) String() string {
	return r.Name()
}

func (r *ContainerDeployerTestRunner) Run() error {
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

func (r *ContainerDeployerTestRunner) createVirtualClusterClient(deployment *lssv1alpha1.LandscaperDeployment) (client.Client, error) {
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

func (r *ContainerDeployerTestRunner) prepare(virtualClient client.Client) (string, error) {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	namespace := &corev1.Namespace{}
	if err := r.clusterClients.TestCluster.Get(r.ctx, types.NamespacedName{Name: containerTestNamespace}, namespace); err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Error(err, "failed to get test namespace", "namespace", containerTestNamespace)
			return "", err
		}
	} else {
		logger.Info("deleting namespace", "name", containerTestNamespace)
		if err := cliutil.DeleteNamespace(r.clusterClients.TestCluster, containerTestNamespace, r.config.SleepTime, r.config.MaxRetries); err != nil {
			logger.Error(err, "failed to delete test namespace", "namespace", containerTestNamespace)
			return "", err
		}
	}

	namespace = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: containerTestNamespace,
		},
	}
	logger.Info("creating namespace in test cluster", "name", containerTestNamespace)
	if err := r.clusterClients.TestCluster.Create(r.ctx, namespace); err != nil {
		logger.Error(err, "failed to create test namespace", "namespace", containerTestNamespace)
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

func (r *ContainerDeployerTestRunner) createTarget(deployment *lssv1alpha1.LandscaperDeployment, virtualClient client.Client, virtualClusterNamespace string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("creating target for deployment", "deploymentName", deployment.Name)
	if _, err := util.BuildKubernetesClusterTarget(r.ctx, virtualClient, r.config.TestClusterKubeconfig, containerTestTargetName, virtualClusterNamespace); err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}
	return nil
}

func (r *ContainerDeployerTestRunner) createInstallation(deployment *lssv1alpha1.LandscaperDeployment, virtualClient client.Client, virtualClusterNamespace string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("creating installation for deployment",
		"deploymentName", deployment.Name,
		"installationName", containerTestInstallationName,
		"installationNamespace", virtualClusterNamespace)

	configMapImport := map[string]interface{}{
		"name":      containerTestConfigmapName,
		"namespace": containerTestNamespace,
		"data": map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}

	configMapImportBytes, err := json.Marshal(configMapImport)
	if err != nil {
		return fmt.Errorf("faield to marshal configmap import: %w", err)
	}

	installation := &lsv1alpha1.Installation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      containerTestInstallationName,
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
					ComponentName: containerTestComponentName,
					Version:       containerTestComponentVersion,
					RepositoryContext: cdv2.NewUnstructuredType(cdv2.OCIRegistryType, map[string]interface{}{
						"baseUrl": containerTestRepositoryContext,
					}),
				},
			},
			Imports: lsv1alpha1.InstallationImports{
				Targets: []lsv1alpha1.TargetImport{
					{
						Name:   "targetCluster",
						Target: fmt.Sprintf("#%s", containerTestTargetName),
					},
				},
			},
			ImportDataMappings: map[string]lsv1alpha1.AnyJSON{
				"configmap": lsv1alpha1.NewAnyJSON(configMapImportBytes),
			},
			Exports: lsv1alpha1.InstallationExports{
				Data: []lsv1alpha1.DataExport{
					{
						Name:    "configMapData",
						DataRef: "configmapdata",
					},
					{
						Name:    "component",
						DataRef: "component",
					},
					{
						Name:    "content",
						DataRef: "content",
					},
					{
						Name:    "state",
						DataRef: "state",
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

func (r *ContainerDeployerTestRunner) verifyInstallation(deployment *lssv1alpha1.LandscaperDeployment, virtualClient client.Client, virtualClusterNamespace string) error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("verifying installation for deployment",
		"deploymentName", deployment.Name,
		"installationName", containerTestInstallationName,
		"installationNamespace", virtualClusterNamespace)

	timeout, err := cliutil.CheckAndWaitUntilLandscaperInstallationSucceeded(
		virtualClient,
		types.NamespacedName{Name: containerTestInstallationName, Namespace: virtualClusterNamespace},
		r.config.SleepTime, r.config.MaxRetries)

	if err != nil || timeout {
		installation := &lsv1alpha1.Installation{}
		if err := virtualClient.Get(r.ctx, types.NamespacedName{Name: containerTestInstallationName, Namespace: virtualClusterNamespace}, installation); err == nil {
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
	if err := r.clusterClients.TestCluster.Get(r.ctx, types.NamespacedName{Name: containerTestConfigmapName, Namespace: containerTestNamespace}, configMap); err != nil {
		return fmt.Errorf("failed to get deployed configmap: %w", err)
	}
	return nil
}
