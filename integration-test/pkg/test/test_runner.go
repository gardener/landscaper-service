// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// SharedTestObjects holds objects that are shared between tests.
type SharedTestObjects struct {
	Installations            map[string]*lsv1alpha1.Installation
	ServiceTargetConfigs     map[string]*lssv1alpha1.ServiceTargetConfig
	LandscaperDeployments    map[string]*lssv1alpha1.LandscaperDeployment
	HostingClusterNamespaces []string
}

// A TestRunner implements an integration test.
type TestRunner interface {
	Init(ctx context.Context, log logr.Logger, kclient client.Client, config *TestConfig, target *lsv1alpha1.Target, testObjects *SharedTestObjects)
	Name() string
	Run() error
}
