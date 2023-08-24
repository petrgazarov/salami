resource_type = "cloudtrail.Trail"
logical_name = "CumuliCloudtrail"
text = (
    "Name: cumuli-cloudtrail\nGlobal service events: True\n"
    "Multi-region: True\nLogging enabled: True"
)
pulumi_aws_imports = ["cloudtrail"]
other_packages = []
uses = ["CumuliCloudtrailLogsBucket"]
