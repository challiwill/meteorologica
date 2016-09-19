package main

import (
	"errors"
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
				fmt.Println("Wrote normalized Azure data to ", normalizedFile.Name())
			}
		}
	}

	// GCP CLIENT
	if getGCP || getAll {
		normalizedGCP, err := getGCPUsage()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			err = gocsv.MarshalFile(&normalizedGCP, normalizedFile)
			if err != nil {
				fmt.Println("Failed to write normalized GCP data to file: ", err.Error())
			} else {
				fmt.Println("Wrote normalized GCP data to ", normalizedFile.Name())
			}
		}
	}

	// AWS CLIENT
	if getAWS || getAll {
		normalizedAWS, err := getAWSUsage()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			err = gocsv.MarshalFile(&normalizedAWS, normalizedFile)
			if err != nil {
				fmt.Println("Failed to write normalized AWS data to file: ", err.Error())
			} else {
				fmt.Println("Wrote normalized AWS data to ", normalizedFile.Name())
			}
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

func getGCPUsage() (datamodels.Reports, error) {
	gcpCredentials, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		fmt.Println("Failed to create GCP credentials: ", err)
		return datamodels.Reports{}, err
	}
	gcpClient, err := gcp.NewClient(gcpCredentials, os.Getenv("GCP_BUCKET_NAME"))
	if err != nil {
		fmt.Println("Failed to create GCP client: ", err)
		return datamodels.Reports{}, err
	}

	fmt.Println("Getting Monthly GCP Usage...")
	gcpMonthlyUsage, err := gcpClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println("Failed to get GCP monthly usage: ", err)
		return datamodels.Reports{}, err
	}

	reports := datamodels.Reports{}
	fmt.Println("Got Monthly GCP Usage:")
	for i, usage := range gcpMonthlyUsage.DailyUsage {
		fileName := "gcp-" + strconv.Itoa(i+1) + ".csv"
		err = ioutil.WriteFile(fileName, usage.CSV, os.ModePerm)
		fmt.Println("Saved GCP Usages to gcp-" + strconv.Itoa(i+1) + ".csv")
		if err != nil {
			fmt.Println("Failed to save GCP Usage to file")
			continue
		}

		gcpDataFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to open GCP file ", fileName)
			continue
		}
		usageReader, err := gcp.NewUsageReader(gcpDataFile)
		if err != nil {
			fmt.Println("Failed to parse GCP file ", fileName)
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

func getAWSUsage() (datamodels.Reports, error) {
	az := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&awssdk.Config{
		Region: awssdk.String(az),
	})
	if err != nil {
		fmt.Println("Failed to create AWS credentails: ", err)
		return datamodels.Reports{}, err
	}

	awsClient := aws.NewClient(os.Getenv("AWS_BUCKET_NAME"), os.Getenv("AWS_MASTER_ACCOUNT_NUMBER"), sess)

	fmt.Println("Getting Monthly AWS Usage...")
	awsMonthlyUsage, err := awsClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println("Failed to get AWS monthly usage: ", err)
		return datamodels.Reports{}, err
	}

	fmt.Println("Got Monthly AWS Usage")
	err = ioutil.WriteFile("aws.csv", awsMonthlyUsage.CSV, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to save AWS Usage to file")
		return datamodels.Reports{}, err
	}
	fmt.Println("AWS Usage saved to aws.csv")

	awsDataFile, err := os.OpenFile("aws.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to open AWS file")
		return datamodels.Reports{}, err
	}
	defer awsDataFile.Close()
	usageReader, err := aws.NewUsageReader(awsDataFile, az)
	if err != nil {
		fmt.Println("Failed to parse AWS file")
		return datamodels.Reports{}, err
	}
	defer os.Remove("aws.csv") // only remove if succeeded to parse
	return usageReader.Normalize(), nil
}
