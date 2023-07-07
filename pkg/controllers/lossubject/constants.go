// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package lossubject

import (
	rbacv1 "k8s.io/api/rbac/v1"
)

const SUBJECTLIST_NAME = "subjectlist"

const ROLE_ADMIN string = "admin"
const ROLEBINDING_ADMIN string = "admin-binding"
const ROLE_MEMBER string = "member"
const ROLEBINDING_MEMBER string = "member-binding"
const ROLE_VIEWER string = "viewer"
const ROLEBINDING_VIEWER string = "viewer-binding"

func ALL_ROLES() [3]RoleInfo {
	return [3]RoleInfo{ADMIN_ROLE_INFO(), MEMBER_ROLE_INFO(), VIEWER_ROLE_INFO()}
}

func ADMIN_ROLE_INFO() RoleInfo {
	return RoleInfo{
		RoleName:        ROLE_ADMIN,
		RoleBindingName: ROLEBINDING_ADMIN,
		PrivilegeList: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"landscaper.gardener.cloud"},
				Resources:     []string{"LosSubjectList"},
				ResourceNames: []string{SUBJECTLIST_NAME},
				Verbs:         []string{"get", "watch", "list", "update"},
			},
		},
	}
}
func MEMBER_ROLE_INFO() RoleInfo {
	return RoleInfo{
		RoleName:        ROLE_MEMBER,
		RoleBindingName: ROLEBINDING_MEMBER,
		PrivilegeList: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"landscaper.gardener.cloud"},
				Resources:     []string{"LosSubjectList"},
				ResourceNames: []string{SUBJECTLIST_NAME},
				Verbs:         []string{"get", "watch", "list"},
			},
		},
	}
}
func VIEWER_ROLE_INFO() RoleInfo {
	return RoleInfo{
		RoleName:        ROLE_VIEWER,
		RoleBindingName: ROLEBINDING_VIEWER,
		PrivilegeList:   []rbacv1.PolicyRule{},
	}
}

type RoleInfo struct {
	RoleName        string
	RoleBindingName string
	PrivilegeList   []rbacv1.PolicyRule
}
