package db

import (
	"database/sql"

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

func NewClientWith(log *logrus.Logger, conn DB) *Client {
	return &Client{
		Log:  log,
		Conn: conn,
	}
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
		if i%1000 == 0 {
			c.Log.Debugf("Saving report to database %d of %d...", i, len(reports))
		}
		_, err := c.Conn.Exec(`
		INSERT INTO resource_billing
		(id, account_number, account_name, day, month, year, service_type, region, resource, usage_quantity, unit_of_measure, cost)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, r.ID, r.AccountNumber, r.AccountName, r.Day, r.Month, r.Year, r.ServiceType, r.Region, r.Resource, r.UsageQuantity, r.UnitOfMeasure, r.Cost)
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

	var (
		accountName   sql.NullString
		region        sql.NullString
		unitOfMeasure sql.NullString
	)

	usageToDate := datamodels.UsageMonthToDate{}
	err := c.Conn.QueryRow(`
		SELECT account_number, account_name, month, year, service_type, SUM(usage_quantity), SUM(cost), region, unit_of_measure, resource
		FROM resource_billing
		WHERE account_number=?
		AND month=?
		AND year=?
		AND service_type=?
		AND region=?
		AND resource=?`,
		id.AccountNumber, id.Month, id.Year, id.ServiceType, id.Region, id.Resource).Scan(
		&usageToDate.AccountNumber,
		&accountName,
		&usageToDate.Month,
		&usageToDate.Year,
		&usageToDate.ServiceType,
		&usageToDate.UsageQuantity,
		&usageToDate.Cost,
		&region,
		&unitOfMeasure,
		&usageToDate.Resource,
	)
	if err == sql.ErrNoRows {
		return datamodels.UsageMonthToDate{}, nil
	}

	if accountName.Valid {
		usageToDate.AccountName = accountName.String
	}
	if region.Valid {
		usageToDate.Region = region.String
	}
	if unitOfMeasure.Valid {
		usageToDate.UnitOfMeasure = unitOfMeasure.String
	}

	return usageToDate, nil
}

func (c *Client) Close() error {
	c.Log.Debug("Entering db.Close")
	defer c.Log.Debug("Returning db.Close")

	return c.Conn.Close()
}
