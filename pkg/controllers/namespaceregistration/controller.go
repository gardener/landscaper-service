// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package namespaceregistration

import (
	"context"
	"fmt"
	"strings"

	"github.com/gardener/landscaper/apis/core/v1alpha1/helper"

	"github.com/gardener/landscaper-service/pkg/utils"

	"github.com/gardener/landscaper/apis/core/v1alpha1"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"
	"github.com/gardener/landscaper-service/pkg/operation"
)

const (
	PHASE_CREATING        = "Creating"
	PHASE_CREATION_FAILED = "CreationFailed"
	PHASE_COMPLETED       = "Completed"
	PHASE_DELETING        = "Deleting"
	PHASE_FAILED          = "DeletingFailed"
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

	logger.Info("start reconcile namespaceRegistration")

	namespaceRegistration := &lssv1alpha1.NamespaceRegistration{}
	if err := c.Client().Get(ctx, req.NamespacedName, namespaceRegistration); err != nil {
		logger.Error(err, "failed loading namespaceregistration cr")
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !strings.HasPrefix(namespaceRegistration.Name, subjectsync.CUSTOM_NS_PREFIX) {
		msg := "InvalidName: name must start with " + subjectsync.CUSTOM_NS_PREFIX
		if namespaceRegistration.Status.LastError != nil && msg != namespaceRegistration.Status.LastError.Message {
			lastError := c.createError(PHASE_CREATION_FAILED, msg, nil)
			c.updateStatus(namespaceRegistration, PHASE_CREATION_FAILED, lastError)
			if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
				logger.Error(err, "failed to update namespaceregistration with invalid name - must start with "+subjectsync.CUSTOM_NS_PREFIX)
				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{}, nil
	}

	// set finalizer
	if namespaceRegistration.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(namespaceRegistration, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(namespaceRegistration, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, namespaceRegistration); err != nil {
			return reconcile.Result{}, err
		}
		// do not return here because the controller only watches for particular events and setting a finalizer is not part of this
	}

	// reconcile delete
	if !namespaceRegistration.DeletionTimestamp.IsZero() {
		return c.HandleDeleteFunc(ctx, namespaceRegistration)
	}

	return c.reconcile(ctx, namespaceRegistration)
}

func (c *Controller) handleDelete(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

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

	return c.removeResourcesAndNamespace(ctx, namespaceRegistration, namespace)
}

func (c *Controller) removeResourcesAndNamespace(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration,
	namespace *corev1.Namespace) (reconcile.Result, error) {

	logger, ctx := logging.FromContextOrNew(ctx, nil)

	if namespaceRegistration.Status.Phase != PHASE_DELETING {
		c.updateStatus(namespaceRegistration, PHASE_DELETING, nil)
		if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
			logger.Error(err, "failed to update namespaceregistration with invalid name - must start with "+subjectsync.CUSTOM_NS_PREFIX)
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// check if installations, executions, deploy items or target sync objects are still there
	installations := &v1alpha1.InstallationList{}
	if err := c.Client().List(ctx, installations, client.InNamespace(namespaceRegistration.GetName())); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed reading installations", err)
	}

	if len(installations.Items) > 0 {
		var tmpErr error
		for i := range installations.Items {
			nextInst := &installations.Items[i]

			// delete root installations with delete without uninstall annotation
			if !utils.HasLabel(&nextInst.ObjectMeta, v1alpha1.EncompassedByLabel) && utils.HasDeleteWithoutUninstallAnnotation(&nextInst.ObjectMeta) {
				if nextInst.GetDeletionTimestamp().IsZero() {
					if err := c.Client().Delete(ctx, nextInst); err != nil {
						tmpErr = err
						logger.Error(err, "failed deleting installations without uninstall: "+client.ObjectKeyFromObject(nextInst).String())
					}
				} else if nextInst.Status.JobID == nextInst.Status.JobIDFinished && !helper.HasOperation(nextInst.ObjectMeta, v1alpha1.ReconcileOperation) {
					// retrigger
					metav1.SetMetaDataAnnotation(&nextInst.ObjectMeta, v1alpha1.OperationAnnotation, string(v1alpha1.ReconcileOperation))
					if err := c.Client().Update(ctx, nextInst); err != nil {
						tmpErr = err
						logger.Error(err, "failed annotating installations without uninstall: "+client.ObjectKeyFromObject(nextInst).String())
					}
				}
			}
		}

		if tmpErr != nil {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed deleting installations", tmpErr)
		} else {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "namespace contains installations", nil)
		}
	}

	executions := &v1alpha1.ExecutionList{}
	if err := c.Client().List(ctx, executions, client.InNamespace(namespaceRegistration.GetName())); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed reading executions", err)
	}

	if len(executions.Items) > 0 {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "namespace contains executions", nil)
	}

	deployItems := &v1alpha1.DeployItemList{}
	if err := c.Client().List(ctx, deployItems, client.InNamespace(namespaceRegistration.GetName())); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed reading deploy items", err)
	}

	if len(deployItems.Items) > 0 {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "namespace contains deploy items", nil)
	}

	targetSyncs := &v1alpha1.TargetSyncList{}
	if err := c.Client().List(ctx, targetSyncs, client.InNamespace(namespaceRegistration.GetName())); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed reading targetsyncs", err)
	}

	if len(targetSyncs.Items) > 0 {
		for i := range targetSyncs.Items {
			nextTargetSync := &targetSyncs.Items[i]
			if err := c.Client().Delete(ctx, nextTargetSync); err != nil {
				return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed removing targetsync", err)
			}
		}

		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "namespace contains targetsyncs", nil)
	}

	return c.removeAccessDataAndNamespace(ctx, namespaceRegistration, namespace)
}

func (c *Controller) removeAccessDataAndNamespace(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration,
	namespace *corev1.Namespace) (reconcile.Result, error) {

	logger, ctx := logging.FromContextOrNew(ctx, nil)

	// delete role binding
	roleBinding := &rbacv1.RoleBinding{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: namespaceRegistration.GetName()}, roleBinding); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("rolebinding in namespace not found")
		} else {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed loading rolebinding", err)
		}
	} else {
		if err := c.Client().Delete(ctx, roleBinding); err != nil {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed deleting rolebinding installations", err)
		}
	}
	//delete role
	role := &rbacv1.Role{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_IN_NAMESPACE, Namespace: namespaceRegistration.GetName()}, role); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("role in namespace not found")
		} else {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed loading role", err)
		}
	} else {
		if err := c.Client().Delete(ctx, role); err != nil {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed deleting role", err)
		}
	}

	if err := c.Client().Delete(ctx, namespace); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed deleting namespace", err)
	}

	controllerutil.RemoveFinalizer(namespaceRegistration, lssv1alpha1.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, namespaceRegistration); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_DELETING, "failed deleting finalizer", err)
	}

	return reconcile.Result{}, nil
}

func (c *Controller) reconcile(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	if namespaceRegistration.Status.Phase == PHASE_COMPLETED {
		logger.Debug("Phase already in Completed")
		return reconcile.Result{}, nil
	}

	if namespaceRegistration.Status.Phase == "" {
		c.updateStatus(namespaceRegistration, PHASE_CREATING, namespaceRegistration.Status.LastError)
		if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
			logger.Error(err, "failed updating status of namespaceregistration when starting namespace creation")
			return reconcile.Result{Requeue: true}, nil
		}
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceRegistration.Name,
		},
	}

	if err := c.Client().Create(ctx, namespace); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_CREATING, "failed creating namespace", err)
		}
	}

	if err := c.createRoleIfNotExistOrUpdate(ctx, namespaceRegistration); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_CREATING, "failed role creation", err)
	}

	if err := c.createRoleBindingIfNotExistOrUpdate(ctx, namespaceRegistration); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PHASE_CREATING, "failed role binding", err)
	}

	c.updateStatus(namespaceRegistration, PHASE_COMPLETED, namespaceRegistration.Status.LastError)
	if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
		logger.Error(err, "failed updating status of namespaceregistration after completion")
		return reconcile.Result{Requeue: true}, nil
	}
	return reconcile.Result{}, nil
}

func (c *Controller) createRoleIfNotExistOrUpdate(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	rules := []rbacv1.PolicyRule{
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
	}

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subjectsync.USER_ROLE_IN_NAMESPACE,
			Namespace: namespaceRegistration.Name,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, c.Client(), role, func() error {
		role.Rules = rules
		return nil
	})
	if err != nil {
		logger.Error(err, "failed ensuring user role")
		return fmt.Errorf("failed ensuring user role %s: %w", role.Name, err)
	}

	return nil
}

func (c *Controller) createRoleBindingIfNotExistOrUpdate(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	// load subjectList from CR
	subjectList := &lssv1alpha1.SubjectList{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: subjectsync.SUBJECT_LIST_NAME, Namespace: subjectsync.LS_USER_NAMESPACE}, subjectList); err != nil {
		logger.Error(err, "failed loading subjectlist cr")
		return fmt.Errorf("failed loading subjectlist %w", err)
	}

	// convert subjects of the SubjectList custom resource into rbac subjects
	subjects := subjectsync.CreateSubjectsForSubjectList(ctx, subjectList)

	//create role binding
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subjectsync.USER_ROLE_BINDING_IN_NAMESPACE,
			Namespace: namespaceRegistration.Name,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, c.Client(), roleBinding, func() error {
		roleBinding.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     subjectsync.USER_ROLE_IN_NAMESPACE,
		}
		roleBinding.Subjects = subjects
		return nil
	})
	if err != nil {
		logger.Error(err, "failed ensuring user role binding")
		return fmt.Errorf("failed ensuring role binding %s: %w", roleBinding.Name, err)
	}

	return nil
}

func (c *Controller) logErrorUpdateAndRetry(ctx context.Context, namespaceRegistration *lssv1alpha1.NamespaceRegistration,
	phase, msg string, err error) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	if err != nil {
		logger.Error(err, msg)
	} else {
		logger.Info(msg)
	}

	lastError := c.createError(phase, msg, err)
	c.updateStatus(namespaceRegistration, phase, lastError)
	if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
		logger.Error(err, "failed updating status of namespaceregistration"+msg)
	}

	return reconcile.Result{Requeue: true}, nil
}

func (c *Controller) updateStatus(namespaceRegistration *lssv1alpha1.NamespaceRegistration, phase string,
	lastError *lssv1alpha1.Error) {
	namespaceRegistration.Status.Phase = phase
	namespaceRegistration.Status.LastError = lastError
}

func (c *Controller) createError(phase, errorDescription string, err error) *lssv1alpha1.Error {
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	return &lssv1alpha1.Error{
		Operation:          phase,
		LastTransitionTime: metav1.Now(),
		LastUpdateTime:     metav1.Now(),
		Reason:             errorDescription,
		Message:            msg,
	}
}
