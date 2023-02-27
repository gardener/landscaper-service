# Multi-Version-Landscaper-Instances

This document contains a proposal how to support multiple versions of Landscaper Instances in one Landscaper as a 
Service (LaaS) landscape.

## Motivation

A LaaS landscape should support different versions of a Landscaper instance (LI), i.e. different LandscaperDeployments
install different versions of a [this](../../.landscaper/landscaper-instance) component. A LaaS landscape should support 
a predefined set of LI versions, which change from release to release (probably if the LaaS component version changes)
and a customer could select when to upgrade which of his instances.

A LI is a complex software artefact used by the customers for their deployments. Upgrading a LI might include complex
things like updating the used Helm version, the kubernetes version of the Landscaper resource shoot cluster and many
more. Due to this complexity, it is impossible to ensure that LI updates might not harm an existing customer deployments.

To reduce this risk, customers must have possibility to test a new LI component version with their development setup 
first, before they decide to upgrade the LI component of their productive environment.

## Requirements

This chapter defines the rules for supported LI versions in a LaaS landscape:

- Versioning:
  - Patch: A new patch version of a LI should only include bugfixes and security fixes. It should not break a customer
    setup.
  - Minor: A new minor version of a LI could contain additional features, minor updates of the most important components,
    like Helm or the Garden shoot cluster. Due to its complexity a minor update might break a customer setup. 
  - Major: Groundbreaking changes. Might break a customer setup.

- Supported LI version of a LaaS landscape:
  - The list of supported versions of a LaaS landscape only contains one patch version (usually the latest) for every 
    supported minor version. 

- Basic upgrading rules
  - A LI instance could be upgraded to the currently supported patch level of its current minor version. 
  - If a LI instance runs on a currently supported patch level or its current minor version is not supported anymore,
    then: 
    - it could be upgraded to the next supported minor version or 
    - if there is no supported higher minor version it could be upgraded to the lowest supported minor version of the 
      next major version.

- An upgrade of a LaaS landscape updates the list of supported LI versions according to the following rules:
  - A supported minor version is set on deprecated, if there is a supported version with a higher major version or
    a supported version with the same major but a higher minor version.
  - A minor version could only be removed if it was deprecated for more than 3 months.
  - It is allowed to replace a minor version with one with a higher patch level but not with a lower one. Such a change 
    will not affect if the minor version was deprecated or not.
  - It is not allowed to add a new supported version with a lower minor level than the highest supported minor level of 
    the same major level.
  - It is not allowed to remove a minor version which is the next version of another supported minor version according
    to the update rules above.

- Upgrading LIs
  - Every LI in a LaaS landscape is automatically upgraded to the currently supported patch version of its current
    minor version.
  - A customer could decide if his LI instance is upgraded automatically to higher minor versions or if this should be
    triggered manually.
  - Upgrades to the next major version must be triggered manually.
  - Customers are automatically informed about newly supported and deprecated LI component versions. 
  - Customers are automatically informed when a LI version will be removed from the supported version list.
  - If a customer LI instance is running on an unsupported minor version, it is automatically upgraded,
    according to upgrading rules described above. The customer is automatically informed about this.

- A customer must be able 
  - to select its intended LI version during onboarding. This must be optional. He gets the latest version
    if nothing is selected.
  - to see
    - the current version of his LI
    - the supported versions
    - the deprecated supported versions
    - the exact timeframes when deprecated LI versions are removed (-> forced update)
    - the upgrade path according to the rules above

## Questions:

- If we release new minor versions too frequently, the customer must test and upgrade quite often.

## Potential architectural consequences

- The Central Landscaper should always be on the latest Landscaper version. The version referenced by the
  LaaS project.

- Do we need to decouple the release of the LaaS and the LI component. Otherwise, the frequency of new
  of LI minor releases is higher than required resulting in more upgrades for the customer. 

- LandscaperDeployments, Instances require the currently deployed version of LI.

- LaaS component version requires list of supported versions, allowed migration paths for LI component.

- We need some object in the ls-user namespace with information about currenty version, supported versions, 
  migrations paths. User upgrades by creating some upgrade object. This is recognized by some controller running in the 
  same namespace as the sidecar controllers, which creates itself an upgrade object in his namespace. Subsequently, a 
  new controller (one for every Target-Shoot-Cluster) of the LaaS watches these objects and updates the 
  LanscaperDeployments accordingly.  

## Important considerations

- Different LI versions needs to be tested especially the Landscaper-Shoot-Version-Combination. 

- Every allowed LI upgrade must be tested. LI upgrade must take into consideration, that only particular kubernetes 
  version upgrades in shoots are allowed.