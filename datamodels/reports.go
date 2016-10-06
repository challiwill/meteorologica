package datamodels

type Report struct {
	AccountNumber string  `csv:"Account Number"`
	AccountName   string  `csv:"Account Name"`
	Day           int     `csv:"Day"`
	Month         string  `csv:"Month"`
	Year          int     `csv:"Year"`
	ServiceType   string  `csv:"Service Type"`
	UsageQuantity float64 `csv:"UsageQuantity"`
	Cost          float64 `csv:"Cost"`
	Region        string  `csv:"Region"`
	UnitOfMeasure string  `csv:"Unit Of Measurement"`
	IAAS          string  `csv:"IAAS"`
}

type Reports []Report
