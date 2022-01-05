// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package servicetargetconfigs

import (
	"context"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper-service/pkg/operation"
)

// Controller is the servicetargetconfig controller
type Controller struct {
	operation.Operation
}

// NewController returns a new servicetargetconfig controller
func NewController(log logr.Logger, c client.Client, scheme *runtime.Scheme) (reconcile.Reconciler, error) {
	ctrl := &Controller{}
	op := operation.NewOperation(log, c, scheme)
	ctrl.Operation = *op
	return ctrl, nil
}

// Reconcile reconciles requests for servicetargetconfigs
func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := c.Log().WithValues("servicetargetconfig", req.NamespacedName.String())
	ctx = logr.NewContext(ctx, log)
	log.V(5).Info("reconcile", "resource", req.NamespacedName)

	config := &lssv1alpha1.ServiceTargetConfig{}
	if err := c.Client().Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			c.Log().V(5).Info(err.Error())
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

	return reconcile.Result{}, c.reconcile(ctx, log, config)
}
