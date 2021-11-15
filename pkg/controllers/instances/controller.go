// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper-service/pkg/operation"
)

// Controller is the instances controller
type Controller struct {
	operation.Operation
}

// NewController returns a new instances controller
func NewController(log logr.Logger, c client.Client, scheme *runtime.Scheme) (reconcile.Reconciler, error) {
	ctrl := &Controller{}
	op := operation.NewOperation(log, c, scheme)
	ctrl.Operation = *op
	return ctrl, nil
}

// Reconcile reconciles requests for instances
func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := c.Log().WithValues("instance", req.NamespacedName.String())
	logger.V(5).Info("reconcile", "resource", req.NamespacedName)

	return reconcile.Result{}, nil
}
