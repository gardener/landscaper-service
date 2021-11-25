// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

const (
	// LandscaperServiceFinalizer is the finalizer used for landscaper-service objects.
	LandscaperServiceFinalizer = "finalizer.landscaper-service.gardener.cloud"
	// LandscaperServiceComponentName is the default component name of the landscaper-service component.
	LandscaperServiceComponentName = "github.com/gardener/landscaper/landscaper-service"
	// LandscaperServiceDefaultContext is the default context name.
	LandscaperServiceDefaultContext = "default"
)
