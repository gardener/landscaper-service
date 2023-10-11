// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gardener/landscaper-service/pkg/apis/provisioning"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gardener/landscaper-service/pkg/apis/validation"
)

const (
	LandscaperDeploymentsResourceType = "landscaperdeployments"
	InstancesResourceType             = "instances"
	ServiceTargetConfigsResourceType  = "servicetargetconfigs"
)

// ValidatorFromResourceType is a helper method that gets a resource type and returns the fitting validator
func ValidatorFromResourceType(log logging.Logger, kubeClient client.Client, scheme *runtime.Scheme, resource string) (GenericValidator, error) {
	abstrVal := newAbstractedValidator(log, kubeClient, scheme)
	var val GenericValidator
	if resource == LandscaperDeploymentsResourceType {
		val = &LandscaperDeploymentValidator{abstrVal}
	} else if resource == InstancesResourceType {
		val = &InstanceValidator{abstrVal}
	} else if resource == ServiceTargetConfigsResourceType {
		val = &ServiceTargetConfigValidator{abstrVal}
	} else {
		return nil, fmt.Errorf("unable to find validator for resource type %q", resource)
	}
	return val, nil
}

type abstractValidator struct {
	Client  client.Client
	decoder runtime.Decoder
	log     logging.Logger
}

// newAbstractedValidator creates a new abstracted validator
func newAbstractedValidator(log logging.Logger, kubeClient client.Client, scheme *runtime.Scheme) abstractValidator {
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
func (dv *LandscaperDeploymentValidator) Handle(_ context.Context, req admission.Request) admission.Response {
	deployment := &provisioning.LandscaperDeployment{}
	if _, _, err := dv.decoder.Decode(req.Object.Raw, nil, deployment); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	var oldDeployment *provisioning.LandscaperDeployment
	if req.Operation == admissionv1.Update && req.OldObject.Raw != nil {
		oldDeployment = &provisioning.LandscaperDeployment{}
		if _, _, err := dv.decoder.Decode(req.OldObject.Raw, nil, oldDeployment); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
	}

	if errs := validation.ValidateLandscaperDeployment(deployment, oldDeployment); len(errs) > 0 {
		return admission.Denied(errs.ToAggregate().Error())
	}

	return admission.Allowed("LandscaperDeployment is valid")
}

// INSTANCE

// InstanceValidator represents a validator for an Instance
type InstanceValidator struct{ abstractValidator }

// Handle handles a request to the webhook
func (iv *InstanceValidator) Handle(_ context.Context, req admission.Request) admission.Response {
	instance := &provisioning.Instance{}
	if _, _, err := iv.decoder.Decode(req.Object.Raw, nil, instance); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	var oldInstance *provisioning.Instance
	if req.Operation == admissionv1.Update && req.OldObject.Raw != nil {
		oldInstance = &provisioning.Instance{}
		if _, _, err := iv.decoder.Decode(req.OldObject.Raw, nil, oldInstance); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
	}

	if errs := validation.ValidateInstance(instance, oldInstance); len(errs) > 0 {
		return admission.Denied(errs.ToAggregate().Error())
	}

	return admission.Allowed("Instance is valid")
}

// SERVICE TARGET CONFIG

// ServiceTargetConfigValidator represents a validator for a ServiceTargetConfig
type ServiceTargetConfigValidator struct{ abstractValidator }

// Handle handles a request to the webhook
func (sv *ServiceTargetConfigValidator) Handle(_ context.Context, req admission.Request) admission.Response {
	config := &provisioning.ServiceTargetConfig{}
	if _, _, err := sv.decoder.Decode(req.Object.Raw, nil, config); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if errs := validation.ValidateServiceTargetConfig(config); len(errs) > 0 {
		return admission.Denied(errs.ToAggregate().Error())
	}

	return admission.Allowed("ServiceTargetConfig is valid")
}
