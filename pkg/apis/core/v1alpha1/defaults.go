// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultsFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_SeedConfig sets the default values for ServiceTargetConfig objects
func SetDefaults_SeedConfig(obj *ServiceTargetConfig) {
	if len(obj.Spec.SecretRef.Namespace) == 0 {
		obj.Spec.SecretRef.Namespace = obj.GetNamespace()
	}
}
