package azure

import (
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
		return nil, fmt.Errorf("NOT OKAY: %s", resp.Status)
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
	return DetailedUsageReport{CSV: []byte(csvFirstTwoLinesRemoved)}
	return DetailedUsageReport{CSV: body}
}
