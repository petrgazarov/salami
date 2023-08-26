resource_type = "ec2.RouteTable"
logical_name = "PublicRouteTable"
text = (
    "Name: cumuli-public-route-table\n"
    'routes: A route with a CIDR block of "0.0.0.0/0" and CumuliInternetGateway internet gateway'
)
pulumi_aws_imports = ["ec2"]
other_packages = []
properties = []
uses = ["CumuliVpc", "CumuliInternetGateway"]
exports = {}
