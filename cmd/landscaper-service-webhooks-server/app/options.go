// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	goflag "flag"
	"fmt"
	"math"
	"strings"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gardener/landscaper-service/pkg/apis/core"
	"github.com/gardener/landscaper-service/pkg/webhook"
)

func defaultWebhookedResources() map[string]webhook.WebhookedResourceDefinition {
	return map[string]webhook.WebhookedResourceDefinition{
		"landscaperdeployments": {
			APIGroup:     core.GroupName,
			APIVersions:  []string{"v1alpha1"},
			ResourceName: "landscaperdeployments",
		},
		"instances": {
			APIGroup:     core.GroupName,
			APIVersions:  []string{"v1alpha1"},
			ResourceName: "instances",
		},
		"servicetargetconfigs": {
			APIGroup:     core.GroupName,
			APIVersions:  []string{"v1alpha1"},
			ResourceName: "servicetargetconfigs",
		},
		"targetschedulings": {
			APIGroup:     core.GroupName,
			APIVersions:  []string{"v1alpha1"},
			ResourceName: "targetschedulings",
		},
	}
}

// Options holds the landscaper service webhook options
type options struct {
	log                         logging.Logger // Log is the logger instance
	port                        int            // port where the webhook server is running
	disabledWebhooks            string         // lists disabled webhooks as a comma-separated string
	webhookServiceNamespaceName string         // webhook service namespace and name in the format <namespace>/<name>
	webhookServicePort          int32          // port of the webhook service
	certificatesNamespace       string         // the namespace in which the webhook credentials are being created/updated

	webhook webhookOptions
}

// options for the webhook (generated from raw CLI options for easier usage)
type webhookOptions struct {
	webhookServiceNamespace string                                // webhook service namespace
	webhookServiceName      string                                // webhook service name
	webhookServicePort      int32                                 // port of the webhook service
	certificatesNamespace   string                                // the certificate namespace
	enabledWebhooks         []webhook.WebhookedResourceDefinition // which resources should be watched by the webhook
}

// NewOptions returns a new options instance
func NewOptions() *options {
	return &options{
		webhook: webhookOptions{},
	}
}

// AddFlags adds flags passed via command line
func (o *options) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&o.port, "port", 9443, "Specify the port of the webhook server")
	fs.StringVar(&o.disabledWebhooks, "disable-webhooks", "", "Specify validation webhooks that should be disabled ('all' to disable validation completely)")
	fs.StringVar(&o.webhookServiceNamespaceName, "webhook-service", "", "Specify namespace and name of the webhook service (format: <namespace>/<name>)")
	fs.Int32Var(&o.webhookServicePort, "webhook-service-port", 9443, "Specify the port of the webhook service")
	logging.InitFlags(fs)

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

// Complete initializes the options instance and validates flags
func (o *options) Complete() error {
	log, err := logging.GetLogger()
	if err != nil {
		return err
	}
	o.log = log

	err = o.validate() // validate options
	if err != nil {
		return err
	}

	allErrs := field.ErrorList{}
	o.webhook.webhookServicePort = o.webhookServicePort
	o.webhook.enabledWebhooks = filterWebhookedResources(defaultWebhookedResources(), stringListToMap(o.disabledWebhooks))
	if len(o.webhook.enabledWebhooks) != 0 && len(o.webhookServiceNamespaceName) != 0 {
		webhookService := strings.Split(o.webhookServiceNamespaceName, "/")
		o.webhook.webhookServiceNamespace = webhookService[0]
		o.webhook.webhookServiceName = webhookService[1]
	}
	o.webhook.certificatesNamespace = getCertificateNamespace(o)
	return allErrs.ToAggregate()
}

func (o *options) validate() error {
	allErrs := field.ErrorList{}
	dwr := defaultWebhookedResources()
	if len(o.disabledWebhooks) != 0 { // something has been disabled
		// validate that no unknown values are in the list of to-be-disabled webhooks
		allowedWebhooks := allowedWebhookDisables()
		disabledWebhooks := strings.Split(o.disabledWebhooks, ",")
		for _, elem := range disabledWebhooks {
			if _, ok := dwr[elem]; (elem != "all") && !ok {
				allErrs = append(allErrs, field.NotSupported(field.NewPath("--disable-webhooks"), elem, allowedWebhooks))
			}
		}
	}

	if len(o.webhookServiceNamespaceName) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("--webhook-service"), "must not be empty"))
	} else {
		ws := strings.Split(o.webhookServiceNamespaceName, "/")
		if len(ws) < 2 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("--webhook-service"), o.webhookServiceNamespaceName, "must have the format '<namespace>/<name>'"))
		} else {
			if len(ws[0]) == 0 || len(ws[1]) == 0 {
				allErrs = append(allErrs, field.Invalid(field.NewPath("--webhook-service"), o.webhookServiceNamespaceName, "neither name nor namespace of the webhook service must be empty"))
			}
		}
	}

	if o.port <= 0 || o.port > math.MaxUint16 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("--port"), o.port, fmt.Sprintf("must be in range [0, %d]", math.MaxUint16)))
	}
	return allErrs.ToAggregate()
}

// filterWebhookedResources returns a slice of WebhookedResourceDefinitions that contains only those of the given webhookedResources whose ResourceName is not specified in disabledWebhooks
func filterWebhookedResources(webhookedResources map[string]webhook.WebhookedResourceDefinition, disabledWebhooks map[string]bool) []webhook.WebhookedResourceDefinition {
	fwr := []webhook.WebhookedResourceDefinition{}
	if _, ok := disabledWebhooks["all"]; ok {
		return fwr // all webhooks disabled, return empty slice
	}
	for _, wr := range webhookedResources {
		if _, ok := disabledWebhooks[wr.ResourceName]; !ok {
			fwr = append(fwr, wr)
		}
	}
	return fwr
}

// allowedWebhookDisables computes a list of allowed values for the '--disable-webhooks' option
func allowedWebhookDisables() []string {
	dwr := defaultWebhookedResources()
	res := make([]string, len(dwr)+1)
	c := 0
	for _, elem := range dwr {
		res[c] = elem.ResourceName
		c++
	}
	res[c] = "all"
	return res
}

// stringListToMap turns a comma-separated list of strings into pseudo-set that maps all elements of the list to true
func stringListToMap(opt string) map[string]bool {
	res := map[string]bool{}
	tmp := strings.Split(opt, ",")
	for _, t := range tmp {
		res[t] = true
	}
	return res
}

// getCertificateNamespace returns the namespace to use for storing the webhooks server certificate
func getCertificateNamespace(opt *options) string {
	if len(opt.certificatesNamespace) != 0 {
		return opt.certificatesNamespace
	} else {
		return opt.webhook.webhookServiceNamespace
	}
}
