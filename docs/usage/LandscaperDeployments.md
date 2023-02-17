<!--
SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# LandscaperDeployments

[LandscaperDeployments](../../pkg/apis/core/type_landscaperdeployment.go) are kubernetes resources, created by 
clients/users, to trigger the deployment of a landscaper instance by the Landscaper as a Service. The Landscaper as a 
Service controller will select a suitable [ServiceTargetConfig](ServiceTargetConfigs.md) 
and creates an [Instance](Instances.md) for the LandscaperDeployment.

### Basic structure:

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: LandscaperDeployment
metadata:
  name: test
  namespace: my-namespace
spec:
  tenantId: "tenant-1"
  purpose: "test"
  landscaperConfiguration:
    deployers:
      - helm
      - manifest
      - container
  oidcConfig: # optional
    clientID: <some client ID>
    issuerURL: <OIDC token issuer URL>
    groupsClaim: groups
    usernameClaim: email

status:
  instanceRef:
    name: test
    namespace: my-namespace
```

## TenantId

The `spec.tenantId` field has to contain the globally unique identifier of the owning tenant.

## Purpose

The `spec.purpose` field should contain the human-readable purpose of the LandscaperDeployment.

## Landscaper Configuration

The `spec.landscaperConfiguration` field contains the configuration of the Landscaper deployment.
Configuration contains the list of the standard deployers that shall be deployed.
For available deployers please check [this documentation](https://github.com/gardener/landscaper/tree/master/docs/deployer).

## Oidc Config

With the optional field OIDC config you specify that the Landscaper resource cluster of a Landscaper instance
(the Gardener shoot cluster on which the user creates Installations, Targets etc,) gets this OIDC configuration such
that user access could be provided via OIDC.

## Instance Reference

The `status.instanceRef` field will be set by the landscaper service controller when the Instance for the LandscaperDeployment has been created.
