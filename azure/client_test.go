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
		client = azure.NewClient(azureServer.URL())
	})

	AfterEach(func() {
		azureServer.Close()
	})

	Describe("Resources", func() {
		var (
			resources azure.Resources
			err       error
		)

		JustBeforeEach(func() {
			resources, err = client.Resources()
		})

		Context("When azure returns resources", func() {
			BeforeEach(func() {
				azureServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/subscriptions/"),
						ghttp.RespondWith(http.StatusOK, resourceJSON),
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
				Expect(resources.Resources).To(HaveLen(2))
			})
		})
	})
})

var resourceJSON = `
{
  "value": [
		{
      "id": "/subscriptions/d657c399-e17c-405d-859e-9f2efb6462e5/providers/Microsoft.Commerce/UsageAggregates/Daily_BRSDT_20150515_0000",
      "name": "Daily_BRSDT_20150515_0000",
      "type": "Microsoft.Commerce/UsageAggregate",
      "properties": {
        "subscriptionId": "d657c399-e17c-405d-859e-9f2efb6462e5",
        "usageStartTime": "2015-05-15T00:00:00+00:00",
        "usageEndTime": "2015-05-16T00:00:00+00:00",
        "instanceData": "{\"Microsoft.Resources\":{\"resourceUri\":\"/subscriptions/d657c399-e17c-405d-859e-9f2efb6462e5/resourceGroups/moinakrg/providers/Microsoft.Storage/storageAccounts/moinakstorage\",\"location\":\"West US\",\"tags\":{\"department\":\"hr\"}}}",
        "meterName": "Storage Transactions (in 10,000s)",
        "meterCategory": "Data Management",
        "unit": "10,000s",
        "meterId": "964c283a-83a3-4dd4-8baf-59511998fe8b",
        "infoFields": {

        },
        "quantity": 9.8390
      }
    },
    {
      "id": "/subscriptions/d657c399-e17c-405d-859e-9f2efb6462e5/providers/Microsoft.Commerce/UsageAggregates/Daily_BRSDT_20150515_0000",
      "name": "Daily_BRSDT_20150515_0000",
      "type": "Microsoft.Commerce/UsageAggregate",
      "properties": {
        "subscriptionId": "d657c399-e17c-405d-859e-9f2efb6462e5",
        "usageStartTime": "2015-05-15T00:00:00+00:00",
        "usageEndTime": "2015-05-16T00:00:00+00:00",
        "instanceData": "{\"Microsoft.Resources\":{\"resourceUri\":\"/subscriptions/d657c399-e17c-405d-859e-9f2efb6462e5/resourceGroups/moinakrg/providers/Microsoft.Storage/storageAccounts/moinakstorage\",\"location\":\"West US\",\"tags\":{\"department\":\"hr\"}}}",
        "meterName": "Data Transfer In (GB)",
        "meterRegion": "Zone 1",
        "meterCategory": "Networking",
        "unit": "GB",
        "meterId": "32c3ebec-1646-49e3-8127-2cafbd3a04d8",
        "infoFields": {

        },
        "quantity": 0.000066
      }
    }
	]}
`
