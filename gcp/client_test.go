package gcp_test

import (
	"github.com/challiwill/meteorologica/gcp"
	"github.com/challiwill/meteorologica/gcp/gcpfakes"
	storage "google.golang.org/api/storage/v1"

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

		Context("when gcp returns buckets", func() {
			BeforeEach(func() {
				service.BucketsReturns(&storage.Buckets{}, nil)
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
