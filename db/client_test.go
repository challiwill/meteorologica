package db_test

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/challiwill/meteorologica/db"
	"github.com/challiwill/meteorologica/db/dbfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		fakedb *dbfakes.FakeDB
		client *db.Client
	)

	BeforeEach(func() {
		fakedb = new(dbfakes.FakeDB)
		client = &db.Client{
			Conn: fakedb,
			Log:  logrus.New(),
		}
	})

	Describe("Close", func() {
		It("Closes the database connection", func() {
			err := client.Close()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakedb.CloseCallCount()).To(Equal(1))
		})
	})

	Describe("SaveReports", func() {
		var (
			reports datamodels.Reports
			err     error
		)

		JustBeforeEach(func() {
			err = client.SaveReports(reports)
		})

		Context("Given valid reports", func() {
			BeforeEach(func() {
				reports = datamodels.Reports{
					datamodels.Report{
						ID:            "some-id",
						AccountNumber: "12345",
						AccountName:   "my-account",
						Day:           17,
						Month:         time.Month(3),
						Year:          1337,
						ServiceType:   "some-service",
						UsageQuantity: 0.65,
						Cost:          12.58,
						Region:        "some-region",
						UnitOfMeasure: "GB",
						Resource:      "MySpecialIAAS",
					},
					datamodels.Report{
						ID:            "some-other-id",
						AccountNumber: "12345",
						AccountName:   "my-account",
						Day:           13,
						Month:         time.Month(1),
						Year:          1905,
						ServiceType:   "special-service",
						UsageQuantity: 0.65,
						Cost:          12.58,
						UnitOfMeasure: "GB",
						Resource:      "LessPreferredIAAS",
					},
				}
			})

			It("returns nil", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("saves the reports in the database", func() {
				Expect(fakedb.ExecCallCount()).To(Equal(len(reports)))
				_, args0 := fakedb.ExecArgsForCall(0)
				Expect(args0[0]).To(Equal("some-id"))
				Expect(args0[1]).To(Equal("12345"))
				Expect(args0[2]).To(Equal("my-account"))
				Expect(args0[3]).To(Equal(17))
				Expect(args0[4]).To(Equal(time.Month(3)))
				Expect(args0[5]).To(Equal(1337))
				Expect(args0[6]).To(Equal("some-service"))
				Expect(args0[7]).To(Equal("some-region"))
				Expect(args0[8]).To(Equal("MySpecialIAAS"))
				Expect(args0[9]).To(Equal(0.65))
				Expect(args0[10]).To(Equal("GB"))
				Expect(args0[11]).To(Equal(12.58))
				_, args1 := fakedb.ExecArgsForCall(1)
				Expect(args1[0]).To(Equal("some-other-id"))
				Expect(args1[1]).To(Equal("12345"))
				Expect(args1[2]).To(Equal("my-account"))
				Expect(args1[3]).To(Equal(13))
				Expect(args1[4]).To(Equal(time.Month(1)))
				Expect(args1[5]).To(Equal(1905))
				Expect(args1[6]).To(Equal("special-service"))
				Expect(args1[7]).To(Equal(""))
				Expect(args1[8]).To(Equal("LessPreferredIAAS"))
				Expect(args1[9]).To(Equal(0.65))
				Expect(args1[10]).To(Equal("GB"))
				Expect(args1[11]).To(Equal(12.58))
			})
		})

		Context("With no reports", func() {
			BeforeEach(func() {
				reports = datamodels.Reports{}
			})

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
