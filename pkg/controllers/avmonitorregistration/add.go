// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package avmonitorregistration

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
)

// AddControllerToManager adds the AvailabilityMonitorRegistrationController to the manager
func AddControllerToManager(logger logging.Logger, mgr manager.Manager, config *coreconfig.LandscaperServiceConfiguration) error {
	log := logger.Reconciles("AvailabilityMonitorRegistrationController", "AvailabilityCollection")
	ctrl, err := NewController(log, mgr.GetClient(), mgr.GetScheme(), config)
	if err != nil {
		return err
	}

	return builder.ControllerManagedBy(mgr).
		For(&v1alpha1.Instance{}).
		Owns(&v1alpha1.AvailabilityCollection{}).
		WithLogConstructor(func(r *reconcile.Request) logr.Logger { return log.Logr() }).
		Complete(ctrl)
}
