// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

var (
	// AutomaticReconcileDefaultDuration specifies the default automatic reconcile duration.
	AutomaticReconcileDefaultDuration = 12 * time.Hour
)

type automaticReconcileHelper struct {
	cl    client.Client
	clock clock.PassiveClock
}

func newAutomaticReconcileHelper(cl client.Client, passiveClock clock.PassiveClock) *automaticReconcileHelper {
	return &automaticReconcileHelper{
		cl:    cl,
		clock: passiveClock,
	}
}

func (r *automaticReconcileHelper) computeAutomaticReconcile(ctx context.Context, instance *lssv1alpha1.Instance, reconcileError error) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	if err := r.updateReconcileStatus(ctx, instance); err != nil {
		logger.Error(err, "failed to update instance status")
		return reconcile.Result{}, err
	}

	if reconcileError == nil {
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: r.getReconcileInterval(instance),
		}, nil
	} else {
		return reconcile.Result{}, reconcileError
	}
}

func (r *automaticReconcileHelper) updateReconcileStatus(ctx context.Context, instance *lssv1alpha1.Instance) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	instance.Status.AutomaticReconcileStatus = &lssv1alpha1.AutomaticReconcileStatus{
		LastReconcileTime: r.metaNow(),
	}

	if err := r.cl.Status().Update(ctx, instance); err != nil {
		logger.Error(err, "failed to update instance status")
		return err
	}

	return nil
}

func (r *automaticReconcileHelper) getReconcileInterval(instance *lssv1alpha1.Instance) time.Duration {
	duration := AutomaticReconcileDefaultDuration
	if instance.Spec.AutomaticReconcile != nil {
		duration = instance.Spec.AutomaticReconcile.Interval.Duration
	}

	return duration
}

func (r *automaticReconcileHelper) now() time.Time {
	return r.clock.Now()
}

func (r *automaticReconcileHelper) metaNow() metav1.Time {
	return metav1.Time{Time: r.now()}
}
