# Keebackup

A CLI binary that back up your KeepassXC file to S3. Also incorporate Terraform and AWS CDK for automatic S3 creation if needed.

## Usage

- Navigate to folder `bin` and run a command in your terminal with these options  
  `keebackup -path <Your KeepassXC Backup>.kdbx -bucket <Bucket name on S3> -region <Region of the bucket>`
- The bucket will be automatically created upon running the command if it hasn't exist at the time of the run.
- To make full use of the tool, you can combine with your system's crontab.
  - For MacOS, open up your terminal and run this command `crontab -e` to edit the crontab jobs.
  - Paste the command you want to run every period of time along with the cron time. For example, this command will backup Passwords.old.kdbx periodically everyday at 23:55.  
    `55 23 * * * /<Path-to-your-binary>/keebackup -path Passwords.old.kdbx -bucket keepassxc-backup-bucket`
  - One disadvantage of this technique is your system need to be online on the time of running.
- Another option is to run the script everytime your system is started up.
  - For MacOS, you can use `launchd`

## Build from source

- Feel free to modify the code all you want, after that to build a new binary, simply run this command in the root folder of this project  
  `go build -o bin/keebackup cmd/main/main.go`
