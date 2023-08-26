resource_type = "ec2.Vpc"
logical_name = "CumuliVpc"
text = (
    "Name: cumuli-vpc\nCidr block: 10.0.0.0/16\n"
    "Enable DNS support: True\nEnable DNS hostnames: True"
)
pulumi_aws_imports = ["ec2"]
other_packages = []
properties = []
uses = []
exports = {}
