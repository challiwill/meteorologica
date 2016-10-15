package azure

import (
	"hash/fnv"
	"strconv"
)

type Usage struct {
	AccountOwnerId         string  `csv:"AccountOwnerId"`
	AccountName            string  `csv:"Account Name"`
	ServiceAdministratorId string  `csv:"ServiceAdministratorId"`
	SubscriptionId         string  `csv:"SubscriptionId"`
	SubscriptionGuid       string  `csv:"SubscriptionGuid"`
	SubscriptionName       string  `csv:"Subscription Name"`
	Date                   string  `csv:"Date"`
	Month                  int     `csv:"Month"`
	Day                    int     `csv:"Day"`
	Year                   int     `csv:"Year"`
	Product                string  `csv:"Product"`
	MeterID                string  `csv:"Meter ID"`
	MeterCategory          string  `csv:"Meter Category"`
	MeterSubCategory       string  `csv:"Meter Sub-Category"`
	MeterRegion            string  `csv:"Meter Region"`
	MeterName              string  `csv:"Meter Name`
	ConsumedQuantity       float64 `csv:"Consumed Quantity"`
	ResourceRate           float64 `csv:"ResourceRate"`
	ExtendedCost           float64 `csv:"ExtendedCost"`
	ResourceLocation       string  `csv:"Resource Location"`
	ConsumedService        string  `csv:"Consumed Service"`
	InstanceID             string  `csv:"Instance ID"`
	ServiceInfo1           string  `csv:"ServiceInfo1"`
	ServiceInfo2           string  `csv:"ServiceInfo2"`
	AdditionalInfo         string  `csv:"AdditionalInfo"`
	Tags                   string  `csv:"Tags"`
	StoreServiceIdentifier string  `csv:"Store Service Identifier"`
	DepartmentName         string  `csv:"Department Name"`
	CostCenter             string  `csv:"Cost Center"`
	UnitOfMeasure          string  `csv:"Unit Of Measure"`
	ResourceGroup          string  `csv:"Resource Group"`
}

func (u Usage) Hash() string {
	h := fnv.New32a()
	h.Write([]byte(u.SubscriptionGuid + strconv.Itoa(u.Year) + strconv.Itoa(u.Month) + strconv.Itoa(u.Day) + u.ConsumedService + u.MeterRegion + IAAS))
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}
