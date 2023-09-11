resource_type = "s3.BucketPolicy"
logical_name = "CumuliAssumedRolesBucketPolicy"
text = (
    'Policy: A JSON policy that allows all principals to perform the "s3:GetObject" '
    "action on all objects in the specified S3 bucket."
)
other_packages = ["json"]
uses = ["CumuliAssumedRolesBucket"]
exports = {}