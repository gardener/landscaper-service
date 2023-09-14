// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package installation

import (
	"encoding/json"
	"fmt"
	"strings"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"

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
	// SidecarConfigImportName is the import for the sidecar configuration.
	SidecarConfigImportName = "sidecarConfig"
	// RotationConfigImportName is the import for the rotation configuration.
	RotationConfigImportName = "rotationConfig"
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
	// AuditPolicyImportName is the import for the audit policy configuration.
	AuditPolicyImportName = "auditPolicy"
	// AuditLogServiceImportName is the import for the audit log service settings.
	AuditLogServiceImportName = "auditLogService"

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
	return toAnyJSON(r)
}

// Landscaper specifies the landscaper controller configuration.
type Landscaper struct {
	// Verbosity defines the logging verbosity level.
	Verbosity string `json:"verbosity,omitempty"`
	// Replicas defines the number of replicas for the landscaper controller deployment.
	Replicas           int                             `json:"replicas,omitempty"`
	Controllers        *lssv1alpha1.Controllers        `json:"controllers,omitempty"`
	DeployItemTimeouts *lssv1alpha1.DeployItemTimeouts `json:"deployItemTimeouts,omitempty"`
	// K8SClientSettings defines k8s client settings like burst and qps.
	K8SClientSettings *lssv1alpha1.K8SClientSettings `json:"k8sClientSettings,omitempty"`
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
	Webhooks      Webhooks               `json:"webhooksServer"`
	Resources     *lssv1alpha1.Resources `json:"resources,omitempty"`
	ResourcesMain *lssv1alpha1.Resources `json:"resourcesMain,omitempty"`
	HPAMain       *lssv1alpha1.HPA       `json:"hpaMain,omitempty"`
	// Deployers specifies the list of landscaper standard deployers that are getting installed.
	Deployers []string `json:"deployers"`
	// DeployersConfig specifies the configuration for the landscaper standard deployers.
	DeployersConfig map[string]*lssv1alpha1.DeployerConfig `json:"deployersConfig,omitempty"`
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
			Replicas:    2,
		},
		Deployers: make([]string, 0),
	}
	return c
}

// ToAnyJSON marshals this landscaper configuration to an AnyJSON object.
func (l *LandscaperConfig) ToAnyJSON() (*lsv1alpha1.AnyJSON, error) {
	return toAnyJSON(l)
}

// SidecarConfig specifies the config for the namespace registration and subject sync controller.
type SidecarConfig struct {
	// Verbosity defines the logging verbosity level.
	Verbosity string `json:"verbosity,omitempty"`
}

// NewSidecarConfig creates a new SidecarConfig.
func NewSidecarConfig() *SidecarConfig {
	c := &SidecarConfig{
		Verbosity: VerbosityDefault.String(),
	}
	return c
}

// ToAnyJSON marshals this SidecarConfig to an AnyJSON object.
func (l *SidecarConfig) ToAnyJSON() (*lsv1alpha1.AnyJSON, error) {
	return toAnyJSON(l)
}

// RotationConfig specifies the config for the rotation of credentials.
type RotationConfig struct {
	// TokenExpirationSeconds defines how long the tokens are valid	which the landscaper and sidecar controllers use
	// to access the resource cluster, e.g. for watching installations, namespace registrations etc.
	TokenExpirationSeconds int64 `json:"tokenExpirationSeconds,omitempty"`
	// AdminKubeconfigExpirationSeconds defines how long the admin kubeconfig for a resource cluster is valid.
	// The kubeconfig is used to deploy RBAC objects on the resource cluster.
	AdminKubeconfigExpirationSeconds int64 `json:"adminKubeconfigExpirationSeconds,omitempty"`
}

// NewRotationConfig creates a new RotationConfig.
func NewRotationConfig(tokenExpirationSeconds, adminKubeconfigExpirationSeconds int64) *RotationConfig {
	return &RotationConfig{
		TokenExpirationSeconds:           tokenExpirationSeconds,
		AdminKubeconfigExpirationSeconds: adminKubeconfigExpirationSeconds,
	}
}

// ToAnyJSON marshals this RotationConfig to an AnyJSON object.
func (r *RotationConfig) ToAnyJSON() (*lsv1alpha1.AnyJSON, error) {
	return toAnyJSON(r)
}

func toAnyJSON(obj any) (*lsv1alpha1.AnyJSON, error) {
	raw, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	anyJSON := lsv1alpha1.NewAnyJSON(raw)
	return &anyJSON, err
}

// GetInstallationExportDataRef returns the export data ref that is dynamically created based on the instance name.
// The data ref string is compatible to be used as a kubernetes object name.
func GetInstallationExportDataRef(instance *lssv1alpha1.Instance, exportName string) string {
	dataRef := fmt.Sprintf("%s-%s", strings.ToLower(exportName), instance.GetName())
	return dataRef
}
