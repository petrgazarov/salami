import re
from typing import ClassVar, Type
from pydantic import Field, BaseModel
from engine.openai import create_chat_completion, get_function_call_from_message
from engine.types import ChatMessage, MessageRole
from engine.models import PulumiResource
from engine.pydantic_to_openai import pydantic_to_openai
from engine.utils.file import open_relative_file


class GetPulumiCodeSchema(BaseModel):
    code: str = Field(
        ...,
        description=("Python code for a Pulumi resource."),
    )
    other_packages: list[str] = Field(
        ...,
        description=(
            "List of other packages referenced in the code. Do not include pulumi_aws. E.g. ['json', 'os']"
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
    result = "These variables are pulumi resources that this resource depends on. These are asynchronous in nature and are of type Output[T]:\n"
    for logical_name in pulumi_resource.uses:
        resource = pulumi_resources[logical_name]
        result += f"{pulumi_resource.variable_name()}: {resource.resource_type}\n"
    result += "\n"
    local_variables = re.findall(r"\{([^}]*)\}", pulumi_resource.text)
    if len(local_variables) > 0:
        result += "These variables are local. They can be accessed directly. No need to await them inside pulumi.Output.all(...):\n"
        for variable in local_variables:
            result += f"{variable}: local variable (string, use directly)\n"
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
            ChatMessage(role=MessageRole.SYSTEM, content=system_prompt),
            ChatMessage(role=MessageRole.USER, content=user_prompt),
        ],
        function_call={"name": GetPulumiCode.name},
        functions=[pydantic_to_openai(GetPulumiCode)],
    )
    function_call = get_function_call_from_message(chat_completion)
    if function_call:
        pulumi_resource.code = function_call.arguments["code"]
        pulumi_resource.other_packages = function_call.arguments["other_packages"]
    return pulumi_resource
