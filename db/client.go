package db

import (
	"database/sql"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"

	"github.com/challiwill/meteorologica/datamodels"
	"github.com/challiwill/meteorologica/errare"
)

//go:generate counterfeiter . DB

type DB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	QueryRow(string, ...interface{}) *sql.Row
	Close() error
	Ping() error
	Begin() (*sql.Tx, error)
}

type Client struct {
	Log  *logrus.Logger
	Conn DB
}

func NewClient(log *logrus.Logger, username, password, address, name string) (*Client, error) {
	if username == "" && password != "" {
		return nil, errare.NewCreationError("database client", "cannot have a database password with a username\n Please set the DB_PASSWORD environment variable")
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
	return "Multiple errors occurred"
}

func (c *Client) SaveReports(reports datamodels.Reports) error {
	c.Log.Debug("Entering db.SaveReports")
	defer c.Log.Debug("Returning db.SaveReports")

	var multiErr MultiErr
	for i, r := range reports {
		c.Log.Debugf("Saving report to database %d of %d...", i, len(reports))
		_, err := c.Conn.Exec(`
		INSERT IGNORE INTO iaas_billing
		(AccountNumber, AccountName, Day, Month, Year, ServiceType, UsageQuantity, Cost, Region, UnitOfMeasure, IAAS)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, r.AccountNumber, r.AccountName, r.Day, r.Month, r.Year, r.ServiceType, r.UsageQuantity, r.Cost, r.Region, r.UnitOfMeasure, r.IAAS)
		if err != nil {
			c.Log.Warn("Failed to save report to database: ", err.Error())
			multiErr.errs = append(multiErr.errs, err)
		}
	}

	if len(multiErr.errs) != 0 && len(multiErr.errs) == len(reports) {
		return multiErr
	}
	return nil
}

func (c *Client) GetUsageMonthToDate(id datamodels.ReportIdentifier) (datamodels.UsageMonthToDate, error) {
	c.Log.Debug("Entering db.GetUsageMonthToDate")
	defer c.Log.Debug("Returning db.GetUsageMonthToDate")

	usageToDate := datamodels.UsageMonthToDate{}
	err := c.Conn.QueryRow(`
		SELECT AccountNumber, AccountName, Month, Year, ServiceType, SUM(UsageQuantity), SUM(Cost), Region, UnitOfMeasure, IAAS
		FROM iaas_billing
		WHERE AccountName=?
		AND Month=?
		AND Year=?
		AND ServiceType=?
		AND Region=?
		AND IAAS=?`,
		id.AccountName, id.Month, id.Year, id.ServiceType, id.Region, id.IAAS).Scan(&usageToDate)
	if err == sql.ErrNoRows {
		err = nil
	}

	return usageToDate, err
}

func (c *Client) Close() error {
	c.Log.Debug("Entering db.Close")
	defer c.Log.Debug("Returning db.Close")

	return c.Conn.Close()
}

func (c *Client) Ping() error {
	c.Log.Debug("Entering db.Ping")
	defer c.Log.Debug("Returning db.Ping")

	return c.Conn.Ping()
}

func (c *Client) Migrate() error {
	c.Log.Debug("Entering db.Migrate")
	defer c.Log.Debug("Returning db.Migrate")

	migration, err := ioutil.ReadFile("migrations/iaas_billing.sql")
	if err != nil {
		return err
	}

	tx, err := c.Conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(migration))
	if err != nil {
		c.Log.Warn("Migration failed, rolling back...")
		er := tx.Rollback()
		if er != nil {
			c.Log.Warn("Rollback failed: ", err.Error())
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		c.Log.Warn("Migration failed, rolling back...")
		err := tx.Rollback()
		if err != nil {
			c.Log.Warn("Rollback failed: ", err.Error())
		}
		return err
	}

	return nil
}
