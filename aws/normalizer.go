package aws

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/datamodels"
)

type Normalizer struct {
	az       string
	log      *logrus.Logger
	location *time.Location
}

func NewNormalizer(log *logrus.Logger, location *time.Location, az string) *Normalizer {
	return &Normalizer{
		az:       az,
		log:      log,
		location: location,
	}
}

func (n *Normalizer) Normalize(usageReports []*Usage) datamodels.Reports {
	n.log.Debug("Entering aws.Normalize")
	defer n.log.Debug("Returnign aws.Normalize")

	var reports datamodels.Reports
	for _, usage := range usageReports {
		if isNotLineItem(usage) {
			continue
		}
		t := time.Now().In(n.location)
		reports = append(reports, datamodels.Report{
			ID:            usage.Hash(n.az),
			AccountNumber: usage.LinkedAccountId,
			AccountName:   usage.LinkedAccountName,
			Day:           t.Day() - 1,
			Month:         t.Month(),
			Year:          t.Year(),
			ServiceType:   usage.ProductName,
			UsageQuantity: usage.UsageQuantity,
			Cost:          usage.TotalCost,
			Region:        n.az,
			UnitOfMeasure: "",
			Resource:      IAAS,
		})
	}
	return reports
}

func isNotLineItem(usage *Usage) bool {
	return usage.RecordType != "LinkedLineItem"
}
