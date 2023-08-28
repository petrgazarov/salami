from pydantic import BaseModel, Field
from engine.types import LogicalName, Resource


class SymbolTable(BaseModel):
    variables: dict[str, str] = Field(
        {},
        description="Variables specified with @variable or @variables decorators.",
    )
    resources: dict[LogicalName, Resource] = Field(
        {},
        description="Resources in the program.",
    )
    locations: dict[LogicalName, str] = Field(
        {},
        description="File paths of resources relative to the root of the generated program.",
    )
