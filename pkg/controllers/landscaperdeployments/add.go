// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployments

import (
	"github.com/go-logr/logr"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	config "github.com/gardener/landscaper-service/pkg/apis/config/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// AddControllerToManager adds the landscaperdeployments controller to the manager
func AddControllerToManager(logger logging.Logger, mgr manager.Manager, config *config.LandscaperServiceConfiguration) error {
	log := logger.Reconciles("landscaperDeployments", "LandscaperDeployments")
	ctrl, err := NewController(log, mgr.GetClient(), mgr.GetScheme(), config)
	if err != nil {
		return err
	}

	return builder.ControllerManagedBy(mgr).
		Named("landscaper-deployment-controller").
		For(&v1alpha1.LandscaperDeployment{}).
		Owns(&v1alpha1.LandscaperDeployment{}).
		Owns(&v1alpha1.Instance{}).
		WithLogConstructor(func(r *reconcile.Request) logr.Logger { return log.Logr() }).
		Complete(ctrl)
}
