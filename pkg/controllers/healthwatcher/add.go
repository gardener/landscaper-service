// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package healthwatcher

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	config "github.com/gardener/landscaper-service/pkg/apis/config/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// AddControllerToManager adds the HealthWatcher controller to the manager
func AddControllerToManager(logger logging.Logger, mgr manager.Manager, config *config.LandscaperServiceConfiguration) error {
	log := logger.Reconciles("HealthWatcher", "AvailabilityCollection")
	ctrl, err := NewController(log, mgr.GetClient(), mgr.GetScheme(), config)
	if err != nil {
		return err
	}

	return builder.ControllerManagedBy(mgr).
		Named("health-watcher-controller").
		For(&v1alpha1.AvailabilityCollection{}).
		WithLogConstructor(func(r *reconcile.Request) logr.Logger { return log.Logr() }).
		Complete(ctrl)
}
