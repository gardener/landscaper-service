// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// ValidatorFromResourceType is a helper method that gets a resource type and returns the fitting validator
func ValidatorFromResourceType(log logr.Logger, kubeClient client.Client, scheme *runtime.Scheme, resource string) (GenericValidator, error) {
	abstrVal := newAbstractedValidator(log, kubeClient, scheme)
	var val GenericValidator
	if resource == "landscaperdeployments" {
		val = &LandscaperDeploymentValidator{abstrVal}
	} else if resource == "instances" {
		val = &InstanceValidator{abstrVal}
	} else if resource == "servicetargetconfigs" {
		val = &ServiceTargetConfigValidator{abstrVal}
	} else {
		return nil, fmt.Errorf("unable to find validator for resource type %q", resource)
	}
	return val, nil
}

type abstractValidator struct {
	Client  client.Client
	decoder runtime.Decoder
	log     logr.Logger
}

// newAbstractedValidator creates a new abstracted validator
func newAbstractedValidator(log logr.Logger, kubeClient client.Client, scheme *runtime.Scheme) abstractValidator {
	return abstractValidator{
		Client:  kubeClient,
		decoder: serializer.NewCodecFactory(scheme).UniversalDecoder(),
		log:     log,
	}
}

// GenericValidator is an abstraction interface that implements admission.Handler and contains additional setter functions for the fields
type GenericValidator interface {
	Handle(context.Context, admission.Request) admission.Response
}

// LANDSCAPER DEPLOYMENT

// LandscaperDeploymentValidator represents a validator for a LandscaperDeployment
type LandscaperDeploymentValidator struct{ abstractValidator }

// Handle handles a request to the webhook
func (iv *LandscaperDeploymentValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	return admission.Allowed("LandscaperDeployment is valid")
}

// INSTANCE

// InstanceValidator represents a validator for an Instance
type InstanceValidator struct{ abstractValidator }

// Handle handles a request to the webhook
func (iv *InstanceValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	return admission.Allowed("Instance is valid")
}

// SERVICE TARGET CONFIG

// ServiceTargetConfigValidator represents a validator for a ServiceTargetConfig
type ServiceTargetConfigValidator struct{ abstractValidator }

// Handle handles a request to the webhook
func (iv *ServiceTargetConfigValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	return admission.Allowed("ServiceTargetConfig is valid")
}
