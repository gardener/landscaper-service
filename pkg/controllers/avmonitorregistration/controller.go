// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package avmonitorregistration

import (
	"context"
	"fmt"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/operation"
)

type Controller struct {
	operation.Operation
}

func NewController(log logr.Logger, c client.Client, scheme *runtime.Scheme, config *coreconfig.LandscaperServiceConfiguration) (reconcile.Reconciler, error) {
	ctrl := &Controller{}
	op := operation.NewOperation(log, c, scheme, config)
	ctrl.Operation = *op
	return ctrl, nil
}

func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := c.Log().WithValues("instance", req.NamespacedName.String())
	ctx = logr.NewContext(ctx, log)
	log.V(5).Info("reconcile", "resource", req.NamespacedName)

	availabilityCollection := &lssv1alpha1.AvailabilityCollection{}
	availabilityCollection.Name = c.Operation.Config().AvailabilityMonitoring.AvailabilityCollectionName
	availabilityCollection.Namespace = c.Operation.Config().AvailabilityMonitoring.AvailabilityCollectionNamespace

	instances := &lssv1alpha1.InstanceList{}
	if err := c.Client().List(ctx, instances); err != nil {
		if apierrors.IsNotFound(err) {
			c.Log().V(5).Info(err.Error())
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	instanceRefsToMonitor := []lssv1alpha1.ObjectReference{}
	for _, instance := range instances.Items {
		//get refered installation
		if instance.Status.InstallationRef == nil || instance.Status.InstallationRef.Name == "" || instance.Status.InstallationRef.Namespace == "" {
			continue
		}
		//get installation
		installation := &lsv1alpha1.Installation{}
		if err := c.Client().Get(ctx, types.NamespacedName{Name: instance.Status.InstallationRef.Name, Namespace: instance.Status.InstallationRef.Namespace}, installation); err != nil {
			if apierrors.IsNotFound(err) {
				c.Log().V(5).Info(err.Error())
				continue
			}
			c.Log().V(5).Info(fmt.Sprintf("could not load installation from installation reference: %w", err))
			continue
		}
		//check if installation not progressing
		if installation.Status.Phase == lsv1alpha1.ComponentPhaseProgressing {
			c.Log().V(5).Info(fmt.Sprintf("installation %s:%s for instance %s:%s is progressing, not health check monitoring", installation.Namespace, installation.Name, instance.Namespace, instance.Name))
			continue
		}

		instanceRefsToMonitor = append(instanceRefsToMonitor, lssv1alpha1.ObjectReference{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		})
	}
	availabilityCollection.Spec = lssv1alpha1.AvailabilityCollectionSpec{
		InstanceRefs: instanceRefsToMonitor,
	}

	_, err := kubernetes.CreateOrUpdate(ctx, c.Client(), availabilityCollection, func() error {
		return nil
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
