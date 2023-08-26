from engine.pulumi_resource import LogicalName, ResourceType, PulumiResource

pulumi_resources = {
    LogicalName("CumuliCloudtrailLogsBucket"): PulumiResource(
        resource_type=ResourceType("s3.Bucket"),
        logical_name=LogicalName("CumuliCloudtrailLogsBucket"),
    )
}
