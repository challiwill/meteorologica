package azure

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"bytes"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/csv"
	"github.com/challiwill/meteorologica/datamodels"
)

type UsageReport struct {
	Month      string `json:"Month"`
	DetailLink string `json:"LinkToDownloadDetailReport"`
}

type UsageReports struct {
	ContractVersion string        `json:"contract_version"`
	Months          []UsageReport `json:"AvailableMonths"`
}

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
	azureMonthlyUsage, err := c.GetBillingData()
	if err != nil {
		return datamodels.Reports{}, err
	}
	c.log.Debug("Got monthly Azure usage")

	readerCleaner, err := csv.NewReaderCleaner(bytes.NewReader(azureMonthlyUsage), 31)
	if err != nil {
		return nil, fmt.Errorf("Failed to Read or Clean Azure reports: %s", err.Error())
	}
	reports := []*Usage{}
	err = csv.GenerateReports(readerCleaner, &reports)
	if err != nil {
		return datamodels.Reports{}, fmt.Errorf("Failed to Generate Reports for Azure: %s", err.Error())
	}

	return NewNormalizer(c.log, c.location).Normalize(reports), nil
}

func (c Client) GetBillingData() ([]byte, error) {
	reqString := strings.Join([]string{c.URL, "rest", c.enrollment, "usage-report?type=detail"}, "/")
	c.log.Debug("Making Azure billing request to address: ", reqString)

	req, err := http.NewRequest("GET", reqString, nil)
	if err != nil {
		return nil, fmt.Errorf("Creating request for Azure failed: %s", err)
	}
	req.Header.Add("authorization", "bearer "+c.accessKey)
	req.Header.Add("api-version", "1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Making request to Azure failed: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Azure responded with error: %s", resp.Status)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
