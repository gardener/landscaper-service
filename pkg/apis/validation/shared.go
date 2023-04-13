// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	lsscore "github.com/gardener/landscaper-service/pkg/apis/core"
)

// ValidateSecretReference validates a secret reference
func ValidateSecretReference(ref *lsscore.SecretReference, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(ref.Key) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("key"), "key may not be empty"))
	}

	allErrs = append(allErrs, ValidateObjectReference(&ref.ObjectReference, fldPath)...)
	return allErrs
}

// ValidateObjectReference validates an object reference
func ValidateObjectReference(ref *lsscore.ObjectReference, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(ref.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "name may not be empty"))
	}

	if len(ref.Namespace) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("namespace"), "namespace may not be empty"))
	}

	return allErrs
}

func ValidateHighAvailabilityConfig(haConfig *lsscore.HighAvailabilityConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if haConfig.ControlPlaneFailureTolerance != "zone" && haConfig.ControlPlaneFailureTolerance != "node" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("controlPlaneFailureTolerance"), haConfig.ControlPlaneFailureTolerance, "allowed values: \"zone\", \"node\""))
	}

	return allErrs
}
