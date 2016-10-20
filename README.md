# Meteorologica
Tool to collect and standardize billing information from multiple IAAS's.
Right now Google Cloud Platform, Amazon Web Services, and Microsoft Azure.

Currently the default behavior is as follows:

* For Each IAAS:
  * Meteorologica collects billing information from the location where it is published (AWS bucket, GCP bucket, or Azure API)
  * Meteorologica normalizes the data
  * Meteorologica inserts the data into the given MySQL database (right now this step can take a long time because it goes row by row)

*NB: Currently if a insert is made and there is a collision (based on some id of max len 25 characters) the usage quantity and cost will be updated to the new values*

## Use
You can use this tool to collect billing info from all your IAAS's just by running the file:
```
go run main.go
```

By default the app is configured to save the data to the configured database.
To keep a local version of the data as a CSV file pass in the `-file` flag:
```
go run main.go -file
```

By default the app collects data from GCP, AWS, and Azure.
To collect billing data from only one (or more) IAAS you can pass the `-resources` flag, for example:
```
go run main.go -resources=aws,gcp
```

All flags:
```
-resources  A comma seperated list of resource to retrieve billing information from. If none are specified the default is AWS, GCP, and Azure
-v          Verbose mode, log at the debug level
-file       Save the generated and normalized data in a local .csv file
-db         Save the data to the database (by default this happens, this flag exists so you can set it to false)
-cron       Run job periodically every day at midnight
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
You need to provide credentials for your MySQL database:
``` yml
db:
  username: account-username
  password: account-password
  address: hostname:port
  name: database-name
```

The recommended schema for the table is:
```
+-----------------+--------------+------+-----+---------+-------+
| Field           | Type         | Null | Key | Default | Extra |
+-----------------+--------------+------+-----+---------+-------+
| id              | varchar(25)  | NO   | PRI | NULL    |       |
| account_number  | varchar(255) | NO   |     | NULL    |       |
| account_name    | varchar(255) | YES  |     | NULL    |       |
| day             | tinyint(2)   | NO   |     | NULL    |       |
| month           | tinyint(2)   | NO   |     | NULL    |       |
| year            | smallint(4)  | NO   |     | NULL    |       |
| service_type    | varchar(255) | NO   |     | NULL    |       |
| region          | varchar(255) | YES  |     | NULL    |       |
| resource        | varchar(255) | NO   |     | NULL    |       |
| usage_quantity  | double       | NO   |     | NULL    |       |
| unit_of_measure | varchar(255) | YES  |     | NULL    |       |
| cost            | double       | NO   |     | NULL    |       |
+-----------------+--------------+------+-----+---------+-------+
```

## Migrations
Migrations are run when the app starts up.
The app protects against conflicting migrations by getting a database lock.
All other app instances will wait for the lock then exit the migrations with a no-op.
Migrations are run from the `db/migrations/migrations.go` file as defined in the same directory.

## Timeline
I am keeping a running list of tasks to accomplish in the [TODO](TODO.md) file.
