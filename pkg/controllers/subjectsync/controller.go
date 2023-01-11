// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package subjectsync

import (
	"context"
	"fmt"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper-service/pkg/operation"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"

	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apitypes "k8s.io/apimachinery/pkg/types"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

const USER_ROLE_IN_NAMESPACE = "user-role"
const USER_ROLE_BINDING_IN_NAMESPACE = "user-role-binding"
const LS_USER_ROLE_IN_NAMESPACE = "ls-user-role"
const LS_USER_ROLE_BINDING_IN_NAMESPACE = "ls-user-role-binding"

type Controller struct {
	operation.Operation
	log logging.Logger

	ReconcileFunc    func(ctx context.Context, subjectList *lssv1alpha1.SubjectList) (reconcile.Result, error)
	HandleDeleteFunc func(ctx context.Context, subjectList *lssv1alpha1.SubjectList) (reconcile.Result, error)
}

func NewController(logger logging.Logger, c client.Client, scheme *runtime.Scheme, config *coreconfig.LandscaperServiceConfiguration) (reconcile.Reconciler, error) {
	ctrl := &Controller{
		log: logger,
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
		Operation: op,
		log:       logger,
	}
	return ctrl
}

func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger, ctx := c.log.StartReconcileAndAddToContext(ctx, req)

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
		return reconcile.Result{}, nil
	}

	// reconcile delete
	if !subjectList.DeletionTimestamp.IsZero() {
		return c.HandleDeleteFunc(ctx, subjectList)
	}

	return c.reconcile(ctx, subjectList)
}

func (c *Controller) handleDelete(ctx context.Context, subjectList *lssv1alpha1.SubjectList) (reconcile.Result, error) {
	logger := c.log

	//TODO: on delete, remove all subjects from the role bindings?
	roleBindings := &rbacv1.RoleBindingList{}
	if err := c.Client().List(ctx, roleBindings); err != nil {
		logger.Error(err, "failed loading role bindings")
		return reconcile.Result{}, err
	}

	for _, roleBinding := range roleBindings.Items {
		logger, ctx := logging.FromContextOrNew(ctx, nil, "roleBinding", apitypes.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}.String())

		//check if it is a matching role binding
		if roleBinding.Name != USER_ROLE_BINDING_IN_NAMESPACE && roleBinding.Name != LS_USER_ROLE_BINDING_IN_NAMESPACE {
			continue
		}

		//remove subject list
		roleBinding.Subjects = []rbacv1.Subject{}

		//	write update role binding
		if err := c.Client().Update(ctx, &roleBinding); err != nil {
			logger.Error(err, "failed updating rolebinding %s %s", roleBinding.Namespace, roleBinding.Name)
			return reconcile.Result{}, err
		}
	}

	controllerutil.RemoveFinalizer(subjectList, lssv1alpha1.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, subjectList); err != nil {
		logger.Error(err, "failed removing finalizer")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil

}

func (c *Controller) reconcile(ctx context.Context, subjectList *lssv1alpha1.SubjectList) (reconcile.Result, error) {
	logger := c.log

	roleBindings := &rbacv1.RoleBindingList{}
	if err := c.Client().List(ctx, roleBindings); err != nil {
		logger.Error(err, "failed loading role bindings")
		return reconcile.Result{}, err
	}

	for _, roleBinding := range roleBindings.Items {
		logger, ctx := logging.FromContextOrNew(ctx, nil, "roleBinding", apitypes.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}.String())

		//check if it is a matching role binding (different naming in ls-user and other namespaces)
		if roleBinding.Name != USER_ROLE_BINDING_IN_NAMESPACE && roleBinding.Name != LS_USER_ROLE_BINDING_IN_NAMESPACE {
			continue
		}

		//remove subject list
		roleBinding.Subjects = []rbacv1.Subject{}

		//add subjects from SubjectList CR
		for _, subject := range subjectList.Spec.Subjects {
			rbacSubject, err := CreateSubjectForSubjectListEntry(subject)
			if err != nil {
				logger.Error(err, "could not create rbac.Subject from SubjectList.spec.subject")
				continue
			}
			roleBinding.Subjects = append(roleBinding.Subjects, *rbacSubject)
		}

		//	write update role binding
		if err := c.Client().Update(ctx, &roleBinding); err != nil {
			logger.Error(err, "failed updating rolebinding %s %s", roleBinding.Namespace, roleBinding.Name)
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}

func CreateSubjectForSubjectListEntry(subjectListEntry lssv1alpha1.Subject) (*rbacv1.Subject, error) {
	switch subjectListEntry.Kind {
	case "User", "Group":
		if subjectListEntry.Namespace != "" {
			return nil, fmt.Errorf("namespace must be empty for subject.Kind==User|Group")
		}
		return &rbacv1.Subject{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     subjectListEntry.Kind,
			Name:     subjectListEntry.Name,
		}, nil
	case "ServiceAccount":
		if subjectListEntry.Namespace == "" {
			return nil, fmt.Errorf("namespace must be set for subject.Kind==ServiceAccount")
		}
		return &rbacv1.Subject{
			APIGroup:  "",
			Kind:      subjectListEntry.Kind,
			Name:      subjectListEntry.Name,
			Namespace: subjectListEntry.Namespace,
		}, nil
	default:
		return nil, fmt.Errorf("subject kind %s unknown", subjectListEntry.Kind)
	}

}
