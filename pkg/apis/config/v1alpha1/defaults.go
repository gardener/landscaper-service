// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_LandscaperServiceConfiguration sets the defaults for the landscaper service configuration.
func SetDefaults_LandscaperServiceConfiguration(obj *LandscaperServiceConfiguration) {
	SetDefaults_CrdManagementConfiguration(&obj.CrdManagement)
}

// SetDefaults_CrdManagementConfiguration sets the defaults for the crd management configuration.
func SetDefaults_CrdManagementConfiguration(obj *CrdManagementConfiguration) {
	if obj.DeployCustomResourceDefinitions == nil {
		obj.DeployCustomResourceDefinitions = pointer.BoolPtr(true)
	}
	if obj.ForceUpdate == nil {
		obj.ForceUpdate = pointer.BoolPtr(true)
	}
}
