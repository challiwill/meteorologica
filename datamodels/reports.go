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

func ConsolidateReports(reports Reports) Reports {
	consolidatedReports := Reports{}
	for _, r := range reports {
		if i, found := find(consolidatedReports, r); found {
			consolidatedReports[i] = sumReports(consolidatedReports[i], r)
			continue
		}
		consolidatedReports = append(consolidatedReports, r)
	}
	return consolidatedReports
}

// find returns a found report if account number and service type match
// TODO it should probably use datamodels.ReportIdentifiers

func find(haystack Reports, needle Report) (int, bool) {
	for i, h := range haystack {
		if h.AccountNumber == needle.AccountNumber &&
			h.ServiceType == needle.ServiceType &&
			h.Region == needle.Region &&
			h.IAAS == needle.IAAS &&
			h.Day == needle.Day &&
			h.Month == needle.Month &&
			h.Year == needle.Year {
			return i, true
		}
	}
	return 0, false
}

func sumReports(one Report, two Report) Report {
	one.UsageQuantity += two.UsageQuantity
	one.Cost += two.Cost
	return one
}
