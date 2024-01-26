// SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// ValidateTargetScheduling validates a TargetScheduling
func ValidateTargetScheduling(scheduling *v1alpha1.TargetScheduling) field.ErrorList {
	allErrs := field.ErrorList{}

	fldPath := field.NewPath("spec", "rules")
	for i := range scheduling.Spec.Rules {
		allErrs = append(allErrs, validateSchedulingRule(&scheduling.Spec.Rules[i], fldPath.Index(i))...)
	}

	return allErrs
}

func validateSchedulingRule(rule *v1alpha1.SchedulingRule, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if rule.Priority < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("priority"), rule.Priority, "priority must be an integer >= 0"))
	}

	configsPath := fldPath.Child("serviceTargetConfigs")
	if len(rule.ServiceTargetConfigs) == 0 {
		allErrs = append(allErrs, field.Required(configsPath, "at least one serviceTargetConfig is required"))
	}

	for i := range rule.ServiceTargetConfigs {
		config := &rule.ServiceTargetConfigs[i]
		configPath := configsPath.Index(i)

		if len(config.Name) == 0 {
			allErrs = append(allErrs, field.Required(configPath.Child("name"), "name needs to be set"))
		}
		if len(config.Namespace) == 0 {
			allErrs = append(allErrs, field.Required(configPath.Child("namespace"), "namespace needs to be set"))
		}
	}

	selectorPath := fldPath.Child("selector")
	for i := range rule.Selector {
		allErrs = append(allErrs, validateSchedulingSelector(&rule.Selector[i], selectorPath.Index(i))...)
	}

	return allErrs
}

func validateSchedulingSelector(selector *v1alpha1.Selector, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if selector == nil {
		return allErrs
	}

	// check that the selector has exactly one term
	count := 0
	if selector.MatchTenant != nil {
		count++
	}
	if selector.MatchLabel != nil {
		count++
	}
	if selector.Or != nil {
		count++
	}
	if selector.And != nil {
		count++
	}
	if selector.Not != nil {
		count++
	}

	if count != 1 {
		allErrs = append(allErrs, field.Invalid(fldPath, selector, "selector must contain exactly one term"))
	}

	// check term
	if selector.MatchTenant != nil {
		if len(selector.MatchTenant.ID) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("matchTenant").Child("id"), "tenant id needs to be set"))
		}

	} else if selector.MatchLabel != nil {
		if len(selector.MatchLabel.Name) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("matchLabel").Child("name"), "label name needs to be set"))
		}
		if len(selector.MatchLabel.Value) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("matchLabel").Child("value"), "label value needs to be set"))
		}

	} else if selector.And != nil {
		for i := range selector.And {
			allErrs = append(allErrs, validateSchedulingSelector(&selector.And[i], fldPath.Child("and").Index(i))...)
		}

	} else if selector.Or != nil {
		for i := range selector.Or {
			allErrs = append(allErrs, validateSchedulingSelector(&selector.Or[i], fldPath.Child("or").Index(i))...)
		}

	} else if selector.Not != nil {
		allErrs = append(allErrs, validateSchedulingSelector(selector.Not, fldPath.Child("not"))...)
	}

	return allErrs
}
