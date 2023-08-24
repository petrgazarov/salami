from engine.types import ChatMessage, Role
from engine.openai import create_chat_completion
from engine.utils.file import open_relative_file
from engine.pulumi_resource import PulumiResource
from engine.functions.get_pulumi_code.function import get_variables


def get_user_prompt(expected, actual):
    return f"Expected:\n{expected}\nActual:\n{actual}"


async def assert_equivalent_resource(expected: str, pulumi_resource: PulumiResource):
    system_prompt = open_relative_file("system_prompt.txt")
    actual = pulumi_resource.code
    res = await create_chat_completion(
        messages=[
            ChatMessage(role=Role.SYSTEM, content=system_prompt),
            ChatMessage(
                role=Role.SYSTEM,
                content=get_variables(pulumi_resource=pulumi_resource),
            ),
            ChatMessage(
                role=Role.USER,
                content=get_user_prompt(expected=expected, actual=actual),
            ),
        ],
    )
    content = res["content"]
    answer_line = content.strip().split("\n")[-1]
    answer = answer_line.split("FINAL ANSWER: ")[-1].strip()

    if answer != "equivalent":
        raise AssertionError(
            f"Equivalence check failed. Full trace: {content}\nExpected:\n{expected}\nActual:\n{actual}"
        )
