import os
import pytest
from engine.program.run_program import run_program
from tests.assertions import assert_equivalent_resource
import tests.ial.cloudtrail as cloudtrail
import tests.ial.servicediscovery as servicediscovery
import tests.ial.ecs as ecs
import tests.ial.s3 as s3
import tests.ial.route53 as route53
import tests.ial.ec2 as ec2
from tests.mongodb import get_current_counter, create_motor_client


async def run_module(module):
    await run_program(
        source_files=[module.source_filepath],
        pulumi_resources=module.pulumi_resources,
    )
    return module.pulumi_resources[module.expected_attributes["logical_name"]]


modules = [
    cloudtrail.trail,
    servicediscovery.private_dns_namespace,
    ecs.cluster,
    s3.bucket,
    s3.bucket_policy.public_bucket,
    s3.bucket_policy.cloudtrail_logs_bucket,
    s3.bucket_public_access_block,
    route53.hosted_zone,
    route53.record,
    ec2.vpc,
    ec2.subnet.public,
    ec2.subnet.private,
    ec2.route_table.public,
    ec2.route_table.private,
]


selected_modules = os.getenv("SELECTED_MODULES")
if selected_modules:
    selected_modules = selected_modules.split(",")
    modules = [
        module
        for module in modules
        if any(
            module.__name__.startswith(selected_module)
            for selected_module in selected_modules
        )
    ]

test_cases = [(run_module(module), module) for module in modules]

current_counter = None


@pytest.fixture(scope="function", autouse=True)
async def setup_database(event_loop):
    global db
    client = create_motor_client(event_loop)
    db = client.salami_test
    yield db
    client.close()


@pytest.mark.parametrize("result_coroutine,module", test_cases)
@pytest.mark.asyncio
async def test_module(setup_database, result_coroutine, module):
    db = await setup_database.__anext__()
    global current_counter
    if current_counter is None:
        current_counter = await get_current_counter(db)
    pulumi_resource = await result_coroutine
    for name, value in module.expected_attributes.items():
        if name == "pulumi_aws_imports":
            assert set(value).issubset(
                set(getattr(pulumi_resource, name))
            ), f"Failed on attribute: {name}"
        else:
            assert value == getattr(
                pulumi_resource, name
            ), f"Failed on attribute: {name}"
    await assert_equivalent_resource(
        db=db,
        expected=module.expected_code,
        pulumi_resource=pulumi_resource,
        pulumi_resources=module.pulumi_resources,
        current_counter=current_counter,
    )
