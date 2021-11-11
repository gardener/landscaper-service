// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package servicetargetconfigs

import (
	"context"

	"github.com/go-logr/logr"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/apis/errors"
)

var (
	initialCapacity int64 = 10
)

// reconcile reconciles a service target config.
func (c *Controller) reconcile(ctx context.Context, log logr.Logger, config *lssv1alpha1.ServiceTargetConfig) error {
	// initialize the capacity
	if config.Status.Capacity == nil {
		config.Status.Capacity = new(int64)
		*config.Status.Capacity = initialCapacity
	}

	// adjust the capacity based on the length of the instance reference list
	*config.Status.Capacity = initialCapacity - int64(len(config.Status.InstanceRefs))

	if err := c.Client().Status().Update(ctx, config); err != nil {
		log.Error(err, "unable to update capacity")
		return errors.NewWrappedError(err, "ReconcileServiceTargetConfig", "UpdateCapacity", err.Error())
	}

	return nil
}
