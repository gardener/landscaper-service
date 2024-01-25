// SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package scheduling

import (
	"fmt"

	"github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/utils"
)

func EvaluateSelectorList(selectors []v1alpha1.Selector, deployment *v1alpha1.LandscaperDeployment) (bool, error) {
	if deployment == nil {
		return false, fmt.Errorf("cannot evaluate selector list: no landscaper deployment specified")
	}

	return evaluateAnd(selectors, deployment)
}

func EvaluateSelector(selector *v1alpha1.Selector, deployment *v1alpha1.LandscaperDeployment) (bool, error) {
	if deployment == nil {
		return false, fmt.Errorf("cannot evaluate selector: no landscaper deployment specified")
	}
	if err := validateSelector(selector); err != nil {
		return false, err
	}

	if selector.MatchTenant != nil {
		return evaluateTenantSelector(selector.MatchTenant, deployment)
	} else if selector.MatchLabel != nil {
		return evaluateLabelSelector(selector.MatchLabel, deployment)
	} else if len(selector.Or) > 0 {
		return evaluateOr(selector.Or, deployment)
	} else if len(selector.And) > 0 {
		return evaluateAnd(selector.And, deployment)
	} else if selector.Not != nil {
		return evaluateNot(selector.Not, deployment)
	}

	return false, fmt.Errorf("selector is empty")
}

func evaluateTenantSelector(selector *v1alpha1.TenantSelector, deployment *v1alpha1.LandscaperDeployment) (bool, error) {
	return selector.ID == deployment.Spec.TenantId, nil
}

// TODO Do we want to support label selectors which check only the existence of the label, ignoring the value?
func evaluateLabelSelector(labelSelector *v1alpha1.LabelSelector, deployment *v1alpha1.LandscaperDeployment) (bool, error) {
	return utils.HasLabelWithValue(&deployment.ObjectMeta, labelSelector.Name, labelSelector.Value), nil
}

func evaluateOr(selectors []v1alpha1.Selector, deployment *v1alpha1.LandscaperDeployment) (bool, error) {
	for i := range selectors {
		selector := &selectors[i]
		match, err := EvaluateSelector(selector, deployment)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}

	return false, nil
}

func evaluateAnd(selectors []v1alpha1.Selector, deployment *v1alpha1.LandscaperDeployment) (bool, error) {
	for i := range selectors {
		selector := &selectors[i]
		match, err := EvaluateSelector(selector, deployment)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}

	return true, nil
}

func evaluateNot(selector *v1alpha1.Selector, deployment *v1alpha1.LandscaperDeployment) (bool, error) {
	match, err := EvaluateSelector(selector, deployment)
	if err != nil {
		return false, err
	}

	return !match, nil
}

func validateSelector(selector *v1alpha1.Selector) error {
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
		return fmt.Errorf("cannot evaluate selector: selector must contain exactly one condition")
	}

	return nil
}
