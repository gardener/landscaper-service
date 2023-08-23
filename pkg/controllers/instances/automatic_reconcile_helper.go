// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"context"
	"reflect"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

const (
	// AutomaticReconcileDefaultDuration specifies the default automatic reconcile duration.
	AutomaticReconcileDefaultDuration = 12 * time.Hour
)

type AutomaticReconcileHelper struct {
	cl    client.Client
	clock clock.PassiveClock
}

func NewAutomaticReconcileHelper(cl client.Client, passiveClock clock.PassiveClock) *AutomaticReconcileHelper {
	return &AutomaticReconcileHelper{
		cl:    cl,
		clock: passiveClock,
	}
}

func (r *AutomaticReconcileHelper) ComputeAutomaticReconcile(ctx context.Context, instance, oldInstance *lssv1alpha1.Instance, reconcileError error) (reconcile.Result, error) {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	// before modifying any objects, calculate the interval for the next automatic reconcile run
	reconcileInterval := r.getReconcileInterval(instance)

	// setting the automatic reconcile status to nil before comparing the old instance and the
	// potentially changed new instance
	oldInstance.Status.AutomaticReconcileStatus = nil
	automaticReconcileStatus := instance.Status.AutomaticReconcileStatus
	instance.Status.AutomaticReconcileStatus = nil

	instanceChanged := !reflect.DeepEqual(oldInstance, instance)
	instance.Status.AutomaticReconcileStatus = automaticReconcileStatus

	// when the instance has changed, update the last reconcile timestamp
	if instanceChanged {
		if err := r.updateReconcileStatus(ctx, instance); err != nil {
			logger.Error(err, "failed to update instance status")
			return reconcile.Result{}, err
		}
	}

	if reconcileError == nil {
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: reconcileInterval,
		}, nil
	} else {
		return reconcile.Result{}, reconcileError
	}
}

func (r *AutomaticReconcileHelper) updateReconcileStatus(ctx context.Context, instance *lssv1alpha1.Instance) error {
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

func (r *AutomaticReconcileHelper) getReconcileInterval(instance *lssv1alpha1.Instance) time.Duration {
	duration := AutomaticReconcileDefaultDuration
	if instance.Spec.AutomaticReconcile != nil {
		duration = instance.Spec.AutomaticReconcile.Interval.Duration
	}

	// when there is a previous last reconcile timestamp, re-calculate the duration until the next
	// automatic reconcile run
	if instance.Status.AutomaticReconcileStatus != nil {
		durationCalculated := duration - r.now().Sub(instance.Status.AutomaticReconcileStatus.LastReconcileTime.Time)
		if durationCalculated > 0 {
			duration = durationCalculated
		}
	}

	return duration
}

func (r *AutomaticReconcileHelper) now() time.Time {
	return r.clock.Now()
}

func (r *AutomaticReconcileHelper) metaNow() metav1.Time {
	return metav1.Time{Time: r.now()}
}
