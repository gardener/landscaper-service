// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"encoding/json"
	"fmt"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	cliutil "github.com/gardener/landscapercli/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lsscore "github.com/gardener/landscaper-service/pkg/apis/core"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lssutils "github.com/gardener/landscaper-service/pkg/utils"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type InstallLAASTestRunner struct {
	BaseTestRunner
}

func (r *InstallLAASTestRunner) Init(
	ctx context.Context, config *test.TestConfig,
	clusterClients *test.ClusterClients, clusterTargets *test.ClusterTargets, testObjects *test.SharedTestObjects) {
	r.BaseInit(r.Name(), ctx, config, clusterClients, clusterTargets, testObjects)
}

func (r *InstallLAASTestRunner) Name() string {
	return "InstallLAAS"
}

func (r *InstallLAASTestRunner) Description() string {
	description := `This test installs the Landscaper-As-A-Service controller component.
The test succeeds when the installation is in phase succeeded before the timeout expires.
Otherwise the test fails.
`
	return description
}

func (r *InstallLAASTestRunner) String() string {
	return r.Name()
}

func (r *InstallLAASTestRunner) Run() error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	logger.Info("creating installation")
	if err := r.createInstallation(); err != nil {
		return err
	}

	logger.Info("creating service target config")
	if err := r.createServiceTargetConfig(); err != nil {
		return err
	}

	return nil
}

func (r *InstallLAASTestRunner) createServiceTargetConfig() error {
	serviceTargetConfig := &lssv1alpha1.ServiceTargetConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.clusterTargets.LaasCluster.Name,
			Namespace: r.config.LaasNamespace,
			Labels: map[string]string{
				lsscore.ServiceTargetConfigVisibleLabelName: "true",
				lsscore.ServiceTargetConfigRegionLabelName:  "eu",
			},
		},
		Spec: lssv1alpha1.ServiceTargetConfigSpec{
			ProviderType: r.config.ProviderType,
			Priority:     0,
			SecretRef: lssv1alpha1.SecretReference{
				ObjectReference: lssv1alpha1.ObjectReference{
					Name:      r.clusterTargets.LaasCluster.Name,
					Namespace: r.config.LaasNamespace,
				},
				Key: "kubeconfig",
			},
		},
	}

	if err := r.clusterClients.TestCluster.Create(r.ctx, serviceTargetConfig); err != nil {
		return fmt.Errorf("failed to create service hostingTarget config: %w", err)
	}

	r.testObjects.ServiceTargetConfigs[types.NamespacedName{Name: serviceTargetConfig.Name, Namespace: serviceTargetConfig.Namespace}.String()] = serviceTargetConfig

	return nil
}

func (r *InstallLAASTestRunner) createInstallation() error {
	logger, _ := logging.FromContextOrNew(r.ctx, nil)

	registryPullSecrets := []lssv1alpha1.ObjectReference{
		{
			Name:      "laas",
			Namespace: r.config.LaasNamespace,
		},
	}

	registryPullSecretsRaw, err := json.Marshal(registryPullSecrets)
	if err != nil {
		return fmt.Errorf("failed to marshal registry pull secrets: %w", err)
	}

	installation := &lsv1alpha1.Installation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "laas",
			Namespace: r.config.LaasNamespace,
			Annotations: map[string]string{
				lsv1alpha1.OperationAnnotation: string(lsv1alpha1.ReconcileOperation),
			},
		},
		Spec: lsv1alpha1.InstallationSpec{
			Context: "laas",
			ComponentDescriptor: &lsv1alpha1.ComponentDescriptorDefinition{
				Reference: &lsv1alpha1.ComponentDescriptorReference{
					RepositoryContext: cdv2.NewUnstructuredType(cdv2.OCIRegistryType, map[string]interface{}{
						"baseUrl": r.config.LaasRepository,
					}),
					ComponentName: r.config.LaasComponent,
					Version:       r.config.LaasVersion,
				},
			},
			Blueprint: lsv1alpha1.BlueprintDefinition{
				Reference: &lsv1alpha1.RemoteBlueprintReference{
					ResourceName: "landscaper-service-blueprint",
				},
			},
			Imports: lsv1alpha1.InstallationImports{
				Targets: []lsv1alpha1.TargetImport{
					{
						Name:   "targetCluster",
						Target: fmt.Sprintf("#%s", r.clusterTargets.HostingCluster.Name),
					},
				},
			},
			ImportDataMappings: map[string]lsv1alpha1.AnyJSON{
				"namespace":           lssutils.StringToAnyJSON(r.config.LaasNamespace),
				"verbosity":           lssutils.StringToAnyJSON(logging.DEBUG.String()),
				"registryPullSecrets": lsv1alpha1.NewAnyJSON(registryPullSecretsRaw),
			},
		},
	}

	if err := r.clusterClients.TestCluster.Create(r.ctx, installation); err != nil {
		return fmt.Errorf("failed to create installation: %w", err)
	}

	timeout, err := cliutil.CheckAndWaitUntilLandscaperInstallationSucceeded(r.clusterClients.TestCluster,
		types.NamespacedName{Name: installation.Name, Namespace: installation.Namespace},
		r.config.SleepTime, r.config.MaxRetries)
	if err != nil || timeout {
		if err := r.clusterClients.TestCluster.Get(r.ctx, types.NamespacedName{Name: installation.Name, Namespace: installation.Namespace}, installation); err == nil {
			logger.Error(fmt.Errorf("installation failed"), "installation", "last error", installation.Status.LastError)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to wait for installation to be ready: %w", err)
	}
	if timeout {
		return fmt.Errorf("waiting for installation timed out")
	}

	r.testObjects.Installations[types.NamespacedName{Name: installation.Name, Namespace: installation.Namespace}.String()] = installation

	return nil
}
