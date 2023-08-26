resource_type = "ec2.Subnet"
logical_name = "PublicSubnetA"
text = (
    "Name: public-subnet-a\nCidr block: 10.0.3.0/24\n"
    "Availability zone: us-west-1a\nMap public IP on launch: True"
)
pulumi_aws_imports = ["ec2"]
other_packages = []
properties = []
uses = ["CumuliVpc"]
exports = {}
