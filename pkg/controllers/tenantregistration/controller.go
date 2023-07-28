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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	"github.com/gardener/landscaper-service/pkg/controllers/lossubject"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/operation"
	"github.com/gardener/landscaper-service/pkg/utils"
)

const TENANT_LABEL_ON_NAMESPACE string = "landscaper-service.gardener.cloud/tenant"

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
	logger, ctx := logging.FromContextOrNew(ctx, nil)
	controllerutil.RemoveFinalizer(tenantRegistration, lssv1alpha1.LandscaperServiceFinalizer)
	if err := c.Client().Update(ctx, tenantRegistration); err != nil {
		logger.Error(err, "Failed removing finalizer")
		return reconcile.Result{}, err
	}
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

	tenantNamespaceName, err := c.createTenantNamespaceIfNotExistAndGetName(ctx, tenantRegistration)
	if err != nil {
		logger.Error(err, "failed ensuring tenant namespace exist")
		return reconcile.Result{}, err
	}

	for _, roleInfo := range lossubject.ALL_ROLES() {
		if err := utils.CreateRoleIfNotExistOrUpdate(ctx, roleInfo.RoleName, *tenantNamespaceName, roleInfo.PrivilegeList, c.Client()); err != nil {
			logger.Error(err, "failed creating/updating role")
			return reconcile.Result{}, err
		}

		//create rolebinding
		if err := utils.CreateRoleBindingIfNotExistOrUpdate(ctx, roleInfo.RoleBindingName, *tenantNamespaceName, roleInfo.RoleName, c.Client()); err != nil {
			logger.Error(err, "failed creating/updating rolebinding", "name", roleInfo.RoleBindingName)
			return reconcile.Result{}, err
		}
	}

	//create ClusterRole+Binding to read TenantRegistration
	if err := c.createTenantReadClusterRoleAndBinding(ctx, tenantRegistration, *tenantNamespaceName); err != nil {
		logger.Error(err, "failed creating ClusterRole and Clusterrolebinding for tenantregistration")
		return reconcile.Result{}, err
	}

	//create subjectsynclist and let other controller handle the initial admin user
	if err := c.createLosSubjectListIfNotExist(ctx, tenantRegistration.Spec.Author, *tenantNamespaceName); err != nil {
		logger.Error(err, "failed creating losSubjectList")
		return reconcile.Result{}, err
	}

	tenantRegistration.Status.Namespace = *tenantNamespaceName
	tenantRegistration.Status.ObservedGeneration = tenantRegistration.Generation
	if err := c.Client().Status().Update(ctx, tenantRegistration); err != nil {
		logger.Error(err, "failed updating status")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

// createTenantNamespaceIfNotExistAndGetName is an idempotent function (running multiple times having the same effect) to only create one namespace no matter the number of calls and returns its name.
// It uses the tenantRegistration.Name as label value to determine, is a namespace for this tenantRegistration has been created before.
func (c *Controller) createTenantNamespaceIfNotExistAndGetName(ctx context.Context, tenantRegistration *lssv1alpha1.TenantRegistration) (*string, error) {
	labelSelector, err := labels.NewRequirement(TENANT_LABEL_ON_NAMESPACE, selection.Equals, []string{tenantRegistration.Name})
	if err != nil {
		return nil, fmt.Errorf("failed constructing namespace list label selctor requirement: %w", err)
	}

	namespaceList := &corev1.NamespaceList{}
	if err := c.Client().List(ctx, namespaceList, client.MatchingLabelsSelector{Selector: labels.NewSelector().Add(*labelSelector)}); err != nil {
		return nil, fmt.Errorf("failed listing namespaces with customer label: %w", err)
	}

	if len(namespaceList.Items) > 1 {
		return nil, fmt.Errorf("listing namespaces with customer label should return 0 or 1, not %d", len(namespaceList.Items))
	} else if len(namespaceList.Items) == 1 {
		return &namespaceList.Items[0].Name, nil
	} else {
		tenantNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "t",
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion:         "landscaper-service.gardener.cloud/v1alpha1", //TODO: cant read it from tenantRegistration since it is not populated. Maybe get the fields populated or use constants or something
						Kind:               "TenantRegistration",
						Name:               tenantRegistration.Name,
						UID:                tenantRegistration.GetUID(),
						BlockOwnerDeletion: utils.Ptr(true),
					},
				},
				Labels: map[string]string{
					TENANT_LABEL_ON_NAMESPACE: tenantRegistration.Name,
				},
			},
		}
		if err := c.Client().Create(ctx, tenantNamespace); err != nil {
			return nil, fmt.Errorf("failed creating tenant namespace: %w", err)
		}
		return &tenantNamespace.Name, nil
	}
}

func (c *Controller) createTenantReadClusterRoleAndBinding(ctx context.Context, tenantRegistration *lssv1alpha1.TenantRegistration, tenantId string) error {
	clusterRoleInfo := lossubject.TENANTREGISTRATION_READ_CLUSTER_ROLE_INFO(tenantRegistration.Name, tenantId)
	if err := utils.CreateClusterRoleIfNotExistOrUpdate(ctx, clusterRoleInfo.RoleName, clusterRoleInfo.PrivilegeList, c.Client()); err != nil {
		return err
	}
	if err := utils.CreateClusterRoleBindingIfNotExistOrUpdate(ctx, clusterRoleInfo.RoleBindingName, clusterRoleInfo.RoleName, c.Client()); err != nil {
		return err
	}
	return nil
}

func (c *Controller) createLosSubjectListIfNotExist(ctx context.Context, initialUser string, namespace string) error {
	losSubjectSpec := lssv1alpha1.LosSubjectListSpec{
		Admins: []lssv1alpha1.LosSubject{
			{
				Kind: "User",
				Name: initialUser,
			},
		},
		Members: []lssv1alpha1.LosSubject{},
		Viewer:  []lssv1alpha1.LosSubject{},
	}

	losSubjectList := lssv1alpha1.LosSubjectList{
		ObjectMeta: metav1.ObjectMeta{
			Name:      lossubject.SUBJECTLIST_NAME,
			Namespace: namespace,
		},
		Spec: losSubjectSpec,
	}

	if err := c.Client().Create(ctx, &losSubjectList); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return fmt.Errorf("failed creating lossubjectlist %s: %w", losSubjectList.Name, err)
	}
	return nil
}
