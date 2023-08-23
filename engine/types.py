from pydantic import BaseModel
import json
from pydantic import BaseModel, Field
from enum import Enum
from typing import Optional


class Role(str, Enum):
    USER = "user"
    SYSTEM = "system"
    ASSISTANT = "assistant"
    FUNCTION = "function"


class FunctionCall(BaseModel):
    name: str = Field(
        ...,
        description="The name of the function to call.",
    )
    arguments: dict = Field(
        ...,
        description="The arguments of the function.",
    )


class ChatMessage(BaseModel):
    role: Role = Field(
        ...,
        description="The role of the message.",
    )
    content: str = Field(
        ...,
        description="The content of the message.",
    )
    name: Optional[str] = None
    function_call: Optional[FunctionCall] = None

    def to_dict(self):
        result: dict[str, str | dict] = {
            "role": self.role.value,
            "content": self.content,
        }
        if self.name:
            result["name"] = self.name
        if self.function_call:
            result["function_call"] = {
                "name": self.function_call.name,
                "arguments": json.dumps(self.function_call.arguments),
            }
        return result
