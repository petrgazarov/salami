resource "aws_s3_bucket" "TestBucket" {
  bucket = var.test_bucket_name
}

variable "test_bucket_name" {
  type = string
  default = "test-bucket-br31m10"
}