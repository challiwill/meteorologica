package aws_test

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	. "github.com/challiwill/meteorologica/aws"
	"github.com/challiwill/meteorologica/aws/awsfakes"
	"github.com/challiwill/meteorologica/datamodels"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

//go:generate counterfeiter io.ReadCloser

var _ = Describe("Client", func() {
	var (
		client    *Client
		log       *logrus.Logger
		s3Client  *awsfakes.FakeS3Client
		logOutput *Buffer
		dbClient  *awsfakes.FakeReportsDatabase
	)

	BeforeEach(func() {
		log = logrus.New()
		logOutput = NewBuffer()
		log.Out = logOutput
		s3Client = new(awsfakes.FakeS3Client)
		dbClient = new(awsfakes.FakeReportsDatabase)
		client = NewClient(log, time.Now().Location(), "my-region", "my-bucket", 1234567890, s3Client, dbClient)
	})

	Describe("Name", func() {
		It("returns the IAAS name", func() {
			Expect(client.Name()).To(Equal("AWS"))
		})
	})

	XDescribe("GetNormalizedUsage", func() {
		It("works", func() {})
	})

	Describe("GetBillingData", func() {
		var (
			usage []byte
			err   error
		)

		JustBeforeEach(func() {
			usage, err = client.GetBillingData()
		})

		Context("when AWS returns a billing file", func() {
			BeforeEach(func() {
				readCloser := new(awsfakes.FakeReadCloser)
				readCloser.ReadStub = func(p []byte) (int, error) {
					return copy(p, dailyUsageResponse), io.EOF
				}

				s3Client.GetObjectReturns(&s3.GetObjectOutput{Body: readCloser}, nil)
			})

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("requests the correct file and bucket", func() {
				Expect(s3Client.GetObjectCallCount()).To(Equal(1))
				object := s3Client.GetObjectArgsForCall(0)
				Expect(object.Bucket).To(Equal(aws.String("my-bucket")))

				expectedFileName := fmt.Sprintf("1234567890-aws-billing-csv-%d-%d.csv", time.Now().Year(), time.Now().Month())
				if time.Now().Month() < 10 {
					expectedFileName = fmt.Sprintf("1234567890-aws-billing-csv-%d-0%d.csv", time.Now().Year(), time.Now().Month())
				}
				Expect(object.Key).To(Equal(aws.String(expectedFileName)))
			})

			It("returns the file", func() {
				Expect(usage).NotTo(BeEmpty())
				Expect(string(usage)).To(Equal(dailyUsageResponse))
			})
		})

		Context("when AWS returns an error", func() {
			BeforeEach(func() {
				s3Client.GetObjectReturns(nil, errors.New("request error"))
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("request error"))
			})
		})
	})

	Describe("CalculateDailyUsages", func() {
		var (
			originalReports  datamodels.Reports
			populatedReports datamodels.Reports
			err              error
		)

		BeforeEach(func() {
			originalReports = datamodels.Reports{
				datamodels.Report{
					AccountName:   "some-account-name",
					AccountNumber: "some-account-number",
					ServiceType:   "some-service-type",
					UsageQuantity: 10,
					Cost:          100,
					IAAS:          "AWS",
					Day:           1,
					Month:         "February",
					Year:          1337,
					Region:        "my-region",
				},
				datamodels.Report{
					AccountName:   "some-account-name",
					AccountNumber: "some-account-number",
					ServiceType:   "some-other-service-type",
					UsageQuantity: 9,
					Cost:          20,
					IAAS:          "AWS",
					Day:           1,
					Month:         "February",
					Year:          1337,
					Region:        "my-region",
				},
			}
		})

		JustBeforeEach(func() {
			populatedReports, err = client.CalculateDailyUsages(originalReports)
		})

		Context("when a database is connected", func() {
			It("fetches cost and usage month to date for each report", func() {
				Expect(dbClient.GetUsageMonthToDateCallCount()).To(Equal(2))
				Expect(dbClient.GetUsageMonthToDateArgsForCall(0)).To(Equal(datamodels.ReportIdentifier{
					AccountNumber: "some-account-number",
					AccountName:   "some-account-name",
					ServiceType:   "some-service-type",
					IAAS:          "AWS",
					Day:           1,
					Month:         "February",
					Year:          1337,
					Region:        "my-region",
				}))
				Expect(dbClient.GetUsageMonthToDateArgsForCall(1)).To(Equal(datamodels.ReportIdentifier{
					AccountNumber: "some-account-number",
					AccountName:   "some-account-name",
					ServiceType:   "some-other-service-type",
					IAAS:          "AWS",
					Day:           1,
					Month:         "February",
					Year:          1337,
					Region:        "my-region",
				}))
			})

			Context("when the database returns found usage", func() {
				var (
					callCount      int
					returnedUsages []datamodels.UsageMonthToDate
				)

				BeforeEach(func() {
					originalReports = append(originalReports,
						datamodels.Report{
							AccountName:   "some-linked-account-name",
							AccountNumber: "some-linked-account-number",
							ServiceType:   "some-other-service-type",
							UsageQuantity: 1,
							Cost:          1,
						},
					)
					returnedUsages = []datamodels.UsageMonthToDate{
						datamodels.UsageMonthToDate{
							AccountNumber: "some-linked-account-number",
							AccountName:   "some-linked-account-name",
							Month:         time.Now().Month().String(),
							Year:          time.Now().Year(),
							ServiceType:   "some-service-type",
							UsageQuantity: 9,
							Cost:          90,
							Region:        "my-region",
							UnitOfMeasure: "GB",
							IAAS:          "AWS",
						},
						datamodels.UsageMonthToDate{
							AccountNumber: "some-account-number",
							AccountName:   "some-account-name",
							Month:         time.Now().Month().String(),
							Year:          time.Now().Year(),
							ServiceType:   "some-other-service-type",
							UsageQuantity: 7,
							Cost:          19,
							Region:        "my-region",
							UnitOfMeasure: "GB",
							IAAS:          "AWS",
						},
					}
					dbClient.GetUsageMonthToDateStub = func(datamodels.ReportIdentifier) (datamodels.UsageMonthToDate, error) {
						if callCount > 1 {
							return datamodels.UsageMonthToDate{}, nil
						}
						retUsage := returnedUsages[callCount]
						callCount++
						return retUsage, nil
					}
				})

				AfterEach(func() {
					callCount = 0
				})

				It("returns the right number of reports", func() {
					Expect(populatedReports).To(HaveLen(3))
				})

				It("calculates daily usages when previous use is found", func() {
					Expect(populatedReports[0].UsageQuantity).To(Equal(float64(1)))
					Expect(populatedReports[0].Cost).To(Equal(float64(10)))
					Expect(populatedReports[1].UsageQuantity).To(Equal(float64(2)))
					Expect(populatedReports[1].Cost).To(Equal(float64(1)))
				})

				It("sets daily amounts to total amounts when no previous usage found", func() {
					Expect(populatedReports[2].Cost).To(Equal(originalReports[2].Cost))
					Expect(populatedReports[2].UsageQuantity).To(Equal(originalReports[2].UsageQuantity))
				})
			})

			Context("when the database fails", func() {
				BeforeEach(func() {
					dbClient.GetUsageMonthToDateReturns(datamodels.UsageMonthToDate{}, errors.New("some-error"))
				})

				It("errors", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("some-error"))
				})
			})
		})

	})
})

var dailyUsageResponse = `
			one, two, three, four, five
			sometimes, you, might, think, you, want json
			but really, we, know, you, want, CSV`
