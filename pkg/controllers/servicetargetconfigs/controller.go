// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package servicetargetconfigs

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	config "github.com/gardener/landscaper-service/pkg/apis/config/v1alpha1"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/operation"
)

// Controller is the servicetargetconfig controller
type Controller struct {
	operation.Operation
	log logging.Logger
}

// NewController returns a new servicetargetconfig controller
func NewController(logger logging.Logger, c client.Client, scheme *runtime.Scheme, config *config.LandscaperServiceConfiguration) (reconcile.Reconciler, error) {
	ctrl := &Controller{
		log: logger,
	}
	op := operation.NewOperation(c, scheme, config)
	ctrl.Operation = *op
	return ctrl, nil
}

// Reconcile reconciles requests for servicetargetconfigs
func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger, ctx := c.log.StartReconcileAndAddToContext(ctx, req)

	config := &lssv1alpha1.ServiceTargetConfig{}
	if err := c.Client().Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info(err.Error())
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	c.Operation.Scheme().Default(config)

	// update observed generation
	if config.Status.ObservedGeneration < config.GetGeneration() {
		config.Status.ObservedGeneration = config.GetGeneration()
		if err := c.Client().Status().Update(ctx, config); err != nil {
			return reconcile.Result{}, err
		}
	}

	// set finalizer
	if config.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(config, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(config, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, config); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	if !config.DeletionTimestamp.IsZero() {
		// TODO: handle delete
		controllerutil.RemoveFinalizer(config, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, config); err != nil {
			return reconcile.Result{}, err
		}

	}

	return reconcile.Result{}, c.reconcile(ctx, config)
}
