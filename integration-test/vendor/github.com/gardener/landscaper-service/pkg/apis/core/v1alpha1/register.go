// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/gardener/landscaper-service/pkg/apis/core"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: core.GroupName, Version: "v1alpha1"}

// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder is a new Schema Builder which registers our API.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes, addDefaultsFuncs)
	// AddToScheme is a reference to the Schema Builder's AddToScheme function.
	AddToScheme = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Schema.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&LandscaperDeployment{},
		&LandscaperDeploymentList{},
		&Instance{},
		&InstanceList{},
		&ServiceTargetConfig{},
		&ServiceTargetConfigList{},
		&AvailabilityCollection{},
		&AvailabilityCollectionList{},
		&NamespaceRegistration{},
		&NamespaceRegistrationList{},
		&SubjectList{},
		&SubjectListList{},
		&TargetScheduling{},
		&TargetSchedulingList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
