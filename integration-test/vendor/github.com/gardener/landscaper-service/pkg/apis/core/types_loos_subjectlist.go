// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceRegistrationList contains a list of NamespaceRegistration
type LosSubjectListList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LosSubjectList `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type LosSubjectList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the SubjectList.
	Spec LosSubjectListSpec `json:"spec"`

	// Status contains the status for the SubjectList.
	// +optional
	Status LosSubjectListStatus `json:"status"`
}

// SubjectListStatus contains the status for the SubjectList.
type LosSubjectListStatus struct {
	Phase              string `json:"phase"`
	SyncedGeneration   int64  `json:"syncedGeneration"`
	ObservedGeneration int64  `json:"observedGeneration"`
}

// SubjectListSpec contains the specification for the SubjectList.
type LosSubjectListSpec struct {
	//Admins contains references to the object or user identities the admin role binding applies to.
	Admins []LosSubject `json:"admins"`
	//Members contains references to the object or user identities the member role binding applies to.
	Members []LosSubject `json:"members"`
	//Viewer contains references to the object or user identities the viewer role binding applies to.
	Viewer []LosSubject `json:"viewer"`
}

// LosSubject is a User, Group or ServiceAccount(with namespace). Similar to rbac.Subject struct but does not depend on it to prevent future k8s version from breaking this logic.
type LosSubject struct {
	// Kind of object being referenced. Values defined by this API group are "User", "Group", and "ServiceAccount".
	// If the Authorizer does not recognized the kind value, the Authorizer should report an error.
	Kind string `json:"kind"`
	// Name of the object being referenced.
	Name string `json:"name"`
}
