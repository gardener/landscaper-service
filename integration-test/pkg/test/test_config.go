// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	k8serrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

var (
	scheme = runtime.NewScheme()
)

const (
	// LaasComponentDefault is the default Landscaper As A Service component name
	LaasComponentDefault = "github.com/gardener/landscaper-service"
	// RepoRootDir is the laas repository root directory
	RepoRootDir = "./.."
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(lsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(lssv1alpha1.AddToScheme(scheme))
}

// TestConfig contains all the configured flags of the integration test.
type TestConfig struct {
	Kubeconfig          string
	RegistryPullSecrets string
	MaxRetries          int
	SleepTime           time.Duration
	LaasComponent       string
	LaasVersion         string
	LaasRepository      string
	LandscaperNamespace string
	LaasNamespace       string
	TestNamespace       string
	LandscaperVersion   string
	ProviderType        string
}

// ParseConfig parses the TestConfig from the command line arguments.
func ParseConfig() *TestConfig {
	var (
		kubeconfig,
		registryPullSecrets,
		laasComponent, laasVersion, LaasRepository,
		landscaperNamespace, laasNamespace, testNamespace,
		providerType string
		maxRetries int
	)

	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to the kubeconfig of the cluster")
	flag.StringVar(&registryPullSecrets, "registry-secrets", "", "registry pull secrets")
	flag.IntVar(&maxRetries, "max-retries", 10, "max retries (every 10s) for all waiting operations")
	flag.StringVar(&laasVersion, "laas-version", "", "landscaper as a service version")
	flag.StringVar(&LaasRepository, "laas-repository", "", "landscaper as a service repository url")
	flag.StringVar(&laasComponent, "laas-component", LaasComponentDefault, "landscaper as a service component")
	flag.StringVar(&landscaperNamespace, "landscaper-namespace", "ls-system", "name of the landscaper namespace")
	flag.StringVar(&laasNamespace, "laas-namespace", "laas-system", "name of the landscaper as a service namespace")
	flag.StringVar(&testNamespace, "test-namespace", "laas-test", "name of the landscaper as a service integration test namespace")
	flag.StringVar(&providerType, "provider-type", "gcp", "the cloud provider of the provided kubeconfig")
	flag.Parse()

	return &TestConfig{
		Kubeconfig:          kubeconfig,
		RegistryPullSecrets: registryPullSecrets,
		MaxRetries:          maxRetries,
		SleepTime:           10 * time.Second,
		LaasComponent:       laasComponent,
		LaasVersion:         laasVersion,
		LaasRepository:      LaasRepository,
		LandscaperNamespace: landscaperNamespace,
		LaasNamespace:       laasNamespace,
		TestNamespace:       testNamespace,
		ProviderType:        providerType,
	}
}

// VerifyConfig validates the given TestConfig.
func VerifyConfig(config *TestConfig) error {
	errorList := make([]error, 0)

	if len(config.Kubeconfig) == 0 {
		errorList = append(errorList, fmt.Errorf("flag \"kubeconfig\" may not be empty"))
	}
	if len(config.RegistryPullSecrets) == 0 {
		errorList = append(errorList, fmt.Errorf("flag \"registry-secrets\" may not be empty"))
	}
	if len(config.LaasVersion) == 0 {
		errorList = append(errorList, fmt.Errorf("flag \"laas-version\" may not be empty"))
	}
	if len(config.LaasRepository) == 0 {
		errorList = append(errorList, fmt.Errorf("flag \"laas-repository\" may not be empty"))
	}
	return k8serrors.NewAggregate(errorList)
}

// Scheme returns the integration test scheme.
func Scheme() *runtime.Scheme {
	return scheme
}
