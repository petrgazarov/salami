resource_type = "route53.Record"
logical_name = "FramerARecords"
text = (
    "Name: {domain_name}\nType: A\n"
    'Records: ["52.223.52.2", "35.71.142.77"]\nTtl: 1800'
)
pulumi_aws_imports = ["route53"]
other_packages = []
properties = []
uses = ["CumuliHostedZone"]
exports = {}
