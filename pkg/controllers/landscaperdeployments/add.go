// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package landscaperdeployments

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// AddControllerToManager adds the landscaperdeployments controller to the manager
func AddControllerToManager(ctx context.Context, logger logr.Logger, mgr manager.Manager) error {
	log := logger.WithName("LandscaperDeployments")
	ctrl, err := NewController(log, mgr.GetClient(), mgr.GetScheme())
	if err != nil {
		return err
	}

	return builder.ControllerManagedBy(mgr).
		For(&v1alpha1.LandscaperDeployment{}).
		Owns(&v1alpha1.LandscaperDeployment{}).
		Owns(&v1alpha1.Instance{}).
		WithLogger(log).
		Complete(ctrl)
}
