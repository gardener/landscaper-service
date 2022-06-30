// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package webhook_test

import (
	"path/filepath"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gardener/landscaper-service/test/utils/envtest"
)

func TestWebHook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webhook Suite")
}

var (
	testenv *envtest.Environment
)

var _ = BeforeSuite(func() {
	var err error
	projectRoot := filepath.Join("../../")
	testenv, err = envtest.NewEnvironment(projectRoot)
	Expect(err).ToNot(HaveOccurred())

	_, err = testenv.Start()
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	Expect(testenv.Stop()).ToNot(HaveOccurred())
})

func CreateAdmissionRequest(obj runtime.Object) admission.Request {
	objData, err := runtime.Encode(&json.Serializer{}, obj)
	Expect(err).ToNot(HaveOccurred())
	request := admission.Request{}

	gvr, _ := meta.UnsafeGuessKindToResource(obj.GetObjectKind().GroupVersionKind())
	groupVersionResource := metav1.GroupVersionResource{
		Group:    gvr.Group,
		Version:  gvr.Version,
		Resource: gvr.Resource,
	}

	request.Resource = groupVersionResource
	request.RequestResource = &groupVersionResource
	request.Object = runtime.RawExtension{Raw: objData}

	if o, ok := obj.(client.Object); ok {
		request.Name = o.GetName()
		request.Namespace = o.GetNamespace()
	}

	return request
}
