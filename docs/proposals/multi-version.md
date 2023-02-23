# Multi-Version-Landscaper-Instances

This document contains a proposal how to support multiple versions of Landscaper Instances in one Landscaper as a 
Service (LaaS) landscape.

## Requirements

- A LaaS landscape should support different versions of a Landscaper instance (LI), i.e. different LandscaperDeployments 
  install different versions of a [this](../../.landscaper/landscaper-instance) component. A more detailed
  selection of particular deployer versions seems to be too complicated and out of scope.

- A LaaS landscape supports a predefined set of LI versions, which change from release to realease (probably if the LaaS 
  component version changes).

- The Central Landscaper should always be on the latest Landscaper version. The version referenced by the 
  LaaS project.  

- A customer must be able 
  - to select its intended LI version during onboarding. This must be optional. He gets the latest version
    if nothing is selected.
  - to upgrade to its intended LI version. Therefore, an allowed upgrade path must be enforced.
  - see which versions are currently supported, which are recommended and which are outdated and what are the possible
    upgrade paths.
  - to see the exact timeframes in which 
    - specific LI versions are supported
    - specific LI versions are set to "deprecated"
    - specific LI versions run out of support -> results in a forced-upgrade?

- How to upgrade a customer LI instance if his current version is not supported anymore? This might become even more 
  critical if a security issue came up.
  - Proposal for an initial solution: 
    - Not supported versions are not upgraded. Only the reconciliation of the corresponding Installation is retriggered
      for token rotation issues, or we enhance the Installation by some automatic reconcile mechanisms.
    - Affected customers are informed via email. This mail is send automatically and if the customer does not upgrade,
      we get informed to personally contact him.

- Do we need some kind of automatic upgrading, which could be selected ba a customer? Could it be selected that only
  patch or minor versions are automatically upgraded?

## Potential architectural consequences

- LanscaperDeployments, Instances require the currently deployed version of LI.

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