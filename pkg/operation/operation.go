// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package operation

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
)

// Operation is the base type for all controller types.
type Operation struct {
	// log is the logger instance.
	log logr.Logger
	// client is the kubernetes client instance
	client client.Client
	// scheme is the controller manager scheme used for serializing and deserializing objects.
	scheme *runtime.Scheme
	// config is the configuration for the landscaper service controller
	config *coreconfig.LandscaperServiceConfiguration
}

// NewOperation creates a new Operation for the given values.
func NewOperation(log logr.Logger, c client.Client, scheme *runtime.Scheme, config *coreconfig.LandscaperServiceConfiguration) *Operation {
	return &Operation{
		log:    log,
		client: c,
		scheme: scheme,
		config: config,
	}
}

// Log returns a logging instance
func (o *Operation) Log() logr.Logger {
	return o.log
}

// Client returns a controller runtime client.Registry
func (o *Operation) Client() client.Client {
	return o.client
}

// Scheme returns a kubernetes scheme
func (o *Operation) Scheme() *runtime.Scheme {
	return o.scheme
}

func (o *Operation) Config() *coreconfig.LandscaperServiceConfiguration {
	return o.config
}
