import re
from typing import ClassVar, Type
from pydantic import Field, BaseModel
from engine.openai import create_chat_completion, get_function_call_from_message
from engine.types import ChatMessage, Role
from engine.pulumi_resource import PulumiResource
from engine.pydantic_to_openai import pydantic_to_openai
from engine.utils.file import open_relative_file


class GetPulumiCodeSchema(BaseModel):
    code: str = Field(
        ...,
        description=("Python code for a Pulumi resource."),
    )
    pulumi_aws_imports: list[str] = Field(
        ...,
        description=(
            "List of imports from pulumi_aws package. The imports must include only those explicitly referenced in the code. E.g. ['ec2']"
        ),
    )
    other_packages: list[str] = Field(
        ...,
        description=(
            "List of other packages referenced in the code. E.g. ['json', 'os']"
        ),
    )


class GetPulumiCode:
    name: ClassVar[str] = "get_pulumi_code"
    description: ClassVar[
        str
    ] = "Converts a description to python code for a Pulumi resource."
    input_schema: ClassVar[Type[GetPulumiCodeSchema]] = GetPulumiCodeSchema


def get_other_properties(pulumi_resource: PulumiResource):
    if not pulumi_resource.properties:
        return ""
    result = "Other resource properties:\n"
    for property in pulumi_resource.properties:
        result += f"{property}\n"
    return result


def get_text(pulumi_resource: PulumiResource):
    if not pulumi_resource.text:
        return ""
    result = "Resource description:\n"
    result += f"{pulumi_resource.text}\n"
    return result


def get_variables(
    pulumi_resource: PulumiResource,
    pulumi_resources: dict[str, PulumiResource],
):
    if not pulumi_resource.uses:
        return ""
    result = "These are variables this resource depends on. These are asynchronous in nature and are of type Output[T]. They behave much like promises and must be awaited inside pulumi.Output.all(...):\n\n"
    for logical_name in pulumi_resource.uses:
        variable_name = "".join(
            ["_" + i.lower() if i.isupper() else i for i in logical_name]
        ).lstrip("_")
        resource = pulumi_resources[logical_name]
        result += f"{variable_name}: {resource.resource_type}\n"
    local_variables = re.findall(r"\{([^}]*)\}", pulumi_resource.text)
    if len(local_variables) > 0:
        result += "The following are local variables of type string. They can be accessed directly. Do not await them inside pulumi.Output.all(...) method.\n"
        for variable in local_variables:
            result += f"{variable}: local variable\n"
    result += "\n"
    return result


async def get_pulumi_code(
    pulumi_resource: PulumiResource,
    pulumi_resources: dict[str, PulumiResource],
):
    user_prompt = f"""Resource type:
{pulumi_resource.resource_type}
Pulumi logical name:
{pulumi_resource.logical_name}
{get_text(pulumi_resource=pulumi_resource)}
{get_other_properties(pulumi_resource=pulumi_resource)}
{get_variables(pulumi_resource=pulumi_resource, pulumi_resources=pulumi_resources)}
"""
    system_prompt = open_relative_file("system_prompt.txt")
    chat_completion = await create_chat_completion(
        messages=[
            ChatMessage(role=Role.SYSTEM, content=system_prompt),
            ChatMessage(role=Role.USER, content=user_prompt),
        ],
        function_call={"name": GetPulumiCode.name},
        functions=[pydantic_to_openai(GetPulumiCode)],
    )
    function_call = get_function_call_from_message(chat_completion)
    if function_call:
        pulumi_resource.code = function_call.arguments["code"]
        pulumi_resource.pulumi_aws_imports = function_call.arguments[
            "pulumi_aws_imports"
        ]
        pulumi_resource.other_packages = function_call.arguments["other_packages"]
    return pulumi_resource
