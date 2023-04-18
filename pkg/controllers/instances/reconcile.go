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
	"time"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/apis/core/v1alpha1/targettypes"
	"github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lssconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/apis/errors"
	lsinstallation "github.com/gardener/landscaper-service/pkg/apis/installation"
	"github.com/gardener/landscaper-service/pkg/utils"
)

var (
	// ClusterKubeconfigExportLabel is the label for the cluster kubeconfig export data object.
	ClusterKubeconfigExportLabel = fmt.Sprintf("%s=%s", lsv1alpha1.DataObjectKeyLabel, lsinstallation.UserKubeconfigExportName)
	// ClusterEndpointExportLabel is the label for the cluster api export data object.
	ClusterEndpointExportLabel = fmt.Sprintf("%s=%s", lsv1alpha1.DataObjectKeyLabel, lsinstallation.ClusterEndpointExportName)
)

const (
	// shootAPIVersion is the gardener shoot version
	shootAPIVersion = "core.gardener.cloud/v1beta1"
	// shootKind is the gardener shoot kind
	shootKind = "Shoot"
	// automaticReconcileSeconds is the number of seconds after which installations of landscaper instances are
	// automatically reconciled. Important: the value must be shorter than tokenExpirationSeconds
	automaticReconcileSeconds = 14 * 24 * 60 * 60
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
func (c *Controller) reconcile(ctx context.Context, instance *lssv1alpha1.Instance) error {
	currOp := "Reconcile"

	if err := c.reconcileContext(ctx, instance); err != nil {
		return errors.NewWrappedError(err, currOp, "ReconcileContextFailed", err.Error())
	}

	if err := c.reconcileGardenerServiceAccountTarget(ctx, instance); err != nil {
		return errors.NewWrappedError(err, currOp, "ReconcileGardenerServiceAccountTargetFailed", err.Error())
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

// reconcileGardenerServiceAccountTarget reconciles the target for the gardener service account.
func (c *Controller) reconcileGardenerServiceAccountTarget(ctx context.Context, instance *lssv1alpha1.Instance) error {
	target := &lsv1alpha1.Target{}
	target.GenerateName = fmt.Sprintf("%s-gardener-sa-", instance.GetName())
	target.Namespace = instance.GetNamespace()

	if instance.Status.GardenerServiceAccountRef != nil && !instance.Status.GardenerServiceAccountRef.IsEmpty() {
		target.Name = instance.Status.GardenerServiceAccountRef.Name
		target.Namespace = instance.Status.GardenerServiceAccountRef.Namespace
	}

	_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), target, func() error {
		return c.mutateGardenerServiceAccountTarget(ctx, target, instance)
	})

	if instance.Status.GardenerServiceAccountRef == nil || !instance.Status.GardenerServiceAccountRef.IsObject(target) {
		instance.Status.GardenerServiceAccountRef = &lssv1alpha1.ObjectReference{
			Name:      target.GetName(),
			Namespace: target.GetNamespace(),
		}

		if err := c.Client().Status().Update(ctx, instance); err != nil {
			return fmt.Errorf("unable to gardener service account target reference for instance: %w", err)
		}
	}

	return err
}

// mutateGardenerServiceAccountTarget creates or updates the target for the gardener service account.
func (c *Controller) mutateGardenerServiceAccountTarget(ctx context.Context, target *lsv1alpha1.Target, instance *lssv1alpha1.Instance) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "mutateGardenerServiceAccountTarget")

	if len(target.Name) > 0 {
		logger.Info("Updating target", lc.KeyResource, client.ObjectKeyFromObject(target).String())
	} else {
		logger.Info("Creating target", lc.KeyResource, types.NamespacedName{Name: target.GenerateName, Namespace: target.Namespace}.String())
	}

	if err := controllerutil.SetControllerReference(instance, target, c.Scheme()); err != nil {
		return fmt.Errorf("unable to set controller reference for target: %w", err)
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
func (c *Controller) reconcileInstallation(ctx context.Context, instance *lssv1alpha1.Instance) error {
	old := instance.DeepCopy()
	if err := c.handleShootName(ctx, instance); err != nil {
		return err
	}

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

	shootConfig := &lssconfig.ShootConfiguration{}
	c.Config().ShootConfiguration.DeepCopyInto(shootConfig)

	if instance.Spec.OIDCConfig != nil {
		// Overwrite the default oidc config with instance specific values
		shootConfig.Kubernetes.KubeAPIServer.OIDCConfig.ClientID = instance.Spec.OIDCConfig.ClientID
		shootConfig.Kubernetes.KubeAPIServer.OIDCConfig.IssuerURL = instance.Spec.OIDCConfig.IssuerURL
		shootConfig.Kubernetes.KubeAPIServer.OIDCConfig.UsernameClaim = instance.Spec.OIDCConfig.UsernameClaim
		shootConfig.Kubernetes.KubeAPIServer.OIDCConfig.GroupsClaim = instance.Spec.OIDCConfig.GroupsClaim
	}

	shootConfigRaw, err := json.Marshal(shootConfig)
	if err != nil {
		return fmt.Errorf("unable to marshal shoot config: %w", err)
	}

	shootLabels := map[string]string{
		lssv1alpha1.ShootTenantIDLabel:          instance.Spec.TenantId,
		lssv1alpha1.ShootInstanceNameLabel:      instance.Name,
		lssv1alpha1.ShootInstanceNamespaceLabel: instance.Namespace,
		lssv1alpha1.ShootInstanceIDLabel:        instance.Spec.ID,
	}
	shootLabelsRaw, err := json.Marshal(shootLabels)
	if err != nil {
		return fmt.Errorf("unable to marshal shoot labels: %w", err)
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
				{
					Name:   "gardenerServiceAccount",
					Target: fmt.Sprintf("#%s", instance.Status.GardenerServiceAccountRef.Name),
				},
			},
		},
		ImportDataMappings: map[string]lsv1alpha1.AnyJSON{
			lsinstallation.HostingClusterNamespaceImportName: utils.StringToAnyJSON(fmt.Sprintf("%s-%s", instance.Spec.TenantId, instance.Spec.ID)),
			lsinstallation.TargetClusterNamespaceImportName:  utils.StringToAnyJSON(lsinstallation.TargetClusterNamespace),
			lsinstallation.RegistryConfigImportName:          *registryConfigRaw,
			lsinstallation.LandscaperConfigImportName:        *landscaperConfigRaw,
			lsinstallation.SidecarConfigImportName:           *sidecarConfigRaw,
			lsinstallation.RotationConfigImportName:          *rotationConfigRaw,
			lsinstallation.ShootNameImportName:               utils.StringToAnyJSON(instance.Status.ShootName),
			lsinstallation.ShootNamespaceImportName:          utils.StringToAnyJSON(instance.Status.ShootNamespace),
			lsinstallation.ShootSecretBindingImportName:      utils.StringToAnyJSON(c.Config().GardenerConfiguration.ShootSecretBindingName),
			lsinstallation.ShootLabelsImportName:             lsv1alpha1.NewAnyJSON(shootLabelsRaw),
			lsinstallation.ShootConfigImportName:             lsv1alpha1.NewAnyJSON(shootConfigRaw),
			lsinstallation.WebhooksHostNameImportName:        utils.StringToAnyJSON(fmt.Sprintf("%s-%s.%s", instance.Spec.TenantId, instance.Spec.ID, config.Spec.IngressDomain)),
		},
		Exports: lsv1alpha1.InstallationExports{
			Data: []lsv1alpha1.DataExport{
				{
					Name:    lsinstallation.ClusterEndpointExportName,
					DataRef: lsinstallation.ClusterEndpointExportName,
				},
				{
					Name:    lsinstallation.UserKubeconfigExportName,
					DataRef: lsinstallation.UserKubeconfigExportName,
				},
				{
					Name:    lsinstallation.AdminKubeconfigExportName,
					DataRef: lsinstallation.AdminKubeconfigExportName,
				},
			},
		},
		AutomaticReconcile: &lsv1alpha1.AutomaticReconcile{
			SucceededReconcile: &lsv1alpha1.SucceededReconcile{
				Interval: &lsv1alpha1.Duration{Duration: time.Duration(automaticReconcileSeconds) * time.Second},
			},
		},
	}

	if c.Config().AuditLogConfig != nil {
		logger.Info("Setting audit log configuration")

		auditPolicyCm := &corev1.ConfigMap{}
		if err := c.Client().Get(ctx, c.Config().AuditLogConfig.AuditPolicy.NamespacedName(), auditPolicyCm); err != nil {
			return fmt.Errorf("unable to retrieve audit policy from config map %q: %w", c.Config().AuditLogConfig.AuditPolicy.NamespacedName(), err)
		}

		auditPolicyRawStr, ok := auditPolicyCm.Data[c.Config().AuditLogConfig.AuditPolicy.Key]
		if !ok {
			return fmt.Errorf("audit policy config map has no key %q", c.Config().AuditLogConfig.AuditPolicy.Key)
		}

		var auditPolicy map[string]interface{}
		decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(auditPolicyRawStr), 512)
		if err := decoder.Decode(&auditPolicy); err != nil {
			return fmt.Errorf("failed to decode audit policy: %w", err)
		}

		auditPolicyRaw, err := json.Marshal(auditPolicy)
		if err != nil {
			return fmt.Errorf("failed to marshal audit policy: %w", err)
		}

		installation.Spec.ImportDataMappings[lsinstallation.SubaccountIdImportName] = utils.StringToAnyJSON(c.Config().AuditLogConfig.SubAccountId)
		installation.Spec.ImportDataMappings[lsinstallation.AuditPolicyImportName] = lsv1alpha1.NewAnyJSON(auditPolicyRaw)
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
func (c *Controller) handleExports(ctx context.Context, instance *lssv1alpha1.Instance, installation *lsv1alpha1.Installation) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "handleExports")

	dataObjects := &lsv1alpha1.DataObjectList{}
	selectorBuilder := strings.Builder{}

	dataObjectSource := fmt.Sprintf("%s=Inst.%s,", lsv1alpha1.DataObjectSourceLabel, installation.GetName())
	dataObjectType := fmt.Sprintf("%s=%s", lsv1alpha1.DataObjectSourceTypeLabel, lsv1alpha1.ExportDataObjectSourceType)
	selectorBuilder.WriteString(dataObjectSource)
	selectorBuilder.WriteString(dataObjectType)

	labelSelector, _ := labels.Parse(selectorBuilder.String())
	listOptions := client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     installation.GetNamespace(),
	}

	if err := c.Client().List(ctx, dataObjects, &listOptions); err != nil {
		return fmt.Errorf("unable to list data objects for ClusterKubeconfig: %w", err)
	}

	if len(dataObjects.Items) > 0 {
		for _, do := range dataObjects.Items {
			key, ok := do.Labels[lsv1alpha1.DataObjectKeyLabel]
			if !ok {
				continue
			}

			switch key {
			case lsinstallation.UserKubeconfigExportName:
				logger.Info("found export data object for user kubeconfig",
					lc.KeyResource, types.NamespacedName{Name: do.Name, Namespace: do.Namespace}.String())
				if err := json.Unmarshal(do.Data.RawMessage, &instance.Status.UserKubeconfig); err != nil {
					return fmt.Errorf("unable to unmarshal user kubeconfig: %w", err)
				}
			case lsinstallation.AdminKubeconfigExportName:
				logger.Info("found export data object for user kubeconfig",
					lc.KeyResource, types.NamespacedName{Name: do.Name, Namespace: do.Namespace}.String())
				if err := json.Unmarshal(do.Data.RawMessage, &instance.Status.AdminKubeconfig); err != nil {
					return fmt.Errorf("unable to unmarshal admin kubeconfig: %w", err)
				}
			case lsinstallation.ClusterEndpointExportName:
				logger.Info("found export data object for cluster endpoint",
					lc.KeyResource, types.NamespacedName{Name: do.Name, Namespace: do.Namespace}.String())
				if err := json.Unmarshal(do.Data.RawMessage, &instance.Status.ClusterEndpoint); err != nil {
					return fmt.Errorf("unable to unmarshal ClusterEndpoint: %w", err)
				}
			default:
				continue
			}
		}
	}

	return nil
}

// handleShootName tries to generate a shoot name if it not already exists.
func (c *Controller) handleShootName(ctx context.Context, instance *lssv1alpha1.Instance) error {
	shootNamespace := fmt.Sprintf("garden-%s", c.Config().GardenerConfiguration.ProjectName)
	instance.Status.ShootNamespace = shootNamespace

	if len(instance.Status.ShootName) > 0 {
		return nil
	}

	shootList, err := c.ListShootsFunc(ctx, instance)
	if err != nil {
		return err
	}

	existingShoots := sets.NewString()
	for _, i := range shootList.Items {
		existingShoots.Insert(i.GetName())
	}

	var shootName string
	for shootName = c.NewUniqueID(); existingShoots.Has(shootName); shootName = c.NewUniqueID() {
	}

	instance.Status.ShootName = shootName

	return nil
}

// listShoots lists all shoot resources for the instances shoot namespace.
func (c *Controller) listShoots(ctx context.Context, instance *lssv1alpha1.Instance) (*unstructured.UnstructuredList, error) {
	shootList := &unstructured.UnstructuredList{
		Object: map[string]interface{}{
			"apiVersion": shootAPIVersion,
			"kind":       shootKind,
		},
	}

	saClient, err := c.getGardenerServiceAccountClient(ctx)
	if err != nil {
		return nil, err
	}

	if err := saClient.List(ctx, shootList, &client.ListOptions{Namespace: instance.Status.ShootNamespace}); err != nil {
		return nil, fmt.Errorf("failed to list shoots: %w", err)
	}

	return shootList, nil
}

// getGardenerServiceAccountClient retrieves and initializes the gardener service account client, if necessary.
func (c *Controller) getGardenerServiceAccountClient(ctx context.Context) (client.Client, error) {
	if c.gardenerServiceAccountClient != nil {
		return c.gardenerServiceAccountClient, nil
	}

	gardenerServiceAccountSecret := &corev1.Secret{}
	if err := c.Client().Get(ctx, c.Config().GardenerConfiguration.ServiceAccountKubeconfig.NamespacedName(), gardenerServiceAccountSecret); err != nil {
		return nil, fmt.Errorf("failed to load gardener service account secret: %w", err)
	}

	key := c.Config().GardenerConfiguration.ServiceAccountKubeconfig.Key
	kubeconfig, ok := gardenerServiceAccountSecret.Data[key]
	if !ok {
		return nil, fmt.Errorf("gardener service account secret has no key %q", key)
	}

	clientCfg, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load gardener service account kubeconfig: %w", err)
	}

	loader := clientcmd.NewDefaultClientConfig(*clientCfg, nil)
	restConfig, err := loader.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load gardener service account rest config: %w", err)
	}

	saClient, err := client.New(restConfig, client.Options{
		Scheme: c.Scheme(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client for gardener service account: %w", err)
	}

	c.gardenerServiceAccountClient = saClient
	return saClient, nil
}
