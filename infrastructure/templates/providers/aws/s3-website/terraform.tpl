terraform {
  backend "remote" {
    organization = "lftechnology"
    token = "{{ info.TERRAFORM_TOKEN }}"
    workspaces {
      name = "{{ info.CLIENT_WORKSPACE }}"
    }
  }
}

provider "aws" {
  region     = "{{ info.AWS_REGION }}"
  access_key = "{{ info.AWS_ACCESS_KEY }}"
  secret_key = "{{ info.AWS_SECRET_KEY }}"
}

resource "aws_s3_bucket" "bucket" {
  bucket = "{{ info.AWS_S3_BUCKET_NAME }}"
  acl    = "public-read"
  force_destroy = "true"

  website {
    index_document = "index.html"
    error_document = "index.html"
  }

  policy = <<POLICY
{
  "Version":"2012-10-17",
  "Statement":[
    {
      "Sid":"AddPerm",
      "Effect":"Allow",
      "Principal": "*",
      "Action":["s3:GetObject"],
      "Resource":["arn:aws:s3:::{{ info.AWS_S3_BUCKET_NAME }}/*"]
    }
  ]
}
POLICY
}

output "bucket_website_endpoint" {
  value = aws_s3_bucket.bucket.website_endpoint
}

output "bucket_name" {
  value = aws_s3_bucket.bucket.bucket
}
