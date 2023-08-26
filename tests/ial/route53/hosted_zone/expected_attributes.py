resource_type = "route53.Zone"
logical_name = "CumuliHostedZone"
text = "Name: {domain_name}"
pulumi_aws_imports = ["route53"]
other_packages = []
properties = []
uses = []
exports = {"name_servers": "name-servers"}
