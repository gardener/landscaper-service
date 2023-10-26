<!--
SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# landscaper-service dependency update

landscaper-service uses various external dependencies. All of these dependencies are transported to the [`eu.gcr.io/gardener-project`](https://console.cloud.google.com/gcr/images/gardener-project/EU/landscaper-service) repository
before they are being used in the blueprints.
To update these dependencies, a set of scripts are available.
These scripts transport the specified version to the repository and update the resources of their respective OCM component.

## gcr.io login

The landscaper-service dependencies are stored in the `eu.gcr.io` repository.
The following steps may be used to set up authentication for docker and helm.

```shell
export HELM_EXPERIMENTAL_OCI=1

gcloud auth login
gcloud auth configure-docker
gcloud auth print-access-token | helm registry login -u oauth2accesstoken --password-stdin https://eu.gcr.io
```

## cert-manager

Check new version on [artifacthub.io / cert-manager](https://artifacthub.io/packages/helm/cert-manager/cert-manager)

```shell
export CERT_MANAGER_VERSION=v1.13.1

./hack/dependency-update/cert-manager.sh ${CERT_MANAGER_VERSION}
```

## SAP BTP Operator

Check new version on [SAP / sap-btp-service-operator](https://github.com/SAP/sap-btp-service-operator/releases/) and [brancz / kube-rbac-proxy](https://quay.io/repository/brancz/kube-rbac-proxy?tab=tags&tag=latest)

```shell
export SAP_BTP_OPERATOR_VERSION=v0.5.3
export KUBE_RBAC_PROXY_VERSION=v0.14.3

./hack/dependency-update/sap-btp-operator.sh ${SAP_BTP_OPERATOR_VERSION} ${KUBE_RBAC_PROXY_VERSION}
```

## fluent-bit

Check new version on [fluentbit.io](https://docs.fluentbit.io/manual/installation/docker#tags-and-versions)

```shell
export FLUENTBIT_VERSION=2.1.9

./hack/dependency-update/fluentbit.sh ${FLUENTBIT_VERSION}
```

## ingress-controller

Check new version for [helm chart & image](https://github.com/kubernetes/ingress-nginx/releases)

```shell
export INGRESS_NGINX_CHART_VERSION=4.8.3
export INGRESS_NGINX_IMAGE_VERSION=v1.9.4

./hack/dependency-update/ingress-controller.sh ${INGRESS_NGINX_CHART_VERSION} ${INGRESS_NGINX_IMAGE_VERSION}
```
