package datamodels_test

import (
	"time"

	. "github.com/challiwill/meteorologica/datamodels"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reports", func() {
	Describe("ConsolidateReports", func() {
		var (
			reports      Reports
			usageReports Reports
		)

		JustBeforeEach(func() {
			usageReports = ConsolidateReports(reports)
		})

		Context("with valid reports", func() {
			BeforeEach(func() {
				reports = Reports{
					Report{
						ID:            "a",
						AccountNumber: "some-account",
						AccountName:   "some-account-name",
						Day:           1,
						Month:         time.Month(2),
						Year:          2016,
						ServiceType:   "some-product",
						UsageQuantity: 0.1,
						Cost:          0.3,
						Region:        "some-region",
						UnitOfMeasure: "some-unit-of-measure",
						Resource:      "my-iaas",
					},
					Report{
						ID:            "b",
						AccountNumber: "some-account",
						AccountName:   "some-account-name",
						Day:           2,
						Month:         time.Month(2),
						Year:          2016,
						ServiceType:   "some-product",
						UsageQuantity: 1.1,
						Cost:          2.3,
						Region:        "some-region",
						UnitOfMeasure: "some-unit-of-measure",
						Resource:      "my-iaas",
					},
					Report{
						ID:            "c",
						AccountNumber: "some-account",
						AccountName:   "some-account-name",
						Day:           1,
						Month:         time.Month(2),
						Year:          2016,
						ServiceType:   "some-other-product",
						UsageQuantity: 0.1,
						Cost:          0.3,
						Region:        "some-region",
						UnitOfMeasure: "some-unit-of-measure",
						Resource:      "my-iaas",
					},
					Report{
						ID:            "a",
						AccountNumber: "some-account",
						AccountName:   "some-account-name",
						Day:           1,
						Month:         time.Month(2),
						Year:          2016,
						ServiceType:   "some-product",
						UsageQuantity: 0.1,
						Cost:          0.3,
						Region:        "some-region",
						UnitOfMeasure: "some-unit-of-measure",
						Resource:      "my-iaas",
					},
					Report{
						ID:            "d",
						AccountNumber: "some-other-account",
						AccountName:   "some-account-name",
						Day:           1,
						Month:         time.Month(2),
						Year:          2016,
						ServiceType:   "some-product",
						UsageQuantity: 0.1,
						Cost:          0.3,
						Region:        "some-region",
						UnitOfMeasure: "some-unit-of-measure",
						Resource:      "my-iaas",
					},
				}
			})

			It("returns one aggregate row for matching service type for same account", func() {
				Expect(len(usageReports)).To(Equal(4))
				Expect(usageReports).To(ContainElement(
					Report{
						ID:            "d",
						AccountNumber: "some-other-account",
						AccountName:   "some-account-name",
						Day:           1,
						Month:         time.Month(2),
						Year:          2016,
						ServiceType:   "some-product",
						UsageQuantity: 0.1,
						Cost:          0.3,
						Region:        "some-region",
						UnitOfMeasure: "some-unit-of-measure",
						Resource:      "my-iaas",
					},
				))
				Expect(usageReports).To(ContainElement(
					Report{
						ID:            "c",
						AccountNumber: "some-account",
						AccountName:   "some-account-name",
						Day:           1,
						Month:         time.Month(2),
						Year:          2016,
						ServiceType:   "some-other-product",
						UsageQuantity: 0.1,
						Cost:          0.3,
						Region:        "some-region",
						UnitOfMeasure: "some-unit-of-measure",
						Resource:      "my-iaas",
					},
				))
				Expect(usageReports).To(ContainElement(
					Report{
						ID:            "a",
						AccountNumber: "some-account",
						AccountName:   "some-account-name",
						Day:           1,
						Month:         time.Month(2),
						Year:          2016,
						ServiceType:   "some-product",
						UsageQuantity: 0.2,
						Cost:          0.6,
						Region:        "some-region",
						UnitOfMeasure: "some-unit-of-measure",
						Resource:      "my-iaas",
					},
				))
				Expect(usageReports).To(ContainElement(
					Report{
						ID:            "b",
						AccountNumber: "some-account",
						AccountName:   "some-account-name",
						Day:           2,
						Month:         time.Month(2),
						Year:          2016,
						ServiceType:   "some-product",
						UsageQuantity: 1.1,
						Cost:          2.3,
						Region:        "some-region",
						UnitOfMeasure: "some-unit-of-measure",
						Resource:      "my-iaas",
					},
				))
			})
		})
	})
})
