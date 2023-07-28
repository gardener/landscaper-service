// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package lossubject

import (
	"context"
	"fmt"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	rbacv1 "k8s.io/api/rbac/v1"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"

	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/operation"
)

type Controller struct {
	operation.Operation
	log logging.Logger

	ReconcileFunc    func(ctx context.Context, losSubjectList *lssv1alpha1.LosSubjectList) (reconcile.Result, error)
	HandleDeleteFunc func(ctx context.Context, losSubjectList *lssv1alpha1.LosSubjectList) (reconcile.Result, error)
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

	losSubjectList := &lssv1alpha1.LosSubjectList{}
	if err := c.Client().Get(ctx, req.NamespacedName, losSubjectList); err != nil {
		logger.Error(err, "failed loading LosSubjectList cr")
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// set finalizer
	if losSubjectList.DeletionTimestamp.IsZero() && !kutils.HasFinalizer(losSubjectList, lssv1alpha1.LandscaperServiceFinalizer) {
		controllerutil.AddFinalizer(losSubjectList, lssv1alpha1.LandscaperServiceFinalizer)
		if err := c.Client().Update(ctx, losSubjectList); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	// reconcile delete
	if !losSubjectList.DeletionTimestamp.IsZero() {
		return c.HandleDeleteFunc(ctx, losSubjectList)
	}

	return c.reconcile(ctx, losSubjectList)
}

func (c *Controller) handleDelete(ctx context.Context, losSubjectList *lssv1alpha1.LosSubjectList) (reconcile.Result, error) {
	// logger, ctx := logging.FromContextOrNew(ctx, nil)

	return reconcile.Result{}, nil

}

func (c *Controller) reconcile(ctx context.Context, losSubjectList *lssv1alpha1.LosSubjectList) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	// if tenantRegistration.Status.SyncedGeneration != tenantRegistration.Generation {
	// 	logger.Info("SyncedGeneration unequal current generation -> sync is not completed")
	// 	return reconcile.Result{}, nil
	// }

	if losSubjectList.Status.ObservedGeneration == losSubjectList.Generation {
		logger.Info("Generation already observed. Nothing to do.")
		return reconcile.Result{}, nil
	}

	if err := c.updateRoleBinding(ctx, ROLEBINDING_ADMIN, losSubjectList.Namespace, losSubjectList.Spec.Admins); err != nil {
		logger.Error(err, "Failed update Role binding", "Role", ROLE_ADMIN)
		return reconcile.Result{}, err
	}
	if err := c.updateRoleBinding(ctx, ROLEBINDING_MEMBER, losSubjectList.Namespace, losSubjectList.Spec.Members); err != nil {
		logger.Error(err, "Failed update Role binding", "Role", ROLE_MEMBER)
		return reconcile.Result{}, err
	}
	if err := c.updateRoleBinding(ctx, ROLEBINDING_VIEWER, losSubjectList.Namespace, losSubjectList.Spec.Viewer); err != nil {
		logger.Error(err, "Failed update Role binding", "Role", ROLE_VIEWER)
		return reconcile.Result{}, err
	}

	//assuming tenantid = namespace name
	tenantId := losSubjectList.Namespace
	if err := c.updateClusterRoleBinding(ctx, TENANTREGISTRATION_READ_CLUSTER_ROLE_BIDNING_NAME(tenantId), losSubjectList.Spec.Admins); err != nil {
		logger.Error(err, "Failed updating ClusterRole binding", "Role", ROLE_VIEWER)
		return reconcile.Result{}, err
	}

	losSubjectList.Status.ObservedGeneration = losSubjectList.Generation
	if err := c.Client().Status().Update(ctx, losSubjectList); err != nil {
		logger.Error(err, "failed updating status")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (c *Controller) updateRoleBinding(ctx context.Context, roleName string, namespace string, subjects []lssv1alpha1.LosSubject) error {
	rbacSubjectList := CreateSubjectsForLosSubjectList(ctx, &subjects)

	roleBinding := &rbacv1.RoleBinding{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: roleName, Namespace: namespace}, roleBinding); err != nil {
		return fmt.Errorf("failed loading role bindings: %v", err)
	}

	roleBinding.Subjects = rbacSubjectList
	if err := c.Client().Update(ctx, roleBinding); err != nil {
		return fmt.Errorf("failed updating rolebinding %s/%s: %v", namespace, roleName, err)
	}

	return nil
}

func (c *Controller) updateClusterRoleBinding(ctx context.Context, clusterRoleName string, subjects []lssv1alpha1.LosSubject) error {
	rbacSubjectList := CreateSubjectsForLosSubjectList(ctx, &subjects)

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	if err := c.Client().Get(ctx, types.NamespacedName{Name: clusterRoleName}, clusterRoleBinding); err != nil {
		return fmt.Errorf("failed loading clusterrolebindings: %v", err)
	}

	clusterRoleBinding.Subjects = rbacSubjectList
	if err := c.Client().Update(ctx, clusterRoleBinding); err != nil {
		return fmt.Errorf("failed updating clusterrolebinding %s: %v", clusterRoleName, err)
	}

	return nil
}

// CreateSubjectsForLosSubjectList converts the subjects of the custom resource into rbac subjects.
func CreateSubjectsForLosSubjectList(ctx context.Context, subjectList *[]lssv1alpha1.LosSubject) []rbacv1.Subject {
	logger, _ := logging.FromContextOrNew(ctx, nil)

	subjects := []rbacv1.Subject{}
	for _, subject := range *subjectList {
		rbacSubject, err := createSubjectForLosSubjectListEntry(subject)
		if err != nil {
			logger.Error(err, "could not create rbac.Subject from LosSubjectList.spec.subject")
			continue
		}
		subjects = append(subjects, *rbacSubject)
	}

	return subjects
}

// createSubjectForSubjectListEntry converts a single subject of the LosSubject custom resource into an rbac subject.
func createSubjectForLosSubjectListEntry(subjectListEntry lssv1alpha1.LosSubject) (*rbacv1.Subject, error) {
	switch subjectListEntry.Kind {
	case subjectsync.SUBJECT_LIST_ENTRY_USER, subjectsync.SUBJECT_LIST_ENTRY_GROUP:
		return &rbacv1.Subject{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     subjectListEntry.Kind,
			Name:     subjectListEntry.Name,
		}, nil
	default:
		return nil, fmt.Errorf("subject kind %s unknown", subjectListEntry.Kind)
	}
}
