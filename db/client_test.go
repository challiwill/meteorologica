package db_test

import (
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
						AccountNumber: "12345",
						AccountName:   "my-account",
						Day:           17,
						Month:         "March",
						Year:          1337,
						ServiceType:   "some-service",
						UsageQuantity: 0.65,
						Cost:          12.58,
						Region:        "some-region",
						UnitOfMeasure: "GB",
						IAAS:          "MySpecialIAAS",
					},
					datamodels.Report{
						AccountNumber: "12345",
						AccountName:   "my-account",
						Day:           13,
						Month:         "January",
						Year:          1905,
						ServiceType:   "special-service",
						UsageQuantity: 0.65,
						Cost:          12.58,
						UnitOfMeasure: "GB",
						IAAS:          "LessPreferredIAAS",
					},
				}
			})

			It("returns nil", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("saves the reports in the database", func() {
				Expect(fakedb.ExecCallCount()).To(Equal(len(reports)))
				query0, args0 := fakedb.ExecArgsForCall(0)
				Expect(query0).To(ContainSubstring("INSERT IGNORE INTO iaas_billing"))
				Expect(query0).To(ContainSubstring("(AccountNumber, AccountName, Day, Month, Year, ServiceType, UsageQuantity, Cost, Region, UnitOfMeasure, IAAS)"))
				Expect(query0).To(ContainSubstring("values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"))
				Expect(args0[0]).To(Equal("12345"))
				Expect(args0[1]).To(Equal("my-account"))
				Expect(args0[2]).To(Equal(17))
				Expect(args0[3]).To(Equal("March"))
				Expect(args0[4]).To(Equal(1337))
				Expect(args0[5]).To(Equal("some-service"))
				Expect(args0[6]).To(Equal(0.65))
				Expect(args0[7]).To(Equal(12.58))
				Expect(args0[8]).To(Equal("some-region"))
				Expect(args0[9]).To(Equal("GB"))
				Expect(args0[10]).To(Equal("MySpecialIAAS"))
				query1, args1 := fakedb.ExecArgsForCall(1)
				Expect(query1).To(ContainSubstring("INSERT IGNORE INTO iaas_billing"))
				Expect(query1).To(ContainSubstring("(AccountNumber, AccountName, Day, Month, Year, ServiceType, UsageQuantity, Cost, Region, UnitOfMeasure, IAAS)"))
				Expect(query1).To(ContainSubstring("values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"))
				Expect(args1[0]).To(Equal("12345"))
				Expect(args1[1]).To(Equal("my-account"))
				Expect(args1[2]).To(Equal(13))
				Expect(args1[3]).To(Equal("January"))
				Expect(args1[4]).To(Equal(1905))
				Expect(args1[5]).To(Equal("special-service"))
				Expect(args1[6]).To(Equal(0.65))
				Expect(args1[7]).To(Equal(12.58))
				Expect(args1[8]).To(Equal(""))
				Expect(args1[9]).To(Equal("GB"))
				Expect(args1[10]).To(Equal("LessPreferredIAAS"))
			})
		})

		Context("With no reports", func() {
			BeforeEach(func() {
				reports = datamodels.Reports{}
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Ping", func() {
		It("pings the database", func() {
			err := client.Ping()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakedb.PingCallCount()).To(Equal(1))
		})
	})

	XDescribe("Migrate", func() {
		It("works", func() {})
	})
})
