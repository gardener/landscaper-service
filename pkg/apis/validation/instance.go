// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"fmt"

	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	lsscore "github.com/gardener/landscaper-service/pkg/apis/core"
)

const (
	// InstanceIdLength is the required length of the Instance id
	InstanceIdLength = 8
)

// ValidateInstance validates an instance
func ValidateInstance(instance *lsscore.Instance, oldInstance *lsscore.Instance) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateInstanceObjectMeta(&instance.ObjectMeta, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateInstanceSpec(&instance.Spec, field.NewPath("spec"))...)
	if oldInstance != nil {
		allErrs = append(allErrs, validateInstanceSpecUpdate(&instance.Spec, &oldInstance.Spec, field.NewPath("spec"))...)
	}
	return allErrs
}

func validateInstanceObjectMeta(objMeta *metav1.ObjectMeta, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, apivalidation.ValidateObjectMeta(objMeta, true, apivalidation.NameIsDNSLabel, fldPath)...)
	return allErrs
}

func validateInstanceSpec(spec *lsscore.InstanceSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateObjectReference(&spec.ServiceTargetConfigRef, fldPath.Child("serviceTargetConfigRef"))...)

	if len(spec.TenantId) != LandscaperDeploymentTenantIdLength {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("tenantId"), spec.TenantId, fmt.Sprintf("must be exactly of size %d", LandscaperDeploymentTenantIdLength)))
	}

	if len(spec.ID) != InstanceIdLength {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("id"), spec.TenantId, fmt.Sprintf("must be exactly of size %d", InstanceIdLength)))
	}

	if spec.HighAvailabilityConfig != nil {
		allErrs = append(allErrs, ValidateHighAvailabilityConfig(spec.HighAvailabilityConfig, fldPath.Child("highAvailabilityConfig"))...)
	}

	return allErrs
}

func validateInstanceSpecUpdate(spec *lsscore.InstanceSpec, oldSpec *lsscore.InstanceSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if spec.TenantId != oldSpec.TenantId {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("tenantId"), "is immutable"))
	}

	if spec.ID != oldSpec.ID {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("id"), "is immutable"))
	}

	if !spec.ServiceTargetConfigRef.Equals(&oldSpec.ServiceTargetConfigRef) {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("serviceTargetConfigRef"), "is immutable"))
	}

	if spec.HighAvailabilityConfig != nil && oldSpec.HighAvailabilityConfig != nil {
		if spec.HighAvailabilityConfig.ControlPlaneFailureTolerance != oldSpec.HighAvailabilityConfig.ControlPlaneFailureTolerance {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("highAvailabilityConfig").Child("controlPlaneFailureTolerance"), "is immutable"))
		}
	}

	return allErrs
}
