// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package mutating

import (
	"context"
	"fmt"
	"path"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/gardener/landscaper-service/pkg/webhook"

	"sigs.k8s.io/controller-runtime/pkg/client"

	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"
)

func UpdateMutatingWebhookConfiguration(ctx context.Context, kubeClient client.Client, o webhook.Options) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyMethod, "UpdateMutatingWebhookConfiguration"})

	// do not deploy or update the webhook if no service name is given
	if len(o.ServiceName) == 0 || len(o.ServiceNamespace) == 0 {
		return nil
	}

	mwc := admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: o.WebhookConfigurationName,
		},
	}

	// construct MutatingWebhookConfiguration
	noSideEffects := admissionregistrationv1.SideEffectClassNone
	failPolicy := admissionregistrationv1.Fail
	mwcWebhooks := []admissionregistrationv1.MutatingWebhook{}

	for _, elem := range o.WebhookedResources {
		rule := admissionregistrationv1.RuleWithOperations{
			Operations: []admissionregistrationv1.OperationType{admissionregistrationv1.Create, admissionregistrationv1.Update},
			Rule:       admissionregistrationv1.Rule{},
		}
		rule.Rule.APIGroups = []string{elem.APIGroup}
		rule.Rule.APIVersions = elem.APIVersions
		rule.Rule.Resources = []string{elem.ResourceName}
		clientConfig := admissionregistrationv1.WebhookClientConfig{
			CABundle: o.CABundle,
		}
		webhookPath := path.Join(o.WebhookBasePath, elem.ResourceName)
		clientConfig.Service = &admissionregistrationv1.ServiceReference{
			Namespace: o.ServiceNamespace,
			Name:      o.ServiceName,
			Path:      &webhookPath,
			Port:      &o.ServicePort,
		}
		mwcWebhook := admissionregistrationv1.MutatingWebhook{
			Name:                    elem.ResourceName + o.WebhookNameSuffix,
			SideEffects:             &noSideEffects,
			FailurePolicy:           &failPolicy,
			ObjectSelector:          &o.ObjectSelector,
			AdmissionReviewVersions: []string{"v1"},
			Rules:                   []admissionregistrationv1.RuleWithOperations{rule},
			ClientConfig:            clientConfig,
		}
		mwcWebhooks = append(mwcWebhooks, mwcWebhook)
	}

	logger.Info("Creating/updating MutatingWebhookConfiguration", lc.KeyResource, o.WebhookConfigurationName, lc.KeyResourceKind, "MutatingWebhookConfiguration")
	_, err := ctrl.CreateOrUpdate(ctx, kubeClient, &mwc, func() error {
		mwc.Webhooks = mwcWebhooks
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to create/update MutatingWebhookConfiguration: %w", err)
	}
	logger.Info("MutatingWebhookConfiguration created/updated", lc.KeyResource, o.WebhookConfigurationName, lc.KeyResourceKind, "MutatingWebhookConfiguration")

	return nil

}

// DeleteMutatingWebhookConfiguration deletes a ValidatingWebhookConfiguration
func DeleteMutatingWebhookConfiguration(ctx context.Context, kubeClient client.Client, name string) error {
	logger, ctx := logging.FromContextOrNew(ctx, []interface{}{lc.KeyMethod, "DeleteMutatingWebhookConfiguration"})

	vwc := admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	logger.Info("Removing MutatingWebhookConfiguration, if it exists", lc.KeyResource, name, lc.KeyResourceKind, "MutatingWebhookConfiguration")
	if err := kubeClient.Delete(ctx, &vwc); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("MutatingWebhookConfiguration not found", lc.KeyResource, name, lc.KeyResourceKind, "MutatingWebhookConfiguration")
		} else {
			return fmt.Errorf("unable to delete MutatingWebhookConfiguration %q: %w", name, err)
		}
	} else {
		logger.Info("MutatingWebhookConfiguration deleted", lc.KeyResource, name, lc.KeyResourceKind, "MutatingWebhookConfiguration")
	}
	return nil
}

// RegisterWebhooks generates certificates and registers the webhooks to the manager
// no-op if WebhookedResources in the given options is either nil or empty
func RegisterWebhooks(ctx context.Context, webhookServer *ctrlwebhook.Server, client client.Client, scheme *runtime.Scheme, o webhook.Options) error {
	logger, _ := logging.FromContextOrNew(ctx, []interface{}{lc.KeyMethod, "RegisterWebhooks"})

	if o.WebhookedResources == nil || len(o.WebhookedResources) == 0 {
		return nil
	}

	// registering webhooks
	for _, elem := range o.WebhookedResources {
		rsLogger := logger.WithName(elem.ResourceName)
		val, err := MutatorFromResourceType(rsLogger, client, scheme, elem.ResourceName)
		if err != nil {
			return fmt.Errorf("unable to register webhooks: %w", err)
		}

		webhookPath := o.WebhookBasePath + elem.ResourceName
		rsLogger.Info("Registering webhook", lc.KeyResource, elem.ResourceName, "path", webhookPath)
		admission := &ctrlwebhook.Admission{Handler: val}
		_ = admission.InjectLogger(rsLogger.Logr())
		if err := admission.InjectScheme(scheme); err != nil {
			return err
		}
		webhookServer.Register(webhookPath, admission)
	}

	return nil
}
