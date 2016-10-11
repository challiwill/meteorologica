package db

import "github.com/challiwill/meteorologica/datamodels"

type NullClient struct{}

func NewNullClient() *NullClient {
	return &NullClient{}
}

func (c *NullClient) GetUsageMonthToDate(datamodels.ReportIdentifier) (datamodels.UsageMonthToDate, error) {
	return datamodels.UsageMonthToDate{}, nil
}
