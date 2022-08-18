// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"context"
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lsserrors "github.com/gardener/landscaper-service/pkg/apis/errors"
	"github.com/gardener/landscaper-service/pkg/operation"
)

// Controller is the instances controller
type Controller struct {
	operation.Operation
	log logging.Logger
}

// NewController returns a new instances controller
func NewController(logger logging.Logger, c client.Client, scheme *runtime.Scheme, config *coreconfig.LandscaperServiceConfiguration) (reconcile.Reconciler, error) {
	ctrl := &Controller{
		log: logger,
	}
	op := operation.NewOperation(c, scheme, config)
	ctrl.Operation = *op
	return ctrl, nil
}

// NewTestActuator creates a new controller for testing purposes.
func NewTestActuator(op operation.Operation, logger logging.Logger) *Controller {
	ctrl := &Controller{
		Operation: op,
		log:       logger,
	}
	return ctrl
}

// Reconcile reconciles requests for instances
func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger, ctx := c.log.StartReconcileAndAddToContext(ctx, req)

	instance := &lssv1alpha1.Instance{}
	if err := c.Client().Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info(err.Error())
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	c.Operation.Scheme().Default(instance)
	errHdl := c.handleErrorFunc(instance)

	// update observed generation
	if instance.Status.ObservedGeneration < instance.GetGeneration() {
		instance.Status.ObservedGeneration = instance.GetGeneration()
		if err := c.Client().Status().Update(ctx, instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	// set finalizer
	if instance.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(instance, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(instance, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, instance); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	// reconcile delete
	if !instance.DeletionTimestamp.IsZero() {
		return reconcile.Result{}, errHdl(ctx, c.handleDelete(ctx, instance))
	}

	// reconcile
	return reconcile.Result{}, errHdl(ctx, c.reconcile(ctx, instance))
}

// handleErrorFunc updates the error status of an instance
func (c *Controller) handleErrorFunc(instance *lssv1alpha1.Instance) func(ctx context.Context, err error) error {
	old := instance.DeepCopy()
	return func(ctx context.Context, err error) error {
		logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(instance).String()})
		instance.Status.LastError = lsserrors.TryUpdateError(instance.Status.LastError, err)

		if !reflect.DeepEqual(old.Status, instance.Status) {
			if err2 := c.Client().Status().Update(ctx, instance); err2 != nil {
				if apierrors.IsConflict(err2) {
					// reduce logging
					logger.Info(fmt.Sprintf("unable to update status: %s", err2.Error()))
				} else {
					logger.Error(err2, "unable to update status")
				}

				// retry on conflict
				if err != nil {
					return err2
				}
			}
		}
		return err
	}
}
