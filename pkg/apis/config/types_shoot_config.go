// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package config

// ShootConfiguration holds the configuration of a gardener shoot cluster.
type ShootConfiguration struct {
	// Provider is the shoot provider configuration.
	Provider ShootProviderConfiguration `json:"provider"`
	// Workers is the shot workers configuration.
	Workers ShootWorkersConfiguration `json:"workers"`

	// Region specifies the region the shoot cluster shall be created in.
	Region string `json:"region"`

	// Kubernetes is the shoot kubernetes configuration.
	Kubernetes ShootKubernetesConfig `json:"kubernetes"`
	// Maintenance is the shoot maintenance configuration.
	Maintenance ShootMaintenanceConfig `json:"maintenance"`
}

// ShootProviderConfiguration is the shoot provider configuration.
type ShootProviderConfiguration struct {
	// Type is the cloud provider type.
	Type string `json:"type"`
	// Zone is the cloud provider specific zone in which the shoot is being created.
	Zone string `json:"zone"`
}

// ShootWorkersConfiguration is the configuration for the shoot worker nodes.
type ShootWorkersConfiguration struct {
	// Machine specifies the machine type used for worker nodes.
	Machine ShootMachineConfiguration `json:"machine"`
	// Volume specifies the volume configuration for the worker nodes.
	Volume ShootWorkerVolumeConfiguration `json:"volume"`

	// Minimum is the minimum amount of worker nodes available.
	Minimum *int32 `json:"minimum"`
	// Maximum is the maximum amount of worker nodes available.
	Maximum *int32 `json:"maximum"`
	// MaxSurge is the amount of worker nodes created during an update.
	MaxSurge *int32 `json:"maxSurge"`
	// MaxUnavailable is the maximum number of nodes unavailable during an update.
	MaxUnavailable *int32 `json:"maxUnavailable"`
}

// ShootMachineConfiguration is the machine specification used for worker nodes.
type ShootMachineConfiguration struct {
	// Type is the cloud provider specific virtual machine type.
	Type string `json:"type"`
	// Image is the virtual machine image used for the worker nodes.
	Image ShootMachineImage `json:"image"`
}

// ShootMachineImage specifies a virtual machine image for worker nodes.
type ShootMachineImage struct {
	// Name of the image.
	Name string `json:"name"`
	// Version of the image.
	Version string `json:"version"`
}

// ShootWorkerVolumeConfiguration is the volume configuration for worker nodes.
type ShootWorkerVolumeConfiguration struct {
	// The cloud provider specific volume type.
	Type string `json:"type"`
	// The size of the volume.
	Size string `json:"size"`
}

// ShootKubernetesConfig is the shoot kubernetes configuration.
type ShootKubernetesConfig struct {
	// The kubernetes version to use.
	Version string `json:"version"`
}

// ShootMaintenanceConfig specifies the maintenance handling for the shoot cluster.
type ShootMaintenanceConfig struct {
	// AutoUpdate specifies which components of the shoot clusters shall be updated automatically
	// during a maintenance time window.
	AutoUpdate ShootAutoUpdateConfig `json:"autoUpdate"`
	// TimeWindow is the time window during which auto updates are performed.
	TimeWindow ShootMaintenanceTimeWindow `json:"timeWindow"`
}

// ShootAutoUpdateConfig specifies which components of the shoot clusters shall be updated automatically
// during a maintenance time window.
type ShootAutoUpdateConfig struct {
	// KubernetesVersion, if set to true, auto updates kubernetes patch versions.
	KubernetesVersion *bool `json:"kubernetesVersion"`
	// MachineImageVersion, if set to true, auto updates machine image versions.
	MachineImageVersion *bool `json:"machineImageVersion"`
}

// ShootMaintenanceTimeWindow is the time window during which auto updates are performed.
type ShootMaintenanceTimeWindow struct {
	// Begin specifies the beginning of the auto update time window.
	Begin string `json:"begin"`
	// Begin specifies the ending of the auto update time window.
	End string `json:"end"`
}
