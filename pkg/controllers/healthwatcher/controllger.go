// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package healthwatcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apitypes "k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/apis/installation"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	"github.com/gardener/landscaper-service/pkg/operation"
)

type Controller struct {
	operation.Operation
	log                 logging.Logger
	kubeClientExtractor ServiceTargetConfigKubeClientExtractorInterface
}

// ServiceTargetConfigKubeClientExtractorInterface implements functionality to create a kubeclient from a servive target config ref
type ServiceTargetConfigKubeClientExtractorInterface interface {
	GetKubeClientFromServiceTargetConfig(ctx context.Context, name string, namespace string, client client.Client) (client.Client, error)
}

func NewController(logger logging.Logger, c client.Client, scheme *runtime.Scheme, config *coreconfig.LandscaperServiceConfiguration) (reconcile.Reconciler, error) {
	ctrl := &Controller{
		log: logger,
	}

	op := operation.NewOperation(c, scheme, config)
	ctrl.Operation = *op
	ctrl.kubeClientExtractor = &ServiceTargetConfigKubeClientExtractor{}
	return ctrl, nil
}

// NewTestActuator creates a new controller for testing purposes.
func NewTestActuator(op operation.Operation, kubeClientExtractor ServiceTargetConfigKubeClientExtractorInterface, logger logging.Logger) *Controller {
	ctrl := &Controller{
		Operation:           op,
		log:                 logger,
		kubeClientExtractor: kubeClientExtractor,
	}
	return ctrl
}

func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger, ctx := c.log.StartReconcileAndAddToContext(ctx, req)

	//get availabilityCollection
	availabilityCollection := &lssv1alpha1.AvailabilityCollection{}
	if err := c.Client().Get(ctx, req.NamespacedName, availabilityCollection); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info(err.Error())
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	//dont run if spec has not changed and we are not in time yet
	if availabilityCollection.ObjectMeta.Generation == availabilityCollection.Status.ObservedGeneration &&
		time.Since(availabilityCollection.Status.LastRun.Time) < c.Operation.Config().AvailabilityMonitoring.PeriodicCheckInterval.Duration {
		return reconcile.Result{}, nil
	}

	//clean status
	availabilityCollection.Status.Instances = []lssv1alpha1.AvailabilityInstance{}

	for _, instanceRefToWatch := range availabilityCollection.Spec.InstanceRefs {
		//get instance
		instance := &lssv1alpha1.Instance{}
		if err := c.Client().Get(ctx, apitypes.NamespacedName{Name: instanceRefToWatch.Name, Namespace: instanceRefToWatch.Namespace}, instance); err != nil {
			if apierrors.IsNotFound(err) {
				logger.Info(err.Error())
				return reconcile.Result{}, nil
			}
			return reconcile.Result{}, err
		}

		availabilityInstance := lssv1alpha1.AvailabilityInstance{
			ObjectReference: lssv1alpha1.ObjectReference{
				Name:      instance.Name,
				Namespace: instance.Namespace,
			},
		}

		//get referred installation
		if instance.Status.InstallationRef == nil || instance.Status.InstallationRef.Name == "" || instance.Status.InstallationRef.Namespace == "" {
			continue
		}
		installation := &lsv1alpha1.Installation{}
		if err := c.Client().Get(ctx, apitypes.NamespacedName{Name: instance.Status.InstallationRef.Name, Namespace: instance.Status.InstallationRef.Namespace}, installation); err != nil {
			if apierrors.IsNotFound(err) {
				logger.Info(err.Error())
				continue
			}
			logger.Info(fmt.Sprintf("could not load installation from installation reference: %s", err.Error()))
			setAvailabilityInstanceStatusToFailed(&availabilityInstance, "could not load installation from installation reference")
			availabilityCollection.Status.Instances = append(availabilityCollection.Status.Instances, availabilityInstance)
			continue
		}

		//check if installation is not progressing
		if installation.Status.Phase == lsv1alpha1.ComponentPhaseProgressing {
			logger.Info(fmt.Sprintf("installation %s:%s for instance %s:%s is progressing, not health check monitoring", installation.Namespace, installation.Name, instance.Namespace, instance.Name))
			continue
		}

		//check that servicetargetconfref exists exists
		if instance.Spec.ServiceTargetConfigRef.Name == "" || instance.Spec.ServiceTargetConfigRef.Namespace == "" {
			logger.Info(fmt.Sprintf("instance %s:%s does not have a ServiceTargetConfig ref", instance.Namespace, instance.Name))
			setAvailabilityInstanceStatusToFailed(&availabilityInstance, "instance does not have a ServiceTargetConfigRef")
			availabilityCollection.Status.Instances = append(availabilityCollection.Status.Instances, availabilityInstance)
			continue
		}

		//get kubeconfig from secret referenced in ServiceTargetConfig so a credential rotation is automatically handled
		targetClient, err := c.kubeClientExtractor.GetKubeClientFromServiceTargetConfig(ctx, instance.Spec.ServiceTargetConfigRef.Name, instance.Spec.ServiceTargetConfigRef.Namespace, c.Client())
		if err != nil {
			logger.Info(err.Error())
			setAvailabilityInstanceStatusToFailed(&availabilityInstance, "could not create k8s client from target config")
			availabilityCollection.Status.Instances = append(availabilityCollection.Status.Instances, availabilityInstance)
			continue
		}

		targetClusterNamespace, err := extractTargetClusterNamespaceFromInstallation(*installation)
		if err != nil {
			logger.Info(err.Error())
			setAvailabilityInstanceStatusToFailed(&availabilityInstance, "could not read target cluster namespace from installation")
			availabilityCollection.Status.Instances = append(availabilityCollection.Status.Instances, availabilityInstance)
			continue
		}

		//collect lshealthcheck
		lsHealthchecks := &lsv1alpha1.LsHealthCheckList{}
		err = targetClient.List(ctx, lsHealthchecks, client.InNamespace(targetClusterNamespace))
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Info(err.Error())
				setAvailabilityInstanceStatusToFailed(&availabilityInstance, "lsHealthCheck not found on target")
				continue
			}
			logger.Info(fmt.Sprintf("could not load lshealthcheck from cluster: %s", err.Error()))
			setAvailabilityInstanceStatusToFailed(&availabilityInstance, "failed retrieving lshealthcheck cr")
			availabilityCollection.Status.Instances = append(availabilityCollection.Status.Instances, availabilityInstance)
			continue
		}

		transferLsHealthCheckStatusToAvailabilityInstance(&availabilityInstance, lsHealthchecks, c.Config().AvailabilityMonitoring.LSHealthCheckTimeout.Duration)
		availabilityCollection.Status.Instances = append(availabilityCollection.Status.Instances, availabilityInstance)
	}
	availabilityCollection.Status.Self = c.getLsHealthCheckFromSelfLandscaper(ctx, c.Config().AvailabilityMonitoring.SelfLandscaperNamespace)
	availabilityCollection.Status.ObservedGeneration = availabilityCollection.ObjectMeta.Generation
	availabilityCollection.Status.LastRun = v1.NewTime(time.Now())

	logFailedInstances(logger, *availabilityCollection)

	//write to status
	if err := c.Client().Status().Update(ctx, availabilityCollection); err != nil {
		return reconcile.Result{}, fmt.Errorf("unable to update availability collection: %w", err)
	}

	//Requeue to run again
	return reconcile.Result{RequeueAfter: c.Config().AvailabilityMonitoring.PeriodicCheckInterval.Duration}, nil

}

func (c *Controller) getLsHealthCheckFromSelfLandscaper(ctx context.Context, namespace string) lssv1alpha1.AvailabilityInstance {
	availabilityInstance := lssv1alpha1.AvailabilityInstance{
		ObjectReference: lssv1alpha1.ObjectReference{
			Name:      "self",
			Namespace: namespace,
		},
	}

	//collect lshealthcheck
	lsHealthchecks := &lsv1alpha1.LsHealthCheckList{}
	err := c.Client().List(ctx, lsHealthchecks, client.InNamespace(namespace))
	if err != nil {
		if apierrors.IsNotFound(err) {
			c.log.Info(err.Error())
			setAvailabilityInstanceStatusToFailed(&availabilityInstance, "lsHealthCheck not found on target")

		}
		c.log.Info(fmt.Sprintf("could not load lshealthcheck from cluster: %s", err.Error()))
		setAvailabilityInstanceStatusToFailed(&availabilityInstance, "failed retrieving lshealthcheck cr")
	}
	transferLsHealthCheckStatusToAvailabilityInstance(&availabilityInstance, lsHealthchecks, c.Config().AvailabilityMonitoring.LSHealthCheckTimeout.Duration)
	return availabilityInstance
}

func logFailedInstances(logger logging.Logger, availabilityCollection lssv1alpha1.AvailabilityCollection) {
	failedInstances := []lssv1alpha1.AvailabilityInstance{}

	for _, inst := range availabilityCollection.Status.Instances {
		if inst.Status == string(lsv1alpha1.LsHealthCheckStatusFailed) {
			failedInstances = append(failedInstances, inst)
		}
	}
	if availabilityCollection.Status.Self.Status == string(lsv1alpha1.LsHealthCheckStatusFailed) {
		failedInstances = append(failedInstances, availabilityCollection.Status.Self)
	}
	if len(failedInstances) > 0 {
		logger.Info("av check failed", "failed instances", failedInstances)
	}
}

func setAvailabilityInstanceStatusToFailed(availabilityInstance *lssv1alpha1.AvailabilityInstance, failedReason string) {
	availabilityInstance.Status = string(lsv1alpha1.LsHealthCheckStatusFailed)
	availabilityInstance.FailedReason = failedReason
}

func transferLsHealthCheckStatusToAvailabilityInstance(availabilityInstance *lssv1alpha1.AvailabilityInstance, lsHealthChecks *lsv1alpha1.LsHealthCheckList, timeout time.Duration) {
	if len(lsHealthChecks.Items) != 1 {
		setAvailabilityInstanceStatusToFailed(availabilityInstance, "number of lsHealthChecks found != 1")
		return
	}

	healthCheck := lsHealthChecks.Items[0]
	if time.Since(healthCheck.LastUpdateTime.Time) > timeout {
		setAvailabilityInstanceStatusToFailed(availabilityInstance, fmt.Sprintf("timeout - last update time not recent enough (timeout %s)", timeout.String()))
	} else {
		availabilityInstance.Status = string(healthCheck.Status)
		availabilityInstance.FailedReason = healthCheck.Description
	}
}

func extractTargetClusterNamespaceFromInstallation(inst lsv1alpha1.Installation) (string, error) {
	hostingClusterNamespaceRaw, ok := inst.Spec.ImportDataMappings[installation.HostingClusterNamespaceImportName]
	if !ok {
		return "", errors.New("could not find hostingClusterNamespace in installation reference")
	}
	var targetClusterNamespace string
	if err := json.Unmarshal(hostingClusterNamespaceRaw.RawMessage, &targetClusterNamespace); err != nil {
		return "", fmt.Errorf("failed to unmarshal hostingClusterNamespace: %w", err)
	}
	return targetClusterNamespace, nil
}

type ServiceTargetConfigKubeClientExtractor struct{}

func (e *ServiceTargetConfigKubeClientExtractor) GetKubeClientFromServiceTargetConfig(ctx context.Context, name string, namespace string, client client.Client) (client.Client, error) {
	if name == "" || namespace == "" {
		return nil, errors.New("name or namespace of serviceTargetConfig is empty")
	}
	serviceTargetConfig := &lssv1alpha1.ServiceTargetConfig{}
	if err := client.Get(ctx, apitypes.NamespacedName{Name: name, Namespace: namespace}, serviceTargetConfig); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed loading ServiceTargetConfig %s:%s - not found: %w", name, namespace, err)
		}
		return nil, fmt.Errorf("could not load ServiceTargetConfig from instance reference: %w", err)
	}
	secretWithKubeconf := &corev1.Secret{}
	if err := client.Get(ctx, apitypes.NamespacedName{Name: serviceTargetConfig.Spec.SecretRef.Name, Namespace: serviceTargetConfig.Spec.SecretRef.Namespace}, secretWithKubeconf); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed loading secret %s:%s for ServiceTargetConfig %s:%s - not found: %w", serviceTargetConfig.Spec.SecretRef.Name, serviceTargetConfig.Spec.SecretRef.Namespace, name, namespace, err)
		}
		return nil, fmt.Errorf("could not load secret %s:%s for ServiceTargetConfig %s:%s: %w", serviceTargetConfig.Spec.SecretRef.Name, serviceTargetConfig.Spec.SecretRef.Namespace, name, namespace, err)
	}

	_, targetClient, _, err := getKubeClientFromSecret(*secretWithKubeconf, serviceTargetConfig.Spec.SecretRef.Key)
	if err != nil {
		return nil, fmt.Errorf("failed building kubeclient for target: %w", err)
	}
	return targetClient, nil
}

func getKubeClientFromSecret(secret corev1.Secret, key string) (*rest.Config, client.Client, kubernetes.Interface, error) {
	kubeconfigBytes, ok := secret.Data[key]
	if !ok {
		return nil, nil, nil, fmt.Errorf("could not found key %s in secret", key)
	}

	kubeconfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		return nil, nil, nil, err
	}
	restConfig, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, nil, nil, err
	}

	kubeClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		return nil, nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, nil, err
	}
	return restConfig, kubeClient, clientset, nil
}
