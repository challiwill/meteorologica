package gcp

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/datamodels"
)

type Normalizer struct {
	log      *logrus.Logger
	location *time.Location
}

func NewNormalizer(log *logrus.Logger, location *time.Location) *Normalizer {
	return &Normalizer{
		log:      log,
		location: location,
	}
}

func (ur *Normalizer) Normalize(usageReports []*Usage) datamodels.Reports {
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
