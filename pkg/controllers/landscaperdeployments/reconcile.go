// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployments

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lsserrors "github.com/gardener/landscaper-service/pkg/apis/errors"
	lssscheduling "github.com/gardener/landscaper-service/pkg/controllers/landscaperdeployments/scheduling"
	"github.com/gardener/landscaper-service/pkg/utils"
)

// reconcile reconciles a landscaper deployment
func (c *Controller) reconcile(ctx context.Context, deployment *lssv1alpha1.LandscaperDeployment) error {
	currOp := "Reconcile"
	oldDeployment := deployment.DeepCopy()

	// reconcile instance
	instance := &lssv1alpha1.Instance{}
	instance.GenerateName = fmt.Sprintf("%s-", deployment.GetName())
	instance.Namespace = deployment.GetNamespace()
	if deployment.Status.InstanceRef != nil && !deployment.Status.InstanceRef.IsEmpty() {
		instance.Name = deployment.Status.InstanceRef.Name
		instance.Namespace = deployment.Status.InstanceRef.Namespace
	}

	_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), instance, func() error {
		return c.mutateInstance(ctx, deployment, instance)
	})

	if err != nil {
		return lsserrors.NewWrappedError(err, currOp, "CreateUpdateInstance", err.Error())
	}

	// set the instance reference for the deployment if not already set
	if deployment.Status.InstanceRef == nil || !deployment.Status.InstanceRef.IsObject(instance) {
		deployment.Status.InstanceRef = &lssv1alpha1.ObjectReference{
			Name:      instance.GetName(),
			Namespace: instance.GetNamespace(),
		}
	}

	if err := c.Client().Get(ctx, client.ObjectKeyFromObject(instance), instance); err != nil {
		return lsserrors.NewWrappedError(err, currOp, "GetInstanceStatus", err.Error())
	}

	deployment.Status.Phase = instance.Status.Phase

	if deployment.IsInternalDataPlane() {
		deployment.Status.DataPlaneType = lssv1alpha1.LandscaperDeploymentDataPlaneTypeInternal
	} else {
		deployment.Status.DataPlaneType = lssv1alpha1.LandscaperDeploymentDataPlaneTypeExternal
	}

	if !reflect.DeepEqual(oldDeployment.Status, deployment.Status) {
		if err := c.Client().Status().Update(ctx, deployment); err != nil {
			return lsserrors.NewWrappedError(err, currOp, "UpdateLandscaperDeploymentStatus", err.Error())
		}
	}

	// if not already added, add the instance reference to the service target configuration
	serviceTargetConf := &lssv1alpha1.ServiceTargetConfig{}
	if err := c.Client().Get(ctx, instance.Spec.ServiceTargetConfigRef.NamespacedName(), serviceTargetConf); err != nil {
		return lsserrors.NewWrappedError(err, currOp, "GetServiceTargetConfig", err.Error())
	}

	instanceRef := &lssv1alpha1.ObjectReference{
		Name:      instance.GetName(),
		Namespace: instance.GetNamespace(),
	}
	if !utils.ContainsReference(serviceTargetConf.Status.InstanceRefs, instanceRef) {
		serviceTargetConf.Status.InstanceRefs = append(serviceTargetConf.Status.InstanceRefs, *instanceRef)

		if err := c.Client().Status().Update(ctx, serviceTargetConf); err != nil {
			return lsserrors.NewWrappedError(err, currOp, "UpdateServiceTargetConfReferences", err.Error())
		}
	}

	return nil
}

// mutateInstance creates/updates the instance for a landscaper deployment
func (c *Controller) mutateInstance(ctx context.Context, deployment *lssv1alpha1.LandscaperDeployment, instance *lssv1alpha1.Instance) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyReconciledResource, client.ObjectKeyFromObject(deployment).String()},
		lc.KeyMethod, "mutateInstance")

	if len(deployment.Name) > 0 {
		logger.Info("Updating instance", lc.KeyResource, client.ObjectKeyFromObject(instance).String())
	} else {
		logger.Info("Creating instance", lc.KeyResource, types.NamespacedName{Name: instance.GenerateName, Namespace: instance.Namespace}.String())
	}

	if err := controllerutil.SetControllerReference(deployment, instance, c.Scheme()); err != nil {
		return fmt.Errorf("unable to set controller reference for instance: %w", err)
	}

	if utils.HasOperationAnnotation(deployment, lssv1alpha1.LandscaperServiceOperationIgnore) {
		utils.SetOperationAnnotation(instance, lssv1alpha1.LandscaperServiceOperationIgnore)
	} else {
		if utils.HasOperationAnnotation(instance, lssv1alpha1.LandscaperServiceOperationIgnore) {
			utils.RemoveOperationAnnotation(instance)
		}
	}

	if len(instance.Spec.ServiceTargetConfigRef.Name) == 0 {
		// try to find a service target configuration that can be used for this landscaper deployment
		var serviceTargetConf *lssv1alpha1.ServiceTargetConfig
		var err error

		serviceTargetConf, err = c.findServiceTargetConfigByScheduling(ctx, deployment)
		if err != nil {
			return err
		}

		instance.Spec.ServiceTargetConfigRef.Name = serviceTargetConf.GetName()
		instance.Spec.ServiceTargetConfigRef.Namespace = serviceTargetConf.GetNamespace()
	}

	if len(instance.Spec.ID) == 0 {
		instanceList := &lssv1alpha1.InstanceList{}
		if err := c.Client().List(ctx, instanceList, &client.ListOptions{Namespace: deployment.Namespace}); err != nil {
			return fmt.Errorf("unable to list instances in namespace %s: %w", deployment.Namespace, err)
		}

		existingIds := sets.NewString()
		for _, i := range instanceList.Items {
			existingIds.Insert(i.Spec.ID)
		}

		var id string
		for id = c.NewUniqueID(); existingIds.Has(id); id = c.NewUniqueID() {
		}
		instance.Spec.ID = id
	}

	instance.Spec.TenantId = deployment.Spec.TenantId
	instance.Spec.LandscaperConfiguration = deployment.Spec.LandscaperConfiguration
	instance.Spec.OIDCConfig = deployment.Spec.OIDCConfig
	instance.Spec.HighAvailabilityConfig = deployment.Spec.HighAvailabilityConfig
	instance.Spec.DataPlane = deployment.Spec.DataPlane

	c.Operation.Scheme().Default(instance)

	return nil
}

func (c *Controller) findServiceTargetConfigByScheduling(ctx context.Context, deployment *lssv1alpha1.LandscaperDeployment) (*lssv1alpha1.ServiceTargetConfig, error) {
	log, ctx := logging.FromContextOrNew(ctx, nil)

	serviceTargetConfigs, err := c.getVisibleServiceTargetConfigs(ctx)
	if err != nil {
		return nil, err
	}

	scheduling, err := c.getSchedulingResource(ctx)
	if err != nil {
		return nil, err
	}

	// determine a matching service target config
	winner, err := lssscheduling.FindServiceTargetConfig(scheduling, deployment, serviceTargetConfigs)
	if err != nil {
		log.Error(err, "unable to find service target config")
		return nil, fmt.Errorf("unable to find service target config: %w", err)
	}

	return winner, nil
}

func (c *Controller) getVisibleServiceTargetConfigs(ctx context.Context) ([]lssv1alpha1.ServiceTargetConfig, error) {
	log, ctx := logging.FromContextOrNew(ctx, nil)

	serviceTargetConfigList := &lssv1alpha1.ServiceTargetConfigList{}
	selectorBuilder := strings.Builder{}
	selectorBuilder.WriteString(fmt.Sprintf("%s=true", lssv1alpha1.ServiceTargetConfigVisibleLabelName))

	labelSelector, _ := labels.Parse(selectorBuilder.String())
	listOptions := client.ListOptions{
		LabelSelector: labelSelector,
	}

	if err := c.Client().List(ctx, serviceTargetConfigList, &listOptions); err != nil {
		log.Error(err, "unable to list service target configs")
		return nil, fmt.Errorf("unable to list service target configs: %w", err)
	}

	return serviceTargetConfigList.Items, nil
}

// getSchedulingResource returns the TargetScheduling resource from the core cluster.
// Returns nil if scheduling is not configured or the scheduling resource does not exist.
func (c *Controller) getSchedulingResource(ctx context.Context) (*lssv1alpha1.TargetScheduling, error) {
	log, ctx := logging.FromContextOrNew(ctx, nil)

	schedulingConfig := c.Config().Scheduling
	if schedulingConfig == nil {
		log.Info("no scheduling configured")
		return nil, nil
	}

	schedulingKey := client.ObjectKey{
		Namespace: c.Config().Scheduling.Namespace,
		Name:      c.Config().Scheduling.Name,
	}
	scheduling := &lssv1alpha1.TargetScheduling{}
	if err := c.Client().Get(ctx, schedulingKey, scheduling); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("no scheduling resource configured")
			return nil, nil
		}

		log.Error(err, "unable to get scheduling object", lc.KeyResource, schedulingKey.String())
		return nil, fmt.Errorf("unable to get scheduling object %s: %w", schedulingKey.String(), err)
	}

	return scheduling, nil
}
