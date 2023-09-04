# Namespaceregistrtion

A user, with access to the Resource-Shoot-Cluster as described before, is only allowed to create Landscaper resources
like Installations, Targets etc. in so-called customer namespaces. A customer namespace is a normal namespace on the
Resource-Shoot-Cluster with a name starting with the prefix *cu-*.

## Create a Customer Namespace

To create such a customer namespace the user must create a
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

If the creation of a customer namespace was successful, the status of the `NamespaceRegistration` looks as follows:

```yaml
status:
  phase: Completed
```

In case of an error you find the corresponding message in the status section.

## Deleting Namespaceregistrations

When deleting a `NamespaceRegistration` the corresponding namespace is deleted. This is executed accoding to the 
following strategy:

- All root Installations with a "delete-without-uninstall" annotation 
  ([see](https://github.com/gardener/landscaper/blob/master/docs/usage/Annotations.md#delete-without-uninstall-annotation))
  are deleted.
- As long as there are still Install...