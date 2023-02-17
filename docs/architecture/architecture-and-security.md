# Landscaper as a Service (LaaS) Architecture and Security Concepts

This document describes the architecture and security concepts of a LaaS landscape.

## 1 Architecture

### Architecture Overview

The overall architecture of one LaaS landscape is depicted in the following picture:

![architecture](images/architecture-01.png)

The main units are:

- Secret Store: The secret store contains the credentials for the pipeline, e.g. to access the Gardener project
  namespaces.

- Garden-LaaS-Project: This is a Garden project where the shoot clusters with the central units of a LaaS landscape
  are running

- Core-Shoot-Cluster: A Garden shoot cluster in the Garden-LaaS-Project. On this shoot cluster the Central Landscaper 
  and the LaaS component/instance is running. The LaaS uses the Central Landscaper to deploy the Landscaper instances for 
  different customers.

- Target-Shoot-Cluster: A Garden shoot cluster where the different Landscaper instances for the customer are running 
  in different namespaces. In principle a LaaS landscape might have multiple Target-Shoot-Clusters.

- Deploy-Pipeline: The Deploy-Pipeline installs the Central Landscaper as a Helm chart on the Core-Shoot-Cluster. The 
  Central Landscaper watches the Core-Shoot-Cluster for Landscaper 
  [Installations](https://github.com/gardener/landscaper/blob/master/docs/usage/Installations.md) resources. The 
  Deploy-Pipeline creates an Installation resource in the Core-Shoot-Cluster to also install the LaaS component 
  in the Core-Shoot-Cluster. The LaaS component is responsible to install Landscaper instances for the customers 
  in the Target-Shoot-Cluster(s). 

- Gardener-Resource-Cluster-Project: The Gardener project for the custom resource shoot cluster. If a customer requests 
  a Landscaper instance, the LaaS, running on the Core-Shoot-Cluster, deploys a new Landscaper instance in a new namespace 
  of the Target-Shoot-Cluster. For every new Landscaper instance, the LaaS also starts a new Resource-Shoot-Cluster
  in the Gardener-Resource-Cluster-Project. For its own deployments, the customer could create Landscaper 
  Installation resources on this Resource-Shoot-Cluster. The new Landscaper instance on the Target-Shoot-Cluster is 
  watching and processing these.

- Resource-Shoot-Cluster: Every Landscaper instance for a customer (running on the Target-Shoot-Cluster) watches 
  a different Resource-Shoot-Cluster, which is a Gardener shoot cluster in the Gardener-Resource-Cluster-Project, where 
  the customer could deploy its Landscaper Installation resources.

### Core-Shoot-Cluster and Target-Shoot-Cluster(s)

The following image gives a more detailed impression of the Core-Shoot-Cluster and the Target-Shoot-Cluster(s) and their
connections.

![architecture](images/core-and-target-cluster.png)

The Deploy-Pipeline installs the Central Landscaper on the Core-Shoot-Cluster as a helm chart deployment. The central 
Landscaper watches the Core-Shoot-Cluster for Landscaper Installations. Since the host and resource cluster of the 
Core Landscaper are equal, we can install it by deploying the 
[bundled `landscaper` helm chart](https://github.com/gardener/landscaper/tree/master/charts/landscaper) 
with its two sub charts `landscaper-rbac` and `landscaper-controller`:
- The `landscaper-rbac` helm chart deploys the ServiceAccounts used by the Landscaper controller internally, for
  watching Installations, and for creating the ValidatingWebhookConfig.
- The `landscaper-controller` chart deploys the landscaper controller and deployers.

Then the Deploy-Pipeline creates an Installation (and Target etc.) for the [LaaS component](../../.landscaper), such 
that the Central Landscaper installs it also on the Core-Shoot-Cluster. The LaaS component has the name 
*github.com/gardener/landscaper/landscaper-service* and comprises the following important parts:

- [blueprint](../../.landscaper/blueprint/blueprint.yaml)
- [resources](../../.landscaper/resources.yaml)
- [component-references](../../.landscaper/component-references.yaml)

To enable the LaaS to deploy e.g. Landscaper Instances for customers to the Target-Shoot-Clusters, a 
[ServiceTargetConfig](../usage/ServiceTargetConfigs.md) object is deployed on the Core-Shoot-Cluster for every 
Target-Shoot-Cluster. Every ServiceTargetConfig references a secret containing a kubeconfig providing access to one 
Target-Shoot-Cluster. The kubeconfig is created for a service account on the Target-Shoot-Cluster.

### Landscaper-Instance Creation

The following image gives a more detailed overview how Landscaper instances are created for a customer.

![architecture](images/landscaper-instance-creation.png)

For every customer there is a dedicated namespace on the Core-Shoot-Cluster. To get a Landscaper Instance a 
user/customer/operator creates a [LandscaperDeployment](../usage/LandscaperDeployments.md) custom resource in this 
namespace. The LaaS instance in the Core-Shoot-Cluster watches for such custom resources and creates an 
[Instance](../usage/Instances.md) custom resource for every LandscaperDeployment. For every Instance, the LaaS deploys an 
[Installation](https://github.com/gardener/landscaper/blob/master/docs/usage/Installations.md) for a Landscaper Instance.

The LaaS furthermore creates two [Targets](https://github.com/gardener/landscaper/blob/master/docs/usage/Targets.md)
as input for the Installation containing the credentials

- for a Target-Shoot-Cluster. These credentials are fetched from one of the ServiceTargetConfig custom resources and 
  are used to install all required artefacts for the Landscaper Instance on a Target-Shoot-Cluster.

- for the Garden project Garden-Resource-Cluster-Project to create a Resource-Shoot-Cluster on which the customer 
  can later deploy its Landscaper resources like Installations, Targets, etc.

The Central Landscaper watches for the Installations and executes them.

### Landscaper Instance Details

A Landscaper Instance is defined by [this](../../.landscaper/landscaper-instance) component. The component 
contains a [blueprint](../../.landscaper/landscaper-instance/blueprint/installation/blueprint.yaml) with five 
[sub installations](https://github.com/gardener/landscaper/blob/master/docs/usage/Blueprints.md#nested-installations), 
which are deploying:

- **[shoot](../../.landscaper/landscaper-instance/blueprint/installation/shoot-cluster-subinst.yaml)**: 
  A new shoot custom resource in the Garden-Resource-Cluster-Project for which the Gardener creates a new 
  Resource-Shoot-Cluster
- **[landscaper](../../.landscaper/landscaper-instance/blueprint/installation/landscaper-deployment-subinst.yaml)**: 
  A new Landscaper and its deployers in new namespace on one of the Target-Shoot-Clusters. This Landscaper is watching 
  and processing the Landscaper resources on the Resource-Shoot-Cluster.
- **[ls-service-target-shoot-sidecar-server](../../.landscaper/landscaper-instance/blueprint/installation/sidecar-subinst.yaml)**: 
  Two controllers handling the access rights of the customer/users to the Resource-Shoot-Cluster as well as the 
  creation of custom namespaces on that cluster.
- **[rbac](../../.landscaper/landscaper-instance/blueprint/installation/landscaper-rbac-subinst.yaml)**: A 
  service account with the right permissions, providing the new Landscaper access to the Resource-Shoot-Cluster.
- **[sidecar-rbac](../../.landscaper/landscaper-instance/blueprint/installation/sidecar-rbac-subinst.yaml)**: 
  A service account with the right permissions providing the controllers installed by 
  ls-service-target-shoot-sidecar-server access to the Resource-Shoot-Cluster.

The following picture gives an overview about the different sub installations:

![see also](images/landscaper-instance-details.png)

The installed Landscaper Instance watches the just created Resource-Shoot-Cluster on which the customer/user could 
deploy and maintain its Installations.

### Details of a Resource-Shoot-Cluster

Initially a Resource-Shoot-Cluster has two namespaces which are interesting with respect to the Landscaper Instance:

- namespace ls-system: This namespace mainly contains the Installations for the deployers of the Landscaper which 
  is part of the Landscaper-Instance. Customers/users have no access to this namespace. 

- namespace ls-user: This namespace allows users to control the access to the Resource-Shoot-Cluster as well as to 
  maintain the customer namespaces. It is created during the startup of the controller deployed by the sub installation
  ls-service-target-shoot-sidecar-server. 

#### Controlling user access to the Resource-Shoot-Cluster

In the namespace ls-users there exists one custom resource *subjects* of type 
[SubjectList](../../pkg/apis/core/v1alpha1/types_subjectlist.go). 

```bash
kubectl get subjectlists -n ls-user subjects
NAME       AGE
subjects   ...
```

At the beginning this custom resource looks as follows:

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: SubjectList
metadata:
  name: subjects
  namespace: ls-user
  ...
spec:
  subjects: []
```

To create an initial access for a user, a cluster administrator creates a service account in the namespace ls-user:

```bash
kubectl create sa -n ls-user <name of service account>
```

Next, the administrator adds an entry for the new service account to the `SubjectList` *subjects* such that it looks 
as follows:

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: SubjectList
metadata:
  name: subjects
  namespace: ls-user
  ...
spec:
  subjects:
  - kind: ServiceAccount
    name: <name of service account>
```

Now, the controllers deployed by ls-service-target-shoot-sidecar-server automatically add this service account to 
particular predefined role-bindings and cluster-role-bindings, with the exactly right permissions a user needs to 
work with a Resource-Shoot-Cluster, e.g. to create Landscaper Installations, Targets etc. 

Next, the cluster admin can fetch a token with a restricted lifetime for the new service account with the
[token request API](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/token-request-v1/#TokenRequest):

```bash
kubectl create token -n ls-user testserviceaccount --duration=48h

eyJhbGciOiJSUzI1NiIsImtpZCI6IkdzM1J...
```

A kubeconfig with this token allows a user to access the Resource-Shoot-Cluster and also to refresh its token when 
required. With this access data, the user can also create additional service accounts, add them to the `SubjectList` 
*subjects*. The user is not allowed to create or modify roles, cluster-roles, role-bindings or cluster-rolebindings to
get more permissions.

As described above, a Landscaper Instance is deployed when a [LandscaperDeployment](../usage/LandscaperDeployments.md)
is created. In a `LandscaperDeployment` you could also specify some OIDC configuration for the Resource-Shoot-Cluster, 
allowing end users to be authenticated via OIDC. To give such an authenticated user the required end user permissions,
you might add its email to the `SubjectList` *subjects*:

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: SubjectList
metadata:
  name: subjects
  namespace: ls-user
  ...
spec:
  subjects:
  - kind: User
    name: <email of user>
```

Again the controllers of ls-service-target-shoot-sidecar-server automatically add this user to the right
predefined role-bindings and cluster-role-bindings, in the same way as for the service accounts above. 

If the OIDC flow also provides group information, another possibility to authorize users is to add some of their groups
to the `SubjectList` *subjects*:

```bash
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: SubjectList
metadata:
  name: subjects
  namespace: ls-user
  ...
spec:
  subjects:
  ...
  - Kind: Group
    Name: <some group>
```

Again the controllers of ls-service-target-shoot-sidecar-server automatically add this group to the right
predefined role-bindings and cluster-role-bindings, in the same way as for the service accounts above.

The following image gives a more detailed descriptions of the involved roles, cluster-roles etc. The namespaces 
*cu-* are so-called customer namespaces on the Resource-Shoot-Cluster and will be described in more detail below:

![architecture](images/resource-cluster.png "Resoure-Shoot-Cluster")

#### Controlling Customer Namespaces on the Resource-Shoot-Cluster

A user, with access to the Resource-Shoot-Cluster as described before, is only allowed to create Landscaper resources 
like Installations, Targets etc. in so-called customer namespaces. A customer namespace is a normal namespace on the
Resource-Shoot-Cluster with a name starting with the prefix *cu-*. 

To create such a namespace the user must create a 
*[namespaceRegistration](../../pkg/apis/core/v1alpha1/types_namespaceregistration.go)* object in the namespace ls-user
with the same name as the namespace. The following manifest for example would create a customer namespace *cu-test*:

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: NamespaceRegistration
metadata:
  name: cu-test
  namespace: ls-user
spec: {}
```

The controllers of ls-service-target-shoot-sidecar-server automatically creates the required roles, role-bindings etc. 
for all entries in the `SubjectList` *subjects* in every newly created customer namespace (see the details in the image 
before). 

## 2 Credentials and Credential Rotation for a Landscaper as a Service (LaaS) landscape

This chapter describes which credentials are used in a LaaS landscape and how they are created and rotated.

According to the Gardener security recommendations we assume that all involved Gardener shoot clusters are switched 
to non-static secrets. 

If possible the following goals should be achieved: 

- Credentials to access a server should be stored on the client side only.
- Credentials have a restricted lifetime of not more than 90 days.
- No static credentials for shoot clusters. (Gardener recommendation).
- No Gardener-Service-Account with Service Account Manager role (Gardener recommendation).

### Prerequisites

The following summarizes the base techniques to get access to Gardener projects and shoot clusters:

- Gardener-Service-Account: A project member with sufficient permissions could create Gardener-Service-Accounts with 
  different roles for a Gardener project. For a description of the roles, see
  [Create a Service Account](https://gardener.cloud/docs/dashboard/usage/gardener-api/#create-a-service-account).

- Gardener-Service-Account-Kubeconfig: To access a Gardener project, you could download a kubeconfig for a 
  Gardener-Service-Account with the help of the Gardener Dashboard. The lifetime of this kubeconfig is 90 days.
  As an alternative to get a token to access a Gardener Project you can use the [token request API](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/token-request-v1/#TokenRequest):
  ```
  kubectl create --kubeconfig=<path to kubeconfig for Gardener-Service-Account of the project with SAM role> \
    --raw "/api/v1/namespaces/garden-<garden-project-name>/serviceaccounts/<service-account-name>/token" \
    -f <(echo '{"spec":{"expirationSeconds": <some-duration>}}') \
    | jq -r .status
  ```

- Shoot-Cluster-Admin-Kubeconfig: To get a cluster-admin kubeconfig for a Gardener shoot cluster you can use the
  [shootsadminkubeconfig-subresource](https://github.com/gardener/gardener/blob/master/docs/usage/shoot_access.md#shootsadminkubeconfig-subresource).
  The lifetime of the tokens is at most 24 hours. To get such a token you must be an admin member of a project or
  you have an admin Gardener-Service-Account-Kubeconfig at hand.

- Shoot-Cluster-Kubeconfig: Kubeconfig with a longer lifetime and particular permissions for accessing a Gardener 
  shoot cluster could be constructed by creating a service account (with RBAC rules) in the shoot cluster and fetching
  kubeconfigs/token using again the 
  [token request API](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/token-request-v1/#TokenRequest).
  With this approach the credentials are not stored in the shoot cluster and have a restricted lifetime. 
  In principle, it is still possible to create the credentials by deploying a secret for a service account. Such 
  credentials (the automatically injected token) have the disadvantages to have no restricted lifetime and exist on the 
  client and server side. Nevertheless, they must be also be rotated (by deleting and recreating the secret). 

In the following we will reference to the names defined here: Gardener-Service-Account, 
Gardener-Service-Account-Kubeconfig, Shoot-Cluster-Admin-Kubeconfig, Shoot-Cluster-Kubeconfig

### Root Credentials

Credential rotation could not be done completely automatically. There must be some root credentials which has to be
created and rotated manually by a user with the corresponding permissions. These root credentials are stored in the
Secret Store.

#### Root Credentials for Gardener Projects

The Gardener-Service-Account-Kubeconfigs for the Gardener-LaaS-Project and the Gardener-Resource-Cluster-Project 
are root credentials for a LaaS landscape. These must be created and rotated manually by a user with 
`Service Account Manager` role.

The corresponding Gardener-Service-Accounts must have the `Admin` role, but not the `Service Account Manager` role. 
This prevents to request new tokens and is a Gardener recommendation. If a service account is allowed to request 
new tokens and one such token leaks, an attacker could create new tokens forever and rotation could never invalidate 
them, except the corresponding service account is replaced.

#### Root Credentials for GCP Projects

The Gardener shoot cluster custom resources in the Gardener-LaaS-Project and the Gardener-Resource-Cluster-Project 
references secrets containing the secret key of a [GCP](https://cloud.google.com/) service account. This is required 
for the creation and maintenance of shoot clusters by the Gardener. It is recommended to use one secret in both Garden 
projects each containing the secret key of a service account of a different GCP project.

The secret keys for the GCP service account are root credentials of a LaaS landscape. The secret keys must be created 
and rotated manually by a human, and stored and replaced in the corresponding k8s secrets of the Gardener-LaaS-Project 
and the Gardener-Resource-Cluster-Project by the Deploy-Pipeline.

#### Image Pull Secrets

Image pull secrets are another group of root credentials stored in the Secret Store which must be rotated manually. 

### Deploy-Pipeline

This Chapter describes which credentials are created and rotated with every execution by the Deploy-Pipeline.

**Credentials for GCP Project:** In the Gardener projects Gardener-LaaS-Project and Gardener-Resource-Cluster-Project
there are secrets containing the secret key of a GCP service account. The shoot recourses in these projects are 
referencing these secrets, and the Gardener uses the corresponding GCP secret key to create and maintain a Gardener 
shoot cluster in the GCP infrastructure. 

The Deploy-Pipeline must update the GCP secrets keys in these secrets of the Gardener projects with
the latest GCP secret keys stored in the Secret Store. The update of the secrets in the Gardener projects is done 
by the Deploy-Pipeline with the Gardener-Service-Account-Kubeconfigs, also stored in the Secret-Store.

**Credentials in ServiceTargetConfigs:** The pipeline must rotate the credentials stored in the ServiceTargetconfigs.
This can be done as follows:

- With the Gardener-Service-Account-Kubeconfig for the Gardener-LaaS-Project stored in the Secret Store, create a 
  Shoot-Cluster-Admin-Kubeconfig for a Target-Shoot-Cluster. 

- With the Shoot-Cluster-Admin-Kubeconfig create a new token for the service account used for the ServiceTargetConfig.
  and replace the old token with the new one.

The ServiceTargetConfigs might be created manually or automatically by the Deploy-Pipeline. The only important point is
that the token in the kubeconfigs must be rotated with every deployment. 

**Credentials when deploying the Central Landscaper:** The Helm deployment of the Central Landscaper is executed by the
Deploy-Pipeline with a Shoot-Cluster-Admin-Kubeconfig for the Core-Shoot-Cluster, which could be fetched as described before.
We deploy the Central Landscaper via the 
[bundled `landscaper` helm chart](https://github.com/gardener/landscaper/tree/master/charts/landscaper).
No kubeconfig is provided as Helm values for the resource and webhook cluster because this is the same as the host cluster.
The landscaper deployment creates the corresponding ServiceAccount and mounts their token into the landscaper pod (see
[ServiceAccount Admission Controller](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#serviceaccount-admission-controller)),
providing automatic token rotation.

**Credentials when deploying the LaaS:** To deploy the LaaS the Deploy-Pipeline creates on the Core-Shoot-Cluster 

- a Target containing the access information for the Core-Shoot-Cluster where the LaaS should be installed

- a Secret containing the access information of the OCI registry containing the Component Descriptor of the LaaS, the
  container images etc. 

- a [Context](https://github.com/gardener/landscaper/blob/master/docs/usage/Context.md) object referencing the 
  secret with the access information of the OCI registry

- a Secret containing the access information for the Garden-Resource-Cluster-Project

- an Installation for the LaaS instance referencing the Target and the Context object

The creation of these objects is executed with a Shoot-Cluster-Admin-Kubeconfig for the Core-Shoot-Cluster, fetched
as described before.

The Target is created and rotated by the Deploy-Pipeline as follows:

- With the Gardener-Service-Account-Kubeconfig for the Gardener-LaaS-Project, create a Shoot-Cluster-Admin-Kubeconfig 
  for Core-Shoot-Cluster.

- Use the Shoot-Cluster-Admin-Kubeconfig for the Core-Shoot-Cluster to create a service account with the required
  permissions for the Core-Shoot-Cluster and create a Shoot-Cluster-Kubeconfig for this service account. This kubeconfig
  is stored in the Target (or the referenced secret).

The Secret containing the access information of the OCI registry just copies the correspponding information from the
Secret-Store.

The Secret containing the access information for the Garden-Resource-Cluster-Project is also created from the 
corresponding data in the Secret-Store. 

The Installation to deploy the LaaS requires some credentials as input data:

- The name of the just created Secret containing the access information for the Garden-Resource-Cluster-Project. 

- The name of the just created Secret containing the access information of the OCI registry containing the 
  Component Descriptor.

- The name of the Secret in the Garden-Resource-Cluster-Project, containing the secret key of a GCP service account
  which is used by Gardener to create and maintain shoot clusters.

### List of Credentials for Rotation

- Manually rotated credentials in Secret Store
  - Gardener-Service-Account-Kubeconfig for Gardener-LaaS-Project
  - Gardener-Service-Account-Kubeconfig for Gardener-Resource-Cluster-Project
  - Secret key for the GCP service account of the two GCP projects referenced by the shoot clusters

- Automatically rotated by the Deploy-Pipeline
  - Secrets containing the GCP secret keys in the Garden-Projects Gardener-LaaS-Project and 
    Gardener-Resource-Cluster-Project referenced by shoot resources 
  - Credentials stored in the ServiceTargetConfigs
  - Target used to install the LaaS
  - Secret containing Gardener-Service-Account-Kubeconfigs for the Gardener-Resource-Cluster-Project provided as input
    to the LaaS installation.

## 3 Questions and open points

### Open Points

- Describe internal credentials of Landscaper

- Describe internal credentials of LaaS + new controller and how to rotate them? 

- Describe logging stack, monitoring, nginx controller

- Remove the secret manifests from the
  [rbac subchart](https://github.com/gardener/landscaper/blob/master/charts/landscaper/charts/rbac/templates/serviceaccount.yaml).
  Tokens will be generated in future by the token request api.
  For a landscaper instance deployed by LaaS, it must be clarified how the token request api is called.
  Alternatives: either a templating function, or a separate component with blueprint.
