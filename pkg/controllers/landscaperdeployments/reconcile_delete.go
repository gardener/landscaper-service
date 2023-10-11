// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployments

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"

	"github.com/gardener/landscaper-service/pkg/apis/constants"
	lsserrors "github.com/gardener/landscaper-service/pkg/apis/provisioning/errors"
	provisioningv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/provisioning/v1alpha2"
)

// handleDelete handles the deletion of a landscaper deployment.
func (c *Controller) handleDelete(ctx context.Context, deployment *provisioningv1alpha2.LandscaperDeployment) error {
	var (
		err             error
		currOp          = "Delete"
		removeFinalizer bool
	)

	if deployment.Status.InstanceRef != nil && !deployment.Status.InstanceRef.IsEmpty() {
		removeFinalizer, err = c.ensureDeleteInstanceForDeployment(ctx, deployment)
		if err != nil {
			return err
		}
	} else {
		removeFinalizer = true
	}

	if removeFinalizer {
		controllerutil.RemoveFinalizer(deployment, constants.LandscaperServiceFinalizer)
		if err = c.Client().Update(ctx, deployment); err != nil {
			return lsserrors.NewWrappedError(err, currOp, "RemoveFinalizer", err.Error())
		}
	}

	return nil
}

// ensureDeleteInstanceForDeployment ensures that the instance referenced by this deployment is deleted.
func (c *Controller) ensureDeleteInstanceForDeployment(ctx context.Context, deployment *provisioningv1alpha2.LandscaperDeployment) (bool, error) {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(deployment).String()},
		lc.KeyMethod, "ensureDeleteInstanceForDeployment")

	logger.Info("Delete instance for landscaper deployment", lc.KeyResource, deployment.Status.InstanceRef.NamespacedName().String())
	instance := &provisioningv1alpha2.Instance{}

	if err := c.Client().Get(ctx, deployment.Status.InstanceRef.NamespacedName(), instance); err != nil {
		if apierrors.IsNotFound(err) {
			deployment.Status.InstanceRef = nil
			if err := c.Client().Status().Update(ctx, deployment); err != nil {
				return false, fmt.Errorf("unable to remove instance reference: %w", err)
			}
			return true, nil
		} else {
			return false, fmt.Errorf("unable to get instance for deployment: %w", err)
		}
	}

	if instance.GetDeletionTimestamp().IsZero() {
		if err := c.Client().Delete(ctx, instance); err != nil {
			return false, fmt.Errorf("unable to delete instance: %w", err)
		}
	}

	return false, nil
}
