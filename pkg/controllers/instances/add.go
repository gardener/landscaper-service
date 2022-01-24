// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// AddControllerToManager adds the instances controller to the manager
func AddControllerToManager(logger logr.Logger, mgr manager.Manager, config *coreconfig.LandscaperServiceConfiguration) error {
	log := logger.WithName("Instances")
	ctrl, err := NewController(log, mgr.GetClient(), mgr.GetScheme(), config)
	if err != nil {
		return err
	}

	return builder.ControllerManagedBy(mgr).
		For(&lssv1alpha1.Instance{}).
		Owns(&lsv1alpha1.Installation{}).
		Owns(&lsv1alpha1.Target{}).
		WithLogger(log).
		Complete(ctrl)
}
