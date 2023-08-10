# Keebackup

A CLI binary that back up your KeepassXC file to S3. Also incorporate Terraform and AWS CDK for automatic S3 creation if needed.

## Usage

- Navigate to folder `bin` and run a command in your terminal with these options  
  `keebackup -path <Your KeepassXC Backup>.kdbx -bucket <Bucket name on S3>`
- The bucket will be automatically created upon running the command using Terraform if it doesn't exists
