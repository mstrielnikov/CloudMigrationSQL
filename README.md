# SQL cloud migration

This Python script helps to backup and import MySQL database to AWS RDS using AWS CLI. The script performs the following steps:

1. Accepts input parameters for MySQL database host, user, password, database name, AWS region, profile, RDS instance identifier, RDS instance class, RDS database name, S3 bucket name, S3 prefix, and IAM role to use while uploading to S3.

2. Creates an AWS STS client and gets details of the current IAM role.

3. Uses the mysqldump command to create a local backup of the MySQL database.

4. Configures AWS RDS and restores the database from the S3 bucket.

5. Removes the local backup file.

# Requirements

* AWS CLI installed and configured with appropriate credentials.

* Python 3.x

* Boto3 library installed

**Usage**

To use the script, run the following command:
```bash
python3 sql_migration_cloud.py \
  --mysql-host $MySQL_database_host \
  --mysql-user $MySQL_database_user \
  --mysql-password $MySQL_database_password \
  --mysql-db $MySQL_database_name \
  --aws-region $AWS_region \
  --aws-profile $AWS_CLI_profile \
  --aws-db-instance $AWS_RDS_instance_identifier \
  --aws-db-instance-class $AWS_RDS_instance_class \
  --aws-db-name $AWS_RDS_database_name \
  --s3-bucket-name $AWS_S3_bucket_name \
  --s3-prefix $AWS_S3_prefix \
  --s3-iam-role $AWS_role_to_use_while_upload_to_S3 \
  --backup-filename $Local_backup_filename
```

Replace the input parameters with the appropriate values.

# Note

* The script requires appropriate permissions to backup and import the database. Ensure that the IAM role used by the script has the necessary permissions.

* The script deletes the local backup file after importing the database to AWS RDS. Ensure that you have a backup of the database before running the script.

* This script does not support SSL connections to the MySQL database.

# References

* [boto3 docs](https://boto3.amazonaws.com/v1/documentation/api/latest/index.html)

* [AWS CLI docs](https://aws.amazon.com/cli/)

* [AWS RDS docs](https://aws.amazon.com/rds/)
