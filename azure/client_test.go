package azure_test

import (
	"net/http"
	"time"

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
		client = azure.NewClient(logrus.New(), time.Now().Location(), azureServer.URL(), "some-key", "1337")
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
			monthlyUsageReport, err = client.GetBillingData()
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

})

var monthlyUsageResponse = `
one, two, three, four, five
sometimes, you, might, think, you, want json
but really, we, know, you, want, CSV
`
