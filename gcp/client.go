package gcp

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
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

type DailyUsageReport struct {
	CSV []byte
}

type Client struct {
	StorageService StorageService
	BucketName     string
	log            *logrus.Logger
}

func NewClient(log *logrus.Logger, jsonCredentials []byte, bucketName string) (*Client, error) {
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
		log:            log,
	}, nil
}

func (c Client) Name() string {
	return "GCP"
}

func (c Client) GetNormalizedUsage() (datamodels.Reports, error) {
	c.log.Info("Getting monthly GCP usage...")
	gcpMonthlyUsage, err := c.MonthlyUsageReport()
	if err != nil {
		c.log.Error("Failed to get GCP monthly usage")
		return datamodels.Reports{}, err
	}

	reports := datamodels.Reports{}
	c.log.Debug("Got monthly GCP usage")
	for i, usage := range gcpMonthlyUsage.DailyUsage {
		usageReader, err := NewUsageReader(c.log, usage.CSV)
		if err != nil {
			c.log.Error("Failed to parse GCP usage for day", i+1)
			continue
		}
		reports = append(reports, usageReader.Normalize()...)
	}

	if len(reports) == 0 {
		return datamodels.Reports{}, errors.New("Failed to parse all GCP usage data")
	}
	return reports, nil
}

func (c Client) MonthlyUsageReport() (DetailedUsageReport, error) {
	monthlyUsageReport := DetailedUsageReport{}

	for i := 1; i < time.Now().Day(); i++ {
		dailyUsage, err := c.DailyUsageReport(i)
		if err != nil {
			c.log.Warnf("Failed to get GCP Daily Usage for %s, %d: %s", time.Now().Month().String(), i, err.Error())
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

func (c Client) PublishFileToBucket(name string) error {
	object := &storage.Object{
		Name:        "a_code_name_saam/" + name,
		ContentType: "text/csv",
	}
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		c.log.Errorf("Failed to open normalized file: %s", name)
		return err
	}

	res, err := c.StorageService.Insert(c.BucketName, object, file)
	if err != nil {
		c.log.Errorf("Objects.Insert to bucket '%s' failed", c.BucketName)
		return err
	}
	c.log.Infof("Created object %v at location %v", res.Name, res.SelfLink)

	return nil
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
