// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"github.com/gardener/landscaper/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LandscaperServiceConfiguration is the configuration for the landscaper service controller
type LandscaperServiceConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// Metrics allows to configure how metrics are exposed
	//+optional
	Metrics *MetricsConfiguration `json:"metrics,omitempty"`

	//AvailabilityMonitoringConfiguration is the configuration for the availability monitoring of the provisioned landscaper
	AvailabilityMonitoring AvailabilityMonitoringConfiguration `json:"availabilityMonitoring"`

	// CrdManagement configures whether the landscaper controller should deploy the CRDs it needs into the cluster
	// +optional
	CrdManagement CrdManagementConfiguration `json:"crdManagement,omitempty"`

	// LandscaperServiceComponent configures the landscaper component that is used by the landscaper service controller.
	LandscaperServiceComponent LandscaperServiceComponentConfiguration `json:"landscaperServiceComponent"`

	// GardenerConfiguration is the gardener specific configuration required for shoot management.
	GardenerConfiguration GardenerConfiguration `json:"gardenerConfiguration"`

	// ShootConfiguration is the specification to the gardener shoots.
	ShootConfiguration ShootConfiguration `json:"shootConfiguration"`

	// AuditLogConfig is the audit log configuration for the created shoots.
	// +optional
	AuditLogConfig *AuditLogConfiguration `json:"auditLogConfig"`

	// Scheduling specifies a reference to a Scheduling resource,
	// which defines rules how ServiceTargetConfigs are assigned to LandscaperDeployments.
	// +optional
	Scheduling *v1alpha1.ObjectReference `json:"scheduling,omitempty"`
}

// AvailabilityMonitoringConfiguration is the configuration for the availability monitoring of the provisioned landscaper
type AvailabilityMonitoringConfiguration struct {
	//AvailabilityCollectionName is the name of the CR containing the av monitoring statuses
	AvailabilityCollectionName string `json:"availabilityCollectionName"`
	//AvailabilityCollectionNamespace is the namespace of the CR containing the av monitoring statuses
	AvailabilityCollectionNamespace string `json:"availabilityCollectionNamespace"`

	//AvailabilityServiceConfiguration configures an external AVS service
	AvailabilityServiceConfiguration *AvailabilityServiceConfiguration `json:"availabilityService"`

	//SelfLandscaperNamespace defines the namespace of the landscaper in the core cluster to be monitored
	SelfLandscaperNamespace string `json:"selfLandscaperNamespace"`

	//PeriodicCheckInterval defines, how often the HealthWatcher controller collects the landscaper health information
	PeriodicCheckInterval v1alpha1.Duration `json:"periodicCheckInterval"`
	//LSHealthCheckTimeout defines the timeout, at which
	// (1) a previously available landscaper is unavailable if no updates occurred
	// (2) a failed landscaper is reported as failed if it does not become available again
	LSHealthCheckTimeout v1alpha1.Duration `json:"lsHealthCheckTimeout"`
}

// AvailabilityServiceConfiguration configures an external AVS service
type AvailabilityServiceConfiguration struct {
	//Url is the full url to the AVS
	Url string `json:"url"`
	//ApiKey is the api key for the AVS
	ApiKey string `json:"apiKey"`
	//Timeout is the timeout for the AVS request
	Timeout string `json:"timeout"`
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

// GardenerConfiguration is the gardener specific configuration required for shoot management.
type GardenerConfiguration struct {
	// ServiceAccountKubeconfig is the reference to the secret containing the service account kubeconfig.
	ServiceAccountKubeconfig v1alpha1.SecretReference `json:"serviceAccountKubeconfig"`

	// ProjectName is the name of gardener project.
	ProjectName string `json:"projectName"`

	// ShootSecretBindingName is the secret binding which is used to allocate resources for shoots.
	ShootSecretBindingName string `json:"shootSecretBindingName"`
}

// AuditLogConfiguration specifies the shoot cluster audit log configuration.
type AuditLogConfiguration struct {
	// AuditLogService contains the audit log service configuration.
	AuditLogService AuditLogService `json:"auditLogService"`
	// AuditPolicy specifies the Kubernetes API server audit policy (audit.k8s.io Policy)
	AuditPolicy v1alpha1.ConfigMapReference `json:"auditPolicy"`
}

// AuditLogService contains the audit log service configuration.
type AuditLogService struct {
	// TenantId is the tenant id of the BTP subaccount in which the audit log service is available.
	TenantId string `json:"tenantId"`
	// Url is the url of the audit log service key credentials.
	Url string `json:"url"`
	// User is the user of the audit log service key credentials.
	User string `json:"user"`
	// Password is the password of the audit log service key credentials.
	Password string `json:"password"`
}
