from engine.pulumi_resource import LogicalName, ResourceType, PulumiResource

pulumi_resources = {
    LogicalName("PrivateSubnetA"): PulumiResource(
        resource_type=ResourceType("ec2.Subnet"),
        logical_name=LogicalName("PrivateSubnetA"),
    ),
    LogicalName("PrivateRouteTable"): PulumiResource(
        resource_type=ResourceType("ec2.RouteTable"),
        logical_name=LogicalName("PrivateRouteTable"),
    ),
}
