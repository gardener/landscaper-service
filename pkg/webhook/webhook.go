// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"context"
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/apis/validation"
)

const (
	LandscaperDeploymentsResourceType = "landscaperdeployments"
	InstancesResourceType             = "instances"
	ServiceTargetConfigsResourceType  = "servicetargetconfigs"
	TargetSchedulingsResourceType     = "targetschedulings"
)

// ValidatorFromResourceType is a helper method that gets a resource type and returns the fitting validator
func ValidatorFromResourceType(log logging.Logger, kubeClient client.Client, scheme *runtime.Scheme, resource string) (GenericValidator, error) {
	abstrVal := newAbstractedValidator(log, kubeClient, scheme)
	var val GenericValidator
	switch resource {
	case LandscaperDeploymentsResourceType:
		val = &LandscaperDeploymentValidator{abstrVal}
	case InstancesResourceType:
		val = &InstanceValidator{abstrVal}
	case ServiceTargetConfigsResourceType:
		val = &ServiceTargetConfigValidator{abstrVal}
	case TargetSchedulingsResourceType:
		val = &TargetSchedulingValidator{abstrVal}
	default:
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
	deployment := &lssv1alpha1.LandscaperDeployment{}
	if _, _, err := dv.decoder.Decode(req.Object.Raw, nil, deployment); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	var oldDeployment *lssv1alpha1.LandscaperDeployment
	if req.Operation == admissionv1.Update && req.OldObject.Raw != nil {
		oldDeployment = &lssv1alpha1.LandscaperDeployment{}
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
	instance := &lssv1alpha1.Instance{}
	if _, _, err := iv.decoder.Decode(req.Object.Raw, nil, instance); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	var oldInstance *lssv1alpha1.Instance
	if req.Operation == admissionv1.Update && req.OldObject.Raw != nil {
		oldInstance = &lssv1alpha1.Instance{}
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
	config := &lssv1alpha1.ServiceTargetConfig{}
	if _, _, err := sv.decoder.Decode(req.Object.Raw, nil, config); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if errs := validation.ValidateServiceTargetConfig(config); len(errs) > 0 {
		return admission.Denied(errs.ToAggregate().Error())
	}

	return admission.Allowed("ServiceTargetConfig is valid")
}

// TARGET SCHEDULING

type TargetSchedulingValidator struct{ abstractValidator }

// Handle handles a request to the webhook
func (sv *TargetSchedulingValidator) Handle(_ context.Context, req admission.Request) admission.Response {
	scheduling := &lssv1alpha1.TargetScheduling{}
	if _, _, err := sv.decoder.Decode(req.Object.Raw, nil, scheduling); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if errs := validation.ValidateTargetScheduling(scheduling); len(errs) > 0 {
		return admission.Denied(errs.ToAggregate().Error())
	}

	return admission.Allowed("TargetScheduling is valid")
}
