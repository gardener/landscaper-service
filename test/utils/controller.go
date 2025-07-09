// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"encoding/json"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"

	config "github.com/gardener/landscaper-service/pkg/apis/config/v1alpha1"

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
func ShouldReconcile(ctx context.Context, reconciler reconcile.Reconciler, req reconcile.Request, optionalDescription ...interface{}) reconcile.Result {
	res, err := reconciler.Reconcile(ctx, req)
	gomega.ExpectWithOffset(1, err).ToNot(gomega.HaveOccurred(), optionalDescription...)
	return res
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
			Name:    "github.com/gardener/landscaper-service/landscaper-instance",
			Version: "v1.1.1",
		},
		AvailabilityMonitoring: config.AvailabilityMonitoringConfiguration{
			AvailabilityCollectionName:      "availability",
			AvailabilityCollectionNamespace: "laas-system",
			SelfLandscaperNamespace:         "landscaper",
			PeriodicCheckInterval:           v1alpha1.Duration{Duration: time.Minute * 1},
			LSHealthCheckTimeout:            v1alpha1.Duration{Duration: time.Minute * 5},
		},
		GardenerConfiguration: config.GardenerConfiguration{
			ShootSecretBindingName: "secret-binding",
			ProjectName:            "test",
			ServiceAccountKubeconfig: v1alpha1.SecretReference{
				ObjectReference: v1alpha1.ObjectReference{
					Name:      "service-account",
					Namespace: "laas-system",
				},
				Key: "kubeconfig",
			},
		},
	}
	repositoryContext, err := json.Marshal(map[string]interface{}{
		"type":    "ociRegistry",
		"baseUrl": "europe-docker.pkg.dev/sap-gcp-cp-k8s-stable-hub/development",
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	cfg.LandscaperServiceComponent.RepositoryContext = v1alpha1.NewAnyJSON(repositoryContext)
	return cfg
}

func DefaultTargetShootConfiguration() *config.TargetShootSidecarConfiguration {
	cfg := &config.TargetShootSidecarConfiguration{}
	return cfg
}

func CreateServiceAccountSecret(ctx context.Context, client client.Client, c *config.LandscaperServiceConfiguration) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.GardenerConfiguration.ServiceAccountKubeconfig.Namespace,
		},
	}
	if err := client.Create(ctx, namespace); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.GardenerConfiguration.ServiceAccountKubeconfig.Name,
			Namespace: c.GardenerConfiguration.ServiceAccountKubeconfig.Namespace,
		},
		StringData: map[string]string{
			c.GardenerConfiguration.ServiceAccountKubeconfig.Key: "kubeconfigcontent",
		},
		Type: "Opaque",
	}
	if err := client.Create(ctx, secret); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}
