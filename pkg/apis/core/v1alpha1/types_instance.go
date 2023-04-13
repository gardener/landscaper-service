// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	lsschema "github.com/gardener/landscaper/apis/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// The Instance is created for each LandscaperDeployment.
// The landscaper service controller selects a suitable/available ServiceTargetConfig and creates
// an Installation.
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the Instance.
	Spec InstanceSpec `json:"spec"`

	// Status contains the status for the Instance.
	// +optional
	Status InstanceStatus `json:"status"`
}

// InstanceSpec contains the specification for an Instance.
type InstanceSpec struct {
	// TenantId is the unique identifier of the owning tenant.
	TenantId string `json:"tenantId"`

	// ID is the id of this instance
	ID string `json:"id"`

	// LandscaperConfiguration contains the configuration for the landscaper service deployment
	LandscaperConfiguration LandscaperConfiguration `json:"landscaperConfiguration"`

	// ServiceTargetConfigRef specifies the target cluster for which the installation is created.
	ServiceTargetConfigRef ObjectReference `json:"serviceTargetConfigRef"`

	// OIDCConfig describes the OIDC config of the customer resource cluster (shoot cluster)
	// +optional
	OIDCConfig *OIDCConfig `json:"oidcConfig,omitempty"`

	// AutomaticReconcile specifies the configuration on when this instance is being automatically reconciled.
	// +optional
	AutomaticReconcile *AutomaticReconcile `json:"automaticReconcile,omitempty"`

	// HighAvailabilityConfig specifies the HA configuration of the resource cluster (shoot cluster)
	// +optional
	HighAvailabilityConfig *HighAvailabilityConfig `json:"highAvailabilityConfig"`
}

// AutomaticReconcile defines the automatic reconcile configuration.
type AutomaticReconcile struct {
	// Interval specifies the interval after which the instance is being automatically reconciled.
	Interval lsv1alpha1.Duration `json:"interval"`
}

// InstanceStatus contains the status for an Instance.
type InstanceStatus struct {
	// ObservedGeneration is the most recent generation observed for this Instance.
	// It corresponds to the Instance generation, which is updated on mutation by the landscaper service controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration"`

	// LastError describes the last error that occurred.
	// +optional
	LastError *Error `json:"lastError,omitempty"`

	// LandscaperServiceComponent define the landscaper server component that is used for this instance.
	// +optional
	LandscaperServiceComponent *LandscaperServiceComponent `json:"landscaperServiceComponent"`

	// ContextRef references the landscaper context for this Instance.
	// +optional
	ContextRef *ObjectReference `json:"contextRef,omitempty"`

	// TargetRef references the Target for this Instance.
	// +optional
	TargetRef *ObjectReference `json:"targetRef,omitempty"`

	// InstallationRef references the Installation for this Instance.
	// +optional
	InstallationRef *ObjectReference `json:"installationRef,omitempty"`

	// ClusterEndpointRef contains the URL at which the landscaper cluster is accessible.
	// +optional
	ClusterEndpoint string `json:"clusterEndpoint,omitempty"`

	// UserKubeconfig contains the user kubeconfig which can be used for accessing the landscaper cluster.
	// +optional
	UserKubeconfig string `json:"userKubeconfig,omitempty"`

	// AdminKubeconfig contains the admin kubeconfig which can be used for accessing the landscaper cluster.
	// +optional
	AdminKubeconfig string `json:"adminKubeconfig,omitempty"`

	// ShootName is the name of the corresponding shoot cluster.
	// +optional
	ShootName string `json:"shootName,omitempty"`

	// ShootNamespace is the namespace in which the shoot resource is being created.
	// +optional
	ShootNamespace string `json:"shootNamespace,omitempty"`

	// AutomaticReconcileStatus contains the status of the automatic reconciliation of this instance.
	// +optional
	AutomaticReconcileStatus *AutomaticReconcileStatus `json:"automaticReconcileStatus,omitempty"`
}

// AutomaticReconcileStatus contains the automatic reconciliation status of an instance.
type AutomaticReconcileStatus struct {
	// LastReconcileTime contains the time at which the instance has been reconciled.
	LastReconcileTime metav1.Time `json:"lastReconcileTime,omitempty"`
}

var InstanceDefinition = lsschema.CustomResourceDefinition{
	Names: lsschema.CustomResourceDefinitionNames{
		Plural:   "instances",
		Singular: "instance",
		ShortNames: []string{
			"instc",
		},
		Kind: "Instance",
	},
	Scope:             lsschema.NamespaceScoped,
	Storage:           true,
	Served:            true,
	SubresourceStatus: true,
	AdditionalPrinterColumns: []lsschema.CustomResourceColumnDefinition{
		{
			Name:     "ServiceTargetConfig",
			Type:     "string",
			JSONPath: ".spec.serviceTargetConfigRef.name",
		},
		{
			Name:     "Installation",
			Type:     "string",
			JSONPath: ".status.installationRef.name",
		},
		{
			Name:     "Age",
			Type:     "date",
			JSONPath: ".metadata.creationTimestamp",
		},
	},
}
