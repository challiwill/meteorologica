package azure_test

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
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
		client = azure.NewClient(logrus.New(), azureServer.URL(), "some-key", "1337")
	})

	AfterEach(func() {
		azureServer.Close()
	})

	Describe("GetCSV", func() {
		var (
			monthlyUsageReport []byte
			err                error
		)

		JustBeforeEach(func() {
			monthlyUsageReport, err = client.GetCSV()
		})

		Context("When azure returns valid data", func() {
			BeforeEach(func() {
				azureServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/rest/1337/usage-report", "type=detail"),
						ghttp.VerifyHeaderKV("authorization", "bearer some-key"),
						ghttp.VerifyHeaderKV("api-version", "1.0"),
						ghttp.RespondWith(http.StatusOK, monthlyUsageResponse),
					),
				)
			})

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("makes a GET request to azure", func() {
				Expect(azureServer.ReceivedRequests()).To(HaveLen(1))
			})

			It("returns monthly usage report", func() {
				Expect(monthlyUsageReport).To(Equal([]byte(monthlyUsageReport)))
			})
		})
	})

	Describe("MakeDetailedUsageReport", func() {
		var (
			azureMonthlyUsage azure.DetailedUsageReport
			azureDataFile     []byte
			err               error
		)

		BeforeEach(func() {
			azureDataFile, err = ioutil.ReadFile("../testfixtures/short-azure.csv")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			// delete azure.csv file
		})

		JustBeforeEach(func() {
			azureMonthlyUsage = azure.MakeDetailedUsageReport(azureDataFile)
		})

		It("returns a parsable csv", func() {
			err = ioutil.WriteFile("/tmp/azure.csv", azureMonthlyUsage.CSV, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			azureDataFile, err := os.OpenFile("/tmp/azure.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
			defer azureDataFile.Close()
			Expect(err).NotTo(HaveOccurred())

			usageReader, err := azure.NewUsageReader(logrus.New(), azureDataFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(usageReader.UsageReports[0].SubscriptionName).NotTo(BeEmpty())
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

var monthlyUsageResponse = `
one, two, three, four, five
sometimes, you, might, think, you, want json
but really, we, know, you, want, CSV
`
