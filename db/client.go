package db

import (
	"database/sql"
	"errors"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"

	"github.com/challiwill/meteorologica/datamodels"
)

//go:generate counterfeiter . DB

type DB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Close() error
}

type Client struct {
	Log  *logrus.Logger
	Conn DB
}

func NewClient(log *logrus.Logger, username, password, address, name string) (*Client, error) {
	if username == "" && password != "" {
		return nil, errors.New("Cannot have a database password without a username. Please set the DB_PASSWORD environment variable.")
	}

	conn, err := sql.Open("mysql", username+":"+password+"@"+"tcp("+address+")/"+name)
	if err != nil {
		return nil, err
	}

	return &Client{
		Log:  log,
		Conn: conn,
	}, nil
}

type MultiErr struct {
	errs []error
}

func (e MultiErr) Error() string {
	errString := "Multiple errors occurred: \n"
	for _, er := range e.errs {
		errString = errString + er.Error()
	}

	return errString
}

func (c *Client) SaveReports(reports datamodels.Reports) error {
	if len(reports) == 0 {
		return errors.New("No reports to save")
	}
	var multiErr MultiErr
	for i, r := range reports {
		c.Log.Debugf("Saving report to database %d of %d...", i, len(reports))
		_, err := c.Conn.Exec(`
		INSERT IGNORE INTO iaas_billing
		(AccountNumber, AccountName, Day, Month, Year, ServiceType, UsageQuantity, Cost, Region, UnitOfMeasure, IAAS)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, r.AccountNumber, r.AccountName, r.Day, r.Month, r.Year, r.ServiceType, r.UsageQuantity, r.Cost, r.Region, r.UnitOfMeasure, r.IAAS)
		if err != nil {
			c.Log.Warn("Failed to save report to database")
			multiErr.errs = append(multiErr.errs, err)
		}
	}

	if len(multiErr.errs) == len(reports) {
		return multiErr
	}
	return nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
