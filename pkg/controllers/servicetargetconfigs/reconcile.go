// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package servicetargetconfigs

import (
	"context"

	"github.com/go-logr/logr"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// reconcile reconciles a service target config.
func (c *Controller) reconcile(ctx context.Context, log logr.Logger, config *lssv1alpha1.ServiceTargetConfig) error {
	return nil
}
