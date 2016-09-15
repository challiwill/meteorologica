package gcp

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

//go:generate counterfeiter . StorageService

type StorageService interface {
	DailyUsage(string, string) (*http.Response, error)
}

type DetailedUsageReport struct {
	DailyUsage []DailyUsageReport
}

type DailyUsageReport struct {
	CSV string
}

type Client struct {
	StorageService StorageService
}

func NewClient(jsonCredentials []byte) (*Client, error) {
	jwtConfig, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/devstorage.read_only")
	if err != nil {
		return nil, err
	}
	service, err := storage.New(jwtConfig.Client(oauth2.NoContext))
	if err != nil {
		return nil, err
	}
	return &Client{
		StorageService: &storageService{service: service},
	}, nil
}

func (c Client) MonthlyUsageReport() (DetailedUsageReport, error) {
	monthlyUsageReport := DetailedUsageReport{}

	for i := 1; i < time.Now().Day(); i++ {
		resp, err := c.StorageService.DailyUsage("pivotal_billing", dailyBillingFileName(i))
		if err != nil {
			return DetailedUsageReport{}, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return DetailedUsageReport{}, err
		}

		monthlyUsageReport.DailyUsage = append(monthlyUsageReport.DailyUsage, DailyUsageReport{CSV: string(body)})
	}
	return monthlyUsageReport, nil

}

func dailyBillingFileName(day int) string {
	year, month, _ := time.Now().Date()
	monthStr := padMonth(month)
	dayStr := padDay(day)
	return url.QueryEscape(strings.Join([]string{"Billing", strconv.Itoa(year), monthStr, dayStr}, "-") + ".csv")
}

func padMonth(month time.Month) string {
	m := strconv.Itoa(int(month))
	if month < 10 {
		return "0" + m
	}
	return m
}

func padDay(day int) string {
	d := strconv.Itoa(day)
	if day < 10 {
		return "0" + d
	}
	return d
}
