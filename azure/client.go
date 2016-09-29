package azure

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
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

type DetailedUsageReport struct {
	CSV []byte
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
	azureMonthlyUsage, err := c.MonthlyUsageReport()
	if err != nil {
		c.log.Error("Failed to get Azure monthly usage")
		return datamodels.Reports{}, err
	}
	c.log.Debug("Got monthly Azure usage")

	usageReader := NewUsageReader(c.log, c.location)
	reports, err := usageReader.GenerateReports(azureMonthlyUsage.CSV)
	if err != nil {
		c.log.Error("Failed to parse Azure usage")
		return datamodels.Reports{}, err
	}

	return usageReader.Normalize(reports), nil
}

func (c Client) MonthlyUsageReport() (DetailedUsageReport, error) {
	csvBody, err := c.GetCSV()
	if err != nil {
		return DetailedUsageReport{}, err
	}
	return MakeDetailedUsageReport(csvBody), nil
}

func (c Client) GetCSV() ([]byte, error) {
	req, err := http.NewRequest("GET", strings.Join([]string{c.URL, "rest", c.enrollment, "usage-report?type=detail"}, "/"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", "bearer "+c.accessKey)
	req.Header.Add("api-version", "1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Azure responded with error: %s", resp.Status)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func MakeDetailedUsageReport(body []byte) DetailedUsageReport {
	csvLines := strings.SplitN(string(body), "\n", 3) // for azure the first two lines are garbage
	csvFirstTwoLinesRemoved := csvLines[2]
	csvStrippedTrailingComma := strings.Replace(csvFirstTwoLinesRemoved, ",\r\n", "\r\n", -1)
	return DetailedUsageReport{CSV: []byte(csvStrippedTrailingComma)}
}
