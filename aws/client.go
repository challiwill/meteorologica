package aws

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/challiwill/meteorologica/csv"
	"github.com/challiwill/meteorologica/datamodels"
)

//go:generate counterfeiter . S3Client

type S3Client interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

type Client struct {
	Bucket        string
	AccountNumber string
	Region        string
	s3            S3Client
	log           *logrus.Logger
	location      *time.Location
}

func NewClient(log *logrus.Logger, location *time.Location, az, bucketName, accountNumber string, s3Client S3Client) *Client {
	return &Client{
		Bucket:        bucketName,
		AccountNumber: accountNumber,
		Region:        az,
		s3:            s3Client,
		log:           log,
		location:      location,
	}
}

func (c Client) Name() string {
	return "AWS"
}

func (c Client) GetNormalizedUsage() (datamodels.Reports, error) {
	c.log.Info("Getting Monthly AWS Usage...")
	c.log.Debug("Entering aws.GetNormalizedUsage")
	defer c.log.Debug("Returning aws.GetNormalizedUsage")

	awsMonthlyUsage, err := c.GetBillingData()
	if err != nil {
		return datamodels.Reports{}, fmt.Errorf("Failed to get monthly AWS billing data: %s", err.Error())
	}
	c.log.Debug("Got Monthly AWS usage")

	readerCleaner, err := csv.NewReaderCleaner(bytes.NewReader(awsMonthlyUsage), 29)
	if err != nil {
		return nil, fmt.Errorf("Failed to Read or Clean AWS reports: %s", err.Error())
	}
	reports := []*Usage{}
	err = csv.GenerateReports(readerCleaner, &reports)
	if err != nil {
		return datamodels.Reports{}, fmt.Errorf("Failed to Generate Reports for AWS: %s", err.Error())
	}

	return NewNormalizer(c.log, c.location, c.Region).Normalize(reports), nil
}

func (c Client) GetBillingData() ([]byte, error) {
	c.log.Debug("Entering aws.GetBillingData")
	defer c.log.Debug("Returning aws.GetBillingData")

	objectInput := &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(c.monthlyBillingFileName()),
	}
	resp, err := c.s3.GetObject(objectInput)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c Client) monthlyBillingFileName() string {
	year, month, _ := time.Now().In(c.location).Date()
	monthStr := padMonth(month)
	return url.QueryEscape(strings.Join([]string{c.AccountNumber, "aws", "billing", "csv", strconv.Itoa(year), monthStr}, "-") + ".csv")
}

func padMonth(month time.Month) string {
	m := strconv.Itoa(int(month))
	if month < 10 {
		return "0" + m
	}
	return m
}
