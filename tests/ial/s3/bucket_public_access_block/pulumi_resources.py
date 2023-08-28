from engine.models.pulumi_resource import LogicalName, ResourceType, PulumiResource

pulumi_resources = {
    LogicalName("CumuliAssumedRolesBucket"): PulumiResource(
        resource_type=ResourceType("s3.Bucket"),
        logical_name=LogicalName("CumuliAssumedRolesBucket"),
    )
}
