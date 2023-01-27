// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"time"

	admissionv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/jsonpath"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/apis/core/v1alpha1/targettypes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"
	cliinstallations "github.com/gardener/landscapercli/cmd/installations"
	cliutil "github.com/gardener/landscapercli/pkg/util"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// GetLandscaperVersion reads the landscaper version out of the LAAS component descriptor
func GetLandscaperVersion(repoRootDir string) (string, error) {

	var landscaperVersion string

	compReferencesFile := path.Join(repoRootDir, ".landscaper", "landscaper-instance", "component-references.yaml")
	raw, err := ioutil.ReadFile(compReferencesFile)
	if err != nil {
		return "", err
	}

	r := bytes.NewReader(raw)
	dec := yaml.NewYAMLOrJSONDecoder(r, 1024)

	var componentReferences map[string]interface{}
	for dec.Decode(&componentReferences) == nil {
		name, ok := componentReferences["name"]
		if !ok {
			continue
		}
		if name != "landscaper" {
			continue
		}
		version, ok := componentReferences["version"]
		if !ok {
			continue
		}
		landscaperVersion = version.(string)
		break
	}

	if len(landscaperVersion) == 0 {
		return "", fmt.Errorf("landscaper version not found")
	}

	return landscaperVersion, nil
}

// ForceDeleteInstallations calls landscaper-cli force delete on the given namespace.
func ForceDeleteInstallations(ctx context.Context, kclient client.Client, kubeConfig, namespace string) error {
	logger, _ := logging.FromContextOrNew(ctx, nil)

	installationList := &lsv1alpha1.InstallationList{}
	if err := kclient.List(ctx, installationList, &client.ListOptions{Namespace: namespace}); err != nil {
		return err
	}

	for _, installation := range installationList.Items {
		if err := kclient.Get(ctx, types.NamespacedName{Name: installation.Name, Namespace: installation.Namespace}, &installation); err != nil {
			if k8serrors.IsNotFound(err) {
				continue
			}
			return err
		}

		logger.Info("Deleting installation", "name", installation.Name)
		forceDeleteCommand := cliinstallations.NewForceDeleteCommand(ctx)
		forceDeleteCommand.SetArgs([]string{
			installation.Name,
			"--kubeconfig",
			kubeConfig,
			"--namespace",
			namespace,
		})
		if err := forceDeleteCommand.Execute(); err != nil {
			return err
		}
	}

	return nil
}

// CleanupLaasResources tries to remove all landscaper deployments, instances and service target configs in the given namespace.
func CleanupLaasResources(ctx context.Context, kclient client.Client, namespace string, sleepTime time.Duration, maxRetries int) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	deploymentList := &lssv1alpha1.LandscaperDeploymentList{}
	if err := kclient.List(ctx, deploymentList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list landscaper deployments: %w", err)
	}

	for _, deployment := range deploymentList.Items {
		logger.Info("Deleting landscaper deployment", lc.KeyResource, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}.String())
		if err := kclient.Delete(ctx, &deployment); err != nil {
			if !k8serrors.IsNotFound(err) {
				return fmt.Errorf("failed to delete landscaper deployment %q: %w", deployment.Name, err)
			}
		}
		_, err := cliutil.CheckAndWaitUntilObjectNotExistAnymore(
			kclient,
			types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace},
			&deployment, sleepTime, maxRetries)

		if err != nil {
			return fmt.Errorf("failed to wait for landscaper deployment %q to be deleted: %w", deployment.Name, err)
		}
	}

	serviceTargetConfigList := &lssv1alpha1.ServiceTargetConfigList{}
	if err := kclient.List(ctx, serviceTargetConfigList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list service target configs: %w", err)
	}

	for _, serviceTargetConfig := range serviceTargetConfigList.Items {
		logger.Info("Deleting service target config", lc.KeyResource, types.NamespacedName{Name: serviceTargetConfig.Name, Namespace: serviceTargetConfig.Namespace}.String())
		if err := kclient.Delete(ctx, &serviceTargetConfig); err != nil {
			if !k8serrors.IsNotFound(err) {
				return fmt.Errorf("failed to delete service target config %q: %w", serviceTargetConfig.Name, err)
			}
		}
		_, err := cliutil.CheckAndWaitUntilObjectNotExistAnymore(
			kclient,
			types.NamespacedName{Name: serviceTargetConfig.Name, Namespace: serviceTargetConfig.Namespace},
			&serviceTargetConfig, sleepTime, maxRetries)

		if err != nil {
			return fmt.Errorf("failed to wait for service target config %q to be deleted: %w", serviceTargetConfig.Name, err)
		}
	}

	return nil
}

// FinalizerPatch removes the finalizer of the given object.
type FinalizerPatch struct {
	Patch map[string]interface{}
}

func (FinalizerPatch) Type() types.PatchType {
	return types.MergePatchType
}

func (FinalizerPatch) Data(_ client.Object) ([]byte, error) {
	patch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"finalizers": nil,
		},
	}

	patchRaw, err := json.Marshal(patch)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal patch: %w", err)
	}

	return patchRaw, nil
}

// RemoveFinalizerLandscaperResources removes the finalizer of all landscaper resources in the given namespace.
func RemoveFinalizerLandscaperResources(ctx context.Context, kclient client.Client, namespace string) error {
	patch := FinalizerPatch{}

	podList := &corev1.PodList{}
	if err := kclient.List(ctx, podList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	for _, pod := range podList.Items {
		if err := kclient.Patch(ctx, &pod, patch); err != nil {
			return fmt.Errorf("failed to patch pod %q: %w", pod.Name, err)
		}
	}

	installationList := &lsv1alpha1.InstallationList{}
	if err := kclient.List(ctx, installationList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list installations: %w", err)
	}

	for _, installation := range installationList.Items {
		if err := kclient.Patch(ctx, &installation, patch); err != nil {
			return fmt.Errorf("failed to patch installation %q: %w", installation.Name, err)
		}
	}

	executionList := &lsv1alpha1.ExecutionList{}
	if err := kclient.List(ctx, executionList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list executions: %w", err)
	}

	for _, execution := range executionList.Items {
		if err := kclient.Patch(ctx, &execution, patch); err != nil {
			return fmt.Errorf("failed to patch execution %q: %w", execution.Name, err)
		}
	}

	deployItemList := &lsv1alpha1.DeployItemList{}
	if err := kclient.List(ctx, deployItemList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list deploy items: %w", err)
	}

	for _, deployItem := range deployItemList.Items {
		if err := kclient.Patch(ctx, &deployItem, patch); err != nil {
			return fmt.Errorf("failed to patch deploy item %q: %w", deployItem.Name, err)
		}
	}

	return nil
}

// RemoveFinalizerLaasResources removes the finalizer of all laas resources in the given namespace.
func RemoveFinalizerLaasResources(ctx context.Context, kclient client.Client, namespace string) error {
	patch := FinalizerPatch{}

	deploymentList := &lssv1alpha1.LandscaperDeploymentList{}
	if err := kclient.List(ctx, deploymentList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list landscaper deployments: %w", err)
	}

	for _, deployment := range deploymentList.Items {
		if err := kclient.Patch(ctx, &deployment, patch); err != nil {
			return fmt.Errorf("failed to patch landscaper deployment: %w", err)
		}
	}

	instanceList := &lssv1alpha1.InstanceList{}
	if err := kclient.List(ctx, instanceList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	for _, instance := range instanceList.Items {
		if err := kclient.Patch(ctx, &instance, patch); err != nil {
			return fmt.Errorf("failed to patch instance: %w", err)
		}
	}

	serviceTargetConfigList := &lssv1alpha1.ServiceTargetConfigList{}
	if err := kclient.List(ctx, serviceTargetConfigList, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list service target configs: %w", err)
	}

	for _, serviceTargetConfig := range serviceTargetConfigList.Items {
		if err := kclient.Patch(ctx, &serviceTargetConfig, patch); err != nil {
			return fmt.Errorf("failed to patch service target config: %w", err)
		}
	}

	return nil
}

// DeleteValidatingWebhookConfiguration deletes the validating webhook configuration of the given name in the given namespace.
func DeleteValidatingWebhookConfiguration(ctx context.Context, kclient client.Client, name, namespace string) error {
	validationConfig := &admissionv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	if err := kclient.Delete(ctx, validationConfig); err != nil {
		if !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete validating webhook configuration %s/%s: %w", name, namespace, err)
		}
	}

	return nil
}

// DeleteTargetClusterNamespaces tries to delete all target cluster namespaces in the cluster.
func DeleteTargetClusterNamespaces(ctx context.Context, kclient client.Client, sleepTime time.Duration, maxRetries int) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	namespaces := &corev1.NamespaceList{}

	if err := kclient.List(ctx, namespaces); err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	for _, namespace := range namespaces.Items {
		if strings.HasPrefix(namespace.Name, "tc-") {
			logger.Info("pruning namespace", lc.KeyResource, namespace.Name)

			if err := RemoveFinalizerLandscaperResources(ctx, kclient, namespace.Name); err != nil {
				return fmt.Errorf("failed to remove finalizers for landscaper resources in namespace %q: %w", namespace.Name, err)
			}

			deployment := &appsv1.Deployment{}
			if err := kclient.DeleteAllOf(ctx, deployment, &client.DeleteAllOfOptions{
				ListOptions: client.ListOptions{
					Namespace: namespace.Name,
				},
				DeleteOptions: client.DeleteOptions{},
			}); err != nil {
				return fmt.Errorf("failed to delete deployments in namespace %q: %w", namespace.Name, err)
			}

			statefulSet := &appsv1.StatefulSet{}
			if err := kclient.DeleteAllOf(ctx, statefulSet, &client.DeleteAllOfOptions{
				ListOptions: client.ListOptions{
					Namespace: namespace.Name,
				},
				DeleteOptions: client.DeleteOptions{},
			}); err != nil {
				return fmt.Errorf("failed to delete stateful sets in namespace %q: %w", namespace.Name, err)
			}

			pod := &corev1.Pod{}
			if err := kclient.DeleteAllOf(ctx, pod, &client.DeleteAllOfOptions{
				ListOptions: client.ListOptions{
					Namespace: namespace.Name,
				},
				DeleteOptions: client.DeleteOptions{},
			}); err != nil {
				return fmt.Errorf("failed to delete pods sets in namespace %q: %w", namespace.Name, err)
			}

			pvc := &corev1.PersistentVolumeClaim{}
			if err := kclient.DeleteAllOf(ctx, pvc, &client.DeleteAllOfOptions{
				ListOptions: client.ListOptions{
					Namespace: namespace.Name,
				},
				DeleteOptions: client.DeleteOptions{},
			}); err != nil {
				return fmt.Errorf("failed to delete persistent volume clains in namespace %q: %w", namespace.Name, err)
			}

			logger.Info("deleting namespace", lc.KeyResource, namespace.Name)
			err := cliutil.DeleteNamespace(kclient, namespace.Name, sleepTime, maxRetries)
			if err != nil {
				return fmt.Errorf("failed to delete namespace %q: %w", namespace.Name, err)
			}
		}
	}

	return nil
}

// BuildKubernetesClusterTargetWithSecretRef builds a landscaper target of the given kubeconfig in the given namespace using a secret reference.
func BuildKubernetesClusterTargetWithSecretRef(ctx context.Context, kclient client.Client, kubeConfig, name, namespace string) (*lsv1alpha1.Target, error) {
	kubeConfigContent, err := ioutil.ReadFile(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot read kubeconfig: %w", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"kubeconfig": string(kubeConfigContent),
		},
	}

	if err := kclient.Create(ctx, secret); err != nil {
		return nil, fmt.Errorf("failed to create default target secret: %w", err)
	}

	target := &lsv1alpha1.Target{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: lsv1alpha1.TargetSpec{
			Type: targettypes.KubernetesClusterTargetType,
			SecretRef: &lsv1alpha1.LocalSecretReference{
				Name: name,
				Key:  "kubeconfig",
			},
		},
	}

	if err := kclient.Create(ctx, target); err != nil {
		return nil, fmt.Errorf("failed to create default target: %w", err)
	}

	return target, nil
}

// BuildKubernetesClusterTarget builds a landscaper target of the given kubeconfig in the given namespace.
func BuildKubernetesClusterTarget(ctx context.Context, kclient client.Client, kubeConfig, name, namespace string) (*lsv1alpha1.Target, error) {
	kubeConfigContent, err := ioutil.ReadFile(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot read kubeconfig: %w", err)
	}
	kubeConfigContentStr := string(kubeConfigContent)

	targetConfig := targettypes.KubernetesClusterTargetConfig{
		Kubeconfig: targettypes.ValueRef{
			StrVal: &kubeConfigContentStr,
		},
	}

	targetConfigRaw, err := json.Marshal(targetConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal target config: %w", err)
	}
	targetConfigAnyJSON := lsv1alpha1.NewAnyJSON(targetConfigRaw)

	target := &lsv1alpha1.Target{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: lsv1alpha1.TargetSpec{
			Type:          targettypes.KubernetesClusterTargetType,
			Configuration: &targetConfigAnyJSON,
		},
	}

	if err := kclient.Create(ctx, target); err != nil {
		return nil, fmt.Errorf("failed to create default target: %w", err)
	}

	return target, nil
}

// BuildLandscaperContext builds a landscaper context containing the given registry pull secrets in the given namespace.
func BuildLandscaperContext(ctx context.Context, kclient client.Client, registryPullSecretsFile, name string, namespaces ...string) error {
	registryPullSecrets, err := ioutil.ReadFile(registryPullSecretsFile)
	if err != nil {
		return fmt.Errorf("failed to read registry pull secret: %w", err)
	}

	for _, namespace := range namespaces {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			StringData: map[string]string{
				corev1.DockerConfigJsonKey: string(registryPullSecrets),
			},
			Type: corev1.SecretTypeDockerConfigJson,
		}

		if err := kclient.Create(ctx, secret); err != nil {
			return fmt.Errorf("failed to create dockerconfigjson secret: %w", err)
		}

		landscaperContext := &lsv1alpha1.Context{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			RegistryPullSecrets: []corev1.LocalObjectReference{
				{
					Name: name,
				},
			},
		}

		if err := kclient.Create(ctx, landscaperContext); err != nil {
			return fmt.Errorf("failed to create landscaper context: %w", err)
		}
	}

	return nil
}

// BuildKubeClientForInstance builds a kubernetes client for the user kubeconfig of an instance.
func BuildKubeClientForInstance(instance *lssv1alpha1.Instance, scheme *runtime.Scheme) (client.Client, error) {
	kubeconfig, err := base64.StdEncoding.DecodeString(instance.Status.UserKubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode kubeconfig of instance %q: %w", instance.Name, err)
	}

	client, err := BuildKubeClient(string(kubeconfig), scheme)
	if err != nil {
		return nil, fmt.Errorf("failed to build kube client for instance %q: %w", instance.Name, err)
	}

	return client, nil
}

// BuildKubeClient builds a kubernetes client from a kubeconfig
func BuildKubeClient(kubeconfig string, scheme *runtime.Scheme) (client.Client, error) {
	clientCfg, err := clientcmd.Load([]byte(kubeconfig))
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	loader := clientcmd.NewDefaultClientConfig(*clientCfg, nil)
	restConfig, err := loader.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load rest config: %e", err)
	}

	client, err := client.New(restConfig, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, fmt.Errorf("failed create client: %e", err)
	}

	return client, nil
}

// DeleteTestShootClusters tries to delete all existing integration test shoot clusters.
func DeleteTestShootClusters(ctx context.Context, gardenerServiceAccountKubeconfigFile, gardenerProject string, scheme *runtime.Scheme) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	kubeconfig, err := ioutil.ReadFile(gardenerServiceAccountKubeconfigFile)
	if err != nil {
		return fmt.Errorf("failed to read gardener service account kubeconfig from file %q: %w", gardenerServiceAccountKubeconfigFile, err)
	}

	saClient, err := BuildKubeClient(string(kubeconfig), scheme)
	if err != nil {
		return fmt.Errorf("failed to build kube client for gardener service account: %w", err)
	}

	shootList := &unstructured.UnstructuredList{
		Object: map[string]interface{}{
			"apiVersion": "core.gardener.cloud/v1beta1",
			"kind":       "Shoot",
		},
	}

	shootNamespace := fmt.Sprintf("garden-%s", gardenerProject)

	labelSelector, _ := labels.Parse(lssv1alpha1.ShootInstanceNameLabel)
	listOptions := client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     shootNamespace,
	}

	if err := saClient.List(ctx, shootList, &listOptions); err != nil {
		return fmt.Errorf("failed to list shoots: %w", err)
	}

	for _, shoot := range shootList.Items {
		if shoot.GetDeletionTimestamp() != nil {
			logger.Info("shoot is already being deleted", lc.KeyResource, shoot.GetName())
			continue
		}

		annotations := shoot.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}
		annotations["confirmation.gardener.cloud/deletion"] = "true"
		shoot.SetAnnotations(annotations)
		if err := saClient.Update(ctx, &shoot); err != nil {
			return fmt.Errorf("failed to annotate shoot %q: %w", shoot.GetName(), err)
		}

		logger.Info("deleting shoot", lc.KeyResource, shoot.GetName())
		if err := saClient.Delete(ctx, &shoot); err != nil {
			if !apierrors.IsNotFound(err) {
				return fmt.Errorf("failed to delete shoot %q: %w", shoot.GetName(), err)
			}
		}
	}

	return nil
}

// ParseIngressDomain parses the ingress domain out of a kubeconfig file.
func ParseIngressDomain(kubeconfigFile string) (string, error) {
	var hostingClusterKubeconfigMap map[string]interface{}
	hostingClusterKubeconfig, err := ioutil.ReadFile(kubeconfigFile)
	if err != nil {
		return "", fmt.Errorf("failed to read hosting cluster kubeconfig: %w", err)
	}
	if err := yaml.Unmarshal(hostingClusterKubeconfig, &hostingClusterKubeconfigMap); err != nil {
		return "", fmt.Errorf("failed to unmarshal hosting cluster kubeconfig: %w", err)
	}

	path := jsonpath.New("clusterServer").AllowMissingKeys(true)
	err = path.Parse("{.clusters[0].cluster.server}")
	if err != nil {
		return "", fmt.Errorf("failed to parse jsonpath: %w", err)
	}

	results, err := path.FindResults(hostingClusterKubeconfigMap)
	if err != nil {
		return "", fmt.Errorf("failed to get cluster server: %w", err)
	}

	value := results[0][0].Interface()
	clusterServer, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("field has an invalid type")
	}

	rexp, err := regexp.Compile(`[^\.]*\.(.*)$`)
	if err != nil {
		return "", fmt.Errorf("failed to parse regex")
	}
	match := rexp.FindStringSubmatch(clusterServer)
	if match == nil || len(match) < 2 {
		return "", fmt.Errorf("failed to parse url")
	}
	ingressDomain := fmt.Sprintf("ingress.%s", match[1])

	return ingressDomain, nil
}
