// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"

	cliutil "github.com/gardener/landscapercli/pkg/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"

	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type VerifyDeleteRunner struct {
	ctx         context.Context
	log         logr.Logger
	kclient     client.Client
	config      *test.TestConfig
	target      *lsv1alpha1.Target
	testObjects *test.SharedTestObjects
}

func (r *VerifyDeleteRunner) Init(
	ctx context.Context, log logr.Logger,
	kclient client.Client, config *test.TestConfig,
	target *lsv1alpha1.Target, testObjects *test.SharedTestObjects) {
	r.ctx = ctx
	r.log = log.WithName(r.Name())
	r.kclient = kclient
	r.config = config
	r.target = target
	r.testObjects = testObjects
}

func (r *VerifyDeleteRunner) Name() string {
	return "VerifyDelete"
}

func (r *VerifyDeleteRunner) Run() error {
	for _, namespace := range r.testObjects.HostingClusterNamespaces {
		if err := r.verifyNamespace(namespace); err != nil {
			return err
		}
	}
	return nil
}

func (r *VerifyDeleteRunner) verifyNamespace(namespaceName string) error {
	r.log.Info("verifying namespace", "name", namespaceName)

	timeout, err := cliutil.CheckConditionPeriodically(func() (bool, error) {
		namespace := &corev1.Namespace{}
		if err := r.kclient.Get(r.ctx, types.NamespacedName{Name: namespaceName}, namespace); err != nil {
			if k8serrors.IsNotFound(err) {
				return true, nil
			} else {
				return false, err
			}
		}

		return false, nil
	}, r.config.SleepTime, r.config.MaxRetries)

	if timeout {
		return fmt.Errorf("timeout while waiting for namespace %q being deleted", namespaceName)
	}
	if err != nil {
		return fmt.Errorf("error while waiting for namespace %q being deleted: %w", namespaceName, err)
	}

	return nil
}
