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

func (ur *Normalizer) Normalize(usageReports []*Usage) datamodels.Reports {
	var reports datamodels.Reports
	for _, usage := range usageReports {
		if usage.ProductCode == "" { // skip lines that total up accounts
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

		ur.log.Debug("Using today's date as the date of retrieval for the AWS billing data")
		t := time.Now().In(ur.location)
		reports = append(reports, datamodels.Report{
			AccountNumber: accountID,
			AccountName:   accountName,
			Day:           strconv.Itoa(t.Day() - 1),
			Month:         t.Month().String(),
			Year:          strconv.Itoa(t.Year()),
			ServiceType:   usage.ProductName,
			UsageQuantity: usage.UsageQuantity,
			Cost:          usage.TotalCost,
			Region:        ur.az,
			UnitOfMeasure: "",
			IAAS:          "AWS",
		})
	}
	return reports
}
