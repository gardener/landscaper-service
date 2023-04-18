// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lsserrors "github.com/gardener/landscaper-service/pkg/apis/errors"
	"github.com/gardener/landscaper-service/pkg/utils"
)

const (
	targetClusterNamespaceDeletionRetryDuration = time.Second * 10
)

// handleDelete handles the deletion of instances
func (c *Controller) handleDelete(ctx context.Context, instance *lssv1alpha1.Instance) (reconcile.Result, error) {
	var (
		err                                 error
		curOp                               = "Delete"
		targetDeleted                       = true
		gardenerServiceAccountTargetDeleted = true
		installationDeleted                 = true
		contextDeleted                      = true
		targetClusterNamespaceDeleted       bool
	)

	if instance.Status.InstallationRef != nil && !instance.Status.InstallationRef.IsEmpty() {
		if installationDeleted, err = c.ensureDeleteInstallationForInstance(ctx, instance); err != nil {
			return reconcile.Result{}, lsserrors.NewWrappedError(err, curOp, "DeleteInstallation", err.Error())
		}
	}

	if !installationDeleted {
		return reconcile.Result{}, nil
	}

	if targetClusterNamespaceDeleted, err = c.ensureDeleteTargetClusterNamespace(ctx, instance); err != nil {
		return reconcile.Result{}, lsserrors.NewWrappedError(err, curOp, "DeleteTargetClusterNamespace", err.Error())
	}

	if !targetClusterNamespaceDeleted {
		// since this namespace is on a different cluster and there is no owner reference set,
		// the retry has to be triggered manually
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: targetClusterNamespaceDeletionRetryDuration,
		}, nil
	}

	if instance.Status.TargetRef != nil && !instance.Status.TargetRef.IsEmpty() {
		if targetDeleted, err = c.ensureDeleteTargetForInstance(ctx, instance); err != nil {
			return reconcile.Result{}, lsserrors.NewWrappedError(err, curOp, "DeleteTarget", err.Error())
		}
	}

	if !targetDeleted {
		return reconcile.Result{}, nil
	}

	if instance.Status.ContextRef != nil && !instance.Status.ContextRef.IsEmpty() {
		if contextDeleted, err = c.ensureDeleteContextForInstance(ctx, instance); err != nil {
			return reconcile.Result{}, lsserrors.NewWrappedError(err, curOp, "DeleteContext", err.Error())
		}
	}

	if instance.Status.GardenerServiceAccountRef != nil && !instance.Status.GardenerServiceAccountRef.IsEmpty() {
		if gardenerServiceAccountTargetDeleted, err = c.ensureDeleteGardenerServiceAccountTargetForInstance(ctx, instance); err != nil {
			return reconcile.Result{}, lsserrors.NewWrappedError(err, curOp, "DeleteGardenerServiceAccountTarget", err.Error())
		}
	}

	if !gardenerServiceAccountTargetDeleted {
		return reconcile.Result{}, nil
	}

	if !contextDeleted {
		return reconcile.Result{}, nil
	}

	serviceTargetConfig := &lssv1alpha1.ServiceTargetConfig{}
	if err := c.Client().Get(ctx, instance.Spec.ServiceTargetConfigRef.NamespacedName(), serviceTargetConfig); err != nil {
		return reconcile.Result{}, lsserrors.NewWrappedError(err, curOp, "GetServiceTargetConfig", err.Error())
	}

	// remove instance reference from service target config
	serviceTargetConfig.Status.InstanceRefs = utils.RemoveReference(serviceTargetConfig.Status.InstanceRefs, &lssv1alpha1.ObjectReference{
		Name:      instance.GetName(),
		Namespace: instance.GetNamespace(),
	})

	if err := c.Client().Status().Update(ctx, serviceTargetConfig); err != nil {
		return reconcile.Result{}, lsserrors.NewWrappedError(err, curOp, "RemoveRefFromServiceTargetConfig", err.Error())
	}

	controllerutil.RemoveFinalizer(instance, lssv1alpha1.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, instance); err != nil {
		return reconcile.Result{}, lsserrors.NewWrappedError(err, curOp, "RemoveFinalizer", err.Error())
	}

	return reconcile.Result{}, nil
}

// ensureDeleteInstallationForInstance ensures that the installation for an instance is deleted
func (c *Controller) ensureDeleteInstallationForInstance(ctx context.Context, instance *lssv1alpha1.Instance) (bool, error) {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "ensureDeleteInstallationForInstance")

	logger.Info("Delete installation for instance", lc.KeyResource, instance.Status.InstallationRef.NamespacedName())
	installation := &lsv1alpha1.Installation{}

	if err := c.Client().Get(ctx, instance.Status.InstallationRef.NamespacedName(), installation); err != nil {
		if apierrors.IsNotFound(err) {
			instance.Status.InstallationRef = nil
			if err := c.Client().Status().Update(ctx, instance); err != nil {
				return false, fmt.Errorf("failed to remove installation reference: %w", err)
			}
			return true, nil
		} else {
			return false, fmt.Errorf("unable to get installation for instance: %w", err)
		}
	}

	if installation.DeletionTimestamp.IsZero() {
		if err := c.Client().Delete(ctx, installation); err != nil {
			return false, fmt.Errorf("unable to delete installation for instance: %w", err)
		}
	}

	return false, nil
}

// ensureDeleteTargetForInstance ensures that the target for an instance is deleted
func (c *Controller) ensureDeleteTargetForInstance(ctx context.Context, instance *lssv1alpha1.Instance) (bool, error) {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "ensureDeleteTargetForInstance")

	logger.Info("Delete target for instance", lc.KeyResource, instance.Status.TargetRef.NamespacedName())
	target := &lsv1alpha1.Target{}

	if err := c.Client().Get(ctx, instance.Status.TargetRef.NamespacedName(), target); err != nil {
		if apierrors.IsNotFound(err) {
			instance.Status.TargetRef = nil
			if err := c.Client().Status().Update(ctx, instance); err != nil {
				return false, fmt.Errorf("failed to remove target reference: %w", err)
			}
			return true, nil
		} else {
			return false, fmt.Errorf("unable to get target for instance: %w", err)
		}
	}

	if target.DeletionTimestamp.IsZero() {
		if err := c.Client().Delete(ctx, target); err != nil {
			return false, fmt.Errorf("unable to delete target for instance: %w", err)
		}
	}

	return false, nil
}

// ensureDeleteTargetForInstance ensures that the target for an instance is deleted
func (c *Controller) ensureDeleteGardenerServiceAccountTargetForInstance(ctx context.Context, instance *lssv1alpha1.Instance) (bool, error) {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "ensureDeleteGardenerServiceAccountTargetForInstance")

	logger.Info("Delete gardener service account target for instance", lc.KeyResource, instance.Status.GardenerServiceAccountRef.NamespacedName())
	target := &lsv1alpha1.Target{}

	if err := c.Client().Get(ctx, instance.Status.GardenerServiceAccountRef.NamespacedName(), target); err != nil {
		if apierrors.IsNotFound(err) {
			instance.Status.GardenerServiceAccountRef = nil
			if err := c.Client().Status().Update(ctx, instance); err != nil {
				return false, fmt.Errorf("failed to remove gardener service account target reference: %w", err)
			}
			return true, nil
		} else {
			return false, fmt.Errorf("unable to get gardener service account target for instance: %w", err)
		}
	}

	if target.DeletionTimestamp.IsZero() {
		if err := c.Client().Delete(ctx, target); err != nil {
			return false, fmt.Errorf("unable to delete gardener service account target for instance: %w", err)
		}
	}

	return false, nil
}

// ensureDeleteContextForInstance ensures that the context for an instance is deleted
func (c *Controller) ensureDeleteContextForInstance(ctx context.Context, instance *lssv1alpha1.Instance) (bool, error) {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "ensureDeleteContextForInstance")

	logger.Info("Delete context for instance", lc.KeyResource, instance.Status.ContextRef.NamespacedName())
	landscaperContext := &lsv1alpha1.Context{}

	if err := c.Client().Get(ctx, instance.Status.ContextRef.NamespacedName(), landscaperContext); err != nil {
		if apierrors.IsNotFound(err) {
			instance.Status.ContextRef = nil
			if err := c.Client().Status().Update(ctx, instance); err != nil {
				return false, fmt.Errorf("failed to remove context reference: %w", err)
			}
			return true, nil
		} else {
			return false, fmt.Errorf("unable to get context for instance: %w", err)
		}
	}

	if err := c.deleteSecretsForContext(ctx, landscaperContext); err != nil {
		return false, err
	}

	if landscaperContext.DeletionTimestamp.IsZero() {
		if err := c.Client().Delete(ctx, landscaperContext); err != nil {
			return false, fmt.Errorf("unable to delete context for instance: %w", err)
		}
	}

	return false, nil
}

func (c *Controller) deleteSecretsForContext(ctx context.Context, landscaperContext *lsv1alpha1.Context) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	errs := make([]error, 0)

	for _, secretRef := range landscaperContext.RegistryPullSecrets {
		key := types.NamespacedName{Name: secretRef.Name, Namespace: landscaperContext.Namespace}
		logger.Info("Delete secrets for context", lc.KeyResource, key.String())

		secret := &corev1.Secret{}
		if err := c.Client().Get(ctx, key, secret); err != nil {
			if !apierrors.IsNotFound(err) {
				errs = append(errs, fmt.Errorf("unable to get secret \"%s\" for context: %w", secretRef.Name, err))
			}
			continue
		}
		if secret.DeletionTimestamp.IsZero() {
			if err := c.Client().Delete(ctx, secret); err != nil {
				errs = append(errs, fmt.Errorf("unable to delete secret \"%s\" for context: %w", secretRef.Name, err))
			}
		}
	}
	return errors.NewAggregate(errs)
}

// ensureDeleteTargetClusterNamespace ensures that the target cluster namespace for an instance has been deleted.
func (c *Controller) ensureDeleteTargetClusterNamespace(ctx context.Context, instance *lssv1alpha1.Instance) (bool, error) {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "ensureDeleteTargetClusterNamespace")

	if len(instance.Spec.TenantId) == 0 || len(instance.Spec.ID) == 0 {
		return true, nil
	}

	targetClusterNamespace := fmt.Sprintf("%s-%s", instance.Spec.TenantId, instance.Spec.ID)

	logger.Info("Delete target cluster namespace for instance", lc.KeyResourceNonNamespaced, targetClusterNamespace)

	targetClusterClient, err := c.kubeClientExtractor.GetKubeClientFromServiceTargetConfig(
		ctx,
		instance.Spec.ServiceTargetConfigRef.Name,
		instance.Spec.ServiceTargetConfigRef.Namespace,
		c.Client())

	if err != nil {
		return false, fmt.Errorf("failed to get client for target cluster: %w", err)
	}

	deleteTargetClusterRBAC(ctx, instance, targetClusterClient)

	namespace := &corev1.Namespace{}

	if err = targetClusterClient.Get(ctx, types.NamespacedName{Name: targetClusterNamespace}, namespace); err != nil {
		if apierrors.IsNotFound(err) {
			return true, nil
		} else {
			return false, fmt.Errorf("failed to get target cluster namespace %q: %w", targetClusterNamespace, err)
		}
	}

	if namespace.DeletionTimestamp.IsZero() {
		if err = targetClusterClient.Delete(ctx, namespace); err != nil {
			return false, fmt.Errorf("failed to delete target cluster namespace %q: %w", targetClusterNamespace, err)
		}
		return false, nil
	}

	return false, nil
}

// deleteTargetClusterRBAC deletes all clusterroles and bindings associated with this instance.
func deleteTargetClusterRBAC(ctx context.Context, instance *lssv1alpha1.Instance, cl client.Client) {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()},
		lc.KeyMethod, "deleteTargetClusterRBAC")

	targetClusterName := fmt.Sprintf("%s-%s", instance.Spec.TenantId, instance.Spec.ID)

	clusterRoleBindings := &rbacv1.ClusterRoleBindingList{}
	if err := cl.List(ctx, clusterRoleBindings, &client.ListOptions{}); err != nil {
		logger.Error(err, "failed to list target clusterrolebindings")
	}

	for _, crb := range clusterRoleBindings.Items {
		if strings.Contains(crb.GetName(), targetClusterName) &&
			crb.DeletionTimestamp.IsZero() {
			logger.Info("deleting clusterrolebinding", lc.KeyResourceNonNamespaced, crb.GetName())
			if err := cl.Delete(ctx, &crb); err != nil {
				logger.Error(err, "failed to delete clusterrolebinding", lc.KeyResourceNonNamespaced, crb.GetName())
			}
		}
	}

	clusterRoles := &rbacv1.ClusterRoleList{}
	if err := cl.List(ctx, clusterRoles, &client.ListOptions{}); err != nil {
		logger.Error(err, "failed to list clusterroles")
	}

	for _, cr := range clusterRoles.Items {
		if strings.Contains(cr.GetName(), targetClusterName) &&
			cr.DeletionTimestamp.IsZero() {
			logger.Info("deleting clusterrole", lc.KeyResourceNonNamespaced, cr.GetName())
			if err := cl.Delete(ctx, &cr); err != nil {
				logger.Error(err, "failed to delete clusterrole", lc.KeyResourceNonNamespaced, cr.GetName())
			}
		}
	}
}
