package azure_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	. "github.com/challiwill/meteorologica/azure"
	"github.com/challiwill/meteorologica/resources"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Azure", func() {
	var (
		client      *Client
		azureServer *ghttp.Server
	)

	BeforeEach(func() {
		azureServer = ghttp.NewServer()
		client = NewClient(logrus.New(), time.Now().Location(), azureServer.URL(), "some-key", 1337)
	})

	AfterEach(func() {
		azureServer.Close()
	})

	Describe("Name", func() {
		It("returns the IAAS name", func() {
			Expect(client.Name()).To(Equal("Azure"))
		})
	})

	XDescribe("GetNormalizedUsage", func() {
		It("works", func() {})
	})

	Describe("GetBillingData", func() {
		var (
			monthlyUsageReport []byte
			err                error
		)

		JustBeforeEach(func() {
			monthlyUsageReport, err = client.GetBillingData()
		})

		//TODO these tests will fail some days, that is bad
		Context("when azure returns valid data", func() {
			BeforeEach(func() {
				azureServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/rest/1337/usage-report", fmt.Sprintf("month=%d-%s&type=detail", time.Now().Year(), resources.PadMonth(time.Now().Month()))),
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

		Context("when azure returns an error", func() {
			BeforeEach(func() {
				azureServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/rest/1337/usage-report", fmt.Sprintf("month=%d-%s&type=detail", time.Now().Year(), resources.PadMonth(time.Now().Month()))),
						ghttp.VerifyHeaderKV("authorization", "bearer some-key"),
						ghttp.VerifyHeaderKV("api-version", "1.0"),
						ghttp.RespondWith(http.StatusInternalServerError, ""),
					),
				)
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("responded with error"))
			})

			It("makes a GET request to azure", func() {
				Expect(azureServer.ReceivedRequests()).To(HaveLen(1))
			})
		})
	})
})

var monthlyUsageResponse = `
one, two, three, four, five
sometimes, you, might, think, you, want json
but really, we, know, you, want, CSV
`
