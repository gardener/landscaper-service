# Advanced Scheduling Proposal

Currently, the scheduling of Landscaper deployments to target clusters is based on a simple priority system.
When scheduling a new Landscaper deployment, the Landscaper Service Controller lists all _visible_ ServiceTargetConfigs (each one representing a target cluster).
The priority noted of each ServiceTargetConfig is divided by the number of already deployed Landscaper deployments + 1.
All ServiceTargetConfigs are the sorted in descending order of the computation result. The ServiceTargetConfig with the highest number is selected.

When onboarding large customers this scheduling algorithm is too limited. For large customers the service operator would like to distribute the load to a defined set of target clusters
instead of evenly distribute the load across all target clusters.

Therefore, an advanced scheduling configuration/algorithm is needed. 
In this proposal there are two approaches for such a configuration. A more complex one which offers a high amount of flexibility in terms of possible configurations
and a simpler one, easier to read by humans and easier to implement.

## Configuring the LAAS Controller

Both approaches are defined in a new custom resource definition, called `Scheduling.landscaper-service.grdener.cloud`. 
This resource is created by the Landscaper Service operator and then referenced in the `LandscaperServiceConfiguration`.
This means, only a single `Scheduling` resource is used by the Landscaper Service controller.
The `Scheduling` resource can be modified without the need of restarting the Landscaper Service controller.
The changes in the `Scheduling` resource are applied during the next scheduling of a Landscaper Deployment.

```yaml
apiVersion: config.landscaper-service.gardener.cloud/v1alpha1
kind: LandscaperServiceConfiguration

schedulingRef:
  name: default
  namespace: laas-system
```

## Complex / Maximum Flexibility

The `Scheduling` resource has a list of rules.
Each rule has a `priority`, a list of ServiceTargetConfigs, `serviceTargetConfigs` and a `selector`.
The selector is a list of terms that are applied to a Landscaper Deployment. A term can match the tenant id of a Landscaper Deployment
and its labels. Additionally, a term can also be a logical `or`, `and` and `not`.
When the selector matches a Landscaper Deployment it is distributed to one of the ServiceTargetConfigs specified in the `servieTargetConfigs` list.
If two or more selectors match, the one with the highest priority wins. If two or more with the same priority match, any of them matches.

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: Scheduling

metadata:
  name: default
  namespace: laas-system
  
rules:
- priority: 10
  serviceTargetConfigs:
    - name: dev-target-01
      namespace: laas-system
    - name: dev-target-02
      namespace: laas-system
  selector:
    - tenantId: tenant0001
    - or:
      - label:
        name: workspace
        value: dev
      - label:
        name: workspace
        value: staging
    - label:
      name: project
      value: product-a

- priority: 5
  serviceTargetConfig:
    name: dev-target-06
    namespace: laas-system
  selector:
    - not:
      - label:
        name: cost-center
        value: 123456
```

## Simple / Less Flexibility

The `Scheduling` resource has a list of rules.
Each rule specifies one tenant id and a list of ServiceTargetConfigs.
When the tenant id of a Landscaper Deployment matches the tenant id of a rule, the deployment is distributed to one of
the ServiceTargetConfigs specified in the list.

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: Scheduling

metadata:
  name: default
  namespace: laas-system
  
rules:
  tenant0001:
    serviceTargetConfigs:
      - name: dev-target-01
        namespace: laas-system
  tenant0002:
    serviceTargetConfigs:
      - name: dev-target-02
        namespace: laas-system
      - name: dev-target-03
        namespace: laas-system
```

## Restrict Scheduling on Service Targets

Service Target configs can be labelled wit a scheduling policy `config.landscaper-service.gardener.cloud/schedulingPolicy`.
If the label value is set to `all`, all LandscaperDeployments can be scheduled on this target even if it hasn't been explicitly selected
by a Scheduling rule.
If the label value is set to `restricted`, LandscaperDeployments can only be scheduled on this target if it has been selected explicitly
by a Scheduling rule.

The default value is `all`.

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: ServiceTargetConfig

metadata:
  name: dev-target-10
  labels:
    config.landscaper-service.gardener.cloud/visible: "true"
    config.landscaper-service.gardener.cloud/schedulingPolicy: "all|resctricted"  
spec:
  priority: 10

  ingressDomain: ingress.mydomain.net

  secretRef:
    name: dev-target-10
    namespace: laas-system
    key: kubeconfig
```