package usagedatagetter

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/gocarina/gocsv"
)

type IaasClient interface {
	Name() string
	GetNormalizedUsage() (datamodels.Reports, error)
}

type BucketClient interface {
	PublishFileToBucket(string) error
}

type UsageDataJob struct {
	log          *logrus.Logger
	IAASClients  []IaasClient
	BucketClient BucketClient
	location     *time.Location
	LastRunTime  time.Time
}

func NewJob(
	iaasClients []IaasClient,
	bucketClient BucketClient,
	log *logrus.Logger,
	location *time.Location,
) UsageDataJob {
	return UsageDataJob{
		log:          log,
		IAASClients:  iaasClients,
		BucketClient: bucketClient,
		location:     location,
	}
}

func (j UsageDataJob) Run() {
	runTime := time.Now().In(j.location)
	j.LastRunTime = runTime
	j.log.Infof("Running periodic job at %s ...", runTime.String())

	normalizedFileName := strings.Join([]string{
		strconv.Itoa(runTime.Year()),
		runTime.Month().String(),
		"normalized-billing-data.csv",
	}, "-")
	normalizedFile, err := os.Create(normalizedFileName)
	if err != nil {
		j.log.Fatal("Failed to create normalized file: ", err.Error())
	}

	for _, iaasClient := range j.IAASClients {
		normalizedData, err := iaasClient.GetNormalizedUsage()
		if err != nil {
			j.log.Errorf("Failed to get %s usage data: %s", iaasClient.Name(), err.Error())
			continue
		}
		err = gocsv.MarshalFile(&normalizedData, normalizedFile)
		if err != nil {
			j.log.Errorf("Failed to write normalized %s data to file: %s", iaasClient.Name(), err.Error())
			continue
		}
		j.log.Infof("Wrote normalized %s data to %s", iaasClient.Name(), normalizedFile.Name())
	}

	err = j.BucketClient.PublishFileToBucket(normalizedFileName)
	if err != nil {
		j.log.Error("Failed to publish data to storage bucket:", err)
	} else {
		err = os.Remove(normalizedFileName) // only remove if succeeded to parse
		if err != nil {
			j.log.Warn("Failed to remove file:", normalizedFile)
		}
	}

	finishedTime := time.Now().In(j.location)
	j.log.Infof("Finished periodic job at %s. It took %s.", finishedTime.String(), finishedTime.Sub(runTime).String())
}
