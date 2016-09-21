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
}

func NewClient(jsonCredentials []byte, bucketName string) (*Client, error) {
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
	}, nil
}

func (c Client) Name() string {
	return "GCP"
}

func (c Client) GetNormalizedUsage(log *logrus.Logger) (datamodels.Reports, error) {
	log.Info("Getting Monthly GCP Usage...")
	gcpMonthlyUsage, err := c.MonthlyUsageReport(log)
	if err != nil {
		log.Error("Failed to get GCP monthly usage")
		return datamodels.Reports{}, err
	}

	reports := datamodels.Reports{}
	log.Debug("Got Monthly GCP Usage:")
	for i, usage := range gcpMonthlyUsage.DailyUsage {
		fileName := "gcp-" + strconv.Itoa(i+1) + ".csv"
		err = ioutil.WriteFile(fileName, usage.CSV, os.ModePerm)
		log.Debug("Saved GCP Usages to gcp-" + strconv.Itoa(i+1) + ".csv")
		if err != nil {
			log.Error("Failed to save GCP Usage to file")
			continue
		}

		gcpDataFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Error("Failed to open GCP file ", fileName)
			continue
		}
		usageReader, err := NewUsageReader(gcpDataFile)
		if err != nil {
			log.Error("Failed to parse GCP file ", fileName)
			continue
		}
		reports = append(reports, usageReader.Normalize()...)
		gcpDataFile.Close()
		defer os.Remove(fileName) // only remove if succeeded to parse
	}

	if len(reports) == 0 {
		return datamodels.Reports{}, errors.New("Failed to parse all GCP usage data")
	}
	return reports, nil
}

func (c Client) MonthlyUsageReport(log *logrus.Logger) (DetailedUsageReport, error) {
	monthlyUsageReport := DetailedUsageReport{}

	for i := 1; i < time.Now().Day(); i++ {
		dailyUsage, err := c.DailyUsageReport(i)
		if err != nil {
			log.Warnf("Failed to get GCP Daily Usage for %d, %s: %s", i, time.Now().Month().String(), err.Error())
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

func (c Client) PublishFileToBucket(log *logrus.Logger, name string) error {
	object := &storage.Object{
		Name:        "a_code_name_saam/" + name,
		ContentType: "text/csv",
	}
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		log.Errorf("Failed to open normalized file: %s", name)
		return err
	}

	res, err := c.StorageService.Insert(c.BucketName, object, file)
	if err != nil {
		log.Errorf("Objects.Insert to bucket '%s' failed", c.BucketName)
		return err
	}
	log.Infof("Created object %v at location %v", res.Name, res.SelfLink)

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
