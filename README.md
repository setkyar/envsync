# envsync

Sync your environment files securely with public/private key encryption via AWS S3.

## Overview

`envsync` is a CLI tool designed to securely synchronize your `.env` files across different machines. It uses public/private key encryption for security and AWS S3 for storage.

![Diagram](diagram.png)

## Commands

- `init`: Initialize your `envsync`. This command sets up public/private keys and configures AWS S3.
- `push`: Push your `.env` file from the current directory to the S3 bucket.
- `pull`: Pull your `.env` file from the S3 bucket to the current directory.

## Getting Started

### Prerequisites

1. **S3 Bucket**: Set up an S3 bucket, e.g., `your-s3-bucket`.
2. **IAM User**: Create an IAM user in AWS and attach the following policy for necessary permissions:

   ```json
   {
       "Version": "2012-10-17",
       "Statement": [
           {
               "Sid": "VisualEditor0",
               "Effect": "Allow",
               "Action": [
                   "s3:PutObject",
                   "s3:PutObjectAcl",
                   "s3:GetObject"
               ],
               "Resource": "arn:aws:s3:::your-s3-bucket/*"
           }
       ]
   }
   ```

3. **AWS Credentials**: Note down the IAM user’s `access_key_id` and `secret_access_key`.

### Initialization

Run `envsync init` and input the AWS configuration when prompted. This will set up the necessary keys and configuration for `envsync`.

### Usage

**Pushing `.env` File**:

To push the `.env` file from your current directory to S3, run:

```sh
envsync push --name=your_project_name
```

This command encrypts your `.env` file and stores it at `your-s3-bucket/your_project_name/.env` in S3.

**Pulling `.env` File**:

To pull the `.env` file from S3, run:

```sh
envsync pull --name=your_project_name
```

### Collaborating with Team Members

For team collaboration, follow these steps:

1. **IAM Permissions**: 
Ensure team members have the necessary IAM permissions (refer to the policy mentioned above).

2. **Key Sharing**: 
Share the public and private keys located in `$HOME/.envsync/` with your team or you create your own `public/private` key pair and configure to use via your own `config.yaml` file and share with the team. You can create your own key pair using the following command...
```
ssh-keygen -t rsa -b 2048 -f private_key.pem && mv private_key.pem.pub public_key.pem && ssh-keygen -p -m PEM -f private_key.pem
```

3. **Team Setup**: Get `private_key.pem` and `public_key.pem` and configure your `config.yaml` like the following.

```yaml
aws:
  region: ap-southeast-1
  s3_bucket: your-s3-bucket
  access_key_id: your-aws-access-key
  secret_access_key: your-aws-secret-key
envsync:
  private_key: ~/.envsync/private_key.pem # Replace with private_key path
  public_key: ~/.envsync/public_key.pem   # Replace with public_key path
```

And run pull or pull like the following
```sh
envsync pull --name=your_project_name --config=config.yaml
```