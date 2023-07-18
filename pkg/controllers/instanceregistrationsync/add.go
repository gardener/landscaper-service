// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instanceregistrationsync

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/controllers/instanceregistrationsync/instanceregistration"
	"github.com/gardener/landscaper-service/pkg/controllers/instanceregistrationsync/landscaperdeployment"
)

// AddControllerToManager adds the InstanceRegistration Controller to the manager
func AddInstanceRegistrationSyncControllerToManager(logger logging.Logger, mgr manager.Manager) error {
	log := logger.Reconciles("InstanceRegistrationController", "InstanceRegistration")
	ctrl, err := instanceregistration.NewController(log, mgr.GetClient(), mgr.GetScheme())
	if err != nil {
		return err
	}

	return builder.ControllerManagedBy(mgr).
		For(&v1alpha1.InstanceRegistration{}).
		WithLogConstructor(func(r *reconcile.Request) logr.Logger { return log.Logr() }).
		Complete(ctrl)
}

// AddLandscaperDeploymentSyncControllerToManager adds the AddLandscaperDeploymentSync Controller to the manager
func AddLandscaperDeploymentSyncControllerToManager(logger logging.Logger, mgr manager.Manager) error {
	log := logger.Reconciles("LandscaperDeploymentSyncController", "LandscaperDeployment")
	ctrl, err := landscaperdeployment.NewController(log, mgr.GetClient(), mgr.GetScheme())
	if err != nil {
		return err
	}

	return builder.ControllerManagedBy(mgr).
		For(&v1alpha1.InstanceRegistration{}). //TODO: restrict on what to react (only status)
		WithLogConstructor(func(r *reconcile.Request) logr.Logger { return log.Logr() }).
		Complete(ctrl)
}
