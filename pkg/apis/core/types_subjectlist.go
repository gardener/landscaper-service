// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceRegistrationList contains a list of NamespaceRegistration
type SubjectListList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SubjectList `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SubjectList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the SubjectList.
	Spec SubjectListSpec `json:"spec"`

	// Status contains the status for the SubjectList.
	// +optional
	Status SubjectListStatus `json:"status"`
}

// SubjectListStatus contains the status for the SubjectList.
type SubjectListStatus struct {
	Phase              string `json:"phase"`
	ObservedGeneration int64  `json:"observedGeneration"`
}

// SubjectListSpec contains the specification for the SubjectList.
type SubjectListSpec struct {
	//Subject contains a reference to the object or user identities a role binding applies to.
	Subjects []Subject `json:"subjects"`
}

// Subject is a User, Group or ServiceAccount(with namespace). Similar to rbac.Subject struct but does not depend on it to prevent future k8s version from breaking this logic.
type Subject struct {
	// Kind of object being referenced. Values defined by this API group are "User", "Group", and "ServiceAccount".
	// If the Authorizer does not recognized the kind value, the Authorizer should report an error.
	Kind string `json:"kind"`
	// Name of the object being referenced.
	Name string `json:"name"`
	// Namespace of the referenced object.  If the object kind is non-namespace, such as "User" or "Group", and this value is not empty
	// the Authorizer should report an error.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}
