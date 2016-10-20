CODE CLIMATE
* backfill tests where necessary (pending tests)
* test time better (eg in gcp client file name generation)

DATABASE
* set up local db for testing
* explore best types for database
* rename lock to be more robust

CONFIGURATION
* merge flags with configuration struct
* properly get DB credentials from env
* 'last job ran' should be stored in DB

FEATURES
* extract each IAAS as a resource or something so that new resources can be added to calculate billing info from (eg Pagerduty)
* get warning about inserting duplicate data

PERFORMANCE
* can use streams of reading from csv and writing to database to do it all concurrently instead of in blocks

OTHER SERVICES
* configurable frontend as a micro service
* usage alerting service

