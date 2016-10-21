package gcp

import (
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

func (n *Normalizer) Normalize(usageReports []*Usage) datamodels.Reports {
	n.log.Debug("Entering gcp.Normalize")
	defer n.log.Debug("Returning gcp.Normalize")

	var reports datamodels.Reports
	for _, usage := range usageReports {
		t, err := time.Parse("2006-01-02T15:04:05-07:00", usage.StartTime)
		if err != nil {
			n.log.Warnf("Could not parse time '%s', defaulting to today '%s'\n", usage.StartTime, time.Now().In(n.location).String())
			t = time.Now().In(n.location)
		}

		reports = append(reports, datamodels.Report{
			ID:            usage.Hash(),
			AccountNumber: usage.ProjectID,
			AccountName:   usage.ProjectName,
			Day:           t.Day(),
			Month:         t.Month(),
			Year:          t.Year(),
			ServiceType:   usage.Description,
			UsageQuantity: usage.Measurement1TotalConsumption,
			Cost:          usage.Cost,
			Region:        "",
			UnitOfMeasure: usage.Measurement1Units,
			Resource:      IAAS,
		})
	}
	return reports
}
