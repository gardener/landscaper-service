// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package lossubject

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// AddControllerToManager adds the TenantRegistration Controller to the manager
func AddControllerToManager(logger logging.Logger, mgr manager.Manager) error {
	log := logger.Reconciles("LosSubjectController", "LosSubjectList")
	ctrl, err := NewController(log, mgr.GetClient(), mgr.GetScheme())
	if err != nil {
		return err
	}

	return builder.ControllerManagedBy(mgr).
		For(&v1alpha1.LosSubjectList{}).
		WithLogConstructor(func(r *reconcile.Request) logr.Logger { return log.Logr() }).
		Complete(ctrl)
}
