package gcp

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/challiwill/meteorologica/datamodels"
	"github.com/gocarina/gocsv"
)

type Usage struct {
	AccountID                    string `csv:"Account ID"`
	LineItem                     string `csv:"Line Item"`
	StartTime                    string `csv:"Start Time"`
	EndTime                      string `csv:"End Time"`
	Project                      string `csv:"Project"`
	Measurement1                 string `csv:"Measurement1"`
	Measurement1TotalConsumption string `csv:"Measurement1 Total Consumption"`
	Measurement1Units            string `csv:"Measurement1 Units"`
	Credit1                      string `csv:"Credit1"`
	Credit1Amount                string `csv:"Credit1 Amount"`
	Credit1Currency              string `csv:"Credit1 Currency"`
	Cost                         string `csv:"Cost"`
	Currency                     string `csv:"Currency"`
	ProjectNumber                string `csv:"Project Number"`
	ProjectID                    string `csv:"Project ID"`
	ProjectName                  string `csv:"Project Name"`
	ProjectLabels                string `csv:"Project Labels"`
	Description                  string `csv:"Description"`
}

type UsageReader struct {
	UsageReports []*Usage
}

func NewUsageReader(monthlyUsage *os.File) (*UsageReader, error) {
	reports, err := generateReports(monthlyUsage)
	if err != nil {
		return nil, err
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
		t, err := time.Parse("2006-01-02T15:04:05-07:00", usage.StartTime)
		if err != nil {
			fmt.Printf("Could not parse time '%s', defaulting to today '%s'\n", usage.StartTime, time.Now().String())
			t = time.Now()
		}
		reports = append(reports, datamodels.Report{
			AccountNumber: usage.ProjectNumber,
			AccountName:   usage.ProjectID,
			Day:           strconv.Itoa(t.Day()),
			Month:         t.Month().String(),
			Year:          strconv.Itoa(t.Year()),
			ServiceType:   usage.Description,
			UsageQuantity: usage.Measurement1TotalConsumption,
			Cost:          usage.Cost,
			Region:        "",
			UnitOfMeasure: usage.Measurement1Units,
			IAAS:          "GCP",
		})
	}
	return reports
}