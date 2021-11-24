// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"context"
	"fmt"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lsserrors "github.com/gardener/landscaper-service/pkg/apis/errors"
	"github.com/gardener/landscaper-service/pkg/utils"
)

// handleDelete handles the deletion of instances
func (c *Controller) handleDelete(ctx context.Context, log logr.Logger, instance *lssv1alpha1.Instance) error {
	curOp := "Delete"

	if instance.Status.InstallationRef != nil && !instance.Status.InstallationRef.IsEmpty() {
		if err := c.ensureDeleteInstallationForInstance(ctx, log, instance); err != nil {
			return lsserrors.NewWrappedError(err, curOp, "DeleteInstallation", err.Error())
		}
		return nil
	}

	if instance.Status.TargetRef != nil && !instance.Status.TargetRef.IsEmpty() {
		if err := c.ensureDeleteTargetForInstance(ctx, log, instance); err != nil {
			return lsserrors.NewWrappedError(err, curOp, "DeleteTarget", err.Error())
		}
		return nil
	}

	serviceTargetConfig := &lssv1alpha1.ServiceTargetConfig{}
	if err := c.Client().Get(ctx, instance.Spec.ServiceTargetConfigRef.NamespacedName(), serviceTargetConfig); err != nil {
		return lsserrors.NewWrappedError(err, curOp, "GetServiceTargetConfig", err.Error())
	}

	// remove instance reference from service target config
	serviceTargetConfig.Status.InstanceRefs = utils.RemoveReference(serviceTargetConfig.Status.InstanceRefs, &lssv1alpha1.ObjectReference{
		Name:      instance.GetName(),
		Namespace: instance.GetNamespace(),
	})

	if err := c.Client().Status().Update(ctx, serviceTargetConfig); err != nil {
		return lsserrors.NewWrappedError(err, curOp, "RemoveRefFromServiceTargetConfig", err.Error())
	}

	controllerutil.RemoveFinalizer(instance, lssv1alpha1.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, instance); err != nil {
		return lsserrors.NewWrappedError(err, curOp, "RemoveFinalizer", err.Error())
	}

	return nil
}

// ensureDeleteInstallationForInstance ensures that the installation for an instance is deleted
func (c *Controller) ensureDeleteInstallationForInstance(ctx context.Context, log logr.Logger, instance *lssv1alpha1.Instance) error {
	log.Info("Delete installation for instance")
	installation := &lsv1alpha1.Installation{}

	if err := c.Client().Get(ctx, instance.Status.InstallationRef.NamespacedName(), installation); err != nil {
		if apierrors.IsNotFound(err) {
			instance.Status.InstallationRef = nil
			if err := c.Client().Status().Update(ctx, instance); err != nil {
				return fmt.Errorf("failed to remove installation reference: %w", err)
			}
		} else {
			return fmt.Errorf("unable to get installation for instance: %w", err)
		}
	}

	if installation.DeletionTimestamp.IsZero() {
		if err := c.Client().Delete(ctx, installation); err != nil {
			return fmt.Errorf("unable to delete installation for instance: %w", err)
		}
	}

	return nil
}

// ensureDeleteTargetForInstance ensures that the target for an instance is deleted
func (c *Controller) ensureDeleteTargetForInstance(ctx context.Context, log logr.Logger, instance *lssv1alpha1.Instance) error {
	log.Info("Delete target for instance")
	target := &lsv1alpha1.Target{}

	if err := c.Client().Get(ctx, instance.Status.TargetRef.NamespacedName(), target); err != nil {
		if apierrors.IsNotFound(err) {
			instance.Status.TargetRef = nil
			if err := c.Client().Status().Update(ctx, instance); err != nil {
				return fmt.Errorf("failed to remove target reference: %w", err)
			}
		} else {
			return fmt.Errorf("unable to get target for instance: %w", err)
		}
	}

	if target.DeletionTimestamp.IsZero() {
		if err := c.Client().Delete(ctx, target); err != nil {
			return fmt.Errorf("unable to delete target for instance: %w", err)
		}
	}

	return nil
}
