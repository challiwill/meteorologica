package main

import (
	"fmt"
	"os"

	"github.com/challiwill/meteorologica/azure"
)

func main() {
	azureClient := azure.NewClient("https://ea.azure.com/", os.Getenv("AZURE_ACCESS_KEY"), os.Getenv("AZURE_ENROLLMENT_NUMBER"))
	mu, err := azureClient.MonthlyUsageReport()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Got Monthly Azure Usage: %s\n", mu.CSV)
	os.Exit(0)
}
