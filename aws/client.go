package aws

import (
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
)

type DetailedUsageReport struct {
	CSV []byte
}

type Client struct {
	Bucket        string
	AccountNumber string
	s3            *s3.S3
}

func NewClient(bucketName, accountNumber string, config client.ConfigProvider) *Client {
	return &Client{
		Bucket:        bucketName,
		AccountNumber: accountNumber,
		s3:            s3.New(config),
	}
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
	return DetailedUsageReport{CSV: body}, nil
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
