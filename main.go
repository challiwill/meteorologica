package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/challiwill/meteorologica/aws"
	"github.com/challiwill/meteorologica/azure"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/challiwill/meteorologica/gcp"
	"github.com/gocarina/gocsv"
)

type Client interface{}

var azureFlag = flag.Bool("azure", false, "")
var gcpFlag = flag.Bool("gcp", false, "")
var awsFlag = flag.Bool("aws", false, "")

func main() {
	flag.Parse()
	getAzure := *azureFlag
	getGCP := *gcpFlag
	getAWS := *awsFlag
	getAll := !getAzure && !getGCP && !getAWS

	normalizedFile, err := os.Create("normalized_iaas_billing_data.csv")
	if err != nil {
		fmt.Println("Failed to create normalized file: ", err.Error())
		os.Exit(1)
	}

	// AZURE CLIENT
	if getAzure || getAll {
		normalizedAzure, err := getAzureUsage()
		if err != nil {
			fmt.Println("Failed to get Azure usage data: ", err.Error())
		} else {
			err = gocsv.MarshalFile(&normalizedAzure, normalizedFile)
			if err != nil {
				fmt.Println("Failed to write normalized Azure data to file: ", err.Error())
			} else {
				fmt.Println("Wrote normalized azure data to ", normalizedFile.Name())
			}
		}
	}

	// GCP CLIENT
	if getGCP || getAll {
		err := getGCPUsage()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// AWS CLIENT
	if getAWS || getAll {
		err := getAWSUsage()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	os.Exit(0)
}

func getAzureUsage() (datamodels.Reports, error) {
	azureClient := azure.NewClient("https://ea.azure.com/", os.Getenv("AZURE_ACCESS_KEY"), os.Getenv("AZURE_ENROLLMENT_NUMBER"))

	fmt.Println("Getting Monthly Azure Usage...")
	azureMonthlyusage, err := azureClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println("Failed to get Azure monthly usage: ", err)
		return datamodels.Reports{}, err
	}

	fmt.Println("Got Monthly Azure Usage")
	err = ioutil.WriteFile("azure.csv", azureMonthlyusage.CSV, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to save Azure Usage to file")
		return datamodels.Reports{}, err
	}
	fmt.Println("Saved Azure Usage to azure.csv")

	azureDataFile, err := os.OpenFile("azure.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to open Azure file")
		return datamodels.Reports{}, err
	}
	defer azureDataFile.Close()
	usageReader, err := azure.NewUsageReader(azureDataFile)
	if err != nil {
		fmt.Println("Failed to parse Azure file")
		return datamodels.Reports{}, err
	}
	defer os.Remove("azure.csv") // only remove if succeeded to parse
	return usageReader.Normalize(), nil
}

func getGCPUsage() error {
	gcpCredentials, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		fmt.Println("Failed to create GCP credentials: ", err)
		return err
	}
	gcpClient, err := gcp.NewClient(gcpCredentials, os.Getenv("GCP_BUCKET_NAME"))
	if err != nil {
		fmt.Println("Failed to create GCP client: ", err)
		return err
	}

	fmt.Println("Getting Monthly GCP Usage...")
	gcpMonthlyUsage, err := gcpClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println("Failed to get GCP monthly usage: ", err)
		return err
	}

	fmt.Println("Got Monthly GCP Usage:")
	for i, usage := range gcpMonthlyUsage.DailyUsage {
		err = ioutil.WriteFile("gcp-"+strconv.Itoa(i+1)+".csv", usage.CSV, os.ModePerm)
		fmt.Println("Saved GCP Usages to gcp-" + strconv.Itoa(i+1) + ".csv")
		if err != nil {
			fmt.Println("Failed to save GCP Usage to file")
			return err
		}
	}

	return nil
}

func getAWSUsage() error {
	sess, err := session.NewSession(&awssdk.Config{
		Region: awssdk.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		fmt.Println("Failed to create AWS credentails: ", err)
		return err
	}

	awsClient := aws.NewClient(os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_MASTER_ACCOUNT_NUMBER"), sess)

	fmt.Println("Getting Monthly AWS Usage...")
	awsMonthlyUsage, err := awsClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println("Failed to get AWS monthly usage: ", err)
		return err
	}

	fmt.Println("Got Monthly AWS Usage")
	err = ioutil.WriteFile("aws.csv", awsMonthlyUsage.CSV, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to save AWS Usage to file")
		return err
	}
	fmt.Println("AWS Usage saved to aws.csv")

	return nil
}
