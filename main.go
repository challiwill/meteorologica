package main

import (
	"fmt"
	"os"

	"github.com/challiwill/meteorologica/azure"
)

func main() {
	azureClient := azure.NewClient("https://ea.azure.com/", os.Getenv("AZURE_ACCESS_KEY"), os.Getenv("AZURE_ENROLLMENT_NUMBER"))
	ur, err := azureClient.UsageReports()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Success: %#v\n", ur)
	os.Exit(0)
}
