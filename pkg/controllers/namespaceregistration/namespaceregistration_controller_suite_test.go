// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package namespaceregistration_test

import (
	"context"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper-service/pkg/controllers/subjectsync"
	"github.com/gardener/landscaper-service/test/utils/envtest"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NamespaceRegistration Controller Test Suite")
}

var (
	testenv *envtest.Environment
)

var _ = BeforeSuite(func() {
	var err error
	projectRoot := filepath.Join("../../../")
	testenv, err = envtest.NewEnvironment(projectRoot)
	Expect(err).ToNot(HaveOccurred())

	_, err = testenv.Start()
	Expect(err).ToNot(HaveOccurred())

	// prepare ls-user namespace
	ctx := context.Background()
	lsUserNamespace := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: subjectsync.LS_USER_NAMESPACE}}
	Expect(testenv.Client.Create(ctx, &lsUserNamespace)).To(Succeed())
	subjectList := lssv1alpha1.SubjectList{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subjectsync.SUBJECT_LIST_NAME,
			Namespace: subjectsync.LS_USER_NAMESPACE,
		},
		Spec: lssv1alpha1.SubjectListSpec{
			Subjects: []lssv1alpha1.Subject{
				{
					Kind: "User",
					Name: "testuser",
				},
			},
		},
	}
	Expect(testenv.Client.Create(ctx, &subjectList)).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(testenv.Stop()).ToNot(HaveOccurred())
})

// edit namespace resource
