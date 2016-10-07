package aws

type Usage struct {
	InvoiceID              string  `csv:"InvoiceID"`
	PayerAccountId         string  `csv:"PayerAccountId"`
	LinkedAccountId        string  `csv:"LinkedAccountId"`
	RecordType             string  `csv:"RecordType"`
	RecordID               string  `csv:"RecordID"`
	BillingPeriodStartDate string  `csv:"BillingPeriodStartDate"`
	BillingPeriodEndDate   string  `csv:"BillingPeriodEndDate"`
	InvoiceDate            string  `csv:"InvoiceDate"`
	PayerAccountName       string  `csv:"PayerAccountName"`
	LinkedAccountName      string  `csv:"LinkedAccountName"`
	TaxationAddress        string  `csv:"TaxationAddress"`
	PayerPONumber          string  `csv:"PayerPONumber"`
	ProductCode            string  `csv:"ProductCode"`
	ProductName            string  `csv:"ProductName"`
	SellerOfRecord         string  `csv:"SellerOfRecord"`
	UsageType              string  `csv:"UsageType"`
	Operation              string  `csv:"Operation"`
	RateId                 string  `csv:"RateId"`
	ItemDescription        string  `csv:"ItemDescription"`
	UsageStartDate         string  `csv:"UsageStartDate"`
	UsageEndDate           string  `csv:"UsageEndDate"`
	UsageQuantity          float64 `csv:"UsageQuantity"` // should get float
	BlendedRate            string  `csv:"BlendedRate"`
	CurrencyCode           string  `csv:"CurrencyCode"`
	CostBeforeTax          string  `csv:"CostBeforeTax"`
	Credits                string  `csv:"Credits"`
	TaxAmount              string  `csv:"TaxAmount"`
	TaxType                string  `csv:"TaxType"`
	TotalCost              float64 `csv:"TotalCost"` // should get float
	DailySpend             float64 `csv:"-"`         // should get float
}
