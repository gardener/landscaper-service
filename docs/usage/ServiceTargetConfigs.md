# ServiceTargetConfigs

ServiceTargetConfigs are kubernetes resources that represent a target cluster on which a Landscaper as a service can be deployed.
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
    config.landscaper-service.gardener.cloud/region: "eu-west-1"

spec:
  providerType: gcp | aws | alicloud
  priority: 10

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

ServiceTargetConfigs support two special labels:

* `config.landscaper-service.gardener.cloud/visible` defines the visibility of the ServiceTargetConfig. 
When set to `true`, the ServiceTargetConfig can be used to schedule new deployments of the Landscaper on the referenced kubernetes cluster.
When not set or set to any other value than `true`, no new deployments can be scheduled on the referenced cluster.
* `config.landscaper-service.gardener.cloud/region` can be used to specify the geo-region of the referenced kubernetes cluster.
LandscaperDeployments can make use of this label to specify on which geo-region the Landscaper shall be deployed.


## Provider Type

The field `spec.providerType` specifies the infrastructure provider of the referenced kubernetes cluster.
Currently, supported values are:

* `gcp`
* `aws`
* `alicloud`


## Priority

The `spec.priority` field is an integer number specifying the scheduling priority for the ServiceTargetConfig. 
To calculate the effective priority when scheduling Instances is calculated by dividing the specified priority by the number of instance references
(`spec.priority/(len(status.instanceRefs) + 1)`).
The more instances that are referenced by a ServiceTargetConfig, the lower the effective priority becomes.


## Secret Reference

The `spec.secretRef` field references a kubernetes secret by a name, a namespace and the key within the secret.
The key must contain the kubeconfig for the kubernetes target cluster on which Landscaper deployments are scheduled.


## Instance References

The `status.instanceRefs` is a list containing references to all Instances using this ServiceTargetConfig.
