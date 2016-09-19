package aws

import (
	"os"

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
	AvailabilityZone       string `csv:"-"`
}

type UsageReader struct {
	UsageReports []*Usage
}

func NewUsageReader(monthlyUsage *os.File, az string) (*UsageReader, error) {
	reports, err := generateReports(monthlyUsage)
	if err != nil {
		return nil, err
	}
	for _, r := range reports {
		r.AvailabilityZone = az
	}
	return &UsageReader{
		UsageReports: reports,
	}, nil
}

func generateReports(monthlyUsage *os.File) ([]*Usage, error) {
	usages := []*Usage{}
	err := gocsv.UnmarshalFile(monthlyUsage, &usages)
	if err != nil {
		return nil, err
	}
	return usages, nil
}

func (ur *UsageReader) Normalize() datamodels.Reports {
	var reports datamodels.Reports
	for _, usage := range ur.UsageReports {
		accountName := usage.LinkedAccountName
		if accountName == "" {
			accountName = usage.PayerAccountName
		}
		accountID := usage.LinkedAccountId
		if accountID == "" {
			accountID = usage.PayerAccountId
		}
		reports = append(reports, datamodels.Report{
			AccountNumber: accountID,
			AccountName:   accountName,
			Day:           usage.BillingPeriodStartDate,
			Month:         usage.BillingPeriodStartDate,
			Year:          usage.BillingPeriodStartDate,
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
