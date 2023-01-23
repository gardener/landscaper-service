<!--
SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# Instances

Instances are kubernetes resources that represent an instance of a [LandscaperDeployment](LandscaperDeployments.md) which has been scheduled on a specific [ServiceTargetConfig](ServiceTargetConfigs.md).
Instances are not directly created by a user. Instances are created by the landscaper service controller out of a LandscaperDeployment when a matching ServiceTargetConfig has been found.
For each Instance the landscaper service controller will create a [Target](https://github.com/gardener/landscaper/blob/master/docs/technical/target_types.md) and an [Installation](https://github.com/gardener/landscaper/blob/master/docs/usage/Installations.md).

### Basic structure:

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: Instance
metadata:
  name: test
  namespace: my-namespace
spec:
  tenantId: "tenant-1"
  id: "1234"
  landscaperConfiguration:
    deployers:
      - helm
      - manifest
      - container
  serviceTargetConfigRef:
    name: default
    namespace: laas-system
    
status:
  installationRef:
    name: test
  namespace: my-namespace

  targetRef:
    name: test
    namespace: my-namespace
    
  contextRef:
    name: test
    namespace: my-namespace
    
  landscaperServiceComponent:
    name: github.com/gardener/landscaper/landscaper-service
    version: v0.19.0

  clusterEndpoint: "10.0.0.1:1234"
  userKubeconfig: "a3ViZWNvbmZpZyBjb250ZW50 ..."
  adminKubeconfig: "a3ViZWNvbmZpZyBjb250ZW50 ..."
  shootName: "a1b2c3d5"
  shootNamespace: "laas"
```

## TenantId

The `spec.tenantId` field contains the globally unique identifier of the owning tenant.

## Id

The `spec.id` field contains the id of the instance unique among all instances within the same namespace.

## Landscaper Configuration

The `spec.landscaperConfiguration` field specifies the Landscaper configuration that is defined by the parent LandscaperDeployment [landscaper configuration](LandscaperDeployments.md#landscaper-configuration).

## Component Reference

The `spec.componentReference` field specifies the component reference that is defined by the parent LandscaperDeployment [component reference](LandscaperDeployments.md#component-reference).

## Service Target Configuration Reference

The `spec.serviceTargetConfigRef` field specified the ServiceTargetConfig that has been selected for this Instance. 
The ServiceTarget config specifies the target kubernetes cluster on which the Landscaper will be deployed.

## Installation Reference

The `status.installationRef` field references the Installation that has been created by the landscaper service controller for this Instance.

## Target Reference

The `status.targetRef` field references the Target that has been created by the landscaper service controller for this Instance.
The target contains the kubeconfig that has been copied from the selected ServiceTargetConfig [secret reference](ServiceTargetConfigs.md#secret-reference).

## Context Reference

The `status.contextRef` field references the Landscaper Context that has been created by the landscaper service controller for this Instance.
The Context contains the repository context that is used by the Installation that has been created for this Instance.

## Landscaper Service Component

The `status.landscaperServiceComponent` field contains the landscaper service component name and version that is being used for the Landscaper instance.
The component name and version are set in the landscaper service controller configuration. 
The landscaper service controller is updated with a different landscaper service component version, all Instances will automatically be reconciled.
During the reconciliation the controller will update the deployed Landscaper to the new version.

## Cluster Endpoint

The `status.clusterEndpoint` field contains the API endpoint of the deployed Landscaper instance, i.e. it is used to
create and maintain Landscaper [Installations](https://github.com/gardener/landscaper/blob/master/docs/usage/Installations.md), 
which are handled by the corresponding Landscaper instance.

## User Kubeconfig

The `status.userKubeconfig` field contains the user kubeconfig which is used to access the deployed Landscaper (user restricted permissions).

## Admin Kubeconfig

The `status.adminKubeconfig` field contains the admin kubeconfig which is used to access the deployed Landscaper (full admin permissions).

## Shoot Name

The `status.shootName` field contains the name of the corresponding shoot resource.

## Shoot Namespace

The `status.shootNamespace` is the namespace in which the shoot resource is created.
