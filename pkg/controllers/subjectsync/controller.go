// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package subjectsync

import (
	"context"
	"fmt"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"

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

	if err := c.createOrUpdateUserClusterRole(ctx); err != nil {
		logger.Error(err, "failed updating user cluster role")
		return reconcile.Result{}, err
	}

	if err := c.createOrUpdateUserClusterRoleBinding(ctx, subjects); err != nil {
		logger.Error(err, "failed updating user cluster role binding")
		return reconcile.Result{}, err
	}

	roleBindings := &rbacv1.RoleBindingList{}
	if err := c.Client().List(ctx, roleBindings); err != nil {
		logger.Error(err, "failed loading role bindings")
		return reconcile.Result{}, err
	}

	for _, roleBinding := range roleBindings.Items {
		logger, ctx := logging.FromContextOrNew(ctx, nil, "roleBinding", apitypes.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}.String())

		//check if it is a matching role binding (different naming in ls-user and other namespaces)
		//only process correct rolebindings
		if !(roleBinding.Name == USER_ROLE_BINDING_IN_NAMESPACE || roleBinding.Name == LS_USER_ROLE_BINDING_IN_NAMESPACE) {
			continue
		}

		if !(strings.HasPrefix(roleBinding.Namespace, CUSTOM_NS_PREFIX) || roleBinding.Namespace == LS_USER_NAMESPACE) {
			logger.Info("user-role/-binding found outside of customer namespace. Reconcile skipped: " + roleBinding.Namespace)
			continue
		}

		// update role binding
		roleBinding.Subjects = subjects
		if err := c.Client().Update(ctx, &roleBinding); err != nil {
			logger.Error(err, "failed updating rolebinding %s %s", roleBinding.Namespace, roleBinding.Name)
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func GetUserPolicyRules() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"namespaces"},
			Verbs:     []string{"get", "list", "watch"},
		},
		{
			APIGroups: []string{"landscaper-service.gardener.cloud"},
			Resources: []string{"subjectlists"},
			Verbs:     []string{"get", "list", "watch"},
		},
		{
			APIGroups: []string{"landscaper.gardener.cloud"},
			Resources: []string{"*"},
			Verbs:     []string{"get", "list", "watch"},
		},
		{
			APIGroups: []string{"apiextensions.k8s.io"},
			Resources: []string{"customresourcedefinitions"},
			Verbs:     []string{"get", "list", "watch"},
		},
	}
}

func (c *Controller) createOrUpdateUserClusterRole(ctx context.Context) error {
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: USER_CLUSTER_ROLE,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, c.Client(), role, func() error {
		role.Rules = GetUserPolicyRules()
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed ensuring user cluster role %s: %w", role.Name, err)
	}
	return nil
}

func (c *Controller) createOrUpdateUserClusterRoleBinding(ctx context.Context, subjects []rbacv1.Subject) error {
	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: USER_CLUSTER_ROLE_BINDING,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, c.Client(), roleBinding, func() error {
		roleBinding.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     USER_CLUSTER_ROLE,
		}
		roleBinding.Subjects = subjects
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed updating user cluster role binding %s: %w", roleBinding.Name, err)
	}

	return nil
}
