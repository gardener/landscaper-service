package namespaceregistration

import (
	"context"
	"fmt"

	"github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/gardener/landscaper/apis/core/v1alpha1/helper"
	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	lc "github.com/gardener/landscaper/controller-utils/pkg/logging/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lssv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha2"
	"github.com/gardener/landscaper-service/pkg/utils"
)

const keyOnDeleteStrategy = "onDeleteStrategy"

type triggerDeletionFunc func(ctx context.Context, cl client.Client, inst *v1alpha1.Installation) error

func getTriggerDeletionFunction(ctx context.Context, namespaceRegistration *lssv1alpha2.NamespaceRegistration) (triggerDeletionFunc, error) {
	logger, _ := logging.FromContextOrNew(ctx, nil)

	strategy, found := namespaceRegistration.Annotations[lssv1alpha2.LandscaperServiceOnDeleteStrategyAnnotation]
	if !found {
		strategy = ""
	}
	logger.Info("determined on-delete-strategy", keyOnDeleteStrategy, strategy)

	switch strategy {
	case "":
		return triggerDeletionByDefaultStrategy, nil
	case lssv1alpha2.LandscaperServiceOnDeleteStrategyDeleteAllInstallations:
		return triggerDeletionWithUninstall, nil
	case lssv1alpha2.LandscaperServiceOnDeleteStrategyDeleteAllInstallationsWithoutUninstall:
		return triggerDeletionWithoutUninstall, nil
	default:
		logger.Info("unknown on-delete-strategy", keyOnDeleteStrategy, strategy)
		return nil, fmt.Errorf("unknown on-delete-strategy %q", strategy)
	}
}

// triggerDeletionByDefaultStrategy deletes the installation if it is a root installation and
// has the delete-without-uninstall annotation.
func triggerDeletionByDefaultStrategy(ctx context.Context, cl client.Client, inst *v1alpha1.Installation) error {
	_, ctx = logging.FromContextOrNew(ctx, nil,
		lc.KeyResource, client.ObjectKeyFromObject(inst).String(),
		keyOnDeleteStrategy, "")

	if utils.HasDeleteWithoutUninstallAnnotation(&inst.ObjectMeta) {
		return deleteOrRetriggerDelete(ctx, cl, inst)
	}

	return nil
}

// triggerDeletionWithUninstall deletes the installation.
func triggerDeletionWithUninstall(ctx context.Context, cl client.Client, inst *v1alpha1.Installation) error {
	_, ctx = logging.FromContextOrNew(ctx, nil,
		lc.KeyResource, client.ObjectKeyFromObject(inst).String(),
		keyOnDeleteStrategy, lssv1alpha2.LandscaperServiceOnDeleteStrategyDeleteAllInstallations)

	return deleteOrRetriggerDelete(ctx, cl, inst)
}

// triggerDeletionWithoutUninstall deletes the installation without uninstall.
func triggerDeletionWithoutUninstall(ctx context.Context, cl client.Client, inst *v1alpha1.Installation) error {
	_, ctx = logging.FromContextOrNew(ctx, nil,
		lc.KeyResource, client.ObjectKeyFromObject(inst).String(),
		keyOnDeleteStrategy, lssv1alpha2.LandscaperServiceOnDeleteStrategyDeleteAllInstallationsWithoutUninstall)

	if err := ensureDeleteWithoutUninstallAnnotation(ctx, cl, inst); err != nil {
		return err
	}

	return deleteOrRetriggerDelete(ctx, cl, inst)
}

func ensureDeleteWithoutUninstallAnnotation(ctx context.Context, cl client.Client, inst *v1alpha1.Installation) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	if !utils.HasDeleteWithoutUninstallAnnotation(&inst.ObjectMeta) {
		metav1.SetMetaDataAnnotation(&inst.ObjectMeta, v1alpha1.DeleteWithoutUninstallAnnotation, "true")
		if err := cl.Update(ctx, inst); err != nil {
			logger.Error(err, "failed adding delete-without-uninstall annotation to installation")
			return err
		}
	}

	return nil
}

func deleteOrRetriggerDelete(ctx context.Context, cl client.Client, inst *v1alpha1.Installation) error {
	logger, ctx := logging.FromContextOrNew(ctx, nil)

	if inst.GetDeletionTimestamp().IsZero() {
		if err := cl.Delete(ctx, inst); err != nil {
			logger.Error(err, "failed deleting installations: "+client.ObjectKeyFromObject(inst).String())
			return err
		}
	} else if inst.Status.JobID == inst.Status.JobIDFinished && !helper.HasOperation(inst.ObjectMeta, v1alpha1.ReconcileOperation) {
		// retrigger
		metav1.SetMetaDataAnnotation(&inst.ObjectMeta, v1alpha1.OperationAnnotation, string(v1alpha1.ReconcileOperation))
		if err := cl.Update(ctx, inst); err != nil {
			logger.Error(err, "failed annotating installations without uninstall")
			return err
		}
	}

	return nil
}
