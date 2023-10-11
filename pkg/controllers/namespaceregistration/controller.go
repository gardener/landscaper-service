// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package namespaceregistration

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	"github.com/gardener/landscaper-service/pkg/apis/constants"
	dataplanev1alpha2 "github.com/gardener/landscaper-service/pkg/apis/dataplane/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"
	"github.com/gardener/landscaper-service/pkg/operation"
	"github.com/gardener/landscaper-service/pkg/utils"
)

const (
	PhaseCreating  = "Creating"
	PhaseFailed    = "Failed"
	PhaseCompleted = "Completed"
	PhaseDeleting  = "Deleting"

	ReasonInvalidName = "invalid name"

	requeueAfterDuration = 30 * time.Second
)

type Controller struct {
	operation.TargetShootSidecarOperation
	log logging.Logger

	ReconcileFunc    func(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration) (reconcile.Result, error)
	HandleDeleteFunc func(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration) (reconcile.Result, error)
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

	namespaceRegistration := &dataplanev1alpha2.NamespaceRegistration{}
	if err := c.Client().Get(ctx, req.NamespacedName, namespaceRegistration); err != nil {
		logger.Error(err, "failed loading namespaceregistration")
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{RequeueAfter: requeueAfterDuration}, nil
	}

	if !strings.HasPrefix(namespaceRegistration.Name, subjectsync.CUSTOM_NS_PREFIX) {
		if namespaceRegistration.Status.Phase != PhaseFailed ||
			namespaceRegistration.Status.LastError == nil ||
			(namespaceRegistration.Status.LastError != nil && namespaceRegistration.Status.LastError.Reason != ReasonInvalidName) {

			err := fmt.Errorf("name must start with %q", subjectsync.CUSTOM_NS_PREFIX)
			lastError := c.createError(namespaceRegistration.Status.Phase, ReasonInvalidName, err)
			c.updateStatus(namespaceRegistration, PhaseFailed, lastError)
			if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
				logger.Error(err, "failed updating namespaceregistration with invalid name - must start with "+subjectsync.CUSTOM_NS_PREFIX)
				return reconcile.Result{RequeueAfter: requeueAfterDuration}, nil
			}
		}

		return reconcile.Result{}, nil
	}

	// set finalizer
	if namespaceRegistration.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(namespaceRegistration, constants.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(namespaceRegistration, constants.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, namespaceRegistration); err != nil {
			logger.Error(err, "failed adding finalizer to namespaceregistration")
			return reconcile.Result{RequeueAfter: requeueAfterDuration}, nil
		}
		// do not return here because the controller only watches for particular events and setting a finalizer is not part of this
	}

	// reconcile delete
	if !namespaceRegistration.DeletionTimestamp.IsZero() {
		return c.HandleDeleteFunc(ctx, namespaceRegistration)
	}

	return c.reconcile(ctx, namespaceRegistration)
}

func (c *Controller) handleDelete(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	namespace := &corev1.Namespace{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: namespaceRegistration.GetName()}, namespace); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("namespace not found, removing namespaceregistration")
			controllerutil.RemoveFinalizer(namespaceRegistration, constants.LandscaperServiceFinalizer)
			if err := c.Client().Update(ctx, namespaceRegistration); err != nil {
				logger.Error(err, "failed removing finalizer")
				return reconcile.Result{RequeueAfter: requeueAfterDuration}, nil
			}
			return reconcile.Result{}, nil
		}
		logger.Error(err, "failed loading namespace")
		return reconcile.Result{RequeueAfter: requeueAfterDuration}, nil
	}

	return c.removeResourcesAndNamespace(ctx, namespaceRegistration, namespace)
}

func (c *Controller) removeResourcesAndNamespace(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration,
	namespace *corev1.Namespace) (reconcile.Result, error) {

	logger, ctx := logging.FromContextOrNew(ctx, nil)

	if namespaceRegistration.Status.Phase != PhaseDeleting {
		c.updateStatus(namespaceRegistration, PhaseDeleting, nil)
		if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
			logger.Error(err, "failed updating status of namespaceregistration when starting deletion")
			return reconcile.Result{RequeueAfter: requeueAfterDuration}, nil
		}
	}

	// check if installations, executions, deploy items or target sync objects are still there
	installations := &v1alpha1.InstallationList{}
	if err := c.Client().List(ctx, installations, client.InNamespace(namespaceRegistration.GetName())); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed reading installations", err)
	}

	if len(installations.Items) > 0 {
		err := c.triggerDeletionOfInstallations(ctx, namespaceRegistration, installations.Items)
		if err != nil {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed deleting installations", err)
		}

		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "namespace contains installations", nil)
	}

	executions := &v1alpha1.ExecutionList{}
	if err := c.Client().List(ctx, executions, client.InNamespace(namespaceRegistration.GetName())); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed reading executions", err)
	}

	if len(executions.Items) > 0 {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "namespace contains executions", nil)
	}

	deployItems := &v1alpha1.DeployItemList{}
	if err := c.Client().List(ctx, deployItems, client.InNamespace(namespaceRegistration.GetName())); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed reading deploy items", err)
	}

	if len(deployItems.Items) > 0 {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "namespace contains deploy items", nil)
	}

	targetSyncs := &v1alpha1.TargetSyncList{}
	if err := c.Client().List(ctx, targetSyncs, client.InNamespace(namespaceRegistration.GetName())); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed reading targetsyncs", err)
	}

	if len(targetSyncs.Items) > 0 {
		for i := range targetSyncs.Items {
			nextTargetSync := &targetSyncs.Items[i]
			if err := c.Client().Delete(ctx, nextTargetSync); err != nil {
				return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed removing targetsync", err)
			}
		}

		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "namespace contains targetsyncs", nil)
	}

	return c.removeAccessDataAndNamespace(ctx, namespaceRegistration, namespace)
}

func (c *Controller) removeAccessDataAndNamespace(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration,
	namespace *corev1.Namespace) (reconcile.Result, error) {

	logger, ctx := logging.FromContextOrNew(ctx, nil)

	// delete role binding
	roleBinding := &rbacv1.RoleBinding{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_BINDING_IN_NAMESPACE, Namespace: namespaceRegistration.GetName()}, roleBinding); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("rolebinding in namespace not found")
		} else {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed loading rolebinding", err)
		}
	} else {
		if err := c.Client().Delete(ctx, roleBinding); err != nil {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed deleting rolebinding", err)
		}
	}

	//delete role
	role := &rbacv1.Role{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: subjectsync.USER_ROLE_IN_NAMESPACE, Namespace: namespaceRegistration.GetName()}, role); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("role in namespace not found")
		} else {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed loading role", err)
		}
	} else {
		if err := c.Client().Delete(ctx, role); err != nil {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed deleting role", err)
		}
	}

	// delete namespace
	if err := c.Client().Delete(ctx, namespace); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed deleting namespace", err)
	}

	controllerutil.RemoveFinalizer(namespaceRegistration, constants.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, namespaceRegistration); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseDeleting, "failed removing finalizer", err)
	}

	return reconcile.Result{}, nil
}

func (c *Controller) reconcile(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	if namespaceRegistration.Status.Phase == PhaseCompleted {
		logger.Debug("Phase already in Completed")
		return reconcile.Result{}, nil
	}

	if namespaceRegistration.Status.Phase == "" {
		c.updateStatus(namespaceRegistration, PhaseCreating, namespaceRegistration.Status.LastError)
		if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
			logger.Error(err, "failed updating status of namespaceregistration when starting namespace creation")
			return reconcile.Result{RequeueAfter: requeueAfterDuration}, nil
		}
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceRegistration.Name,
		},
	}

	if err := c.Client().Create(ctx, namespace); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseCreating, "failed creating namespace", err)
		}
	}

	if err := c.createRoleIfNotExistOrUpdate(ctx, namespaceRegistration); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseCreating, "failed creating role", err)
	}

	if err := c.createRoleBindingIfNotExistOrUpdate(ctx, namespaceRegistration); err != nil {
		return c.logErrorUpdateAndRetry(ctx, namespaceRegistration, PhaseCreating, "failed creating rolebinding", err)
	}

	c.updateStatus(namespaceRegistration, PhaseCompleted, nil)
	if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
		logger.Error(err, "failed updating status of namespaceregistration after completion")
		return reconcile.Result{RequeueAfter: requeueAfterDuration}, nil
	}
	return reconcile.Result{}, nil
}

func (c *Controller) createRoleIfNotExistOrUpdate(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration) error {
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

func (c *Controller) createRoleBindingIfNotExistOrUpdate(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	// load subjectList from CR
	subjectList := &dataplanev1alpha2.SubjectList{}
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

func (c *Controller) triggerDeletionOfInstallations(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration, installations []v1alpha1.Installation) error {
	triggerDeletion, err := getTriggerDeletionFunction(ctx, namespaceRegistration)
	if err != nil {
		return err
	}

	// trigger deletion of root installations according to deletion strategy
	var triggerErr error
	for i := range installations {
		inst := &installations[i]
		if !utils.HasLabel(&inst.ObjectMeta, v1alpha1.EncompassedByLabel) {
			if tmpErr := triggerDeletion(ctx, c.Client(), inst); tmpErr != nil {
				triggerErr = tmpErr
			}
		}
	}

	return triggerErr
}

func (c *Controller) logErrorUpdateAndRetry(ctx context.Context, namespaceRegistration *dataplanev1alpha2.NamespaceRegistration,
	phase, msg string, err error) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	if err != nil {
		logger.Error(err, msg)
	} else {
		logger.Info(msg)
	}

	lastError := c.createError(namespaceRegistration.Status.Phase, msg, err)
	c.updateStatus(namespaceRegistration, phase, lastError)
	if err := c.Client().Status().Update(ctx, namespaceRegistration); err != nil {
		logger.Error(err, "failed updating status of namespaceregistration after error: "+msg)
	}

	return reconcile.Result{RequeueAfter: requeueAfterDuration}, nil
}

func (c *Controller) updateStatus(namespaceRegistration *dataplanev1alpha2.NamespaceRegistration, phase string,
	lastError *dataplanev1alpha2.Error) {
	namespaceRegistration.Status.Phase = phase
	namespaceRegistration.Status.LastError = lastError
}

func (c *Controller) createError(phase, reason string, err error) *dataplanev1alpha2.Error {
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	return &dataplanev1alpha2.Error{
		Operation:          phase,
		LastTransitionTime: metav1.Now(),
		LastUpdateTime:     metav1.Now(),
		Reason:             reason,
		Message:            msg,
	}
}
