package aws

import (
	"bytes"
	"errors"
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
	AccountNumber string
	Region        string
	s3            S3Client
	log           *logrus.Logger
	location      *time.Location
	db            ReportsDatabase
}

func NewClient(log *logrus.Logger, location *time.Location, az, bucketName, accountNumber string, s3Client S3Client, db ReportsDatabase) *Client {
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

	reports = c.ConsolidateReports(reports)
	reports, err = c.CalculateDailyUsages(reports)
	if err != nil {
		return datamodels.Reports{}, err
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

func (c Client) ConsolidateReports(reports []*Usage) []*Usage {
	consolidatedReports := []*Usage{}
	for _, r := range reports {
		if i, found := find(consolidatedReports, r); found {
			consolidatedReports[i] = sumReports(consolidatedReports[i], r)
			continue
		}
		consolidatedReports = append(consolidatedReports, r)
	}
	return consolidatedReports
}

// find returns a found report if account number and service type match
// TODO it should probably use datamodels.ReportIdentifier's

func find(haystack []*Usage, needle *Usage) (int, bool) {
	for i, h := range haystack {
		if h.PayerAccountName == needle.PayerAccountName &&
			h.LinkedAccountName == needle.LinkedAccountName &&
			h.ProductName == needle.ProductName {
			return i, true
		}
	}
	return 0, false
}

// sumReports only adds the additive fields that we currenlty care about. This
// might change in the future if we start storing more fields in the database

func sumReports(one *Usage, two *Usage) *Usage {
	one.UsageQuantity += two.UsageQuantity
	one.TotalCost += two.TotalCost
	return one
}

func (c Client) CalculateDailyUsages(reports []*Usage) ([]*Usage, error) {
	// TODO this should become part of the normalizer in some ways (like a
	// NormalizeReportIdentifier() function)
	if c.db == nil {
		return nil, errors.New("no database connected")
	}
	for i, report := range reports {
		accountName := report.LinkedAccountName
		if accountName == "" {
			accountName = report.PayerAccountName
		}
		accountID := report.LinkedAccountId
		if accountID == "" {
			accountID = report.PayerAccountId
		}

		usageToDate, err := c.db.GetUsageMonthToDate(datamodels.ReportIdentifier{
			AccountNumber: accountID,
			AccountName:   accountName,
			ServiceType:   report.ProductName,
			Day:           time.Now().Day(),
			Month:         time.Now().Month().String(),
			Year:          time.Now().Year(),
			IAAS:          IAAS,
			Region:        c.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to get usage month-to-date: %s", err.Error())
		}

		reports[i].DailyUsage = reports[i].UsageQuantity - usageToDate.UsageQuantity
		reports[i].DailySpend = reports[i].TotalCost - usageToDate.Cost
	}
	return reports, nil
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
