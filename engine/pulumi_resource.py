from pydantic import BaseModel, ConfigDict


class ResourceType(str):
    pass


class LogicalName(str):
    pass


class PulumiResource(BaseModel):
    resource_type: ResourceType
    logical_name: LogicalName
    code: str = ""
    other_packages: list[str] = []
    text: str = ""
    properties: list[str] = []
    uses: list[LogicalName] = []
    exports: dict[str, str] = {}

    model_config = ConfigDict(arbitrary_types_allowed=True)
