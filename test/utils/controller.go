// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"encoding/json"
	"time"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	v1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"

	"github.com/gardener/landscaper-service/pkg/apis/config"

	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// RequestFromObject creates a new reconcile.for a object
func RequestFromObject(obj client.Object) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: obj.GetNamespace(),
			Name:      obj.GetName(),
		},
	}
}

// ShouldReconcile reconciles the given reconciler with the given request
// and expects that no error occurred
func ShouldReconcile(ctx context.Context, reconciler reconcile.Reconciler, req reconcile.Request, optionalDescription ...interface{}) {
	_, err := reconciler.Reconcile(ctx, req)
	gomega.ExpectWithOffset(1, err).ToNot(gomega.HaveOccurred(), optionalDescription...)
}

// ShouldNotReconcile reconciles the given reconciler with the given request
// and expects that an error occurred
func ShouldNotReconcile(ctx context.Context, reconciler reconcile.Reconciler, req reconcile.Request, optionalDescription ...interface{}) {
	_, err := reconciler.Reconcile(ctx, req)
	gomega.ExpectWithOffset(1, err).To(gomega.HaveOccurred(), optionalDescription...)
}

func DefaultControllerConfiguration() *config.LandscaperServiceConfiguration {
	cfg := &config.LandscaperServiceConfiguration{
		LandscaperServiceComponent: config.LandscaperServiceComponentConfiguration{
			Name:    "github.com/gardener/landscaper/landscaper-service",
			Version: "v1.1.1",
		},
		AvailabilityMonitoring: config.AvailabilityMonitoringConfiguration{
			AvailabilityCollectionName:      "availability",
			AvailabilityCollectionNamespace: "laas-system",
			SelfLandscaperNamespace:         "landscaper",
			PeriodicCheckInterval:           v1alpha1.Duration{Duration: time.Minute * 1},
			LSHealthCheckTimeout:            v1alpha1.Duration{Duration: time.Minute * 5},
		},
	}
	repositoryContext, err := json.Marshal(map[string]interface{}{
		"type":    "ociRegistry",
		"baseUrl": "eu.gcr.io/gardener-project/development",
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	cfg.LandscaperServiceComponent.RepositoryContext = lsv1alpha1.NewAnyJSON(repositoryContext)
	return cfg
}
