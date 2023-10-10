// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package avuploader_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"

	lssv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha2"
	"github.com/gardener/landscaper-service/test/utils/envtest"

	avuploader "github.com/gardener/landscaper-service/pkg/controllers/avuploader"
)

var _ = Describe("Reconcile", func() {
	var (
		ctx   context.Context
		state *envtest.State
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	AfterEach(func() {
		defer ctx.Done()
		if state != nil {
			Expect(testenv.CleanupResources(ctx, state)).ToNot(HaveOccurred())
		}
	})
})

var _ = Describe("avs request construction", func() {
	It("should construct a UP avs request", func() {
		availabilityCollection := lssv1alpha2.AvailabilityCollection{
			Status: lssv1alpha2.AvailabilityCollectionStatus{
				Instances: []lssv1alpha2.AvailabilityInstance{
					{
						Status: string(lsv1alpha1.LsHealthCheckStatusOk),
						ObjectReference: lssv1alpha2.ObjectReference{
							Name:      "instance1",
							Namespace: "instance1-namespace",
						},
					},
				},
			},
		}
		request := avuploader.ExportConstructAvsRequest(availabilityCollection)
		Expect(request.Status).To(Equal(avuploader.AVS_STATUS_UP))
		Expect(len(request.Instances)).To(Equal(0))
	})

	It("should construct a DOWN avs request", func() {
		availabilityCollection := lssv1alpha2.AvailabilityCollection{
			Status: lssv1alpha2.AvailabilityCollectionStatus{
				Instances: []lssv1alpha2.AvailabilityInstance{
					{
						Status: string(lsv1alpha1.LsHealthCheckStatusOk),
						ObjectReference: lssv1alpha2.ObjectReference{
							Name:      "instance1",
							Namespace: "instance1-namespace",
						},
					},
					{
						Status:       string(lsv1alpha1.LsHealthCheckStatusFailed),
						FailedReason: "timeout",
						ObjectReference: lssv1alpha2.ObjectReference{
							Name:      "instance2",
							Namespace: "instance2-namespace",
						},
					},
				},
				Self: lssv1alpha2.AvailabilityInstance{
					Status:       string(lsv1alpha1.LsHealthCheckStatusFailed),
					FailedReason: "timeout2",
					ObjectReference: lssv1alpha2.ObjectReference{
						Name:      "self",
						Namespace: "landscaper",
					},
				},
			},
		}
		request := avuploader.ExportConstructAvsRequest(availabilityCollection)
		Expect(request.Status).To(Equal(avuploader.AVS_STATUS_DOWN))
		Expect(request.OutageReason).To(ContainSubstring("2/3"))
		Expect(len(request.Instances)).To(Equal(2))
		Expect(request.Instances[0].InstanceId).To(Equal("instance2"))
		Expect(request.Instances[0].OutageReason).To(ContainSubstring("timeout"))
		Expect(request.Instances[1].InstanceId).To(Equal("Self"))
		Expect(request.Instances[1].OutageReason).To(ContainSubstring("timeout2"))
	})
})
