<!--
SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# Configure ServiceTargetConfigs

In this Chapter we deploy a [ServiceTargetConfig](../usage/ServiceTargetConfigs.md) manifest, which defines a kubernetes
cluster into which the Landscaper as a Service could deploy Landscaper instances.

## Kubeconfig Secret

To create a new ServiceTargetConfig a kubernetes secret of type `Opaque` needs to be created.
This secret needs have a key that contains the kubeconfig of the target cluster that the ServiceTargetConfig represents.
The secret can be created in any namespace. It is however recommended, to create it in a namespace that is only accessible to landscaper service administrators.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: default-target
  namespace: laas-system
type: Opaque
stringData:
  kubeconfig: |
    apiVersion: v1
    kind: Config
    ...
```

```sh
kubectl apply -f secret.yaml
```

## ServiceTargetConfig

The secret that was just created, is then referenced in a ServiceTargetConfig.
The `spec.secretRef` field contains the name, the namespace of the secret and the key that contains the kubeconfig.
The visibility label `config.landscaper-service.gardener.cloud/visible` needs to be set with the value `"true"`.

:warning: Attention: It is important that the `spec.providerType` matches the infrastructure provider type of the targeted kubernetes cluster.

The ServiceTargetConfig can be created in any namespace. It is however recommended, to create it in a namespace that is only accessible to landscaper service administrators.

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: ServiceTargetConfig

metadata:
  name: default
  namespace: laas-system
  labels:
    config.landscaper-service.gardener.cloud/visible: "true"

spec:
  priority: 10

  secretRef:
    name: default-target
    namespace: laas-system
    key: kubeconfig
```

```sh
kubectl apply -f servicetargetconfig.yaml
```

```sh
kubectl -n laas-system get servicetargetconfigs.landscaper-service.gardener.cloud

NAME      VISIBLE   PRIORITY   AGE
default   true      10         1h10m
```
