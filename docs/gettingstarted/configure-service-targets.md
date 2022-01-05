# Configure ServiceTargetConfigs

## Kubeconfig Secret

To create a new [ServiceTargetConfig](../usage/ServiceTargetConfigs.md) a kubernetes secret of type `Opaque` needs to be created.
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
The label `config.landscaper-service.gardener.cloud/region` is optional and can be used by LandscaperDeployments for selecting a target cluster with a specific geo-region.

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
    config.landscaper-service.gardener.cloud/region: "eu"

spec:
  providerType: gcp
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

NAME      PROVIDERTYPE   REGION   VISIBLE   PRIORITY   AGE
default   gcp            eu       true      10         1h10m
```

## Managing Visibility

A ServiceTargetConfig can be set to invisible state. When invisible, no new Landscaper deployments can be scheduled on the referenced kubernetes target cluster.
To set a ServiceTargetConfig to invisible, do the following:

```sh
kubectl -n laas-system label --overwrite=true servicetargetconfigs.landscaper-service.gardener.cloud default config.landscaper-service.gardener.cloud/visible=false
```

To set a ServiceTargetConfig to visible, do the following:

```sh
kubectl -n laas-system label --overwrite=true servicetargetconfigs.landscaper-service.gardener.cloud default config.landscaper-service.gardener.cloud/visible=true
```
