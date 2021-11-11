// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployments

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lsserrors "github.com/gardener/landscaper-service/pkg/apis/errors"
)

// handleDelete handles the deletion of a landscaper deployment.
func (c *Controller) handleDelete(ctx context.Context, log logr.Logger, deployment *lssv1alpha1.LandscaperDeployment) error {
	currOp := "Delete"
	if deployment.Status.InstanceRef != nil {
		return c.ensureDeleteInstanceForDeployment(ctx, log, deployment)
	}

	controllerutil.RemoveFinalizer(deployment, lssv1alpha1.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, deployment); err != nil {
		return lsserrors.NewWrappedError(err, currOp, "RemoveFinalizer", err.Error())
	}

	return nil
}

// ensureDeleteInstanceForDeployment ensures that the instance referenced by this deployment is deleted.
func (c *Controller) ensureDeleteInstanceForDeployment(ctx context.Context, log logr.Logger, deployment *lssv1alpha1.LandscaperDeployment) error {
	instance := &lssv1alpha1.Instance{}

	if err := c.Client().Get(ctx, deployment.Status.InstanceRef.NamespacedName(), instance); err != nil {
		if apierrors.IsNotFound(err) {
			deployment.Status.InstanceRef = nil
			if err := c.Client().Status().Update(ctx, deployment); err != nil {
				return fmt.Errorf("unable to remove instance reference: %w", err)
			}
		}
		log.Error(err, "unable to get instance for deployment")
		return fmt.Errorf("unable to get instance for deployment")
	}

	if instance.GetDeletionTimestamp().IsZero() {
		if err := c.Client().Delete(ctx, instance); err != nil {
			return fmt.Errorf("unable to delete instance: %w", err)
		}
	}

	return nil
}
