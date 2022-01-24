<!--
SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# Installation and configuration of the Landscaper Service

This document describes the installation of the landscaper service.
Landscaper service is a Kubernetes controller that reconciles [landscaper](https://github.com/gardener/landscaper) deployments as a service.

## Prerequisites

For the landscaper service to work, a landscaper needs to be deployed on the same cluster as the landscaper service.
Please consult the [landscaper documentation](https://github.com/gardener/landscaper/tree/master/docs) on how to install the landscaper.


## Installation

The landscaper service can be installed via a landscaper installation using the _landscaper-service-blueprint_ blueprint of the _github.com/gardener/landscaper-service_ component.

```


### Installation blueprint

To install the landscaper service via a [landscaper installation](https://github.com/gardener/landscaper/blob/master/docs/usage/Installations.md), a [target resource](https://github.com/gardener/landscaper/blob/master/docs/technical/target_types.md) needs to be created.
The target defines the kubernetes cluster on which the landscaper service is installed.

```yaml
apiVersion: landscaper.gardener.cloud/v1alpha1
kind: Target
metadata:
  name: laas-target-cluster
spec:
  type: landscaper.gardener.cloud/kubernetes-cluster
  config:
    kubeconfig: |
      apiVersion: v1
      kind: Config
      ...
```

```sh
kubectl create namespace laas-system
kubectl apply -n laas-system -f target.yaml
```

This target is referenced by an installation resource, which specifies the landscaper service installation with its configuration.

```yaml
apiVersion: landscaper.gardener.cloud/v1alpha1
kind: Installation
metadata:
  name: landscaper-service
spec:
  componentDescriptor:
    ref:
      repositoryContext:
        type: ociRegistry
        baseUrl: eu.gcr.io/gardener-project/development
      componentName: github.com/gardener/landscaper-service
      version: v0.1.0

  blueprint:
    ref:
      resourceName: landscaper-service-blueprint

  imports:
    targets:
      - name: targetCluster
        target: "#laas-target-cluster"

  importDataMappings:
    namespace: laas-system
    verbosity: 2

    # optional: registry pull secrets, list of secrets in "kubernetes.io/dockerconfigjson" format
    # registryPullSecrets:
    #  - name: secret1
    #    namespace: laas-system
    #  - name: secret2
    #    namespace: laas-system
```

The specification of the `registryPullSecrets` is optional and is only needed when the landscaper service component can't be pulled anonymously.
The `registryPullSecrets` field contains a list of secrets referenced by name and namespace.
See [this documentation](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#registry-secret-existing-credentials) for the required secret format.

```sh
kubectl apply -n laas-system -f installation.yaml 
```

Once the installation has been created, the landscaper will perform the necessary steps to perform the deployment of the landscaper service.
The status of the installation can be inspected with the [landscaper-cli](https://github.com/gardener/landscapercli).

```sh
landscaper-cli installations inspect -n laas-system

[✅ Succeeded] Installation landscaper-service
    └── [✅ Succeeded] DeployItem landscaper-service-landscaper-service-2dv4x

```

## Update Process

The _github.com/gardener/landscaper-service_ component contains a component reference to the supported landscaper version.
When the landscaper service controller is updated, all currently deployed landscaper instances will be automatically updated to the new supported version.
