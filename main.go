package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {

	// Initialize a session in us-east-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String("us-east-2")},
	// )

	// Initialize a session in us-east-2 that the SDK will use to load
	// credentials are read from environement variable
	awsRegion := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secreteKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	roleArn := os.Getenv("AWS_ROLE_ARN")
	kmsKey := os.Getenv("AWS_KEY")
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(accessKey, secreteKey, ""),
	})

	// Create RDS service client
	svc := rds.New(sess)
	result, err := svc.DescribeDBSnapshots(nil)
	if err != nil {
		exitErrorf("Unable to list snapshots, %v", err)
	}

	currentTime := time.Now()
	sanpshotTime := currentTime.Format("2006-01-02")
	gitPrimeSnapshot := "<snapshot prefix>" + sanpshotTime
	IndentifierSanpshotTime := currentTime.Format("20060102")
	exportIdentifier := "export snapshot prefix>" + IndentifierSanpshotTime

	// Getting latest snapsnot from rds snapshot list
	for _, s := range result.DBSnapshots {

		if strings.Contains(*s.DBSnapshotArn, gitPrimeSnapshot) {
			fmt.Printf("* %s with status %s\n",
				aws.StringValue(s.DBSnapshotArn), aws.StringValue(s.Status))

			exportList := []*string{}
			input := &rds.StartExportTaskInput{
				ExportOnly:           exportList, // Optional
				ExportTaskIdentifier: aws.String(exportIdentifier),
				IamRoleArn:           aws.String(roleArn),
				KmsKeyId:             aws.String(kmsKey),
				S3BucketName:         aws.String(bucketName),
				//S3Prefix:             aws.String("/"),  // optional
				SourceArn: aws.String(*s.DBSnapshotArn),
			}

			result, err := svc.StartExportTask(input)

			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case rds.ErrCodeDBSnapshotAlreadyExistsFault:
						fmt.Println(rds.ErrCodeDBSnapshotAlreadyExistsFault, aerr.Error())
					case rds.ErrCodeDBSnapshotNotFoundFault:
						fmt.Println(rds.ErrCodeDBSnapshotNotFoundFault, aerr.Error())
					case rds.ErrCodeInvalidDBSnapshotStateFault:
						fmt.Println(rds.ErrCodeInvalidDBSnapshotStateFault, aerr.Error())
					case rds.ErrCodeSnapshotQuotaExceededFault:
						fmt.Println(rds.ErrCodeSnapshotQuotaExceededFault, aerr.Error())
					case rds.ErrCodeKMSKeyNotAccessibleFault:
						fmt.Println(rds.ErrCodeKMSKeyNotAccessibleFault, aerr.Error())
					default:
						fmt.Println(aerr.Error())
					}
				} else {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
				}
				return
			}

			fmt.Println(result)

		}

	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
