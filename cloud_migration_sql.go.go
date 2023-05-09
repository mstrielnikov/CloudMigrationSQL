package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func main() {
	// Parse command-line arguments
	mysqlHost := flag.String("mysql-host", "", "MySQL database host.")
	mysqlUser := flag.String("mysql-user", "", "MySQL database user.")
	mysqlPassword := flag.String("mysql-password", "", "MySQL database password.")
	mysqlDb := flag.String("mysql-db", "", "MySQL database name.")

	awsDbName := flag.String("aws-db-name", "", "AWS RDS database name.")
	awsDbEngine := flag.String("aws-db-engine", "mysql", "AWS RDS database engine name.")
	awsDbInstanceClass := flag.String("aws-db-instance-class", "", "AWS RDS instance class.")

	s3BucketName := flag.String("s3-bucket-name", "", "AWS S3 bucket name where dump uploaded.")
	s3BucketPrefix := flag.String("s3-prefix", "", "AWS S3 prefix where dump uploaded.")

	backupFilepath := flag.String("backup-filename", "backup.sql", "Local backup filename.")
	flag.Parse()

	// Load the SDK's configuration from environment and shared config (~/.aws/config), and
	// create the client with this.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	// Create an STS client using the loaded SDK configuration
	stsClient := sts.NewFromConfig(cfg)

	// Call the GetCallerIdentity API to get details of the current IAM role
	stsOutput, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		panic(fmt.Sprintf("failed to get caller identity, %v", err))
	}

	// Print the ARN of the current IAM role
	fmt.Printf("Current IAM role ARN: %s\n", *stsOutput.Arn)

	// Make sqldump backup
	backupCmd := exec.Command("mysqldump", "-h", *mysqlHost, "-u", *mysqlUser, "-p"+*mysqlPassword, *mysqlDb, "--result-file="+*backupFilepath)
	if err := backupCmd.Run(); err != nil {
		log.Fatalf("Error backing up MySQL database: %v", err)
		os.Exit(1)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Create S3 uploader
	s3Uploader := manager.NewUploader(s3Client)

	// Upload backup to S3
	backupFile, err := os.Open(*backupFilepath)
	if err != nil {
		log.Fatalf("Error opening backup file: %v:", err)
		os.Exit(1)
	}
	defer backupFile.Close()

	s3output, err := s3Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(*s3BucketName),
		Key:    aws.String(*s3BucketPrefix),
		Body:   backupFile,
	})

	if err != nil {
		log.Fatalf("Error downloading S3 object: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Uploaded file stats")
	fmt.Println("The object key of the newly created object:", s3output.Key)
	fmt.Println("Entity tag for the uploaded object:", s3output.ETag)
	fmt.Println("The base64-encoded, 160-bit SHA-1 digest of the object:", s3output.ChecksumSHA1)
	fmt.Println("Does uploaded object uses a S3 Bucket Key for AWS KMS (SSE-KMS):", s3output.BucketKeyEnabled)

	// Create RDS client
	rdsClient := rds.NewFromConfig(cfg)

	// Create a new RDS database instance from the MySQL dump file
	restoreInput := &rds.RestoreDBInstanceFromS3Input{
		DBInstanceIdentifier: aws.String(*awsDbName),
		DBInstanceClass:      aws.String(*awsDbInstanceClass),
		Engine:               aws.String(*awsDbEngine),
		AllocatedStorage:     aws.Int32(10),
		S3BucketName:         aws.String(*s3BucketName),
		S3IngestionRoleArn:   aws.String(*stsOutput.Arn),
		S3Prefix:             aws.String(*s3BucketPrefix),
		DBName:               aws.String(*awsDbName),
	}
	_, err = rdsClient.RestoreDBInstanceFromS3(context.TODO(), restoreInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error restoring RDS instance: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("MySQL dump file imported successfully.")
	// Add descriptive output
}
