package aws

import (
	"hash/fnv"
	"strconv"
	"time"
)

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
	UsageQuantity          float64 `csv:"UsageQuantity"`
	BlendedRate            string  `csv:"BlendedRate"`
	CurrencyCode           string  `csv:"CurrencyCode"`
	CostBeforeTax          string  `csv:"CostBeforeTax"`
	Credits                string  `csv:"Credits"`
	TaxAmount              string  `csv:"TaxAmount"`
	TaxType                string  `csv:"TaxType"`
	TotalCost              float64 `csv:"TotalCost"`
}

func (u Usage) Hash(az string) string {
	yr, mn, dy := time.Now().Date()
	h := fnv.New64a()
	h.Write([]byte(u.LinkedAccountId + u.ProductName + az + IAAS))
	return strconv.FormatUint(uint64(h.Sum64()), 10) + strconv.Itoa(yr) + strconv.Itoa(int(mn)) + strconv.Itoa(dy)
}
