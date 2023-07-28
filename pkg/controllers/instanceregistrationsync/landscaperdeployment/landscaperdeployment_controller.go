// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployment

import (
	"context"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	"github.com/gardener/landscaper-service/pkg/controllers/instanceregistrationsync/instanceregistration"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/operation"
)

const LANDSCAPEDEPLOYMENT_REGISTRATION_LABEL string = "landscaper-service.gardener.cloud/landscapedeployment"

type Controller struct {
	operation.Operation
	log logging.Logger

	ReconcileFunc    func(ctx context.Context, landscaperDeployment *lssv1alpha1.LandscaperDeployment) (reconcile.Result, error)
	HandleDeleteFunc func(ctx context.Context, landscaperDeployment *lssv1alpha1.LandscaperDeployment) (reconcile.Result, error)
}

func NewController(logger logging.Logger, c client.Client, scheme *runtime.Scheme) (reconcile.Reconciler, error) {
	ctrl := &Controller{
		log: logger,
	}
	ctrl.ReconcileFunc = ctrl.reconcile
	ctrl.HandleDeleteFunc = ctrl.handleDelete
	op := operation.NewOperation(c, scheme, &coreconfig.LandscaperServiceConfiguration{})
	ctrl.Operation = *op
	return ctrl, nil
}

// NewTestActuator creates a new controller for testing purposes.
func NewTestActuator(op operation.Operation, logger logging.Logger) *Controller {
	ctrl := &Controller{
		Operation: op,
		log:       logger,
	}
	ctrl.ReconcileFunc = ctrl.reconcile
	ctrl.HandleDeleteFunc = ctrl.handleDelete
	return ctrl
}

func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger, ctx := c.log.StartReconcileAndAddToContext(ctx, req)

	landscaperDeployment := &lssv1alpha1.LandscaperDeployment{}
	if err := c.Client().Get(ctx, req.NamespacedName, landscaperDeployment); err != nil {
		logger.Error(err, "failed loading LandscaperDeployment cr")
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// set finalizer
	if landscaperDeployment.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(landscaperDeployment, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(landscaperDeployment, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, landscaperDeployment); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	// reconcile delete
	if !landscaperDeployment.DeletionTimestamp.IsZero() {
		return c.HandleDeleteFunc(ctx, landscaperDeployment)
	}

	return c.reconcile(ctx, landscaperDeployment)
}

func (c *Controller) handleDelete(ctx context.Context, landscaperDeployment *lssv1alpha1.LandscaperDeployment) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)
	controllerutil.RemoveFinalizer(landscaperDeployment, lssv1alpha1.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, landscaperDeployment); err != nil {
		logger.Error(err, "Failed removing finalizer")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil

}

func (c *Controller) reconcile(ctx context.Context, landscaperDeployment *lssv1alpha1.LandscaperDeployment) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	//we filter with the controlelrruntime to only be called on status changes, so we assume it changed when called

	if landscaperDeployment.Status.InstanceRef == nil {
		return reconcile.Result{}, nil //maybe write a pending/progressing information
	}

	//get corresponding instance to landscaperdeployment
	instance := lssv1alpha1.Instance{}
	if err := c.Client().Get(ctx, landscaperDeployment.Status.InstanceRef.NamespacedName(), &instance); err != nil {
		return reconcile.Result{}, err
	}

	//write InstanceRegistration status
	instanceRegistration := &lssv1alpha1.InstanceRegistration{}

	landscaperDeploymentLabels := landscaperDeployment.GetLabels()
	if landscaperDeploymentLabels == nil {
		return reconcile.Result{}, nil
	} else {
		name, ok := landscaperDeploymentLabels[instanceregistration.INSTANCE_REGISTRATION_LABEL_NAME]
		if !ok {
			logger.Info("missing label on landscaperDeployment", "label", instanceregistration.INSTANCE_REGISTRATION_LABEL_NAME)
			return reconcile.Result{}, nil
		}
		namespace, ok := landscaperDeploymentLabels[instanceregistration.INSTANCE_REGISTRATION_LABEL_NAMESPACE]
		if !ok {
			logger.Info("missing label on landscaperDeployment", "label", instanceregistration.INSTANCE_REGISTRATION_LABEL_NAMESPACE)
			return reconcile.Result{}, nil
		}
		instanceRegistration.SetNamespace(namespace)
		instanceRegistration.SetName(name)
	}

	instanceRegistration.Status.LastError = instance.Status.LastError
	instanceRegistration.Status.UserKubeconfig = instance.Status.UserKubeconfig
	if err := c.Client().Status().Update(ctx, instanceRegistration); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
