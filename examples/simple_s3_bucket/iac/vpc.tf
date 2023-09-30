variable "test_bucket_name" {
  description = "The name of the S3 bucket"
  type        = string
}

resource "aws_s3_bucket" "TestBucket" {
  bucket = var.test_bucket_name
  acl    = "private"

  tags = {
    Name = var.test_bucket_name
  }
}

variable "test_bucket_name" {
  type = string
  default = "test-bucket-br31m9"
}