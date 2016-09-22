Tool to collect and standardize billing information from multiple IAAS's.

Right now Google Cloud Platform, Amazon Web Services, and Microsoft Azure.

## Use
You can use this tool to collect billing info from all your IAAS's just by running the file:
```
go run main.go
```

The app is configured to collect, standardize, and upload a consolidated csv data file at midnight PST each day. If you would like it to run immediately pass it the `-now` flag, for example:
```
go run main.go -now
```

By default the app is configured to send the standardized file to the given GCP bucket. To keep the file locally and not ship it to GCP pass in the `-file` flag:
```
go run main.go -file
```

By default the app collects data from GCP, AWS, and Azure. To collect billing data from only one (or more) IAAS you can pass a flag (currently `-gcp`, `-azure`, or `-aws`), for example:
```
go run main.go -aws
```


All flags:
```
-azure    Retrive Azure data
-aws      Retrieve AWS data
-gcp      Retrieve GCP data
-v        Verbose mode, log at the debug level
-file     Save the generated and normalized data in a local .csv file
-local    Do not connect to any services, specifically do not send the data to the database or the GCP bucket (this overrides '-db' and '-bucket')
-db       Save the data to the database (by default this happens, you would only set this flag to send the data to the database and not the GCP bucket)
-bucket   Save the data to the GCP bucket (by default this happens, you would only set this flag to send the data to the GCP bucket and not the database)
```

## Deployment

To push to cloudfoundry run the following command from within the app directory:
```
cf push meteorologica -b https://github.com/cloudfoundry/go-buildpack.git
```

There is a healthcheck that you can use to confirm the app is running, see when the last data collection job ran, and when the next job will run. You can access it at [/healthcheck](http://meteorologica.cfapps.io/healthcheck).

Metrics and logs available at [https://metrics.run.pivotal.io](https://metrics.run.pivotal.io)

##Environment Needed:
Be careful not to upload any credentials to Github as this repository is Public.

###GCP:
You need to generate and download a
[service_account_credential](https://cloud.google.com/storage/docs/authentication#service_accounts).
Provide a path to the file as an environment variable. The file should probably be uploaded to wherever the app is running along with the app (for example in a `credentials/` directory).

You must provide the name of the bucket that holds the billing information files. The billing files are assumed to have the naming format `Billing-YYYY-MM-DD.csv`.

```
GOOGLE_APPLICATION_CREDENTIALS=./path/to/service_account_credential.json
GCP_BUCKET_NAME=my-bucket
```

###AWS:
You need to provide credentials and configuration options. The master account number is the account number for the billing management account.

```
AWS_REGION=us-east-1
AWS_MASTER_ACCOUNT_NUMBER=12345
AWS_BUCKET_NAME=bucket-name
AWS_ACCESS_KEY_ID=acess-key-id
AWS_SECRET_ACCESS_KEY=secret-access-key
```


###Azure:
You need to provide the API Access Key and your Enrollment Number as environment variables.

```
AZURE_ENROLLMENT_NUMBER=12345
AZURE_ACCESS_KEY=api-access-key
```

