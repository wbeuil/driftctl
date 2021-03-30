provider "aws" {
    region  = "us-east-1"
}
terraform {
    required_providers {
        aws = {
            version = "~> 3.19.0"
        }
    }
}

resource "random_string" "prefix" {
  length  = 6
  upper   = false
  special = false
}

resource "aws_s3_bucket" "foobar" {
  bucket = "${random_string.prefix.result}.driftctl-test.com"
  policy = jsonencode({
    Version = "2012-10-17"
    Id      = "MYBUCKETPOLICY"
    Statement = [
        {
            Sid       = "IPAllow"
            Effect    = "Deny"
            Principal = "*"
            Action    = "s3:*"
            Resource = [
                "arn:aws:s3:::${random_string.prefix.result}.driftctl-test.com",
                "arn:aws:s3:::${random_string.prefix.result}.driftctl-test.com/*"
            ]
            Condition = {
                IpAddress = {
                    "aws:SourceIp" = "8.8.8.8/32"
                }
            }
        },
    ]
  })
}
