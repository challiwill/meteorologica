package gcp

import (
	"fmt"
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
	CSV []byte
}

type Client struct {
	StorageService StorageService
	BucketName     string
}

func NewClient(jsonCredentials []byte, bucketName string) (*Client, error) {
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
		BucketName:     bucketName,
	}, nil
}

func (c Client) MonthlyUsageReport() (DetailedUsageReport, error) {
	monthlyUsageReport := DetailedUsageReport{}

	for i := 1; i < time.Now().Day(); i++ {
		dailyUsage, err := c.DailyUsageReport(i)
		if err != nil {
			fmt.Printf("Failed to get GCP Daily Usage for day %d this month: %s\n", i, err.Error())
			continue
		}
		monthlyUsageReport.DailyUsage = append(monthlyUsageReport.DailyUsage, dailyUsage)
	}
	return monthlyUsageReport, nil

}

func (c Client) DailyUsageReport(day int) (DailyUsageReport, error) {
	resp, err := c.StorageService.DailyUsage(c.BucketName, dailyBillingFileName(day))
	if err != nil {
		return DailyUsageReport{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DailyUsageReport{}, err
	}
	return DailyUsageReport{CSV: body}, nil
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
