// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package operation

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Operation struct {
	log    logr.Logger
	client client.Client
	scheme *runtime.Scheme
}

func NewOperation(log logr.Logger, c client.Client, scheme *runtime.Scheme) *Operation {
	return &Operation{
		log:    log,
		client: c,
		scheme: scheme,
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
