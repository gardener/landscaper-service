// SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package scheduling

import (
	"fmt"
	"sort"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// PickServiceTargetConfig selects one of the ServiceTargetConfigs, considering their priority and usage.
// For each ServiceTargetConfig, its priority is divided by the number of already deployed LandscaperDeployments + 1.
// The ServiceTargetConfigs are sorted descending by these numbers.
// The ServiceTargetConfig with the highest number is returned, i.e. the first in the sorted list.
func PickServiceTargetConfig(configs []*lssv1alpha1.ServiceTargetConfig) (*lssv1alpha1.ServiceTargetConfig, error) {
	if len(configs) == 0 {
		err := fmt.Errorf("no service target available")
		return nil, err
	}

	SortServiceTargetConfigs(configs)
	return configs[0], nil
}

// SortServiceTargetConfigs sorts the ServiceTargetConfigs by priority and usage.
func SortServiceTargetConfigs(configs []*lssv1alpha1.ServiceTargetConfig) {
	if len(configs) == 0 {
		return
	}

	// sort the configurations by priority and capacity
	sort.SliceStable(configs, func(i, j int) bool {
		l := configs[i]
		r := configs[j]

		lPrio := l.Spec.Priority / int64(len(l.Status.InstanceRefs)+1)
		rPrio := r.Spec.Priority / int64(len(r.Status.InstanceRefs)+1)

		return lPrio > rPrio
	})
}
