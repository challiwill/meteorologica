CODE CLIMATE
* backfill tests where necessary (pending tests)
* test time better (eg in gcp client file name generation)
* make Usage structs have good types

DATABASE
* set up proper db migration framework
* set up local db for testing
* find better way to make unique key - currently there could be rows that appear to be duplicates but aren't
* make columns not null

CONFIGURATION
* merge flags with configuration struct
* properly get DB credentials from env
* allow storage bucket to be different type (eg aws) if desired (or remove storage bucket entirely)
* 'last job ran' should be stored in DB

FEATURES
* extract each IAAS as a resource or something so that new resources can be added to calculate billing info from (eg Pagerduty)

PERFORMANCE
* can use streams of reading from csv and writing to database to do it all concurrently instead of in blocks
* can use better types for slices of usages (eg in AWS) that have smarter methods for serching for duplicates

OTHER SERVICES
* configurable frontend as a micro service
* usage alerting service
