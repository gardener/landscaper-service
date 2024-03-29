<!--
SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# ServiceTargetConfigs

ServiceTargetConfigs are kubernetes resources that represent a target cluster on which a Landscaper as a Service can
deploy Landscaper instances.

Each ServiceTargetConfig has a reference to a Secret that contains the Kubeconfig of the target cluster.
ServiceTargetConfigs doesn't need to reside in the namespaces as [LandscaperDeployments](LandscaperDeployments.md) and [Instances](Instances.md) which are referencing it.
It is advised to create the ServiceTargetConfigs in a separate namespace that is only accessible by administrators.

### Basic structure:

````yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: ServiceTargetConfig

metadata:
  name: default
  labels:
    config.landscaper-service.gardener.cloud/visible: "true"

spec:
  priority: 10

  ingressDomain: ingress.mydomain.net

  secretRef:
    name: default-target
    namespace: laas-system
    key: kubeconfig

status:
  instanceRefs:
    - name: test
      namespace: my-namespace
````

## Labels

ServiceTargetConfig supported labels:

* `config.landscaper-service.gardener.cloud/visible` defines the visibility of the ServiceTargetConfig. 
When set to `true`, the ServiceTargetConfig can be used to schedule new deployments of the Landscaper on the referenced kubernetes cluster.
When not set or set to any other value than `true`, no new deployments can be scheduled on the referenced cluster.

### Managing Visibility

A ServiceTargetConfig can be set to invisible state. When invisible, no new Landscaper deployments can be scheduled on the referenced kubernetes target cluster.
To set a ServiceTargetConfig to invisible, do the following:

```sh
kubectl -n laas-system label --overwrite=true servicetargetconfigs.landscaper-service.gardener.cloud default config.landscaper-service.gardener.cloud/visible=false
```

To set a ServiceTargetConfig to visible, do the following:

```sh
kubectl -n laas-system label --overwrite=true servicetargetconfigs.landscaper-service.gardener.cloud default config.landscaper-service.gardener.cloud/visible=true
```

## Priority

The `spec.priority` field is an integer number specifying the scheduling priority for the ServiceTargetConfig. 
To calculate the effective priority when scheduling Instances is calculated by dividing the specified priority by the number of instance references
(`spec.priority/(len(status.instanceRefs) + 1)`).
The more instances that are referenced by a ServiceTargetConfig, the lower the effective priority becomes.

## Ingress Domain

The `spec.ingressDomain` field is a string specifying the ingress domain of the referenced target cluster.
The ingress domain of the target cluster has to be configured so that an ingress resource can be created as an endpoint for the Webhook Server of every Landscaper instance.

## Secret Reference

The `spec.secretRef` field references a kubernetes secret by a name, a namespace and the key within the secret.
The key must contain the kubeconfig for the kubernetes target cluster on which Landscaper deployments are scheduled.


## Instance References

The `status.instanceRefs` is a list containing references to all Instances using this ServiceTargetConfig.
