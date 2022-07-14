// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"fmt"

	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	lsscore "github.com/gardener/landscaper-service/pkg/apis/core"
)

var (
	// AllowedProviderTypes specifies the allowed provider types.
	AllowedProviderTypes = sets.NewString(
		"alicloud",
		"aws",
		"gcp",
	)
)

// ValidateServiceTargetConfig validates a ServiceTargetConfig
func ValidateServiceTargetConfig(config *lsscore.ServiceTargetConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateServiceTargetConfigObjectMeta(&config.ObjectMeta, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateServiceTargetConfigSpec(&config.Spec, field.NewPath("spec"))...)
	return allErrs
}

func validateServiceTargetConfigObjectMeta(objMeta *metav1.ObjectMeta, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, apivalidation.ValidateObjectMeta(objMeta, true, apivalidation.NameIsDNSLabel, fldPath)...)

	labelsPath := fldPath.Child("labels")

	regionLabelPath := labelsPath.Child(lsscore.ServiceTargetConfigRegionLabelName)
	regionLabelValue, ok := objMeta.Labels[lsscore.ServiceTargetConfigRegionLabelName]
	if !ok {
		allErrs = append(allErrs, field.Required(regionLabelPath, "label needs to be set"))
	} else if len(regionLabelValue) == 0 {
		allErrs = append(allErrs, field.Required(regionLabelPath, "label value may not be empty"))
	}

	visibleLabelPath := labelsPath.Child(lsscore.ServiceTargetConfigVisibleLabelName)
	visibleLabelValue, ok := objMeta.Labels[lsscore.ServiceTargetConfigVisibleLabelName]
	if !ok {
		allErrs = append(allErrs, field.Required(visibleLabelPath, "label needs to be set"))
	} else if visibleLabelValue != "true" && visibleLabelValue != "false" {
		allErrs = append(allErrs, field.Invalid(visibleLabelPath, visibleLabelValue, "invalid label value, allowed values: \"true\", \"false\""))
	}

	return allErrs
}

func validateServiceTargetConfigSpec(spec *lsscore.ServiceTargetConfigSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(spec.ProviderType) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("providerType"), "providerType may not be empty"))
	} else if !AllowedProviderTypes.Has(spec.ProviderType) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("providerType"), spec.ProviderType, fmt.Sprintf("providerType must be one of the following: %v", AllowedProviderTypes.List())))
	}

	allErrs = append(allErrs, ValidateSecretReference(&spec.SecretRef, fldPath.Child("secretRef"))...)

	return allErrs
}
