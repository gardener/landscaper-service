// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubernetesscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	webhookcert "github.com/gardener/landscaper/controller-utils/pkg/webhook"

	"github.com/gardener/landscaper-service/pkg/apis/core/install"
	"github.com/gardener/landscaper-service/pkg/version"
	"github.com/gardener/landscaper-service/pkg/webhook"
)

// NewLandscaperServiceWebhooksCommand creates a new command for the landscaper service webhooks server
func NewLandscaperServiceWebhooksCommand(ctx context.Context) *cobra.Command {
	options := NewOptions()

	cmd := &cobra.Command{
		Use:   "landscaper-service-webhooks",
		Short: "Landscaper Service webhooks serves the landscaper service validation, mutating and defaulting webhooks.",

		Run: func(cmd *cobra.Command, args []string) {
			if err := options.Complete(); err != nil {
				fmt.Print(err)
				os.Exit(1)
			}
			if err := options.run(ctx); err != nil {
				options.log.Error(err, "unable to run landscaper service webhooks server")
				os.Exit(1)
			}
		},
	}
	options.AddFlags(cmd.Flags())
	return cmd
}

func (o *options) run(ctx context.Context) error {
	o.log.Info(fmt.Sprintf("Start Landscaper Service Webhooks Server with version %q", version.Get().String()))

	opts := ctrlwebhook.Options{}
	opts.Port = o.port
	opts.WebhookMux = http.NewServeMux()
	opts.WebhookMux.Handle("/healthz", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		if _, err := writer.Write([]byte("Ok")); err != nil {
			o.log.Error(err, "unable to send health response")
		}
	}))
	opts.CertDir = filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs")
	webhookServer := ctrlwebhook.NewServer(opts)

	ctrl.SetLogger(o.log.Logr())

	restConfig := ctrl.GetConfigOrDie()
	scheme := runtime.NewScheme()
	utilruntime.Must(kubernetesscheme.AddToScheme(scheme))
	install.Install(scheme)

	kubeClient, err := webhook.GetCachelessClient(restConfig)
	if err != nil {
		return fmt.Errorf("unable to get client: %w", err)
	}

	// create ValidatingWebhookConfiguration and register webhooks, if validation is enabled, delete it otherwise
	if err := registerWebhooks(ctx, webhookServer, kubeClient, scheme, opts.CertDir, o); err != nil {
		return fmt.Errorf("unable to register validation webhook: %w", err)
	}

	if err := webhookServer.Start(ctx); err != nil {
		o.log.Error(err, "error while starting webhook server")
		os.Exit(1)
	}
	return nil
}

func registerWebhooks(ctx context.Context,
	webhookServer ctrlwebhook.Server,
	kubeClient client.Client,
	scheme *runtime.Scheme,
	certDir string,
	o *options) error {

	webhookLogger := logging.Wrap(ctrl.Log.WithName("webhook").WithName("validation"))
	ctx = logging.NewContext(ctx, webhookLogger)

	webhookConfigurationName := "landscaper-service-validation-webhook"
	// noop if all webhooks are disabled
	if len(o.webhook.enabledWebhooks) == 0 {
		webhookLogger.Info("Validation disabled")
		return webhook.DeleteValidatingWebhookConfiguration(ctx, kubeClient, webhookConfigurationName)
	}

	webhookLogger.Info("Validation enabled")

	// initialize webhook options
	wo := webhook.Options{
		WebhookConfigurationName: webhookConfigurationName,
		WebhookBasePath:          "/webhook/validate/",
		WebhookNameSuffix:        ".validation.landscaper-service.gardener.cloud",
		ObjectSelector: metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Operator: metav1.LabelSelectorOpNotIn,
					Key:      "validation.landscaper-service.gardener.cloud/skip-validation",
					Values:   []string{"true"},
				},
			},
		},
		ServicePort:        o.webhook.webhookServicePort,
		ServiceName:        o.webhook.webhookServiceName,
		ServiceNamespace:   o.webhook.webhookServiceNamespace,
		WebhookedResources: o.webhook.enabledWebhooks,
	}

	// generate certificates
	var err error
	dnsNames := webhookcert.GeDNSNamesFromNamespacedName(wo.ServiceNamespace, wo.ServiceName)
	caCert, _, err := webhookcert.GenerateCertificates(ctx, kubeClient, certDir, o.webhook.certificatesNamespace,
		"landscaper-service-webhook", "landscaper-service-webhook-cert", dnsNames)
	if err != nil {
		return fmt.Errorf("unable to generate webhook certificates: %w", err)
	}

	wo.CABundle = caCert.CertificatePEM

	// log which resources are being watched
	webhookedResourcesLog := []string{}
	for _, elem := range wo.WebhookedResources {
		webhookedResourcesLog = append(webhookedResourcesLog, elem.ResourceName)
	}
	webhookLogger.Info("Enabling validation", "resources", webhookedResourcesLog)

	// create/update/delete ValidatingWebhookConfiguration
	if err := webhook.UpdateValidatingWebhookConfiguration(ctx, kubeClient, wo); err != nil {
		return err
	}
	// register webhooks
	if err := webhook.RegisterWebhooks(ctx, webhookServer, kubeClient, scheme, wo); err != nil {
		return err
	}

	return nil
}
