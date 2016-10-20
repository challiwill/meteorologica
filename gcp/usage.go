package gcp

import (
	"hash/fnv"
	"strconv"
	"time"
)

type Usage struct {
	AccountID                    string  `csv:"Account ID"`
	LineItem                     string  `csv:"Line Item"`
	StartTime                    string  `csv:"Start Time"`
	EndTime                      string  `csv:"End Time"`
	Project                      string  `csv:"Project"`
	Measurement1                 string  `csv:"Measurement1"`
	Measurement1TotalConsumption float64 `csv:"Measurement1 Total Consumption"`
	Measurement1Units            string  `csv:"Measurement1 Units"`
	Credit1                      string  `csv:"Credit1"`
	Credit1Amount                string  `csv:"Credit1 Amount"`
	Credit1Currency              string  `csv:"Credit1 Currency"`
	Cost                         float64 `csv:"Cost"`
	Currency                     string  `csv:"Currency"`
	ProjectNumber                string  `csv:"Project Number"`
	ProjectID                    string  `csv:"Project ID"`
	ProjectName                  string  `csv:"Project Name"`
	ProjectLabels                string  `csv:"Project Labels"`
	Description                  string  `csv:"Description"`
}

func (u Usage) Hash() string {
	h := fnv.New64a()
	h.Write([]byte(u.ProjectNumber + u.Description + IAAS))
	t, err := time.Parse("2006-01-02T15:04:05-07:00", u.StartTime)
	if err != nil {
		t = time.Now()
	}
	return strconv.FormatUint(uint64(h.Sum64()), 10) + strconv.Itoa(t.Year()) + strconv.Itoa(int(t.Month())) + strconv.Itoa(t.Day())
}
