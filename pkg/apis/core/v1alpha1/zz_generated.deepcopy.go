//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

SPDX-License-Identifier: Apache-2.0
*/
// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutomaticReconcile) DeepCopyInto(out *AutomaticReconcile) {
	*out = *in
	out.Interval = in.Interval
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutomaticReconcile.
func (in *AutomaticReconcile) DeepCopy() *AutomaticReconcile {
	if in == nil {
		return nil
	}
	out := new(AutomaticReconcile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutomaticReconcileStatus) DeepCopyInto(out *AutomaticReconcileStatus) {
	*out = *in
	in.LastReconcileTime.DeepCopyInto(&out.LastReconcileTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutomaticReconcileStatus.
func (in *AutomaticReconcileStatus) DeepCopy() *AutomaticReconcileStatus {
	if in == nil {
		return nil
	}
	out := new(AutomaticReconcileStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AvailabilityCollection) DeepCopyInto(out *AvailabilityCollection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AvailabilityCollection.
func (in *AvailabilityCollection) DeepCopy() *AvailabilityCollection {
	if in == nil {
		return nil
	}
	out := new(AvailabilityCollection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AvailabilityCollection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AvailabilityCollectionList) DeepCopyInto(out *AvailabilityCollectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AvailabilityCollection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AvailabilityCollectionList.
func (in *AvailabilityCollectionList) DeepCopy() *AvailabilityCollectionList {
	if in == nil {
		return nil
	}
	out := new(AvailabilityCollectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AvailabilityCollectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AvailabilityCollectionSpec) DeepCopyInto(out *AvailabilityCollectionSpec) {
	*out = *in
	if in.InstanceRefs != nil {
		in, out := &in.InstanceRefs, &out.InstanceRefs
		*out = make([]ObjectReference, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AvailabilityCollectionSpec.
func (in *AvailabilityCollectionSpec) DeepCopy() *AvailabilityCollectionSpec {
	if in == nil {
		return nil
	}
	out := new(AvailabilityCollectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AvailabilityCollectionStatus) DeepCopyInto(out *AvailabilityCollectionStatus) {
	*out = *in
	in.LastRun.DeepCopyInto(&out.LastRun)
	in.LastReported.DeepCopyInto(&out.LastReported)
	if in.Instances != nil {
		in, out := &in.Instances, &out.Instances
		*out = make([]AvailabilityInstance, len(*in))
		copy(*out, *in)
	}
	out.Self = in.Self
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AvailabilityCollectionStatus.
func (in *AvailabilityCollectionStatus) DeepCopy() *AvailabilityCollectionStatus {
	if in == nil {
		return nil
	}
	out := new(AvailabilityCollectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AvailabilityInstance) DeepCopyInto(out *AvailabilityInstance) {
	*out = *in
	out.ObjectReference = in.ObjectReference
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AvailabilityInstance.
func (in *AvailabilityInstance) DeepCopy() *AvailabilityInstance {
	if in == nil {
		return nil
	}
	out := new(AvailabilityInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Error) DeepCopyInto(out *Error) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Error.
func (in *Error) DeepCopy() *Error {
	if in == nil {
		return nil
	}
	out := new(Error)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HighAvailabilityConfig) DeepCopyInto(out *HighAvailabilityConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HighAvailabilityConfig.
func (in *HighAvailabilityConfig) DeepCopy() *HighAvailabilityConfig {
	if in == nil {
		return nil
	}
	out := new(HighAvailabilityConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Instance) DeepCopyInto(out *Instance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Instance.
func (in *Instance) DeepCopy() *Instance {
	if in == nil {
		return nil
	}
	out := new(Instance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Instance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceList) DeepCopyInto(out *InstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Instance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceList.
func (in *InstanceList) DeepCopy() *InstanceList {
	if in == nil {
		return nil
	}
	out := new(InstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceRegistration) DeepCopyInto(out *InstanceRegistration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceRegistration.
func (in *InstanceRegistration) DeepCopy() *InstanceRegistration {
	if in == nil {
		return nil
	}
	out := new(InstanceRegistration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstanceRegistration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceRegistrationList) DeepCopyInto(out *InstanceRegistrationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InstanceRegistration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceRegistrationList.
func (in *InstanceRegistrationList) DeepCopy() *InstanceRegistrationList {
	if in == nil {
		return nil
	}
	out := new(InstanceRegistrationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstanceRegistrationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceRegistrationSpec) DeepCopyInto(out *InstanceRegistrationSpec) {
	*out = *in
	in.LandscaperDeploymentSpec.DeepCopyInto(&out.LandscaperDeploymentSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceRegistrationSpec.
func (in *InstanceRegistrationSpec) DeepCopy() *InstanceRegistrationSpec {
	if in == nil {
		return nil
	}
	out := new(InstanceRegistrationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceRegistrationStatus) DeepCopyInto(out *InstanceRegistrationStatus) {
	*out = *in
	if in.LandscaperDeploymentInfo != nil {
		in, out := &in.LandscaperDeploymentInfo, &out.LandscaperDeploymentInfo
		*out = new(types.NamespacedName)
		**out = **in
	}
	if in.LastError != nil {
		in, out := &in.LastError, &out.LastError
		*out = new(Error)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceRegistrationStatus.
func (in *InstanceRegistrationStatus) DeepCopy() *InstanceRegistrationStatus {
	if in == nil {
		return nil
	}
	out := new(InstanceRegistrationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceSpec) DeepCopyInto(out *InstanceSpec) {
	*out = *in
	in.LandscaperConfiguration.DeepCopyInto(&out.LandscaperConfiguration)
	out.ServiceTargetConfigRef = in.ServiceTargetConfigRef
	if in.OIDCConfig != nil {
		in, out := &in.OIDCConfig, &out.OIDCConfig
		*out = new(OIDCConfig)
		**out = **in
	}
	if in.AutomaticReconcile != nil {
		in, out := &in.AutomaticReconcile, &out.AutomaticReconcile
		*out = new(AutomaticReconcile)
		**out = **in
	}
	if in.HighAvailabilityConfig != nil {
		in, out := &in.HighAvailabilityConfig, &out.HighAvailabilityConfig
		*out = new(HighAvailabilityConfig)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceSpec.
func (in *InstanceSpec) DeepCopy() *InstanceSpec {
	if in == nil {
		return nil
	}
	out := new(InstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceStatus) DeepCopyInto(out *InstanceStatus) {
	*out = *in
	if in.LastError != nil {
		in, out := &in.LastError, &out.LastError
		*out = new(Error)
		(*in).DeepCopyInto(*out)
	}
	if in.LandscaperServiceComponent != nil {
		in, out := &in.LandscaperServiceComponent, &out.LandscaperServiceComponent
		*out = new(LandscaperServiceComponent)
		**out = **in
	}
	if in.ContextRef != nil {
		in, out := &in.ContextRef, &out.ContextRef
		*out = new(ObjectReference)
		**out = **in
	}
	if in.TargetRef != nil {
		in, out := &in.TargetRef, &out.TargetRef
		*out = new(ObjectReference)
		**out = **in
	}
	if in.GardenerServiceAccountRef != nil {
		in, out := &in.GardenerServiceAccountRef, &out.GardenerServiceAccountRef
		*out = new(ObjectReference)
		**out = **in
	}
	if in.InstallationRef != nil {
		in, out := &in.InstallationRef, &out.InstallationRef
		*out = new(ObjectReference)
		**out = **in
	}
	if in.AutomaticReconcileStatus != nil {
		in, out := &in.AutomaticReconcileStatus, &out.AutomaticReconcileStatus
		*out = new(AutomaticReconcileStatus)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceStatus.
func (in *InstanceStatus) DeepCopy() *InstanceStatus {
	if in == nil {
		return nil
	}
	out := new(InstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LandscaperConfiguration) DeepCopyInto(out *LandscaperConfiguration) {
	*out = *in
	if in.Deployers != nil {
		in, out := &in.Deployers, &out.Deployers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LandscaperConfiguration.
func (in *LandscaperConfiguration) DeepCopy() *LandscaperConfiguration {
	if in == nil {
		return nil
	}
	out := new(LandscaperConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LandscaperDeployment) DeepCopyInto(out *LandscaperDeployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LandscaperDeployment.
func (in *LandscaperDeployment) DeepCopy() *LandscaperDeployment {
	if in == nil {
		return nil
	}
	out := new(LandscaperDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LandscaperDeployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LandscaperDeploymentList) DeepCopyInto(out *LandscaperDeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LandscaperDeployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LandscaperDeploymentList.
func (in *LandscaperDeploymentList) DeepCopy() *LandscaperDeploymentList {
	if in == nil {
		return nil
	}
	out := new(LandscaperDeploymentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LandscaperDeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LandscaperDeploymentSpec) DeepCopyInto(out *LandscaperDeploymentSpec) {
	*out = *in
	in.LandscaperConfiguration.DeepCopyInto(&out.LandscaperConfiguration)
	if in.OIDCConfig != nil {
		in, out := &in.OIDCConfig, &out.OIDCConfig
		*out = new(OIDCConfig)
		**out = **in
	}
	if in.HighAvailabilityConfig != nil {
		in, out := &in.HighAvailabilityConfig, &out.HighAvailabilityConfig
		*out = new(HighAvailabilityConfig)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LandscaperDeploymentSpec.
func (in *LandscaperDeploymentSpec) DeepCopy() *LandscaperDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(LandscaperDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LandscaperDeploymentStatus) DeepCopyInto(out *LandscaperDeploymentStatus) {
	*out = *in
	if in.LastError != nil {
		in, out := &in.LastError, &out.LastError
		*out = new(Error)
		(*in).DeepCopyInto(*out)
	}
	if in.InstanceRef != nil {
		in, out := &in.InstanceRef, &out.InstanceRef
		*out = new(ObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LandscaperDeploymentStatus.
func (in *LandscaperDeploymentStatus) DeepCopy() *LandscaperDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(LandscaperDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LandscaperServiceComponent) DeepCopyInto(out *LandscaperServiceComponent) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LandscaperServiceComponent.
func (in *LandscaperServiceComponent) DeepCopy() *LandscaperServiceComponent {
	if in == nil {
		return nil
	}
	out := new(LandscaperServiceComponent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LosSubject) DeepCopyInto(out *LosSubject) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LosSubject.
func (in *LosSubject) DeepCopy() *LosSubject {
	if in == nil {
		return nil
	}
	out := new(LosSubject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LosSubjectList) DeepCopyInto(out *LosSubjectList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LosSubjectList.
func (in *LosSubjectList) DeepCopy() *LosSubjectList {
	if in == nil {
		return nil
	}
	out := new(LosSubjectList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LosSubjectList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LosSubjectListList) DeepCopyInto(out *LosSubjectListList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LosSubjectList, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LosSubjectListList.
func (in *LosSubjectListList) DeepCopy() *LosSubjectListList {
	if in == nil {
		return nil
	}
	out := new(LosSubjectListList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LosSubjectListList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LosSubjectListSpec) DeepCopyInto(out *LosSubjectListSpec) {
	*out = *in
	if in.Admins != nil {
		in, out := &in.Admins, &out.Admins
		*out = make([]LosSubject, len(*in))
		copy(*out, *in)
	}
	if in.Members != nil {
		in, out := &in.Members, &out.Members
		*out = make([]LosSubject, len(*in))
		copy(*out, *in)
	}
	if in.Viewer != nil {
		in, out := &in.Viewer, &out.Viewer
		*out = make([]LosSubject, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LosSubjectListSpec.
func (in *LosSubjectListSpec) DeepCopy() *LosSubjectListSpec {
	if in == nil {
		return nil
	}
	out := new(LosSubjectListSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LosSubjectListStatus) DeepCopyInto(out *LosSubjectListStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LosSubjectListStatus.
func (in *LosSubjectListStatus) DeepCopy() *LosSubjectListStatus {
	if in == nil {
		return nil
	}
	out := new(LosSubjectListStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceRegistration) DeepCopyInto(out *NamespaceRegistration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceRegistration.
func (in *NamespaceRegistration) DeepCopy() *NamespaceRegistration {
	if in == nil {
		return nil
	}
	out := new(NamespaceRegistration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NamespaceRegistration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceRegistrationList) DeepCopyInto(out *NamespaceRegistrationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NamespaceRegistration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceRegistrationList.
func (in *NamespaceRegistrationList) DeepCopy() *NamespaceRegistrationList {
	if in == nil {
		return nil
	}
	out := new(NamespaceRegistrationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NamespaceRegistrationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceRegistrationSpec) DeepCopyInto(out *NamespaceRegistrationSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceRegistrationSpec.
func (in *NamespaceRegistrationSpec) DeepCopy() *NamespaceRegistrationSpec {
	if in == nil {
		return nil
	}
	out := new(NamespaceRegistrationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceRegistrationStatus) DeepCopyInto(out *NamespaceRegistrationStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceRegistrationStatus.
func (in *NamespaceRegistrationStatus) DeepCopy() *NamespaceRegistrationStatus {
	if in == nil {
		return nil
	}
	out := new(NamespaceRegistrationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OIDCConfig) DeepCopyInto(out *OIDCConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OIDCConfig.
func (in *OIDCConfig) DeepCopy() *OIDCConfig {
	if in == nil {
		return nil
	}
	out := new(OIDCConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectReference) DeepCopyInto(out *ObjectReference) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectReference.
func (in *ObjectReference) DeepCopy() *ObjectReference {
	if in == nil {
		return nil
	}
	out := new(ObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretReference) DeepCopyInto(out *SecretReference) {
	*out = *in
	out.ObjectReference = in.ObjectReference
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretReference.
func (in *SecretReference) DeepCopy() *SecretReference {
	if in == nil {
		return nil
	}
	out := new(SecretReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceTargetConfig) DeepCopyInto(out *ServiceTargetConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceTargetConfig.
func (in *ServiceTargetConfig) DeepCopy() *ServiceTargetConfig {
	if in == nil {
		return nil
	}
	out := new(ServiceTargetConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceTargetConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceTargetConfigList) DeepCopyInto(out *ServiceTargetConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceTargetConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceTargetConfigList.
func (in *ServiceTargetConfigList) DeepCopy() *ServiceTargetConfigList {
	if in == nil {
		return nil
	}
	out := new(ServiceTargetConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceTargetConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceTargetConfigSpec) DeepCopyInto(out *ServiceTargetConfigSpec) {
	*out = *in
	out.SecretRef = in.SecretRef
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceTargetConfigSpec.
func (in *ServiceTargetConfigSpec) DeepCopy() *ServiceTargetConfigSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceTargetConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceTargetConfigStatus) DeepCopyInto(out *ServiceTargetConfigStatus) {
	*out = *in
	if in.InstanceRefs != nil {
		in, out := &in.InstanceRefs, &out.InstanceRefs
		*out = make([]ObjectReference, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceTargetConfigStatus.
func (in *ServiceTargetConfigStatus) DeepCopy() *ServiceTargetConfigStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceTargetConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Subject) DeepCopyInto(out *Subject) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Subject.
func (in *Subject) DeepCopy() *Subject {
	if in == nil {
		return nil
	}
	out := new(Subject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubjectList) DeepCopyInto(out *SubjectList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubjectList.
func (in *SubjectList) DeepCopy() *SubjectList {
	if in == nil {
		return nil
	}
	out := new(SubjectList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SubjectList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubjectListList) DeepCopyInto(out *SubjectListList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SubjectList, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubjectListList.
func (in *SubjectListList) DeepCopy() *SubjectListList {
	if in == nil {
		return nil
	}
	out := new(SubjectListList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SubjectListList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubjectListSpec) DeepCopyInto(out *SubjectListSpec) {
	*out = *in
	if in.Subjects != nil {
		in, out := &in.Subjects, &out.Subjects
		*out = make([]Subject, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubjectListSpec.
func (in *SubjectListSpec) DeepCopy() *SubjectListSpec {
	if in == nil {
		return nil
	}
	out := new(SubjectListSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubjectListStatus) DeepCopyInto(out *SubjectListStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubjectListStatus.
func (in *SubjectListStatus) DeepCopy() *SubjectListStatus {
	if in == nil {
		return nil
	}
	out := new(SubjectListStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TenantRegistration) DeepCopyInto(out *TenantRegistration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TenantRegistration.
func (in *TenantRegistration) DeepCopy() *TenantRegistration {
	if in == nil {
		return nil
	}
	out := new(TenantRegistration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TenantRegistration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TenantRegistrationList) DeepCopyInto(out *TenantRegistrationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TenantRegistration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TenantRegistrationList.
func (in *TenantRegistrationList) DeepCopy() *TenantRegistrationList {
	if in == nil {
		return nil
	}
	out := new(TenantRegistrationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TenantRegistrationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TenantRegistrationSpec) DeepCopyInto(out *TenantRegistrationSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TenantRegistrationSpec.
func (in *TenantRegistrationSpec) DeepCopy() *TenantRegistrationSpec {
	if in == nil {
		return nil
	}
	out := new(TenantRegistrationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TenantRegistrationStatus) DeepCopyInto(out *TenantRegistrationStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TenantRegistrationStatus.
func (in *TenantRegistrationStatus) DeepCopy() *TenantRegistrationStatus {
	if in == nil {
		return nil
	}
	out := new(TenantRegistrationStatus)
	in.DeepCopyInto(out)
	return out
}
