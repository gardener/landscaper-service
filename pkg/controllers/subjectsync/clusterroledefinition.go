package subjectsync

import (
	"context"
	"fmt"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterRoleDefinition struct {
	roleName    string
	bindingName string
	rules       []rbacv1.PolicyRule
}

func GetUserClusterRoleDefinition() *ClusterRoleDefinition {
	return &ClusterRoleDefinition{
		roleName:    USER_CLUSTER_ROLE,
		bindingName: USER_CLUSTER_ROLE_BINDING,
		rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"landscaper-service.gardener.cloud"},
				Resources: []string{"subjectlists"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"landscaper.gardener.cloud"},
				Resources: []string{"*"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
}

func GetViewerClusterRoleDefinition() *ClusterRoleDefinition {
	return &ClusterRoleDefinition{
		roleName:    VIEWER_CLUSTER_ROLE,
		bindingName: VIEWER_CLUSTER_ROLE_BINDING,
		rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"landscaper.gardener.cloud"},
				Resources: []string{"installations", "executions", "deployitems"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
}

func (r *ClusterRoleDefinition) PolicyRules() []rbacv1.PolicyRule {
	return r.rules
}

func (r *ClusterRoleDefinition) CreateOrUpdateClusterRole(ctx context.Context, cl client.Client) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.roleName,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, cl, role, func() error {
		role.Rules = r.rules
		return nil
	})
	if err != nil {
		logger.Error(err, "failed ensuring cluster role", lc.KeyResource, r.roleName)
		return fmt.Errorf("failed ensuring cluster role %s: %w", r.roleName, err)
	}
	return nil
}

func (r *ClusterRoleDefinition) CreateOrUpdateClusterRoleBinding(ctx context.Context, cl client.Client, subjects []rbacv1.Subject) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.bindingName,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, cl, roleBinding, func() error {
		roleBinding.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     r.roleName,
		}
		roleBinding.Subjects = subjects
		return nil
	})
	if err != nil {
		logger.Error(err, "failed ensuring cluster role binding", lc.KeyResource, r.bindingName)
		return fmt.Errorf("failed ensuring cluster role binding %s: %w", r.bindingName, err)
	}

	return nil
}
