// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

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
	// +optional
	Landscaper *Landscaper `json:"landscaper,omitempty"`
	// Resources configures the resources of the "central" landscaper pod, i.e. the pod responsible for crds creation,
	// deployer management, context controller.
	// +optional
	Resources *Resources `json:"resources,omitempty"`
	// ResourcesMain configures the resources of the "main" landscaper pods, i.e. the pods of installation and execution controller.
	// +optional
	ResourcesMain *Resources `json:"resourcesMain,omitempty"`
	// HPAMain configures the horizontal pod autoscaling of the "main" landscaper pods, i.e. the pods of installation and execution controller.
	// +optional
	HPAMain *HPA `json:"hpaMain,omitempty"`
	// Deployers is the list of deployers that are getting installed alongside with this Instance.
	Deployers []string `json:"deployers"`
	// DeployersConfig specifies the configuration for the landscaper standard deployers.
	// +optional
	DeployersConfig map[string]*DeployerConfig `json:"deployersConfig,omitempty"`
}

type Landscaper struct {
	Controllers        *Controllers        `json:"controllers,omitempty"`
	K8SClientSettings  *K8SClientSettings  `json:"k8sClientSettings,omitempty"`
	DeployItemTimeouts *DeployItemTimeouts `json:"deployItemTimeouts,omitempty"`
}

// Controllers specifies the config for the "main" landscaper controllers, i.e. the installation and execution controller.
type Controllers struct {
	Installations *Controller `json:"installations,omitempty"`
	Executions    *Controller `json:"executions,omitempty"`
}

// Controller specifies the config for a landscaper controller.
type Controller struct {
	Workers int32 `json:"workers,omitempty"`
}

// K8SClientSettings specifies the settings for the k8s clients which landscaper uses to access host and resource cluster.
type K8SClientSettings struct {
	HostClient     *K8SClientLimits `json:"hostClient,omitempty"`
	ResourceClient *K8SClientLimits `json:"resourceClient,omitempty"`
}

// K8SClientLimits specifies the settings for a k8s client.
type K8SClientLimits struct {
	Burst int32 `json:"burst,omitempty"`
	QPS   int32 `json:"qps,omitempty"`
}

// DeployItemTimeouts configures the timeout controller.
type DeployItemTimeouts struct {
	Pickup             string `json:"pickup,omitempty"`
	ProgressingDefault string `json:"progressingDefault,omitempty"`
}

// DeployerConfig configures a deployer.
type DeployerConfig struct {
	Deployer  *Deployer  `json:"deployer,omitempty"`
	Resources *Resources `json:"resources,omitempty"`
	HPA       *HPA       `json:"hpa,omitempty"`
}

type Deployer struct {
	Controller        *Controller        `json:"controller,omitempty"`
	K8SClientSettings *K8SClientSettings `json:"k8sClientSettings,omitempty"`
}

// Resources configures the resources of pods (requested cpu and memory)
type Resources struct {
	Requests ResourceRequests `json:"requests,omitempty"`
}

type ResourceRequests struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// HPA configures the horizontal pod autoscaling of pods.
type HPA struct {
	MaxReplicas              int32 `json:"maxReplicas,omitempty"`
	AverageMemoryUtilization int32 `json:"averageMemoryUtilization,omitempty"`
	AverageCpuUtilization    int32 `json:"averageCpuUtilization,omitempty"`
}

// LandscaperServiceComponent defines the landscaper service component that is being used.
type LandscaperServiceComponent struct {
	// Name defines the component name of the landscaper service component.
	Name string `json:"name"`

	// Version defines the version of the landscaper service component.
	Version string `json:"version"`
}

type DataPlane struct {
	SecretRef  SecretReference `json:"secretRef"`
	Kubeconfig string          `json:"kubeconfig"`
}
