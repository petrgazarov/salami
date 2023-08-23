from engine.parse_source import parse_source
from engine.program.pulumi_resources import pulumi_resources
from engine.functions.get_pulumi_code import get_pulumi_code


async def run_program(source_files: list[str]):
    for file in source_files:
        with open(file, "r") as f:
            content = f.read()
        text_blocks = content.split("\n\n")
        for text_block in text_blocks:
            pulumi_resource = parse_source(text=text_block)
            pulumi_resource = await get_pulumi_code(pulumi_resource=pulumi_resource)
            pulumi_resources[pulumi_resource.logical_name] = pulumi_resource
