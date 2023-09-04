// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package subjectsync

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// AddControllerToManager adds the SubjectList Controller to the manager
func AddControllerToManager(logger logging.Logger, mgr manager.Manager, config *coreconfig.TargetShootSidecarConfiguration) error {
	log := logger.Reconciles("SubjectSyncController", "SubjectList")
	ctrl, err := NewController(log, mgr.GetClient(), mgr.GetScheme(), config)
	if err != nil {
		return err
	}

	predicates := builder.WithPredicates(predicate.Or(predicate.LabelChangedPredicate{},
		predicate.GenerationChangedPredicate{}, predicate.AnnotationChangedPredicate{}))

	return builder.ControllerManagedBy(mgr).
		For(&v1alpha1.SubjectList{}, predicates).
		WithLogConstructor(func(r *reconcile.Request) logr.Logger { return log.Logr() }).
		Complete(ctrl)
}
