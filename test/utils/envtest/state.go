// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package envtest

import (
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"

	dataplanev1alpha1 "github.com/gardener/landscaper-service/pkg/apis/dataplane/v1alpha1"
	provisioningv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/provisioning/v1alpha2"
)

// State holds the information used within a single test.
type State struct {
	// Namespace is the kubernetes namespace where objects are located for testing.
	Namespace string

	// Deployments contains all landscaper deployments in this test environment.
	Deployments map[string]*provisioningv1alpha2.LandscaperDeployment
	// Instances contains all instances in this test environment.
	Instances map[string]*provisioningv1alpha2.Instance
	// Configs contains all service target configs in this test environment.
	Configs map[string]*provisioningv1alpha2.ServiceTargetConfig
	// Secrets contains all secrets in this test environment.
	Secrets map[string]*corev1.Secret
	// ConfigMaps contains all config maps in this test environment.
	ConfigMaps map[string]*corev1.ConfigMap
	// Installations contains all installations in this test environment.
	Installations map[string]*lsv1alpha1.Installation
	// Executions contains all executions in this test environment.
	Executions map[string]*lsv1alpha1.Execution
	// DeployItems contains all DeployItems in this test environment.
	DeployItems map[string]*lsv1alpha1.DeployItem
	// TargetSync contains all targetsyncs in this test environment.
	TargetSync map[string]*lsv1alpha1.TargetSync
	// Targets contains all targets in this test environment.
	Targets map[string]*lsv1alpha1.Target
	// Contexts contains all contexts in this test environment
	Contexts map[string]*lsv1alpha1.Context
	// AvailabilityCollections contains all availabilityCollections in this test environment
	AvailabilityCollections map[string]*provisioningv1alpha2.AvailabilityCollection
	// LsHealthChecks contains all LsHealthCheck in this test environment
	LsHealthChecks map[string]*lsv1alpha1.LsHealthCheck
	// NamespaceRegistrations contains all NamespaceRegistration in this test environment
	NamespaceRegistrations map[string]*dataplanev1alpha1.NamespaceRegistration
	// SubjectLists contains all SubjectList in this test environment
	SubjectLists map[string]*dataplanev1alpha1.SubjectList
}

// NewState creates a new state.
func NewState(namespace string) *State {
	return &State{
		Namespace:               namespace,
		Deployments:             make(map[string]*provisioningv1alpha2.LandscaperDeployment),
		Instances:               make(map[string]*provisioningv1alpha2.Instance),
		Configs:                 make(map[string]*provisioningv1alpha2.ServiceTargetConfig),
		Secrets:                 make(map[string]*corev1.Secret),
		ConfigMaps:              make(map[string]*corev1.ConfigMap),
		Installations:           make(map[string]*lsv1alpha1.Installation),
		Executions:              make(map[string]*lsv1alpha1.Execution),
		DeployItems:             make(map[string]*lsv1alpha1.DeployItem),
		TargetSync:              make(map[string]*lsv1alpha1.TargetSync),
		Targets:                 make(map[string]*lsv1alpha1.Target),
		Contexts:                make(map[string]*lsv1alpha1.Context),
		AvailabilityCollections: make(map[string]*provisioningv1alpha2.AvailabilityCollection),
		LsHealthChecks:          make(map[string]*lsv1alpha1.LsHealthCheck),
		NamespaceRegistrations:  make(map[string]*dataplanev1alpha1.NamespaceRegistration),
		SubjectLists:            make(map[string]*dataplanev1alpha1.SubjectList),
	}
}

// GetDeployment retrieves a landscaper deployment by the given name.
func (s *State) GetDeployment(name string) *provisioningv1alpha2.LandscaperDeployment {
	return s.Deployments[s.Namespace+"/"+name]
}

// GetInstance retrieves an instance by the given name.
func (s *State) GetInstance(name string) *provisioningv1alpha2.Instance {
	return s.Instances[s.Namespace+"/"+name]
}

// GetConfig retrieves a landscaper target config by the given name.
func (s *State) GetConfig(name string) *provisioningv1alpha2.ServiceTargetConfig {
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

func (s *State) GetExecution(name string) *lsv1alpha1.Execution {
	return s.Executions[s.Namespace+"/"+name]
}

func (s *State) GetDeployItem(name string) *lsv1alpha1.DeployItem {
	return s.DeployItems[s.Namespace+"/"+name]
}

func (s *State) GetTargetSync(name string) *lsv1alpha1.TargetSync {
	return s.TargetSync[s.Namespace+"/"+name]
}

// GetTarget retrieves a target by the given name.
func (s *State) GetTarget(name string) *lsv1alpha1.Target {
	return s.Targets[s.Namespace+"/"+name]
}

// GetContext retrieves a context by the given name.
func (s *State) GetContext(name string) *lsv1alpha1.Context {
	return s.Contexts[s.Namespace+"/"+name]
}

// GetAvailabilityCollection retrieves a AvailabilityCollection by the given name
func (s *State) GetAvailabilityCollection(name string) *provisioningv1alpha2.AvailabilityCollection {
	return s.AvailabilityCollections[s.Namespace+"/"+name]
}

// GetLsHealthCheck retrieves a LsHealthCheck by the given name
func (s *State) GetLsHealthCheck(name string) *lsv1alpha1.LsHealthCheck {
	return s.LsHealthChecks[s.Namespace+"/"+name]
}

// GetLsHealthCheckInNamespace retrieves a LsHealthCheck by the given name in the given namespace
func (s *State) GetLsHealthCheckInNamespace(name string, namespace string) *lsv1alpha1.LsHealthCheck {
	return s.LsHealthChecks[namespace+"/"+name]
}

// GetNamespaceRegistration retrieves a NamespaceRegistration by the given name in the given namespace
func (s *State) GetNamespaceRegistration(name string) *dataplanev1alpha1.NamespaceRegistration {
	return s.NamespaceRegistrations[s.Namespace+"/"+name]
}

// GetSubjectListInNamespace retrieves a SubjectList by the given name in the given namespace
func (s *State) GetSubjectList(name string) *dataplanev1alpha1.SubjectList {
	return s.SubjectLists[subjectsync.LS_USER_NAMESPACE+"/"+name]
}

// AddObject adds a client.Object to the state.
func (s *State) AddObject(object client.Object) {
	switch o := object.(type) {
	case *provisioningv1alpha2.LandscaperDeployment:
		s.Deployments[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *provisioningv1alpha2.Instance:
		s.Instances[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *provisioningv1alpha2.ServiceTargetConfig:
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
	case *lsv1alpha1.Execution:
		s.Executions[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lsv1alpha1.DeployItem:
		s.DeployItems[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lsv1alpha1.TargetSync:
		s.TargetSync[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lsv1alpha1.Target:
		s.Targets[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lsv1alpha1.Context:
		s.Contexts[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *provisioningv1alpha2.AvailabilityCollection:
		s.AvailabilityCollections[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *lsv1alpha1.LsHealthCheck:
		s.LsHealthChecks[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *dataplanev1alpha1.NamespaceRegistration:
		s.NamespaceRegistrations[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	case *dataplanev1alpha1.SubjectList:
		s.SubjectLists[types.NamespacedName{Name: o.Name, Namespace: o.Namespace}.String()] = o.DeepCopy()
	}
}
