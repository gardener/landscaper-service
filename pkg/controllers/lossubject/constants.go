// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package lossubject

const SUBJECTLIST_NAME = "subjectlist"

const ROLE_ADMIN string = "admin"
const ROLEBINDING_ADMIN string = "admin-binding"
const ROLE_MEMBER string = "member"
const ROLEBINDING_MEMBER string = "member-binding"
const ROLE_VIEWER string = "viewer"
const ROLEBINDING_VIEWER string = "viewer-binding"

func ALL_ROLES() [3]string {
	return [3]string{ROLE_ADMIN, ROLE_MEMBER, ROLE_VIEWER}
}
