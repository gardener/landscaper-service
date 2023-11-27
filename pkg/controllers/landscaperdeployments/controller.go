// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployments

import (
	"context"
	"fmt"
	"reflect"

	guuid "github.com/google/uuid"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"

	config "github.com/gardener/landscaper-service/pkg/apis/config/v1alpha1"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lsserrors "github.com/gardener/landscaper-service/pkg/apis/errors"
	"github.com/gardener/landscaper-service/pkg/operation"
)

// Controller is the landscaperdeployments controller
type Controller struct {
	operation.Operation
	log logging.Logger

	UniqueIDFunc func() string

	ReconcileFunc    func(ctx context.Context, deployment *lssv1alpha1.LandscaperDeployment) error
	HandleDeleteFunc func(ctx context.Context, deployment *lssv1alpha1.LandscaperDeployment) error
}

// NewUniqueID creates a new unique id string with a length of 8.
func (c *Controller) NewUniqueID() string {
	id := c.UniqueIDFunc()
	if len(id) > 8 {
		id = id[:8]
	}
	return id
}

func defaultUniqueIdFunc() string {
	return guuid.New().String()
}

// NewController returns a new landscaperdeployments controller
func NewController(logger logging.Logger, c client.Client, scheme *runtime.Scheme, config *config.LandscaperServiceConfiguration) (reconcile.Reconciler, error) {
	ctrl := &Controller{
		log:          logger,
		UniqueIDFunc: defaultUniqueIdFunc,
	}
	ctrl.ReconcileFunc = ctrl.reconcile
	ctrl.HandleDeleteFunc = ctrl.handleDelete
	op := operation.NewOperation(c, scheme, config)
	ctrl.Operation = *op
	return ctrl, nil
}

// NewTestActuator creates a new controller for testing purposes.
func NewTestActuator(op operation.Operation, logger logging.Logger) *Controller {
	ctrl := &Controller{
		Operation:    op,
		log:          logger,
		UniqueIDFunc: defaultUniqueIdFunc,
	}
	ctrl.ReconcileFunc = ctrl.reconcile
	ctrl.HandleDeleteFunc = ctrl.handleDelete
	return ctrl
}

// Reconcile reconciles requests for landscaperdeployments
func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger, ctx := c.log.StartReconcileAndAddToContext(ctx, req)

	deployment := &lssv1alpha1.LandscaperDeployment{}
	if err := c.Client().Get(ctx, req.NamespacedName, deployment); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info(err.Error())
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	c.Operation.Scheme().Default(deployment)
	errHdl := c.handleErrorFunc(deployment)

	// update observed generation
	if deployment.Status.ObservedGeneration < deployment.GetGeneration() {
		deployment.Status.ObservedGeneration = deployment.GetGeneration()
		if err := c.Client().Status().Update(ctx, deployment); err != nil {
			return reconcile.Result{}, err
		}
	}

	// set finalizer
	if deployment.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(deployment, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(deployment, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, deployment); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	// reconcile delete
	if !deployment.DeletionTimestamp.IsZero() {
		return reconcile.Result{}, errHdl(ctx, c.HandleDeleteFunc(ctx, deployment))
	}

	// reconcile
	return reconcile.Result{}, errHdl(ctx, c.ReconcileFunc(ctx, deployment))
}

// handleErrorFunc updates the error status of a landscaper deployment
func (c *Controller) handleErrorFunc(deployment *lssv1alpha1.LandscaperDeployment) func(ctx context.Context, err error) error {
	old := deployment.DeepCopy()
	return func(ctx context.Context, err error) error {
		logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(deployment).String()})
		deployment.Status.LastError = lsserrors.TryUpdateError(deployment.Status.LastError, err)

		if !reflect.DeepEqual(old.Status, deployment.Status) {
			if err2 := c.Client().Status().Update(ctx, deployment); err2 != nil {
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
