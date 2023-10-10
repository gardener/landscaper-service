// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package envtest

import (
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubernetescheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	lscoreinstall "github.com/gardener/landscaper/apis/core/install"

	lsscoreinstall "github.com/gardener/landscaper-service/pkg/apis/core/install"
	lssv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha2"
)

var (
	// LandscaperServiceScheme is the scheme used for testing.
	LandscaperServiceScheme = runtime.NewScheme()
	// DeploymentGVK is the GVK for landscaper deployments.
	DeploymentGVK schema.GroupVersionKind
	// InstanceGVK is the GVK for instances.
	InstanceGVK schema.GroupVersionKind
	// ConfigGVK is the GVK for service target configs.
	ConfigGVK schema.GroupVersionKind
	// SecretGVK is the GVK for secrets.
	SecretGVK schema.GroupVersionKind
	// ConfigMapGVK is the GVK for config maps.
	ConfigMapGVK schema.GroupVersionKind
	// InstallationGVK is the GVK for installations.
	InstallationGVK schema.GroupVersionKind
	// ExecutionGVK is the GVK for executions.
	ExecutionGVK schema.GroupVersionKind
	// DeployItemGVK is the GVK for deploy items.
	DeployItemGVK schema.GroupVersionKind
	// TargetSyncGVK is the GVK for target syncs.
	TargetSyncGVK schema.GroupVersionKind
	// TargetGVK is the GVK for targets.
	TargetGVK schema.GroupVersionKind
	// ContextGVK is the GVK for contexts.
	ContextGVK schema.GroupVersionKind
	// AvailabilityCollectionGVK is the GVK for AvailabilityCollection.
	AvailabilityCollectionGVK schema.GroupVersionKind
	// LsHealthCheckGVK is the GVK for LsHealthCheck.
	LsHealthCheckGVK schema.GroupVersionKind
	// NamespaceRegistrationGVK is the GVK for SubjectList.
	NamespaceRegistrationGVK schema.GroupVersionKind
	// SubjectListGVK is the GVK for SubjectList.
	SubjectListGVK schema.GroupVersionKind
)

func init() {
	var err error

	lsscoreinstall.Install(LandscaperServiceScheme)
	lscoreinstall.Install(LandscaperServiceScheme)
	utilruntime.Must(kubernetescheme.AddToScheme(LandscaperServiceScheme))

	DeploymentGVK, err = apiutil.GVKForObject(&lssv1alpha2.LandscaperDeployment{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	InstanceGVK, err = apiutil.GVKForObject(&lssv1alpha2.Instance{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	ConfigGVK, err = apiutil.GVKForObject(&lssv1alpha2.ServiceTargetConfig{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	SecretGVK, err = apiutil.GVKForObject(&corev1.Secret{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	ConfigMapGVK, err = apiutil.GVKForObject(&corev1.ConfigMap{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	InstallationGVK, err = apiutil.GVKForObject(&lsv1alpha1.Installation{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	ExecutionGVK, err = apiutil.GVKForObject(&lsv1alpha1.Execution{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	DeployItemGVK, err = apiutil.GVKForObject(&lsv1alpha1.DeployItem{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	TargetSyncGVK, err = apiutil.GVKForObject(&lsv1alpha1.TargetSync{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	TargetGVK, err = apiutil.GVKForObject(&lsv1alpha1.Target{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	ContextGVK, err = apiutil.GVKForObject(&lsv1alpha1.Context{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	AvailabilityCollectionGVK, err = apiutil.GVKForObject(&lssv1alpha2.AvailabilityCollection{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	LsHealthCheckGVK, err = apiutil.GVKForObject(&lsv1alpha1.LsHealthCheck{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	NamespaceRegistrationGVK, err = apiutil.GVKForObject(&lssv1alpha2.NamespaceRegistration{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	SubjectListGVK, err = apiutil.GVKForObject(&lssv1alpha2.SubjectList{}, LandscaperServiceScheme)
	utilruntime.Must(err)
}
