Tool to collect and standardize billing information from multiple IAAS's.

Right now Google Cloud Platform, Amazon Web Services, and Microsoft Azure.

##Environment Needed:

###GCP:
You need to generate and download a
[service_account_credential](https://cloud.google.com/storage/docs/authentication#service_accounts).
Provide a path to the file as an environment variable.

You must provide the bucket that holds the billing files. Files are assumed to have the format `Billing-YYYY-MM-DD.csv`.

```
GOOGLE_APPLICATION_CREDENTIALS=./path/to/service_account_credential.json
GCP_BUCKET_NAME=my-bucket
```

###AWS:


###Azure:
You need to provide the API Access Key and your Enrollment Number as environment variables.

```
AZURE_ENROLLMENT_NUMBER=12345
AZURE_ACCESS_KEY=api-access-key
```

