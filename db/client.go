package db

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"

	"github.com/challiwill/meteorologica/datamodels"
)

//go:generate counterfeiter . DB

type DB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Close() error
}

type Client struct {
	Conn DB
}

func NewClient(username, password, address, name string) (*Client, error) {
	if username == "" && password != "" {
		return nil, errors.New("Cannot have a database password without a username. Please set the DB_PASSWORD environment variable.")
	}

	conn, err := sql.Open("mysql", username+":"+password+"@"+address+"/"+name+"?reconnect=true")
	if err != nil {
		return nil, err
	}

	return &Client{
		Conn: conn,
	}, nil
}

func (c *Client) SaveReports(reports datamodels.Reports) error {
	if len(reports) == 0 {
		return errors.New("No reports to save")
	}
	for _, r := range reports {
		c.Conn.Exec(`
		INSERT INTO iaas_billing
		(AccountNumber, AccountName, Day, Month, Year, ServiceType, UsageQuantity, Cost, Region, UnitOfMeasure, IAAS)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, r.AccountNumber, r.AccountName, r.Day, r.Month, r.Year, r.ServiceType, r.UsageQuantity, r.Cost, r.Region, r.UnitOfMeasure, r.IAAS)
	}
	return nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
