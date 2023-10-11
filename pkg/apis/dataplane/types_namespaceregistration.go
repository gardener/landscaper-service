// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package dataplane

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceRegistrationList contains a list of NamespaceRegistration
type NamespaceRegistrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceRegistration `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NamespaceRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the NamespaceRegistration.
	Spec NamespaceRegistrationSpec `json:"spec"`

	// Status contains the status for the NamespaceRegistration.
	// +optional
	Status NamespaceRegistrationStatus `json:"status"`
}

type NamespaceRegistrationStatus struct {
	Phase string `json:"phase"`
	// +optional
	LastError *Error `json:"lastError,omitempty"`
}

type NamespaceRegistrationSpec struct {
}
