// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"os"

	lsinstall "github.com/gardener/landscaper/apis/core/install"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	lssinstall "github.com/gardener/landscaper-service/pkg/apis/core/install"
	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/controllers/namespaceregistration"
	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"
	"github.com/gardener/landscaper-service/pkg/crdmanager"
	"github.com/gardener/landscaper-service/pkg/utils"
	"github.com/gardener/landscaper-service/pkg/version"
)

// NewResourceClusterControllerCommand creates a new command for the landscaper service controller
func NewResourceClusterControllerCommand(ctx context.Context) *cobra.Command {
	options := NewOptions()

	cmd := &cobra.Command{
		Use:   "resource-cluster-controller",
		Short: "resource-cluster-controller manages the creation/deletion of namespaces and subject sync to rolebindings",

		Run: func(cmd *cobra.Command, args []string) {
			if err := options.Complete(ctx); err != nil {
				fmt.Print(err)
				os.Exit(1)
			}
			if err := options.run(ctx); err != nil {
				options.Log.Error(err, "unable to run resource cluster controller")
			}
		},
	}

	options.AddFlags(cmd.Flags())

	return cmd
}

func (o *options) run(ctx context.Context) error {
	o.Log.Info(fmt.Sprintf("Start Resource Cluster Controller with version %q", version.Get().String()))

	opts := manager.Options{
		LeaderElection: false,
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
		NewClient: utils.NewUncachedClient,
	}

	if o.Config.Metrics != nil {
		opts.Metrics.BindAddress = fmt.Sprintf(":%d", o.Config.Metrics.Port)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), opts)
	if err != nil {
		return fmt.Errorf("unable to setup manager: %w", err)
	}

	if err := o.ensureCRDs(ctx, mgr); err != nil {
		return err
	}

	lssinstall.Install(mgr.GetScheme())
	lsinstall.Install(mgr.GetScheme())

	o.Log.Info("Fetching client for init")
	initClient, err := createClientForInit(mgr.GetConfig())
	if err != nil {
		return err
	}

	if err := createLsUserNamespaceIfNotExist(ctx, initClient); err != nil {
		return fmt.Errorf("failed creating initial required namespace: %w", err)
	}

	lsUserRoleDef := subjectsync.GetLsUserRoleDefinition()

	if err := lsUserRoleDef.CreateOrUpdateRole(ctx, initClient); err != nil {
		return fmt.Errorf("failed creating initial required role: %w", err)
	}
	if err := lsUserRoleDef.CreateRoleBindingWithoutSubjectsIfNotExist(ctx, initClient); err != nil {
		return fmt.Errorf("failed creating initial required rolebinding: %w", err)
	}
	if err := createSubjectsListIfNotExist(ctx, initClient); err != nil {
		return fmt.Errorf("failed creating initial required subjectlist: %w", err)
	}

	ctrlLogger := o.Log.WithName("controllers")
	if err := namespaceregistration.AddControllerToManager(ctrlLogger, mgr, o.Config); err != nil {
		return fmt.Errorf("unable to setup namespaceregistration controller: %w", err)
	}
	if err := subjectsync.AddControllerToManager(ctrlLogger, mgr, o.Config); err != nil {
		return fmt.Errorf("unable to setup subjectsync controller: %w", err)
	}

	o.Log.Info("starting the controllers")
	if err := mgr.Start(ctx); err != nil {
		o.Log.Error(err, "error while running manager")
		os.Exit(1)
	}
	return nil
}

func (o *options) ensureCRDs(ctx context.Context, mgr manager.Manager) error {
	ctx = logging.NewContext(ctx, logging.Wrap(ctrl.Log.WithName("crdManager")))
	crdmgr, err := crdmanager.NewCrdManager(mgr, o.Config.CrdManagement)
	if err != nil {
		return fmt.Errorf("unable to setup CRD manager: %w", err)
	}

	if err := crdmgr.EnsureCRDs(ctx); err != nil {
		return fmt.Errorf("failed to handle CRDs: %w", err)
	}

	return nil
}

func createClientForInit(config *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	lssinstall.Install(scheme)
	lsinstall.Install(scheme)
	utilruntime.Must(rbacv1.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))

	c, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to create init client: %w", err)
	}

	return c, nil
}

func createLsUserNamespaceIfNotExist(ctx context.Context, c client.Client) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: subjectsync.LS_USER_NAMESPACE,
		},
	}
	if err := c.Create(ctx, namespace); err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed creating namespace %s:  %w", namespace.Name, err)
	}
	return nil
}

func createSubjectsListIfNotExist(ctx context.Context, c client.Client) error {
	subjectList := &lssv1alpha1.SubjectList{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subjectsync.SUBJECT_LIST_NAME,
			Namespace: subjectsync.LS_USER_NAMESPACE,
		},
		Spec: lssv1alpha1.SubjectListSpec{
			Subjects: []lssv1alpha1.Subject{},
		},
	}
	if err := c.Create(ctx, subjectList); err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed creating subjectlist %s:  %w", subjectList.Name, err)
	}
	return nil
}
