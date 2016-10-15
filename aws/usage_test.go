package aws_test

import (
	. "github.com/challiwill/meteorologica/aws"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Usage", func() {
	Describe("Hash", func() {
		var (
			usage      Usage
			otherUsage Usage
		)

		BeforeEach(func() {
			usage = Usage{
				InvoiceID:              "some-invoice-id",
				PayerAccountId:         "some-payer-account-id",
				LinkedAccountId:        "some-linked-account-id",
				RecordType:             "some-record-type",
				RecordID:               "some-record-id",
				BillingPeriodStartDate: "some-start-date",
				BillingPeriodEndDate:   "some-end-date",
				InvoiceDate:            "some-invoice-date",
				PayerAccountName:       "some-payer-account-name",
				LinkedAccountName:      "some-linked-account-name",
				TaxationAddress:        "some-address",
				PayerPONumber:          "some-numnber",
				ProductCode:            "some-code",
				ProductName:            "some-name",
				SellerOfRecord:         "some-record",
				UsageType:              "some-type",
				Operation:              "some-operation",
				RateId:                 "some-rate-id",
				ItemDescription:        "some-description",
				UsageStartDate:         "some-usage-start-date",
				UsageEndDate:           "some-usage-end-date",
				UsageQuantity:          1.2,
				BlendedRate:            "some-rate",
				CurrencyCode:           "some-currency",
				CostBeforeTax:          "some-cost",
				Credits:                "some-credit",
				TaxAmount:              "some-tax",
				TaxType:                "some-type",
				TotalCost:              3.4,
			}
			otherUsage = Usage{
				InvoiceID:              "some-invoice-id",
				PayerAccountId:         "some-payer-account-id",
				LinkedAccountId:        "some-other-linked-account-id",
				RecordType:             "some-record-type",
				RecordID:               "some-record-id",
				BillingPeriodStartDate: "some-start-date",
				BillingPeriodEndDate:   "some-end-date",
				InvoiceDate:            "some-invoice-date",
				PayerAccountName:       "some-payer-account-name",
				LinkedAccountName:      "some-linked-account-name",
				TaxationAddress:        "some-address",
				PayerPONumber:          "some-numnber",
				ProductCode:            "some-code",
				ProductName:            "some-name",
				SellerOfRecord:         "some-record",
				UsageType:              "some-type",
				Operation:              "some-operation",
				RateId:                 "some-rate-id",
				ItemDescription:        "some-description",
				UsageStartDate:         "some-usage-start-date",
				UsageEndDate:           "some-usage-end-date",
				UsageQuantity:          1.2,
				BlendedRate:            "some-rate",
				CurrencyCode:           "some-currency",
				CostBeforeTax:          "some-cost",
				Credits:                "some-credit",
				TaxAmount:              "some-tax",
				TaxType:                "some-type",
				TotalCost:              3.4,
			}
		})

		It("does not return empty string", func() {
			Expect(usage.Hash("some-region")).NotTo(BeEmpty())
		})

		It("returns different hash for different structs", func() {
			Expect(usage.Hash("some-region")).NotTo(Equal(otherUsage.Hash("some-region")))
		})

		It("returns the same hash each time", func() {
			Expect(usage.Hash("some-region")).To(Equal(usage.Hash("some-region")))
		})
	})
})
