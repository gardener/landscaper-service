// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package avuploader

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// AddControllerToManager adds the AVUploader controller to the manager
func AddControllerToManager(ctx context.Context, logger logr.Logger, mgr manager.Manager, config *coreconfig.LandscaperServiceConfiguration) error {
	log := logger.WithName("AVUploader")
	if config.AvailabilityMonitoring.AvailabilityServiceConfiguration == nil {
		log.Info("AvailabilityServiceConfiguration missing, not starting AVUploader")
		return nil
	}

	ctrl, err := NewController(log, mgr.GetClient(), mgr.GetScheme(), config)
	if err != nil {
		return err
	}

	return builder.ControllerManagedBy(mgr).
		For(&v1alpha1.AvailabilityCollection{}).
		WithLogConstructor(func(r *reconcile.Request) logr.Logger { return log }).
		Complete(ctrl)
}
