package gcp_test

import (
	"github.com/challiwill/meteorologica/gcp"
	"github.com/challiwill/meteorologica/gcp/gcpfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gcp", func() {
	var (
		client  *gcp.Client
		service *gcpfakes.FakeStorageService
	)

	BeforeEach(func() {
		service = new(gcpfakes.FakeStorageService)
		client = &gcp.Client{
			StorageService: service,
		}
	})

	Describe("MonthlyUsageReport", func() {
		var (
			usageReport gcp.DetailedUsageReport
			err         error
		)

		JustBeforeEach(func() {
			usageReport, err = client.MonthlyUsageReport()
		})

		XContext("when gcp returns the bucket", func() {
			BeforeEach(func() {
			})

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the buckets contents", func() {
				Expect(usageReport).NotTo(BeEmpty())
			})
		})
	})
})
