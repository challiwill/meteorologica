package aws

import (
	"strconv"
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
		if usage.ProductCode == "" { // skip lines that tally up accounts
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

		t := time.Now().In(n.location)
		reports = append(reports, datamodels.Report{
			AccountNumber: accountID,
			AccountName:   accountName,
			Day:           strconv.Itoa(t.Day() - 1),
			Month:         t.Month().String(),
			Year:          strconv.Itoa(t.Year()),
			ServiceType:   usage.ProductName,
			UsageQuantity: usage.UsageQuantity,
			Cost:          usage.TotalCost,
			Region:        n.az,
			UnitOfMeasure: "",
			IAAS:          "AWS",
		})
	}
	return reports
}
