from engine.models.pulumi_resource import LogicalName, ResourceType, PulumiResource

pulumi_resources = {
    LogicalName("CumuliEcsCluster"): PulumiResource(
        resource_type=ResourceType("ecs.Cluster"),
        logical_name=LogicalName("CumuliEcsCluster"),
    ),
    LogicalName("CumuliServerTaskDefinition"): PulumiResource(
        resource_type=ResourceType("iam.TaskDefinition"),
        logical_name=LogicalName("CumuliServerTaskDefinition"),
    ),
    LogicalName("PublicSubnetA"): PulumiResource(
        resource_type=ResourceType("ecr.Subnet"),
        logical_name=LogicalName("PublicSubnetA"),
    ),
    LogicalName("PublicSubnetB"): PulumiResource(
        resource_type=ResourceType("logs.Subnet"),
        logical_name=LogicalName("PublicSubnetB"),
    ),
    LogicalName("CumuliServerEcsSecurityGroup"): PulumiResource(
        resource_type=ResourceType("ec2.SecurityGroup"),
        logical_name=LogicalName("CumuliServerEcsSecurityGroup"),
    ),
    LogicalName("CumuliServerTargetGroup"): PulumiResource(
        resource_type=ResourceType("ec2.TargetGroup"),
        logical_name=LogicalName("CumuliServerTargetGroup"),
    ),
}
