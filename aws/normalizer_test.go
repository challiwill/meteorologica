package aws_test

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	. "github.com/challiwill/meteorologica/aws"
	"github.com/challiwill/meteorologica/datamodels"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Normalizer", func() {
	var (
		normalizer *Normalizer
		log        *logrus.Logger
		loc        *time.Location
	)

	BeforeEach(func() {
		log = logrus.New()
		log.Out = NewBuffer()
		loc = time.Now().Location()
		normalizer = NewNormalizer(log, loc, "my-region")
	})

	Describe("Normalize", func() {
		var (
			usageReports []*Usage
			reports      datamodels.Reports
		)

		BeforeEach(func() {
			usageReports = []*Usage{
				&Usage{
					InvoiceID:              "some-invoice-id",
					PayerAccountId:         "some-payer-account-id",
					LinkedAccountId:        "some-linked-account-id",
					RecordType:             "LinkedLineItem",
					RecordID:               "some-record-id",
					BillingPeriodStartDate: "some-start-date",
					BillingPeriodEndDate:   "some-end-date",
					InvoiceDate:            "some-invoice-date",
					PayerAccountName:       "some-payer-account-name",
					LinkedAccountName:      "some-linked-account-name",
					TaxationAddress:        "some-address",
					PayerPONumber:          "some-payer-number",
					ProductCode:            "some-product-code",
					ProductName:            "some-product-name",
					SellerOfRecord:         "some-seller-record",
					UsageType:              "some-usage-type",
					Operation:              "some-operation",
					RateId:                 "some-rate-id",
					ItemDescription:        "some-item-description",
					UsageStartDate:         "some-usage-start-date",
					UsageEndDate:           "some-usage-end-date",
					UsageQuantity:          "some-usage-quantity",
					BlendedRate:            "some-blended-rate",
					CurrencyCode:           "some-currency-code",
					CostBeforeTax:          "some-cost",
					Credits:                "some-credits",
					TaxAmount:              "some-tax-amount",
					TaxType:                "some-tax-type",
					TotalCost:              "some-total-cost",
					DailySpend:             "some-daily-spend",
				},
				&Usage{
					InvoiceID:              "some-invoice-id",
					PayerAccountId:         "some-payer-account-id",
					LinkedAccountId:        "",
					RecordType:             "PayerLineItem",
					RecordID:               "some-other-record-id",
					BillingPeriodStartDate: "some-start-date",
					BillingPeriodEndDate:   "some-end-date",
					InvoiceDate:            "some-invoice-date",
					PayerAccountName:       "some-payer-account-name",
					LinkedAccountName:      "",
					TaxationAddress:        "some-address",
					PayerPONumber:          "some-payer-number",
					ProductCode:            "some-other-product-code",
					ProductName:            "some-other-product-name",
					SellerOfRecord:         "some-other-seller-record",
					UsageType:              "some-other-usage-type",
					Operation:              "some-other-operation",
					RateId:                 "some-other-rate-id",
					ItemDescription:        "some-other-item-description",
					UsageStartDate:         "some-other-usage-start-date",
					UsageEndDate:           "some-other-usage-end-date",
					UsageQuantity:          "some-other-usage-quantity",
					BlendedRate:            "some-other-blended-rate",
					CurrencyCode:           "some-other-currency-code",
					CostBeforeTax:          "some-other-cost",
					Credits:                "some-other-credits",
					TaxAmount:              "some-other-tax-amount",
					TaxType:                "some-other-tax-type",
					TotalCost:              "some-other-total-cost",
					DailySpend:             "some-other-daily-spend",
				},
			}
		})

		JustBeforeEach(func() {
			reports = normalizer.Normalize(usageReports)
		})

		Context("with at least one report", func() {
			Context("with valid data", func() {
				It("returns the same number of reports", func() {
					Expect(reports).To(HaveLen((len(usageReports))))
				})

				It("returns properly converted reports", func() {
					Expect(reports[0]).To(Equal(datamodels.Report{
						AccountNumber: "some-linked-account-id",
						AccountName:   "some-linked-account-name",
						Day:           strconv.Itoa(time.Now().Day() - 1),
						Month:         time.Now().Month().String(),
						Year:          strconv.Itoa(time.Now().Year()),
						ServiceType:   "some-product-name",
						UsageQuantity: "some-usage-quantity",
						Cost:          "some-total-cost",
						Region:        "my-region",
						UnitOfMeasure: "",
						IAAS:          "AWS",
					}))
					Expect(reports[1]).To(Equal(datamodels.Report{
						AccountNumber: "some-payer-account-id",
						AccountName:   "some-payer-account-name",
						Day:           strconv.Itoa(time.Now().Day() - 1),
						Month:         time.Now().Month().String(),
						Year:          strconv.Itoa(time.Now().Year()),
						ServiceType:   "some-other-product-name",
						UsageQuantity: "some-other-usage-quantity",
						Cost:          "some-other-total-cost",
						Region:        "my-region",
						UnitOfMeasure: "",
						IAAS:          "AWS",
					}))
				})
			})

			Context("with rows that are tallys", func() {
				XIt("works", func() {})
			})
		})

		Context("with no reports", func() {
			It("returns empty", func() {
				reports := normalizer.Normalize(nil)

				Expect(reports).To(HaveLen(0))
			})
		})
	})
})
