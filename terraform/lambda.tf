terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

provider "aws" {
  alias  = "us_east_1"
  region = "us-east-1"
}

# S3 Buckets
resource "aws_s3_bucket" "frontend" {
  bucket = "excel-viazov-dev"
}

resource "aws_s3_bucket" "files" {
  bucket = "excel-viazov-dev-files"
}

resource "aws_s3_bucket_website_configuration" "frontend" {
  bucket = aws_s3_bucket.frontend.id
  index_document { suffix = "index.html" }
  error_document { key = "index.html" }
}

resource "aws_s3_bucket_public_access_block" "frontend" {
  bucket = aws_s3_bucket.frontend.id
  block_public_acls = false
  block_public_policy = false
  ignore_public_acls = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "frontend" {
  bucket = aws_s3_bucket.frontend.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = "*"
      Action = "s3:GetObject"
      Resource = "${aws_s3_bucket.frontend.arn}/*"
    }]
  })
  depends_on = [aws_s3_bucket_public_access_block.frontend]
}

# Lambda
resource "aws_iam_role" "lambda_role" {
  name = "${var.project_name}-lambda-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy" "lambda_policy" {
  name = "${var.project_name}-lambda-policy"
  role = aws_iam_role.lambda_role.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = ["logs:*"]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow"
        Action = ["s3:*"]
        Resource = "${aws_s3_bucket.files.arn}/*"
      }
    ]
  })
}

resource "aws_lambda_function" "app" {
  filename = "lambda.zip"
  function_name = var.project_name
  role = aws_iam_role.lambda_role.arn
  handler = "main"
  runtime = "provided.al2023"
  timeout = 300
  memory_size = 512
  environment {
    variables = { S3_BUCKET = aws_s3_bucket.files.bucket }
  }
}

# API Gateway
resource "aws_api_gateway_rest_api" "api" {
  name = "${var.project_name}-api"
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id = aws_api_gateway_rest_api.api.root_resource_id
  path_part = "{proxy+}"
}

resource "aws_api_gateway_method" "proxy_method" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method.proxy_method.http_method
  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = aws_lambda_function.app.invoke_arn
}

resource "aws_api_gateway_deployment" "api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name = "prod"
  depends_on = [aws_api_gateway_integration.lambda_integration]
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id = "AllowExecutionFromAPIGateway"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.app.function_name
  principal = "apigateway.amazonaws.com"
  source_arn = "${aws_api_gateway_rest_api.api.execution_arn}/*/*"
}

# SSL Certificates
resource "aws_acm_certificate" "frontend" {
  provider = aws.us_east_1
  domain_name = "excel.viazov.dev"
  validation_method = "DNS"
  lifecycle { create_before_destroy = true }
}

# CloudFront
resource "aws_cloudfront_distribution" "frontend" {
  enabled = true
  is_ipv6_enabled = true
  default_root_object = "index.html"
  aliases = ["excel.viazov.dev"]

  origin {
    domain_name = aws_s3_bucket_website_configuration.frontend.website_endpoint
    origin_id = "S3-${aws_s3_bucket.frontend.id}"
    custom_origin_config {
      http_port = 80
      https_port = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols = ["TLSv1.2"]
    }
  }

  origin {
    domain_name = "${aws_api_gateway_rest_api.api.id}.execute-api.${var.aws_region}.amazonaws.com"
    origin_id = "API-Gateway"
    custom_origin_config {
      http_port = 80
      https_port = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols = ["TLSv1.2"]
    }
  }

  default_cache_behavior {
    allowed_methods = ["GET", "HEAD", "OPTIONS"]
    cached_methods = ["GET", "HEAD"]
    target_origin_id = "S3-${aws_s3_bucket.frontend.id}"
    compress = true
    viewer_protocol_policy = "redirect-to-https"
    forwarded_values {
      query_string = false
      cookies { forward = "none" }
    }
    min_ttl = 0
    default_ttl = 3600
    max_ttl = 86400
  }

  ordered_cache_behavior {
    path_pattern = "/api/*"
    allowed_methods = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods = ["GET", "HEAD"]
    target_origin_id = "API-Gateway"
    compress = true
    viewer_protocol_policy = "redirect-to-https"
    forwarded_values {
      query_string = true
      headers = ["Authorization", "Content-Type"]
      cookies { forward = "none" }
    }
    min_ttl = 0
    default_ttl = 0
    max_ttl = 0
  }

  restrictions {
    geo_restriction { restriction_type = "none" }
  }

  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate.frontend.arn
    ssl_support_method = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
}

# Outputs
output "cloudfront_domain" {
  value = aws_cloudfront_distribution.frontend.domain_name
}

output "certificate_validation" {
  value = {
    for dvo in aws_acm_certificate.frontend.domain_validation_options : dvo.domain_name => {
      name = dvo.resource_record_name
      type = dvo.resource_record_type
      value = dvo.resource_record_value
    }
  }
}