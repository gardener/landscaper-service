// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/types"
)

// ObjectReference is the reference to a kubernetes object.
type ObjectReference struct {
	// Name is the name of the kubernetes object.
	Name string `json:"name"`

	// Namespace is the namespace of kubernetes object.
	// +optional
	Namespace string `json:"namespace"`
}

// NamespacedName returns the namespaced name for the object reference.
func (r *ObjectReference) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Name:      r.Name,
		Namespace: r.Namespace,
	}
}

// Error holds information about an error that occurred.
type Error struct {
	// Operation describes the operator where the error occurred.
	Operation string `json:"operation"`

	// A human-readable message indicating details about the transition.
	Message string `json:"message"`
}
