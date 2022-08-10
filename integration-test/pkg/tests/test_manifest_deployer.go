// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	cliutil "github.com/gardener/landscapercli/pkg/util"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lssutils "github.com/gardener/landscaper-service/pkg/utils"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
	"github.com/gardener/landscaper-service/test/integration/pkg/util"
)

const (
	manifestTestNamespace         = "example"
	manifestTestTargetName        = "default-target"
	manifestTestInstallationName  = "manifest-test"
	manifestTestComponentName     = "github.com/gardener/landscaper-examples/manifest-deployer/create-configmap"
	manifestTestComponentVersion  = "v0.1.0"
	manifestTestRepositoryContext = "eu.gcr.io/gardener-project/landscaper/examples"
	manifestTestConfigmapName     = "test-configmap"
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
	description := `This test creates an installation on the tenant virtual cluster using the Landscaper Manifest Deployer.
The target used by the installation points to the test cluster. The test succeeds when the installation is in phase succeeded
before the timeout expires and the configmap is correctly created in the target cluster.
`
	return description
}

func (r *ManifestDeployerTestRunner) String() string {
	return r.Name()
}

func (r *ManifestDeployerTestRunner) Run() error {
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

func (r *ManifestDeployerTestRunner) createVirtualClusterClient(deployment *lssv1alpha1.LandscaperDeployment) (client.Client, error) {
	r.log.Info("creating virtual cluster client for deployment", "deploymentName", deployment.Name)

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

func (r *ManifestDeployerTestRunner) prepare(virtualClient client.Client) (string, error) {
	namespace := &corev1.Namespace{}
	if err := r.clusterClients.TestCluster.Get(r.ctx, types.NamespacedName{Name: manifestTestNamespace}, namespace); err != nil {
		if !apierrors.IsNotFound(err) {
			r.log.Error(err, "failed to get test namespace", "namespace", manifestTestNamespace)
			return "", err
		}
	} else {
		r.log.Info("deleting namespace", "name", manifestTestNamespace)
		if err := cliutil.DeleteNamespace(r.clusterClients.TestCluster, manifestTestNamespace, r.config.SleepTime, r.config.MaxRetries); err != nil {
			r.log.Error(err, "failed to delete test namespace", "namespace", manifestTestNamespace)
			return "", err
		}
	}

	namespace = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: manifestTestNamespace,
		},
	}
	r.log.Info("creating namespace in test cluster", "name", manifestTestNamespace)
	if err := r.clusterClients.TestCluster.Create(r.ctx, namespace); err != nil {
		r.log.Error(err, "failed to create test namespace", "namespace", manifestTestNamespace)
		return "", err
	}

	namespace = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "manifest-test-",
		},
	}
	r.log.Info("creating namespace in virtual cluster", "generateName", namespace.GenerateName)
	if err := virtualClient.Create(r.ctx, namespace); err != nil {
		r.log.Error(err, "failed to create namespace in virtual cluster", "generateName", namespace.GenerateName)
		return "", err
	}

	return namespace.Name, nil
}

func (r *ManifestDeployerTestRunner) createTarget(deployment *lssv1alpha1.LandscaperDeployment, virtualClient client.Client, virtualClusterNamespace string) error {
	r.log.Info("creating target for deployment", "deploymentName", deployment.Name)
	if _, err := util.BuildKubernetesClusterTarget(r.ctx, virtualClient, r.config.TestClusterKubeconfig, manifestTestTargetName, virtualClusterNamespace); err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}
	return nil
}

func (r *ManifestDeployerTestRunner) createInstallation(deployment *lssv1alpha1.LandscaperDeployment, virtualClient client.Client, virtualClusterNamespace string) error {
	r.log.Info("creating installation for deployment",
		"deploymentName", deployment.Name,
		"installationName", manifestTestInstallationName,
		"installationNamespace", virtualClusterNamespace)

	installation := &lsv1alpha1.Installation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      manifestTestInstallationName,
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
					ComponentName: manifestTestComponentName,
					Version:       manifestTestComponentVersion,
					RepositoryContext: cdv2.NewUnstructuredType(cdv2.OCIRegistryType, map[string]interface{}{
						"baseUrl": manifestTestRepositoryContext,
					}),
				},
			},
			Imports: lsv1alpha1.InstallationImports{
				Targets: []lsv1alpha1.TargetImport{
					{
						Name:   "cluster",
						Target: fmt.Sprintf("#%s", manifestTestTargetName),
					},
				},
			},
			ImportDataMappings: map[string]lsv1alpha1.AnyJSON{
				"configmapName": lssutils.StringToAnyJSON(manifestTestConfigmapName),
			},
		},
		Status: lsv1alpha1.InstallationStatus{},
	}

	if err := virtualClient.Create(r.ctx, installation); err != nil {
		return fmt.Errorf("failed to create installation: %w", err)
	}

	return nil
}

func (r *ManifestDeployerTestRunner) verifyInstallation(deployment *lssv1alpha1.LandscaperDeployment, virtualClient client.Client, virtualClusterNamespace string) error {
	r.log.Info("verifying installation for deployment",
		"deploymentName", deployment.Name,
		"installationName", manifestTestInstallationName,
		"installationNamespace", virtualClusterNamespace)

	timeout, err := cliutil.CheckAndWaitUntilLandscaperInstallationSucceeded(
		virtualClient,
		types.NamespacedName{Name: manifestTestInstallationName, Namespace: virtualClusterNamespace},
		r.config.SleepTime, r.config.MaxRetries)

	if err != nil || timeout {
		installation := &lsv1alpha1.Installation{}
		if err := virtualClient.Get(r.ctx, types.NamespacedName{Name: manifestTestInstallationName, Namespace: virtualClusterNamespace}, installation); err == nil {
			r.log.Error(fmt.Errorf("installation failed"), "installation", "last error", installation.Status.LastError)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to wait for installation to be ready: %w", err)
	}
	if timeout {
		return fmt.Errorf("waiting for installation timed out")
	}

	configMap := &corev1.ConfigMap{}
	if err := r.clusterClients.TestCluster.Get(r.ctx, types.NamespacedName{Name: manifestTestConfigmapName, Namespace: manifestTestNamespace}, configMap); err != nil {
		return fmt.Errorf("failed to get deployed configmap: %w", err)
	}
	return nil
}
