package frontend_architecture

const InfrastructureTemplate = `
// Terraform State Backend Initialization
terraform {
  backend "remote" {
    organization = "lftechnology"
    token = "{{ info.Token }}"
    workspaces {
      name = "{{ info.Client.Project }}-{{ info.Client.Deployment.Name }}-{{ info.Client.Deployment.Type }}"
    }
  }
}

variable "region" {
  default = "us-east-1"
}

variable "bucket_name" {
  default = "com.shift.{{ info.Client.Project|lower  }}.{{ info.Client.Deployment.Name|lower }}" 
}
// Provider Initialization
provider "aws" {
  region     = var.region
  access_key = "{{ info.Client.Deployment.AccessKey }}"
  secret_key = "{{ info.Client.Deployment.SecretKey }}"
}

//Bucket Initialization
resource "aws_s3_bucket" "bucket" {
  bucket = var.bucket_name
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
      "Resource":["arn:aws:s3:::${var.bucket_name}/*"]
    }
  ]
}
POLICY
}

resource "aws_cloudfront_distribution" "www_distribution" {
  // origin is where CloudFront gets its content from.
  origin {
    custom_origin_config {
      // These are all the defaults.
      http_port              = "80"
      https_port             = "443"
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1", "TLSv1.1", "TLSv1.2"]
    }

    // Here we're using our S3 bucket's URL!
    domain_name = aws_s3_bucket.bucket.website_endpoint
    // This can be any name to identify this origin.
    origin_id = var.bucket_name
  }

  enabled             = true
  default_root_object = "index.html"

  // All values are defaults from the AWS console.
  default_cache_behavior {
    viewer_protocol_policy = "redirect-to-https"
    compress               = true
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    // This needs to match the origin above.
    target_origin_id = var.bucket_name
    min_ttl          = 0
    default_ttl      = 0
    max_ttl          = 0

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  wait_for_deployment = false
  tags = {
    Name = var.bucket_name
    Project = "{{ info.Client.Project }}-{{ info.Client.Deployment.Name }}-{{ info.Client.Deployment.Type }}"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  // 404 Handling
  custom_error_response {
    error_code = 404
    response_code = 200
    response_page_path = "/"
  }
}

// Outputs

output "bucket_name" {
  value = aws_s3_bucket.bucket.bucket
}

output "frontend_web_url" {
  value = aws_cloudfront_distribution.www_distribution.domain_name
}
`
