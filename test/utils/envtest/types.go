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

	lsscoreinstall "github.com/gardener/landscaper-service/pkg/apis/core/install"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	lscoreinstall "github.com/gardener/landscaper/apis/core/install"
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
	// InstallationGVK is the GVK for installations.
	InstallationGVK schema.GroupVersionKind
	// TargetGVK is the GVK for targets.
	TargetGVK schema.GroupVersionKind
)

func init() {
	var err error

	lsscoreinstall.Install(LandscaperServiceScheme)
	lscoreinstall.Install(LandscaperServiceScheme)
	utilruntime.Must(kubernetescheme.AddToScheme(LandscaperServiceScheme))

	DeploymentGVK, err = apiutil.GVKForObject(&lssv1alpha1.LandscaperDeployment{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	InstanceGVK, err = apiutil.GVKForObject(&lssv1alpha1.Instance{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	ConfigGVK, err = apiutil.GVKForObject(&lssv1alpha1.ServiceTargetConfig{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	SecretGVK, err = apiutil.GVKForObject(&corev1.Secret{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	InstallationGVK, err = apiutil.GVKForObject(&lsv1alpha1.Installation{}, LandscaperServiceScheme)
	utilruntime.Must(err)
	TargetGVK, err = apiutil.GVKForObject(&lsv1alpha1.Target{}, LandscaperServiceScheme)
	utilruntime.Must(err)
}