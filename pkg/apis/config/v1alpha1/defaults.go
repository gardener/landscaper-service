// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_LandscaperServiceConfiguration sets the defaults for the landscaper service configuration.
func SetDefaults_LandscaperServiceConfiguration(obj *LandscaperServiceConfiguration) {
	SetDefaults_CrdManagementConfiguration(&obj.CrdManagement)
	SetDefaults_AvailabilityMonitoringConfiguration(&obj.AvailabilityMonitoring)
}

// SetDefaults_CrdManagementConfiguration sets the defaults for the crd management configuration.
func SetDefaults_CrdManagementConfiguration(obj *CrdManagementConfiguration) {
	if obj.DeployCustomResourceDefinitions == nil {
		obj.DeployCustomResourceDefinitions = pointer.Bool(true)
	}
	if obj.ForceUpdate == nil {
		obj.ForceUpdate = pointer.Bool(true)
	}
}

// AvailabilityMonitoringConfiguration sets the defaults for the availability monitoring configuration.
func SetDefaults_AvailabilityMonitoringConfiguration(obj *AvailabilityMonitoringConfiguration) {
	if obj.AvailabilityCollectionName == "" {
		obj.AvailabilityCollectionName = "availability"
	}
	if obj.AvailabilityCollectionNamespace == "" {
		obj.AvailabilityCollectionNamespace = "laas-system"
	}
	if obj.SelfLandscaperNamespace == "" {
		obj.SelfLandscaperNamespace = "landscaper"
	}
	if obj.PeriodicCheckInterval.Duration == 0 {
		obj.PeriodicCheckInterval.Duration = time.Minute * 1
	}
	if obj.LSHealthCheckTimeout.Duration == 0 {
		obj.LSHealthCheckTimeout.Duration = time.Minute * 5
	}
	if obj.AvailabilityServiceConfiguration != nil {
		if obj.AvailabilityServiceConfiguration.Timeout == "" {
			obj.AvailabilityServiceConfiguration.Timeout = "30s"
		}
	}
}

// SetDefaults_ShootConfiguration sets the defaults for the shoot configuration.
func SetDefaults_ShootConfiguration(obj *ShootConfiguration) {
	maintenance := &obj.Maintenance

	if maintenance.AutoUpdate.KubernetesVersion == nil {
		obj.Maintenance.AutoUpdate.KubernetesVersion = pointer.Bool(false)
	}

	if maintenance.AutoUpdate.MachineImageVersion == nil {
		obj.Maintenance.AutoUpdate.MachineImageVersion = pointer.Bool(false)
	}
	workers := &obj.Workers

	if workers.Minimum == nil {
		workers.Minimum = pointer.Int32(1)
	}

	if workers.Maximum == nil {
		workers.Maximum = pointer.Int32(1)
	}

	if workers.MaxSurge == nil {
		workers.MaxSurge = pointer.Int32(1)
	}

	if workers.MaxUnavailable == nil {
		workers.MaxUnavailable = pointer.Int32(0)
	}
}
