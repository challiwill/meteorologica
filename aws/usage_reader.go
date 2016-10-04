package aws

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/gocarina/gocsv"
)

type Usage struct {
	InvoiceID              string `csv:"InvoiceID"`
	PayerAccountId         string `csv:"PayerAccountId"`
	LinkedAccountId        string `csv:"LinkedAccountId"`
	RecordType             string `csv:"RecordType"`
	RecordID               string `csv:"RecordID"`
	BillingPeriodStartDate string `csv:"BillingPeriodStartDate"`
	BillingPeriodEndDate   string `csv:"BillingPeriodEndDate"`
	InvoiceDate            string `csv:"InvoiceDate"`
	PayerAccountName       string `csv:"PayerAccountName"`
	LinkedAccountName      string `csv:"LinkedAccountName"`
	TaxationAddress        string `csv:"TaxationAddress"`
	PayerPONumber          string `csv:"PayerPONumber"`
	ProductCode            string `csv:"ProductCode"`
	ProductName            string `csv:"ProductName"`
	SellerOfRecord         string `csv:"SellerOfRecord"`
	UsageType              string `csv:"UsageType"`
	Operation              string `csv:"Operation"`
	RateId                 string `csv:"RateId"`
	ItemDescription        string `csv:"ItemDescription"`
	UsageStartDate         string `csv:"UsageStartDate"`
	UsageEndDate           string `csv:"UsageEndDate"`
	UsageQuantity          string `csv:"UsageQuantity"`
	BlendedRate            string `csv:"BlendedRate"`
	CurrencyCode           string `csv:"CurrencyCode"`
	CostBeforeTax          string `csv:"CostBeforeTax"`
	Credits                string `csv:"Credits"`
	TaxAmount              string `csv:"TaxAmount"`
	TaxType                string `csv:"TaxType"`
	TotalCost              string `csv:"TotalCost"`
	DailySpend             string `csv:"-"`
	AvailabilityZone       string `csv:"-"`
}

type UsageReader struct {
	az       string
	log      *logrus.Logger
	location *time.Location
}

func NewUsageReader(log *logrus.Logger, location *time.Location, az string) *UsageReader {
	return &UsageReader{
		az:       az,
		log:      log,
		location: location,
	}
}

func (ur *UsageReader) GenerateReports(monthlyUsage []byte) ([]*Usage, error) {
	usages := []*Usage{}
	err := gocsv.UnmarshalBytes(monthlyUsage, &usages)
	if err != nil {
		return nil, err
	}

	for _, usage := range usages {
		usage.AvailabilityZone = ur.az
	}
	return usages, nil
}

func (ur *UsageReader) Normalize(usageReports []*Usage) datamodels.Reports {
	var reports datamodels.Reports
	for _, usage := range usageReports {
		if usage.ProductCode == "" { // skip lines that total up accounts
			continue
		}
		accountName := usage.LinkedAccountName
		if accountName == "" {
			accountName = usage.PayerAccountName
		}
		accountID := usage.LinkedAccountId
		if accountID == "" {
			accountID = usage.PayerAccountId
		}

		ur.log.Debug("Using today's date as the date of retrieval for the AWS billing data")
		t := time.Now().In(ur.location)
		reports = append(reports, datamodels.Report{
			AccountNumber: accountID,
			AccountName:   accountName,
			Day:           strconv.Itoa(t.Day() - 1),
			Month:         t.Month().String(),
			Year:          strconv.Itoa(t.Year()),
			ServiceType:   usage.ProductName,
			UsageQuantity: usage.UsageQuantity,
			Cost:          usage.TotalCost,
			Region:        usage.AvailabilityZone,
			UnitOfMeasure: "",
			IAAS:          "AWS",
		})
	}
	return reports
}
