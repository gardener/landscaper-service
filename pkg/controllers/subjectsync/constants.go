package subjectsync

const (
	USER_CLUSTER_ROLE                 = "landscaper-service:namespace-registrator"
	USER_CLUSTER_ROLE_BINDING         = "landscaper-service:namespace-registrator"
	LS_USER_ROLE_IN_NAMESPACE         = "landscaper-service:namespace-registrator"
	LS_USER_ROLE_BINDING_IN_NAMESPACE = "landscaper-service:namespace-registrator"
	USER_ROLE_IN_NAMESPACE            = "landscaper-service:landscaper-user"
	USER_ROLE_BINDING_IN_NAMESPACE    = "landscaper-service:landscaper-user"

	SUBJECT_LIST_NAME = "subjects"
	LS_USER_NAMESPACE = "ls-user"

	SUBJECT_LIST_ENTRY_USER            = "User"
	SUBJECT_LIST_ENTRY_GROUP           = "Group"
	SUBJECT_LIST_ENTRY_SERVICE_ACCOUNT = "ServiceAccount"

	CUSTOM_NS_PREFIX = "cu-"
)
