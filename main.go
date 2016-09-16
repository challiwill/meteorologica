package main

import (
	"fmt"
	"io/ioutil"
	"os"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/challiwill/meteorologica/aws"
	"github.com/challiwill/meteorologica/azure"
	"github.com/challiwill/meteorologica/gcp"
)

type Client interface{}

func main() {
	// AZURE CLIENT
	azureClient := azure.NewClient("https://ea.azure.com/", os.Getenv("AZURE_ACCESS_KEY"), os.Getenv("AZURE_ENROLLMENT_NUMBER"))

	// GCP CLIENT
	gcpCredentials, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		fmt.Println("Failed to create GCP credentials: ", err)
		os.Exit(1)
	}
	gcpClient, err := gcp.NewClient(gcpCredentials, os.Getenv("GCP_BUCKET_NAME"))
	if err != nil {
		fmt.Println("Failed to create GCP client: ", err)
		os.Exit(1)
	}

	// AWS CLIENT
	sess, err := session.NewSession(&awssdk.Config{
		Region: awssdk.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		fmt.Println("Failed to create AWS credentails: ", err)
		os.Exit(1)
	}
	awsClient := aws.NewClient(os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_MASTER_ACCOUNT_NUMBER"), sess)

	// GET USAGES
	azureMonthlyusage, err := azureClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println("Failed to get Azure monthly usage: ", err)
	} else {
		fmt.Printf("Got Monthly Azure Usage: %s\n", string(azureMonthlyusage.CSV))
	}

	gcpMonthlyUsage, err := gcpClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println("Failed to get GCP monthly usage: ", err)
	} else {
		fmt.Println("Got Monthly GCP Usage:")
		for _, usage := range gcpMonthlyUsage.DailyUsage {
			fmt.Println(string(usage.CSV))
		}
	}

	awsMonthlyusage, err := awsClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println("Failed to get AWS monthly usage: ", err)
	} else {
		fmt.Printf("Got Monthly AWS Usage: %s\n", string(awsMonthlyusage.CSV))
	}

	os.Exit(0)
}
