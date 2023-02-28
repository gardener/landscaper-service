# Add-On Support

This document contains a proposal on how to support add-ons to user Landscaper Instances (LI).

## Motivation

A LI provides a user with the Landscaper functionalities. This might be sufficient for some or even most of the users.
For extended use-cases which are not part of the Landscaper, additional functionality could be provided by so called
Landscaper Instance Add-Ons (LIA).
A LIA can be of different forms. For example:

* A set of Kubernetes controllers, providing new APIs for the user to use or extending existing APIs.
* A set of configurations applied to a LI.
* A set of policies applied to a LI.

A LI user can choose a set of LIA that shall be installed and updated alongside the LI.
To give the LI user the ability to choose and configure LIAs, there must be an extension of the LI provisioning API.

## Requirements

This chapter describes the requirements for LIAs.

- Kubernetes Compatible:
  - All LIAs must be deployable to a Kubernetes cluster target.
- Defined/Described in OCM:
  - A LIA must bring its component description in the OCM component descriptor format. In the OCM format it can define its resources and
    dependencies to other components.
- Configuration described as a Landscaper Blueprint:
  - The deployment configuration of a LIA must be described and provided as a Landscaper Blueprint resource in the component descriptor.
- Self-Contained:
  - A LIA component must describe all its required dependencies. This means, no implicit requirements of functionality
    that is not provided by the LIA component descriptor itself.
- Requirements-Description:
  - A LIA must describe its requirements of the following in a standardized format:
    - Supported Kubernetes versions
    - Supported Landscaper versions

## Add-On Types

Even when the functionality provided by LIAs can be of any form, in the perspective of the LaaS landscape management,
there are two distinct types of LIAs.

1. LIAs running on resource-shoot-clusters:
   Tenant isolation is achieved by deploying the LIA directly on the user resource-shoot-cluster. 
   A LIA of this type is provided the Landscaper Target pointing to the resource-shoot-cluster.
2. LIAs running on target shoot clusters:
   Tenant isolation is achieved by deploying the LIA into the tenant namespace int the target-shoot-cluster.
   A LIA of this type is provided a Landscaper Target pointing to the target-shoot-cluster, a Landscaper Target pointing to
   the resource-shoot-cluster and the name of the tenant namespace in the target-shoot-cluster.

## LIAs running on resource-shoot-clusters

LIAs that are running on resource-shoot-cluster are being configured at the Landscaper Deployment custom resource for creating a LI.
Multiple LIAs can be configured for one Landscaper Deployment custom resource. Each LIA which is being added to the Landscaper Deployment custom resource
must provide a configuration.
For the configuration of the LIA, a Landscaper Installation template is specified in a Kubernetes configmap or secret. 
This Installation template is then referenced in the Landscaper Deployment custom resource.
The Installation template contains also the component descriptor reference, the component version and the Blueprint resources that is being invoked
for installing the LIA.
The LaaS controller then creates an instance of the Installation template alongside the LI Installation.
In the Installation template, the Landscaper Target, pointing to the resource-shoot-cluster, is being injected by the LaaS controller.

All LIAs are being installed after the LI installation has succeeded. However, the ordering in which the LIAs are being installed,
is not guaranteed.

Because all LIAs are specifying the compatible Kubernetes and Landscaper versions, the LaaS controller can verify if the LIA
is compatible with the Landscaper and Kubernetes version deployed with the LI. If the combination is not compatible, the LaaS controller
issues a error and refuses to install this combination of LI and LIAs.

When a LIA installs deployments/pods on the resource-shoot-cluster, the LIA must ensure themselves that required Service Accounts
with associated Roles and RoleCollections are being created.

### Important Considerations

* LIA providers must ensure that they are providing updates that provide a clear upgrade-path from the current set of supported
Landscaper and Kubernetes versions to newer released versions. Otherwise, a LIA might prevent updating a LI.
* When a LIA deploys pods onto the resource-shoot-clusters, Node resources are being consumed. The Node groups should be configured in a way, that
the resource-shoot-cluster can schedule new Nodes to provide the required resources for LIAs. However, it might be impossible to install any possible number
and combinations of LIAs.

## LIAs running on resource-shoot-clusters

TBD