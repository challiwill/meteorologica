# Meteorologica
Tool to collect and standardize billing information from multiple IAAS's.
Right now Google Cloud Platform, Amazon Web Services, and Microsoft Azure.

Currently the default behavior is as follows:
* Meteorologica collects billing information from the location where it is published (AWS bucket, GCP bucket, Azure API)
* Meteorologica normalizes the data
* Meteorologica inserts the data into the given MySQL database (right now this step can take a long time because it goes row by row)
* Metorologica saves the data to a local file `YEAR-MONTH-normalized-billing-data.csv`
* Meteorologica uploads this file to the specified bucket (currently a GCP bucket only)

*NB: Eventually we would like to upload the billing information to the database as a load from file using the csv file.
This should make the save to database step a lot faster.
Currently some database as a service providers do not yet support this feature.*

## Use
You can use this tool to collect billing info from all your IAAS's just by running the file:
```
go run main.go
```

The app is configured to collect, standardize, and upload a consolidated csv data file at midnight PST each day.
If you would like it to run immediately pass it the `-now` flag, for example:
```
go run main.go -now
```

By default the app is configured to send the standardized file to the given GCP bucket.
To keep the file locally and not delete it after processing pass in the `-file` flag:
```
go run main.go -file
```

By default the app collects data from GCP, AWS, and Azure.
To collect billing data from only one (or more) IAAS you can pass a flag (currently `-gcp`, `-azure`, or `-aws`), for example:
```
go run main.go -aws -gcp
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

There is a healthcheck that you can use to confirm the app is running, see when the last data collection job ran, and when the next job will run.
You can access it at [/healthcheck](http://meteorologica.cfapps.io/healthcheck).

Metrics and logs available at [https://metrics.run.pivotal.io](https://metrics.run.pivotal.io)

##Environment Needed:
Be careful not to upload any credentials to Github as this repository is Public.

###GCP:
You need to generate and download a
[service_account_credential](https://cloud.google.com/storage/docs/authentication#service_accounts).
Provide a path to the file as an environment variable.
The file should probably be uploaded to wherever the app is running along with the app (for example in a `credentials/` directory).

You must provide the name of the bucket that holds the billing information files. The billing files are assumed to have the naming format `Billing-YYYY-MM-DD.csv`.
```
GOOGLE_APPLICATION_CREDENTIALS=./path/to/service_account_credential.json
GCP_BUCKET_NAME=my-bucket
```

###AWS:
You need to provide credentials and configuration options.
The master account number is the account number for the billing management account.
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

### MySQL Database:
If you would like to connect to a MySQL database the following environment variables must be set as needed:
```
DB_USERNAME=account-username
DB_PASSWORD=account-password
DB_ADDRESS=hotname:port
DB_NAME=database-name
```

The expected schema is (`VARCHAR` and `CHAR` should be adjusted as necessary):
```
CREATE TABLE iaas_billing (
  AccountNumber VARCHAR(15),
  AccountName VARCHAR(30),
  Day CHAR(2),
  Month CHAR(9),
  Year CHAR(5),
  ServiceType VARCHAR(30),
  UsageQuantity VARCHAR(10),
  Cost VARCHAR(10),
  Region VARCHAR(10),
  UnitOfMeasure VARCHAR(10),
  IAAS VARCHAR(10),
  UNIQUE KEY(AccountNumber, Day, Month, Year, ServiceType, UsageQuantity, Region, IAAS)
);
```
