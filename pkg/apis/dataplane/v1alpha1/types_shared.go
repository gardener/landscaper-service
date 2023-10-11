// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Error holds information about an error that occurred.
type Error struct {
	// Operation describes the operator where the error occurred.
	Operation string `json:"operation"`

	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// Last time the condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`

	// The reason for the condition's last transition.
	Reason string `json:"reason"`

	// A human-readable message indicating details about the transition.
	Message string `json:"message"`
}
