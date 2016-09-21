package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/challiwill/meteorologica/aws"
	"github.com/challiwill/meteorologica/azure"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/challiwill/meteorologica/gcp"
	"github.com/gocarina/gocsv"
	"github.com/robfig/cron"
)

type Client interface{}

var azureFlag = flag.Bool("azure", false, "")
var gcpFlag = flag.Bool("gcp", false, "")
var awsFlag = flag.Bool("aws", false, "")
var nowFlag = flag.Bool("now", false, "")

func main() {
	flag.Parse()

	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel

	sfTime, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Error("Failed to load San Francisco time, using local time instead")
		sfTime = time.Now().Location()
	} else {
		log.Info("Using San Francisco time. Current SF time is: ", time.Now().In(sfTime).String())
	}
	usageDataJob := UsageDataJob{
		getAzure: *azureFlag,
		getGCP:   *gcpFlag,
		getAWS:   *awsFlag,
		log:      log,
		location: sfTime,
	}

	runTime := time.Time{}

	// BILLING DATA
	if *nowFlag {
		usageDataJob.Run()
		os.Exit(0)
	}

	c := cron.NewWithLocation(sfTime)
	err = c.AddJob("@midnight", usageDataJob)
	if err != nil {
		log.Fatal("Could not create cron job:", err.Error())
	}
	c.Start()

	// HEALTHCHECK
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Meteorologica is deployed\n\n Last job ran at %s\n Next job will run in roughly %s", runTime.String(), c.Entries()[0].Next.Sub(time.Now().In(sfTime)).String())
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
	os.Exit(0)
}

type UsageDataJob struct {
	log      *logrus.Logger
	getAzure bool
	getGCP   bool
	getAWS   bool
	location *time.Location
}

func (j UsageDataJob) Run() {
	getAll := !j.getAzure && !j.getGCP && !j.getAWS
	runTime := time.Now().In(j.location)
	j.log.Infof("Running periodic job at %s ...", runTime.String())
	normalizedFileName := strings.Join([]string{
		strconv.Itoa(runTime.Year()),
		runTime.Month().String(),
		strconv.Itoa(runTime.Day()),
		"normalized-billing-data.csv",
	}, "-")
	normalizedFile, err := os.Create(normalizedFileName)
	if err != nil {
		j.log.Fatal("Failed to create normalized file: ", err.Error())
	}

	// AZURE CLIENT
	if j.getAzure || getAll {
		normalizedAzure, err := getAzureUsage(j.log)
		if err != nil {
			j.log.Error("Failed to get Azure usage data: ", err.Error())
		} else {
			err = gocsv.MarshalFile(&normalizedAzure, normalizedFile)
			if err != nil {
				j.log.Error("Failed to write normalized Azure data to file: ", err.Error())
			} else {
				j.log.Info("Wrote normalized Azure data to ", normalizedFile.Name())
			}
		}
	}

	// GCP CLIENT
	if j.getGCP || getAll {
		normalizedGCP, err := getGCPUsage(j.log)
		if err != nil {
			j.log.Error("Failed to get GCP usage data: ", err.Error())
		} else {
			err = gocsv.MarshalFile(&normalizedGCP, normalizedFile)
			if err != nil {
				j.log.Error("Failed to write normalized GCP data to file: ", err.Error())
			} else {
				j.log.Info("Wrote normalized GCP data to ", normalizedFile.Name())
			}
		}
	}

	// AWS CLIENT
	if j.getAWS || getAll {
		normalizedAWS, err := getAWSUsage(j.log)
		if err != nil {
			j.log.Error("Failed to get AWS usage data: ", err.Error())
		} else {
			err = gocsv.MarshalFile(&normalizedAWS, normalizedFile)
			if err != nil {
				j.log.Error("Failed to write normalized AWS data to file: ", err.Error())
			} else {
				j.log.Info("Wrote normalized AWS data to ", normalizedFile.Name())
			}
		}
	}

	gcpCredentials, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		j.log.Error("Failed to create GCP credentials to publish normalized data file")
	}
	gcpClient, err := gcp.NewClient(gcpCredentials, os.Getenv("GCP_BUCKET_NAME"))
	if err != nil {
		j.log.Error("Failed to create GCP client to publish normalized data file:", err)
	}
	err = gcpClient.PublishFileToBucket(j.log, normalizedFileName)
	if err != nil {
		j.log.Error("Failed to publish data to GCP Bucket:", err)
	}
	finishedTime := time.Now().In(j.location)
	j.log.Infof("Finished periodic job at %s. It took %s.", finishedTime.String(), finishedTime.Sub(runTime).String())
}

func getAzureUsage(log *logrus.Logger) (datamodels.Reports, error) {
	azureClient := azure.NewClient("https://ea.azure.com/", os.Getenv("AZURE_ACCESS_KEY"), os.Getenv("AZURE_ENROLLMENT_NUMBER"))

	log.Info("Getting Monthly Azure Usage...")
	azureMonthlyusage, err := azureClient.MonthlyUsageReport()
	if err != nil {
		log.Error("Failed to get Azure monthly usage")
		return datamodels.Reports{}, err
	}

	log.Debug("Got Monthly Azure Usage")
	err = ioutil.WriteFile("azure.csv", azureMonthlyusage.CSV, os.ModePerm)
	if err != nil {
		log.Error("Failed to save Azure Usage to file")
		return datamodels.Reports{}, err
	}
	log.Debug("Saved Azure Usage to azure.csv")

	azureDataFile, err := os.OpenFile("azure.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error("Failed to open Azure file")
		return datamodels.Reports{}, err
	}
	defer azureDataFile.Close()
	usageReader, err := azure.NewUsageReader(azureDataFile)
	if err != nil {
		log.Error("Failed to parse Azure file")
		return datamodels.Reports{}, err
	}
	defer os.Remove("azure.csv") // only remove if succeeded to parse
	return usageReader.Normalize(), nil
}

func getGCPUsage(log *logrus.Logger) (datamodels.Reports, error) {
	gcpCredentials, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		log.Error("Failed to create GCP credentials")
		return datamodels.Reports{}, err
	}
	gcpClient, err := gcp.NewClient(gcpCredentials, os.Getenv("GCP_BUCKET_NAME"))
	if err != nil {
		log.Error("Failed to create GCP client")
		return datamodels.Reports{}, err
	}

	log.Info("Getting Monthly GCP Usage...")
	gcpMonthlyUsage, err := gcpClient.MonthlyUsageReport()
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
		usageReader, err := gcp.NewUsageReader(gcpDataFile)
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

func getAWSUsage(log *logrus.Logger) (datamodels.Reports, error) {
	az := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&awssdk.Config{
		Region: awssdk.String(az),
	})
	if err != nil {
		log.Error("Failed to create AWS credentails")
		return datamodels.Reports{}, err
	}

	awsClient := aws.NewClient(os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_MASTER_ACCOUNT_NUMBER"), sess)

	log.Info("Getting Monthly AWS Usage...")
	awsMonthlyUsage, err := awsClient.MonthlyUsageReport()
	if err != nil {
		log.Error("Failed to get AWS monthly usage: ", err)
		return datamodels.Reports{}, err
	}

	log.Debug("Got Monthly AWS Usage")
	err = ioutil.WriteFile("aws.csv", awsMonthlyUsage.CSV, os.ModePerm)
	if err != nil {
		log.Error("Failed to save AWS Usage to file")
		return datamodels.Reports{}, err
	}
	log.Debug("AWS Usage saved to aws.csv")

	awsDataFile, err := os.OpenFile("aws.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error("Failed to open AWS file")
		return datamodels.Reports{}, err
	}
	defer awsDataFile.Close()
	usageReader, err := aws.NewUsageReader(awsDataFile, az)
	if err != nil {
		log.Error("Failed to parse AWS file")
		return datamodels.Reports{}, err
	}
	defer os.Remove("aws.csv") // only remove if succeeded to parse
	return usageReader.Normalize(), nil
}
