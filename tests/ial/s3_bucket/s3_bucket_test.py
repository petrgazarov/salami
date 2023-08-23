import pytest
from engine.program.run_program import run_program
from engine.program.pulumi_resources import pulumi_resources
from engine.pulumi_resource import LogicalName
from tests.assertions import assert_equivalent_resource
from engine.utils.file import get_full_path

source_filepath = get_full_path("source.txt")


@pytest.mark.asyncio
async def test_s3_bucket():
    await run_program(source_files=[source_filepath])
    pulumi_resource = pulumi_resources[LogicalName("CumuliCloudtrailLogsBucket")]
    assert pulumi_resource.resource_type == "s3.Bucket"
    assert pulumi_resource.logical_name == "CumuliCloudtrailLogsBucket"
    assert pulumi_resource.text == "Bucket name: cumuli-cloudtrail-logs"
    assert pulumi_resource.pulumi_aws_imports == ["s3"]
    assert pulumi_resource.other_packages == []
    await assert_equivalent_resource(
        expected='s3.Bucket("CumuliCloudtrailLogsBucket", bucket="cumuli-cloudtrail-logs")',
        actual=pulumi_resource.code,
        pulumi_resource=pulumi_resource,
    )
    assert pulumi_resource.text == "Bucket name: cumuli-cloudtrail-logs"
    assert pulumi_resource.properties == []
    assert pulumi_resource.uses == []
