# Landscaper as a Service High Availability (HA) Concept

This document describes the HA concepts of a LaaS landscape. The
[architecture document](./architecture.md) is a prerequisite of the following description.

## Concepts

As described in the architecture, every customer has access to its Resource-Cluster on which he maintains k8s Custom 
Resources (Installations, Targets ...) and Secrets to specify his installation tasks. 

- For a productive scenario, the Resource-Cluster is set up with Control Plane HA which is a Gardener mechanism to
  guarantee HA of their clusters. 

- All components required for the direct customer interaction, i.e. for the maintenance of its k8s resources, are 
  setup such that at least 2 replicas are running on two different nodes in different zones. The relevant LaaS 
  components are the nginx-controller and the webhook backends on the Target-Clusters. The Target-Clusters itself
  are again set up with Control Plane HA. 

- The maintained k8s Custom Resources trigger the relevant actions as asynchronous background jobs. Every component
  required for this, is restarted within a few seconds, if it fails and can continue its work immediately. The relevant
  components are the landscaper and the different deployer.