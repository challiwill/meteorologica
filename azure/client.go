package azure

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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

type Client struct {
	URL        string
	client     *http.Client
	accessKey  string
	enrollment int
}

func NewClient(serverURL, key string, enrollment int) *Client {
	return &Client{
		URL:        serverURL,
		client:     new(http.Client),
		accessKey:  key,
		enrollment: enrollment,
	}
}

func (c Client) UsageReports() (UsageReports, error) {
	req, err := http.NewRequest("GET", strings.Join([]string{c.URL, "rest", strconv.Itoa(c.enrollment), "usage-reports"}, "/"), nil)
	if err != nil {
		return UsageReports{}, err
	}
	req.Header.Add("authorization", c.accessKey)
	req.Header.Add("api-version", "1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return UsageReports{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return UsageReports{}, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UsageReports{}, err
	}

	usageReporst := &UsageReports{}
	err = json.Unmarshal(body, usageReporst)
	if err != nil {
		return UsageReports{}, err
	}
	return *usageReporst, nil
}
