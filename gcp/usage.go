package gcp

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
