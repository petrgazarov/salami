from engine.pulumi_resource import LogicalName, ResourceType, PulumiResource

pulumi_resources = {
    LogicalName("CumuliHostedZone"): PulumiResource(
        resource_type=ResourceType("route53.Zone"),
        logical_name=LogicalName("CumuliHostedZone"),
    )
}
