# cp-aws-rds-snapshot-to-s3

This utility helps to copy the current day aws rds snapshot to the S3 bucket. This is to maintain the backup of rds snapshot in another region for DR needs.

**Note: This utility uses the aws go sdk to perform operation in asws**

# Getting Started

Clone this package to your workspace

go get github.com/gaurav1der/cp-aws-rds-snapshot-to-s3

## Configuration

Follow environmaent parameter need to be set for this ulity to work

  
    export AWS_ACCESS_KEY_ID=<>
    export AWS_SECRET_ACCESS_KEY=<>
    export AWS_REGION=<>
    export AWS_ROLE_ARN=<>
    export AWS_KEY=<>
    export AWS_BUCKET_NAME=<>
    export RDS_SNAPSHOT=<>
    export EXPORT_RDS_SNAPSHOT=<>
    
## Run    

go run main.go
    
    
  


