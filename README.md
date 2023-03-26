# psql-cloud-migration

This Python script helps to backup and import MySQL database to AWS RDS using AWS CLI. The script performs the following steps:

1. Accepts input parameters for MySQL database host, user, password, database name, AWS region, profile, RDS instance identifier, RDS instance class, RDS database name, S3 bucket name, S3 prefix, and IAM role to use while uploading to S3.

2. Creates an AWS STS client and gets details of the current IAM role.

3. Uses the mysqldump command to create a local backup of the MySQL database.

4. Configures AWS RDS and restores the database from the S3 bucket.

5. Removes the local backup file.

Requirements:

* AWS CLI installed and configured with appropriate credentials.

* Python 3.x

* Boto3 library installed
