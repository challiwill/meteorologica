package azure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
}

func NewClient(serverURL, key, enrollment string) *Client {
	return &Client{
		URL:        serverURL,
		client:     new(http.Client),
		accessKey:  key,
		enrollment: enrollment,
	}
}

func (c Client) UsageReports() (UsageReports, error) {
	req, err := http.NewRequest("GET", strings.Join([]string{c.URL, "rest", c.enrollment, "usage-reports"}, "/"), nil)
	if err != nil {
		return UsageReports{}, err
	}
	req.Header.Add("authorization", "bearer "+c.accessKey)
	req.Header.Add("api-version", "1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return UsageReports{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return UsageReports{}, fmt.Errorf("NOT OKAY: %s", resp.Status)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UsageReports{}, err
	}

	usageReports := &UsageReports{}
	err = json.Unmarshal(body, usageReports)
	if err != nil {
		return UsageReports{}, err
	}
	return *usageReports, nil
}

func (c Client) MonthlyUsageReport() (DetailedUsageReport, error) {
	req, err := http.NewRequest("GET", strings.Join([]string{c.URL, "rest", c.enrollment, "usage-report?type=detail"}, "/"), nil)
	if err != nil {
		return DetailedUsageReport{}, err
	}
	req.Header.Add("authorization", "bearer "+c.accessKey)
	req.Header.Add("api-version", "1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return DetailedUsageReport{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return DetailedUsageReport{}, fmt.Errorf("NOT OKAY: %s", resp.Status)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DetailedUsageReport{}, err
	}

	return DetailedUsageReport{CSV: body}, nil
}
