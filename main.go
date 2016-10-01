package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/challiwill/meteorologica/aws"
	"github.com/challiwill/meteorologica/azure"
	"github.com/challiwill/meteorologica/db"
	"github.com/challiwill/meteorologica/gcp"
	"github.com/challiwill/meteorologica/usagedatajob"
	"github.com/robfig/cron"
)

type Client interface{}

var (
	azureFlag     = flag.Bool("azure", false, "Only retrieve Azure data (by default Azure, AWS, and GCP data is retrieved)")
	gcpFlag       = flag.Bool("gcp", false, "Only retrieve GCP data (by default Azure, AWS, and GCP data is retrieved)")
	awsFlag       = flag.Bool("aws", false, "Only retrieve AWS data (by default Azure, AWS, and GCP data is retrieved)")
	nowFlag       = flag.Bool("now", false, "Run job now (instead of waiting for cron job to kick off at midnight)")
	verboseFlag   = flag.Bool("v", false, "Log at Debug level")
	fileFlag      = flag.Bool("file", false, "Keep the generated, normalized CSV file locally")
	localOnlyFlag = flag.Bool("local", false, "Do not connect to any services (overrides -db and -bucket)")
	dbFlag        = flag.Bool("db", false, "Save the data to the database")
	bucketFlag    = flag.Bool("bucket", false, "Save the data as a .csv to the provided GCP bucket")
	migrateFlag   = flag.Bool("migrate", false, "Run migrations then exit")
)

func main() {
	flag.Parse()

	getAll := !(*azureFlag) && !(*gcpFlag) && !(*awsFlag)
	getAzure := *azureFlag || getAll
	getGCP := *gcpFlag || getAll
	getAWS := *awsFlag || getAll

	keepFile := *fileFlag
	localOnly := *localOnlyFlag
	dbf := *dbFlag
	bkf := *bucketFlag
	saveToDB := (dbf || (!dbf && !bkf)) && !localOnly
	saveToBucket := (bkf || (!dbf && !bkf)) && !localOnly

	migrate := *migrateFlag

	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
	if *verboseFlag {
		log.Level = logrus.DebugLevel
	}

	sfTime, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		sfTime = time.Now().Location()
		log.Warn("Failed to load San Francisco time, using local time instead. Current local time is: ", time.Now().In(sfTime).String())
	} else {
		log.Info("Using San Francisco time. Current SF time is: ", time.Now().In(sfTime).String())
	}

	// DB Client
	var dbClient *db.Client
	if saveToDB || migrate {
		log.Debug("Creating DB Client")
		dbClient, err = db.NewClient(log, os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_ADDRESS"), os.Getenv("DB_NAME"))
		if err != nil {
			log.Fatal("Failed to create database client: ", err.Error())
		}
		err = dbClient.Ping()
		if err != nil {
			log.Fatal("Failed to create database connection: ", err.Error())
		}
	}

	if migrate {
		err := dbClient.Migrate()
		if err != nil {
			log.Error(err.Error())
			log.Fatal("Failed to migrate database: ", err.Error())
		}
		os.Exit(0)
	}

	var iaasClients []usagedatajob.IaasClient
	var bucketClient usagedatajob.BucketClient

	// Azure Client
	if getAzure {
		log.Debug("Creating Azure Client")
		accessKey := os.Getenv("AZURE_ACCESS_KEY")
		enrollmentNumber := os.Getenv("AZURE_ENROLLMENT_NUMBER")
		if accessKey == "" || enrollmentNumber == "" {
			log.Fatal("Azure requires AZURE_ACCESS_KEY and AZURE_ENROLLMENT_NUMBER environment variables to be set")
		}
		azureClient := azure.NewClient(log, sfTime, "https://ea.azure.com/", accessKey, enrollmentNumber)
		iaasClients = append(iaasClients, azureClient)
	}

	// GCP Client
	if getGCP {
		log.Debug("Creating GCP Client")
		credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		bucketName := os.Getenv("GCP_BUCKET_NAME")
		if credentialsFile == "" || bucketName == "" {
			log.Fatal("GCP requires GCP_BUCKET_NAME and GOOGLE_APPLICATION_CREDENTIALS environment variables to be set")
		}
		gcpCredentials, err := ioutil.ReadFile(credentialsFile)
		if err != nil {
			log.Fatal("Failed to create GCP credentials: ", err.Error())
		}
		gcpClient, err := gcp.NewClient(log, sfTime, gcpCredentials, bucketName)
		if err != nil {
			log.Fatal("Failed to create GCP client: ", err.Error())
		}
		iaasClients = append(iaasClients, gcpClient)
	}

	// BucketClient
	if saveToBucket {
		log.Debug("Creating Bucket Client (GCP)")
		credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if credentialsFile == "" {
			log.Fatal("To store the file Meteorologica requires GOOGLE_APPLICATION_CREDENTIALS environment variable to be set")
		}
		bucketName := os.Getenv("STORAGE_BUCKET_NAME")
		if credentialsFile == "" {
			log.Fatal("To store the file Meteorologica requires STORAGE_BUCKET_NAME environment variable to be set")
		}
		gcpCredentials, err := ioutil.ReadFile(credentialsFile)
		if err != nil {
			log.Fatal("Failed to create Bucket (GCP) credentials: ", err.Error())
		}
		gcpClient, err := gcp.NewClient(log, sfTime, gcpCredentials, bucketName)
		if err != nil {
			log.Fatal("Failed to create Bucket (GCP) client: ", err.Error())
		}
		bucketClient = gcpClient
	}

	// AWS Client
	if getAWS {
		log.Debug("Creating AWS Client")
		az := os.Getenv("AWS_REGION")
		bucketName := os.Getenv("AWS_BUCKET_NAME")
		accountNumber := os.Getenv("AWS_MASTER_ACCOUNT_NUMBER")
		if az == "" || bucketName == "" || accountNumber == "" {
			log.Fatal("AWS requires AWS_REGION, AWS_BUCKET_NAME, and AWS_MASTER_ACCOUNT_NUMBER environment variables to be set")
		}
		sess, err := session.NewSession(&awssdk.Config{Region: awssdk.String(az)})
		if err != nil {
			log.Fatal("Failed to create AWS credentials: ", err.Error())
		}
		awsClient := aws.NewClient(log, sfTime, az, bucketName, accountNumber, sess)
		iaasClients = append(iaasClients, awsClient)
	}

	usageDataJob := usagedatajob.NewJob(log, sfTime, iaasClients, bucketClient, dbClient, keepFile, saveToBucket, saveToDB)

	// BILLING DATA
	if *nowFlag {
		usageDataJob.Run()
		if dbClient != nil {
			dbClient.Close()
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
