// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	lsscore "github.com/gardener/landscaper-service/pkg/apis/core"
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

	allErrs = append(allErrs, ValidateSecretReference(&spec.SecretRef, fldPath.Child("secretRef"))...)

	if len(spec.IngressDomain) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("ingressDomain"), "ingressDomain may not be empty"))
	}

	return allErrs
}
