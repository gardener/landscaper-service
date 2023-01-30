# landscaper-service dependency update

landscaper-service uses various external dependencies. All of these dependencies are transported to the `eu.gcr.io/gardener-project` repository
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

```shell
export CERT_MANAGER_VERSION=v1.11.0

./hack/dependency-update/cert-manager.sh ${CERT_MANAGER_VERSION}
```

## SAP BTP Operator

```shell
export SAP_BTP_OPERATOR_VERSION=v0.3.6

./hack/dependency-update/sap-btp-operator.sh ${SAP_BTP_OPERATOR_VERSION}
```

## fluent-bit
```shell
export FLUENTBIT_VERSION=2.0.8

./hack/dependency-update/fluentbit.sh ${FLUENTBIT_VERSION}
```

## ingress-controller
```shell
export INGRESS_NGINX_CHART_VERSION=4.4.2
export INGRESS_NGINX_IMAGE_VERSION=v1.5.1

./hack/dependency-update/ingress-controller.sh ${INGRESS_NGINX_CHART_VERSION} ${INGRESS_NGINX_IMAGE_VERSION}
```