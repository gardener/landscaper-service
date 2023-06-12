package mutating

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	lsscore "github.com/gardener/landscaper-service/pkg/apis/core"
)

const (
	TenantRegistrationResourceType = "tenantregistrations"
)

// MutatorFromResourceType is a helper method that gets a resource type and returns the fitting mutator
func MutatorFromResourceType(log logging.Logger, kubeClient client.Client, scheme *runtime.Scheme, resource string) (GenericMutator, error) {
	abstrMut := newAbstractMutator(log, kubeClient, scheme)
	var val GenericMutator
	if resource == TenantRegistrationResourceType {
		val = &TenantRegistrationMutator{abstrMut}
	} else {
		return nil, fmt.Errorf("unable to find mutator for resource type %q", resource)
	}
	return val, nil
}

type abstractMutator struct {
	Client  client.Client
	decoder runtime.Decoder
	log     logging.Logger
}

// newAbstractMutator creates a new abstracted mutator
func newAbstractMutator(log logging.Logger, kubeClient client.Client, scheme *runtime.Scheme) abstractMutator {
	return abstractMutator{
		Client:  kubeClient,
		decoder: serializer.NewCodecFactory(scheme).UniversalDecoder(),
		log:     log,
	}
}

// GenericMutator is an abstraction interface that implements admission.Handler and contains additional setter functions for the fields
type GenericMutator interface {
	Handle(context.Context, admission.Request) admission.Response
}

// TENANT REGISTRATION

// TenantRegistrationMutator represents a mutator for a TenantRegistrations
type TenantRegistrationMutator struct{ abstractMutator }

// Handle handles a request to the webhook
func (tm *TenantRegistrationMutator) Handle(_ context.Context, req admission.Request) admission.Response {
	log := tm.log.WithName("TenantRegistrationMutator").WithValues("name", req.Name, "namespace", req.Namespace)

	tenantRegistration := &lsscore.TenantRegistration{}
	_, gkv, err := tm.decoder.Decode(req.Object.Raw, nil, tenantRegistration)
	if err != nil {
		log.Error(err, "failed decoding TenantRegistration from admision request")
		return admission.Errored(http.StatusBadRequest, err)
	}
	//since tenantRegistration GKV is empty after decoding, it needs to be set for later diff calculation. Otherwise a JSON patch of removing it would be created and fails
	tenantRegistration.SetGroupVersionKind(*gkv)

	if tenantRegistration.Spec.Author == req.UserInfo.Username {
		log.Info("Author already set correctly")
		return admission.Allowed("Author already set correctly")
	}
	tenantRegistration.Spec.Author = req.UserInfo.Username
	userAddedTenantRegistration, err := json.Marshal(tenantRegistration)
	if err != nil {
		log.Error(err, "failed marshaling Tenant Registration with author added")
		return admission.Errored(http.StatusBadRequest, errors.New("failed marshaling user attributed Tenant Registration"))
	}
	res := admission.PatchResponseFromRaw(req.Object.Raw, userAddedTenantRegistration)
	return res

	//TODO: or shall we manually create a path reuquest, since with the delta calculation, we get some serialised defaults into it, but if value is omitted, replace is wrong and add would be the right operation
	// patchForUsername := jsonpatch.NewOperation("replace", "/spec/author", req.UserInfo.Username)
	// res := admission.Patched("add username as author", patchForUsername)
	// return res

}
