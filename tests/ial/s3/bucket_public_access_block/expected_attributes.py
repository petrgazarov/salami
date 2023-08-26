resource_type = "s3.BucketPublicAccessBlock"
logical_name = "CumuliAssetsPublicAccessBlock"
text = (
    "Block public ACLs: True\nBlock public policy: False\n"
    "Ignore public ACLs: True\nRestrict public buckets: False"
)
pulumi_aws_imports = ["s3"]
other_packages = []
uses = ["CumuliAssumedRolesBucket"]
exports = {}
