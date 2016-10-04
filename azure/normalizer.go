package azure

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
	ur.log.Debug("Normalizing Azure data...")
	var reports datamodels.Reports
	for _, usage := range usageReports {
		month := time.Now().In(ur.location).Month()
		m, _ := strconv.Atoi(usage.Month)
		if m < 1 || m > 12 {
			ur.log.Warnf("%s month is invalid, defaulting to this %s", usage.Month, time.Now().In(ur.location).Month().String())
		} else {
			month = time.Month(m)
		}
		reports = append(reports, datamodels.Report{
			AccountNumber: usage.SubscriptionGuid,
			AccountName:   usage.SubscriptionName,
			Day:           usage.Day,
			Month:         month.String(),
			Year:          usage.Year,
			ServiceType:   usage.ConsumedService,
			UsageQuantity: usage.ConsumedQuantity,
			Cost:          usage.ExtendedCost,
			Region:        usage.MeterRegion,
			UnitOfMeasure: usage.UnitOfMeasure,
			IAAS:          "Azure",
		})
	}
	return reports
}
