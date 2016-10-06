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

func (n *Normalizer) Normalize(usageReports []*Usage) datamodels.Reports {
	n.log.Debug("Entering azure.Normalize")
	defer n.log.Debug("Returning azure.Normalize")

	n.log.Debug("Normalizing Azure data...")
	var reports datamodels.Reports
	for _, usage := range usageReports {
		month := time.Now().In(n.location).Month()
		m, err := strconv.Atoi(usage.Month)
		if err != nil || m < 1 || m > 12 {
			n.log.Warnf("%s month is invalid, defaulting to this %s", usage.Month, time.Now().In(n.location).Month().String())
		} else {
			month = time.Month(m)
		}

		day, err := strconv.Atoi(usage.Day)
		if err != nil || day < 1 || day > 31 {
			day = time.Now().In(n.location).Day()
			n.log.Warnf("%s day is invalid, defaulting to this %d", usage.Day, day)
		}
		year, err := strconv.Atoi(usage.Year)
		if err != nil {
			year = time.Now().In(n.location).Year()
			n.log.Warnf("%s year is invalid, defaulting to this %d", usage.Month, year)
		}

		quantity, err := strconv.ParseFloat(usage.ConsumedQuantity, 64)
		if err != nil {
			quantity = 0
			n.log.Warnf("consumed quantity '%s' is invalid, defaulting to 0", usage.ConsumedQuantity)
		}
		cost, err := strconv.ParseFloat(usage.ExtendedCost, 64)
		if err != nil {
			cost = 0
			n.log.Warnf("extended cost '%s' is invalid, defaulting to 0", usage.ExtendedCost)
		}

		reports = append(reports, datamodels.Report{
			AccountNumber: usage.SubscriptionGuid,
			AccountName:   usage.SubscriptionName,
			Day:           day,
			Month:         month.String(),
			Year:          year,
			ServiceType:   usage.ConsumedService,
			UsageQuantity: quantity,
			Cost:          cost,
			Region:        usage.MeterRegion,
			UnitOfMeasure: usage.UnitOfMeasure,
			IAAS:          "Azure",
		})
	}
	return reports
}
