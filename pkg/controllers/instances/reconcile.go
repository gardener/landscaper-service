// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/apis/core/v1alpha1/targettypes"
	"github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/apis/errors"
	lsinstallation "github.com/gardener/landscaper-service/pkg/apis/installation"
	"github.com/gardener/landscaper-service/pkg/utils"
)

var (
	// ClusterKubeconfigExportLabel is the label for the cluster kubeconfig export data object.
	ClusterKubeconfigExportLabel = fmt.Sprintf("%s=%s", lsv1alpha1.DataObjectKeyLabel, lsinstallation.ClusterKubeconfigExportName)
	// ClusterEndpointExportLabel is the label for the cluster api export data object.
	ClusterEndpointExportLabel = fmt.Sprintf("%s=%s", lsv1alpha1.DataObjectKeyLabel, lsinstallation.ClusterEndpointExportName)
)

// reconcile reconciles an instance.
func (c *Controller) reconcile(ctx context.Context, instance *lssv1alpha1.Instance) error {
	currOp := "Reconcile"

	if err := c.reconcileContext(ctx, instance); err != nil {
		return errors.NewWrappedError(err, currOp, "ReconcileContextFailed", err.Error())
	}

	if err := c.reconcileTarget(ctx, instance); err != nil {
		return errors.NewWrappedError(err, currOp, "ReconcileTargetFailed", err.Error())
	}

	if err := c.reconcileInstallation(ctx, instance); err != nil {
		return errors.NewWrappedError(err, currOp, "ReconcileInstallationFailed", err.Error())
	}

	return nil
}

// reconcileContext reconciles the context for an instance.
func (c *Controller) reconcileContext(ctx context.Context, instance *lssv1alpha1.Instance) error {
	landscaperContext := &lsv1alpha1.Context{}
	landscaperContext.GenerateName = fmt.Sprintf("%s-", instance.GetName())
	landscaperContext.Namespace = instance.GetNamespace()

	if instance.Status.ContextRef != nil && !instance.Status.ContextRef.IsEmpty() {
		landscaperContext.Name = instance.Status.ContextRef.Name
		landscaperContext.Namespace = instance.Status.ContextRef.Namespace
	}

	_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), landscaperContext, func() error {
		return c.mutateContext(ctx, landscaperContext, instance)
	})

	if err != nil {
		return fmt.Errorf("unable to create/update landscaperContext: %w", err)
	}

	if instance.Status.ContextRef == nil || !instance.Status.ContextRef.IsObject(landscaperContext) {
		instance.Status.ContextRef = &lssv1alpha1.ObjectReference{
			Name:      landscaperContext.GetName(),
			Namespace: landscaperContext.GetNamespace(),
		}

		if err := c.Client().Status().Update(ctx, instance); err != nil {
			return fmt.Errorf("unable to update target reference for instance: %w", err)
		}
	}

	return nil
}

// mutateTarget creates or updates the context for an instance.
func (c *Controller) mutateContext(ctx context.Context, context *lsv1alpha1.Context, instance *lssv1alpha1.Instance) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "mutateContext")

	if len(context.Name) > 0 {
		logger.Info("Updating context", lc.KeyResource, client.ObjectKeyFromObject(context).String())
	} else {
		logger.Info("Creating context", lc.KeyResource, types.NamespacedName{Name: context.GenerateName, Namespace: context.Namespace}.String())
	}

	repositoryContext := &cdv2.UnstructuredTypedObject{}
	err := json.Unmarshal(c.Config().LandscaperServiceComponent.RepositoryContext.RawMessage, repositoryContext)

	if err != nil {
		return fmt.Errorf("failed to unmarshal repository context: %w", err)
	}

	context.RepositoryContext = repositoryContext
	context.RegistryPullSecrets = make([]corev1.LocalObjectReference, 0, len(c.Config().LandscaperServiceComponent.RegistryPullSecrets))

	for i, secretRef := range c.Config().LandscaperServiceComponent.RegistryPullSecrets {
		configuredSecret := &corev1.Secret{}
		if err := c.Client().Get(ctx, client.ObjectKey{Name: secretRef.Name, Namespace: secretRef.Namespace}, configuredSecret); err != nil {
			return fmt.Errorf("unable to get registry pull secret \"%s/%s\": %w", secretRef.Namespace, secretRef.Name, err)
		}

		registryPullSecret := &corev1.Secret{}
		registryPullSecret.Name = fmt.Sprintf("%s-regsecret-%d", instance.Name, i)
		registryPullSecret.Namespace = context.Namespace

		_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), registryPullSecret, func() error {
			registryPullSecret.Type = configuredSecret.Type
			registryPullSecret.Data = configuredSecret.Data
			return nil
		})

		if err != nil {
			return fmt.Errorf("unable to create/update registry pull configuredSecret: %w", err)
		}

		context.RegistryPullSecrets = append(context.RegistryPullSecrets, corev1.LocalObjectReference{
			Name: registryPullSecret.Name,
		})
	}

	return nil
}

// reconcileTarget reconciles the target for an instance.
func (c *Controller) reconcileTarget(ctx context.Context, instance *lssv1alpha1.Instance) error {
	target := &lsv1alpha1.Target{}
	target.GenerateName = fmt.Sprintf("%s-", instance.GetName())
	target.Namespace = instance.GetNamespace()

	if instance.Status.TargetRef != nil && !instance.Status.TargetRef.IsEmpty() {
		target.Name = instance.Status.TargetRef.Name
		target.Namespace = instance.Status.TargetRef.Namespace
	}

	_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), target, func() error {
		return c.mutateTarget(ctx, target, instance)
	})

	if err != nil {
		return fmt.Errorf("unable to create/update target: %w", err)
	}

	if instance.Status.TargetRef == nil || !instance.Status.TargetRef.IsObject(target) {
		instance.Status.TargetRef = &lssv1alpha1.ObjectReference{
			Name:      target.GetName(),
			Namespace: target.GetNamespace(),
		}

		if err := c.Client().Status().Update(ctx, instance); err != nil {
			return fmt.Errorf("unable to update target reference for instance: %w", err)
		}
	}

	return nil
}

// mutateTarget creates or updates the target for an instance.
func (c *Controller) mutateTarget(ctx context.Context, target *lsv1alpha1.Target, instance *lssv1alpha1.Instance) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "mutateTarget")

	if len(target.Name) > 0 {
		logger.Info("Updating target", lc.KeyResource, client.ObjectKeyFromObject(target).String())
	} else {
		logger.Info("Creating target", lc.KeyResource, types.NamespacedName{Name: target.GenerateName, Namespace: target.Namespace}.String())
	}

	if err := controllerutil.SetControllerReference(instance, target, c.Scheme()); err != nil {
		return fmt.Errorf("unable to set controller reference for target: %w", err)
	}

	config := &lssv1alpha1.ServiceTargetConfig{}
	if err := c.Client().Get(ctx, instance.Spec.ServiceTargetConfigRef.NamespacedName(), config); err != nil {
		return fmt.Errorf("unable to get service target config for instance: %w", err)
	}

	secret := &corev1.Secret{}
	if err := c.Client().Get(ctx, config.Spec.SecretRef.NamespacedName(), secret); err != nil {
		return fmt.Errorf("unable to get kubeconfig secret for service target config: %w", err)
	}

	kubeconfig, ok := secret.Data[config.Spec.SecretRef.Key]
	if !ok {
		return fmt.Errorf("unable to read kubeconfig from secret: missing key %q", config.Spec.SecretRef.Key)
	}
	kubeconfigStr := string(kubeconfig)

	targetConfig := targettypes.KubernetesClusterTargetConfig{
		Kubeconfig: targettypes.ValueRef{
			StrVal: &kubeconfigStr,
		},
	}

	targetConfigRaw, err := json.Marshal(targetConfig)
	if err != nil {
		return fmt.Errorf("unable to marshal kubeconfig: %w", err)
	}
	targetConfigAnyJSON := lsv1alpha1.NewAnyJSON(targetConfigRaw)

	target.Spec = lsv1alpha1.TargetSpec{
		Type:          targettypes.KubernetesClusterTargetType,
		Configuration: &targetConfigAnyJSON,
	}

	return nil
}

// reconcileInstallation reconciles the installation for an instance
func (c *Controller) reconcileInstallation(ctx context.Context, instance *lssv1alpha1.Instance) error {
	installation := &lsv1alpha1.Installation{}
	installation.GenerateName = fmt.Sprintf("%s-", instance.GetName())
	installation.Namespace = instance.GetNamespace()

	if instance.Status.InstallationRef != nil && !instance.Status.InstallationRef.IsEmpty() {
		installation.Name = instance.Status.InstallationRef.Name
		installation.Namespace = instance.Status.InstallationRef.Namespace

		if err := c.Client().Get(ctx, instance.Status.InstallationRef.NamespacedName(), installation); err != nil {
			return fmt.Errorf("unable to get installation: %w", err)
		}
	}

	_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), installation, func() error {
		return c.mutateInstallation(ctx, installation, instance)
	})

	if err != nil {
		return fmt.Errorf("unable to create/update installation: %w", err)
	}

	old := instance.DeepCopy()
	instance.Status.InstallationRef = &lssv1alpha1.ObjectReference{
		Name:      installation.GetName(),
		Namespace: installation.GetNamespace(),
	}

	instance.Status.LandscaperServiceComponent = &lssv1alpha1.LandscaperServiceComponent{
		Name:    installation.Spec.ComponentDescriptor.Reference.ComponentName,
		Version: installation.Spec.ComponentDescriptor.Reference.Version,
	}

	if err := c.handleExports(ctx, instance, installation); err != nil {
		return err
	}

	if !reflect.DeepEqual(old.Status, instance.Status) {
		if err := c.Client().Status().Update(ctx, instance); err != nil {
			return fmt.Errorf("unable to update instance status: %w", err)
		}
	}

	return nil
}

// mutateInstallation creates or updates the installation for an instance.
func (c *Controller) mutateInstallation(ctx context.Context, installation *lsv1alpha1.Installation, instance *lssv1alpha1.Instance) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "mutateInstallation")

	if len(installation.Name) > 0 {
		logger.Info("Updating installation", lc.KeyResource, client.ObjectKeyFromObject(installation).String())
	} else {
		logger.Info("Creating installation", lc.KeyResource, types.NamespacedName{Name: installation.GenerateName, Namespace: installation.Namespace}.String())
	}

	// create a copy of the current installation spec for deciding whether a reconcile-annotation has to be set.
	oldInstallationSpec := installation.Spec.DeepCopy()

	if err := controllerutil.SetControllerReference(instance, installation, c.Scheme()); err != nil {
		return fmt.Errorf("unable to set owner reference for installation: %w", err)
	}

	config := &lssv1alpha1.ServiceTargetConfig{}
	if err := c.Client().Get(ctx, instance.Spec.ServiceTargetConfigRef.NamespacedName(), config); err != nil {
		return fmt.Errorf("unable to get service target config for instance: %w", err)
	}

	registryConfig := lsinstallation.NewRegistryConfig()
	registryConfigRaw, err := registryConfig.ToAnyJSON()
	if err != nil {
		return fmt.Errorf("unable to marshal registry config: %w", err)
	}

	landscaperConfig := lsinstallation.NewLandscaperConfig()
	landscaperConfig.Deployers = instance.Spec.LandscaperConfiguration.Deployers
	landscaperConfig.Landscaper.Verbosity = logging.INFO.String()
	landscaperConfigRaw, err := landscaperConfig.ToAnyJSON()
	if err != nil {
		return fmt.Errorf("unable to marshal landscaper config: %w", err)
	}

	installation.Spec = lsv1alpha1.InstallationSpec{
		Context: instance.Status.ContextRef.Name,
		ComponentDescriptor: &lsv1alpha1.ComponentDescriptorDefinition{
			Reference: &lsv1alpha1.ComponentDescriptorReference{
				ComponentName: c.Operation.Config().LandscaperServiceComponent.Name,
				Version:       c.Operation.Config().LandscaperServiceComponent.Version,
			},
		},
		Blueprint: lsv1alpha1.BlueprintDefinition{
			Reference: &lsv1alpha1.RemoteBlueprintReference{
				ResourceName: "installation-blueprint",
			},
		},
		Imports: lsv1alpha1.InstallationImports{
			Targets: []lsv1alpha1.TargetImport{
				{
					Name:   "hostingCluster",
					Target: fmt.Sprintf("#%s", instance.Status.TargetRef.Name),
				},
			},
		},
		ImportDataMappings: map[string]lsv1alpha1.AnyJSON{
			lsinstallation.HostingClusterNamespaceImportName: utils.StringToAnyJSON(fmt.Sprintf("%s-%s", instance.Spec.TenantId, instance.Spec.ID)),
			lsinstallation.DeleteHostingClusterImportName:    utils.BoolToAnyJSON(true),
			lsinstallation.VirtualClusterNamespaceImportName: utils.StringToAnyJSON(lsinstallation.VirtualClusterNamespace),
			lsinstallation.ProviderTypeImportName:            utils.StringToAnyJSON(config.Spec.ProviderType),
			lsinstallation.RegistryConfigImportName:          *registryConfigRaw,
			lsinstallation.LandscaperConfigImportName:        *landscaperConfigRaw,
			// TODO: how should this be configured?
			lsinstallation.DNSAccessDomainImportName: utils.StringToAnyJSON(""),
		},
		Exports: lsv1alpha1.InstallationExports{
			Data: []lsv1alpha1.DataExport{
				{
					Name:    lsinstallation.ClusterEndpointExportName,
					DataRef: lsinstallation.ClusterEndpointExportName,
				},
				{
					Name:    lsinstallation.ClusterKubeconfigExportName,
					DataRef: lsinstallation.ClusterKubeconfigExportName,
				},
			},
		},
	}

	if !deepEqualInstallationSpec(oldInstallationSpec, installation.Spec.DeepCopy()) {
		// set reconcile annotation to start/update the installation
		logger.Info("Setting reconcile operation annotation")
		if installation.Annotations == nil {
			installation.Annotations = make(map[string]string)
		}
		installation.Annotations[lsv1alpha1.OperationAnnotation] = string(lsv1alpha1.ReconcileOperation)
	}

	return nil
}

// deepEqualInstallationSpec tests whether two landscaper service installation specs are equal
func deepEqualInstallationSpec(specA, specB *lsv1alpha1.InstallationSpec) bool {
	landscaperConfigA := make(map[string]interface{})
	if err := json.Unmarshal(specA.ImportDataMappings[lsinstallation.LandscaperConfigImportName].RawMessage, &landscaperConfigA); err != nil {
		return false
	}
	landscaperConfigB := make(map[string]interface{})
	if err := json.Unmarshal(specB.ImportDataMappings[lsinstallation.LandscaperConfigImportName].RawMessage, &landscaperConfigB); err != nil {
		return false
	}
	if !reflect.DeepEqual(landscaperConfigA, landscaperConfigB) {
		return false
	}

	registryConfigA := make(map[string]interface{})
	if err := json.Unmarshal(specA.ImportDataMappings[lsinstallation.RegistryConfigImportName].RawMessage, &registryConfigA); err != nil {
		return false
	}
	registryConfigB := make(map[string]interface{})
	if err := json.Unmarshal(specA.ImportDataMappings[lsinstallation.RegistryConfigImportName].RawMessage, &registryConfigB); err != nil {
		return false
	}
	if !reflect.DeepEqual(registryConfigA, registryConfigB) {
		return false
	}

	delete(specA.ImportDataMappings, lsinstallation.LandscaperConfigImportName)
	delete(specB.ImportDataMappings, lsinstallation.LandscaperConfigImportName)
	delete(specA.ImportDataMappings, lsinstallation.RegistryConfigImportName)
	delete(specB.ImportDataMappings, lsinstallation.RegistryConfigImportName)

	return reflect.DeepEqual(specA, specB)
}

// handleExports tries to find the exports of the installation and update the instance status accordingly.
func (c *Controller) handleExports(ctx context.Context, instance *lssv1alpha1.Instance, installation *lsv1alpha1.Installation) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "handleExports")

	dataObjects := &lsv1alpha1.DataObjectList{}
	selectorBuilder := strings.Builder{}

	// Cluster Kubeconfig
	dataObjectSource := fmt.Sprintf("%s=Inst.%s,", lsv1alpha1.DataObjectSourceLabel, installation.GetName())
	selectorBuilder.WriteString(dataObjectSource)
	selectorBuilder.WriteString(ClusterKubeconfigExportLabel)

	labelSelector, _ := labels.Parse(selectorBuilder.String())
	listOptions := client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     installation.GetNamespace(),
	}

	if err := c.Client().List(ctx, dataObjects, &listOptions); err != nil {
		return fmt.Errorf("unable to list data objects for ClusterKubeconfig: %w", err)
	}

	if len(dataObjects.Items) > 0 {
		logger.Info("found export data object for cluster kubeconfig",
			lc.KeyResource, types.NamespacedName{Name: dataObjects.Items[0].Name, Namespace: dataObjects.Items[0].Namespace}.String())
		if err := json.Unmarshal(dataObjects.Items[0].Data.RawMessage, &instance.Status.ClusterKubeconfig); err != nil {
			return fmt.Errorf("unable to unmarshal ClusterKubeconfig: %w", err)
		}
	}

	// Cluster Endpoint
	selectorBuilder.Reset()
	selectorBuilder.WriteString(dataObjectSource)
	selectorBuilder.WriteString(ClusterEndpointExportLabel)

	labelSelector, _ = labels.Parse(selectorBuilder.String())
	listOptions = client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     installation.GetNamespace(),
	}

	if err := c.Client().List(ctx, dataObjects, &listOptions); err != nil {
		return fmt.Errorf("unable to list data objects for ClusterEndpoint: %w", err)
	}

	if len(dataObjects.Items) > 0 {
		logger.Info("found export data object for cluster endpoint",
			lc.KeyResource, types.NamespacedName{Name: dataObjects.Items[0].Name, Namespace: dataObjects.Items[0].Namespace}.String())
		if err := json.Unmarshal(dataObjects.Items[0].Data.RawMessage, &instance.Status.ClusterEndpoint); err != nil {
			return fmt.Errorf("unable to unmarshal ClusterEndpoint: %w", err)
		}
	}

	return nil
}
