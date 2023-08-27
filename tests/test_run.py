import os
import pytest
from engine.program.run_program import run_program
from tests.equivalent_resource import assert_equivalent_resource
import tests.ial.cloudtrail as cloudtrail
import tests.ial.s3 as s3
import tests.ial.route53 as route53
import tests.ial.ec2 as ec2
import tests.ial.ecs as ecs
from tests.mongodb import get_current_counter, create_motor_client, create_test_result


async def run_module(module):
    await run_program(
        source_files=[module.source_filepath],
        pulumi_resources=module.pulumi_resources,
    )
    return module.pulumi_resources[module.expected_attributes["logical_name"]]


modules = [
    cloudtrail.trail,
    s3.bucket_policy.public_bucket,
    s3.bucket_policy.cloudtrail_logs_bucket,
    s3.bucket_public_access_block,
    route53.record,
    ec2.subnet,
    ec2.route_table,
    ec2.route_table_association,
    ec2.vpc_endpoint.interface,
    ec2.vpc_endpoint.gateway,
    ec2.security_group,
    ecs.task_definition,
    ecs.service,
]


selected_modules = os.getenv("SELECTED_MODULES")
if selected_modules:
    selected_modules = selected_modules.split(",")
    selected_modules = [module.replace("/", ".") for module in selected_modules]
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

    test_result = await create_test_result(
        db=db,
        test_run_counter=current_counter,
        expected=module.expected_code,
        actual=pulumi_resource.code,
    )

    for name, value in module.expected_attributes.items():
        if isinstance(value, list):
            assert set(value) == set(
                getattr(pulumi_resource, name)
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
        test_result=test_result,
    )
