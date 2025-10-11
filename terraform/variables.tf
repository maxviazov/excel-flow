variable "aws_region" {
  default = "us-east-1"
}

variable "domain_name" {
  description = "Domain name (e.g., viazov.dev)"
  type        = string
}

variable "route53_zone_id" {
  description = "Route53 hosted zone ID (optional, not used with Cloudflare)"
  type        = string
  default     = ""
}
