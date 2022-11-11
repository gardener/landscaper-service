// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

const (
	// LandscaperServiceFinalizer is the finalizer used for landscaper-service objects.
	LandscaperServiceFinalizer = "finalizer.landscaper-service.gardener.cloud"

	ShootTenantIDLabel          = "shoot.landscaper-service.gardener.cloud/tenantId"
	ShootInstanceNameLabel      = "shoot.landscaper-service.gardener.cloud/instanceName"
	ShootInstanceNamespaceLabel = "shoot.landscaper-service.gardener.cloud/instanceNamespace"
	ShootInstanceIDLabel        = "shoot.landscaper-service.gardener.cloud/instanceId"
)
