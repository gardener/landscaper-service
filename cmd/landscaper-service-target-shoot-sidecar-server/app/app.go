// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	lsinstall "github.com/gardener/landscaper/apis/core/install"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"

	lssinstall "github.com/gardener/landscaper-service/pkg/apis/core/install"
	"github.com/gardener/landscaper-service/pkg/controllers/namespaceregistration"
	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"

	"github.com/gardener/landscaper-service/pkg/crdmanager"
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
		LeaderElection:     false,
		Port:               9443,
		MetricsBindAddress: "0",
	}

	if o.Config.Metrics != nil {
		opts.MetricsBindAddress = fmt.Sprintf(":%d", o.Config.Metrics.Port)
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
