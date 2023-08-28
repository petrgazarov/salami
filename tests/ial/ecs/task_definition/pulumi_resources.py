from engine.models.pulumi_resource import LogicalName, ResourceType, PulumiResource

pulumi_resources = {
    LogicalName("CumuliServerTaskRole"): PulumiResource(
        resource_type=ResourceType("ecs.TaskRole"),
        logical_name=LogicalName("CumuliServerTaskRole"),
    ),
    LogicalName("CumuliEcsExecutionRole"): PulumiResource(
        resource_type=ResourceType("iam.Role"),
        logical_name=LogicalName("CumuliEcsExecutionRole"),
    ),
    LogicalName("CumuliServerRepository"): PulumiResource(
        resource_type=ResourceType("ecr.Repository"),
        logical_name=LogicalName("CumuliServerRepository"),
    ),
    LogicalName("CumuliServerLogGroup"): PulumiResource(
        resource_type=ResourceType("logs.LogGroup"),
        logical_name=LogicalName("CumuliServerLogGroup"),
    ),
}
