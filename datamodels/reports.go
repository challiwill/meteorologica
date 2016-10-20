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
	consolidatedReports := make(map[string]Report)
	for _, r := range reports {
		if _, ok := consolidatedReports[r.ID]; ok {
			consolidatedReports[r.ID] = sumReports(consolidatedReports[r.ID], r)
			continue
		}
		consolidatedReports[r.ID] = r
	}

	consolidatedReportsSlice := Reports{}
	for _, v := range consolidatedReports {
		consolidatedReportsSlice = append(consolidatedReportsSlice, v)
	}
	return consolidatedReportsSlice
}

func sumReports(one Report, two Report) Report {
	one.UsageQuantity += two.UsageQuantity
	one.Cost += two.Cost
	return one
}
