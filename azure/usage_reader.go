package azure

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/gocarina/gocsv"
)

type Usage struct {
	AccountOwnerId         string `csv:"AccountOwnerId"`
	AccountName            string `csv:"Account Name"`
	ServiceAdministratorId string `csv:"ServiceAdministratorId"`
	SubscriptionId         string `csv:"SubscriptionId"`
	SubscriptionGuid       string `csv:"SubscriptionGuid"`
	SubscriptionName       string `csv:"Subscription Name"`
	Date                   string `csv:"Date"`
	Month                  string `csv:"Month"`
	Day                    string `csv:"Day"`
	Year                   string `csv:"Year"`
	Product                string `csv:"Product"`
	MeterID                string `csv:"Meter ID"`
	MeterCategory          string `csv:"Meter Category"`
	MeterSubCategory       string `csv:"Meter Sub-Category"`
	MeterRegion            string `csv:"Meter Region"`
	MeterName              string `csv:"Meter Name`
	ConsumedQuantity       string `csv:"Consumed Quantity"`
	ResourceRate           string `csv:"ResourceRate"`
	ExtendedCost           string `csv:"ExtendedCost"`
	ResourceLocation       string `csv:"Resource Location"`
	ConsumedService        string `csv:"Consumed Service"`
	InstanceID             string `csv:"Instance ID"`
	ServiceInfo1           string `csv:"ServiceInfo1"`
	ServiceInfo2           string `csv:"ServiceInfo2"`
	AdditionalInfo         string `csv:"AdditionalInfo"`
	Tags                   string `csv:"Tags"`
	StoreServiceIdentifier string `csv:"Store Service Identifier"`
	DepartmentName         string `csv:"Department Name"`
	CostCenter             string `csv:"Cost Center"`
	UnitOfMeasure          string `csv:"Unit Of Measure"`
	ResourceGroup          string `csv:"Resource Group"`
}

type UsageReader struct {
	UsageReports []*Usage
	log          *logrus.Logger
}

func NewUsageReader(log *logrus.Logger, monthlyUsage []byte) (*UsageReader, error) {
	reports, err := generateReports(monthlyUsage)
	if err != nil {
		return nil, err
	}
	return &UsageReader{
		UsageReports: reports,
		log:          log,
	}, nil
}

func generateReports(monthlyUsage []byte) ([]*Usage, error) {
	usages := []*Usage{}
	err := gocsv.UnmarshalBytes(monthlyUsage, &usages)
	if err != nil {
		return nil, err
	}
	return usages, nil
}

func (ur *UsageReader) Normalize() datamodels.Reports {
	var reports datamodels.Reports
	for _, usage := range ur.UsageReports {
		month := time.Now().Month()
		m, _ := strconv.Atoi(usage.Month)
		if m < 1 || m > 12 {
			ur.log.Warn("%s month is invalid, defaulting to this %s\n", usage.Month, time.Now().Month().String())
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
