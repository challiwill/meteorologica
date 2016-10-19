package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/challiwill/meteorologica/aws"
	"github.com/challiwill/meteorologica/azure"
	"github.com/challiwill/meteorologica/db"
	"github.com/challiwill/meteorologica/db/migrations"
	"github.com/challiwill/meteorologica/gcp"
	"github.com/challiwill/meteorologica/usagedatajob"
	"github.com/heroku/rollrus"
	"github.com/jinzhu/configor"
	"github.com/robfig/cron"
)

type Client interface{}

var Config = struct {
	Port int `default:"8080"`

	Azure struct {
		AccessKey        string `yaml:"access-key" env:"M_AZURE_ACCESS_KEY"`
		EnrollmentNumber int    `yaml:"enrollment-number" env:"M_AZURE_ENROLLMENT_NUMBER"`
	}

	GCP struct {
		BucketName                 string `yaml:"bucket-name" env:"M_GCP_BUCKET_NAME"`
		ApplicationCredentialsPath string `yaml:"application-credentials-path" env:"M_GCP_APPLICATION_CREDENTIALS_PATH"`
	}

	AWS struct {
		Region              string
		MasterAccountNumber int64  `yaml:"master-account-number" env:"M_AWS_MASTER_ACCOUNT_NUMBER"`
		BucketName          string `yaml:"bucket-name" env:"M_AWS_BUCKET_NAME"`
		AccessKeyID         string `yaml:"access-key-id" env:"M_AWS_ACCESS_KEY_ID"`
		SecretAccessKey     string `yaml:"secret-access-key" env:"M_AWS_SECRET_ACCESS_KEY"`
	}

	DB struct {
		Username string
		Password string
		Address  string
		Name     string
	}

	Rollbar struct {
		Token string
	}
}{}

var (
	resourcesFlag = flag.String("resources", "aws,gcp,azure", "A comma seperated list of resource to retrieve billing information from. If none are specified the default is AWS, GCP, and Azure")
	nowFlag       = flag.Bool("now", false, "Run job now (instead of waiting for cron job to kick off at midnight)")
	verboseFlag   = flag.Bool("v", false, "Log at Debug level")
	fileFlag      = flag.Bool("file", false, "Save a local copy of the data as a normalized CSV file")
	dbFlag        = flag.Bool("db", true, "Save the data to the database")
)

func main() {
	_ = os.Setenv("CONFIGOR_ENV_PREFIX", "M")
	err := configor.Load(&Config, "configuration/meteorologica.yml")
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %s", err.Error())
	}
	flag.Parse()
	keepFile, saveToDB, resources := parseFlags()
	log := configureLog()

	sfTime, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		sfTime = time.Now().Location()
		log.Warn("Failed to load San Francisco time, using local time instead. Current local time is: ", time.Now().In(sfTime).String())
	} else {
		log.Info("Using San Francisco time. Current SF time is: ", time.Now().In(sfTime).String())
	}

	// DB Client
	var dbClient *db.Client
	if saveToDB {
		log.Debug("Migrating DB and Creating DB Client")
		dbClient, err = migrations.LockDBAndMigrate(log, "mysql", Config.DB.Username, Config.DB.Password, Config.DB.Address, Config.DB.Name)
		if err != nil {
			log.Fatalf("database migration exited with error: %s", err.Error())
		}
	}

	var iaasClients []usagedatajob.IaasClient

	// Azure Client
	if caseInsensitiveContains(resources, "Azure") {
		log.Debug("Creating Azure Client")
		if Config.Azure.AccessKey == "" || Config.Azure.EnrollmentNumber == 0 {
			log.Fatal("Azure requires access-key and enrollment-number to be configured")
		}
		azureClient := azure.NewClient(log, sfTime, "https://ea.azure.com/", Config.Azure.AccessKey, Config.Azure.EnrollmentNumber)
		iaasClients = append(iaasClients, azureClient)
	}

	// GCP Client
	if caseInsensitiveContains(resources, "GCP") {
		log.Debug("Creating GCP Client")
		if Config.GCP.ApplicationCredentialsPath == "" || Config.GCP.BucketName == "" {
			log.Fatal("GCP requires bucket-name and application-credentials-path to be configured")
		}
		gcpCredentials, err := ioutil.ReadFile(Config.GCP.ApplicationCredentialsPath)
		if err != nil {
			log.Fatal("Failed to create GCP credentials: ", err.Error())
		}
		gcpClient, err := gcp.NewClient(log, sfTime, gcpCredentials, Config.GCP.BucketName)
		if err != nil {
			log.Fatal("Failed to create GCP client: ", err.Error())
		}
		iaasClients = append(iaasClients, gcpClient)
	}

	// AWS Client
	if caseInsensitiveContains(resources, "AWS") {
		log.Debug("Creating AWS Client")
		if Config.AWS.Region == "" {
			log.Fatal("AWS requires region to be configured")
		}
		if Config.AWS.BucketName == "" {
			log.Fatal("AWS requires bucket-name to be configured")
		}
		if Config.AWS.MasterAccountNumber == int64(0) {
			log.Fatal("AWS requires master_account_number to be configured")
		}
		_ = os.Setenv("AWS_ACCESS_KEY_ID", Config.AWS.AccessKeyID)
		_ = os.Setenv("AWS_SECRET_ACCESS_KEY", Config.AWS.SecretAccessKey)
		sess, err := session.NewSession(&awssdk.Config{Region: awssdk.String(Config.AWS.Region)})
		if err != nil {
			log.Fatal("Failed to create AWS credentials: ", err.Error())
		}
		var reportsDatabase aws.ReportsDatabase
		reportsDatabase = dbClient
		if dbClient == nil {
			reportsDatabase = db.NewNullClient()
		}
		awsClient := aws.NewClient(log, sfTime, Config.AWS.Region, Config.AWS.BucketName, Config.AWS.MasterAccountNumber, s3.New(sess), reportsDatabase)
		iaasClients = append(iaasClients, awsClient)
	}

	usageDataJob := usagedatajob.NewJob(log, sfTime, iaasClients, dbClient, keepFile, saveToDB)

	// BILLING DATA
	if *nowFlag {
		usageDataJob.Run()
		if dbClient != nil {
			_ = dbClient.Close()
		}
		os.Exit(0)
	}

	c := cron.NewWithLocation(sfTime)
	err = c.AddJob("@midnight", usageDataJob)
	if err != nil {
		log.Fatal("Could not create cron job: ", err.Error())
	}
	c.Start()

	// HEALTHCHECK
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Meteorologica is deployed\n\n Last job ran at %s\n\n Next job will run in roughly %s\n    at %s\n\nThere are %d jobs scheduled.",
			c.Entries()[0].Prev.In(sfTime).String(),
			c.Entries()[0].Next.In(sfTime).Sub(time.Now().In(sfTime)).String(),
			c.Entries()[0].Next.In(sfTime).String(),
			len(c.Entries()),
		)
	})
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Config.Port), nil))
}

func parseFlags() (bool, bool, []string) {
	resources := strings.Split(*resourcesFlag, ",")
	return *fileFlag, *dbFlag, resources
}

func configureLog() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
	env := configor.ENV()
	if (*verboseFlag || env == "development") && *verboseFlag != false {
		log.Level = logrus.DebugLevel
	}
	if Config.Rollbar.Token != "" {
		log.Infof("Creating Rollbar hook for %s environment", env)
		log.Hooks.Add(rollrus.NewHook(Config.Rollbar.Token, env))
	}
	return log
}

func caseInsensitiveContains(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if strings.ToLower(hay) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}
