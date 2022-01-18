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

The landscaper service can be installed via [Helm](https://helm.sh) using the Helm chart [charts/landscaper-service](../../charts/landscaper-service)
or via a landscaper installation using the _landscaper-service-blueprint_ blueprint of the _github.com/gardener/landscaper-service_ component.


### Helm

To install the landscaper service via Helm, the chart has to be pulled from the landscaper service OCI registry.

:warning: Attention: Helm version v3.7.0 or later is required for this to work.

```sh
export HELM_EXPERIMENTAL_OCI=1
export LAAS_VERSION="v0.1.0" # use the latest available version

kubectl create namespace laas-system
helm pull oci://eu.gcr.io/gardener-project/landscaper-service/charts/landscaper-service --version $LAAS_VERSION
helm install -n laas-system landscaper-service ./landscaper-service-${LAAS_VERSION}.tgz
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
```

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
