// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployments

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"

	"github.com/gardener/landscaper-service/pkg/apis/constants"
	"github.com/gardener/landscaper-service/pkg/apis/provisioning"
	lsserrors "github.com/gardener/landscaper-service/pkg/apis/provisioning/errors"
	provisioningv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/provisioning/v1alpha2"
	"github.com/gardener/landscaper-service/pkg/utils"
)

// reconcile reconciles a landscaper deployment
func (c *Controller) reconcile(ctx context.Context, deployment *provisioningv1alpha2.LandscaperDeployment) error {
	currOp := "Reconcile"

	// reconcile instance
	instance := &provisioningv1alpha2.Instance{}
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
		deployment.Status.InstanceRef = &provisioningv1alpha2.ObjectReference{
			Name:      instance.GetName(),
			Namespace: instance.GetNamespace(),
		}

		if err := c.Client().Status().Update(ctx, deployment); err != nil {
			return lsserrors.NewWrappedError(err, currOp, "UpdateInstanceRefForDeployment", err.Error())
		}
	}

	// if not already added, add the instance reference to the service target configuration
	serviceTargetConf := &provisioningv1alpha2.ServiceTargetConfig{}
	if err := c.Client().Get(ctx, instance.Spec.ServiceTargetConfigRef.NamespacedName(), serviceTargetConf); err != nil {
		return lsserrors.NewWrappedError(err, currOp, "GetServiceTargetConfig", err.Error())
	}

	instanceRef := &provisioningv1alpha2.ObjectReference{
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
func (c *Controller) mutateInstance(ctx context.Context, deployment *provisioningv1alpha2.LandscaperDeployment, instance *provisioningv1alpha2.Instance) error {
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

	if utils.HasOperationAnnotation(deployment, constants.LandscaperServiceOperationIgnore) {
		utils.SetOperationAnnotation(instance, constants.LandscaperServiceOperationIgnore)
	} else {
		if utils.HasOperationAnnotation(instance, constants.LandscaperServiceOperationIgnore) {
			utils.RemoveOperationAnnotation(instance)
		}
	}

	if len(instance.Spec.ServiceTargetConfigRef.Name) == 0 {
		// try to find a service target configuration that can be used for this landscaper deployment
		serviceTargetConf, err := c.findServiceTargetConfig(ctx)
		if err != nil {
			return err
		}

		instance.Spec.ServiceTargetConfigRef.Name = serviceTargetConf.GetName()
		instance.Spec.ServiceTargetConfigRef.Namespace = serviceTargetConf.GetNamespace()
	}

	if len(instance.Spec.ID) == 0 {
		instanceList := &provisioningv1alpha2.InstanceList{}
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
	instance.Spec.DataPlane = deployment.Spec.DataPlane

	c.Operation.Scheme().Default(instance)

	return nil
}

// findServiceTargetConfig tries to find a service target configuration that applies to the deployment requirements and has capacity available.
func (c *Controller) findServiceTargetConfig(ctx context.Context) (*provisioningv1alpha2.ServiceTargetConfig, error) {
	serviceTargetConfigs := &provisioningv1alpha2.ServiceTargetConfigList{}
	selectorBuilder := strings.Builder{}
	selectorBuilder.WriteString(fmt.Sprintf("%s=true", provisioning.ServiceTargetConfigVisibleLabelName))

	labelSelector, _ := labels.Parse(selectorBuilder.String())
	listOptions := client.ListOptions{
		LabelSelector: labelSelector,
	}

	if err := c.Client().List(ctx, serviceTargetConfigs, &listOptions); err != nil {
		return nil, fmt.Errorf("unable to list service target configs: %w", err)
	}

	SortServiceTargetConfigs(serviceTargetConfigs)

	if len(serviceTargetConfigs.Items) == 0 {
		err := fmt.Errorf("no service target with remaining capacity available")
		return nil, err
	}

	return &serviceTargetConfigs.Items[0], nil
}

// SortServiceTargetConfigs sorts all configs by priority and usage.
func SortServiceTargetConfigs(configs *provisioningv1alpha2.ServiceTargetConfigList) {
	if len(configs.Items) == 0 {
		return
	}

	// sort the configurations by priority and capacity
	sort.SliceStable(configs.Items, func(i, j int) bool {
		l := &configs.Items[i]
		r := &configs.Items[j]

		lPrio := l.Spec.Priority / int64(len(l.Status.InstanceRefs)+1)
		rPrio := r.Spec.Priority / int64(len(r.Status.InstanceRefs)+1)

		return lPrio > rPrio
	})
}
