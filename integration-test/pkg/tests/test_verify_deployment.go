// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	cliutil "github.com/gardener/landscapercli/pkg/util"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/test/integration/pkg/test"
)

type VerifyDeploymentRunner struct {
	ctx         context.Context
	log         logr.Logger
	kclient     client.Client
	config      *test.TestConfig
	target      *lsv1alpha1.Target
	testObjects *test.SharedTestObjects
}

func (r *VerifyDeploymentRunner) Init(
	ctx context.Context, log logr.Logger,
	kclient client.Client, config *test.TestConfig,
	target *lsv1alpha1.Target, testObjects *test.SharedTestObjects) {
	r.ctx = ctx
	r.log = log.WithName(r.Name())
	r.kclient = kclient
	r.config = config
	r.target = target
	r.testObjects = testObjects
}

func (r *VerifyDeploymentRunner) Name() string {
	return "VerifyDeployment"
}

func (r *VerifyDeploymentRunner) Description() string {
	description := `
This test verifies that a tenant Landscaper deployment has been installed correctly.
The test succeeds when all required pods (api server, etcd, landscaper ...) are running in the tenant namespace and
the connection to the virtual cluster can be established. Otherwise the test fails.
`
	return description
}

func (r *VerifyDeploymentRunner) String() string {
	return r.Name()
}

func (r *VerifyDeploymentRunner) Run() error {
	for _, deployment := range r.testObjects.LandscaperDeployments {
		if err := r.verifyDeployment(deployment); err != nil {
			return err
		}
	}
	return nil
}

func (r *VerifyDeploymentRunner) verifyDeployment(deployment *lssv1alpha1.LandscaperDeployment) error {
	instance := &lssv1alpha1.Instance{}
	if err := r.kclient.Get(
		r.ctx,
		types.NamespacedName{Name: deployment.Status.InstanceRef.Name, Namespace: deployment.Status.InstanceRef.Namespace},
		instance); err != nil {
		return fmt.Errorf("failed to get instance for deployment %q: %w", deployment.Name, err)
	}

	installation := &lsv1alpha1.Installation{}
	if err := r.kclient.Get(
		r.ctx,
		types.NamespacedName{Name: instance.Status.InstallationRef.Name, Namespace: instance.Status.InstallationRef.Namespace},
		installation); err != nil {
		return fmt.Errorf("failed to get installation for instance %q: %w", instance.Name, err)
	}

	hostingClusterNamespaceRaw, ok := installation.Spec.ImportDataMappings["hostingClusterNamespace"]

	if !ok {
		return fmt.Errorf("installation has no hostingClusterNamespace setting")
	}

	var hostingClusterNamespace string
	if err := json.Unmarshal(hostingClusterNamespaceRaw.RawMessage, &hostingClusterNamespace); err != nil {
		return fmt.Errorf("failed to unmarshal hostingClusterNamespace: %w", err)
	}

	if err := r.verifyPods(hostingClusterNamespace, len(instance.Spec.LandscaperConfiguration.Deployers)); err != nil {
		return err
	}

	if err := r.verifyKubeconfig(instance); err != nil {
		return err
	}

	r.testObjects.HostingClusterNamespaces = append(r.testObjects.HostingClusterNamespaces, hostingClusterNamespace)

	return nil
}

func (r *VerifyDeploymentRunner) verifyPods(namespace string, numDeployers int) error {
	expectedPods := []string{
		"landscaper-controller",
		"landscaper-webhooks",
		"etcd-main",
		"etcd-events",
		"apiserver",
		"controller-manager,",
	}

	r.log.Info("waiting for pods to be created")

	podList := &corev1.PodList{}
	timeout, err := cliutil.CheckConditionPeriodically(func() (bool, error) {
		if err := r.kclient.List(r.ctx, podList, &client.ListOptions{Namespace: namespace}); err != nil {
			return false, fmt.Errorf("failed to list pods in namespace %q: %w", namespace, err)
		}
		if len(podList.Items) >= (len(expectedPods) + numDeployers) {
			return true, nil
		}
		return false, nil
	}, r.config.SleepTime, r.config.MaxRetries*5)

	if err != nil {
		return err
	}
	if timeout {
		return fmt.Errorf("incomplete number pods in namespace %q, expected %d, actual %d", namespace, len(expectedPods)+numDeployers, len(podList.Items))
	}

	r.log.Info("waiting for pods to become running")

	timeout, err = cliutil.CheckConditionPeriodically(func() (bool, error) {
		if err := r.kclient.List(r.ctx, podList, &client.ListOptions{Namespace: namespace}); err != nil {
			return false, fmt.Errorf("failed to list pods in namespace %q: %w", namespace, err)
		}

		for _, pod := range podList.Items {
			if pod.Status.Phase != corev1.PodRunning {
				return false, nil
			}
		}

		return true, nil
	}, r.config.SleepTime, r.config.MaxRetries)

	if err != nil {
		return err
	}
	if timeout {
		return fmt.Errorf("not all pods in namespace %q are running", namespace)
	}

	return nil
}

func (r *VerifyDeploymentRunner) verifyKubeconfig(instance *lssv1alpha1.Instance) error {
	r.log.Info("verifying kubeconfig for instance", "name", instance.Name)

	timeout, err := cliutil.CheckConditionPeriodically(func() (bool, error) {
		if err := r.kclient.Get(r.ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, instance); err != nil {
			return false, err
		}

		return len(instance.Status.ClusterKubeconfig) > 0, nil
	}, r.config.SleepTime, r.config.MaxRetries)

	if timeout {
		return fmt.Errorf("timeout while reading ClusterKubeconfig for instance %q", instance.Name)
	}
	if err != nil {
		return fmt.Errorf("error while reading ClusterKubeconfig for instance %q: %w", instance.Name, err)
	}

	kubeconfig, err := base64.StdEncoding.DecodeString(instance.Status.ClusterKubeconfig)
	if err != nil {
		return fmt.Errorf("failed to decode kubeconfig of instance %q: %w", instance.Name, err)
	}

	laasClientCfg, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig of instance %q: %w", instance.Name, err)
	}

	loader := clientcmd.NewDefaultClientConfig(*laasClientCfg, nil)
	laasRestConfig, err := loader.ClientConfig()
	if err != nil {
		return fmt.Errorf("failed to load rest config of instance %q: %err", instance.Name, err)
	}

	laasClient, err := client.New(laasRestConfig, client.Options{
		Scheme: test.Scheme(),
	})
	if err != nil {
		return fmt.Errorf("failed create client for instance %q: %err", instance.Name, err)
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "itest-",
		},
	}

	if err := laasClient.Create(r.ctx, namespace); err != nil {
		return fmt.Errorf("failed to create namespace on cluster for instance %q: %w", instance.Name, err)
	}

	installationList := &lsv1alpha1.InstallationList{}
	if err := laasClient.List(r.ctx, installationList, &client.ListOptions{Namespace: namespace.Name}); err != nil {
		return fmt.Errorf("failed to list installations on cluster for instance %q: %w", instance.Name, err)
	}

	return nil
}
