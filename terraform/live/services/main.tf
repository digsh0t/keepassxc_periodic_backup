provider "aws" {
  region = "us-east-2"
}

terraform {
  required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

module "s3_bucket" {
  source = "../../modules/services/s3"

  bucket_name = var.bucket_name
}
