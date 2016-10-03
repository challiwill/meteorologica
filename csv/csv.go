package csv

import "github.com/gocarina/gocsv"

type CSV [][]string

func GenerateReports(monthlyUsageReader *ReaderCleaner, usages interface{}) error {
	err := gocsv.UnmarshalCSV(monthlyUsageReader, usages)
	if err != nil {
		return err
	}
	return nil
}
