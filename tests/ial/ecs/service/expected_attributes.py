resource_type = "ecs.Service"
logical_name = "CumuliServerEcsService"
text = (
    "Name: cumuli-server\nDesired_count: 1\nLaunch type: FARGATE\n"
    "Network configuration: Assign public IP enabled. The subnets are PublicSubnetA and PublicSubnetB. "
    "Security group is CumuliServerEcsSecurityGroup.\nLoad balancers: A single target group "
    "CumuliServerTargetGroup, {server_container_name} container name, and {container_port} container_port.\n"
    "Wait for steady state: True"
)
other_packages = []
properties = [
    "ECS type deployment controller",
    "Enabled deployment circuit breaker with rollback",
]
uses = [
    "CumuliEcsCluster",
    "CumuliServerTaskDefinition",
    "PublicSubnetA",
    "PublicSubnetB",
    "CumuliServerEcsSecurityGroup",
    "CumuliServerTargetGroup",
]
exports = {}
