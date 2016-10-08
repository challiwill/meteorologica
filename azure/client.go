package azure

import (
	"net/http"
	"strings"
	"time"

	"bytes"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/csv"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/challiwill/meteorologica/errare"
)

type Client struct {
	URL        string
	client     *http.Client
	accessKey  string
	enrollment string
	log        *logrus.Logger
	location   *time.Location
}

func NewClient(log *logrus.Logger, location *time.Location, serverURL, key, enrollment string) *Client {
	return &Client{
		URL:        serverURL,
		client:     new(http.Client),
		accessKey:  key,
		enrollment: enrollment,
		log:        log,
		location:   location,
	}
}

func (c Client) Name() string {
	return "Azure"
}

func (c Client) GetNormalizedUsage() (datamodels.Reports, error) {
	c.log.Info("Getting monthly Azure usage...")
	c.log.Debug("Entering azure.GetNormalizedUsage")
	defer c.log.Debug("Returning azure.GetNormalizedUsage")

	azureMonthlyUsage, err := c.GetBillingData()
	if err != nil {
		c.log.Error("Failed to get Azure monthly usage")
		return datamodels.Reports{}, err
	}
	c.log.Debug("Got monthly Azure usage")

	readerCleaner, err := csv.NewReaderCleaner(bytes.NewReader(azureMonthlyUsage), 31)
	if err != nil {
		return datamodels.Reports{}, csv.NewReadCleanError("Azure", err)
	}
	reports := []*Usage{}
	err = csv.GenerateReports(readerCleaner, &reports)
	if err != nil {
		return datamodels.Reports{}, csv.NewReportParseError("Azure", err)
	}

	return NewNormalizer(c.log, c.location).Normalize(reports), nil
}

func (c Client) GetBillingData() ([]byte, error) {
	c.log.Debug("Entering azure.GetBillingData")
	defer c.log.Debug("Returning azure.GetBillingData")

	reqString := strings.Join([]string{c.URL, "rest", c.enrollment, "usage-report?type=detail"}, "/")
	c.log.Debug("Making Azure billing request to address: ", reqString)

	req, err := http.NewRequest("GET", reqString, nil)
	if err != nil {
		return nil, errare.NewCreationError("Azure request", err.Error())
	}
	req.Header.Add("authorization", "bearer "+c.accessKey)
	req.Header.Add("api-version", "1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errare.NewRequestError(err, "Azure")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errare.NewResponseError(resp.Status, "Azure")
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
