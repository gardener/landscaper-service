//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

SPDX-License-Identifier: Apache-2.0
*/
// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	json "encoding/json"
	unsafe "unsafe"

	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"

	config "github.com/gardener/landscaper/apis/config"
	core "github.com/gardener/landscaper/apis/core"
	corev1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	helm "github.com/gardener/landscaper/apis/deployer/helm"
	continuousreconcile "github.com/gardener/landscaper/apis/deployer/utils/continuousreconcile"
	managedresource "github.com/gardener/landscaper/apis/deployer/utils/managedresource"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*ArchiveAccess)(nil), (*helm.ArchiveAccess)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ArchiveAccess_To_helm_ArchiveAccess(a.(*ArchiveAccess), b.(*helm.ArchiveAccess), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.ArchiveAccess)(nil), (*ArchiveAccess)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_ArchiveAccess_To_v1alpha1_ArchiveAccess(a.(*helm.ArchiveAccess), b.(*ArchiveAccess), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Auth)(nil), (*helm.Auth)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Auth_To_helm_Auth(a.(*Auth), b.(*helm.Auth), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.Auth)(nil), (*Auth)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_Auth_To_v1alpha1_Auth(a.(*helm.Auth), b.(*Auth), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Chart)(nil), (*helm.Chart)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Chart_To_helm_Chart(a.(*Chart), b.(*helm.Chart), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.Chart)(nil), (*Chart)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_Chart_To_v1alpha1_Chart(a.(*helm.Chart), b.(*Chart), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Configuration)(nil), (*helm.Configuration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Configuration_To_helm_Configuration(a.(*Configuration), b.(*helm.Configuration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.Configuration)(nil), (*Configuration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_Configuration_To_v1alpha1_Configuration(a.(*helm.Configuration), b.(*Configuration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Controller)(nil), (*helm.Controller)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Controller_To_helm_Controller(a.(*Controller), b.(*helm.Controller), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.Controller)(nil), (*Controller)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_Controller_To_v1alpha1_Controller(a.(*helm.Controller), b.(*Controller), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ExportConfiguration)(nil), (*helm.ExportConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ExportConfiguration_To_helm_ExportConfiguration(a.(*ExportConfiguration), b.(*helm.ExportConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.ExportConfiguration)(nil), (*ExportConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_ExportConfiguration_To_v1alpha1_ExportConfiguration(a.(*helm.ExportConfiguration), b.(*ExportConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HPAConfiguration)(nil), (*helm.HPAConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_HPAConfiguration_To_helm_HPAConfiguration(a.(*HPAConfiguration), b.(*helm.HPAConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.HPAConfiguration)(nil), (*HPAConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_HPAConfiguration_To_v1alpha1_HPAConfiguration(a.(*helm.HPAConfiguration), b.(*HPAConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HelmChartRepo)(nil), (*helm.HelmChartRepo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_HelmChartRepo_To_helm_HelmChartRepo(a.(*HelmChartRepo), b.(*helm.HelmChartRepo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.HelmChartRepo)(nil), (*HelmChartRepo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_HelmChartRepo_To_v1alpha1_HelmChartRepo(a.(*helm.HelmChartRepo), b.(*HelmChartRepo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HelmChartRepoCredentials)(nil), (*helm.HelmChartRepoCredentials)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_HelmChartRepoCredentials_To_helm_HelmChartRepoCredentials(a.(*HelmChartRepoCredentials), b.(*helm.HelmChartRepoCredentials), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.HelmChartRepoCredentials)(nil), (*HelmChartRepoCredentials)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_HelmChartRepoCredentials_To_v1alpha1_HelmChartRepoCredentials(a.(*helm.HelmChartRepoCredentials), b.(*HelmChartRepoCredentials), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HelmDeploymentConfiguration)(nil), (*helm.HelmDeploymentConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_HelmDeploymentConfiguration_To_helm_HelmDeploymentConfiguration(a.(*HelmDeploymentConfiguration), b.(*helm.HelmDeploymentConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.HelmDeploymentConfiguration)(nil), (*HelmDeploymentConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_HelmDeploymentConfiguration_To_v1alpha1_HelmDeploymentConfiguration(a.(*helm.HelmDeploymentConfiguration), b.(*HelmDeploymentConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HelmInstallConfiguration)(nil), (*helm.HelmInstallConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_HelmInstallConfiguration_To_helm_HelmInstallConfiguration(a.(*HelmInstallConfiguration), b.(*helm.HelmInstallConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.HelmInstallConfiguration)(nil), (*HelmInstallConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_HelmInstallConfiguration_To_v1alpha1_HelmInstallConfiguration(a.(*helm.HelmInstallConfiguration), b.(*HelmInstallConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HelmUninstallConfiguration)(nil), (*helm.HelmUninstallConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_HelmUninstallConfiguration_To_helm_HelmUninstallConfiguration(a.(*HelmUninstallConfiguration), b.(*helm.HelmUninstallConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.HelmUninstallConfiguration)(nil), (*HelmUninstallConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_HelmUninstallConfiguration_To_v1alpha1_HelmUninstallConfiguration(a.(*helm.HelmUninstallConfiguration), b.(*HelmUninstallConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ProviderConfiguration)(nil), (*helm.ProviderConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ProviderConfiguration_To_helm_ProviderConfiguration(a.(*ProviderConfiguration), b.(*helm.ProviderConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.ProviderConfiguration)(nil), (*ProviderConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_ProviderConfiguration_To_v1alpha1_ProviderConfiguration(a.(*helm.ProviderConfiguration), b.(*ProviderConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ProviderStatus)(nil), (*helm.ProviderStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ProviderStatus_To_helm_ProviderStatus(a.(*ProviderStatus), b.(*helm.ProviderStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.ProviderStatus)(nil), (*ProviderStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_ProviderStatus_To_v1alpha1_ProviderStatus(a.(*helm.ProviderStatus), b.(*ProviderStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*RemoteArchiveAccess)(nil), (*helm.RemoteArchiveAccess)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RemoteArchiveAccess_To_helm_RemoteArchiveAccess(a.(*RemoteArchiveAccess), b.(*helm.RemoteArchiveAccess), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.RemoteArchiveAccess)(nil), (*RemoteArchiveAccess)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_RemoteArchiveAccess_To_v1alpha1_RemoteArchiveAccess(a.(*helm.RemoteArchiveAccess), b.(*RemoteArchiveAccess), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*RemoteChartReference)(nil), (*helm.RemoteChartReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RemoteChartReference_To_helm_RemoteChartReference(a.(*RemoteChartReference), b.(*helm.RemoteChartReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*helm.RemoteChartReference)(nil), (*RemoteChartReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_helm_RemoteChartReference_To_v1alpha1_RemoteChartReference(a.(*helm.RemoteChartReference), b.(*RemoteChartReference), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_ArchiveAccess_To_helm_ArchiveAccess(in *ArchiveAccess, out *helm.ArchiveAccess, s conversion.Scope) error {
	out.Raw = in.Raw
	out.Remote = (*helm.RemoteArchiveAccess)(unsafe.Pointer(in.Remote))
	return nil
}

// Convert_v1alpha1_ArchiveAccess_To_helm_ArchiveAccess is an autogenerated conversion function.
func Convert_v1alpha1_ArchiveAccess_To_helm_ArchiveAccess(in *ArchiveAccess, out *helm.ArchiveAccess, s conversion.Scope) error {
	return autoConvert_v1alpha1_ArchiveAccess_To_helm_ArchiveAccess(in, out, s)
}

func autoConvert_helm_ArchiveAccess_To_v1alpha1_ArchiveAccess(in *helm.ArchiveAccess, out *ArchiveAccess, s conversion.Scope) error {
	out.Raw = in.Raw
	out.Remote = (*RemoteArchiveAccess)(unsafe.Pointer(in.Remote))
	return nil
}

// Convert_helm_ArchiveAccess_To_v1alpha1_ArchiveAccess is an autogenerated conversion function.
func Convert_helm_ArchiveAccess_To_v1alpha1_ArchiveAccess(in *helm.ArchiveAccess, out *ArchiveAccess, s conversion.Scope) error {
	return autoConvert_helm_ArchiveAccess_To_v1alpha1_ArchiveAccess(in, out, s)
}

func autoConvert_v1alpha1_Auth_To_helm_Auth(in *Auth, out *helm.Auth, s conversion.Scope) error {
	out.URL = in.URL
	out.CustomCAData = in.CustomCAData
	out.AuthHeader = in.AuthHeader
	out.SecretRef = (*corev1alpha1.LocalSecretReference)(unsafe.Pointer(in.SecretRef))
	return nil
}

// Convert_v1alpha1_Auth_To_helm_Auth is an autogenerated conversion function.
func Convert_v1alpha1_Auth_To_helm_Auth(in *Auth, out *helm.Auth, s conversion.Scope) error {
	return autoConvert_v1alpha1_Auth_To_helm_Auth(in, out, s)
}

func autoConvert_helm_Auth_To_v1alpha1_Auth(in *helm.Auth, out *Auth, s conversion.Scope) error {
	out.URL = in.URL
	out.CustomCAData = in.CustomCAData
	out.AuthHeader = in.AuthHeader
	out.SecretRef = (*corev1alpha1.LocalSecretReference)(unsafe.Pointer(in.SecretRef))
	return nil
}

// Convert_helm_Auth_To_v1alpha1_Auth is an autogenerated conversion function.
func Convert_helm_Auth_To_v1alpha1_Auth(in *helm.Auth, out *Auth, s conversion.Scope) error {
	return autoConvert_helm_Auth_To_v1alpha1_Auth(in, out, s)
}

func autoConvert_v1alpha1_Chart_To_helm_Chart(in *Chart, out *helm.Chart, s conversion.Scope) error {
	out.Ref = in.Ref
	out.FromResource = (*helm.RemoteChartReference)(unsafe.Pointer(in.FromResource))
	out.Archive = (*helm.ArchiveAccess)(unsafe.Pointer(in.Archive))
	out.HelmChartRepo = (*helm.HelmChartRepo)(unsafe.Pointer(in.HelmChartRepo))
	out.ResourceRef = in.ResourceRef
	return nil
}

// Convert_v1alpha1_Chart_To_helm_Chart is an autogenerated conversion function.
func Convert_v1alpha1_Chart_To_helm_Chart(in *Chart, out *helm.Chart, s conversion.Scope) error {
	return autoConvert_v1alpha1_Chart_To_helm_Chart(in, out, s)
}

func autoConvert_helm_Chart_To_v1alpha1_Chart(in *helm.Chart, out *Chart, s conversion.Scope) error {
	out.Ref = in.Ref
	out.FromResource = (*RemoteChartReference)(unsafe.Pointer(in.FromResource))
	out.Archive = (*ArchiveAccess)(unsafe.Pointer(in.Archive))
	out.HelmChartRepo = (*HelmChartRepo)(unsafe.Pointer(in.HelmChartRepo))
	out.ResourceRef = in.ResourceRef
	return nil
}

// Convert_helm_Chart_To_v1alpha1_Chart is an autogenerated conversion function.
func Convert_helm_Chart_To_v1alpha1_Chart(in *helm.Chart, out *Chart, s conversion.Scope) error {
	return autoConvert_helm_Chart_To_v1alpha1_Chart(in, out, s)
}

func autoConvert_v1alpha1_Configuration_To_helm_Configuration(in *Configuration, out *helm.Configuration, s conversion.Scope) error {
	out.Identity = in.Identity
	out.OCI = (*config.OCIConfiguration)(unsafe.Pointer(in.OCI))
	out.TargetSelector = *(*[]corev1alpha1.TargetSelector)(unsafe.Pointer(&in.TargetSelector))
	if err := Convert_v1alpha1_ExportConfiguration_To_helm_ExportConfiguration(&in.Export, &out.Export, s); err != nil {
		return err
	}
	out.HPAConfiguration = (*helm.HPAConfiguration)(unsafe.Pointer(in.HPAConfiguration))
	if err := Convert_v1alpha1_Controller_To_helm_Controller(&in.Controller, &out.Controller, s); err != nil {
		return err
	}
	out.UseOCMLib = in.UseOCMLib
	return nil
}

// Convert_v1alpha1_Configuration_To_helm_Configuration is an autogenerated conversion function.
func Convert_v1alpha1_Configuration_To_helm_Configuration(in *Configuration, out *helm.Configuration, s conversion.Scope) error {
	return autoConvert_v1alpha1_Configuration_To_helm_Configuration(in, out, s)
}

func autoConvert_helm_Configuration_To_v1alpha1_Configuration(in *helm.Configuration, out *Configuration, s conversion.Scope) error {
	out.Identity = in.Identity
	out.OCI = (*config.OCIConfiguration)(unsafe.Pointer(in.OCI))
	out.TargetSelector = *(*[]corev1alpha1.TargetSelector)(unsafe.Pointer(&in.TargetSelector))
	if err := Convert_helm_ExportConfiguration_To_v1alpha1_ExportConfiguration(&in.Export, &out.Export, s); err != nil {
		return err
	}
	out.HPAConfiguration = (*HPAConfiguration)(unsafe.Pointer(in.HPAConfiguration))
	if err := Convert_helm_Controller_To_v1alpha1_Controller(&in.Controller, &out.Controller, s); err != nil {
		return err
	}
	out.UseOCMLib = in.UseOCMLib
	return nil
}

// Convert_helm_Configuration_To_v1alpha1_Configuration is an autogenerated conversion function.
func Convert_helm_Configuration_To_v1alpha1_Configuration(in *helm.Configuration, out *Configuration, s conversion.Scope) error {
	return autoConvert_helm_Configuration_To_v1alpha1_Configuration(in, out, s)
}

func autoConvert_v1alpha1_Controller_To_helm_Controller(in *Controller, out *helm.Controller, s conversion.Scope) error {
	out.CommonControllerConfig = in.CommonControllerConfig
	return nil
}

// Convert_v1alpha1_Controller_To_helm_Controller is an autogenerated conversion function.
func Convert_v1alpha1_Controller_To_helm_Controller(in *Controller, out *helm.Controller, s conversion.Scope) error {
	return autoConvert_v1alpha1_Controller_To_helm_Controller(in, out, s)
}

func autoConvert_helm_Controller_To_v1alpha1_Controller(in *helm.Controller, out *Controller, s conversion.Scope) error {
	out.CommonControllerConfig = in.CommonControllerConfig
	return nil
}

// Convert_helm_Controller_To_v1alpha1_Controller is an autogenerated conversion function.
func Convert_helm_Controller_To_v1alpha1_Controller(in *helm.Controller, out *Controller, s conversion.Scope) error {
	return autoConvert_helm_Controller_To_v1alpha1_Controller(in, out, s)
}

func autoConvert_v1alpha1_ExportConfiguration_To_helm_ExportConfiguration(in *ExportConfiguration, out *helm.ExportConfiguration, s conversion.Scope) error {
	out.DefaultTimeout = (*corev1alpha1.Duration)(unsafe.Pointer(in.DefaultTimeout))
	return nil
}

// Convert_v1alpha1_ExportConfiguration_To_helm_ExportConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_ExportConfiguration_To_helm_ExportConfiguration(in *ExportConfiguration, out *helm.ExportConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_ExportConfiguration_To_helm_ExportConfiguration(in, out, s)
}

func autoConvert_helm_ExportConfiguration_To_v1alpha1_ExportConfiguration(in *helm.ExportConfiguration, out *ExportConfiguration, s conversion.Scope) error {
	out.DefaultTimeout = (*corev1alpha1.Duration)(unsafe.Pointer(in.DefaultTimeout))
	return nil
}

// Convert_helm_ExportConfiguration_To_v1alpha1_ExportConfiguration is an autogenerated conversion function.
func Convert_helm_ExportConfiguration_To_v1alpha1_ExportConfiguration(in *helm.ExportConfiguration, out *ExportConfiguration, s conversion.Scope) error {
	return autoConvert_helm_ExportConfiguration_To_v1alpha1_ExportConfiguration(in, out, s)
}

func autoConvert_v1alpha1_HPAConfiguration_To_helm_HPAConfiguration(in *HPAConfiguration, out *helm.HPAConfiguration, s conversion.Scope) error {
	out.MaxReplicas = in.MaxReplicas
	return nil
}

// Convert_v1alpha1_HPAConfiguration_To_helm_HPAConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_HPAConfiguration_To_helm_HPAConfiguration(in *HPAConfiguration, out *helm.HPAConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_HPAConfiguration_To_helm_HPAConfiguration(in, out, s)
}

func autoConvert_helm_HPAConfiguration_To_v1alpha1_HPAConfiguration(in *helm.HPAConfiguration, out *HPAConfiguration, s conversion.Scope) error {
	out.MaxReplicas = in.MaxReplicas
	return nil
}

// Convert_helm_HPAConfiguration_To_v1alpha1_HPAConfiguration is an autogenerated conversion function.
func Convert_helm_HPAConfiguration_To_v1alpha1_HPAConfiguration(in *helm.HPAConfiguration, out *HPAConfiguration, s conversion.Scope) error {
	return autoConvert_helm_HPAConfiguration_To_v1alpha1_HPAConfiguration(in, out, s)
}

func autoConvert_v1alpha1_HelmChartRepo_To_helm_HelmChartRepo(in *HelmChartRepo, out *helm.HelmChartRepo, s conversion.Scope) error {
	out.HelmChartRepoUrl = in.HelmChartRepoUrl
	out.HelmChartName = in.HelmChartName
	out.HelmChartVersion = in.HelmChartVersion
	return nil
}

// Convert_v1alpha1_HelmChartRepo_To_helm_HelmChartRepo is an autogenerated conversion function.
func Convert_v1alpha1_HelmChartRepo_To_helm_HelmChartRepo(in *HelmChartRepo, out *helm.HelmChartRepo, s conversion.Scope) error {
	return autoConvert_v1alpha1_HelmChartRepo_To_helm_HelmChartRepo(in, out, s)
}

func autoConvert_helm_HelmChartRepo_To_v1alpha1_HelmChartRepo(in *helm.HelmChartRepo, out *HelmChartRepo, s conversion.Scope) error {
	out.HelmChartRepoUrl = in.HelmChartRepoUrl
	out.HelmChartName = in.HelmChartName
	out.HelmChartVersion = in.HelmChartVersion
	return nil
}

// Convert_helm_HelmChartRepo_To_v1alpha1_HelmChartRepo is an autogenerated conversion function.
func Convert_helm_HelmChartRepo_To_v1alpha1_HelmChartRepo(in *helm.HelmChartRepo, out *HelmChartRepo, s conversion.Scope) error {
	return autoConvert_helm_HelmChartRepo_To_v1alpha1_HelmChartRepo(in, out, s)
}

func autoConvert_v1alpha1_HelmChartRepoCredentials_To_helm_HelmChartRepoCredentials(in *HelmChartRepoCredentials, out *helm.HelmChartRepoCredentials, s conversion.Scope) error {
	out.Auths = *(*[]helm.Auth)(unsafe.Pointer(&in.Auths))
	return nil
}

// Convert_v1alpha1_HelmChartRepoCredentials_To_helm_HelmChartRepoCredentials is an autogenerated conversion function.
func Convert_v1alpha1_HelmChartRepoCredentials_To_helm_HelmChartRepoCredentials(in *HelmChartRepoCredentials, out *helm.HelmChartRepoCredentials, s conversion.Scope) error {
	return autoConvert_v1alpha1_HelmChartRepoCredentials_To_helm_HelmChartRepoCredentials(in, out, s)
}

func autoConvert_helm_HelmChartRepoCredentials_To_v1alpha1_HelmChartRepoCredentials(in *helm.HelmChartRepoCredentials, out *HelmChartRepoCredentials, s conversion.Scope) error {
	out.Auths = *(*[]Auth)(unsafe.Pointer(&in.Auths))
	return nil
}

// Convert_helm_HelmChartRepoCredentials_To_v1alpha1_HelmChartRepoCredentials is an autogenerated conversion function.
func Convert_helm_HelmChartRepoCredentials_To_v1alpha1_HelmChartRepoCredentials(in *helm.HelmChartRepoCredentials, out *HelmChartRepoCredentials, s conversion.Scope) error {
	return autoConvert_helm_HelmChartRepoCredentials_To_v1alpha1_HelmChartRepoCredentials(in, out, s)
}

func autoConvert_v1alpha1_HelmDeploymentConfiguration_To_helm_HelmDeploymentConfiguration(in *HelmDeploymentConfiguration, out *helm.HelmDeploymentConfiguration, s conversion.Scope) error {
	out.Install = *(*map[string]core.AnyJSON)(unsafe.Pointer(&in.Install))
	out.Upgrade = *(*map[string]core.AnyJSON)(unsafe.Pointer(&in.Upgrade))
	out.Uninstall = *(*map[string]core.AnyJSON)(unsafe.Pointer(&in.Uninstall))
	return nil
}

// Convert_v1alpha1_HelmDeploymentConfiguration_To_helm_HelmDeploymentConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_HelmDeploymentConfiguration_To_helm_HelmDeploymentConfiguration(in *HelmDeploymentConfiguration, out *helm.HelmDeploymentConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_HelmDeploymentConfiguration_To_helm_HelmDeploymentConfiguration(in, out, s)
}

func autoConvert_helm_HelmDeploymentConfiguration_To_v1alpha1_HelmDeploymentConfiguration(in *helm.HelmDeploymentConfiguration, out *HelmDeploymentConfiguration, s conversion.Scope) error {
	out.Install = *(*map[string]corev1alpha1.AnyJSON)(unsafe.Pointer(&in.Install))
	out.Upgrade = *(*map[string]corev1alpha1.AnyJSON)(unsafe.Pointer(&in.Upgrade))
	out.Uninstall = *(*map[string]corev1alpha1.AnyJSON)(unsafe.Pointer(&in.Uninstall))
	return nil
}

// Convert_helm_HelmDeploymentConfiguration_To_v1alpha1_HelmDeploymentConfiguration is an autogenerated conversion function.
func Convert_helm_HelmDeploymentConfiguration_To_v1alpha1_HelmDeploymentConfiguration(in *helm.HelmDeploymentConfiguration, out *HelmDeploymentConfiguration, s conversion.Scope) error {
	return autoConvert_helm_HelmDeploymentConfiguration_To_v1alpha1_HelmDeploymentConfiguration(in, out, s)
}

func autoConvert_v1alpha1_HelmInstallConfiguration_To_helm_HelmInstallConfiguration(in *HelmInstallConfiguration, out *helm.HelmInstallConfiguration, s conversion.Scope) error {
	out.Atomic = in.Atomic
	out.Timeout = (*corev1alpha1.Duration)(unsafe.Pointer(in.Timeout))
	return nil
}

// Convert_v1alpha1_HelmInstallConfiguration_To_helm_HelmInstallConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_HelmInstallConfiguration_To_helm_HelmInstallConfiguration(in *HelmInstallConfiguration, out *helm.HelmInstallConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_HelmInstallConfiguration_To_helm_HelmInstallConfiguration(in, out, s)
}

func autoConvert_helm_HelmInstallConfiguration_To_v1alpha1_HelmInstallConfiguration(in *helm.HelmInstallConfiguration, out *HelmInstallConfiguration, s conversion.Scope) error {
	out.Atomic = in.Atomic
	out.Timeout = (*corev1alpha1.Duration)(unsafe.Pointer(in.Timeout))
	return nil
}

// Convert_helm_HelmInstallConfiguration_To_v1alpha1_HelmInstallConfiguration is an autogenerated conversion function.
func Convert_helm_HelmInstallConfiguration_To_v1alpha1_HelmInstallConfiguration(in *helm.HelmInstallConfiguration, out *HelmInstallConfiguration, s conversion.Scope) error {
	return autoConvert_helm_HelmInstallConfiguration_To_v1alpha1_HelmInstallConfiguration(in, out, s)
}

func autoConvert_v1alpha1_HelmUninstallConfiguration_To_helm_HelmUninstallConfiguration(in *HelmUninstallConfiguration, out *helm.HelmUninstallConfiguration, s conversion.Scope) error {
	out.Timeout = (*corev1alpha1.Duration)(unsafe.Pointer(in.Timeout))
	return nil
}

// Convert_v1alpha1_HelmUninstallConfiguration_To_helm_HelmUninstallConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_HelmUninstallConfiguration_To_helm_HelmUninstallConfiguration(in *HelmUninstallConfiguration, out *helm.HelmUninstallConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_HelmUninstallConfiguration_To_helm_HelmUninstallConfiguration(in, out, s)
}

func autoConvert_helm_HelmUninstallConfiguration_To_v1alpha1_HelmUninstallConfiguration(in *helm.HelmUninstallConfiguration, out *HelmUninstallConfiguration, s conversion.Scope) error {
	out.Timeout = (*corev1alpha1.Duration)(unsafe.Pointer(in.Timeout))
	return nil
}

// Convert_helm_HelmUninstallConfiguration_To_v1alpha1_HelmUninstallConfiguration is an autogenerated conversion function.
func Convert_helm_HelmUninstallConfiguration_To_v1alpha1_HelmUninstallConfiguration(in *helm.HelmUninstallConfiguration, out *HelmUninstallConfiguration, s conversion.Scope) error {
	return autoConvert_helm_HelmUninstallConfiguration_To_v1alpha1_HelmUninstallConfiguration(in, out, s)
}

func autoConvert_v1alpha1_ProviderConfiguration_To_helm_ProviderConfiguration(in *ProviderConfiguration, out *helm.ProviderConfiguration, s conversion.Scope) error {
	out.Kubeconfig = in.Kubeconfig
	out.UpdateStrategy = helm.UpdateStrategy(in.UpdateStrategy)
	out.ReadinessChecks = in.ReadinessChecks
	if err := Convert_v1alpha1_Chart_To_helm_Chart(&in.Chart, &out.Chart, s); err != nil {
		return err
	}
	out.Name = in.Name
	out.Namespace = in.Namespace
	out.CreateNamespace = in.CreateNamespace
	out.Values = *(*json.RawMessage)(unsafe.Pointer(&in.Values))
	out.ExportsFromManifests = *(*[]managedresource.Export)(unsafe.Pointer(&in.ExportsFromManifests))
	out.Exports = (*managedresource.Exports)(unsafe.Pointer(in.Exports))
	out.ContinuousReconcile = (*continuousreconcile.ContinuousReconcileSpec)(unsafe.Pointer(in.ContinuousReconcile))
	out.HelmDeployment = (*bool)(unsafe.Pointer(in.HelmDeployment))
	out.HelmDeploymentConfig = (*helm.HelmDeploymentConfiguration)(unsafe.Pointer(in.HelmDeploymentConfig))
	out.DeletionGroups = *(*[]managedresource.DeletionGroupDefinition)(unsafe.Pointer(&in.DeletionGroups))
	out.DeletionGroupsDuringUpdate = *(*[]managedresource.DeletionGroupDefinition)(unsafe.Pointer(&in.DeletionGroupsDuringUpdate))
	return nil
}

// Convert_v1alpha1_ProviderConfiguration_To_helm_ProviderConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_ProviderConfiguration_To_helm_ProviderConfiguration(in *ProviderConfiguration, out *helm.ProviderConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_ProviderConfiguration_To_helm_ProviderConfiguration(in, out, s)
}

func autoConvert_helm_ProviderConfiguration_To_v1alpha1_ProviderConfiguration(in *helm.ProviderConfiguration, out *ProviderConfiguration, s conversion.Scope) error {
	out.Kubeconfig = in.Kubeconfig
	out.ReadinessChecks = in.ReadinessChecks
	out.UpdateStrategy = UpdateStrategy(in.UpdateStrategy)
	if err := Convert_helm_Chart_To_v1alpha1_Chart(&in.Chart, &out.Chart, s); err != nil {
		return err
	}
	out.Name = in.Name
	out.Namespace = in.Namespace
	out.CreateNamespace = in.CreateNamespace
	out.Values = *(*json.RawMessage)(unsafe.Pointer(&in.Values))
	out.ExportsFromManifests = *(*[]managedresource.Export)(unsafe.Pointer(&in.ExportsFromManifests))
	out.Exports = (*managedresource.Exports)(unsafe.Pointer(in.Exports))
	out.ContinuousReconcile = (*continuousreconcile.ContinuousReconcileSpec)(unsafe.Pointer(in.ContinuousReconcile))
	out.HelmDeployment = (*bool)(unsafe.Pointer(in.HelmDeployment))
	out.HelmDeploymentConfig = (*HelmDeploymentConfiguration)(unsafe.Pointer(in.HelmDeploymentConfig))
	out.DeletionGroups = *(*[]managedresource.DeletionGroupDefinition)(unsafe.Pointer(&in.DeletionGroups))
	out.DeletionGroupsDuringUpdate = *(*[]managedresource.DeletionGroupDefinition)(unsafe.Pointer(&in.DeletionGroupsDuringUpdate))
	return nil
}

// Convert_helm_ProviderConfiguration_To_v1alpha1_ProviderConfiguration is an autogenerated conversion function.
func Convert_helm_ProviderConfiguration_To_v1alpha1_ProviderConfiguration(in *helm.ProviderConfiguration, out *ProviderConfiguration, s conversion.Scope) error {
	return autoConvert_helm_ProviderConfiguration_To_v1alpha1_ProviderConfiguration(in, out, s)
}

func autoConvert_v1alpha1_ProviderStatus_To_helm_ProviderStatus(in *ProviderStatus, out *helm.ProviderStatus, s conversion.Scope) error {
	out.ManagedResources = *(*managedresource.ManagedResourceStatusList)(unsafe.Pointer(&in.ManagedResources))
	return nil
}

// Convert_v1alpha1_ProviderStatus_To_helm_ProviderStatus is an autogenerated conversion function.
func Convert_v1alpha1_ProviderStatus_To_helm_ProviderStatus(in *ProviderStatus, out *helm.ProviderStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_ProviderStatus_To_helm_ProviderStatus(in, out, s)
}

func autoConvert_helm_ProviderStatus_To_v1alpha1_ProviderStatus(in *helm.ProviderStatus, out *ProviderStatus, s conversion.Scope) error {
	out.ManagedResources = *(*managedresource.ManagedResourceStatusList)(unsafe.Pointer(&in.ManagedResources))
	return nil
}

// Convert_helm_ProviderStatus_To_v1alpha1_ProviderStatus is an autogenerated conversion function.
func Convert_helm_ProviderStatus_To_v1alpha1_ProviderStatus(in *helm.ProviderStatus, out *ProviderStatus, s conversion.Scope) error {
	return autoConvert_helm_ProviderStatus_To_v1alpha1_ProviderStatus(in, out, s)
}

func autoConvert_v1alpha1_RemoteArchiveAccess_To_helm_RemoteArchiveAccess(in *RemoteArchiveAccess, out *helm.RemoteArchiveAccess, s conversion.Scope) error {
	out.URL = in.URL
	return nil
}

// Convert_v1alpha1_RemoteArchiveAccess_To_helm_RemoteArchiveAccess is an autogenerated conversion function.
func Convert_v1alpha1_RemoteArchiveAccess_To_helm_RemoteArchiveAccess(in *RemoteArchiveAccess, out *helm.RemoteArchiveAccess, s conversion.Scope) error {
	return autoConvert_v1alpha1_RemoteArchiveAccess_To_helm_RemoteArchiveAccess(in, out, s)
}

func autoConvert_helm_RemoteArchiveAccess_To_v1alpha1_RemoteArchiveAccess(in *helm.RemoteArchiveAccess, out *RemoteArchiveAccess, s conversion.Scope) error {
	out.URL = in.URL
	return nil
}

// Convert_helm_RemoteArchiveAccess_To_v1alpha1_RemoteArchiveAccess is an autogenerated conversion function.
func Convert_helm_RemoteArchiveAccess_To_v1alpha1_RemoteArchiveAccess(in *helm.RemoteArchiveAccess, out *RemoteArchiveAccess, s conversion.Scope) error {
	return autoConvert_helm_RemoteArchiveAccess_To_v1alpha1_RemoteArchiveAccess(in, out, s)
}

func autoConvert_v1alpha1_RemoteChartReference_To_helm_RemoteChartReference(in *RemoteChartReference, out *helm.RemoteChartReference, s conversion.Scope) error {
	out.ComponentDescriptorDefinition = in.ComponentDescriptorDefinition
	out.ResourceName = in.ResourceName
	return nil
}

// Convert_v1alpha1_RemoteChartReference_To_helm_RemoteChartReference is an autogenerated conversion function.
func Convert_v1alpha1_RemoteChartReference_To_helm_RemoteChartReference(in *RemoteChartReference, out *helm.RemoteChartReference, s conversion.Scope) error {
	return autoConvert_v1alpha1_RemoteChartReference_To_helm_RemoteChartReference(in, out, s)
}

func autoConvert_helm_RemoteChartReference_To_v1alpha1_RemoteChartReference(in *helm.RemoteChartReference, out *RemoteChartReference, s conversion.Scope) error {
	out.ComponentDescriptorDefinition = in.ComponentDescriptorDefinition
	out.ResourceName = in.ResourceName
	return nil
}

// Convert_helm_RemoteChartReference_To_v1alpha1_RemoteChartReference is an autogenerated conversion function.
func Convert_helm_RemoteChartReference_To_v1alpha1_RemoteChartReference(in *helm.RemoteChartReference, out *RemoteChartReference, s conversion.Scope) error {
	return autoConvert_helm_RemoteChartReference_To_v1alpha1_RemoteChartReference(in, out, s)
}
