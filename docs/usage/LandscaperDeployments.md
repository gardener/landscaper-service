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
    
  highAvailabilityConfig:
    controlPlaneFailureTolerance: "zone"
      
status:
  instanceRef:
    name: test
    namespace: my-namespace

  phase: Succeeded
  dataPlaneType: Internal
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

## High Availability Config

With this optional field the high availability mode of the Landscaper resource cluster of a Landscaper instance is configured.
Allowed values: `node`, `zone`.
For more details please check [this documentation](https://gardener.cloud/docs/gardener/high-availability/#shoot-clusters)

## DataPlane

The LandscaperDeployment can be created with an external data plane reference.
The reference can be either specified as an inline configuration or a Kubernetes secret reference.
This field can't be combined with `spec.oidcConfig` and `spec.highAvailabilityConfig`.
For more details please check [this documentation](../gettingstarted/create-landscaper-deployment.md#create-a-landscaper-deployment-with-an-external-data-plane).

```yaml
spec:
  dataPlane:
    kubeconfig: |
      apiVersion: v1
      kind: Config
      ...
```

```yaml
spec:
  dataPlane:
    secretRef:
      name: dataplane
      namespace: test
      key: kubeconfig
```

## Instance Reference

The `status.instanceRef` field will be set by the landscaper service controller when the Instance for the LandscaperDeployment has been created.

## Phase

The `status.phase` field mirrors the phase of the corresponding Landscaper Installation.

## DataPlaneType

The `status.dataPlaneType` shows the user whether an internal resource Shoot cluster is used (_Internal_) or an external data plane is used (_External_).
