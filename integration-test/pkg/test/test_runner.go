// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"fmt"

	"github.com/gardener/landscaper-service/test/integration/pkg/util"
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// SharedTestObjects holds objects that are shared between tests.
type SharedTestObjects struct {
	Installations            map[string]*lsv1alpha1.Installation
	ServiceTargetConfigs     map[string]*lssv1alpha1.ServiceTargetConfig
	LandscaperDeployments    map[string]*lssv1alpha1.LandscaperDeployment
	HostingClusterNamespaces []string
}

type ClusterClients struct {
	TestCluster    client.Client
	HostingCluster client.Client
}

type ClusterTargets struct {
	HostingCluster *lsv1alpha1.Target
	LaasCluster    *lsv1alpha1.Target
}

func NewClusterClients(config *TestConfig) (*ClusterClients, error) {
	testClusterCfg, err := clientcmd.BuildConfigFromFlags("", config.TestClusterKubeconfig)
	if err != nil {
		return nil, err
	}

	testClusterClient, err := client.New(testClusterCfg, client.Options{
		Scheme: Scheme(),
	})
	if err != nil {
		return nil, err
	}

	hostingClusterCfg, err := clientcmd.BuildConfigFromFlags("", config.HostingClusterKubeconfig)
	if err != nil {
		return nil, err
	}

	hostingClusterClient, err := client.New(hostingClusterCfg, client.Options{
		Scheme: Scheme(),
	})
	if err != nil {
		return nil, err
	}

	return &ClusterClients{
		TestCluster:    testClusterClient,
		HostingCluster: hostingClusterClient,
	}, nil
}

func NewClusterTargets(ctx context.Context, kclient client.Client, config *TestConfig) (*ClusterTargets, error) {
	hostingClusterTarget, err := util.BuildKubernetesClusterTarget(ctx, kclient, config.TestClusterKubeconfig, "test-target", config.LaasNamespace)
	if err != nil {
		return nil, fmt.Errorf("cannot build hosting-cluster target: %w", err)
	}

	laasClusterTarget, err := util.BuildKubernetesClusterTarget(ctx, kclient, config.HostingClusterKubeconfig, "hosting-target", config.LaasNamespace)
	if err != nil {
		return nil, fmt.Errorf("cannot build hosting-cluster target: %w", err)
	}

	return &ClusterTargets{
		HostingCluster: hostingClusterTarget,
		LaasCluster:    laasClusterTarget,
	}, nil
}

// A TestRunner implements an integration test.
type TestRunner interface {
	Init(ctx context.Context, log logr.Logger, config *TestConfig, clusterEndpoints *ClusterClients, clusterTargets *ClusterTargets, testObjects *SharedTestObjects)
	Name() string
	Description() string
	String() string
	Run() error
}
