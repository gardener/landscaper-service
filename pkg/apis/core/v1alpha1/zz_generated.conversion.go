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

	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"

	core "github.com/gardener/landscaper-service/pkg/apis/core"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Error)(nil), (*core.Error)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Error_To_core_Error(a.(*Error), b.(*core.Error), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Error)(nil), (*Error)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Error_To_v1alpha1_Error(a.(*core.Error), b.(*Error), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Instance)(nil), (*core.Instance)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Instance_To_core_Instance(a.(*Instance), b.(*core.Instance), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Instance)(nil), (*Instance)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Instance_To_v1alpha1_Instance(a.(*core.Instance), b.(*Instance), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InstanceList)(nil), (*core.InstanceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InstanceList_To_core_InstanceList(a.(*InstanceList), b.(*core.InstanceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.InstanceList)(nil), (*InstanceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_InstanceList_To_v1alpha1_InstanceList(a.(*core.InstanceList), b.(*InstanceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InstanceSpec)(nil), (*core.InstanceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InstanceSpec_To_core_InstanceSpec(a.(*InstanceSpec), b.(*core.InstanceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.InstanceSpec)(nil), (*InstanceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_InstanceSpec_To_v1alpha1_InstanceSpec(a.(*core.InstanceSpec), b.(*InstanceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InstanceStatus)(nil), (*core.InstanceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InstanceStatus_To_core_InstanceStatus(a.(*InstanceStatus), b.(*core.InstanceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.InstanceStatus)(nil), (*InstanceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_InstanceStatus_To_v1alpha1_InstanceStatus(a.(*core.InstanceStatus), b.(*InstanceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LandscaperConfiguration)(nil), (*core.LandscaperConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LandscaperConfiguration_To_core_LandscaperConfiguration(a.(*LandscaperConfiguration), b.(*core.LandscaperConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LandscaperConfiguration)(nil), (*LandscaperConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LandscaperConfiguration_To_v1alpha1_LandscaperConfiguration(a.(*core.LandscaperConfiguration), b.(*LandscaperConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LandscaperDeployment)(nil), (*core.LandscaperDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LandscaperDeployment_To_core_LandscaperDeployment(a.(*LandscaperDeployment), b.(*core.LandscaperDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LandscaperDeployment)(nil), (*LandscaperDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LandscaperDeployment_To_v1alpha1_LandscaperDeployment(a.(*core.LandscaperDeployment), b.(*LandscaperDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LandscaperDeploymentList)(nil), (*core.LandscaperDeploymentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LandscaperDeploymentList_To_core_LandscaperDeploymentList(a.(*LandscaperDeploymentList), b.(*core.LandscaperDeploymentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LandscaperDeploymentList)(nil), (*LandscaperDeploymentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LandscaperDeploymentList_To_v1alpha1_LandscaperDeploymentList(a.(*core.LandscaperDeploymentList), b.(*LandscaperDeploymentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LandscaperDeploymentSpec)(nil), (*core.LandscaperDeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LandscaperDeploymentSpec_To_core_LandscaperDeploymentSpec(a.(*LandscaperDeploymentSpec), b.(*core.LandscaperDeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LandscaperDeploymentSpec)(nil), (*LandscaperDeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LandscaperDeploymentSpec_To_v1alpha1_LandscaperDeploymentSpec(a.(*core.LandscaperDeploymentSpec), b.(*LandscaperDeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LandscaperDeploymentStatus)(nil), (*core.LandscaperDeploymentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LandscaperDeploymentStatus_To_core_LandscaperDeploymentStatus(a.(*LandscaperDeploymentStatus), b.(*core.LandscaperDeploymentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LandscaperDeploymentStatus)(nil), (*LandscaperDeploymentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LandscaperDeploymentStatus_To_v1alpha1_LandscaperDeploymentStatus(a.(*core.LandscaperDeploymentStatus), b.(*LandscaperDeploymentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LandscaperServiceComponentReference)(nil), (*core.LandscaperServiceComponentReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LandscaperServiceComponentReference_To_core_LandscaperServiceComponentReference(a.(*LandscaperServiceComponentReference), b.(*core.LandscaperServiceComponentReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LandscaperServiceComponentReference)(nil), (*LandscaperServiceComponentReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LandscaperServiceComponentReference_To_v1alpha1_LandscaperServiceComponentReference(a.(*core.LandscaperServiceComponentReference), b.(*LandscaperServiceComponentReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ObjectReference)(nil), (*core.ObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ObjectReference_To_core_ObjectReference(a.(*ObjectReference), b.(*core.ObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ObjectReference)(nil), (*ObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ObjectReference_To_v1alpha1_ObjectReference(a.(*core.ObjectReference), b.(*ObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*SecretReference)(nil), (*core.SecretReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_SecretReference_To_core_SecretReference(a.(*SecretReference), b.(*core.SecretReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SecretReference)(nil), (*SecretReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SecretReference_To_v1alpha1_SecretReference(a.(*core.SecretReference), b.(*SecretReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ServiceTargetConfig)(nil), (*core.ServiceTargetConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ServiceTargetConfig_To_core_ServiceTargetConfig(a.(*ServiceTargetConfig), b.(*core.ServiceTargetConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceTargetConfig)(nil), (*ServiceTargetConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceTargetConfig_To_v1alpha1_ServiceTargetConfig(a.(*core.ServiceTargetConfig), b.(*ServiceTargetConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ServiceTargetConfigList)(nil), (*core.ServiceTargetConfigList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ServiceTargetConfigList_To_core_ServiceTargetConfigList(a.(*ServiceTargetConfigList), b.(*core.ServiceTargetConfigList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceTargetConfigList)(nil), (*ServiceTargetConfigList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceTargetConfigList_To_v1alpha1_ServiceTargetConfigList(a.(*core.ServiceTargetConfigList), b.(*ServiceTargetConfigList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ServiceTargetConfigSpec)(nil), (*core.ServiceTargetConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ServiceTargetConfigSpec_To_core_ServiceTargetConfigSpec(a.(*ServiceTargetConfigSpec), b.(*core.ServiceTargetConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceTargetConfigSpec)(nil), (*ServiceTargetConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceTargetConfigSpec_To_v1alpha1_ServiceTargetConfigSpec(a.(*core.ServiceTargetConfigSpec), b.(*ServiceTargetConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ServiceTargetConfigStatus)(nil), (*core.ServiceTargetConfigStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ServiceTargetConfigStatus_To_core_ServiceTargetConfigStatus(a.(*ServiceTargetConfigStatus), b.(*core.ServiceTargetConfigStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceTargetConfigStatus)(nil), (*ServiceTargetConfigStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceTargetConfigStatus_To_v1alpha1_ServiceTargetConfigStatus(a.(*core.ServiceTargetConfigStatus), b.(*ServiceTargetConfigStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_Error_To_core_Error(in *Error, out *core.Error, s conversion.Scope) error {
	out.Operation = in.Operation
	out.LastTransitionTime = in.LastTransitionTime
	out.LastUpdateTime = in.LastUpdateTime
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

// Convert_v1alpha1_Error_To_core_Error is an autogenerated conversion function.
func Convert_v1alpha1_Error_To_core_Error(in *Error, out *core.Error, s conversion.Scope) error {
	return autoConvert_v1alpha1_Error_To_core_Error(in, out, s)
}

func autoConvert_core_Error_To_v1alpha1_Error(in *core.Error, out *Error, s conversion.Scope) error {
	out.Operation = in.Operation
	out.LastTransitionTime = in.LastTransitionTime
	out.LastUpdateTime = in.LastUpdateTime
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

// Convert_core_Error_To_v1alpha1_Error is an autogenerated conversion function.
func Convert_core_Error_To_v1alpha1_Error(in *core.Error, out *Error, s conversion.Scope) error {
	return autoConvert_core_Error_To_v1alpha1_Error(in, out, s)
}

func autoConvert_v1alpha1_Instance_To_core_Instance(in *Instance, out *core.Instance, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_InstanceSpec_To_core_InstanceSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_InstanceStatus_To_core_InstanceStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Instance_To_core_Instance is an autogenerated conversion function.
func Convert_v1alpha1_Instance_To_core_Instance(in *Instance, out *core.Instance, s conversion.Scope) error {
	return autoConvert_v1alpha1_Instance_To_core_Instance(in, out, s)
}

func autoConvert_core_Instance_To_v1alpha1_Instance(in *core.Instance, out *Instance, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_core_InstanceSpec_To_v1alpha1_InstanceSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_core_InstanceStatus_To_v1alpha1_InstanceStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_core_Instance_To_v1alpha1_Instance is an autogenerated conversion function.
func Convert_core_Instance_To_v1alpha1_Instance(in *core.Instance, out *Instance, s conversion.Scope) error {
	return autoConvert_core_Instance_To_v1alpha1_Instance(in, out, s)
}

func autoConvert_v1alpha1_InstanceList_To_core_InstanceList(in *InstanceList, out *core.InstanceList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]core.Instance)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_InstanceList_To_core_InstanceList is an autogenerated conversion function.
func Convert_v1alpha1_InstanceList_To_core_InstanceList(in *InstanceList, out *core.InstanceList, s conversion.Scope) error {
	return autoConvert_v1alpha1_InstanceList_To_core_InstanceList(in, out, s)
}

func autoConvert_core_InstanceList_To_v1alpha1_InstanceList(in *core.InstanceList, out *InstanceList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]Instance)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_core_InstanceList_To_v1alpha1_InstanceList is an autogenerated conversion function.
func Convert_core_InstanceList_To_v1alpha1_InstanceList(in *core.InstanceList, out *InstanceList, s conversion.Scope) error {
	return autoConvert_core_InstanceList_To_v1alpha1_InstanceList(in, out, s)
}

func autoConvert_v1alpha1_InstanceSpec_To_core_InstanceSpec(in *InstanceSpec, out *core.InstanceSpec, s conversion.Scope) error {
	if err := Convert_v1alpha1_LandscaperConfiguration_To_core_LandscaperConfiguration(&in.LandscaperConfiguration, &out.LandscaperConfiguration, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_LandscaperServiceComponentReference_To_core_LandscaperServiceComponentReference(&in.ComponentReference, &out.ComponentReference, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ObjectReference_To_core_ObjectReference(&in.ServiceTargetConfigRef, &out.ServiceTargetConfigRef, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_InstanceSpec_To_core_InstanceSpec is an autogenerated conversion function.
func Convert_v1alpha1_InstanceSpec_To_core_InstanceSpec(in *InstanceSpec, out *core.InstanceSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_InstanceSpec_To_core_InstanceSpec(in, out, s)
}

func autoConvert_core_InstanceSpec_To_v1alpha1_InstanceSpec(in *core.InstanceSpec, out *InstanceSpec, s conversion.Scope) error {
	if err := Convert_core_LandscaperConfiguration_To_v1alpha1_LandscaperConfiguration(&in.LandscaperConfiguration, &out.LandscaperConfiguration, s); err != nil {
		return err
	}
	if err := Convert_core_LandscaperServiceComponentReference_To_v1alpha1_LandscaperServiceComponentReference(&in.ComponentReference, &out.ComponentReference, s); err != nil {
		return err
	}
	if err := Convert_core_ObjectReference_To_v1alpha1_ObjectReference(&in.ServiceTargetConfigRef, &out.ServiceTargetConfigRef, s); err != nil {
		return err
	}
	return nil
}

// Convert_core_InstanceSpec_To_v1alpha1_InstanceSpec is an autogenerated conversion function.
func Convert_core_InstanceSpec_To_v1alpha1_InstanceSpec(in *core.InstanceSpec, out *InstanceSpec, s conversion.Scope) error {
	return autoConvert_core_InstanceSpec_To_v1alpha1_InstanceSpec(in, out, s)
}

func autoConvert_v1alpha1_InstanceStatus_To_core_InstanceStatus(in *InstanceStatus, out *core.InstanceStatus, s conversion.Scope) error {
	out.ObservedGeneration = in.ObservedGeneration
	out.LastError = (*core.Error)(unsafe.Pointer(in.LastError))
	out.TargetRef = (*core.ObjectReference)(unsafe.Pointer(in.TargetRef))
	out.InstallationRef = (*core.ObjectReference)(unsafe.Pointer(in.InstallationRef))
	out.ClusterEndpoint = in.ClusterEndpoint
	out.ClusterKubeconfig = in.ClusterKubeconfig
	return nil
}

// Convert_v1alpha1_InstanceStatus_To_core_InstanceStatus is an autogenerated conversion function.
func Convert_v1alpha1_InstanceStatus_To_core_InstanceStatus(in *InstanceStatus, out *core.InstanceStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_InstanceStatus_To_core_InstanceStatus(in, out, s)
}

func autoConvert_core_InstanceStatus_To_v1alpha1_InstanceStatus(in *core.InstanceStatus, out *InstanceStatus, s conversion.Scope) error {
	out.ObservedGeneration = in.ObservedGeneration
	out.LastError = (*Error)(unsafe.Pointer(in.LastError))
	out.TargetRef = (*ObjectReference)(unsafe.Pointer(in.TargetRef))
	out.InstallationRef = (*ObjectReference)(unsafe.Pointer(in.InstallationRef))
	out.ClusterEndpoint = in.ClusterEndpoint
	out.ClusterKubeconfig = in.ClusterKubeconfig
	return nil
}

// Convert_core_InstanceStatus_To_v1alpha1_InstanceStatus is an autogenerated conversion function.
func Convert_core_InstanceStatus_To_v1alpha1_InstanceStatus(in *core.InstanceStatus, out *InstanceStatus, s conversion.Scope) error {
	return autoConvert_core_InstanceStatus_To_v1alpha1_InstanceStatus(in, out, s)
}

func autoConvert_v1alpha1_LandscaperConfiguration_To_core_LandscaperConfiguration(in *LandscaperConfiguration, out *core.LandscaperConfiguration, s conversion.Scope) error {
	out.Deployers = *(*[]string)(unsafe.Pointer(&in.Deployers))
	return nil
}

// Convert_v1alpha1_LandscaperConfiguration_To_core_LandscaperConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_LandscaperConfiguration_To_core_LandscaperConfiguration(in *LandscaperConfiguration, out *core.LandscaperConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_LandscaperConfiguration_To_core_LandscaperConfiguration(in, out, s)
}

func autoConvert_core_LandscaperConfiguration_To_v1alpha1_LandscaperConfiguration(in *core.LandscaperConfiguration, out *LandscaperConfiguration, s conversion.Scope) error {
	out.Deployers = *(*[]string)(unsafe.Pointer(&in.Deployers))
	return nil
}

// Convert_core_LandscaperConfiguration_To_v1alpha1_LandscaperConfiguration is an autogenerated conversion function.
func Convert_core_LandscaperConfiguration_To_v1alpha1_LandscaperConfiguration(in *core.LandscaperConfiguration, out *LandscaperConfiguration, s conversion.Scope) error {
	return autoConvert_core_LandscaperConfiguration_To_v1alpha1_LandscaperConfiguration(in, out, s)
}

func autoConvert_v1alpha1_LandscaperDeployment_To_core_LandscaperDeployment(in *LandscaperDeployment, out *core.LandscaperDeployment, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_LandscaperDeploymentSpec_To_core_LandscaperDeploymentSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_LandscaperDeploymentStatus_To_core_LandscaperDeploymentStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_LandscaperDeployment_To_core_LandscaperDeployment is an autogenerated conversion function.
func Convert_v1alpha1_LandscaperDeployment_To_core_LandscaperDeployment(in *LandscaperDeployment, out *core.LandscaperDeployment, s conversion.Scope) error {
	return autoConvert_v1alpha1_LandscaperDeployment_To_core_LandscaperDeployment(in, out, s)
}

func autoConvert_core_LandscaperDeployment_To_v1alpha1_LandscaperDeployment(in *core.LandscaperDeployment, out *LandscaperDeployment, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_core_LandscaperDeploymentSpec_To_v1alpha1_LandscaperDeploymentSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_core_LandscaperDeploymentStatus_To_v1alpha1_LandscaperDeploymentStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_core_LandscaperDeployment_To_v1alpha1_LandscaperDeployment is an autogenerated conversion function.
func Convert_core_LandscaperDeployment_To_v1alpha1_LandscaperDeployment(in *core.LandscaperDeployment, out *LandscaperDeployment, s conversion.Scope) error {
	return autoConvert_core_LandscaperDeployment_To_v1alpha1_LandscaperDeployment(in, out, s)
}

func autoConvert_v1alpha1_LandscaperDeploymentList_To_core_LandscaperDeploymentList(in *LandscaperDeploymentList, out *core.LandscaperDeploymentList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]core.LandscaperDeployment)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_LandscaperDeploymentList_To_core_LandscaperDeploymentList is an autogenerated conversion function.
func Convert_v1alpha1_LandscaperDeploymentList_To_core_LandscaperDeploymentList(in *LandscaperDeploymentList, out *core.LandscaperDeploymentList, s conversion.Scope) error {
	return autoConvert_v1alpha1_LandscaperDeploymentList_To_core_LandscaperDeploymentList(in, out, s)
}

func autoConvert_core_LandscaperDeploymentList_To_v1alpha1_LandscaperDeploymentList(in *core.LandscaperDeploymentList, out *LandscaperDeploymentList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]LandscaperDeployment)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_core_LandscaperDeploymentList_To_v1alpha1_LandscaperDeploymentList is an autogenerated conversion function.
func Convert_core_LandscaperDeploymentList_To_v1alpha1_LandscaperDeploymentList(in *core.LandscaperDeploymentList, out *LandscaperDeploymentList, s conversion.Scope) error {
	return autoConvert_core_LandscaperDeploymentList_To_v1alpha1_LandscaperDeploymentList(in, out, s)
}

func autoConvert_v1alpha1_LandscaperDeploymentSpec_To_core_LandscaperDeploymentSpec(in *LandscaperDeploymentSpec, out *core.LandscaperDeploymentSpec, s conversion.Scope) error {
	out.Purpose = in.Purpose
	if err := Convert_v1alpha1_LandscaperConfiguration_To_core_LandscaperConfiguration(&in.LandscaperConfiguration, &out.LandscaperConfiguration, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_LandscaperServiceComponentReference_To_core_LandscaperServiceComponentReference(&in.ComponentReference, &out.ComponentReference, s); err != nil {
		return err
	}
	out.Region = in.Region
	return nil
}

// Convert_v1alpha1_LandscaperDeploymentSpec_To_core_LandscaperDeploymentSpec is an autogenerated conversion function.
func Convert_v1alpha1_LandscaperDeploymentSpec_To_core_LandscaperDeploymentSpec(in *LandscaperDeploymentSpec, out *core.LandscaperDeploymentSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_LandscaperDeploymentSpec_To_core_LandscaperDeploymentSpec(in, out, s)
}

func autoConvert_core_LandscaperDeploymentSpec_To_v1alpha1_LandscaperDeploymentSpec(in *core.LandscaperDeploymentSpec, out *LandscaperDeploymentSpec, s conversion.Scope) error {
	out.Purpose = in.Purpose
	if err := Convert_core_LandscaperConfiguration_To_v1alpha1_LandscaperConfiguration(&in.LandscaperConfiguration, &out.LandscaperConfiguration, s); err != nil {
		return err
	}
	if err := Convert_core_LandscaperServiceComponentReference_To_v1alpha1_LandscaperServiceComponentReference(&in.ComponentReference, &out.ComponentReference, s); err != nil {
		return err
	}
	out.Region = in.Region
	return nil
}

// Convert_core_LandscaperDeploymentSpec_To_v1alpha1_LandscaperDeploymentSpec is an autogenerated conversion function.
func Convert_core_LandscaperDeploymentSpec_To_v1alpha1_LandscaperDeploymentSpec(in *core.LandscaperDeploymentSpec, out *LandscaperDeploymentSpec, s conversion.Scope) error {
	return autoConvert_core_LandscaperDeploymentSpec_To_v1alpha1_LandscaperDeploymentSpec(in, out, s)
}

func autoConvert_v1alpha1_LandscaperDeploymentStatus_To_core_LandscaperDeploymentStatus(in *LandscaperDeploymentStatus, out *core.LandscaperDeploymentStatus, s conversion.Scope) error {
	out.ObservedGeneration = in.ObservedGeneration
	out.LastError = (*core.Error)(unsafe.Pointer(in.LastError))
	out.InstanceRef = (*core.ObjectReference)(unsafe.Pointer(in.InstanceRef))
	return nil
}

// Convert_v1alpha1_LandscaperDeploymentStatus_To_core_LandscaperDeploymentStatus is an autogenerated conversion function.
func Convert_v1alpha1_LandscaperDeploymentStatus_To_core_LandscaperDeploymentStatus(in *LandscaperDeploymentStatus, out *core.LandscaperDeploymentStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_LandscaperDeploymentStatus_To_core_LandscaperDeploymentStatus(in, out, s)
}

func autoConvert_core_LandscaperDeploymentStatus_To_v1alpha1_LandscaperDeploymentStatus(in *core.LandscaperDeploymentStatus, out *LandscaperDeploymentStatus, s conversion.Scope) error {
	out.ObservedGeneration = in.ObservedGeneration
	out.LastError = (*Error)(unsafe.Pointer(in.LastError))
	out.InstanceRef = (*ObjectReference)(unsafe.Pointer(in.InstanceRef))
	return nil
}

// Convert_core_LandscaperDeploymentStatus_To_v1alpha1_LandscaperDeploymentStatus is an autogenerated conversion function.
func Convert_core_LandscaperDeploymentStatus_To_v1alpha1_LandscaperDeploymentStatus(in *core.LandscaperDeploymentStatus, out *LandscaperDeploymentStatus, s conversion.Scope) error {
	return autoConvert_core_LandscaperDeploymentStatus_To_v1alpha1_LandscaperDeploymentStatus(in, out, s)
}

func autoConvert_v1alpha1_LandscaperServiceComponentReference_To_core_LandscaperServiceComponentReference(in *LandscaperServiceComponentReference, out *core.LandscaperServiceComponentReference, s conversion.Scope) error {
	out.Context = in.Context
	out.ComponentName = in.ComponentName
	out.Version = in.Version
	return nil
}

// Convert_v1alpha1_LandscaperServiceComponentReference_To_core_LandscaperServiceComponentReference is an autogenerated conversion function.
func Convert_v1alpha1_LandscaperServiceComponentReference_To_core_LandscaperServiceComponentReference(in *LandscaperServiceComponentReference, out *core.LandscaperServiceComponentReference, s conversion.Scope) error {
	return autoConvert_v1alpha1_LandscaperServiceComponentReference_To_core_LandscaperServiceComponentReference(in, out, s)
}

func autoConvert_core_LandscaperServiceComponentReference_To_v1alpha1_LandscaperServiceComponentReference(in *core.LandscaperServiceComponentReference, out *LandscaperServiceComponentReference, s conversion.Scope) error {
	out.Context = in.Context
	out.ComponentName = in.ComponentName
	out.Version = in.Version
	return nil
}

// Convert_core_LandscaperServiceComponentReference_To_v1alpha1_LandscaperServiceComponentReference is an autogenerated conversion function.
func Convert_core_LandscaperServiceComponentReference_To_v1alpha1_LandscaperServiceComponentReference(in *core.LandscaperServiceComponentReference, out *LandscaperServiceComponentReference, s conversion.Scope) error {
	return autoConvert_core_LandscaperServiceComponentReference_To_v1alpha1_LandscaperServiceComponentReference(in, out, s)
}

func autoConvert_v1alpha1_ObjectReference_To_core_ObjectReference(in *ObjectReference, out *core.ObjectReference, s conversion.Scope) error {
	out.Name = in.Name
	out.Namespace = in.Namespace
	return nil
}

// Convert_v1alpha1_ObjectReference_To_core_ObjectReference is an autogenerated conversion function.
func Convert_v1alpha1_ObjectReference_To_core_ObjectReference(in *ObjectReference, out *core.ObjectReference, s conversion.Scope) error {
	return autoConvert_v1alpha1_ObjectReference_To_core_ObjectReference(in, out, s)
}

func autoConvert_core_ObjectReference_To_v1alpha1_ObjectReference(in *core.ObjectReference, out *ObjectReference, s conversion.Scope) error {
	out.Name = in.Name
	out.Namespace = in.Namespace
	return nil
}

// Convert_core_ObjectReference_To_v1alpha1_ObjectReference is an autogenerated conversion function.
func Convert_core_ObjectReference_To_v1alpha1_ObjectReference(in *core.ObjectReference, out *ObjectReference, s conversion.Scope) error {
	return autoConvert_core_ObjectReference_To_v1alpha1_ObjectReference(in, out, s)
}

func autoConvert_v1alpha1_SecretReference_To_core_SecretReference(in *SecretReference, out *core.SecretReference, s conversion.Scope) error {
	if err := Convert_v1alpha1_ObjectReference_To_core_ObjectReference(&in.ObjectReference, &out.ObjectReference, s); err != nil {
		return err
	}
	out.Key = in.Key
	return nil
}

// Convert_v1alpha1_SecretReference_To_core_SecretReference is an autogenerated conversion function.
func Convert_v1alpha1_SecretReference_To_core_SecretReference(in *SecretReference, out *core.SecretReference, s conversion.Scope) error {
	return autoConvert_v1alpha1_SecretReference_To_core_SecretReference(in, out, s)
}

func autoConvert_core_SecretReference_To_v1alpha1_SecretReference(in *core.SecretReference, out *SecretReference, s conversion.Scope) error {
	if err := Convert_core_ObjectReference_To_v1alpha1_ObjectReference(&in.ObjectReference, &out.ObjectReference, s); err != nil {
		return err
	}
	out.Key = in.Key
	return nil
}

// Convert_core_SecretReference_To_v1alpha1_SecretReference is an autogenerated conversion function.
func Convert_core_SecretReference_To_v1alpha1_SecretReference(in *core.SecretReference, out *SecretReference, s conversion.Scope) error {
	return autoConvert_core_SecretReference_To_v1alpha1_SecretReference(in, out, s)
}

func autoConvert_v1alpha1_ServiceTargetConfig_To_core_ServiceTargetConfig(in *ServiceTargetConfig, out *core.ServiceTargetConfig, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_ServiceTargetConfigSpec_To_core_ServiceTargetConfigSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ServiceTargetConfigStatus_To_core_ServiceTargetConfigStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_ServiceTargetConfig_To_core_ServiceTargetConfig is an autogenerated conversion function.
func Convert_v1alpha1_ServiceTargetConfig_To_core_ServiceTargetConfig(in *ServiceTargetConfig, out *core.ServiceTargetConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_ServiceTargetConfig_To_core_ServiceTargetConfig(in, out, s)
}

func autoConvert_core_ServiceTargetConfig_To_v1alpha1_ServiceTargetConfig(in *core.ServiceTargetConfig, out *ServiceTargetConfig, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_core_ServiceTargetConfigSpec_To_v1alpha1_ServiceTargetConfigSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_core_ServiceTargetConfigStatus_To_v1alpha1_ServiceTargetConfigStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_core_ServiceTargetConfig_To_v1alpha1_ServiceTargetConfig is an autogenerated conversion function.
func Convert_core_ServiceTargetConfig_To_v1alpha1_ServiceTargetConfig(in *core.ServiceTargetConfig, out *ServiceTargetConfig, s conversion.Scope) error {
	return autoConvert_core_ServiceTargetConfig_To_v1alpha1_ServiceTargetConfig(in, out, s)
}

func autoConvert_v1alpha1_ServiceTargetConfigList_To_core_ServiceTargetConfigList(in *ServiceTargetConfigList, out *core.ServiceTargetConfigList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]core.ServiceTargetConfig)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_ServiceTargetConfigList_To_core_ServiceTargetConfigList is an autogenerated conversion function.
func Convert_v1alpha1_ServiceTargetConfigList_To_core_ServiceTargetConfigList(in *ServiceTargetConfigList, out *core.ServiceTargetConfigList, s conversion.Scope) error {
	return autoConvert_v1alpha1_ServiceTargetConfigList_To_core_ServiceTargetConfigList(in, out, s)
}

func autoConvert_core_ServiceTargetConfigList_To_v1alpha1_ServiceTargetConfigList(in *core.ServiceTargetConfigList, out *ServiceTargetConfigList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]ServiceTargetConfig)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_core_ServiceTargetConfigList_To_v1alpha1_ServiceTargetConfigList is an autogenerated conversion function.
func Convert_core_ServiceTargetConfigList_To_v1alpha1_ServiceTargetConfigList(in *core.ServiceTargetConfigList, out *ServiceTargetConfigList, s conversion.Scope) error {
	return autoConvert_core_ServiceTargetConfigList_To_v1alpha1_ServiceTargetConfigList(in, out, s)
}

func autoConvert_v1alpha1_ServiceTargetConfigSpec_To_core_ServiceTargetConfigSpec(in *ServiceTargetConfigSpec, out *core.ServiceTargetConfigSpec, s conversion.Scope) error {
	out.ProviderType = in.ProviderType
	out.Region = in.Region
	out.Priority = in.Priority
	out.Visible = in.Visible
	if err := Convert_v1alpha1_SecretReference_To_core_SecretReference(&in.SecretRef, &out.SecretRef, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_ServiceTargetConfigSpec_To_core_ServiceTargetConfigSpec is an autogenerated conversion function.
func Convert_v1alpha1_ServiceTargetConfigSpec_To_core_ServiceTargetConfigSpec(in *ServiceTargetConfigSpec, out *core.ServiceTargetConfigSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_ServiceTargetConfigSpec_To_core_ServiceTargetConfigSpec(in, out, s)
}

func autoConvert_core_ServiceTargetConfigSpec_To_v1alpha1_ServiceTargetConfigSpec(in *core.ServiceTargetConfigSpec, out *ServiceTargetConfigSpec, s conversion.Scope) error {
	out.ProviderType = in.ProviderType
	out.Region = in.Region
	out.Priority = in.Priority
	out.Visible = in.Visible
	if err := Convert_core_SecretReference_To_v1alpha1_SecretReference(&in.SecretRef, &out.SecretRef, s); err != nil {
		return err
	}
	return nil
}

// Convert_core_ServiceTargetConfigSpec_To_v1alpha1_ServiceTargetConfigSpec is an autogenerated conversion function.
func Convert_core_ServiceTargetConfigSpec_To_v1alpha1_ServiceTargetConfigSpec(in *core.ServiceTargetConfigSpec, out *ServiceTargetConfigSpec, s conversion.Scope) error {
	return autoConvert_core_ServiceTargetConfigSpec_To_v1alpha1_ServiceTargetConfigSpec(in, out, s)
}

func autoConvert_v1alpha1_ServiceTargetConfigStatus_To_core_ServiceTargetConfigStatus(in *ServiceTargetConfigStatus, out *core.ServiceTargetConfigStatus, s conversion.Scope) error {
	out.ObservedGeneration = in.ObservedGeneration
	out.Capacity = (*int64)(unsafe.Pointer(in.Capacity))
	out.InstanceRefs = *(*[]core.ObjectReference)(unsafe.Pointer(&in.InstanceRefs))
	return nil
}

// Convert_v1alpha1_ServiceTargetConfigStatus_To_core_ServiceTargetConfigStatus is an autogenerated conversion function.
func Convert_v1alpha1_ServiceTargetConfigStatus_To_core_ServiceTargetConfigStatus(in *ServiceTargetConfigStatus, out *core.ServiceTargetConfigStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_ServiceTargetConfigStatus_To_core_ServiceTargetConfigStatus(in, out, s)
}

func autoConvert_core_ServiceTargetConfigStatus_To_v1alpha1_ServiceTargetConfigStatus(in *core.ServiceTargetConfigStatus, out *ServiceTargetConfigStatus, s conversion.Scope) error {
	out.ObservedGeneration = in.ObservedGeneration
	out.Capacity = (*int64)(unsafe.Pointer(in.Capacity))
	out.InstanceRefs = *(*[]ObjectReference)(unsafe.Pointer(&in.InstanceRefs))
	return nil
}

// Convert_core_ServiceTargetConfigStatus_To_v1alpha1_ServiceTargetConfigStatus is an autogenerated conversion function.
func Convert_core_ServiceTargetConfigStatus_To_v1alpha1_ServiceTargetConfigStatus(in *core.ServiceTargetConfigStatus, out *ServiceTargetConfigStatus, s conversion.Scope) error {
	return autoConvert_core_ServiceTargetConfigStatus_To_v1alpha1_ServiceTargetConfigStatus(in, out, s)
}
