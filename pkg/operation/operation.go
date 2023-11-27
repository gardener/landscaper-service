// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package operation

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/landscaper-service/pkg/apis/config/v1alpha1"
)

// Operation is the base type for all controller types.
type Operation struct {
	// client is the kubernetes client instance
	client client.Client
	// scheme is the controller manager scheme used for serializing and deserializing objects.
	scheme *runtime.Scheme
	// config is the configuration for the landscaper service controller
	config *v1alpha1.LandscaperServiceConfiguration
}

// NewOperation creates a new Operation for the given values.
func NewOperation(c client.Client, scheme *runtime.Scheme, config *v1alpha1.LandscaperServiceConfiguration) *Operation {
	return &Operation{
		client: c,
		scheme: scheme,
		config: config,
	}
}

// Client returns a controller runtime client.Registry
func (o *Operation) Client() client.Client {
	return o.client
}

// Scheme returns a kubernetes scheme
func (o *Operation) Scheme() *runtime.Scheme {
	return o.scheme
}

func (o *Operation) Config() *v1alpha1.LandscaperServiceConfiguration {
	return o.config
}
