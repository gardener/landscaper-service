// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instanceregistration

import (
	"context"
	"fmt"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/operation"
)

const INSTANCE_REGISTRATION_LABEL string = "landscaper-service.gardener.cloud/instanceregistration"

type Controller struct {
	operation.Operation
	log logging.Logger

	ReconcileFunc    func(ctx context.Context, instanceRegistration *lssv1alpha1.InstanceRegistration) (reconcile.Result, error)
	HandleDeleteFunc func(ctx context.Context, instanceRegistration *lssv1alpha1.InstanceRegistration) (reconcile.Result, error)
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

	instanceRegistration := &lssv1alpha1.InstanceRegistration{}
	if err := c.Client().Get(ctx, req.NamespacedName, instanceRegistration); err != nil {
		logger.Error(err, "failed loading InstanceRegistration cr")
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// set finalizer
	if instanceRegistration.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(instanceRegistration, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(instanceRegistration, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, instanceRegistration); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	// reconcile delete
	if !instanceRegistration.DeletionTimestamp.IsZero() {
		return c.HandleDeleteFunc(ctx, instanceRegistration)
	}

	return c.reconcile(ctx, instanceRegistration)
}

func (c *Controller) handleDelete(ctx context.Context, instanceRegistration *lssv1alpha1.InstanceRegistration) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)
	controllerutil.RemoveFinalizer(instanceRegistration, lssv1alpha1.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, instanceRegistration); err != nil {
		logger.Error(err, "Failed removing finalizer")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil

}

func (c *Controller) reconcile(ctx context.Context, instanceRegistration *lssv1alpha1.InstanceRegistration) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	// if instanceRegistration.Status.SyncedGeneration != instanceRegistration.Generation {
	// 	logger.Info("SyncedGeneration unequal current generation -> sync is not completed")
	// 	return reconcile.Result{}, nil
	// }

	if instanceRegistration.Status.ObservedGeneration == instanceRegistration.Generation {
		logger.Info("Generation already observed. Nothing to do.")
		return reconcile.Result{}, nil
	}

	tenantId := instanceRegistration.Namespace

	//check if namespace exists and create if necessary
	targetNamespace := &corev1.Namespace{}
	targetNamespace.SetName(tenantId)
	if err := c.Client().Get(ctx, client.ObjectKeyFromObject(targetNamespace), targetNamespace); err != nil {
		if apierrors.IsNotFound(err) {
			if err := c.Client().Create(ctx, targetNamespace); err != nil {
				return reconcile.Result{}, err
			}
		} else {
			return reconcile.Result{}, err
		}
	}

	//create landscaperdeployment if not exist, else update spec
	landscaperDeployment := &lssv1alpha1.LandscaperDeployment{ObjectMeta: v1.ObjectMeta{
		Name:      instanceRegistration.Name,
		Namespace: instanceRegistration.Namespace,
	}}
	if _, err := controllerutil.CreateOrUpdate(ctx, c.Client(), landscaperDeployment, func() error {
		landscaperDeployment.Spec = instanceRegistration.Spec.LandscaperDeploymentSpec

		labels := landscaperDeployment.GetLabels()
		if labels == nil {
			labels = map[string]string{}
		}
		labels[INSTANCE_REGISTRATION_LABEL] = fmt.Sprintf("%s/%s", instanceRegistration.Namespace, instanceRegistration.Name)
		landscaperDeployment.SetLabels(labels)

		landscaperDeployment.Spec.TenantId = tenantId
		return nil
	}); err != nil {
		return reconcile.Result{}, err
	}

	instanceRegistration.Status.ObservedGeneration = instanceRegistration.Generation
	instanceRegistration.Status.LandscaperDeploymentInfo = &lssv1alpha1.LandscaperDeploymentInfo{
		Name:      landscaperDeployment.Name,
		Namespace: landscaperDeployment.Namespace,
	}
	if err := c.Client().Status().Update(ctx, instanceRegistration); err != nil {
		logger.Error(err, "failed updating status")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
