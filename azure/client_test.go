package azure_test

import (
	"net/http"

	"github.com/challiwill/meteorologica/azure"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Azure", func() {
	var (
		client      *azure.Client
		azureServer *ghttp.Server
	)

	BeforeEach(func() {
		azureServer = ghttp.NewServer()
		client = azure.NewClient(azureServer.URL(), "some-key", "1337")
	})

	AfterEach(func() {
		azureServer.Close()
	})

	Describe("UsageReports", func() {
		var (
			usageReports azure.UsageReports
			err          error
		)

		JustBeforeEach(func() {
			usageReports, err = client.UsageReports()
		})

		Context("When azure returns usage list", func() {
			BeforeEach(func() {
				azureServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/rest/1337/usage-reports"),
						ghttp.VerifyHeaderKV("authorization", "bearer some-key"),
						ghttp.VerifyHeaderKV("api-version", "1.0"),
						ghttp.RespondWith(http.StatusOK, availableMonthsResponse),
					),
				)
			})

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("makes a GET request to azure", func() {
				Expect(azureServer.ReceivedRequests()).To(HaveLen(1))
			})

			It("returns resources", func() {
				Expect(usageReports.Months).To(HaveLen(2))
			})
		})
	})
})

var availableMonthsResponse = `
{
	"object_type" : "Usage",
	"contract_version" : "1.0",
	"AvailableMonths":
	[
		{
			"Month":"2014-02",
			"LinkToDownloadSummaryReport":"/rest/100100/usagereport?month=2014-02&type=summary",
			"LinkToDownloadDetailReport":"/rest/100100/usagereport?month=2014-02&type=detail"
		}
		,{
			"Month":"2014-03",
			"LinkToDownloadSummaryReport":"/rest/100100/usagereport?month=2014-03&type=summary",
			"LinkToDownloadDetailReport":"/rest/100100/usage-report?month=2014-03&type=detail"
		}
	]
}
`
