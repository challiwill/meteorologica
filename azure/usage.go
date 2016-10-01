package azure

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
