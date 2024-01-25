// SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package scheduling

import (
	"fmt"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/utils"
)

func FindServiceTargetConfig(
	scheduling *lssv1alpha1.TargetScheduling,
	deployment *lssv1alpha1.LandscaperDeployment,
	serviceTargetConfigs []lssv1alpha1.ServiceTargetConfig) (*lssv1alpha1.ServiceTargetConfig, error) {

	// Find the ServiceTargetConfigs which match the deployment according to the scheduling rules.
	configRefs := make([]lssv1alpha1.ObjectReference, 0)
	if scheduling != nil {
		var err error
		configRefs, err = evaluateRules(scheduling, deployment)
		if err != nil {
			return nil, err
		}
	}

	// If scheduling is not configured, or no scheduling rules match, there are no configRefs so far.
	// In this case, we continue with the unrestricted ServiceTargetConfigs.
	if len(configRefs) == 0 {
		configRefs = getUnrestricted(serviceTargetConfigs)
	}

	// Remove duplicates and not existing ServiceTargetConfigs.
	configs := convertAndFilter(configRefs, serviceTargetConfigs)
	if len(configs) == 0 {
		err := fmt.Errorf("no service target config available")
		return nil, err
	}

	// Pick one of the ServiceTargetConfigs.
	return PickServiceTargetConfig(configs)
}

func evaluateRules(
	scheduling *lssv1alpha1.TargetScheduling,
	deployment *lssv1alpha1.LandscaperDeployment,
) ([]lssv1alpha1.ObjectReference, error) {

	var highestFoundPrio int64 = -1
	candidates := make([]lssv1alpha1.ObjectReference, 0)

	for i := range scheduling.Spec.Rules {
		rule := &scheduling.Spec.Rules[i]

		if len(rule.ServiceTargetConfigs) == 0 {
			return nil, fmt.Errorf("rule must contain at least one service target config")
		}
		if rule.Priority < 0 {
			return nil, fmt.Errorf("rule priority must not be negative")
		}

		if rule.Priority < highestFoundPrio {
			// we have already found a candidate with a higher prio
			continue
		}

		match, err := EvaluateSelectorList(rule.Selector, deployment)
		if err != nil {
			return nil, err
		}

		if !match {
			// rule does not match
			continue
		}

		if rule.Priority > highestFoundPrio {
			// rule has higher prio: replace candidates
			highestFoundPrio = rule.Priority
			candidates = rule.ServiceTargetConfigs
		} else {
			// rule has same prio: append candidates
			candidates = append(candidates, rule.ServiceTargetConfigs...)
		}
	}

	return candidates, nil
}

func getUnrestricted(serviceTargetConfigs []lssv1alpha1.ServiceTargetConfig) []lssv1alpha1.ObjectReference {
	result := make([]lssv1alpha1.ObjectReference, 0)

	for i := range serviceTargetConfigs {
		serviceTargetConfig := &serviceTargetConfigs[i]
		if !serviceTargetConfig.Spec.Restricted {
			ref := lssv1alpha1.ObjectReference{
				Name:      serviceTargetConfig.Name,
				Namespace: serviceTargetConfig.Namespace,
			}
			result = append(result, ref)
		}
	}

	return result
}

// convertAndFilter converts ObjectReferences to ServiceTargetConfigs.
// It skips duplicates and ObjectReferences for which there exists no ServiceTargetConfig.
func convertAndFilter(
	configRefs []lssv1alpha1.ObjectReference,
	serviceTargetConfigs []lssv1alpha1.ServiceTargetConfig,
) []*lssv1alpha1.ServiceTargetConfig {

	m := map[lssv1alpha1.ObjectReference]*lssv1alpha1.ServiceTargetConfig{}

	for _, ref := range configRefs {
		for k := range serviceTargetConfigs {
			serviceTargetConfig := &serviceTargetConfigs[k]
			if ref.Name == serviceTargetConfig.Name && ref.Namespace == serviceTargetConfig.Namespace {
				m[ref] = serviceTargetConfig
			}
		}
	}

	return utils.GetMapValues(m)
}
