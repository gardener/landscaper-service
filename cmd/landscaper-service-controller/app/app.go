// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"os"

	lsinstall "github.com/gardener/landscaper/apis/core/install"
	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	lssinstall "github.com/gardener/landscaper-service/pkg/apis/core/install"
	instancesctrl "github.com/gardener/landscaper-service/pkg/controllers/instances"
	landscaperdeploymentsctrl "github.com/gardener/landscaper-service/pkg/controllers/landscaperdeployments"
	servicetargetconfigsctrl "github.com/gardener/landscaper-service/pkg/controllers/servicetargetconfigs"
	"github.com/gardener/landscaper-service/pkg/crdmanager"
	"github.com/gardener/landscaper-service/pkg/version"
)

// NewLandscaperServiceControllerCommand creates a new command for the landscaper service controller
func NewLandscaperServiceControllerCommand(ctx context.Context) *cobra.Command {
	options := NewOptions()

	cmd := &cobra.Command{
		Use:   "landscaper-service-controller",
		Short: "Landscaper Service controller manages the installation and lifecycle of Landscaper installations",

		Run: func(cmd *cobra.Command, args []string) {
			if err := options.Complete(ctx); err != nil {
				fmt.Print(err)
				os.Exit(1)
			}
			if err := options.run(ctx); err != nil {
				options.Log.Error(err, "unable to run landscaper service controller")
			}
		},
	}

	options.AddFlags(cmd.Flags())

	return cmd
}

func (o *options) run(ctx context.Context) error {
	o.Log.Info(fmt.Sprintf("Start Landscaper Service Controller with version %q", version.Get().String()))

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

	crdmgr, err := crdmanager.NewCrdManager(ctrl.Log.WithName("setup").WithName("CRDManager"), mgr, o.Config.CrdManagement)
	if err != nil {
		return fmt.Errorf("unable to setup CRD manager: %w", err)
	}

	if err := crdmgr.EnsureCRDs(ctx); err != nil {
		return fmt.Errorf("failed to handle CRDs: %w", err)
	}

	lssinstall.Install(mgr.GetScheme())
	lsinstall.Install(mgr.GetScheme())

	ctrlLogger := o.Log.WithName("controllers")
	if err := landscaperdeploymentsctrl.AddControllerToManager(ctrlLogger, mgr, o.Config); err != nil {
		return fmt.Errorf("unable to setup landscaper deployments controller: %w", err)
	}
	if err := instancesctrl.AddControllerToManager(ctrlLogger, mgr, o.Config); err != nil {
		return fmt.Errorf("unable to setup instances controller: %w", err)
	}
	if err := servicetargetconfigsctrl.AddControllerToManager(ctrlLogger, mgr, o.Config); err != nil {
		return fmt.Errorf("unable to setup service target configs controller: %w", err)
	}

	o.Log.Info("starting the controllers")
	if err := mgr.Start(ctx); err != nil {
		o.Log.Error(err, "error while running manager")
		os.Exit(1)
	}
	return nil
}
