terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

resource "aws_s3_bucket" "create_buckets" {
  for_each = toset(var.bucket_names)
  bucket = each.value
  tags = {
    Name = each.value
  }
}
