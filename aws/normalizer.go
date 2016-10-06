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
		if isNotLineItem(usage) {
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
		usageQuantity, err := strconv.ParseFloat(usage.UsageQuantity, 64)
		if err != nil {
			n.log.Warnf("Usage Quantity '%s' invalid, setting to 0", usage.UsageQuantity)
			usageQuantity = 0
		}
		cost, err := strconv.ParseFloat(usage.TotalCost, 64)
		if err != nil {
			n.log.Warnf("Total Cost '%s' invalid, setting to 0", usage.TotalCost)
			cost = 0
		}
		t := time.Now().In(n.location)
		reports = append(reports, datamodels.Report{
			AccountNumber: accountID,
			AccountName:   accountName,
			Day:           t.Day() - 1,
			Month:         t.Month().String(),
			Year:          t.Year(),
			ServiceType:   usage.ProductName,
			UsageQuantity: usageQuantity,
			Cost:          cost,
			Region:        n.az,
			UnitOfMeasure: "",
			IAAS:          "AWS",
		})
	}
	return reports
}

func isNotLineItem(usage *Usage) bool {
	return usage.RecordType != "PayerLineItem" && usage.RecordType != "LinkedLineItem"
}
