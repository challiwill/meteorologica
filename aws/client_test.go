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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

//go:generate counterfeiter io.ReadCloser

var _ = Describe("Client", func() {
	var (
		client   *Client
		log      *logrus.Logger
		s3Client *awsfakes.FakeS3Client
	)

	BeforeEach(func() {
		log = logrus.New()
		log.Out = NewBuffer()
		s3Client = new(awsfakes.FakeS3Client)
		client = NewClient(log, time.Now().Location(), "my-region", "my-bucket", "my-account-number", s3Client)
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

				expectedFileName := fmt.Sprintf("my-account-number-aws-billing-csv-%d-%d.csv", time.Now().Year(), time.Now().Month())
				if time.Now().Month() < 10 {
					expectedFileName = fmt.Sprintf("my-account-number-aws-billing-csv-%d-0%d.csv", time.Now().Year(), time.Now().Month())
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
})

var dailyUsageResponse = `
			one, two, three, four, five
			sometimes, you, might, think, you, want json
			but really, we, know, you, want, CSV`
