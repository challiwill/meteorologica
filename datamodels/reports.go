package datamodels

type ReportIdentifier struct {
	AccountNumber string
	AccountName   string
	ServiceType   string
	Day           int
	Month         string
	Year          int
	IAAS          string
	Region        string
}

type UsageMonthToDate struct {
	AccountNumber string
	AccountName   string
	Month         string
	Year          int
	ServiceType   string
	UsageQuantity float64
	Cost          float64
	Region        string
	UnitOfMeasure string
	IAAS          string
}

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
