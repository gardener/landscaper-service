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
	// LandscaperDeploymentTenantIdLength is the required length of the LandscaperDeployment tenant id
	LandscaperDeploymentTenantIdLength = 8
)

// ValidateLandscaperDeployment validates a LandscaperDeployment
func ValidateLandscaperDeployment(deployment *lsscore.LandscaperDeployment, oldDeployment *lsscore.LandscaperDeployment) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateLandscaperDeploymentObjectMeta(&deployment.ObjectMeta, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateLandscaperDeploymentSpec(&deployment.Spec, field.NewPath("spec"))...)
	if oldDeployment != nil {
		allErrs = append(allErrs, validateLandscaperDeploymentSpecUpdate(&deployment.Spec, &oldDeployment.Spec, field.NewPath("spec"))...)
	}
	return allErrs
}

func validateLandscaperDeploymentObjectMeta(objMeta *metav1.ObjectMeta, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, apivalidation.ValidateObjectMeta(objMeta, true, apivalidation.NameIsDNSLabel, fldPath)...)
	return allErrs
}

func validateLandscaperDeploymentSpec(spec *lsscore.LandscaperDeploymentSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(spec.TenantId) != LandscaperDeploymentTenantIdLength {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("tenantId"), spec.TenantId, fmt.Sprintf("must be exactly of size %d", LandscaperDeploymentTenantIdLength)))
	}

	if len(spec.Purpose) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("purpose"), "purpose may not be empty"))
	}

	if spec.HighAvailabilityConfig != nil {
		allErrs = append(allErrs, ValidateHighAvailabilityConfig(spec.HighAvailabilityConfig, fldPath.Child("highAvailabilityConfig"))...)
	}

	return allErrs
}

func validateLandscaperDeploymentSpecUpdate(spec *lsscore.LandscaperDeploymentSpec, oldSpec *lsscore.LandscaperDeploymentSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if spec.TenantId != oldSpec.TenantId {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("tenantId"), "is immutable"))
	}

	if spec.HighAvailabilityConfig != nil && oldSpec.HighAvailabilityConfig != nil {
		if spec.HighAvailabilityConfig.ControlPlaneFailureTolerance != oldSpec.HighAvailabilityConfig.ControlPlaneFailureTolerance {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("highAvailabilityConfig").Child("controlPlaneFailureTolerance"), "is immutable"))
		}
	}

	return allErrs
}
