# Meteorologica
Tool to collect and standardize billing information from multiple IAAS's.
Right now Google Cloud Platform, Amazon Web Services, and Microsoft Azure.

Currently the default behavior is as follows:
* Meteorologica collects billing information from the location where it is published (AWS bucket, GCP bucket, Azure API)
* Meteorologica normalizes the data
* Meteorologica inserts the data into the given MySQL database (right now this step can take a long time because it goes row by row)
* Metorologica saves the data to a local file `YEAR-MONTH-normalized-billing-data.csv` (eg `2016-September-normalized-billing-data.csv`)
* Meteorologica uploads this file to the specified bucket (currently a GCP bucket only)

*NB: Currently if a insert is made and there is a collision (same hash) the usage and cost will be updated to the new values*

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
-now      Run the task now instead of waiting for next scheduled job (next midnight)
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

To set the environment between `production`, or the default `development` set the environment variable
(this only matters for logging at the moment, `development` logs at the `debug` level, `production` logs at `warn` level):
```
M_ENV=production
```

Configuration of the following integrations is available through `Environment Variables`, `configuration/meteorologica.{ENVIRONMENT}.yml`, `configuration/meteorologica.yml` in that order of priority.
If setting an environment variable it is all caps, and of the form `M_SERVICE_VARIABLE_NAME` where `M` is a prefix for this application,
`SERVICE` is the respective service (eg `GCP`, `AZURE`, `AWS`, `ROLLBAR`),
and `VARIABLE` is the value that needs to get conveyed as described below, where camel case is translated to underscore separated (eg: `M_AZURE_ACCESS_KEY=my-access-key`).
If creating a `.yml` configuration file set the variables as shown below.

###GCP:
You need to generate and download a
[service_account_credential](https://cloud.google.com/storage/docs/authentication#service_accounts).
Provide a path to the file as a variable.
The file should probably be uploaded to wherever the app is running along with the app (for example in a `credentials/` directory).

You must provide the name of the bucket that holds the billing information files, and the name of the bucket where you would like the final .csv to end up.
The billing files are assumed to have the naming format `Billing-YYYY-MM-DD.csv`.
``` yml
gcp:
  bucket-name: my-bucket
  application-credentials-path: ./path/to/service_account_credential.json
  storage-bucket-name: my-final-bucket
```

###AWS:
You need to provide credentials and configuration options.
The master account number is the account number for the billing management account.
``` yml
aws:
  region: us-east-1
  master-account-number: 12345
  bucket-name: bucket-name
  access-key-id: access-key-id
  secret-access-key: secret-access-key
```


###Azure:
You need to provide the API Access Key and your Enrollment Number.
``` yml
azure:
  enrollment-number: 12345
  access-key: api-access-key
```

### MySQL Database:
If you would like to connect to a MySQL database the following variables must be set as needed:
``` yml
db:
  username: account-username
  password: account-password
  address: hostname:port
  name: database-name
```

The expected schema is (`VARCHAR` and `CHAR` should be adjusted as necessary):
``` sql
CREATE TABLE database-name.iaas_billing (
  AccountNumber VARCHAR(15),
  AccountName VARCHAR(30),
  Day TINYINT(2),
  Month CHAR(9),
  Year SMALLINT(4),
  ServiceType VARCHAR(30),
  UsageQuantity DOUBLE,
  Cost DECIMAL(15,2),
  Region VARCHAR(10),
  UnitOfMeasure VARCHAR(10),
  IAAS VARCHAR(10),
  UNIQUE KEY(AccountNumber, Day, Month, Year, ServiceType, UsageQuantity, Region, IAAS)
);
```

## Migrations
The current way to run migrations is rather janky. Create a file `migrations/iaas_billing.sql` and then push the app with the `-migrate` flag set:
```
cf push meteorologica -c "meteorologica -migrate"
```
After this push the app again or restage as appropriate.
Someday this will be made better.

## Timeline
I am keeping a running list of tasks to accomplish in the [TODO](TODO.md) file.
