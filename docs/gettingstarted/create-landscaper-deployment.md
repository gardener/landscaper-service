<!--
SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# Create a Landscaper Deployment

To trigger the deployment of a Landscaper instance by the Landscaper as a Service, 
a [LandscaperDeployment](../usage/LandscaperDeployments.md) resource needs to be created.
The LandscaperDeployment specifies the configuration of the Landscaper deployment as well as the version of the Landscaper to deploy.
If not already existing, a namespace for the deployment needs to be created.

```sh
kubectl create namespace laas-user
```

In the next step, the LandscaperDeployment resource is created. The field `spec.landscaperConfiguration.deployers` has to contain the list of active deployers.
The field `spec.tenantId` has to contain a globally unique tenant id.

```yaml
apiVersion: landscaper-service.gardener.cloud/v1alpha1
kind: LandscaperDeployment
metadata:
  name: test
spec:
  tenantId: "tenant-1"
  purpose: "test"
  landscaperConfiguration:
    deployers:
      - helm
      - manifest
      - container
```

```sh
kubectl -n laas-user apply -f deployment.yaml
```

After the LandscaperDeployment has been created, the landscaper controller select an appropriate [ServiceTargetConfig](../usage/ServiceTargetConfigs.md) and create an [Instance](../usage/Instances.md) .

```sh
kubectl -n laas-user get landscaperdeployments.landscaper-service.gardener.cloud test

NAME   INSTANCE     AGE
test   test-8qh5w   13m
```

The Instance will show the selected ServiceTargetConfig as well as the Installation that has been automatically created by the landscaper service controller.

```sh
kubectl -n laas-user get instances.landscaper-service.gardener.cloud test-8qh5w

NAME         SERVICETARGETCONFIG   INSTALLATION       AGE
test-8qh5w   default               test-8qh5w-hmzrp   15m
```

The installation will automatically start several sub-installations. Once all installations are in phase `Succeeded`, the Landscaper has been deployed successfully.

```sh
kubectl -n laas-user get installations 

NAME                          PHASE       EXECUTION                     AGE
landscaper-deployment-25mxk   Succeeded   landscaper-deployment-25mxk   6m43s
landscaper-rbac-d5j98         Succeeded   landscaper-rbac-d5j98         6m43s
test-8qh5w-hmzrp              Succeeded                                 6m45s
virtual-garden-x95xd          Succeeded   virtual-garden-x95xd          6m43s
```

The installation status can also be inspected with the [landscaper-cli](https://github.com/gardener/landscapercli).

```sh
landscaper-cli installations inspect -n laas-user test-8qh5w-hmzrp

[✅ Succeeded] Installation test-8qh5w-hmzrp
    ├── [✅ Succeeded] Installation virtual-garden-x95xd
    │   └── [✅ Succeeded] DeployItem virtual-garden-x95xd-virtual-garden-container-deployer-6h2qv
    ├── [✅ Succeeded] Installation landscaper-rbac-d5j98
    │   └── [✅ Succeeded] DeployItem landscaper-rbac-d5j98-landscaper-rbac-nqw4j
    └── [✅ Succeeded] Installation landscaper-deployment-25mxk
        └── [✅ Succeeded] DeployItem landscaper-deployment-25mxk-landscaper-8c5cz
```

Once the installation has successfully finished, the landscaper service controller will update the Instance status with the `clusterEndpoint` and `clusterKubeconfig` information.

```sh
kubectl -n laas-user get instances.landscaper-service.gardener.cloud test-8qh5w -o jsonpath="{.status}" | jq
{
  "clusterEndpoint": "10.0.0.1",
  "clusterKubeconfig": "a3ViZWNvbmZpZyBjb250ZW50 ...",
  "installationRef": {
    "name": "test-8qh5w-hmzrp",
    "namespace": "laas-user"
  },
  "observedGeneration": 1,
  "targetRef": {
    "name": "test-8qh5w-k88bs",
    "namespace": "laas-user"
  },
  "contextRef": {
    "name": "test-8qh5w-a5w3s",
    "namespace": "laas-user"
  }
}
```

The `status.clusterKubeconfig` field is base64 encode and can be exported into a local kubeconfig file.

```sh
kubectl -n laas-user get instances.landscaper-service.gardener.cloud test-8qh5w -o jsonpath="{.status.clusterKubeconfig}" | base64 -d > landscaper-kubeconfig.yaml
```

This kubeconfig file can be used to authenticate at the deployed Landscaper instance.