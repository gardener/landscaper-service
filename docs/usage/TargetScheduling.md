<!--
SPDX-FileCopyrightText: 2024 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# Target Scheduling

Target scheduling is the process of assigning with each [LandscaperDeployment][1] a [ServiceTargetConfig][2].
The ServiceTargetConfig defines on which target cluster the Landscaper instance will be deployed.
The target scheduling is done once for each new LandscaperDeployment by the LandscaperDeployment controller.

There is a [default scheduling](#default-scheduling) algorithm based on a simple priority system.
This requires no further configuration.

When onboarding large customers, the service operator might want to distribute the load to a defined set of
target clusters instead of evenly distribute the load across all target clusters. This can be achieved with the
[advanced scheduling](#advanced-scheduling) algorithm.


## Definitions

A ServiceTargetConfig is called **visible** if it has the annotation 
`config.landscaper-service.gardener.cloud/visible: true`. Otherwise, it is called **invisible**.
Invisible ServiceTargetConfigs are excluded from the scheduling, just as if they would not exist.

A ServiceTargetConfig is called **restricted** if its boolean field `spec.restricted` has the value `true`.
Otherwise, it is called **unrestricted**. By default, a ServiceTargetConfig is unrestricted. 
Restricted ServiceTargetConfigs can only be assigned by rules of the advanced scheduling.


## Default Scheduling

When scheduling a new LandscaperDeployment, the LandscaperDeployment controller lists all visible and unrestricted 
ServiceTargetConfigs. The priority (field `spec.priority`) of each ServiceTargetConfig is divided by the number 
of already deployed LandscaperDeployments + 1. All ServiceTargetConfigs are then sorted in descending order of the 
computation result. The ServiceTargetConfig with the highest number is selected.

## Advanced Scheduling

Advanced scheduling is defined in a new custom resource called `targetscheduling.landscaper-service.gardener.cloud`.
For advanced scheduling, a single instance of this custom resource must be created 
on the LaaS core cluster, in the LaaS system namespace (usually `laas-system`), with name `scheduling`.

The `TargetScheduling` resource can be created and modified without the need of restarting the Landscaper Service
controller. The changes in the `Scheduling` resource are applied during the next scheduling of a Landscaper Deployment.

**Remark:** The name and namespace of the TargetScheduling can be found in the config secret of the LaaS,
which contains the LandscaperServiceConfiguration:
```yaml
apiVersion: config.landscaper-service.gardener.cloud/v1alpha1
kind: LandscaperServiceConfiguration
...
scheduling:
  name: scheduling
  namespace: laas-system
```

### Structure of the TargetScheduling Resource

The spec of the TargetScheduling resource has a list of rules. Each rule has:

- a priority,  
- a list of ServiceTargetConfigs,  
- and a [selector](#selectors-matching-landscaperdeployments),
  which is a list of [terms](#terms-matching-landscaperdeployments).  

Therefore, the general structure is the following:

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: Scheduling

metadata:
  name: scheduling
  namespace: laas-system
  
spec:
  rules:
    - priority: 10
      serviceTargetConfigs:
        - name:      # name of a ServiceTargetConfig
          namespace: # namespace of a ServiceTargetConfig
        ...
      selector:
        - # term 1
        - # term 2
        ...
    
    ... # more rules 
```


### Choosing a ServiceTargetConfig for a LandscaperDeployment

[Selectors](#selectors) and [terms](#terms) either **match** a LandscaperDeployment or not, 
i.e. they describe properties that a LandscaperDeployment either has or does not have.

When distributing a LandscaperDeployment, we consider all rules whose selector matches the LandscaperDeployment.
Of these rules, we restrict ourselves to those with the highest priority. (It might be more than one, because different 
rules can have the same priority.) Finally, we choose one of the ServiceTargetConfigs from these rules. This is done
according to the same ranking as in the default scheduling.

It can happen that no rule applies, or that none of the ServiceTargetConfigs in the found rules exist.
In this case, we fall back to the default scheduling algorithm.


### Terms

There are different types of terms:

- terms that match a LandscaperDeployment with a certain tenant ID:
  ```yaml
  matchTenant:
    id: tenant0001
  ```
- terms that match LandscaperDeployments with a certain label:
  ```yaml
  matchLabel:
    name: workspace  # label name
    value: dev       # label value
  ```
- terms that combine a list of other terms with a logical "or":
  ```yaml
  or:
    - # term 1
    - # term 2
    ...
  ```
- terms that combine a list of other terms with a logical "and":
  ```yaml
  and:
    - # term 1
    - # term 2
    ...
  ```
- terms that negate another term:
  ```yaml
  not:
    # some term
  ```

### Selectors

A selector matches a LandscaperDeployment if all its terms match the LandscaperDeployment. This means, a selector
consisting of several terms like this:

```yaml
selector:
  - # term 1
  - # term 2
  - # term 3
```

is equivalent with:

```yaml
selector:
  - and:
      - # term 1
      - # term 2
      - # term 3
```


### Example

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: Scheduling

metadata:
  name: scheduling
  namespace: laas-system
  
spec:
  rules:
    - priority: 10
      serviceTargetConfigs:
        - name: dev-target-01
          namespace: laas-system
        - name: dev-target-02
          namespace: laas-system
      selector:
        - matchTenant:
            id: tenant0001
        - or:
            - matchLabel:
                name: workspace
                value: dev
            - matchLabel:
                name: workspace
                value: staging
        - matchLabel:
            name: project
            value: product-a
    
    - priority: 5
      serviceTargetConfig:
        - name: dev-target-06
          namespace: laas-system
      selector:
        - not:
            matchLabel:
              name: cost-center
              value: 123456
```


<!-- References -->

[1]: ./LandscaperDeployments.md
[2]: ./ServiceTargetConfigs.md
