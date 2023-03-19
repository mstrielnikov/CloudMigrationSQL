import os
import subprocess
import boto3

parser = argparse.ArgumentParser(description="Backup and import MySQL database to AWS RDS.")
parser.add_argument("--mysql-host",            required=True,  help="MySQL database host.")
parser.add_argument("--mysql-user",            required=True,  help="MySQL database user.")
parser.add_argument("--mysql-password",        required=True,  help="MySQL database password.")
parser.add_argument("--mysql-db",              required=True,  help="MySQL database name.")

parser.add_argument("--aws-region",            required=True,  help="AWS region where the RDS instance is located.")
parser.add_argument("--aws-profile",           required=True,  help="AWS CLI profile to use.")
parser.add_argument("--aws-db-instance",       required=True,  help="AWS RDS instance identifier.")
parser.add_argument("--aws-db-instance-class", required=True,  help="AWS RDS instance identifier.")
parser.add_argument("--aws-db-name",           required=True,  help="AWS RDS database name.")

parser.add_argument("--s3-bucket-name",        required=True,  help="AWS S3 bucket name where dump uploaded.")
parser.add_argument("--s3-prefix",             required=True,  help="AWS S3 prefix where dump iploaded.")
parser.add_argument("--s3-iam-role",           required=True,  help="AWS role to use while upload to S3")

parser.add_argument("--backup-filename",       required=True,  help="Local backup filename.", default="backup.sql",)

args = parser.parse_args()

MySqlHost          = args.mysql_host
MySqlUser          = args.mysql_user
MySqlPassword      = args.mysql_password
MySqlDb            = args.mysql_db

AwsRegion          = args.aws_region
AwsProfile         = args.aws_profile
AwsDbInstance      = args.aws_db_instance
AwsDbInstanceClass = args.aws_db_instance_class
AwsDbName          = args.aws_db_name

AwsS3BucketName    = args.aws_bucket_name
AwsS3Prefix        = args.aws_s3_prefix
AwsS3IamRole       = args.aws_s3_iam_role

BackupFilename     = args.backup_filename

# Create an STS client
sts = boto3.client('sts')

# Get the details of the current IAM role
response = sts.get_caller_identity()

# ARN of the current IAM role
AwsRoleArn  = response['Arn']
AwsRoleName = response['Name']

# Create local backup
dump_cmd = f"mysqldump -h {MySqlHost} -u {MySqlUser} -p{MySqlPassword} {MySqlDb} > {BackupFilename}"
subprocess.run(dump_cmd, shell=True, check=True)

# AWS RDS configuration
rds = boto3.client("rds", region_name=AwsRegion, profile_name=AwsProfile)
restore_args = {
    "DBInstanceIdentifier": AwsDbInstance,
    "DBInstanceClass": AwsDbInstanceClass,
    "Engine": "mysql",
    "AllocatedStorage": 10,
    "S3BucketName": AwsS3BucketName,
    "S3IngestionRoleArn": f"arn:aws:iam::{AwsRoleArn}:role/{AwsRoleName}",
    "S3Prefix": f"{AwsS3Prefix}{BackupFilename}",
    "DBName": AwsDbName,
}
response = rds.restore_db_instance_from_s3(**restore_args)

# Remove local backup
os.remove(BACKUP_FILENAME)

print("Database backup and import successful!")

