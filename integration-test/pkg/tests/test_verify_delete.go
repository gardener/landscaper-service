// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"

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

func (r *VerifyDeleteRunner) verifyNamespace(namespace string) error {
	r.log.Info("verifying namespace", "name", namespace)

	pvcList := &corev1.PersistentVolumeClaimList{}
	if err := r.kclient.List(r.ctx, pvcList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list persistent volume claims in namespace %q: %w", namespace, err)
	}

	if len(pvcList.Items) > 0 {
		return fmt.Errorf("there are persistent volume claims existing in namespace %q", namespace)
	}

	podList := &corev1.PodList{}
	if err := r.kclient.List(r.ctx, podList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list pods in namespace %q: %w", namespace, err)
	}

	if len(podList.Items) > 0 {
		return fmt.Errorf("there are pods existing in namespace %q", namespace)
	}

	serviceList := &corev1.ServiceList{}
	if err := r.kclient.List(r.ctx, serviceList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list services in namespace %q: %w", namespace, err)
	}

	if len(serviceList.Items) > 0 {
		return fmt.Errorf("there are services existing in namespace %q", namespace)
	}

	return nil
}
