from engine.models.pulumi_resource import LogicalName, ResourceType, PulumiResource

pulumi_resources = {
    LogicalName("CumuliVpc"): PulumiResource(
        resource_type=ResourceType("ec2.Vpc"),
        logical_name=LogicalName("CumuliVpc"),
    ),
    LogicalName("PrivateSubnetA"): PulumiResource(
        resource_type=ResourceType("ec2.Subnet"),
        logical_name=LogicalName("PrivateSubnetA"),
    ),
    LogicalName("PrivateSubnetB"): PulumiResource(
        resource_type=ResourceType("ec2.Subnet"),
        logical_name=LogicalName("PrivateSubnetB"),
    ),
    LogicalName("EcrVpcEndpointSG"): PulumiResource(
        resource_type=ResourceType("ec2.SecurityGroup"),
        logical_name=LogicalName("EcrVpcEndpointSG"),
    ),
    LogicalName("CumuliEcsExecutionRole"): PulumiResource(
        resource_type=ResourceType("iam.Role"),
        logical_name=LogicalName("CumuliEcsExecutionRole"),
    ),
    LogicalName("PythonExecEcsExecutionRole"): PulumiResource(
        resource_type=ResourceType("iam.Role"),
        logical_name=LogicalName("PythonExecEcsExecutionRole"),
    ),
}
