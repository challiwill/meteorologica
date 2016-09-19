package datamodels

type Report struct {
	AccountNumber string `csv:"Account Number"`
	AccountName   string `csv:"Account Name"`
	Day           string `csv:"Day"`
	Month         string `csv:"Month"`
	Year          string `csv:"Year"`
	ServiceType   string `csv:"Service Type"`
	UsageQuantity string `csv:"UsageQuantity"`
	Cost          string `csv:"Cost"`
	Region        string `csv:"Region"`
	UnitOfMeasure string `csv:"Unit Of Measurement"`
	IAAS          string `csv:"IAAS"`
}

type Reports []Report
