package aws

import (
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/challiwill/meteorologica/datamodels"
)

type DetailedUsageReport struct {
	CSV []byte
}

type Client struct {
	Bucket        string
	AccountNumber string
	Region        string
	s3            *s3.S3
	log           *logrus.Logger
}

func NewClient(log *logrus.Logger, az, bucketName, accountNumber string, config client.ConfigProvider) *Client {
	return &Client{
		Bucket:        bucketName,
		AccountNumber: accountNumber,
		Region:        az,
		s3:            s3.New(config),
		log:           log,
	}
}

func (c Client) Name() string {
	return "AWS"
}

func (c Client) GetNormalizedUsage() (datamodels.Reports, error) {
	c.log.Info("Getting Monthly AWS Usage...")
	awsMonthlyUsage, err := c.MonthlyUsageReport()
	if err != nil {
		c.log.Error("Failed to get AWS monthly usage: ", err)
		return datamodels.Reports{}, err
	}

	c.log.Debug("Got Monthly AWS Usage")
	err = ioutil.WriteFile("aws.csv", awsMonthlyUsage.CSV, os.ModePerm)
	if err != nil {
		c.log.Error("Failed to save AWS Usage to file")
		return datamodels.Reports{}, err
	}
	c.log.Debug("AWS Usage saved to aws.csv")

	awsDataFile, err := os.OpenFile("aws.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		c.log.Error("Failed to open AWS file")
		return datamodels.Reports{}, err
	}
	defer awsDataFile.Close()
	usageReader, err := NewUsageReader(c.log, awsDataFile, c.Region)
	if err != nil {
		c.log.Error("Failed to parse AWS file")
		return datamodels.Reports{}, err
	}
	defer os.Remove("aws.csv") // only remove if succeeded to parse
	return usageReader.Normalize(), nil
}

func (c Client) MonthlyUsageReport() (DetailedUsageReport, error) {
	objectInput := &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(c.monthlyBillingFileName()),
	}
	resp, err := c.s3.GetObject(objectInput)
	if err != nil {
		return DetailedUsageReport{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DetailedUsageReport{}, err
	}

	csvLines := strings.SplitN(string(body), "\n", -1)
	for csvLines[len(csvLines)-1] == "" { // remove empty lines
		csvLines = csvLines[:len(csvLines)-1]
	}
	csvLines = csvLines[:len(csvLines)-1] // the last filled in line is a warning
	csvStr := strings.Join(csvLines, "\n")
	return DetailedUsageReport{CSV: []byte(csvStr)}, nil
}

func (c Client) monthlyBillingFileName() string {
	year, month, _ := time.Now().Date()
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
