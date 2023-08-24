resource_type = "s3.BucketPolicy"
logical_name = "CumuliCloudtrailLogsBucketPolicy"
text = (
    'Policy: A JSON policy with two statements. The first statement allows the "cloudtrail.amazonaws.com" '
    "service to put objects to /AWSLogs/{aws_account_id}/* path inside the bucket, with the condition that the bucket owner has full control. "
    "The second statement allows the same service to get the bucket's ACL."
)
pulumi_aws_imports = ["s3"]
other_packages = ["json"]
uses = ["CumuliCloudtrailLogsBucket"]
