// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package namespaceregistration

import (
	"context"
	"fmt"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"
	"github.com/gardener/landscaper-service/pkg/operation"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

type Controller struct {
	operation.TargetShootSidecarOperation
	log logging.Logger

	ReconcileFunc    func(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) (reconcile.Result, error)
	HandleDeleteFunc func(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) (reconcile.Result, error)
}

func NewController(logger logging.Logger, c client.Client, scheme *runtime.Scheme, config *coreconfig.TargetShootSidecarConfiguration) (reconcile.Reconciler, error) {
	ctrl := &Controller{
		log: logger,
	}
	ctrl.ReconcileFunc = ctrl.reconcile
	ctrl.HandleDeleteFunc = ctrl.handleDelete
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
	ctrl.ReconcileFunc = ctrl.reconcile
	ctrl.HandleDeleteFunc = ctrl.handleDelete
	return ctrl
}

func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger, ctx := c.log.StartReconcileAndAddToContext(ctx, req)

	namespaceRegistration := &lssv1alpha1.NamespaceRegistration{}
	if err := c.Client().Get(ctx, req.NamespacedName, namespaceRegistration); err != nil {
		logger.Error(err, "failed loading namespaceregistration cr")
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// set finalizer
	if namespaceRegistration.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(namespaceRegistration, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(namespaceRegistration, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, namespaceRegistration); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	// reconcile delete
	if !namespaceRegistration.DeletionTimestamp.IsZero() {
		return c.HandleDeleteFunc(ctx, namespaceRegistration)
	}

	return c.reconcile(ctx, namespaceRegistration)
}

func (c *Controller) handleDelete(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) (reconcile.Result, error) {
	logger := c.log

	namespace := &corev1.Namespace{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: namespaceRegistration.GetName()}, namespace); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("namespace not found, removing namespaceregistration")
			controllerutil.RemoveFinalizer(namespaceRegistration, lssv1alpha1.LandscaperServiceFinalizer)
			if err := c.Client().Update(ctx, namespaceRegistration); err != nil {
				logger.Error(err, "failed removing finalizer")
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
		logger.Error(err, "failed loading namespace cr")
		return reconcile.Result{}, err
	}

	//delete role binding
	roleBinding := &rbacv1.RoleBinding{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: namespaceRegistration.GetName()}, roleBinding); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("rolebinding in namespace not found")
		} else {
			logger.Error(err, "failed loading rolebinding in namespace")
			return reconcile.Result{}, err
		}
	} else {
		if err := c.Client().Delete(ctx, roleBinding); err != nil {
			logger.Error(err, "failed deleting rolebinding in namespace")
			return reconcile.Result{}, err //TODO
		}
	}
	//delete role
	role := &rbacv1.Role{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_IN_NAMESPACE, Namespace: namespaceRegistration.GetName()}, role); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("role in namespace not found")
		} else {
			logger.Error(err, "failed loading role in namespace")
			return reconcile.Result{}, err
		}
	} else {
		if err := c.Client().Delete(ctx, role); err != nil {
			logger.Error(err, "failed deleting role in namespace")
			return reconcile.Result{}, err //TODO
		}
	}

	if err := c.Client().Delete(ctx, namespace); err != nil {
		logger.Error(err, "failed loading namespace cr")
		return reconcile.Result{}, err
	}
	controllerutil.RemoveFinalizer(namespaceRegistration, lssv1alpha1.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, namespaceRegistration); err != nil {
		logger.Error(err, "failed removing finalizer")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil

}

func (c *Controller) reconcile(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) (reconcile.Result, error) {
	logger := c.log

	if namespaceRegistration.Status.Phase == "Completed" {
		logger.Info("Phase already in Completed")
		return reconcile.Result{}, nil
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceRegistration.Name,
		},
	}
	if err := c.Client().Create(ctx, namespace); err != nil {
		if apierrors.IsAlreadyExists(err) {
			namespaceRegistration.Status.Phase = "Completed"
			if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
				logger.Error(err, "failed updating status of namespaceregistration")
				return reconcile.Result{}, err
			}
		}
		logger.Error(err, "failed creating namespace")
		return reconcile.Result{}, err
	}

	if err := c.createRoleIfNotExistOrUpdate(ctx, namespaceRegistration); err != nil {
		namespaceRegistration.Status.Phase = "Failed Role Creation"
	}
	if err := c.createRoleBindingIfNotExistOrUpdate(ctx, namespaceRegistration); err != nil {
		namespaceRegistration.Status.Phase = "Failed Role Binding"
	}

	namespaceRegistration.Status.Phase = "Completed"
	if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
		logger.Error(err, "failed updating status of namespaceregistration")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (c *Controller) createRoleIfNotExistOrUpdate(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) error {
	logger := c.log

	//create role
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subjectsync.USER_ROLE_IN_NAMESPACE,
			Namespace: namespaceRegistration.Name,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"landscaper.gardener.cloud"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets", "configmaps"},
				Verbs:     []string{"*"},
			},
		},
	}

	if err := c.Client().Create(ctx, role); err != nil {
		if apierrors.IsAlreadyExists(err) {
			if err := c.Client().Update(ctx, role); err != nil {
				logger.Error(err, "failed updating role")
				return err
			}
		}
		logger.Error(err, "failed creating role")
		return err
	}
	return nil
}

func (c *Controller) createRoleBindingIfNotExistOrUpdate(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) error {
	logger := c.log

	// load subjectList from CR
	subjectList := &lssv1alpha1.SubjectList{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: subjectsync.SUBJECT_LIST_NAME, Namespace: subjectsync.LS_USER_NAMESPACE}, subjectList); err != nil {
		logger.Error(err, "failed loading subjectlist cr")
		return fmt.Errorf("failed loading subjectlist %w", err)
	}

	//create role binding
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subjectsync.USER_ROLE_BINDING_IN_NAMESPACE,
			Namespace: namespaceRegistration.Name,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     subjectsync.USER_ROLE_IN_NAMESPACE,
		},
	}

	//add subjectlist subjects to rolebinding
	for _, subject := range subjectList.Spec.Subjects {
		rbacSubject, err := subjectsync.CreateSubjectForSubjectListEntry(subject)
		if err != nil {
			return fmt.Errorf("could not create rbac.Subject from SubjectList.spec.subject: %w", err) //TODO: change to continue?
		}
		roleBinding.Subjects = append(roleBinding.Subjects, *rbacSubject)
	}

	if err := c.Client().Create(ctx, roleBinding); err != nil {
		if apierrors.IsAlreadyExists(err) {
			if err := c.Client().Update(ctx, roleBinding); err != nil {
				logger.Error(err, "failed updating role binding")
				return err
			}
		}
		logger.Error(err, "failed creating role binding")
		return err
	}

	return nil
}
