// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// IsEmpty checks whether this reference has an empty name or empty namespace.
func (r *ObjectReference) IsEmpty() bool {
	return len(r.Name) == 0 || len(r.Namespace) == 0
}

// Equals test whether this object reference equals the given object reference.
func (r *ObjectReference) Equals(other *ObjectReference) bool {
	return r.Name == other.Name && r.Namespace == other.Namespace
}

// IsObject tests whether this object reference references the given object.
func (r *ObjectReference) IsObject(o metav1.Object) bool {
	return r.Name == o.GetName() && r.Namespace == o.GetNamespace()
}

// SecretReference is a reference to data in a secret.
type SecretReference struct {
	ObjectReference `json:",inline"`

	// Key is the name of the key in the secret that holds the data.
	// +optional
	Key string `json:"key"`
}

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

// LandscaperConfiguration contains the configuration for a landscaper service deployment.
type LandscaperConfiguration struct {
	// Deployers is the list of deployers that are getting installed alongside with this Instance.
	Deployers []string `json:"deployers"`
}

// LandscaperServiceComponent defines the landscaper service component that is being used.
type LandscaperServiceComponent struct {
	// Name defines the component name of the landscaper service component.
	Name string `json:"name"`

	// Version defines the version of the landscaper service component.
	Version string `json:"version"`
}

// OIDCConfig defines the OIDC configuration
type OIDCConfig struct {
	ClientID      string `json:"clientID,omitempty"`
	IssuerURL     string `json:"issuerURL,omitempty"`
	UsernameClaim string `json:"usernameClaim,omitempty"`
	GroupsClaim   string `json:"groupsClaim,omitempty"`
}
