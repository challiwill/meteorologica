package gcp

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
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
	log      *logrus.Logger
	location *time.Location
}

func NewUsageReader(log *logrus.Logger, location *time.Location) *UsageReader {
	return &UsageReader{
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
	return usages, nil
}

func (ur *UsageReader) Normalize(usageReports []*Usage) datamodels.Reports {
	var reports datamodels.Reports
	for _, usage := range usageReports {
		t, err := time.Parse("2006-01-02T15:04:05-07:00", usage.StartTime)
		if err != nil {
			ur.log.Warnf("Could not parse time '%s', defaulting to today '%s'\n", usage.StartTime, time.Now().In(ur.location).String())
			t = time.Now().In(ur.location)
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
