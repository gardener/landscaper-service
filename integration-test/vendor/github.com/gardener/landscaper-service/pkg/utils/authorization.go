// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"fmt"

	kutils "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateRoleIfNotExistOrUpdate(ctx context.Context, name string, namespace string, rules []rbacv1.PolicyRule, client client.Client) error {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, client, role, func() error {
		role.Rules = rules
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed ensuring role %s: %w", role.Name, err)
	}

	return nil
}

func CreateRoleBindingIfNotExistOrUpdate(ctx context.Context, name string, namespace string, roleName string, client client.Client) error {
	//create role binding
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	_, err := kutils.CreateOrUpdate(ctx, client, roleBinding, func() error {
		roleBinding.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roleName,
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed ensuring role binding %s: %w", roleBinding.Name, err)
	}

	return nil
}

func CreateClusterRoleIfNotExistOrUpdate(ctx context.Context, name string, rules []rbacv1.PolicyRule, client client.Client) error {
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err := kutils.CreateOrUpdate(ctx, client, role, func() error {
		role.Rules = rules
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed ensuring cluster role %s: %w", role.Name, err)
	}

	return nil
}

func CreateClusterRoleBindingIfNotExistOrUpdate(ctx context.Context, name string, clusterRoleName string, client client.Client) error {
	//create role binding
	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := kutils.CreateOrUpdate(ctx, client, roleBinding, func() error {
		roleBinding.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed ensuring cluster role binding %s: %w", roleBinding.Name, err)
	}

	return nil
}
