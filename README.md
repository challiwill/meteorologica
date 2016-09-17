Tool to collect and standardize billing information from multiple IAAS's.

Right now Google Cloud Platform, Amazon Web Services, and Microsoft Azure.

## Use
You can use this tool to collect billing info from all your IAAS's just by running the file:
```
go run main.go
```

It also accepts flags to collect billing data from only one (or more) IAAS: 
```
go run main.go -aws -azure
```

##Environment Needed:

###GCP:
You need to generate and download a
[service_account_credential](https://cloud.google.com/storage/docs/authentication#service_accounts).
Provide a path to the file as an environment variable.

You must provide the name of the bucket that holds the billing information files. The billing files are assumed to have the naming format `Billing-YYYY-MM-DD.csv`.

```
GOOGLE_APPLICATION_CREDENTIALS=./path/to/service_account_credential.json
GCP_BUCKET_NAME=my-bucket
```

###AWS:
You need to provide credentials and configuration options. The master account number is the account number for the billing management account.

```
export AWS_REGION=us-east-1
export AWS_MASTER_ACCOUNT_NUMBER=12345
export AWS_BUCKET_NAME=bucket-name
export AWS_ACCESS_KEY_ID=acess-key-id
export AWS_SECRET_ACCESS_KEY=secret-access-key
```


###Azure:
You need to provide the API Access Key and your Enrollment Number as environment variables.

```
AZURE_ENROLLMENT_NUMBER=12345
AZURE_ACCESS_KEY=api-access-key
```

