// SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TargetSchedulingList contains a list of Scheduling
type TargetSchedulingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TargetScheduling `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TargetScheduling defines the rules according to which a LandscaperDeployment is assigned a ServiceTargetConfig.
// +kubebuilder:resource:singular="targetscheduling",path="targetschedulings",shortName="ts",scope="Namespaced"
// +kubebuilder:storageversion
type TargetScheduling struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the Scheduling
	Spec TargetSchedulingSpec `json:"spec"`
}

type TargetSchedulingSpec struct {
	Rules []SchedulingRule `json:"rules,omitempty"`
}

type SchedulingRule struct {

	// The Priority of this SchedulingRule.
	// SchedulingRules with a higher priority number will be preferred over SchedulingRules with a lower priority number.
	Priority int64 `json:"priority,omitempty"`

	ServiceTargetConfigs []ObjectReference `json:"serviceTargetConfigs,omitempty"`

	Selector []Selector `json:"selector,omitempty"`
}

type Selector struct {

	// +optional
	MatchTenant *TenantSelector `json:"matchTenant,omitempty"`

	// +optional
	MatchLabel *LabelSelector `json:"matchLabel,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +optional
	Or []Selector `json:"or,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +optional
	And []Selector `json:"and,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +optional
	Not *Selector `json:"not,omitempty"`
}

type TenantSelector struct {
	ID string `json:"id,omitempty"`
}

type LabelSelector struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
