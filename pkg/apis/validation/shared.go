// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"slices"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// ValidateSecretReference validates a secret reference
func ValidateSecretReference(ref *v1alpha1.SecretReference, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(ref.Key) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("key"), "key may not be empty"))
	}

	allErrs = append(allErrs, ValidateObjectReference(&ref.ObjectReference, fldPath)...)
	return allErrs
}

// ValidateObjectReference validates an object reference
func ValidateObjectReference(ref *v1alpha1.ObjectReference, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(ref.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "name may not be empty"))
	}

	if len(ref.Namespace) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("namespace"), "namespace may not be empty"))
	}

	return allErrs
}

func ValidateHighAvailabilityConfig(haConfig *v1alpha1.HighAvailabilityConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if haConfig.ControlPlaneFailureTolerance != "zone" && haConfig.ControlPlaneFailureTolerance != "node" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("controlPlaneFailureTolerance"), haConfig.ControlPlaneFailureTolerance, "allowed values: \"zone\", \"node\""))
	}

	return allErrs
}

func ValidateDataPlane(dataPlane *v1alpha1.DataPlane, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if dataPlane.SecretRef == nil && len(dataPlane.Kubeconfig) == 0 {
		allErrs = append(allErrs, field.Forbidden(fldPath, "either secretRef or kubeconfig must be specified"))
	}

	if dataPlane.SecretRef != nil && len(dataPlane.Kubeconfig) > 0 {
		allErrs = append(allErrs, field.Forbidden(fldPath, "secretRef or kubeconfig must not be specified at the same time"))
	}

	if dataPlane.SecretRef != nil {
		allErrs = append(allErrs, ValidateSecretReference(dataPlane.SecretRef, fldPath.Child("secretRef"))...)
	}

	return allErrs
}

var supportedDeployers = []string{"helm", "manifest", "container"}

func ValidateLandscaperConfiguration(landscaperConfiguration *v1alpha1.LandscaperConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(landscaperConfiguration.Deployers) != 0 {
		for _, deployer := range landscaperConfiguration.Deployers {
			if !slices.Contains(supportedDeployers, deployer) {
				allErrs = append(allErrs, field.NotSupported(fldPath, deployer, supportedDeployers))
			}
		}
	}

	return allErrs
}
