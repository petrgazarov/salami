resource_type = "ec2.SecurityGroup"
logical_name = "PythonExecEcsSecurityGroup"
text = (
    "Name: python-exec-ecs-security-group\nDescription: Security group for python exec ECS service\n"
    "Egress: open to tcp traffic on port 443\nIngress: A list containing one rule allowing access on "
    '"tcp" protocol, {container_port} port, and limited to CumuliServerEcsSecurityGroup security group'
)
other_packages = []
properties = []
uses = ["CumuliVpc", "CumuliServerEcsSecurityGroup"]
exports = {}
