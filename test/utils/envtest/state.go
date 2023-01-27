// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package envtest

import (
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// State holds the information used within a single test.
type State struct {
	// Namespace is the kubernetes namespace where objects are located for testing.
	Namespace string

	// Deployments contains all landscaper deployments in this test environment.
	Deployments map[string]*lssv1alpha1.LandscaperDeployment
	// Instances contains all instances in this test environment.
	Instances map[string]*lssv1alpha1.Instance
	// Configs contains all service target configs in this test environment.
	Configs map[string]*lssv1alpha1.ServiceTargetConfig
	// Secrets contains all secrets in this test environment.
	Secrets map[string]*corev1.Secret
	// ConfigMaps contains all config maps in this test environment.
	ConfigMaps map[string]*corev1.ConfigMap
	// Installations contains all installations in this test environment.
	Installations map[string]*lsv1alpha1.Installation
	// Targets contains all targets in this test environment.
	Targets map[string]*lsv1alpha1.Target
	// Contexts contains all contexts in this test environment
	Contexts map[string]*lsv1alpha1.Context
	// AvailabilityCollections contains all availabilityCollections in this test environment
	AvailabilityCollections map[string]*lssv1alpha1.AvailabilityCollection
	// LsHealthChecks contains all LsHealthCheck in this test environment
	LsHealthChecks map[string]*lsv1alpha1.LsHealthCheck
}

// NewState creates a new state.
func NewState(namespace string) *State {
	return &State{
		Namespace:               namespace,
		Deployments:             make(map[string]*lssv1alpha1.LandscaperDeployment),
		Instances:               make(map[string]*lssv1alpha1.Instance),
		Configs:                 make(map[string]*lssv1alpha1.ServiceTargetConfig),
		Secrets:                 make(map[string]*corev1.Secret),
		ConfigMaps:              make(map[string]*corev1.ConfigMap),
		Installations:           make(map[string]*lsv1alpha1.Installation),
		Targets:                 make(map[string]*lsv1alpha1.Target),
		Contexts:                make(map[string]*lsv1alpha1.Context),
		AvailabilityCollections: make(map[string]*lssv1alpha1.AvailabilityCollection),
		LsHealthChecks:          make(map[string]*lsv1alpha1.LsHealthCheck),
	}
}

// GetDeployment retrieves a landscaper deployment by the given name.
func (s *State) GetDeployment(name string) *lssv1alpha1.LandscaperDeployment {
	return s.Deployments[s.Namespace+"/"+name]
}

// GetInstance retrieves an instance by the given name.
func (s *State) GetInstance(name string) *lssv1alpha1.Instance {
	return s.Instances[s.Namespace+"/"+name]
}

// GetConfig retrieves a landscaper target config by the given name.
func (s *State) GetConfig(name string) *lssv1alpha1.ServiceTargetConfig {
	return s.Configs[s.Namespace+"/"+name]
}

// GetSecret retrieves a secret by the given name.
func (s *State) GetSecret(name string) *corev1.Secret {
	return s.Secrets[s.Namespace+"/"+name]
}

// GetConfigMap retrieves a config map by the given name.
func (s *State) GetConfigMap(name string) *corev1.ConfigMap {
	return s.ConfigMaps[s.Namespace+"/"+name]
}

// GetInstallation retrieves an installation by the given name.
func (s *State) GetInstallation(name string) *lsv1alpha1.Installation {
	return s.Installations[s.Namespace+"/"+name]
}

// GetTarget retrieves a target by the given name.
func (s *State) GetTarget(name string) *lsv1alpha1.Target {
	return s.Targets[s.Namespace+"/"+name]
}

// GetContext retrieves a context by the given name.
func (s *State) GetContext(name string) *lsv1alpha1.Context {
	return s.Contexts[s.Namespace+"/"+name]
}

// GetAvailabilityCollectionretrieves a AvailabilityCollection by the given name
func (s *State) GetAvailabilityCollection(name string) *lssv1alpha1.AvailabilityCollection {
	return s.AvailabilityCollections[s.Namespace+"/"+name]
}

// GetLsHealthCheck a LsHealthCheck by the given name
func (s *State) GetLsHealthCheck(name string) *lsv1alpha1.LsHealthCheck {
	return s.LsHealthChecks[s.Namespace+"/"+name]
}

// GetLsHealthCheckInNamespace a LsHealthCheck by the given name in the given namespace
func (s *State) GetLsHealthCheckInNamespace(name string, namespace string) *lsv1alpha1.LsHealthCheck {
	return s.LsHealthChecks[namespace+"/"+name]
}

// AddObject adds a client.Object to the state.
func (s *State) AddObject(object client.Object) {
	switch o := object.(type) {
	case *lssv1alpha1.LandscaperDeployment:
		s.Deployments[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lssv1alpha1.Instance:
		s.Instances[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lssv1alpha1.ServiceTargetConfig:
		s.Configs[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *corev1.Secret:
		if o.Data == nil {
			o.Data = make(map[string][]byte)
		}
		for key, data := range o.StringData {
			o.Data[key] = []byte(data)
		}
		s.Secrets[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *corev1.ConfigMap:
		s.ConfigMaps[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lsv1alpha1.Installation:
		s.Installations[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lsv1alpha1.Target:
		s.Targets[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lsv1alpha1.Context:
		s.Contexts[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lssv1alpha1.AvailabilityCollection:
		s.AvailabilityCollections[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lsv1alpha1.LsHealthCheck:
		s.LsHealthChecks[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	}
}
