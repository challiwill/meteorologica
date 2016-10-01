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
		c.log.Error("Failed to get Azure monthly usage")
		return datamodels.Reports{}, err
	}
	c.log.Debug("Got monthly Azure usage")

	usageReader, err := csv.NewReaderCleaner(bytes.NewReader(azureMonthlyUsage), 30)
	if err != nil {
		return nil, err
	}
	reports := []*Usage{}
	err = csv.GenerateReports(usageReader, reports)
	if err != nil {
		c.log.Error("Failed to parse Azure usage")
		return datamodels.Reports{}, err
	}

	return NewNormalizer(c.log, c.location).Normalize(reports), nil
}

func (c Client) GetBillingData() ([]byte, error) {
	req, err := http.NewRequest("GET", strings.Join([]string{c.URL, "rest", c.enrollment, "usage-report?type=detail"}, "/"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", "bearer "+c.accessKey)
	req.Header.Add("api-version", "1.0")

	resp, err := c.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Azure responded with error: %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
