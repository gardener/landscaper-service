package subjectsync

import (
	"context"
	"fmt"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RoleDefinition struct {
	namespace   string
	roleName    string
	bindingName string
	rules       []rbacv1.PolicyRule
}

// GetLsUserRoleDefinition defines the admin role for the "ls-user" namespace.
func GetLsUserRoleDefinition() *RoleDefinition {
	return &RoleDefinition{
		namespace:   LS_USER_NAMESPACE,
		roleName:    LS_USER_ROLE_IN_NAMESPACE,
		bindingName: LS_USER_ROLE_BINDING_IN_NAMESPACE,
		rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"landscaper-service.gardener.cloud"},
				Resources:     []string{"subjectlists"},
				ResourceNames: []string{SUBJECT_LIST_NAME},
				Verbs:         []string{"get", "update", "patch", "list", "watch"},
			},
			{
				APIGroups:     []string{"landscaper-service.gardener.cloud"},
				Resources:     []string{"subjectlists/status"},
				ResourceNames: []string{SUBJECT_LIST_NAME},
				Verbs:         []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"landscaper-service.gardener.cloud"},
				Resources: []string{"namespaceregistrations"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"landscaper-service.gardener.cloud"},
				Resources: []string{"namespaceregistrations/status"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts/token"},
				Verbs:     []string{"create"},
			},
		},
	}
}

// GetUserRoleDefinition defines the admin role for a customer namespace generated from a NamespaceRegistration.
func GetUserRoleDefinition(namespace string) *RoleDefinition {
	return &RoleDefinition{
		namespace:   namespace,
		roleName:    USER_ROLE_IN_NAMESPACE,
		bindingName: USER_ROLE_BINDING_IN_NAMESPACE,
		rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"landscaper.gardener.cloud"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets", "configmaps"},
				Verbs:     []string{"*"},
			},
		},
	}
}

// GetViewerRoleDefinition defines the viewer role for a customer namespace generated from a NamespaceRegistration.
func GetViewerRoleDefinition(namespace string) *RoleDefinition {
	return &RoleDefinition{
		namespace:   namespace,
		roleName:    VIEWER_ROLE_IN_NAMESPACE,
		bindingName: VIEWER_ROLE_BINDING_IN_NAMESPACE,
		rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"landscaper.gardener.cloud"},
				Resources: []string{"installations", "executions", "deployitems"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
}

func (r *RoleDefinition) PolicyRules() []rbacv1.PolicyRule {
	return r.rules
}

func (r *RoleDefinition) roleString() string {
	return client.ObjectKey{Namespace: r.namespace, Name: r.roleName}.String()
}

func (r *RoleDefinition) roleBindingString() string {
	return client.ObjectKey{Namespace: r.namespace, Name: r.bindingName}.String()
}

func (r *RoleDefinition) CreateOrUpdateRole(ctx context.Context, cl client.Client) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.roleName,
			Namespace: r.namespace,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, cl, role, func() error {
		role.Rules = r.rules
		return nil
	})
	if err != nil {
		logger.Error(err, "failed ensuring role", lc.KeyResource, r.roleString())
		return fmt.Errorf("failed ensuring role %s: %w", r.roleString(), err)
	}

	return nil
}

func (r *RoleDefinition) CreateOrUpdateRoleBinding(ctx context.Context, cl client.Client, subjects []rbacv1.Subject) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	//create role binding
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.bindingName,
			Namespace: r.namespace,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, cl, roleBinding, func() error {
		roleBinding.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     r.roleName,
		}
		roleBinding.Subjects = subjects
		return nil
	})
	if err != nil {
		logger.Error(err, "failed ensuring role binding", lc.KeyResource, r.roleBindingString())
		return fmt.Errorf("failed ensuring role binding %s: %w", r.roleBindingString(), err)
	}

	return nil
}

func (r *RoleDefinition) CreateRoleBindingWithoutSubjectsIfNotExist(ctx context.Context, cl client.Client) error {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.bindingName,
			Namespace: r.namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     r.roleName,
		},
	}

	if err := cl.Create(ctx, roleBinding); err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed creating role binding %s: %w", r.roleBindingString(), err)
	}

	return nil
}

func UpdateRoleBindingSubjects(ctx context.Context, cl client.Client, binding *rbacv1.RoleBinding, subjects []rbacv1.Subject) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	binding.Subjects = subjects
	if err := cl.Update(ctx, binding); err != nil {
		logger.Error(err, "failed updating role binding")
		return fmt.Errorf("failed updating role binding %s %s: %w", binding.Namespace, binding.Name, err)
	}

	return nil
}

func (r *RoleDefinition) DeleteRole(ctx context.Context, cl client.Client) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	role := &rbacv1.Role{}
	if err := cl.Get(ctx, types.NamespacedName{Name: r.roleName, Namespace: r.namespace}, role); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("role in namespace not found", lc.KeyResource, r.roleString())
		} else {
			logger.Error(err, "failed loading role", lc.KeyResource, r.roleString())
			return fmt.Errorf("failed loading role %s: %w", r.roleString(), err)
		}
	} else {
		if err := cl.Delete(ctx, role); err != nil {
			logger.Error(err, "failed deleting role", lc.KeyResource, r.roleString())
			return fmt.Errorf("failed deleting role %s: %w", r.roleString(), err)
		}
	}

	return nil
}

func (r *RoleDefinition) DeleteRoleBinding(ctx context.Context, cl client.Client) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	roleBinding := &rbacv1.RoleBinding{}
	if err := cl.Get(ctx, types.NamespacedName{Name: r.bindingName, Namespace: r.namespace}, roleBinding); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("rolebinding in namespace not found", lc.KeyResource, r.roleBindingString())
		} else {
			logger.Error(err, "failed loading rolebinding", lc.KeyResource, r.roleBindingString())
			return fmt.Errorf("failed loading rolebinding %s: %w", r.roleBindingString(), err)
		}
	} else {
		if err := cl.Delete(ctx, roleBinding); err != nil {
			logger.Error(err, "failed deleting rolebinding", lc.KeyResource, r.roleBindingString())
			return fmt.Errorf("failed deleting rolebinding %s: %w", r.roleBindingString(), err)
		}
	}

	return nil
}
