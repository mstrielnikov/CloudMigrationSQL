# CloudMigrationSQL

This Go script helps to backup and import MySQL database to AWS RDS using AWS CLI. The script performs the following steps:

1. Accepts input parameters for MySQL database host, user, password, database name, AWS region, profile, RDS instance identifier, RDS instance class, RDS database name, S3 bucket name, S3 prefix, and IAM role to use while uploading to S3.

2. Creates an AWS STS client and gets details of the current IAM role.

3. Uses the mysqldump command to create a local backup of the MySQL database.

4. Configures AWS RDS and restores the database from the S3 bucket.

5. Removes the local backup file.

# Requirements

* AWS CLI installed and configured with appropriate credentials.

* Boto3 library installed

**Usage**

To use the script, run the following command:
```bash
./cloud_migration_sql.go \
  --mysql-host $MySQL_database_host \
  --mysql-user $MySQL_database_user \
  --mysql-password $MySQL_database_password \
  --mysql-db $MySQL_database_name \
  --aws-profile $AWS_CLI_profile \
  --aws-db-engine $AWS_RDS_instance_identifier \
  --aws-db-engine-v $AWS_RDS_instance_class \
  --aws-db-name $AWS_RDS_database_name \
  --s3-bucket-name $AWS_S3_bucket_name \
  --s3-prefix $AWS_S3_prefix \
  --backup-filename $Local_backup_filename
```

Replace the input parameters with the appropriate values.

# Note

* The script requires appropriate permissions to backup and import the database. Ensure that the IAM role used by the script has the necessary permissions.

* This script does not support SSL connections to the MySQL database.

# References

* [AWS SDK Go v2 docs](https://aws.github.io/aws-sdk-go-v2/docs/)
* [Github: AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2#resources)
* [AWS SDK Go v2: STS](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sts)
* [AWS SDK Go v2: KMS](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/kms)
* [AWS SDK Go v2: S3](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3)
* [AWS SDK Go v2: RDS](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/rds)
