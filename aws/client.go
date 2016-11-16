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
	"github.com/challiwill/meteorologica/calendar"
	"github.com/challiwill/meteorologica/csv"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/challiwill/meteorologica/errare"
)

var IAAS = "AWS"

//go:generate counterfeiter . S3Client

type S3Client interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

//go:generate counterfeiter . ReportsDatabase

type ReportsDatabase interface {
	GetUsageMonthToDate(datamodels.ReportIdentifier) (datamodels.UsageMonthToDate, error)
}

type Client struct {
	Bucket        string
	AccountNumber int64
	Region        string
	s3            S3Client
	log           *logrus.Logger
	location      *time.Location
	db            ReportsDatabase
}

func NewClient(log *logrus.Logger, location *time.Location, az, bucketName string, accountNumber int64, s3Client S3Client, db ReportsDatabase) *Client {
	return &Client{
		Bucket:        bucketName,
		AccountNumber: accountNumber,
		Region:        az,
		s3:            s3Client,
		log:           log,
		location:      location,
		db:            db,
	}
}

func (c Client) Name() string {
	return IAAS
}

func (c Client) GetNormalizedUsage() (datamodels.Reports, error) {
	c.log.Info("Getting Monthly AWS Usage...")
	c.log.Debug("Entering aws.GetNormalizedUsage")
	defer c.log.Debug("Returning aws.GetNormalizedUsage")

	awsMonthlyUsage, err := c.GetBillingData()
	if err != nil {
		c.log.Error("Failed to get AWS monthly usage")
		return datamodels.Reports{}, errare.NewRequestError(err, "AWS")
	}
	c.log.Debug("Got Monthly AWS usage")

	readerCleaner, err := csv.NewReaderCleaner(bytes.NewReader(awsMonthlyUsage), 29)
	if err != nil {
		return datamodels.Reports{}, csv.NewReadCleanError("AWS", err)
	}
	reports := []*Usage{}
	err = csv.GenerateReports(readerCleaner, &reports)
	if err != nil {
		return datamodels.Reports{}, csv.NewReportParseError("AWS", err)
	}

	normalizer := NewNormalizer(c.log, c.location, c.Region)
	normalizedReports := normalizer.Normalize(reports)
	normalizedReports = datamodels.ConsolidateReports(normalizedReports)
	normalizedReports, err = c.CalculateDailyUsages(normalizedReports)
	if err != nil {
		return datamodels.Reports{}, err
	}

	return normalizedReports, nil
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

func (c Client) CalculateDailyUsages(reports datamodels.Reports) (datamodels.Reports, error) {
	for i, report := range reports {
		usageToDate, err := c.db.GetUsageMonthToDate(datamodels.ReportIdentifier{
			AccountNumber: report.AccountNumber,
			AccountName:   report.AccountName,
			ServiceType:   report.ServiceType,
			Day:           report.Day,
			Month:         report.Month,
			Year:          report.Year,
			Resource:      report.Resource,
			Region:        report.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to get usage month-to-date: %s", err.Error())
		}

		reports[i].UsageQuantity = reports[i].UsageQuantity - usageToDate.UsageQuantity
		reports[i].Cost = reports[i].Cost - usageToDate.Cost
	}
	return reports, nil
}

func (c Client) monthlyBillingFileName() string {
	year, month, _ := calendar.YesterdaysDate(c.location)
	monthStr := calendar.PadMonth(month)
	return url.QueryEscape(strings.Join([]string{strconv.FormatInt(c.AccountNumber, 10), "aws", "billing", "csv", strconv.Itoa(year), monthStr}, "-") + ".csv")
}
