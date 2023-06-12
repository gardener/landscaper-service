// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TenantRegistrationList contains a list of TenantRegistration
type TenantRegistrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TenantRegistration `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TenantRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the TenantRegistration.
	Spec TenantRegistrationSpec `json:"spec"`

	// Status contains the status for the TenantRegistration.
	// +optional
	Status TenantRegistrationStatus `json:"status"`
}

// TenantRegistrationStatus contains the status for the TenantRegistration.
type TenantRegistrationStatus struct {
	Synced bool `json:"synced"`
	Ready  bool `json:"ready"`
}

// TenantRegistrationSpec contains the specification for the TenantRegistration.
type TenantRegistrationSpec struct {
	Author string `json:"author"`
}
