package usagedatajob

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
	log      *logrus.Logger
	location *time.Location

	IAASClients []IaasClient

	saveFile     bool
	saveToBucket bool
	saveToDB     bool
	BucketClient BucketClient

	LastRunTime time.Time
}

func NewJob(
	log *logrus.Logger,
	location *time.Location,
	iaasClients []IaasClient,
	bucketClient BucketClient,
	saveFile bool,
	saveToBucket bool,
	saveToDB bool,
) *UsageDataJob {
	return &UsageDataJob{
		log:      log,
		location: location,

		IAASClients:  iaasClients,
		BucketClient: bucketClient,

		saveFile:     saveFile,
		saveToBucket: saveToBucket,
		saveToDB:     saveToDB,
	}
}

func (j *UsageDataJob) Run() {
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

	for i, iaasClient := range j.IAASClients {
		normalizedData, err := iaasClient.GetNormalizedUsage()
		if err != nil {
			j.log.Errorf("Failed to get %s usage data: %s", iaasClient.Name(), err.Error())
			continue
		}

		if j.saveFile || j.saveToBucket { // Append to file
			if i == 0 {
				err = gocsv.Marshal(&normalizedData, normalizedFile)
			} else {
				err = gocsv.MarshalWithoutHeaders(&normalizedData, normalizedFile)
			}
			if err != nil {
				j.log.Errorf("Failed to write normalized %s data to file: %s", iaasClient.Name(), err.Error())
				continue
			}
			j.log.Infof("Wrote normalized %s data to %s", iaasClient.Name(), normalizedFile.Name())
		}
	}

	if j.saveToBucket { // Send file to bucket
		err = j.BucketClient.PublishFileToBucket(normalizedFileName)
		if err != nil {
			j.log.Error("Failed to publish data to storage bucket:", err)
		} else {
			if !j.saveFile {
				err = os.Remove(normalizedFileName) // only remove if succeeded to parse
				if err != nil {
					j.log.Warn("Failed to remove file:", normalizedFile)
				}
			}
		}
	}

	finishedTime := time.Now().In(j.location)
	j.log.Infof("Finished periodic job at %s. It took %s.", finishedTime.String(), finishedTime.Sub(runTime).String())
}
