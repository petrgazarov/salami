import pytest
from engine.program.run_program import run_program
from engine.program.pulumi_resources import pulumi_resources
from engine.pulumi_resource import LogicalName, ResourceType, PulumiResource
from tests.assertions import assert_equivalent_resource
from engine.utils.file import get_full_path, open_relative_file
from tests.utils.expected_attributes import expected_attributes

source_filepath = get_full_path("source.txt")
expected_code = open_relative_file("expected_code.txt")

pulumi_resources[LogicalName("CumuliCloudtrailLogsBucket")] = PulumiResource(
    resource_type=ResourceType("s3.Bucket"),
    logical_name=LogicalName("CumuliCloudtrailLogsBucket"),
    text="Bucket name: cumuli-cloudtrail-logs",
    code="",
    exports={},
    pulumi_aws_imports=["s3"],
    other_packages=[],
    properties=[],
    uses=[],
)


@pytest.mark.asyncio
async def test_cloudtrail():
    await run_program(source_files=[source_filepath])
    pulumi_resource = pulumi_resources[LogicalName("CumuliCloudtrail")]
    for name, value in expected_attributes():
        assert value == getattr(pulumi_resource, name), f"Failed on attribute: {name}"
    await assert_equivalent_resource(
        expected=expected_code,
        pulumi_resource=pulumi_resource,
    )
