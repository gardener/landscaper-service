// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package subjectsync

import (
	"context"
	"strings"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	config "github.com/gardener/landscaper-service/pkg/apis/config/v1alpha1"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/operation"
)

type Controller struct {
	operation.TargetShootSidecarOperation
	log logging.Logger

	ReconcileFunc func(ctx context.Context, subjectList *lssv1alpha1.SubjectList) (reconcile.Result, error)
}

func NewController(logger logging.Logger, c client.Client, scheme *runtime.Scheme, config *config.TargetShootSidecarConfiguration) (reconcile.Reconciler, error) {
	ctrl := &Controller{
		log: logger,
	}
	ctrl.ReconcileFunc = ctrl.reconcile
	op := operation.NewTargetShootSidecarOperation(c, scheme, config)
	ctrl.TargetShootSidecarOperation = *op
	return ctrl, nil
}

// NewTestActuator creates a new controller for testing purposes.
func NewTestActuator(op operation.TargetShootSidecarOperation, logger logging.Logger) *Controller {
	ctrl := &Controller{
		TargetShootSidecarOperation: op,
		log:                         logger,
	}
	return ctrl
}

func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger, ctx := c.log.StartReconcileAndAddToContext(ctx, req)

	logger.Info("start reconcile subjectList")

	subjectList := &lssv1alpha1.SubjectList{}
	if err := c.Client().Get(ctx, req.NamespacedName, subjectList); err != nil {
		logger.Error(err, "failed loading subjectlist cr")
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// set finalizer
	if subjectList.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(subjectList, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(subjectList, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, subjectList); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	if !subjectList.DeletionTimestamp.IsZero() {
		// deletion of the subjectlist is not allowed
		return reconcile.Result{}, nil
	}

	return c.reconcile(ctx, subjectList)
}

func (c *Controller) reconcile(ctx context.Context, subjectList *lssv1alpha1.SubjectList) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	// convert subjects of the SubjectList custom resource into rbac subjects
	subjects := CreateSubjectsForSubjectList(ctx, subjectList)
	viewerSubjects := CreateViewerSubjectsForSubjectList(ctx, subjectList)

	userClusterRoleDef := GetUserClusterRoleDefinition()

	if err := userClusterRoleDef.CreateOrUpdateClusterRole(ctx, c.Client()); err != nil {
		logger.Error(err, "failed updating user cluster role")
		return reconcile.Result{}, err
	}

	if err := userClusterRoleDef.CreateOrUpdateClusterRoleBinding(ctx, c.Client(), subjects); err != nil {
		logger.Error(err, "failed updating user cluster role binding")
		return reconcile.Result{}, err
	}

	viewerClusterRoleDef := GetViewerClusterRoleDefinition()

	if err := viewerClusterRoleDef.CreateOrUpdateClusterRole(ctx, c.Client()); err != nil {
		logger.Error(err, "failed updating viewer cluster role")
		return reconcile.Result{}, err
	}

	if err := viewerClusterRoleDef.CreateOrUpdateClusterRoleBinding(ctx, c.Client(), viewerSubjects); err != nil {
		logger.Error(err, "failed updating viewer cluster role binding")
		return reconcile.Result{}, err
	}

	roleBindings := &rbacv1.RoleBindingList{}
	if err := c.Client().List(ctx, roleBindings); err != nil {
		logger.Error(err, "failed loading role bindings")
		return reconcile.Result{}, err
	}

	for _, roleBinding := range roleBindings.Items {
		logger, ctx := logging.FromContextOrNew(ctx, nil, "roleBinding", apitypes.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}.String())

		switch roleBinding.Name {
		case LS_USER_ROLE_BINDING_IN_NAMESPACE:
			if roleBinding.Namespace != LS_USER_NAMESPACE {
				logger.Info("ls-user role binding found outside of ls-user namespace. Reconcile skipped: " + roleBinding.Namespace)
				continue
			}

			if err := UpdateRoleBindingSubjects(ctx, c.Client(), &roleBinding, subjects); err != nil {
				return reconcile.Result{}, err
			}

		case USER_ROLE_BINDING_IN_NAMESPACE:
			if !strings.HasPrefix(roleBinding.Namespace, CUSTOM_NS_PREFIX) {
				logger.Info("user role binding found outside of customer namespace. Reconcile skipped: " + roleBinding.Namespace)
				continue
			}

			if err := UpdateRoleBindingSubjects(ctx, c.Client(), &roleBinding, subjects); err != nil {
				return reconcile.Result{}, err
			}

		case VIEWER_ROLE_BINDING_IN_NAMESPACE:
			if !strings.HasPrefix(roleBinding.Namespace, CUSTOM_NS_PREFIX) {
				logger.Info("viewer role binding found outside of customer namespace. Reconcile skipped: " + roleBinding.Namespace)
				continue
			}

			if err := UpdateRoleBindingSubjects(ctx, c.Client(), &roleBinding, viewerSubjects); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil
}
