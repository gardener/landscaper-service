<!--
SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->

# landscaper-service dependency update

landscaper-service uses some external dependencies.

To update these dependencies check new versions for

- [landscaper](https://github.com/gardener/landscaper/releases)
- [ingress-nginx helm chart & image](https://github.com/kubernetes/ingress-nginx/releases)

and adapt versions in file [ocm-settings](../../.landscaper/ocm-settings.yaml)

## OCM reuse logging-stack

The [logging-stack](https://github.tools.sap/ocm-reuse/logging-stack) is now used from [OCM reuse](https://github.tools.sap/ocm-reuse)

To update it's dependencies check new versions for

- [artifacthub.io / cert-manager](https://artifacthub.io/packages/helm/cert-manager/cert-manager)
- [SAP / sap-btp-service-operator](https://github.com/SAP/sap-btp-service-operator/releases/)
- [brancz / kube-rbac-proxy](https://quay.io/repository/brancz/kube-rbac-proxy?tab=tags&tag=latest)
- [fluentbit.io](https://docs.fluentbit.io/manual/installation/docker#tags-and-versions)

and adapt versions in file [logging-stack/settings.yaml](https://github.tools.sap/ocm-reuse/logging-stack/blob/main/logging-stack/settings.yaml)

