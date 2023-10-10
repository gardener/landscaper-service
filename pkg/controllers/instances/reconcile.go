// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/apis/core/v1alpha1/targettypes"
	"github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lssv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha2"
	"github.com/gardener/landscaper-service/pkg/apis/errors"
	lsinstallation "github.com/gardener/landscaper-service/pkg/apis/installation"
	"github.com/gardener/landscaper-service/pkg/utils"
)

const (
	// shootAPIVersion is the gardener shoot version
	shootAPIVersion = "core.gardener.cloud/v1beta1"
	// shootKind is the gardener shoot kind
	shootKind = "Shoot"
	// automaticReconcileSeconds is the number of seconds after which installations of landscaper instances are
	// automatically reconciled. Important: the value must be shorter than tokenExpirationSeconds
	automaticReconcileSeconds = 14 * 24 * 60 * 60
	// failedReconcileSeconds is the number of seconds after which a failed landscaper instance is automatically reconciled.
	failedReconcileSeconds = 60 * 10
	// tokenExpirationSeconds defines how long the tokens are valid	which the landscaper and sidecar controllers use
	// to access the resource cluster, e.g. for watching installations, namespace registrations etc.
	// Important: the value must be larger than automaticReconcileSeconds.
	tokenExpirationSeconds = int64(90 * 24 * 60 * 60)
	// adminKubeconfigExpirationSeconds defines how long the admin kubeconfig for a resource cluster is valid.
	// This kubeconfig is used to deploy RBAC objects on the resource cluster. Maximum: 86400 (1 day).
	// Each reconcile uses a new kubeconfig, so that the short duration suffices.
	adminKubeconfigExpirationSeconds = int64(24 * 60 * 60)
)

// reconcile reconciles an instance.
func (c *Controller) reconcile(ctx context.Context, instance *lssv1alpha2.Instance) error {
	currOp := "Reconcile"

	if err := c.reconcileContext(ctx, instance); err != nil {
		return errors.NewWrappedError(err, currOp, "ReconcileContextFailed", err.Error())
	}

	if err := c.reconcileDataPlaneClusterTarget(ctx, instance); err != nil {
		return errors.NewWrappedError(err, currOp, "ReconcileGardenerServiceAccountTargetFailed", err.Error())
	}

	if err := c.reconcileTargetClusterTarget(ctx, instance); err != nil {
		return errors.NewWrappedError(err, currOp, "ReconcileTargetFailed", err.Error())
	}

	if err := c.reconcileInstallation(ctx, instance); err != nil {
		return errors.NewWrappedError(err, currOp, "ReconcileInstallationFailed", err.Error())
	}

	return nil
}

// reconcileContext reconciles the context for an instance.
func (c *Controller) reconcileContext(ctx context.Context, instance *lssv1alpha2.Instance) error {
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
		instance.Status.ContextRef = &lssv1alpha2.ObjectReference{
			Name:      landscaperContext.GetName(),
			Namespace: landscaperContext.GetNamespace(),
		}

		if err := c.Client().Status().Update(ctx, instance); err != nil {
			return fmt.Errorf("unable to update target reference for instance: %w", err)
		}
	}

	return nil
}

// mutateTargetClusterTarget creates or updates the context for an instance.
func (c *Controller) mutateContext(ctx context.Context, context *lsv1alpha1.Context, instance *lssv1alpha2.Instance) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "mutateContext")

	if len(context.Name) > 0 {
		logger.Info("Updating context", lc.KeyResource, client.ObjectKeyFromObject(context).String())
	} else {
		logger.Info("Creating context", lc.KeyResource, types.NamespacedName{Name: context.GenerateName, Namespace: context.Namespace}.String())
	}

	if err := controllerutil.SetControllerReference(instance, context, c.Scheme()); err != nil {
		return fmt.Errorf("unable to set controller reference for target: %w", err)
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

// reconcileTargetClusterTarget reconciles the target for an instance.
func (c *Controller) reconcileTargetClusterTarget(ctx context.Context, instance *lssv1alpha2.Instance) error {
	target := &lsv1alpha1.Target{}
	target.GenerateName = fmt.Sprintf("%s-", instance.GetName())
	target.Namespace = instance.GetNamespace()

	if instance.Status.TargetClusterRef != nil && !instance.Status.TargetClusterRef.IsEmpty() {
		target.Name = instance.Status.TargetClusterRef.Name
		target.Namespace = instance.Status.TargetClusterRef.Namespace
	}

	_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), target, func() error {
		return c.mutateTargetClusterTarget(ctx, target, instance)
	})

	if err != nil {
		return fmt.Errorf("unable to create/update target: %w", err)
	}

	if instance.Status.TargetClusterRef == nil || !instance.Status.TargetClusterRef.IsObject(target) {
		instance.Status.TargetClusterRef = &lssv1alpha2.ObjectReference{
			Name:      target.GetName(),
			Namespace: target.GetNamespace(),
		}

		if err := c.Client().Status().Update(ctx, instance); err != nil {
			return fmt.Errorf("unable to update target cluster target reference for instance: %w", err)
		}
	}

	return nil
}

// mutateTargetClusterTarget creates or updates the target for an instance.
func (c *Controller) mutateTargetClusterTarget(ctx context.Context, target *lsv1alpha1.Target, instance *lssv1alpha2.Instance) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "mutateTargetClusterTarget")

	if len(target.Name) > 0 {
		logger.Info("Updating target cluster target", lc.KeyResource, client.ObjectKeyFromObject(target).String())
	} else {
		logger.Info("Creating target cluster target", lc.KeyResource, types.NamespacedName{Name: target.GenerateName, Namespace: target.Namespace}.String())
	}

	if err := controllerutil.SetControllerReference(instance, target, c.Scheme()); err != nil {
		return fmt.Errorf("unable to set controller reference for target: %w", err)
	}

	config := &lssv1alpha2.ServiceTargetConfig{}
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

// reconcileDataPlaneClusterTarget reconciles the target for the gardener service account.
func (c *Controller) reconcileDataPlaneClusterTarget(ctx context.Context, instance *lssv1alpha2.Instance) error {
	target := &lsv1alpha1.Target{}
	target.GenerateName = fmt.Sprintf("%s-data-plane-", instance.GetName())
	target.Namespace = instance.GetNamespace()

	if instance.Status.DataPlaneClusterRef != nil && !instance.Status.DataPlaneClusterRef.IsEmpty() {
		target.Name = instance.Status.DataPlaneClusterRef.Name
		target.Namespace = instance.Status.DataPlaneClusterRef.Namespace
	}

	_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), target, func() error {
		return c.mutateDataPlaneClusterTarget(ctx, target, instance)
	})

	if instance.Status.DataPlaneClusterRef == nil || !instance.Status.DataPlaneClusterRef.IsObject(target) {
		instance.Status.DataPlaneClusterRef = &lssv1alpha2.ObjectReference{
			Name:      target.GetName(),
			Namespace: target.GetNamespace(),
		}

		if err := c.Client().Status().Update(ctx, instance); err != nil {
			return fmt.Errorf("unable to update data plane cluster target reference for instance: %w", err)
		}
	}

	return err
}

// mutateDataPlaneClusterTarget creates or updates the target for the gardener service account.
func (c *Controller) mutateDataPlaneClusterTarget(ctx context.Context, target *lsv1alpha1.Target, instance *lssv1alpha2.Instance) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "mutateDataPlaneClusterTarget")

	if len(target.Name) > 0 {
		logger.Info("Updating data plane cluster target", lc.KeyResource, client.ObjectKeyFromObject(target).String())
	} else {
		logger.Info("Creating data plane cluster target", lc.KeyResource, types.NamespacedName{Name: target.GenerateName, Namespace: target.Namespace}.String())
	}

	if err := controllerutil.SetControllerReference(instance, target, c.Scheme()); err != nil {
		return fmt.Errorf("unable to set controller reference for target: %w", err)
	}

	if len(instance.Spec.DataPlane.Kubeconfig) > 0 {

	} else {

	}

	saConfig := c.Config().GardenerConfiguration.ServiceAccountKubeconfig

	secret := &corev1.Secret{}
	if err := c.Client().Get(ctx, saConfig.NamespacedName(), secret); err != nil {
		return fmt.Errorf("unable to get kubeconfig secret for gardener service account: %w", err)
	}

	kubeconfig, ok := secret.Data[saConfig.Key]
	if !ok {
		return fmt.Errorf("unable to read kubeconfig from secret: missing key %q", saConfig.Key)
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
func (c *Controller) reconcileInstallation(ctx context.Context, instance *lssv1alpha2.Instance) error {
	old := instance.DeepCopy()

	if !reflect.DeepEqual(old.Status, instance.Status) {
		if err := c.Client().Status().Update(ctx, instance); err != nil {
			return fmt.Errorf("unable to update instance status: %w", err)
		}
	}

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

	old = instance.DeepCopy()
	instance.Status.InstallationRef = &lssv1alpha2.ObjectReference{
		Name:      installation.GetName(),
		Namespace: installation.GetNamespace(),
	}

	instance.Status.LandscaperServiceComponent = &lssv1alpha2.LandscaperServiceComponent{
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
func (c *Controller) mutateInstallation(ctx context.Context, installation *lsv1alpha1.Installation, instance *lssv1alpha2.Instance) error {
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

	config := &lssv1alpha2.ServiceTargetConfig{}
	if err := c.Client().Get(ctx, instance.Spec.ServiceTargetConfigRef.NamespacedName(), config); err != nil {
		return fmt.Errorf("unable to get service target config for instance: %w", err)
	}

	registryConfig := lsinstallation.NewRegistryConfig()
	registryConfigRaw, err := registryConfig.ToAnyJSON()
	if err != nil {
		return fmt.Errorf("unable to marshal registry config: %w", err)
	}

	landscaperConfig := lsinstallation.NewLandscaperConfig()
	landscaperConfig.Resources = instance.Spec.LandscaperConfiguration.Resources
	landscaperConfig.ResourcesMain = instance.Spec.LandscaperConfiguration.ResourcesMain
	landscaperConfig.HPAMain = instance.Spec.LandscaperConfiguration.HPAMain
	landscaperConfig.Deployers = instance.Spec.LandscaperConfiguration.Deployers
	landscaperConfig.DeployersConfig = instance.Spec.LandscaperConfiguration.DeployersConfig
	landscaperConfig.Landscaper.Verbosity = logging.INFO.String()
	if instance.Spec.LandscaperConfiguration.Landscaper != nil {
		landscaperConfig.Landscaper.Controllers = instance.Spec.LandscaperConfiguration.Landscaper.Controllers
		landscaperConfig.Landscaper.DeployItemTimeouts = instance.Spec.LandscaperConfiguration.Landscaper.DeployItemTimeouts
		landscaperConfig.Landscaper.K8SClientSettings = instance.Spec.LandscaperConfiguration.Landscaper.K8SClientSettings
	}

	landscaperConfigRaw, err := landscaperConfig.ToAnyJSON()
	if err != nil {
		return fmt.Errorf("unable to marshal landscaper config: %w", err)
	}

	sidecarConfig := lsinstallation.NewSidecarConfig()
	sidecarConfigRaw, err := sidecarConfig.ToAnyJSON()
	if err != nil {
		return fmt.Errorf("unable to marshal sidecar config: %w", err)
	}

	rotationConfig := lsinstallation.NewRotationConfig(tokenExpirationSeconds, adminKubeconfigExpirationSeconds)
	rotationConfigRaw, err := rotationConfig.ToAnyJSON()
	if err != nil {
		return fmt.Errorf("unable to marshal rotation config: %w", err)
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
					Name:   lsinstallation.TargetClusterImportName,
					Target: instance.Status.TargetClusterRef.Name,
				},
				{
					Name:   lsinstallation.DataPlaneClusterNamespace,
					Target: instance.Status.GardenerServiceAccountRef.Name,
				},
			},
		},
		ImportDataMappings: map[string]lsv1alpha1.AnyJSON{
			lsinstallation.TargetClusterNamespaceImportName:    utils.StringToAnyJSON(fmt.Sprintf("%s-%s", instance.Spec.TenantId, instance.Spec.ID)),
			lsinstallation.DataPlaneClusterNamespaceImportName: utils.StringToAnyJSON(lsinstallation.DataPlaneClusterNamespace),
			lsinstallation.RegistryConfigImportName:            *registryConfigRaw,
			lsinstallation.LandscaperConfigImportName:          *landscaperConfigRaw,
			lsinstallation.SidecarConfigImportName:             *sidecarConfigRaw,
			lsinstallation.RotationConfigImportName:            *rotationConfigRaw,
			lsinstallation.WebhooksHostNameImportName:          utils.StringToAnyJSON(fmt.Sprintf("%s-%s.%s", instance.Spec.TenantId, instance.Spec.ID, config.Spec.IngressDomain)),
		},
		Exports: lsv1alpha1.InstallationExports{
			Data: []lsv1alpha1.DataExport{
				{
					Name:    lsinstallation.AdminKubeconfigExportName,
					DataRef: lsinstallation.GetInstallationExportDataRef(instance, lsinstallation.AdminKubeconfigExportName),
				},
			},
		},
		AutomaticReconcile: &lsv1alpha1.AutomaticReconcile{
			SucceededReconcile: &lsv1alpha1.SucceededReconcile{
				Interval: &lsv1alpha1.Duration{Duration: time.Duration(automaticReconcileSeconds) * time.Second},
			},
			FailedReconcile: &lsv1alpha1.FailedReconcile{
				Interval: &lsv1alpha1.Duration{Duration: time.Duration(failedReconcileSeconds) * time.Second},
			},
		},
	}

	if !InstallationSpecDeepEquals(oldInstallationSpec, installation.Spec.DeepCopy()) {
		// set reconcile annotation to start/update the installation
		logger.Info("Setting reconcile operation annotation")
		if installation.Annotations == nil {
			installation.Annotations = make(map[string]string)
		}
		installation.Annotations[lsv1alpha1.OperationAnnotation] = string(lsv1alpha1.ReconcileOperation)
	}

	return nil
}

// handleExports tries to find the exports of the installation and update the instance status accordingly.
func (c *Controller) handleExports(ctx context.Context, instance *lssv1alpha2.Instance, installation *lsv1alpha1.Installation) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "handleExports")

	dataObjects := &lsv1alpha1.DataObjectList{}

	labelSelector := labels.SelectorFromSet(map[string]string{
		lsv1alpha1.DataObjectSourceLabel:     fmt.Sprintf("Inst.%s", installation.GetName()),
		lsv1alpha1.DataObjectSourceTypeLabel: string(lsv1alpha1.ExportDataObjectSourceType),
	})

	listOptions := client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     installation.GetNamespace(),
	}

	if err := c.Client().List(ctx, dataObjects, &listOptions); err != nil {
		return fmt.Errorf("unable to list data objects for ClusterKubeconfig: %w", err)
	}

	if len(dataObjects.Items) > 0 {
		adminKubeconfigExportName := lsinstallation.GetInstallationExportDataRef(instance, lsinstallation.AdminKubeconfigExportName)

		for _, do := range dataObjects.Items {
			key, ok := do.Labels[lsv1alpha1.DataObjectKeyLabel]
			if !ok {
				continue
			}

			switch key {
			case adminKubeconfigExportName:
				logger.Info("found export data object for user kubeconfig",
					lc.KeyResource, types.NamespacedName{Name: do.Name, Namespace: do.Namespace}.String())
				if err := json.Unmarshal(do.Data.RawMessage, &instance.Status.AdminKubeconfig); err != nil {
					return fmt.Errorf("unable to unmarshal admin kubeconfig: %w", err)
				}
			default:
				continue
			}
		}
	}

	return nil
}
