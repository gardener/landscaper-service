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
  - Patch: A new patch version of a LI should only include bugfixes and security fixes. It must not break a customer
    setup.
  - Minor: A new minor version of a LI could contain additional features, minor updates of the most important components,
    like Helm or the Garden shoot cluster. Due to its complexity a minor update might break a customer setup. 
  - Major: Groundbreaking changes. Customer setup is very likely to be broken.

- Supported LI version of a LaaS landscape:
  - The list of supported versions of LI of a LaaS landscape is a list of LI versions whereby for every combination of
    major and minor version there is at most one entry. 
    - Examples:
      - v0.1.15, v0.3.0, v1.0.0, v3.2.3, v5.2.3: valid list
      - v0.1.15, v0.3.0, v0.3.3, v1.0.0, v3.2.3: invalid list because there are two patch versions of v0.3.* listed
    - The goal of this is that only the latest patch versions of a minor version is supported on a LaaS landscape.

- Basic upgrading rules
  - A LI instance is allowed to be upgraded to the currently supported patch level of its current minor version. 
  - If a LI instance runs on a currently supported version or its current minor version is not supported anymore,
    then: 
    - it is allowed to be upgraded to the next supported minor version or 
    - if there is no supported higher minor version it is allowed to be upgraded to the lowest supported minor version 
      of the next major version.

- An upgrade of a LaaS landscape updates the list of supported LI versions according to the following rules:
  - A supported minor version is set on deprecated, if there is a supported version with a higher major version or
    a supported version with the same major but a higher minor version.
  - A supported minor version can only be removed if the deprecation duration has ended.
  - It is allowed to replace a minor version with one with a higher patch level but not with a lower one. If the minor 
    version was already deprecated before, this holds also for the minor version with the higher patch level and also the
    deprecation duration is not affected.
  - It is not allowed to add a new supported version with a lower minor level than the highest supported minor level of 
    the same major level.
  - It is not allowed to remove a minor version which is the next version of another supported minor version according
    to the update rules above.
  - Deprecation duration: The deprecation duration is currently 3 months.

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

- The version of the Central Landscaper is completely independent of the supported LI versions and their internal 
  Landscaper versions.

- Do we need to decouple the release of the LaaS and the LI component? Otherwise, the frequency of new
  of LI minor releases is higher than required resulting in more upgrades for the customer. 

- LandscaperDeployments, Instances require the currently deployed version of LI.

- LaaS component version requires list of supported versions, allowed migration paths for LI component.

- We need some object in the ls-user namespace with information about currenty version, supported versions, 
  migrations paths. User upgrades by creating some upgrade object. This is recognized by some controller running in the 
  same namespace as the sidecar controllers, which creates itself an upgrade object in his namespace. Subsequently, a 
  new controller (one for every Target-Shoot-Cluster) of the LaaS watches these objects and updates the 
  LanscaperDeployments accordingly.  

### Project Setup

The critical component with respect to the support of multiple version is the LI component. The more different 
minor (and of course major) versions of this component exits, the more critical upgrades must be done by customers. 
Therefore, it should be as small as possible to reduce the number of potential minor releases.

- To reduce the complexity of the LI component, the 
[sidecar component](../../.landscaper/landscaper-instance/blueprint/installation/sidecar-subinst.yaml) should be removed.
The sidecar component is under our control and new features should not result in a new version of the LI component.

- There is no reason for the LI component to be coupled with the releases of the LaaS project especially if the sidecar 
  component is removed from it. Decoupling the LI component from the LaaS releases means that minor releases of the LaaS 
  do not result in minor releases of the LI.  Therefore, it is a natural step that the LI component should not be part 
  of the LaaS project.

- The LI component could be either part of the Landscaper project of get its own project. Making the LI component part 
  of the Landscaper means that every PR/release must not only execute the Landscaper integration tests but also LI
  integration tests. The LI integration tests are very time-consuming because they require the creation of a garden shoot 
  cluster. Therefore, it might be a good idea to decouple LI and Landscaper in different projects. 

- Testing:
  - The LI component integration tests must test a minimal feature creating a new shoot cluster 
  - The LaaS integration tests must include: 
    - check if upgrades to patch version are possible and still working
    - upgrade to next minor versions are possible and still working

## Important considerations

- Different LI versions needs to be tested especially the Landscaper-Shoot-Version-Combination. 

- Every allowed LI upgrade must be tested. LI upgrade must take into consideration, that only particular kubernetes 
  version upgrades in shoots are allowed.