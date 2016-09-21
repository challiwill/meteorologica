package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/challiwill/meteorologica/aws"
	"github.com/challiwill/meteorologica/azure"
	"github.com/challiwill/meteorologica/gcp"
	"github.com/challiwill/meteorologica/usagedatagetter"
	"github.com/robfig/cron"
)

type Client interface{}

var azureFlag = flag.Bool("azure", false, "")
var gcpFlag = flag.Bool("gcp", false, "")
var awsFlag = flag.Bool("aws", false, "")
var nowFlag = flag.Bool("now", false, "")
var verboseFlag = flag.Bool("v", false, "")

func main() {
	flag.Parse()

	getAzure := *azureFlag
	getAWS := *awsFlag
	getGCP := *gcpFlag
	getAll := !getAzure && !getGCP && !getAWS

	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
	if *verboseFlag {
		log.Level = logrus.DebugLevel
	}

	sfTime, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		sfTime = time.Now().Location()
		log.Warn("Failed to load San Francisco time, using local time instead. Current local time is: ", time.Now().In(sfTime).String())
	} else {
		log.Info("Using San Francisco time. Current SF time is: ", time.Now().In(sfTime).String())
	}

	var iaasClients []usagedatagetter.IaasClient
	var bucketClient usagedatagetter.BucketClient

	// Azure Client
	if getAzure || getAll {
		log.Debug("Creating Azure Client")
		accessKey := os.Getenv("AZURE_ACCESS_KEY")
		enrollmentNumber := os.Getenv("AZURE_ENROLLMENT_NUMBER")
		if accessKey == "" || enrollmentNumber == "" {
			log.Error("Azure requires AZURE_ACCESS_KEY and AZURE_ENROLLMENT_NUMBER environment variables to be set")
		} else {
			azureClient := azure.NewClient(log, "https://ea.azure.com/", accessKey, enrollmentNumber)
			iaasClients = append(iaasClients, azureClient)
		}
	}

	// GCP Client
	log.Debug("Creating GCP Client")
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	bucketName := os.Getenv("GCP_BUCKET_NAME")
	if credentialsFile == "" || bucketName == "" {
		log.Fatal("GCP requires GCP_BUCKET_NAME and GOOGLE_APPLICATION_CREDENTIALS environment variables to be set")
	}
	gcpCredentials, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatal("Failed to create GCP credentials:", err.Error())
	} else {
		gcpClient, err := gcp.NewClient(log, gcpCredentials, bucketName)
		if err != nil {
			log.Fatal("Failed to create GCP client:", err.Error())
		}
		if getGCP || getAll {
			iaasClients = append(iaasClients, gcpClient)
		}
		bucketClient = gcpClient
	}

	// AWS Client
	if getAWS || getAll {
		log.Debug("Creating AWS Client")
		az := os.Getenv("AWS_REGION")
		bucketName := os.Getenv("AWS_BUCKET_NAME")
		accountNumber := os.Getenv("AWS_MASTER_ACCOUNT_NUMBER")
		if az == "" || bucketName == "" || accountNumber == "" {
			log.Error("AWS requires AWS_REGION, AWS_BUCKET_NAME, and AWS_MASTER_ACCOUNT_NUMBER environment variables to be set")
		} else {
			sess, err := session.NewSession(&awssdk.Config{
				Region: awssdk.String(az),
			})
			if err != nil {
				log.Error("Failed to create AWS credentials:", err.Error())
			} else {
				awsClient := aws.NewClient(log, az, bucketName, accountNumber, sess)
				iaasClients = append(iaasClients, awsClient)
			}
		}
	}

	usageDataJob := usagedatagetter.NewJob(iaasClients, bucketClient, log, sfTime)

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
		fmt.Fprintf(w, "Meteorologica is deployed\n\n Last job ran at %s\n Next job will run in roughly %s", usageDataJob.LastRunTime.String(), c.Entries()[0].Next.Sub(time.Now().In(sfTime)).String())
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
