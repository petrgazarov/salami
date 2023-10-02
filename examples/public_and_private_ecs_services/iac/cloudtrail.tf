resource "aws_s3_bucket" "SalamiCloudtrailLogsBucket" {
  bucket = "salami-cloudtrail-logs"
}

resource "aws_s3_bucket_policy" "SalamiCloudtrailLogsBucketPolicy" {
  bucket = aws_s3_bucket.SalamiCloudtrailLogsBucket.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowPutObject",
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudtrail.amazonaws.com"
      },
      "Action": "s3:PutObject",
      "Resource": "arn:aws:s3:::${aws_s3_bucket.SalamiCloudtrailLogsBucket.bucket}/*",
      "Condition": {
        "StringEquals": {
          "s3:x-amz-acl": "bucket-owner-full-control"
        }
      }
    },
    {
      "Sid": "AllowGetBucketAcl",
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudtrail.amazonaws.com"
      },
      "Action": "s3:GetBucketAcl",
      "Resource": "arn:aws:s3:::${aws_s3_bucket.SalamiCloudtrailLogsBucket.bucket}"
    }
  ]
}
POLICY
}

resource "aws_cloudtrail" "Cloudtrail" {
  name                          = "cloudtrail"
  s3_bucket_name                = aws_s3_bucket.SalamiCloudtrailLogsBucket.bucket
  include_global_service_events = true
  is_multi_region_trail         = true
  enable_logging                = true
}