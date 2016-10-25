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

type DBClient interface {
	SaveReports(datamodels.Reports) error
}

type UsageDataJob struct {
	log      *logrus.Logger
	location *time.Location

	IAASClients []IaasClient

	saveFile bool
	DBClient DBClient
}

func NewJob(
	log *logrus.Logger,
	location *time.Location,
	iaasClients []IaasClient,
	dbClient DBClient,
	saveFile bool,
) *UsageDataJob {
	return &UsageDataJob{
		log:      log,
		location: location,

		IAASClients: iaasClients,
		DBClient:    dbClient,

		saveFile: saveFile,
	}
}

func (j *UsageDataJob) Run() {
	j.log.Debug("Entering usagedatajob.Run")
	defer j.log.Debug("Returning usagedatajob.Run")

	runTime := time.Now().In(j.location)
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

		j.log.Debugf("Saving %s data to database...", iaasClient.Name())
		err = j.DBClient.SaveReports(normalizedData)
		if err != nil {
			j.log.Errorf("Failed to save %s usage data to the database: %s", iaasClient.Name(), err.Error())
		} else {
			j.log.Debugf("Saved %s data to database", iaasClient.Name())
		}

		if j.saveFile { // Append to file
			j.log.Debugf("Writing %s data to file...", iaasClient.Name())
			if i == 0 {
				err = gocsv.MarshalFile(&normalizedData, normalizedFile)
			} else {
				err = gocsv.MarshalWithoutHeaders(&normalizedData, normalizedFile)
			}
			if err != nil {
				j.log.Errorf("Failed to write normalized %s data to file: %s", iaasClient.Name(), err.Error())
			} else {
				j.log.Debugf("Wrote normalized %s data to %s", iaasClient.Name(), normalizedFile.Name())
			}
		}
	}

	if !j.saveFile {
		err = os.Remove(normalizedFileName)
		if err != nil {
			j.log.Warn("Failed to remove file:", normalizedFile)
		}
	}

	finishedTime := time.Now().In(j.location)
	j.log.Infof("Finished periodic job at %s. It took %s.", finishedTime.String(), finishedTime.Sub(runTime).String())
}
