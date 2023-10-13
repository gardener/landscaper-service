// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AvailabilityCollectionList contains a list of AvailabilityCollection
type AvailabilityCollectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AvailabilityCollection `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AvailabilityCollection is created/updated by the AvilabilityMonitoringRegistrationController.
// It contains a list of references to Instances that should be monitored for availability
type AvailabilityCollection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the AvailabilityCollection.
	Spec AvailabilityCollectionSpec `json:"spec"`

	// Status contains the status for the AvailabilityCollection.
	// +optional
	Status AvailabilityCollectionStatus `json:"status"`
}

// AvailabilityCollectionStatus contains the status for the AvailabilityCollection.
type AvailabilityCollectionStatus struct {
	// metadata.generation observed by the HealthWatcher controller.
	// Used to distinguish between a necessary reconcile (scheduled or spec change)
	// and unnecessary reconcile (status change)
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastRun is the last time, the HealthWatcher collected all instance status.
	// +optional
	LastRun metav1.Time `json:"lastRun"`

	// LastReported is the last time, the AV Uploader uploaded all instance status. Prevents multi upload of the same status.
	// +optional
	LastReported metav1.Time `json:"lastReported"`

	// Instances collects the status for all instances specified in spec.instanceRefs
	Instances []AvailabilityInstance `json:"instances"`

	// Self collects the status the own landscaper
	Self AvailabilityInstance `json:"self"`
}

// AvailabilityInstance contains the availability status for one instance.
type AvailabilityInstance struct {
	ObjectReference `json:",inline"`
	// Status is the availability status of the instance.
	Status string `json:"status"`
	// FailedReason is the reason the status is in failed.
	FailedReason string `json:"failedReason"`

	// FailedSince contains the timestamp since the object is in failed status
	// +optional
	FailedSince *metav1.Time `json:"failedSince,omitempty"`
}

// AvailabilityCollectionSpec contains the spec for the AvailabilityCollection.
type AvailabilityCollectionSpec struct {
	// InstanceRefs specifies all instances to monitor
	InstanceRefs []ObjectReference `json:"instanceRefs"`
}
