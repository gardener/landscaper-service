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
  landscaperConfiguration:
    deployers:
      - helm
      - manifest
      - container
  componentReference:
    context: mycontext
    componentName: github.com/gardener/landscaper-service
    version: v0.16.0
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

  clusterEndpoint: "10.0.0.1:1234"
  clusterKubeconfig: "a3ViZWNvbmZpZyBjb250ZW50 ..."
```

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

## Cluster Endpoint

The `status.clusterEndpoint` field contains the API endpoint of the deployed Landscaper instance, i.e. it is used to
create and maintain Landscaper [Installations](https://github.com/gardener/landscaper/blob/master/docs/usage/Installations.md), 
which are handled by the corresponding Landscaper instance.

## Cluster Kubeconfig

The `status.clusterKubeconfig` field contains the kubeconfig which is used to access the deployed Landscaper.
