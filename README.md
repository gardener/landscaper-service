# Landscaper Service

[![reuse compliant](https://reuse.software/badge/reuse-compliant.svg)](https://reuse.software/)

🚧 This repository is heavily under construction and should be considered experimental.

The _Landscaper Service_ aims to implement a Kubernetes API extension, which provides and manages Landscaper installations as a service.
Its goal is to install the [Landscaper](https://github.com/gardener/landscaper) on configurable target Kubernetes clusters. Landscaper installations are isolated for each tenant, as each installation is hosted in its own virtual cluster (node-less Kubernetes cluster with its own Kube API server and etcd).  
