//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

SPDX-License-Identifier: Apache-2.0
*/
// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	v1 "k8s.io/api/core/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"

	config "github.com/gardener/landscaper-service/pkg/apis/config"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*AvailabilityMonitoringConfiguration)(nil), (*config.AvailabilityMonitoringConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AvailabilityMonitoringConfiguration_To_config_AvailabilityMonitoringConfiguration(a.(*AvailabilityMonitoringConfiguration), b.(*config.AvailabilityMonitoringConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.AvailabilityMonitoringConfiguration)(nil), (*AvailabilityMonitoringConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AvailabilityMonitoringConfiguration_To_v1alpha1_AvailabilityMonitoringConfiguration(a.(*config.AvailabilityMonitoringConfiguration), b.(*AvailabilityMonitoringConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AvailabilityServiceConfiguration)(nil), (*config.AvailabilityServiceConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AvailabilityServiceConfiguration_To_config_AvailabilityServiceConfiguration(a.(*AvailabilityServiceConfiguration), b.(*config.AvailabilityServiceConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.AvailabilityServiceConfiguration)(nil), (*AvailabilityServiceConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AvailabilityServiceConfiguration_To_v1alpha1_AvailabilityServiceConfiguration(a.(*config.AvailabilityServiceConfiguration), b.(*AvailabilityServiceConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CrdManagementConfiguration)(nil), (*config.CrdManagementConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CrdManagementConfiguration_To_config_CrdManagementConfiguration(a.(*CrdManagementConfiguration), b.(*config.CrdManagementConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.CrdManagementConfiguration)(nil), (*CrdManagementConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_CrdManagementConfiguration_To_v1alpha1_CrdManagementConfiguration(a.(*config.CrdManagementConfiguration), b.(*CrdManagementConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LandscaperServiceComponentConfiguration)(nil), (*config.LandscaperServiceComponentConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LandscaperServiceComponentConfiguration_To_config_LandscaperServiceComponentConfiguration(a.(*LandscaperServiceComponentConfiguration), b.(*config.LandscaperServiceComponentConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.LandscaperServiceComponentConfiguration)(nil), (*LandscaperServiceComponentConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_LandscaperServiceComponentConfiguration_To_v1alpha1_LandscaperServiceComponentConfiguration(a.(*config.LandscaperServiceComponentConfiguration), b.(*LandscaperServiceComponentConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LandscaperServiceConfiguration)(nil), (*config.LandscaperServiceConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LandscaperServiceConfiguration_To_config_LandscaperServiceConfiguration(a.(*LandscaperServiceConfiguration), b.(*config.LandscaperServiceConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.LandscaperServiceConfiguration)(nil), (*LandscaperServiceConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_LandscaperServiceConfiguration_To_v1alpha1_LandscaperServiceConfiguration(a.(*config.LandscaperServiceConfiguration), b.(*LandscaperServiceConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MetricsConfiguration)(nil), (*config.MetricsConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MetricsConfiguration_To_config_MetricsConfiguration(a.(*MetricsConfiguration), b.(*config.MetricsConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.MetricsConfiguration)(nil), (*MetricsConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_MetricsConfiguration_To_v1alpha1_MetricsConfiguration(a.(*config.MetricsConfiguration), b.(*MetricsConfiguration), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_AvailabilityMonitoringConfiguration_To_config_AvailabilityMonitoringConfiguration(in *AvailabilityMonitoringConfiguration, out *config.AvailabilityMonitoringConfiguration, s conversion.Scope) error {
	out.AvailabilityCollectionName = in.AvailabilityCollectionName
	out.AvailabilityCollectionNamespace = in.AvailabilityCollectionNamespace
	if err := Convert_v1alpha1_AvailabilityServiceConfiguration_To_config_AvailabilityServiceConfiguration(&in.AvailabilityServiceConfiguration, &out.AvailabilityServiceConfiguration, s); err != nil {
		return err
	}
	out.SelfLandscaperNamespace = in.SelfLandscaperNamespace
	out.PeriodicCheckInterval = in.PeriodicCheckInterval
	out.LSHealthCheckTimeout = in.LSHealthCheckTimeout
	return nil
}

// Convert_v1alpha1_AvailabilityMonitoringConfiguration_To_config_AvailabilityMonitoringConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_AvailabilityMonitoringConfiguration_To_config_AvailabilityMonitoringConfiguration(in *AvailabilityMonitoringConfiguration, out *config.AvailabilityMonitoringConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_AvailabilityMonitoringConfiguration_To_config_AvailabilityMonitoringConfiguration(in, out, s)
}

func autoConvert_config_AvailabilityMonitoringConfiguration_To_v1alpha1_AvailabilityMonitoringConfiguration(in *config.AvailabilityMonitoringConfiguration, out *AvailabilityMonitoringConfiguration, s conversion.Scope) error {
	out.AvailabilityCollectionName = in.AvailabilityCollectionName
	out.AvailabilityCollectionNamespace = in.AvailabilityCollectionNamespace
	if err := Convert_config_AvailabilityServiceConfiguration_To_v1alpha1_AvailabilityServiceConfiguration(&in.AvailabilityServiceConfiguration, &out.AvailabilityServiceConfiguration, s); err != nil {
		return err
	}
	out.SelfLandscaperNamespace = in.SelfLandscaperNamespace
	out.PeriodicCheckInterval = in.PeriodicCheckInterval
	out.LSHealthCheckTimeout = in.LSHealthCheckTimeout
	return nil
}

// Convert_config_AvailabilityMonitoringConfiguration_To_v1alpha1_AvailabilityMonitoringConfiguration is an autogenerated conversion function.
func Convert_config_AvailabilityMonitoringConfiguration_To_v1alpha1_AvailabilityMonitoringConfiguration(in *config.AvailabilityMonitoringConfiguration, out *AvailabilityMonitoringConfiguration, s conversion.Scope) error {
	return autoConvert_config_AvailabilityMonitoringConfiguration_To_v1alpha1_AvailabilityMonitoringConfiguration(in, out, s)
}

func autoConvert_v1alpha1_AvailabilityServiceConfiguration_To_config_AvailabilityServiceConfiguration(in *AvailabilityServiceConfiguration, out *config.AvailabilityServiceConfiguration, s conversion.Scope) error {
	out.Url = in.Url
	out.ApiKey = in.ApiKey
	return nil
}

// Convert_v1alpha1_AvailabilityServiceConfiguration_To_config_AvailabilityServiceConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_AvailabilityServiceConfiguration_To_config_AvailabilityServiceConfiguration(in *AvailabilityServiceConfiguration, out *config.AvailabilityServiceConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_AvailabilityServiceConfiguration_To_config_AvailabilityServiceConfiguration(in, out, s)
}

func autoConvert_config_AvailabilityServiceConfiguration_To_v1alpha1_AvailabilityServiceConfiguration(in *config.AvailabilityServiceConfiguration, out *AvailabilityServiceConfiguration, s conversion.Scope) error {
	out.Url = in.Url
	out.ApiKey = in.ApiKey
	return nil
}

// Convert_config_AvailabilityServiceConfiguration_To_v1alpha1_AvailabilityServiceConfiguration is an autogenerated conversion function.
func Convert_config_AvailabilityServiceConfiguration_To_v1alpha1_AvailabilityServiceConfiguration(in *config.AvailabilityServiceConfiguration, out *AvailabilityServiceConfiguration, s conversion.Scope) error {
	return autoConvert_config_AvailabilityServiceConfiguration_To_v1alpha1_AvailabilityServiceConfiguration(in, out, s)
}

func autoConvert_v1alpha1_CrdManagementConfiguration_To_config_CrdManagementConfiguration(in *CrdManagementConfiguration, out *config.CrdManagementConfiguration, s conversion.Scope) error {
	out.DeployCustomResourceDefinitions = (*bool)(unsafe.Pointer(in.DeployCustomResourceDefinitions))
	out.ForceUpdate = (*bool)(unsafe.Pointer(in.ForceUpdate))
	return nil
}

// Convert_v1alpha1_CrdManagementConfiguration_To_config_CrdManagementConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_CrdManagementConfiguration_To_config_CrdManagementConfiguration(in *CrdManagementConfiguration, out *config.CrdManagementConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_CrdManagementConfiguration_To_config_CrdManagementConfiguration(in, out, s)
}

func autoConvert_config_CrdManagementConfiguration_To_v1alpha1_CrdManagementConfiguration(in *config.CrdManagementConfiguration, out *CrdManagementConfiguration, s conversion.Scope) error {
	out.DeployCustomResourceDefinitions = (*bool)(unsafe.Pointer(in.DeployCustomResourceDefinitions))
	out.ForceUpdate = (*bool)(unsafe.Pointer(in.ForceUpdate))
	return nil
}

// Convert_config_CrdManagementConfiguration_To_v1alpha1_CrdManagementConfiguration is an autogenerated conversion function.
func Convert_config_CrdManagementConfiguration_To_v1alpha1_CrdManagementConfiguration(in *config.CrdManagementConfiguration, out *CrdManagementConfiguration, s conversion.Scope) error {
	return autoConvert_config_CrdManagementConfiguration_To_v1alpha1_CrdManagementConfiguration(in, out, s)
}

func autoConvert_v1alpha1_LandscaperServiceComponentConfiguration_To_config_LandscaperServiceComponentConfiguration(in *LandscaperServiceComponentConfiguration, out *config.LandscaperServiceComponentConfiguration, s conversion.Scope) error {
	out.Name = in.Name
	out.Version = in.Version
	out.RepositoryContext = in.RepositoryContext
	out.RegistryPullSecrets = *(*[]v1.SecretReference)(unsafe.Pointer(&in.RegistryPullSecrets))
	return nil
}

// Convert_v1alpha1_LandscaperServiceComponentConfiguration_To_config_LandscaperServiceComponentConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_LandscaperServiceComponentConfiguration_To_config_LandscaperServiceComponentConfiguration(in *LandscaperServiceComponentConfiguration, out *config.LandscaperServiceComponentConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_LandscaperServiceComponentConfiguration_To_config_LandscaperServiceComponentConfiguration(in, out, s)
}

func autoConvert_config_LandscaperServiceComponentConfiguration_To_v1alpha1_LandscaperServiceComponentConfiguration(in *config.LandscaperServiceComponentConfiguration, out *LandscaperServiceComponentConfiguration, s conversion.Scope) error {
	out.Name = in.Name
	out.Version = in.Version
	out.RepositoryContext = in.RepositoryContext
	out.RegistryPullSecrets = *(*[]v1.SecretReference)(unsafe.Pointer(&in.RegistryPullSecrets))
	return nil
}

// Convert_config_LandscaperServiceComponentConfiguration_To_v1alpha1_LandscaperServiceComponentConfiguration is an autogenerated conversion function.
func Convert_config_LandscaperServiceComponentConfiguration_To_v1alpha1_LandscaperServiceComponentConfiguration(in *config.LandscaperServiceComponentConfiguration, out *LandscaperServiceComponentConfiguration, s conversion.Scope) error {
	return autoConvert_config_LandscaperServiceComponentConfiguration_To_v1alpha1_LandscaperServiceComponentConfiguration(in, out, s)
}

func autoConvert_v1alpha1_LandscaperServiceConfiguration_To_config_LandscaperServiceConfiguration(in *LandscaperServiceConfiguration, out *config.LandscaperServiceConfiguration, s conversion.Scope) error {
	out.Metrics = (*config.MetricsConfiguration)(unsafe.Pointer(in.Metrics))
	if err := Convert_v1alpha1_AvailabilityMonitoringConfiguration_To_config_AvailabilityMonitoringConfiguration(&in.AvailabilityMonitoring, &out.AvailabilityMonitoring, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_CrdManagementConfiguration_To_config_CrdManagementConfiguration(&in.CrdManagement, &out.CrdManagement, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_LandscaperServiceComponentConfiguration_To_config_LandscaperServiceComponentConfiguration(&in.LandscaperServiceComponent, &out.LandscaperServiceComponent, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_LandscaperServiceConfiguration_To_config_LandscaperServiceConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_LandscaperServiceConfiguration_To_config_LandscaperServiceConfiguration(in *LandscaperServiceConfiguration, out *config.LandscaperServiceConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_LandscaperServiceConfiguration_To_config_LandscaperServiceConfiguration(in, out, s)
}

func autoConvert_config_LandscaperServiceConfiguration_To_v1alpha1_LandscaperServiceConfiguration(in *config.LandscaperServiceConfiguration, out *LandscaperServiceConfiguration, s conversion.Scope) error {
	out.Metrics = (*MetricsConfiguration)(unsafe.Pointer(in.Metrics))
	if err := Convert_config_AvailabilityMonitoringConfiguration_To_v1alpha1_AvailabilityMonitoringConfiguration(&in.AvailabilityMonitoring, &out.AvailabilityMonitoring, s); err != nil {
		return err
	}
	if err := Convert_config_CrdManagementConfiguration_To_v1alpha1_CrdManagementConfiguration(&in.CrdManagement, &out.CrdManagement, s); err != nil {
		return err
	}
	if err := Convert_config_LandscaperServiceComponentConfiguration_To_v1alpha1_LandscaperServiceComponentConfiguration(&in.LandscaperServiceComponent, &out.LandscaperServiceComponent, s); err != nil {
		return err
	}
	return nil
}

// Convert_config_LandscaperServiceConfiguration_To_v1alpha1_LandscaperServiceConfiguration is an autogenerated conversion function.
func Convert_config_LandscaperServiceConfiguration_To_v1alpha1_LandscaperServiceConfiguration(in *config.LandscaperServiceConfiguration, out *LandscaperServiceConfiguration, s conversion.Scope) error {
	return autoConvert_config_LandscaperServiceConfiguration_To_v1alpha1_LandscaperServiceConfiguration(in, out, s)
}

func autoConvert_v1alpha1_MetricsConfiguration_To_config_MetricsConfiguration(in *MetricsConfiguration, out *config.MetricsConfiguration, s conversion.Scope) error {
	out.Port = in.Port
	return nil
}

// Convert_v1alpha1_MetricsConfiguration_To_config_MetricsConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_MetricsConfiguration_To_config_MetricsConfiguration(in *MetricsConfiguration, out *config.MetricsConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_MetricsConfiguration_To_config_MetricsConfiguration(in, out, s)
}

func autoConvert_config_MetricsConfiguration_To_v1alpha1_MetricsConfiguration(in *config.MetricsConfiguration, out *MetricsConfiguration, s conversion.Scope) error {
	out.Port = in.Port
	return nil
}

// Convert_config_MetricsConfiguration_To_v1alpha1_MetricsConfiguration is an autogenerated conversion function.
func Convert_config_MetricsConfiguration_To_v1alpha1_MetricsConfiguration(in *config.MetricsConfiguration, out *MetricsConfiguration, s conversion.Scope) error {
	return autoConvert_config_MetricsConfiguration_To_v1alpha1_MetricsConfiguration(in, out, s)
}
