package azure

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
	n.log.Debug("Entering azure.Normalize")
	defer n.log.Debug("Returning azure.Normalize")

	n.log.Debug("Normalizing Azure data...")
	var reports datamodels.Reports
	for _, usage := range usageReports {
		reports = append(reports, datamodels.Report{
			ID:            usage.Hash(),
			AccountNumber: usage.SubscriptionGuid,
			AccountName:   usage.SubscriptionName,
			Day:           usage.Day,
			Month:         time.Month(usage.Month),
			Year:          usage.Year,
			ServiceType:   usage.ConsumedService,
			UsageQuantity: usage.ConsumedQuantity,
			Cost:          usage.ExtendedCost,
			Region:        usage.MeterRegion,
			UnitOfMeasure: usage.UnitOfMeasure,
			Resource:      IAAS,
		})
	}
	return reports
}
