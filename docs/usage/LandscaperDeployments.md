<!--
SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# LandscaperDeployments

LandscaperDeployments are kubernetes resources, created by clients/users, to trigger the deployment of a landscaper instance by 
the Landscaper as a Service. The Landscaper as a Service controller will select a suitable [ServiceTargetConfig](ServiceTargetConfigs.md) 
and creates an [Instance](Instances.md) for the LandscaperDeployment.

### Basic structure:

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: LandscaperDeployment
metadata:
  name: test
  namespace: my-namespace
spec:
  purpose: "test"
  region: "eu-west-1"
  landscaperConfiguration:
    deployers:
      - helm
      - manifest
      - container

status:
  instanceRef:
    name: test
    namespace: my-namespace
```

## Purpose

The `spec.purpose` field should contain the human-readable purpose of the LandscaperDeployment. 


## Region

The `spec.region` field is optional and can be used to restrict the selected target kubernetes cluster to a specific geo-region.


## Landscaper Configuration

The `spec.landscaperConfiguration` field contains the configuration of the Landscaper deployment.
Configuration contains the list of the standard deployers that shall be deployed.
For available deployers please check [this documentation](https://github.com/gardener/landscaper/tree/master/docs/deployer).

## Instance Reference

The `status.instanceRef` field will be set by the landscaper service controller when the Instance for the LandscaperDeployment has been created.
