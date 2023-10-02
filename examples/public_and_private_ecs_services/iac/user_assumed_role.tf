resource "aws_s3_bucket" "AssumedRolesBucket" {
  bucket = "assumed-roles-hg12a"
  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket_public_access_block" "AssetsPublicAccessBlock" {
  bucket = aws_s3_bucket.AssumedRolesBucket.id

  block_public_acls       = true
  block_public_policy     = false
  ignore_public_acls      = true
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "AssumedRolesBucketPolicy" {
  bucket = aws_s3_bucket.AssumedRolesBucket.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::${aws_s3_bucket.AssumedRolesBucket.bucket}/*"
    }
  ]
}
POLICY
}