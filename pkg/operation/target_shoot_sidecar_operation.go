// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package operation

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
)

// Operation is the base type for all controller types.
type TargetShootSidecarOperation struct {
	// client is the kubernetes client instance
	client client.Client
	// scheme is the controller manager scheme used for serializing and deserializing objects.
	scheme *runtime.Scheme
	// config is the configuration for the landscaper service controller
	config *coreconfig.TargetShootSidecarConfiguration
}

// NewTargetShootSidecarOperation creates a new TargetShootSidecarOperation for the given values.
func NewTargetShootSidecarOperation(c client.Client, scheme *runtime.Scheme, config *coreconfig.TargetShootSidecarConfiguration) *TargetShootSidecarOperation {
	return &TargetShootSidecarOperation{
		client: c,
		scheme: scheme,
		config: config,
	}
}

// Client returns a controller runtime client.Registry
func (o *TargetShootSidecarOperation) Client() client.Client {
	return o.client
}

// Scheme returns a kubernetes scheme
func (o *TargetShootSidecarOperation) Scheme() *runtime.Scheme {
	return o.scheme
}

func (o *TargetShootSidecarOperation) Config() *coreconfig.TargetShootSidecarConfiguration {
	return o.config
}
