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
landscaper-deployment-qmjgx   Succeeded   landscaper-deployment-qmjgx   11m
landscaper-rbac-4b4wg         Succeeded   landscaper-rbac-4b4wg         11m
shoot-cluster-x4cgs           Succeeded   shoot-cluster-x4cgs           11m
test-h97dx-mtklg              Succeeded                                 11m

```

The installation status can also be inspected with the [landscaper-cli](https://github.com/gardener/landscapercli).

```sh
landscaper-cli installations inspect -n laas-user test-h97dx-mtklg

[✅ Succeeded] Installation test-h97dx-mtklg
    ├── [✅ Succeeded] Installation landscaper-deployment-qmjgx
    │   └── [✅ Succeeded] DeployItem landscaper-deployment-qmjgx-landscaper-6df5f
    ├── [✅ Succeeded] Installation landscaper-rbac-4b4wg
    │   └── [✅ Succeeded] DeployItem landscaper-rbac-4b4wg-landscaper-rbac-4dkm9
    └── [✅ Succeeded] Installation shoot-cluster-x4cgs
        └── [✅ Succeeded] DeployItem shoot-cluster-x4cgs-shoot-cluster-qppwx

```

Once the installation has successfully finished, the landscaper service controller will update the Instance status with the `clusterEndpoint` and `clusterKubeconfig` information.

```sh
kubectl -n laas-user get instances.landscaper-service.gardener.cloud test-h97dx -o jsonpath="{.status}" | jq
{
  "adminKubeconfig": "a3ViZWNvbmZpZyBjb250ZW50 ...",
  "clusterEndpoint": "https://api.ef5818d3.laas.shoot.mydomain.com",
  "contextRef": {
    "name": "test-h97dx-c9bx9",
    "namespace": "laas-user"
  },
  "installationRef": {
    "name": "test-h97dx-mtklg",
    "namespace": "laas-user"
  },
  "landscaperServiceComponent": {
    "name": "github.com/gardener/landscaper-service/landscaper-instance",
    "version": "v0.41.0"
  },
  "observedGeneration": 2,
  "shootName": "ef5818d3",
  "shootNamespace": "garden-laas",
  "targetRef": {
    "name": "test-h97dx-w4pkc",
    "namespace": "laas-user"
  },
  "userKubeconfig": "a3ViZWNvbmZpZyBjb250ZW50 ..."
}
```

The `status.userKubeconfig` and `status.adminKubeconfig` fields are base64 encoded and can be exported into a local kubeconfig file.

```sh
kubectl -n laas-user get instances.landscaper-service.gardener.cloud test-8qh5w -o jsonpath="{.status.userKubeconfig}" | base64 -d > user-kubeconfig.yaml
kubectl -n laas-user get instances.landscaper-service.gardener.cloud test-8qh5w -o jsonpath="{.status.adminKubeconfig}" | base64 -d > admin-kubeconfig.yaml
```

These kubeconfig files can be used to authenticate at the deployed Landscaper instance.