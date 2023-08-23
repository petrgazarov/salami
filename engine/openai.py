import os
import openai
import json
from enum import Enum
from typing import List, Optional
from engine.types import ChatMessage, FunctionCall
from engine.logger import logger

openai.api_key = os.getenv("OPENAI_API_KEY")
openai.log = "debug"
temperature = 0.0


class ModelName(str, Enum):
    GPT3_5 = "gpt-3.5-turbo-0613"
    GPT4 = "gpt-4-0613"


async def create_chat_completion(
    messages: List[ChatMessage],
    functions=None,
    function_call=None,
    model_name: ModelName = ModelName.GPT4,
) -> dict:
    args = {
        "model": model_name.value,
        "messages": [message.to_dict() for message in messages],
        "temperature": temperature,
    }
    if function_call is not None:
        args["function_call"] = function_call
    if functions is not None:
        args["functions"] = functions
    response = await openai.ChatCompletion.acreate(**args)
    message = response["choices"][0]["message"]  # type: ignore
    print(f"OpenAI response: {message}")
    return message


def get_function_call_from_message(message: dict) -> Optional[FunctionCall]:
    function_call = message.get("function_call")
    if function_call is None:
        return None

    logger.info(f"OpenAI response function: {function_call}")
    return FunctionCall(
        name=function_call["name"],
        arguments=json.loads(function_call["arguments"]),
    )
