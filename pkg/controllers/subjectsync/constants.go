package subjectsync

const (
	// USER_CLUSTER_ROLE is the cluster-wide admin role
	USER_CLUSTER_ROLE         = "landscaper-service:namespace-registrator"
	USER_CLUSTER_ROLE_BINDING = "landscaper-service:namespace-registrator"

	// LS_USER_ROLE_IN_NAMESPACE is the admin role for namespace "ls-user"
	LS_USER_ROLE_IN_NAMESPACE         = "landscaper-service:namespace-registrator"
	LS_USER_ROLE_BINDING_IN_NAMESPACE = "landscaper-service:namespace-registrator"

	// USER_ROLE_IN_NAMESPACE is the admin role for registered customer namespaces
	USER_ROLE_IN_NAMESPACE         = "landscaper-service:landscaper-user"
	USER_ROLE_BINDING_IN_NAMESPACE = "landscaper-service:landscaper-user"

	// VIEWER_CLUSTER_ROLE is the cluster-wide viewer role
	VIEWER_CLUSTER_ROLE         = "landscaper-service:landscaper-cluster-viewer"
	VIEWER_CLUSTER_ROLE_BINDING = "landscaper-service:landscaper-cluster-viewer"

	// VIEWER_ROLE_IN_NAMESPACE is the viewer role for registered customer namespaces
	VIEWER_ROLE_IN_NAMESPACE         = "landscaper-service:landscaper-viewer"
	VIEWER_ROLE_BINDING_IN_NAMESPACE = "landscaper-service:landscaper-viewer"

	SUBJECT_LIST_NAME = "subjects"
	LS_USER_NAMESPACE = "ls-user"

	SUBJECT_LIST_ENTRY_USER            = "User"
	SUBJECT_LIST_ENTRY_GROUP           = "Group"
	SUBJECT_LIST_ENTRY_SERVICE_ACCOUNT = "ServiceAccount"

	CUSTOM_NS_PREFIX = "cu-"
)
