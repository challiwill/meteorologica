package gcp

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/csv"
	"github.com/challiwill/meteorologica/datamodels"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

//go:generate counterfeiter . StorageService

type StorageService interface {
	DailyUsage(string, string) (*http.Response, error)
	Insert(string, *storage.Object, *os.File) (*storage.Object, error)
}

type DetailedUsageReport struct {
	DailyUsage []DailyUsageReport
}

type DailyUsageReport []byte

type Client struct {
	StorageService StorageService
	BucketName     string
	Log            *logrus.Logger
	Location       *time.Location
}

func NewClient(log *logrus.Logger, location *time.Location, jsonCredentials []byte, bucketName string) (*Client, error) {
	jwtConfig, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/devstorage.read_write")
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
		Log:            log,
		Location:       location,
	}, nil
}

func (c Client) Name() string {
	return "GCP"
}

func (c Client) GetNormalizedUsage() (datamodels.Reports, error) {
	c.Log.Info("Getting monthly GCP usage...")
	gcpMonthlyUsage, err := c.GetBillingData()
	if err != nil {
		c.Log.Error("Failed to get GCP monthly usage")
		return datamodels.Reports{}, err
	}

	reports := datamodels.Reports{}
	c.Log.Debug("Got monthly GCP usage")
	usageReader := NewNormalizer(c.Log, c.Location)
	for i, usage := range gcpMonthlyUsage.DailyUsage {
		readerCleaner, err := csv.NewReaderCleaner(bytes.NewReader(usage), 15)
		if err != nil {
			return nil, err
		}
		dailyReport := []*Usage{}
		err = csv.GenerateReports(readerCleaner, &dailyReport)
		if err != nil {
			// try again with 18 columns as it is sometimes
			readerCleaner, err := csv.NewReaderCleaner(bytes.NewReader(usage), 18)
			if err != nil {
				return nil, err
			}
			dailyReport := []*Usage{}
			err = csv.GenerateReports(readerCleaner, &dailyReport)
			if err != nil {
				c.Log.Errorf("Failed to parse GCP usage for day: %d %s: %s", i+1, time.Now().In(c.Location).Month().String(), err.Error())
				continue
			}
		}
		reports = append(reports, usageReader.Normalize(dailyReport)...)
	}

	if len(reports) == 0 {
		return datamodels.Reports{}, errors.New("Failed to parse all GCP usage data")
	}
	return reports, nil
}

func (c Client) GetBillingData() (DetailedUsageReport, error) {
	monthlyUsageReport := DetailedUsageReport{}

	for i := 1; i < time.Now().In(c.Location).Day(); i++ {
		dailyUsage, err := c.DailyUsageReport(i)
		if err != nil {
			c.Log.Warnf("Failed to get GCP Daily Usage for %s, %d: %s", time.Now().In(c.Location).Month().String(), i, err.Error())
			continue
		}
		monthlyUsageReport.DailyUsage = append(monthlyUsageReport.DailyUsage, dailyUsage)
	}
	return monthlyUsageReport, nil

}

func (c Client) DailyUsageReport(day int) (DailyUsageReport, error) {
	resp, err := c.StorageService.DailyUsage(c.BucketName, c.dailyBillingFileName(day))
	if err != nil {
		return nil, fmt.Errorf("Making request to GCP failed: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GCP responded with error: %s", resp.Status)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c Client) PublishFileToBucket(name string) error {
	object := &storage.Object{
		Name:        name,
		ContentType: "text/csv",
	}
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		c.Log.Errorf("Failed to open normalized file: %s", name)
		return err
	}

	res, err := c.StorageService.Insert(c.BucketName, object, file)
	if err != nil {
		c.Log.Errorf("Objects.Insert to bucket '%s' failed", c.BucketName)
		return err
	}
	c.Log.Infof("Created object %v at location %v", res.Name, res.SelfLink)

	return nil
}

func (c Client) dailyBillingFileName(day int) string {
	year, month, _ := time.Now().In(c.Location).Date()
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
