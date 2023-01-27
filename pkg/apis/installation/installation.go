// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package installation

import (
	"encoding/json"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
)

// This file contains the API of the landscaper-service component "installation-blueprint" resource.
// github.com/gardener/landscaper-service/landscaper-instance

const (
	// UseInMemoryOverlayDefault is default value for the landscaper cache in memory overlay configuration.
	UseInMemoryOverlayDefault = false
	// AllowPlainHttpRegistriesDefault is the default value for the landscaper registry config. allow plain http configuration.
	AllowPlainHttpRegistriesDefault = false
	// InsecureSkipVerifyDefault is the default value for the landscaper registry config. insecure skip verify configuration.
	InsecureSkipVerifyDefault = false

	// ReplicasDefault is the default number of replicas for the landscaper controller deployment.
	ReplicasDefault = 1
	// VerbosityDefault is the default verbose level for the landscaper controller deployment.
	VerbosityDefault = logging.INFO
	// WebhooksServicePortDefault is the default service port for the landscaper webhooks server deployment.
	WebhooksServicePortDefault = 9443

	// HostingClusterNamespaceImportName is the import name for the hosting cluster namespace.
	HostingClusterNamespaceImportName = "hostingClusterNamespace"
	// TargetClusterNamespaceImportName is the import for the target cluster namespace.
	TargetClusterNamespaceImportName = "targetClusterNamespace"
	// RegistryConfigImportName is the import for the registry configuration.
	RegistryConfigImportName = "registryConfig"
	// LandscaperConfigImportName is the import for the landscaper configuration.
	LandscaperConfigImportName = "landscaperConfig"
	// ShootNameImportName is the import for the shoot name.
	ShootNameImportName = "shootName"
	// ShootNamespaceImportName is the import for the shoot namespace.
	ShootNamespaceImportName = "shootNamespace"
	// ShootSecretBindingImportName is the import for the name of the shoot secret binding.
	ShootSecretBindingImportName = "shootSecretBindingName"
	// ShootLabelsImportName is the import for the shoot labels.
	ShootLabelsImportName = "shootLabels"
	// ShootConfigImportName is the shoot configuration import.
	ShootConfigImportName = "shootConfig"
	//WebhooksHostNameImportName is the import for the webhooks host name.
	WebhooksHostNameImportName = "webhooksHostName"

	// TargetClusterNamespace is the target cluster namespace used for landscaper internals.
	TargetClusterNamespace = "ls-system"

	// ClusterEndpointExportName is the name of the cluster endpoint export.
	ClusterEndpointExportName = "landscaperClusterEndpoint"
	// UserKubeconfigExportName is the name of the user kubeconfig export.
	UserKubeconfigExportName = "landscaperUserKubeconfig"
	// AdminKubeconfigExportName is the name of the admin kubeconfig export.
	AdminKubeconfigExportName = "landscaperAdminKubeconfig"
)

// CacheConfig specifies the landscaper registry cache configuration.
type CacheConfig struct {
	// UseInMemoryOverly - see github.com/gardener/landscaper/apis/config OCICacheConfiguration.UseInMemoryOverly
	UseInMemoryOverly bool `json:"useInMemoryOverly"`
}

// RegistryConfig specifies the landscaper registry configuration.
type RegistryConfig struct {
	// Cache is the cache configuration.
	Cache CacheConfig `json:"cache"`
	// AllowPlainHttpRegistries - see github.com/gardener/landscaper/apis/config OCIConfiguration.AllowPlainHttp
	AllowPlainHttpRegistries bool `json:"allowPlainHttpRegistries"`
	// InsecureSkipVerify - see github.com/gardener/landscaper/apis/config OCIConfiguration.InsecureSkipVerify
	InsecureSkipVerify bool `json:"insecureSkipVerify"`
}

// NewRegistryConfig creates a new registry configuration initialized with default values.
func NewRegistryConfig() *RegistryConfig {
	r := &RegistryConfig{
		Cache: CacheConfig{
			UseInMemoryOverly: UseInMemoryOverlayDefault,
		},
		AllowPlainHttpRegistries: AllowPlainHttpRegistriesDefault,
		InsecureSkipVerify:       InsecureSkipVerifyDefault,
	}
	return r
}

// ToAnyJSON marshals this registry configuration to an AnyJSON object.
func (r *RegistryConfig) ToAnyJSON() (*lsv1alpha1.AnyJSON, error) {
	raw, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	anyJSON := lsv1alpha1.NewAnyJSON(raw)
	return &anyJSON, err
}

// Landscaper specifies the landscaper controller configuration.
type Landscaper struct {
	// Verbosity defines the logging verbosity level.
	Verbosity string `json:"verbosity,omitempty"`
	// Replicas defines the number of replicas for the landscaper controller deployment.
	Replicas int `json:"replicas,omitempty"`
}

// Webhooks specifies the landscaper webhooks server configuration.
type Webhooks struct {
	// ServicePort specifies the landscaper webhooks service port.
	ServicePort int `json:"servicePort,omitempty"`
	// Replicas defines the number of replicas for the landscaper webhooks server deployment.
	Replicas int `json:"replicas,omitempty"`
}

// LandscaperConfig specifies the landscaper deployment configuration for the API of the "landscaper-as-a-service" component.
type LandscaperConfig struct {
	// Landscaper specifies the landscaper controller configuration.
	Landscaper Landscaper `json:"landscaper"`
	// Webhooks specifies the landscaper webhooks server configuration.
	Webhooks Webhooks `json:"webhooksServer"`
	// Deployers specifies the list of landscaper standard deployers that are getting installed.
	Deployers []string `json:"deployers"`
	// DeployersConfig specifies the configuration for the landscaper standard deployers.
	DeployersConfig lsv1alpha1.AnyJSON `json:"deployersConfig,omitempty"`
}

// NewLandscaperConfig creates a new landscaper configuration initialized with default values.
func NewLandscaperConfig() *LandscaperConfig {
	c := &LandscaperConfig{
		Landscaper: Landscaper{
			Verbosity: VerbosityDefault.String(),
			Replicas:  ReplicasDefault,
		},
		Webhooks: Webhooks{
			ServicePort: WebhooksServicePortDefault,
			Replicas:    ReplicasDefault,
		},
		Deployers:       make([]string, 0),
		DeployersConfig: lsv1alpha1.NewAnyJSON([]byte("{}")),
	}
	return c
}

// ToAnyJSON marshals this landscaper configuration to an AnyJSON object.
func (l *LandscaperConfig) ToAnyJSON() (*lsv1alpha1.AnyJSON, error) {
	raw, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	anyJSON := lsv1alpha1.NewAnyJSON(raw)
	return &anyJSON, err
}
