// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TargetShootSidecarConfiguration is the configuration for the landscaper service controller
type TargetShootSidecarConfiguration struct {
	metav1.TypeMeta

	// Metrics allows to configure how metrics are exposed
	//+optional
	Metrics *MetricsConfiguration `json:"metrics,omitempty"`

	// CrdManagement configures whether the landscaper controller should deploy the CRDs it needs into the cluster
	// +optional
	CrdManagement CrdManagementConfiguration `json:"crdManagement,omitempty"`
}
