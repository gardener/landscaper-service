// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/gardener/landscaper/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LandscaperServiceConfiguration is the configuration for the landscaper service controller
type LandscaperServiceConfiguration struct {
	metav1.TypeMeta

	// Metrics allows to configure how metrics are exposed
	//+optional
	Metrics *MetricsConfiguration `json:"metrics,omitempty"`

	// CrdManagement configures whether the landscaper controller should deploy the CRDs it needs into the cluster
	// +optional
	CrdManagement CrdManagementConfiguration `json:"crdManagement,omitempty"`

	// LandscaperServiceComponent configures the landscaper component that is used by the landscaper service controller.
	LandscaperServiceComponent LandscaperServiceComponentConfiguration `json:"landscaperServiceComponent"`
}

// MetricsConfiguration allows to configure how metrics are exposed
type MetricsConfiguration struct {
	// Port specifies the port on which metrics are published
	Port int32 `json:"port"`
}

// CrdManagementConfiguration contains the configuration of the CRD management
type CrdManagementConfiguration struct {
	// DeployCustomResourceDefinitions specifies if CRDs should be deployed
	DeployCustomResourceDefinitions *bool `json:"deployCrd"`

	// ForceUpdate specifies whether existing CRDs should be updated
	// +optional
	ForceUpdate *bool `json:"forceUpdate,omitempty"`
}

// LandscaperServiceComponentConfiguration contains the configuration for the landscaper service component.
type LandscaperServiceComponentConfiguration struct {
	// Name is the component name
	Name string `json:"name"`

	// Version is the component version
	Version string `json:"version"`

	// RepositoryContext specifies the repository context for accessing the landscaper service component.
	RepositoryContext v1alpha1.AnyJSON `json:"repositoryContext"`

	// RegistryPullSecrets can be used to specify secrets that are needed to access the repository context.
	// +optional
	RegistryPullSecrets []corev1.SecretReference `json:"registryPullSecrets,omitempty"`
}
