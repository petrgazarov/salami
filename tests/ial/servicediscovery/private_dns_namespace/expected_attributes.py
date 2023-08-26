resource_type = "servicediscovery.PrivateDnsNamespace"
logical_name = "CumuliEcsPrivateDnsNamespace"
text = "Description: Private namespace for Cumuli ECS cluster\nName: {local_dns_namespace_name}"
pulumi_aws_imports = ["servicediscovery"]
other_packages = []
uses = ["CumuliVpc"]
exports = {}
