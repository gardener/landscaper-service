// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tenantregistration

import (
	"context"
	"fmt"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/operation"
	"github.com/gardener/landscaper-service/pkg/utils"
)

type Controller struct {
	operation.Operation
	log logging.Logger

	ReconcileFunc    func(ctx context.Context, tenantRegistration *lssv1alpha1.TenantRegistration) (reconcile.Result, error)
	HandleDeleteFunc func(ctx context.Context, tenantRegistration *lssv1alpha1.TenantRegistration) (reconcile.Result, error)
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

	tenantRegistration := &lssv1alpha1.TenantRegistration{}
	if err := c.Client().Get(ctx, req.NamespacedName, tenantRegistration); err != nil {
		logger.Error(err, "failed loading TenantRegistration cr")
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// set finalizer
	if tenantRegistration.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(tenantRegistration, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(tenantRegistration, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, tenantRegistration); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	// reconcile delete
	if !tenantRegistration.DeletionTimestamp.IsZero() {
		return c.HandleDeleteFunc(ctx, tenantRegistration)
	}

	return c.reconcile(ctx, tenantRegistration)
}

func (c *Controller) handleDelete(ctx context.Context, tenantRegistration *lssv1alpha1.TenantRegistration) (reconcile.Result, error) {
	// logger, ctx := logging.FromContextOrNew(ctx, nil)

	return reconcile.Result{}, nil

}

func (c *Controller) reconcile(ctx context.Context, tenantRegistration *lssv1alpha1.TenantRegistration) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	// if tenantRegistration.Status.SyncedGeneration != tenantRegistration.Generation {
	// 	logger.Info("SyncedGeneration unequal current generation -> sync is not completed")
	// 	return reconcile.Result{}, nil
	// }

	if tenantRegistration.Status.ObservedGeneration == tenantRegistration.Generation {
		logger.Info("Generation already observed. Nothing to do.")
		return reconcile.Result{}, nil
	}

	tenantNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "t",
		},
	}
	if err := c.Client().Create(ctx, tenantNamespace); err != nil {
		logger.Error(err, "failed creating tenant namespace")
		return reconcile.Result{}, err
	}

	for _, roleName := range []string{"admin", "member", "editor"} {
		if err := utils.CreateRoleIfNotExistOrUpdate(ctx, roleName, tenantNamespace.Name, getRolePolicyRules(), c.Client()); err != nil {
			logger.Error(err, "failed creating/updating role")
			return reconcile.Result{}, err
		}

		//create rolebinding
		roleBindingName := fmt.Sprintf("%s-binding", roleName)
		if err := utils.CreateRoleBindingIfNotExistOrUpdate(ctx, roleBindingName, tenantNamespace.Name, roleName, c.Client()); err != nil {
			logger.Error(err, "failed creating/updating rolebinding", "name", roleBindingName)
			return reconcile.Result{}, err
		}
	}

	//create subjectsynclist and let other controller write the admin user
	//TODO

	tenantRegistration.Status.Namespace = tenantNamespace.Name
	if err := c.Client().Update(ctx, tenantRegistration); err != nil {
		logger.Error(err, "failed updating tenantRegistration.Status.Namespace")
		return reconcile.Result{}, err
	}

	tenantRegistration.Status.ObservedGeneration = tenantNamespace.Generation
	if err := c.Client().Status().Update(ctx, tenantRegistration); err != nil {
		logger.Error(err, "failed updating observed generation")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func getRolePolicyRules() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{
		// TODO: add rules that are required
		// {
		// 	APIGroups: []string{"landscaper-service.gardener.cloud"},
		// 	Resources: []string{"*"},
		// 	Verbs:     []string{"*"},
		// },
	}
}
