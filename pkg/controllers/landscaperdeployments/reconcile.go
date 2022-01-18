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

	"github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lsserrors "github.com/gardener/landscaper-service/pkg/apis/errors"
	"github.com/gardener/landscaper-service/pkg/utils"
)

// reconcile reconciles a landscaper deployment
func (c *Controller) reconcile(ctx context.Context, log logr.Logger, deployment *lssv1alpha1.LandscaperDeployment) error {
	currOp := "Reconcile"
	log.Info("Reconcile deployment", "name", deployment.GetName(), "namespace", deployment.GetNamespace())

	// reconcile instance
	instance := &lssv1alpha1.Instance{}
	instance.GenerateName = fmt.Sprintf("%s-", deployment.GetName())
	instance.Namespace = deployment.GetNamespace()
	if deployment.Status.InstanceRef != nil && !deployment.Status.InstanceRef.IsEmpty() {
		instance.Name = deployment.Status.InstanceRef.Name
		instance.Namespace = deployment.Status.InstanceRef.Namespace
	}

	_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), instance, func() error {
		return c.mutateInstance(ctx, log, deployment, instance)
	})

	if err != nil {
		return lsserrors.NewWrappedError(err, currOp, "CreateUpdateInstance", err.Error())
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

	// set the instance reference for the deployment if not already set
	if deployment.Status.InstanceRef == nil || !deployment.Status.InstanceRef.IsObject(instance) {
		deployment.Status.InstanceRef = &lssv1alpha1.ObjectReference{
			Name:      instance.GetName(),
			Namespace: instance.GetNamespace(),
		}

		if err := c.Client().Status().Update(ctx, deployment); err != nil {
			return lsserrors.NewWrappedError(err, currOp, "UpdateInstanceRefForDeployment", err.Error())
		}
	}
	return nil
}

// mutateInstance creates/updates the instance for a landscaper deployment
func (c *Controller) mutateInstance(ctx context.Context, log logr.Logger, deployment *lssv1alpha1.LandscaperDeployment, instance *lssv1alpha1.Instance) error {
	log.Info("Create/Update instance for deployment", "name", deployment.GetName(), "namespace", deployment.GetNamespace())

	if err := controllerutil.SetControllerReference(deployment, instance, c.Scheme()); err != nil {
		return fmt.Errorf("unable to set controller reference for instance: %w", err)
	}

	if len(instance.Spec.ServiceTargetConfigRef.Name) == 0 {
		// try to find a service target configuration that can be used for this landscaper deployment
		serviceTargetConf, err := c.findServiceTargetConfig(ctx, log, deployment)
		if err != nil {
			return err
		}

		instance.Spec.ServiceTargetConfigRef.Name = serviceTargetConf.GetName()
		instance.Spec.ServiceTargetConfigRef.Namespace = serviceTargetConf.GetNamespace()
	}

	instance.Spec.LandscaperConfiguration = deployment.Spec.LandscaperConfiguration
	c.Operation.Scheme().Default(instance)

	return nil
}

// findServiceTargetConfig tries to find a service target configuration that applies to the deployment requirements and has capacity available.
func (c *Controller) findServiceTargetConfig(ctx context.Context, log logr.Logger, deployment *lssv1alpha1.LandscaperDeployment) (*lssv1alpha1.ServiceTargetConfig, error) {
	serviceTargetConfigs := &lssv1alpha1.ServiceTargetConfigList{}
	selectorBuilder := strings.Builder{}
	selectorBuilder.WriteString(fmt.Sprintf("%s=true", lssv1alpha1.ServiceTargetConfigVisibleLabelName))

	if len(deployment.Spec.Region) > 0 {
		log.V(5).Info("region filter active", "region", deployment.Spec.Region)
		selectorBuilder.WriteString(fmt.Sprintf(",%s=%s", lssv1alpha1.ServiceTargetConfigRegionLabelName, deployment.Spec.Region))
	}

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
func SortServiceTargetConfigs(configs *lssv1alpha1.ServiceTargetConfigList) {
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
