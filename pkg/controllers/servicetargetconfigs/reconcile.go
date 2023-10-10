// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package servicetargetconfigs

import (
	"context"

	lssv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha2"
)

// reconcile reconciles a service target config.
func (c *Controller) reconcile(_ context.Context, _ *lssv1alpha2.ServiceTargetConfig) error {
	return nil
}
