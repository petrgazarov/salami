from engine.models.pulumi_resource import LogicalName, ResourceType, PulumiResource

pulumi_resources = {
    LogicalName("CumuliVpc"): PulumiResource(
        resource_type=ResourceType("ec2.Vpc"),
        logical_name=LogicalName("CumuliVpc"),
    ),
    LogicalName("CumuliInternetGateway"): PulumiResource(
        resource_type=ResourceType("ec2.InternetGateway"),
        logical_name=LogicalName("CumuliInternetGateway"),
    ),
}
