package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/challiwill/meteorologica/azure"
	"github.com/challiwill/meteorologica/gcp"
)

type Client interface{}

func main() {
	azureClient := azure.NewClient("https://ea.azure.com/", os.Getenv("AZURE_ACCESS_KEY"), os.Getenv("AZURE_ENROLLMENT_NUMBER"))

	gcpCredentials, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gcpClient, err := gcp.NewClient(gcpCredentials)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	azureMonthlyusage, err := azureClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Got Monthly Azure Usage: %s\n", azureMonthlyusage.CSV)

	gcpMonthlyUsage, err := gcpClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Got Monthly GCP Usage:")
	for _, usage := range gcpMonthlyUsage.DailyUsage {
		fmt.Println(usage.CSV)
	}

	os.Exit(0)
}
