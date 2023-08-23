import pytest
from engine.program.run_program import run_program
from engine.program.pulumi_resources import pulumi_resources
from engine.pulumi_resource import LogicalName, ResourceType, PulumiResource
from tests.assertions import assert_equivalent_resource
from engine.utils.file import get_full_path, open_relative_file

source_filepath = get_full_path("source.txt")

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

expected_code = open_relative_file("pulumi_code.txt")


@pytest.mark.asyncio
async def test_s3_bucket_policy():
    await run_program(source_files=[source_filepath])
    pulumi_resource = pulumi_resources[LogicalName("CumuliCloudtrailLogsBucketPolicy")]
    assert pulumi_resource.resource_type == "s3.BucketPolicy"
    assert pulumi_resource.logical_name == "CumuliCloudtrailLogsBucketPolicy"
    assert pulumi_resource.text == (
        'Policy: A JSON policy with two statements. The first statement allows the "cloudtrail.amazonaws.com" '
        "service to put objects to /AWSLogs/{aws_account_id}/* path inside the bucket, with the condition that the bucket owner has full control. "
        "The second statement allows the same service to get the bucket's ACL."
    )
    assert pulumi_resource.pulumi_aws_imports == ["s3"]
    assert pulumi_resource.other_packages == ["json"]
    await assert_equivalent_resource(
        expected=expected_code,
        actual=pulumi_resource.code,
        pulumi_resource=pulumi_resource,
    )
