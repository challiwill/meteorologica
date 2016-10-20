package db

import (
	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/datamodels"
)

type NullClient struct {
	log *logrus.Logger
}

func NewNullClient(log *logrus.Logger) *NullClient {
	return &NullClient{log: log}
}

func (c *NullClient) GetUsageMonthToDate(datamodels.ReportIdentifier) (datamodels.UsageMonthToDate, error) {
	c.log.Debug("No-op: using db.NullClient")
	return datamodels.UsageMonthToDate{}, nil
}

func (c *NullClient) SaveReports(datamodels.Reports) error {
	c.log.Debug("No-op: using db.NullClient")
	return nil
}

func (c *NullClient) Close() error {
	c.log.Debug("No-op: using db.NullClient")
	return nil
}
