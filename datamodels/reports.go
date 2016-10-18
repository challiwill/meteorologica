package datamodels

import "time"

type ReportIdentifier struct {
	AccountNumber string
	AccountName   string
	ServiceType   string
	Day           int
	Month         time.Month
	Year          int
	Resource      string
	Region        string
}

type UsageMonthToDate struct {
	AccountNumber string
	AccountName   string
	Month         time.Month
	Year          int
	ServiceType   string
	UsageQuantity float64
	Cost          float64
	Region        string
	UnitOfMeasure string
	Resource      string
}

type Report struct {
	ID            string     `csv:"ID"`
	AccountNumber string     `csv:"Account Number"`
	AccountName   string     `csv:"Account Name"`
	Day           int        `csv:"Day"`
	Month         time.Month `csv:"Month"`
	Year          int        `csv:"Year"`
	ServiceType   string     `csv:"Service Type"`
	Region        string     `csv:"Region"`
	Resource      string     `csv:"Resource"`
	UsageQuantity float64    `csv:"Usage Quantity"`
	UnitOfMeasure string     `csv:"Unit Of Measurement"`
	Cost          float64    `csv:"Cost"`
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
		if h.ID == needle.ID {
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
