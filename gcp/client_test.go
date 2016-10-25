package gcp_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	. "github.com/challiwill/meteorologica/gcp"
	"github.com/challiwill/meteorologica/gcp/gcpfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

//go:generate counterfeiter io.ReadCloser

var _ = Describe("Gcp", func() {
	var (
		client  *Client
		service *gcpfakes.FakeStorageService
		log     *logrus.Logger
	)

	BeforeEach(func() {
		service = new(gcpfakes.FakeStorageService)
		log = logrus.New()
		log.Out = NewBuffer()
		client = &Client{
			StorageService: service,
			BucketName:     "my-bucket",
			Log:            log,
			Location:       time.Now().Location(),
		}
	})

	Describe("Name", func() {
		It("returns the IAAS name", func() {
			Expect(client.Name()).To(Equal("GCP"))
		})
	})

	XDescribe("GetNormalizedUsage", func() {
		It("works", func() {})
	})

	XDescribe("GetBillingData", func() {
		It("works", func() {})
	})

	Describe("DailyUsageReport", func() {
		var (
			report []byte
			err    error
			day    int
		)

		JustBeforeEach(func() {
			report, err = client.DailyUsageReport(day)
		})

		Context("when the storage service returns a successful response", func() {
			var (
				dailyUsageResponse = `
one, two, three, four, five
sometimes, you, might, think, you, want json
but really, we, know, you, want, CSV
				`
			)

			BeforeEach(func() {
				day = 12
				readCloser := new(gcpfakes.FakeReadCloser)
				readCloser.ReadStub = func(p []byte) (int, error) {
					return copy(p, dailyUsageResponse), io.EOF
				}

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       readCloser,
				}
				service.DailyUsageReturns(resp, nil)
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("requests the given bucketname and appropriate file name", func() {
				Expect(service.DailyUsageCallCount()).To(Equal(1))
				bucketName, fileName := service.DailyUsageArgsForCall(0)
				Expect(bucketName).To(Equal("my-bucket"))
				expectedFileName := fmt.Sprintf("Billing-%d-%d-12.csv", time.Now().Year(), time.Now().Month())
				if time.Now().Month() < 10 {
					expectedFileName = fmt.Sprintf("Billing-%d-0%d-12.csv", time.Now().Year(), time.Now().Month())
				}
				Expect(fileName).To(Equal(expectedFileName))
			})

			It("returns the file", func() {
				Expect(report).NotTo(BeEmpty())
				Expect(string(report)).To(Equal(dailyUsageResponse))
			})

			Context("when the day is a single digit", func() {
				BeforeEach(func() {
					day = 2
				})

				It("requests the given bucketname and appropriate file name", func() {
					bucketName, fileName := service.DailyUsageArgsForCall(0)
					Expect(bucketName).To(Equal("my-bucket"))
					expectedFileName := fmt.Sprintf("Billing-%d-%d-02.csv", time.Now().Year(), time.Now().Month())
					if time.Now().Month() < 10 {
						expectedFileName = fmt.Sprintf("Billing-%d-0%d-02.csv", time.Now().Year(), time.Now().Month())
					}
					Expect(fileName).To(Equal(expectedFileName))
				})

			})
		})

		Context("when the storage service returns a failed http response", func() {
			BeforeEach(func() {
				resp := &http.Response{
					StatusCode: http.StatusInternalServerError,
				}
				service.DailyUsageReturns(resp, nil)
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("GCP responded with error"))
			})
		})

		Context("when the storage service returns an error", func() {
			BeforeEach(func() {
				service.DailyUsageReturns(nil, errors.New("some-error"))
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some-error"))
			})
		})
	})
})
